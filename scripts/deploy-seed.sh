#!/bin/bash
set -e

BINARY="heyad"
BINARY_PATH="$(go env GOPATH)/bin/heyad"
HEYA_DIR="$(dirname "$(dirname "$(realpath "$0")")")"
CHAIN_ID="heya-1"
SEED_NODE_ID="e1d96e06e0844b787e94393f3ab5594c39c5b234"
SEED_IP="178.63.164.6"

print_step() { echo -e "\n\e[1;34m>>> $1\e[0m"; }

print_step "Checking dependencies..."
if ! command -v go &>/dev/null; then
    echo "Installing Go 1.26..."
    wget -q https://go.dev/dl/go1.26.4.linux-amd64.tar.gz -O /tmp/go.tar.gz
    tar -C /usr/local -xzf /tmp/go.tar.gz
    export PATH="/usr/local/go/bin:$PATH"
    echo 'export PATH="/usr/local/go/bin:$PATH"' >> /root/.bashrc
fi

if [ ! -f /usr/lib/x86_64-linux-gnu/libwasmvm.x86_64.so ]; then
    echo "Installing libwasmvm..."
    wget -q "https://github.com/CosmWasm/wasmvm/releases/download/v2.2.7/libwasmvm.x86_64.so" \
        -O /usr/lib/x86_64-linux-gnu/libwasmvm.x86_64.so
    ldconfig
fi

print_step "Cloning / updating source..."
if [ ! -d "$HEYA_DIR" ]; then
    git clone https://github.com/heyanetwork/heya.git "$HEYA_DIR"
fi

cd "$HEYA_DIR"
git pull

print_step "Building binary..."
CGO_ENABLED=1 go build -ldflags "-X github.com/cosmos/cosmos-sdk/version.Name=heya \
    -X github.com/cosmos/cosmos-sdk/version.AppName=heyad \
    -X github.com/cosmos/cosmos-sdk/version.Version=v1.0.0 \
    -X github.com/cosmos/cosmos-sdk/version.Commit=$(git rev-parse HEAD) \
    -X 'github.com/cosmos/cosmos-sdk/version.BuildTags=cosmwasm wasm'" \
    -o "$BINARY_PATH" ./cmd/heyad/

print_step "Initializing as seed node..."
$BINARY init "heya-seed" --chain-id "$CHAIN_ID"

print_step "Configuring seed_mode..."
sed -i 's/^seed_mode = .*/seed_mode = true/' ~/.heya/config/config.toml
sed -i 's/^seeds = .*/seeds = ""/' ~/.heya/config/config.toml
sed -i 's/^persistent_peers = .*/persistent_peers = ""/' ~/.heya/config/config.toml
sed -i 's/^laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:26656"/' ~/.heya/config/config.toml

print_step "Configuring app.toml..."
sed -i 's/^minimum-gas-prices = .*/minimum-gas-prices = "0.025uheya"/' ~/.heya/config/app.toml

print_step "Copying genesis.json..."
cp "$HEYA_DIR/genesis.json" ~/.heya/config/genesis.json

print_step "Systemd service..."
cat > /etc/systemd/system/heyad.service <<EOF
[Unit]
Description=Heya Seed Node
After=network-online.target

[Service]
User=root
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

print_step "Seed node started!"
echo "ID:    $($BINARY tendermint show-node-id)"
echo "Logs:  journalctl -u heyad -f"
