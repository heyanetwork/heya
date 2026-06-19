#!/bin/bash
set -e

NEBULA_DIR="/root/nebula"
BINARY="nebulad"
BINARY_PATH="/usr/local/bin/nebulad"
CHAIN_ID="nebula-1"
SEED_NODE_ID="1efe4ede5860cd60a36d0161df60fc3e31c2a038"
SEED_IP="178.63.164.6"

print_step() { echo -e "\n\e[1;34m>>> $1\e[0m"; }

print_step "Sprawdzanie zależności..."
if ! command -v go &>/dev/null; then
    echo "Instalowanie Go 1.26..."
    wget -q https://go.dev/dl/go1.26.4.linux-amd64.tar.gz -O /tmp/go.tar.gz
    tar -C /usr/local -xzf /tmp/go.tar.gz
    export PATH="/usr/local/go/bin:$PATH"
    echo 'export PATH="/usr/local/go/bin:$PATH"' >> /root/.bashrc
fi

if [ ! -f /usr/lib/x86_64-linux-gnu/libwasmvm.x86_64.so ]; then
    echo "Instalowanie libwasmvm..."
    wget -q "https://github.com/CosmWasm/wasmvm/releases/download/v2.2.7/libwasmvm.x86_64.so" \
        -O /usr/lib/x86_64-linux-gnu/libwasmvm.x86_64.so
    ldconfig
fi

print_step "Klonowanie / aktualizacja źródła..."
if [ ! -d "$NEBULA_DIR" ]; then
    git clone <repo-url> "$NEBULA_DIR"
fi

cd "$NEBULA_DIR"
git pull

print_step "Budowanie binary..."
CGO_ENABLED=1 go build -ldflags "-X github.com/cosmos/cosmos-sdk/version.Name=nebula \
    -X github.com/cosmos/cosmos-sdk/version.AppName=nebulad \
    -X github.com/cosmos/cosmos-sdk/version.Version=v1.0.0 \
    -X github.com/cosmos/cosmos-sdk/version.Commit=$(git rev-parse HEAD) \
    -X 'github.com/cosmos/cosmos-sdk/version.BuildTags=cosmwasm wasm'" \
    -o "$BINARY_PATH" ./cmd/nebulad/

print_step "Inicjalizacja jako seed node..."
$BINARY init "nebula-seed" --chain-id "$CHAIN_ID"

print_step "Konfiguracja seed_mode..."
sed -i 's/^seed_mode = .*/seed_mode = true/' ~/.nebula/config/config.toml
sed -i 's/^seeds = .*/seeds = ""/' ~/.nebula/config/config.toml
sed -i 's/^persistent_peers = .*/persistent_peers = ""/' ~/.nebula/config/config.toml
sed -i 's/^laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:26656"/' ~/.nebula/config/config.toml

print_step "Konfiguracja app.toml..."
sed -i 's/^minimum-gas-prices = .*/minimum-gas-prices = "0.025unebula"/' ~/.nebula/config/app.toml

print_step "Kopiowanie genesis.json..."
cp "$NEBULA_DIR/genesis.json" ~/.nebula/config/genesis.json

print_step "Systemd service..."
cat > /etc/systemd/system/nebulad.service <<EOF
[Unit]
Description=Nebula Seed Node
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
systemctl enable nebulad
systemctl start nebulad

print_step "Seed node uruchomiony!"
echo "ID:    $($BINARY tendermint show-node-id)"
echo "Logi:  journalctl -u nebulad -f"