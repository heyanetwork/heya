# Nebula Chain - Summary

## Goal
Build a Cosmos SDK blockchain (Nebula) with a native coin (NEB), max supply 10B, and CosmWasm smart contracts.

## Constraints & Preferences
- Run locally with systemd (`nebulad.service`)
- Address prefix: `nebula`
- Denom: `unebula` (1 NEB = 1,000,000 unebula)
- Chain ID: `nebula-1`
- Max supply: 10,000,000,000 NEB (hard cap via custom supplycap module)
- Genesis circulation: 5,000,000,000 NEB (50%)
- Seed node baked into binary: `1efe4ede5860cd60a36d0161df60fc3e31c2a038@178.63.164.6:26656`
- Block time: ~5s (default `timeout_commit`)
- User IP for peer discovery: 178.63.164.6

## Progress
### Done
- Scaffolded chain with Ignite CLI v28.5.0 (Cosmos SDK v0.50.8)
- Built binary `nebulad` and installed to `/root/go/bin/nebulad`
- Configured `app/config.go` (address prefix, min gas price `0.025unebula`)
- Created systemd service at `/etc/systemd/system/nebulad.service` (active, enabled)
- Created custom `app/supplycap` module – zeros inflation params when total supply ≥ 10B NEB
- Registered supplycap module in `app.RegisterModules()` and in `beginBlockers` before mint
- Genesis configured: `unebula` denom everywhere, 3 accounts:
  - alice: 2,000,000,000 NEB available + 500,000,000 NEB staked (validator NEBULA 1)
  - bob: 1,500,000,000 NEB
  - community: 1,000,000,000 NEB
- Added validator description (identity, website, security-contact, details)
- Created `/root/nebula/join-network.sh` for new nodes
- Enabled `seed_mode = true`, PEX on main node
- Embedded seeds/persistent_peers in `initCometBFTConfig()` so new nodes auto-connect
- Integrated CosmWasm: added `github.com/CosmWasm/wasmd@v0.54.8`, `github.com/CosmWasm/wasmvm/v2@v2.2.7`
- Installed `libwasmvm.x86_64.so` (pre-built shared library from wasmvm v2.2.7 release)
- Updated `app/app.go`, `app/app_config.go`, `cmd/nebulad/cmd/root.go` for wasm integration
- Reset chain state via `nebulad unsafe-reset-all` to handle new wasm store key
- Set binary version via ldflags: `-X .../version.Name=nebula -X .../version.Version=v1.0.0 -X .../version.BuildTags='cosmwasm wasm'`
- Chain running and producing blocks; wasm module functional
- **Dependency audit completed** – 9 packages upgraded successfully
- **Upgrade handler framework** – `app/upgrades.go` with `setupUpgradeHandlers()` for future store upgrades (uses `StoreUpgrades` + `UpgradeStoreLoader` on scheduled upgrades; no more `unsafe-reset-all` needed)
- **Deployment scripts** – `scripts/deploy-seed.sh` (new seed/production node setup), `scripts/export-genesis.sh` (genesis export), `join-network.sh` updated (build from source):

### Upgraded Packages (SDK v0.50.15 compatible)
| Package | Old | New |
|---------|-----|-----|
| cosmossdk.io/math | v1.5.0 | v1.5.3 |
| cosmossdk.io/store | v1.1.1 | v1.1.2 |
| cosmossdk.io/depinject | v1.1.0 | v1.2.1 |
| github.com/cometbft/cometbft | v0.38.21 | v0.38.23 |
| github.com/cosmos/gogoproto | v1.7.0 | v1.7.2 |
| github.com/spf13/cobra | v1.9.1 | v1.10.2 |
| github.com/spf13/pflag | v1.0.6 | v1.0.10 |
| github.com/spf13/viper | v1.19.0 | v1.21.0 |
| github.com/spf13/cast | v1.7.1 | v1.10.0 |

