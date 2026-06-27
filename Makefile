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

tools:
	cd tools && GONOSUMCHECK=* GONOSUMDB=* GOPROXY=https://goproxy.io,direct \
		go install \
			github.com/bufbuild/buf/cmd/buf \
			golang.org/x/tools/cmd/goimports \
			google.golang.org/grpc/cmd/protoc-gen-go-grpc \
			github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

deps-update:
	go get github.com/hashicorp/go-getter@v1.8.6
	go get github.com/ulikunitz/xz@v0.5.15
	go get github.com/dvsekhvalnov/jose2go@v1.7.0

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(GOPATH)/bin/$(BINARY)
	go clean -cache
