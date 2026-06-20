BUILD_DIR ?= build
BINARY ?= heyad

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "unknown")
COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

LD_FLAGS = -X github.com/cosmos/cosmos-sdk/version.Name=heya \
	-X github.com/cosmos/cosmos-sdk/version.AppName=heyad \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

build:
	go build -trimpath -ldflags "-s -w $(LD_FLAGS)" -o $(BUILD_DIR)/$(BINARY) ./cmd/heyad

install:
	go install -trimpath -ldflags "-s -w $(LD_FLAGS)" ./cmd/heyad
