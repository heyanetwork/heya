#!/bin/bash
set -e

if [ -z "$1" ]; then
    echo "Usage: $0 <tag>"
    echo "Example: $0 v1.1.0"
    exit 1
fi

TAG="$1"
VERSION="${TAG#v}"
BINARY="heyad"

echo "Building release ${TAG}..."

for os in linux darwin; do
    for arch in amd64 arm64; do
        echo "  ${os}-${arch}..."
        GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 go build \
            -ldflags "-s -w -X github.com/cosmos/cosmos-sdk/version.Name=heya \
                -X github.com/cosmos/cosmos-sdk/version.AppName=heyad \
                -X github.com/cosmos/cosmos-sdk/version.Version=${TAG} \
                -X github.com/cosmos/cosmos-sdk/version.Commit=$(git rev-parse HEAD)" \
            -o "build/${BINARY}-${os}-${arch}" ./cmd/heyad/
        tar -czf "build/heya-${TAG}-${os}-${arch}.tar.gz" \
            -C build "${BINARY}-${os}-${arch}"
        rm "build/${BINARY}-${os}-${arch}"
    done
done

cd build
sha256sum heya-${TAG}-*.tar.gz > checksums.txt
echo ""
echo "Artifacts in build/:"
ls -lh
echo ""
echo "Upload with:"
echo "  gh release upload ${TAG} build/* --clobber"
