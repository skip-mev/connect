#!/usr/bin/make -f

export VERSION := $(shell echo $(shell git describe --always --match "v*") | sed 's/^v//')
export COMMIT := $(shell git log -1 --format='%H')
export COMETBFT_VERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::')

BIN_DIR ?= $(GOPATH)/bin
BUILD_DIR ?= $(CURDIR)/build
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)
HTTPS_GIT := https://github.com/skip-mev/slinky.git
DOCKER := $(shell which docker)
DOCKER_COMPOSE := $(shell which docker-compose)
ORACLE_CONFIG_FILE ?= $(CURDIR)/config/local/oracle.json
MARKET_CONFIG_FILE ?= $(CURDIR)/config/local/market.json
CONFIG_DIR ?= $(CURDIR)/config
HOMEDIR ?= $(CURDIR)/tests/.slinkyd
GENESIS ?= $(HOMEDIR)/config/genesis.json
GENESIS_TMP ?= $(HOMEDIR)/config/genesis_tmp.json
COVER_FILE ?= cover.out
BENCHMARK_ITERS ?= 10
DEFI_PROVIDERS_ENABLED ?= false
SOLANA_NODE_ENDPOINT ?= https://api.devnet.solana.com

LEVANT_VAR_FILE:=$(shell mktemp -d)/levant.yaml
NOMAD_FILE_SLINKY:=contrib/nomad/slinky.nomad

###############################################################################
###                               build                                     ###
###############################################################################

build: tidy
	@go build -o ./build/ ./...

run-oracle-server: build
	@./build/slinky --oracle-config-path ${ORACLE_CONFIG_FILE} --market-config-path ${MARKET_CONFIG_FILE}

run-oracle-client: build
	@./build/client --host localhost --port 8080

run-prom-client: 
	@$(DOCKER) run \
		-p 9090:9090 \
		-v ./contrib/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml \
		prom/prometheus

update-local-configs: build
	@echo "Updating local config..."
	@./build/slinky-config --oracle-config-path ${ORACLE_CONFIG_FILE} --market-config-path ${MARKET_CONFIG_FILE} --raydium-enabled ${DEFI_ORACLE_ENABLED} \
		--solana-node-endpoint ${SOLANA_NODE_ENDPOINT}

generate-defi-configs:
	@echo "Generating defi configs..."
	@go run ./scripts/raydium/generate-config/main.go

start-oracle:
	@echo "Starting oracle side-car, blockchain, and prometheus dashboard..."
	@$(DOCKER_COMPOSE) -f docker-compose.yml up -d

stop-oracle:
	@echo "Stopping network..."
	@$(DOCKER_COMPOSE) -f docker-compose.yml down

start-sidecar:
	@echo "Starting oracle side-car and prometheus dashboard..."
	@$(DOCKER_COMPOSE) -f docker-compose.yml up -d oracle prometheus

stop-sidecar:
	@echo "Stopping network..."
	@$(DOCKER_COMPOSE) -f docker-compose.yml down

install: tidy
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/slinky
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/slinky-config

.PHONY: build run-oracle-server install

###############################################################################
##                                  Docker                                   ##
###############################################################################

docker-build:
	@echo "Building E2E Docker image..."
	@DOCKER_BUILDKIT=1 $(DOCKER) build -t skip-mev/slinky-e2e -f contrib/images/slinky.e2e.Dockerfile .
	@DOCKER_BUILDKIT=1 $(DOCKER) build -t skip-mev/slinky-e2e-oracle -f contrib/images/slinky.sidecar.dev.Dockerfile .

.PHONY: docker-build

###############################################################################
###                                Test App                                 ###
###############################################################################

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=testapp \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=testappd \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X github.com/cometbft/cometbft/version.TMCoreSemVer=$(COMETBFT_VERSION)

