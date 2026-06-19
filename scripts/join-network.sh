#!/bin/bash
set -e

SEED_NODE_ID="1efe4ede5860cd60a36d0161df60fc3e31c2a038"
SEED_IP="178.63.164.6"
CHAIN_ID="nebula-1"
DENOM="unebula"
BINARY="nebulad"
BINARY_PATH="/usr/local/bin/nebulad"
NEBULA_DIR="/root/nebula"

print_step() { echo -e "\n\e[1;34m>>> $1\e[0m"; }

print_step "Budowanie Nebula binary z zrodel..."
if [ ! -f "$NEBULA_DIR/go.mod" ]; then
    echo "Blad: Nie znaleziono kodu zrodlowego w $NEBULA_DIR"
    echo "Sklonuj repo: git clone <repo-url> $NEBULA_DIR"
    exit 1
fi

cd "$NEBULA_DIR"
CGO_ENABLED=1 go build -ldflags "-X github.com/cosmos/cosmos-sdk/version.Name=nebula \
    -X github.com/cosmos/cosmos-sdk/version.AppName=nebulad \
    -X github.com/cosmos/cosmos-sdk/version.Version=v1.0.0 \
    -X github.com/cosmos/cosmos-sdk/version.Commit=$(git rev-parse HEAD 2>/dev/null || echo 'unknown') \
    -X 'github.com/cosmos/cosmos-sdk/version.BuildTags=cosmwasm wasm'" \
    -o "$BINARY_PATH" ./cmd/nebulad/

print_step "Inicjalizacja node'a..."
$BINARY init "$(hostname)" --chain-id $CHAIN_ID

print_step "Konfiguracja persistent_peers..."
PEERS="${SEED_NODE_ID}@${SEED_IP}:26656"
sed -i "s/^persistent_peers = .*/persistent_peers = \"$PEERS\"/" ~/.nebula/config/config.toml

print_step "Konfiguracja seed_peers..."
sed -i "s/^seeds = .*/seeds = \"$PEERS\"/" ~/.nebula/config/config.toml

print_step "Konfiguracja app.toml (min gas price)..."
sed -i 's/^minimum-gas-prices = .*/minimum-gas-prices = "0.025unebula"/' ~/.nebula/config/app.toml

print_step "Systemd service..."
cat > /etc/systemd/system/nebulad.service <<EOF
[Unit]
Description=Nebula Node
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

print_step "Gotowe! Node synchronizuje sie z siecia Nebula."
echo "Sprawdz: journalctl -u nebulad -f"
echo "Status:  $BINARY status"
