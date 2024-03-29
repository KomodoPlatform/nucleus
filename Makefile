#!/usr/bin/make -f

VERSION := 0.0.1
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true

export GO111MODULE = on

# Build tags
build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# Linker flags
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=nucleus \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=nucleusd \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

clean:
	rm -rf build

build: go.sum
ifeq ($(OS),Windows_NT)
	go build $(BUILD_FLAGS) -o build/nucleusd.exe ./cmd/nucleusd
else
	go build $(BUILD_FLAGS) -o build/nucleusd ./cmd/nucleusd
endif

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/nucleusd

# Make sure following packages are installed in the system
# - protobuf
# - protobuf-devel  (or dev)
# - libprotoc-devel (or dev)

# Install `protoc-gen-gocosmos` and `protoc-gen-grpc-gateway` from source
# go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos
# go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0
gen-proto:
	proto_dirs=$$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
	for dir in $$proto_dirs; do
	  protoc \
	  -I "proto" \
	  -I "third_party/proto" \
	  --gocosmos_out=plugins=interfacetype+grpc,\
	Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
	  $$(find "$${dir}" -maxdepth 1 -name '*.proto')

	  # command to generate gRPC gateway (*.pb.gw.go in respective modules) files
	  protoc \
	  -I "proto" \
	  -I "third_party/proto" \
	  --grpc-gateway_out=logtostderr=true:. \
	  $$(find "$${dir}" -maxdepth 1 -name '*.proto')
	done
	cp -rf ./nucleus/x/* ./x
	rm -rf ./nucleus

fmt-check:
	find . -name '*.go' -type f -not -path "*.git*" | xargs gofmt -d -s
	find . -name '*.go' -type f -not -path "*.git*" | xargs goimports -d -e

lint-check:
	golangci-lint run

check:
	$(MAKE) fmt-check
	$(MAKE) lint-check
	go mod verify

benchmark:
	@go test -mod=readonly -bench=. ./...

test: test-unit
test-all: test-unit test-race test-cover

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./...

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -tags='ledger test_ledger_mock' ./...

test-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

cpu-profile-simulation-test:
	@VERSION=$(VERSION) go test -benchmem -run ^BenchmarkSimulation -bench ^BenchmarkSimulation ./app -cpuprofile cpu.out -Commit=true -Verbose=true -Enabled=true

.PHONY: build install test check
.ONESHELL:
.SILENT:
