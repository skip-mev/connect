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
HOMEDIR ?= $(CURDIR)/tests/.slinkyd
GENESIS ?= $(HOMEDIR)/config/genesis.json
GENESIS_TMP ?= $(HOMEDIR)/config/genesis_tmp.json
APP_TOML ?= $(HOMEDIR)/config/app.toml
CONFIG_TOML ?= $(HOMEDIR)/config/config.toml
COVER_FILE ?= cover.out
BENCHMARK_ITERS ?= 10
USE_CORE_MARKETS ?= true
USE_RAYDIUM_MARKETS ?= false
USE_UNISWAPV3_BASE_MARKETS ?= false
USE_COINGECKO_MARKETS ?= false
USE_COINMARKETCAP_MARKETS ?= false
USE_OSMOSIS_MARKETS ?= false
SCRIPT_DIR := $(CURDIR)/scripts
DEV_COMPOSE ?= $(CURDIR)/contrib/compose/docker-compose-dev.yml

LEVANT_VAR_FILE:=$(shell mktemp -d)/levant.yaml
NOMAD_FILE_SLINKY:=contrib/nomad/slinky.nomad

TAG := $(shell git describe --tags --always --dirty)

export HOMEDIR := $(HOMEDIR)
export APP_TOML := $(APP_TOML)
export GENESIS := $(GENESIS)
export GENESIS_TMP := $(GENESIS_TMP)
export USE_CORE_MARKETS ?= $(USE_CORE_MARKETS)
export USE_RAYDIUM_MARKETS ?= $(USE_RAYDIUM_MARKETS)
export USE_UNISWAPV3_BASE_MARKETS ?= $(USE_UNISWAPV3_BASE_MARKETS)
export USE_COINGECKO_MARKETS ?= $(USE_COINGECKO_MARKETS)
export USE_COINMARKETCAP_MARKETS ?= $(USE_COINMARKETCAP_MARKETS)
export USE_OSMOSIS_MARKETS ?= $(USE_OSMOSIS_MARKETS)
export SCRIPT_DIR := $(SCRIPT_DIR)

BUILD_TAGS := -X github.com/skip-mev/slinky/cmd/build.Build=$(TAG)

###############################################################################
###                               build                                     ###
###############################################################################

build: tidy
	go build -ldflags="$(BUILD_TAGS)" \
	 -o ./build/ ./...

run-oracle-client: build
	@./build/client --host localhost --port 8080

start-all-dev:
	@echo "Starting development oracle side-car, blockchain, grafana, and prometheus dashboard..."
	@$(DOCKER_COMPOSE) -f $(DEV_COMPOSE) --profile all up -d --build

stop-all-dev:
	@echo "Stopping development network..."
	@$(DOCKER_COMPOSE) -f $(DEV_COMPOSE) --profile all down

start-sidecar-dev:
	@echo "Starting development oracle side-car, grafana, and prometheus dashboard..."
	@$(DOCKER_COMPOSE) -f $(DEV_COMPOSE) --profile sidecar up -d --build

stop-sidecar-dev:
	@echo "Stopping development oracle..."
	@$(DOCKER_COMPOSE) -f $(DEV_COMPOSE) --profile sidecar down

install: tidy
	@go install -ldflags="$(BUILD_TAGS)" -mod=readonly ./cmd/slinky
	@go install -mod=readonly $(BUILD_FLAGS) ./tests/simapp/slinkyd

.PHONY: build install run-oracle-client start-all-dev stop-all-dev

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

delete-configs:
	@rm -rf ./tests/.slinkyd/

build-market-map:
	@echo "Building market map..."
	@sh ./scripts/genesis.sh

# build-configs builds a slinky simulation application binary in the build folder (/test/.slinkyd)
build-configs: delete-configs build-market-map
	@dasel put -r toml 'instrumentation.enabled' -f $(CONFIG_TOML) -t bool -v true
	@dasel put -r toml 'rpc.laddr' -f $(CONFIG_TOML) -t string -v "tcp://0.0.0.0:26657"
	@dasel put -r toml 'telemetry.enabled' -f $(APP_TOML) -t bool -v true
	@dasel put -r toml 'api.enable' -f $(APP_TOML) -t bool -v true
	@dasel put -r toml 'grpc.address' -f $(APP_TOML) -t string -v "0.0.0.0:9090"
	@dasel put -r toml 'api.address' -f $(APP_TOML) -t string -v "tcp://0.0.0.0:1317"
	@dasel put -r toml 'api.enabled-unsafe-cors' -f $(APP_TOML) -t bool -v true

# start-app starts a slinky simulation application binary in the build folder (/test/.slinkyd)
# this will set the environment variable for running locally
start-app:
	@./build/slinkyd start --log_level info --home $(HOMEDIR)

# build-and-start-app builds a slinky simulation application binary in the build folder
# and initializes a single validator configuration. If desired, users can supplement
# other addresses using "genesis add-genesis-account address 10000000000000000000000000stake".
# This will allow users to bootstrap their wallet with a balance.
build-and-start-app: build-configs start-app

.PHONY: build-test-app build-configs build-and-start-app start-app delete-configs

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
	@cd ./tests/petri && go mod tidy

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
	@cd ./providers/apis/defi/osmosis && go generate ./...

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

