#!/bin/bash
set -e

SEED_NODE_ID="e1d96e06e0844b787e94393f3ab5594c39c5b234"
SEED_IP="178.63.164.6"
CHAIN_ID="heya-1"
DENOM="uheya"
BINARY="heyad"
BINARY_PATH="/usr/local/bin/heyad"
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
GENESIS_FILE="$SCRIPT_DIR/genesis.json"

print_step() { echo -e "\n\e[1;34m>>> $1\e[0m"; }

BUILD_MODE="download"
if [ "$#" -eq 0 ]; then
    echo ""
    echo "Select installation method:"
    echo "  1) Download pre-built binary (default)"
    echo "  2) Build from source"
    read -r -p "Choice [1/2]: " CHOICE
    case "$CHOICE" in
        2|build) BUILD_MODE="build" ;;
        *)       BUILD_MODE="download" ;;
    esac
else
    for arg in "$@"; do
        case "$arg" in
            --build) BUILD_MODE="build" ;;
            --download) BUILD_MODE="download" ;;
        esac
    done
fi

ARCH="$(uname -m)"
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64)   ARCH="arm64" ;;
    *)       echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

print_step "Fetching latest release info from GitHub..."
API_URL="https://api.github.com/repos/heyanetwork/heya/releases/latest"
TAG=$(curl -sS "$API_URL" | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": "\(.*\)",/\1/')
VERSION="${TAG#v}"

if [ "$BUILD_MODE" = "build" ]; then
    print_step "Building Heya binary from source (${TAG})..."
    if ! command -v go &>/dev/null; then
        echo "Error: Go is not installed. Install it first: https://go.dev/doc/install"
        exit 1
    fi
    SRC_URL="https://github.com/heyanetwork/heya/archive/refs/tags/${TAG}.tar.gz"
    echo "Downloading source ${TAG}..."
    curl -sSL "$SRC_URL" -o /tmp/heya-src.tar.gz
    tar -xzf /tmp/heya-src.tar.gz -C /tmp/
    cd "/tmp/heya-${VERSION}"
    CGO_ENABLED=1 go build -ldflags "-X github.com/cosmos/cosmos-sdk/version.Name=heya \
        -X github.com/cosmos/cosmos-sdk/version.AppName=heyad \
        -X github.com/cosmos/cosmos-sdk/version.Version=${TAG} \
        -X github.com/cosmos/cosmos-sdk/version.Commit=$(git rev-parse HEAD 2>/dev/null || echo 'unknown') \
        -X 'github.com/cosmos/cosmos-sdk/version.BuildTags=cosmwasm wasm'" \
        -o "$BINARY_PATH" ./cmd/heyad/
    rm -f /tmp/heya-src.tar.gz
    rm -rf "/tmp/heya-${VERSION}"
    echo "Built $BINARY ${TAG} from source to $BINARY_PATH"
else
    print_step "Downloading latest Heya release from GitHub..."
    FILENAME="heya-${VERSION}-${OS}-${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/heyanetwork/heya/releases/download/${TAG}/${FILENAME}"
    echo "Downloading $FILENAME..."
    curl -sSL "$DOWNLOAD_URL" -o /tmp/heya.tar.gz
    tar -xzf /tmp/heya.tar.gz -C /tmp/
    cp "/tmp/heya-${VERSION}-${OS}-${ARCH}/heyad" "$BINARY_PATH"
    chmod +x "$BINARY_PATH"
    rm -f /tmp/heya.tar.gz
    rm -rf "/tmp/heya-${VERSION}-${OS}-${ARCH}"
    echo "Downloaded $BINARY v$VERSION to $BINARY_PATH"
fi

print_step "Checking for existing home directory..."
HEYA_HOME="$HOME/.heya"
if [ -d "$HEYA_HOME/config" ]; then
    echo "Found existing config, backing up keys..."
    [ -f "$HEYA_HOME/config/priv_validator_key.json" ] && cp "$HEYA_HOME/config/priv_validator_key.json" "$HEYA_HOME/config/priv_validator_key.json.bak"
    [ -f "$HEYA_HOME/config/node_key.json" ] && cp "$HEYA_HOME/config/node_key.json" "$HEYA_HOME/config/node_key.json.bak"
fi

print_step "Initializing node (generates config files)..."
"$BINARY_PATH" init "$(hostname)" --chain-id $CHAIN_ID --overwrite 2>/dev/null

print_step "Restoring keys and genesis from repository..."
if [ -f "$HEYA_HOME/config/priv_validator_key.json.bak" ]; then
    mv "$HEYA_HOME/config/priv_validator_key.json.bak" "$HEYA_HOME/config/priv_validator_key.json"
fi
if [ -f "$HEYA_HOME/config/node_key.json.bak" ]; then
    mv "$HEYA_HOME/config/node_key.json.bak" "$HEYA_HOME/config/node_key.json"
fi
if [ ! -f "$GENESIS_FILE" ]; then
    echo "Error: $GENESIS_FILE not found!"
    exit 1
fi
cp "$GENESIS_FILE" "$HEYA_HOME/config/genesis.json"
echo "Genesis copied from repository"

print_step "Configuring persistent_peers..."
PEERS="${SEED_NODE_ID}@${SEED_IP}:26656"
sed -i "s/^persistent_peers = .*/persistent_peers = \"$PEERS\"/" ~/.heya/config/config.toml

print_step "Configuring seed_peers..."
sed -i "s/^seeds = .*/seeds = \"$PEERS\"/" ~/.heya/config/config.toml

print_step "Configuring app.toml (min gas price)..."
sed -i 's/^minimum-gas-prices = .*/minimum-gas-prices = "0.025uheya"/' ~/.heya/config/app.toml

print_step "Setting up systemd service..."
cat > /etc/systemd/system/heyad.service <<EOF
[Unit]
Description=Heya Node
After=network-online.target

[Service]
User=$(whoami)
ExecStart=$BINARY_PATH start
Restart=on-failure
RestartSec=5
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable heyad
systemctl start heyad

print_step "Done! Node is syncing with Heya network."
echo "Check: journalctl -u heyad -f"
echo "Status:  $BINARY status"
