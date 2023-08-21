#!/usr/bin/make -f

export VERSION := $(shell echo $(shell git describe --always --match "v*") | sed 's/^v//')
export COMMIT := $(shell git log -1 --format='%H')
export COMETBFT_VERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::')

BIN_DIR ?= $(GOPATH)/bin
BUILD_DIR ?= $(CURDIR)/build
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)
HTTPS_GIT := https://github.com/skip-mev/slinky.git
DOCKER := $(shell which docker)
CONFIG_FILE ?= $(CURDIR)/conf/dev/config.toml
HOMEDIR ?= $(CURDIR)/tests/.slinkyd
GENESIS ?= $(HOMEDIR)/config/genesis.json
GENESIS_TMP ?= $(HOMEDIR)/config/genesis_tmp.json

###############################################################################
###                               build                                     ###
###############################################################################

build:
	go build -o ./build/ ./...

run-oracle-server: build
	./build/oracle -config ${CONFIG_FILE}

.PHONY: build run-oracle-server

###############################################################################
##                                  Docker                                   ##
###############################################################################

docker-build:
	@echo "Building E2E Docker image..."
	@DOCKER_BUILDKIT=1 docker build -t skip-mev/slinky-e2e -f contrib/images/slinky.e2e.Dockerfile .
	@DOCKER_BUILDKIT=1 docker build -t skip-mev/slinky-e2e-oracle -f contrib/images/slinky.e2e.oracle.Dockerfile .

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
	cd $(CURDIR)/tests/simapp && go build -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILD_DIR)/:
	mkdir -p $(BUILD_DIR)/

# build-configs builds a slinky simulation application binary in the build folder (/test/.slinkyd)
build-configs: build-test-app
	rm -rf $(HOMEDIR)

	./build/slinkyd init validator --chain-id skip-1 --home $(HOMEDIR)
	./build/slinkyd keys add validator --home $(HOMEDIR)
	./build/slinkyd genesis add-genesis-account validator 10000000000000000000000000stake --home $(HOMEDIR)
	./build/slinkyd genesis add-genesis-account cosmos1see0htr47uapjvcvh0hu6385rp8lw3em24hysg 10000000000000000000000000stake --home $(HOMEDIR)
	./build/slinkyd genesis gentx validator 1000000000stake --chain-id skip-1 --home $(HOMEDIR)
	./build/slinkyd genesis collect-gentxs --home $(HOMEDIR)
	jq '.consensus["params"]["abci"]["vote_extensions_enable_height"] = "2"' $(GENESIS) > $(GENESIS_TMP) && mv $(GENESIS_TMP) $(GENESIS) 
	jq '.app_state["oracle"]["currency_pair_genesis"] += [{"currency_pair": {"Base": "BITCOIN", "Quote": "USD"},"currency_pair_price": null,"nonce": "0"}]' $(GENESIS) > $(GENESIS_TMP) && mv $(GENESIS_TMP) $(GENESIS)

# start-app starts a slinky simulation application binary in the build folder (/test/.slinkyd)
start-app:
	./build/slinkyd start --api.enable true --api.enabled-unsafe-cors true --log_level debug --home $(HOMEDIR)


# build-and-start-app builds a slinky simulation application binary in the build folder
# and initializes a single validator configuration. If desired, users can suppliment
# other addresses using "genesis add-genesis-account address 10000000000000000000000000stake".
# This will allow users to bootstrap their wallet with a balance.
build-and-start-app: build-test-app build-configs start-app

.PHONY: build-test-app build-configs build-and-start-app start-app

###############################################################################
###                               Testing                                   ###
###############################################################################

test-integration: docker-build
	@echo "Running integration tests..."
	@cd ./tests/integration && go mod tidy &&  go test -p 1 -v -race ./...

test: tidy
	@go test -v -race $(shell go list ./... | grep -v tests/)

.PHONY: test test-e2e

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.13.5
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-pulsar-gen:
	@echo "Generating Dep-Inj Protobuf files"
	@$(protoImage) sh ./scripts/protocgen-pulsar.sh

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

proto-check-breaking:
	@$(protoImage) buf breaking --against $(HTTPS_GIT)#branch=main

proto-update-deps:
	@echo "Updating Protobuf dependencies"
	$(DOCKER) run --rm -v $(CURDIR)/proto:/workspace --workdir /workspace $(protoImageName) buf mod update

.PHONY: proto-all proto-gen proto-pulsar-gen proto-format proto-lint proto-check-breaking proto-update-deps


###############################################################################
###                              Formatting                                 ###
###############################################################################

tidy:
	go mod tidy

.PHONY: tidy

###############################################################################
###                                Linting                                  ###
###############################################################################

golangci_lint_cmd=golangci-lint
golangci_version=v1.53.3

lint:
	@echo "--> Running linters"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@golangci-lint run

format:
	@echo "--> Running formatters"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@golangci-lint run --fix

lint-markdown:
	@echo "--> Running markdown linter"
	@markdownlint **/*.md

.PHONY: lint lint-fix lint-markdown