# DB backend selection
ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc
endif
ifeq (badgerdb,$(findstring badgerdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += badgerdb
endif
# handle rocksdb
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  CGO_ENABLED=1
  build_tags += rocksdb
endif
# handle boltdb
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += boltdb
endif

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

BUILD_TARGETS := build-test-app

build-test-app: BUILD_ARGS=-o $(BUILD_DIR)/

$(BUILD_TARGETS): $(BUILD_DIR)/
	@cd $(CURDIR)/tests/simapp && go build -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILD_DIR)/:
	@mkdir -p $(BUILD_DIR)/

# build-configs builds a slinky simulation application binary in the build folder (/test/.slinkyd)
build-configs:
	@./build/slinkyd init validator --chain-id skip-1 --home $(HOMEDIR)
	@./build/slinkyd keys add validator --home $(HOMEDIR) --keyring-backend test
	@./build/slinkyd genesis add-genesis-account validator 10000000000000000000000000stake --home $(HOMEDIR) --keyring-backend test
	@./build/slinkyd genesis add-genesis-account cosmos1see0htr47uapjvcvh0hu6385rp8lw3em24hysg 10000000000000000000000000stake --home $(HOMEDIR) --keyring-backend test
	@./build/slinkyd genesis gentx validator 1000000000stake --chain-id skip-1 --home $(HOMEDIR) --keyring-backend test
	@./build/slinkyd genesis collect-gentxs --home $(HOMEDIR)
	@jq '.consensus["params"]["abci"]["vote_extensions_enable_height"] = "2"' $(GENESIS) > $(GENESIS_TMP) && mv $(GENESIS_TMP) $(GENESIS)
	@jq '.app_state["oracle"]["currency_pair_genesis"] += [{"currency_pair": {"Base": "BTC", "Quote": "USD"},"currency_pair_price": null,"nonce": "0"}]' $(GENESIS) > $(GENESIS_TMP) && mv $(GENESIS_TMP) $(GENESIS)
	@jq '.app_state["oracle"]["next_id"] = "2"' $(GENESIS) > $(GENESIS_TMP) && mv $(GENESIS_TMP) $(GENESIS)

# start-app starts a slinky simulation application binary in the build folder (/test/.slinkyd)
# this will set the environment variable for running locally
start-app:
	@./build/slinkyd start --api.enable true --api.enabled-unsafe-cors true --log_level info --home $(HOMEDIR)


# build-and-start-app builds a slinky simulation application binary in the build folder
# and initializes a single validator configuration. If desired, users can supplement
# other addresses using "genesis add-genesis-account address 10000000000000000000000000stake".
# This will allow users to bootstrap their wallet with a balance.
build-and-start-app: build-configs start-app

.PHONY: build-test-app build-configs build-and-start-app start-app

###############################################################################
###                               Testing                                   ###
###############################################################################

test-integration: tidy docker-build
	@echo "Running integration tests..."
	@cd ./tests/integration &&  go test -p 1 -v -race -timeout 30m

test-petri-integ: tidy docker-build
	@echo "Running petri integration tests..."
	@cd ./tests/petri &&  go test -p 1 -v -race -timeout 30m

test: tidy
	@go test -v -race $(shell go list ./... | grep -v tests/)

test-bench: tidy
	@go test -count=$(BENCHMARK_ITERS) -benchmem -run notest -bench . ./... | grep Benchmark

test-cover: tidy
	@echo Running unit tests and creating coverage report...
	@go test -mod=readonly -v -timeout 30m -coverprofile=$(COVER_FILE) -covermode=atomic $(shell go list ./... | grep -v tests/)
	@sed -i'.bak' -e '/.pb.go/d' $(COVER_FILE)
	@sed -i'.bak' -e '/.pulsar.go/d' $(COVER_FILE)
	@sed -i'.bak' -e '/.proto/d' $(COVER_FILE)
	@sed -i'.bak' -e '/.pb.gw.go/d' $(COVER_FILE)
	@sed -i'.bak' -e '/mocks/d' $(COVER_FILE)

.PHONY: test test-e2e test-petri-integ

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.14.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: tidy proto-format proto-gen proto-pulsar-gen format

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-pulsar-gen:
	@echo "Generating Dep-Inj Protobuf files"
	@$(protoImage) sh ./scripts/protocgen-pulsar.sh

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(DOCKER) run --rm -v $(CURDIR)/proto:/workspace --workdir /workspace $(protoImageName) buf lint --error-format=json

proto-check-breaking:
	@$(protoImage) buf breaking --against $(HTTPS_GIT)#branch=main

proto-update-deps:
	@echo "Updating Protobuf dependencies"
	@$(DOCKER) run --rm -v $(CURDIR)/proto:/workspace --workdir /workspace $(protoImageName) buf mod update

.PHONY: proto-all proto-gen proto-pulsar-gen proto-format proto-lint proto-check-breaking proto-update-deps


###############################################################################
###                              Formatting                                 ###
###############################################################################

tidy:
	@go mod tidy
	@cd ./tests/integration && go mod tidy

.PHONY: tidy

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	@echo "--> Running linter"
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run --out-format=tab

lint-fix:
	@echo "--> Running linter"
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run --fix --out-format=tab --issues-exit-code=0

lint-markdown:
	@echo "--> Running markdown linter"
	@markdownlint **/*.md

govulncheck:
	@echo "--> Running govulncheck"
	@go run golang.org/x/vuln/cmd/govulncheck -test ./...

.PHONY: lint lint-fix lint-markdown govulncheck

###############################################################################
###                                Mocks                                    ###
###############################################################################

mocks: gen-mocks format

gen-mocks:
	@echo "--> generating mocks"
	@go install github.com/vektra/mockery/v2
	@go generate ./...

###############################################################################
###                                Formatting                               ###
###############################################################################

format:
	@find . -name '*.go' -type f -not -path "*.git*" -not -path "*/mocks/*" -not -name '*.pb.go' -not -name '*.pulsar.go' -not -name '*.gw.go' | xargs go run mvdan.cc/gofumpt -w .
	@find . -name '*.go' -type f -not -path "*.git*" -not -path "*/mocks/*" -not -name '*.pb.go' -not -name '*.pulsar.go' -not -name '*.gw.go' | xargs go run github.com/client9/misspell/cmd/misspell -w
	@find . -name '*.go' -type f -not -path "*.git*" -not -path "/*mocks/*" -not -name '*.pb.go' -not -name '*.pulsar.go' -not -name '*.gw.go' | xargs go run golang.org/x/tools/cmd/goimports -w -local github.com/skip-mev/slinky

.PHONY: format

###############################################################################
###                                dev-deploy                               ###
###############################################################################

deploy-dev:
	@touch ${LEVANT_VAR_FILE}
	@yq e -i '.sidecar_image |= "${SIDECAR_IMAGE}"' ${LEVANT_VAR_FILE}
	@yq e -i '.chain_image |= "${CHAIN_IMAGE}"' ${LEVANT_VAR_FILE}
	@levant deploy -force -force-count -var-file=${LEVANT_VAR_FILE} ${NOMAD_FILE_SLINKY}

.PHONY: deploy-dev

