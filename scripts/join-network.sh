#!/bin/bash
set -e

SEED_NODE_ID="e1d96e06e0844b787e94393f3ab5594c39c5b234"
SEED_IP="178.63.164.6"
CHAIN_ID="heya-1"
DENOM="uheya"
BINARY="heyad"
BINARY_PATH="$(go env GOPATH)/bin/heyad"
HEYA_DIR="/root/heya"

print_step() { echo -e "\n\e[1;34m>>> $1\e[0m"; }

print_step "Budowanie Heya binary z zrodel..."
if [ ! -f "$HEYA_DIR/go.mod" ]; then
    echo "Blad: Nie znaleziono kodu zrodlowego w $HEYA_DIR"
    echo "Sklonuj repo: git clone https://github.com/heya-protocol/heya.git $HEYA_DIR"
    exit 1
fi

cd "$HEYA_DIR"
CGO_ENABLED=1 go build -ldflags "-X github.com/cosmos/cosmos-sdk/version.Name=heya \
    -X github.com/cosmos/cosmos-sdk/version.AppName=heyad \
    -X github.com/cosmos/cosmos-sdk/version.Version=v1.0.0 \
    -X github.com/cosmos/cosmos-sdk/version.Commit=$(git rev-parse HEAD 2>/dev/null || echo 'unknown') \
    -X 'github.com/cosmos/cosmos-sdk/version.BuildTags=cosmwasm wasm'" \
    -o "$BINARY_PATH" ./cmd/heyad/

print_step "Inicjalizacja node'a..."
$BINARY init "$(hostname)" --chain-id $CHAIN_ID

print_step "Konfiguracja persistent_peers..."
PEERS="${SEED_NODE_ID}@${SEED_IP}:26656"
sed -i "s/^persistent_peers = .*/persistent_peers = \"$PEERS\"/" ~/.heya/config/config.toml

print_step "Konfiguracja seed_peers..."
sed -i "s/^seeds = .*/seeds = \"$PEERS\"/" ~/.heya/config/config.toml

print_step "Konfiguracja app.toml (min gas price)..."
sed -i 's/^minimum-gas-prices = .*/minimum-gas-prices = "0.025uheya"/' ~/.heya/config/app.toml

print_step "Systemd service..."
cat > /etc/systemd/system/heyad.service <<EOF
[Unit]
Description=Heya Node
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

print_step "Gotowe! Node synchronizuje sie z siecia Heya."
echo "Sprawdz: journalctl -u heyad -f"
echo "Status:  $BINARY status"
