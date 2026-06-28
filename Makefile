BUILD_DIR ?= build
BINARY ?= heyad

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "unknown")
COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

LD_FLAGS = -X github.com/cosmos/cosmos-sdk/version.Name=heya \
	-X github.com/cosmos/cosmos-sdk/version.AppName=heyad \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

# release: stripped, PIE, no debug info
RELEASE_FLAGS = -trimpath -buildmode=pie -ldflags "-s -w $(LD_FLAGS)"

# debug: with symbols, optimizations off, inlining off
DEBUG_FLAGS = -gcflags "all=-N -l" -ldflags "$(LD_FLAGS)"

# PGO: Go 1.21+ auto-picks cmd/heyad/default.pgo if present (~2-7% CPU boost)
build: cmd/heyad/default.pgo
	go build -mod=readonly $(RELEASE_FLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/heyad

install: cmd/heyad/default.pgo
	go install -mod=readonly $(RELEASE_FLAGS) ./cmd/heyad

build-debug:
	go build -mod=mod $(DEBUG_FLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/heyad

install-debug:
	go install -mod=mod $(DEBUG_FLAGS) ./cmd/heyad

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

# regenerate PGO profile (run after protocol changes)
pgo-profile:
	go test -benchmem -run='^$$' -bench '^BenchmarkFullAppSimulation$$' ./app \
		-Commit=true -cpuprofile=cmd/heyad/default.pgo \
		-Enabled=true -Period=5 -NumBlocks=50

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(GOPATH)/bin/$(BINARY)
	go clean -cache