### Skipped (already latest or incompatible)
- `cosmossdk.io/core` v0.12.0 retracted (incompatible with v0.50)
- `cosmossdk.io/api` v0.9.2 – needs SDK v0.52+
- `cosmossdk.io/x/*` v0.2.0 – needs SDK v0.52+
- `client/v2` – no stable for SDK v0.50
- `wasmd` v0.55.x – needs SDK v0.52+
- `cosmos-sdk` – already at latest v0.50.15
- `ibc-go/v8` – already at latest v8.8.0
- `log`, `confix`, `testify`, `protobuf` – already at latest

### In Progress
(none)

### Blocked
(none)

## Key Decisions
- Max supply of 10B NEB enforced via custom `app/supplycap` module (BeginBlocker checks supply, zeros inflation params), not at SDK level.
- 50% supply at genesis (5B NEB) to allow inflation-based emission for remaining 5B.
- Seed node data baked into binary source (`cmd/nebulad/cmd/config.go`), not config file.
- Validator renamed to "NEBULA 1" via `tx staking edit-validator`.
- Use `GOPROXY='https://proxy.golang.org,https://goproxy.io,direct'` – proxy.golang.org needed for some packages that goproxy.io can't serve, and vice versa.
- Store upgrades now handled via `StoreUpgrades` + `UpgradeStoreLoader` in `app/upgrades.go`; new modules can be added without resetting chain state.

## Next Steps
- Set up GitHub releases for binary distribution (update `join-network.sh` download URL).
- Integrate IBC wasm light client (ics721) for cross-chain wasm contract migration.
- Add monitoring/alerting (prometheus metrics already available).

## Critical Context
- **Go**: 1.26.4 (`/usr/local/go/bin/go`)
- **Proxy**: `GOPROXY='https://proxy.golang.org,https://goproxy.io,direct'` – neither proxy alone serves all packages; proxy.golang.org returns 403 for some, goproxy.io fails on others (e.g. old `golang.org/x/sys`). Use both with fallback.
- **Build**:
  ```
  cd /root/nebula
  CGO_ENABLED=1 go build -ldflags "-X github.com/cosmos/cosmos-sdk/version.Name=nebula \
    -X github.com/cosmos/cosmos-sdk/version.AppName=nebulad \
    -X github.com/cosmos/cosmos-sdk/version.Version=v1.0.0 \
    -X github.com/cosmos/cosmos-sdk/version.Commit=$(git rev-parse HEAD) \
    -X 'github.com/cosmos/cosmos-sdk/version.BuildTags=cosmwasm wasm'" \
    -o /root/go/bin/nebulad ./cmd/nebulad/
  ```
- **libwasmvm**: loaded from `/root/go/pkg/mod/github.com/\!cosm\!wasm/wasmvm/v2@v2.2.7/internal/api/libwasmvm.x86_64.so` or system `/usr/lib/x86_64-linux-gnu/libwasmvm.x86_64.so`
- **Store upgrades**: New module store keys added via `app.RegisterStores()`. For existing chains, schedule upgrade via governance → handler in `app/upgrades.go` applies `StoreUpgrades` with `UpgradeStoreLoader`.
- **Validator**: NEBULA 1, power 500M unebula, moniker `nebula-test`, node ID `1efe4ede5860cd60a36d0161df60fc3e31c2a038` (seed mode).

## Relevant Files
- `/root/nebula/`: chain root
- `/root/nebula/app/supplycap/module.go`: max supply enforcement (10B NEB cap)
- `/root/nebula/cmd/nebulad/cmd/config.go`: default seeds and persistent_peers baked into binary
- `/root/nebula/app/app.go`: app structure, wasm keeper creation
- `/root/nebula/app/app_config.go`: module ordering, permissions
- `/root/nebula/cmd/nebulad/cmd/root.go`: CLI registration
- `/root/nebula/config.yml`: Ignite config
- `/root/nebula/app/upgrades.go`: upgrade handler framework (v2 upgrade with StoreUpgrades)
- `/root/nebula/join-network.sh`: automated new-node join script (builds from source)
- `/root/nebula/scripts/deploy-seed.sh`: production seed node deployment
- `/root/nebula/scripts/export-genesis.sh`: genesis export from live chain
- `/root/nebula/go.mod`: dependency list (audited and upgraded)
- `/etc/systemd/system/nebulad.service`: systemd unit
- `/root/.nebula/config/`: data dir
