# Heya Chain - Cosmos SDK Local Network

Name: **Heya** (HEYA)
Denom: `uheya` (1 HEYA = 1,000,000 uheya)
Chain ID: `heya-1`
Addresses: `heya1...`

## Project Structure

```
/root/heya/
  app/              # Main Cosmos SDK application
  cmd/heyad/      # Node binary (heyad)
  config.yml        # Ignite CLI configuration
  docs/             # API documentation
  proto/            # Protobuf files
  testutil/         # Test helpers
```

## Run (release)

```bash
# Build with PGO and PIE
make build

# Run with GC optimizations
GOMEMLIMIT=$(( $(grep MemTotal /proc/meminfo | awk '{print $2}') * 85 / 100 / 1024 ))MiB \
  GOGC=100 \
  ./build/heyad start
```

## Development (debug)

```bash
make install-debug
heyad start
```

# Run (init + start)

```bash
# Initialize (if starting fresh)
rm -rf ~/.heya
heyad init my-moniker --chain-id heya-1

# Add accounts
heyad keys add alice --keyring-backend test
heyad keys add bob --keyring-backend test

# Add to genesis
heyad genesis add-genesis-account alice 100000000000uheya --keyring-backend test
heyad genesis add-genesis-account bob 50000000000uheya --keyring-backend test

# Create validator
heyad genesis gentx alice 50000000000uheya --keyring-backend test --chain-id heya-1
heyad genesis collect-gentxs

# Configure denom in genesis (stake -> uheya)
jq '.app_state.crisis.constant_fee.denom = "uheya" |
    .app_state.gov.params.min_deposit[0].denom = "uheya" |
    .app_state.mint.params.mint_denom = "uheya" |
    .app_state.staking.params.bond_denom = "uheya"' \
    ~/.heya/config/genesis.json > tmp.json && mv tmp.json ~/.heya/config/genesis.json

# Start
heyad start
```

## Useful Commands

```bash
# Node status
heyad status

# Account balance
heyad query bank balances heya1...

# Send tokens
heyad tx bank send alice heya1... 1000000uheya --keyring-backend test --chain-id heya-1

# Staking (delegate)
heyad tx staking delegate heyavaloper1... 1000000uheya --from alice --keyring-backend test --chain-id heya-1

# Validator rewards
heyad query distribution rewards heya1... --chain-id heya-1

# Governance proposal
heyad tx gov submit-proposal --title "Test" --description "Test" --deposit 10000000uheya --from alice --chain-id heya-1 --keyring-backend test
```

## Ports (default)
- RPC: 26657
- P2P: 26656
- gRPC: 9090
- REST API: 1317
