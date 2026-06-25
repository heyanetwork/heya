# Heya Blockchain

**heya** is a Cosmos SDK-based blockchain (v0.50.15) with IBC-go v8.8.0, CosmWasm v0.54.8, and CometBFT v0.38.23.

## Tokenomics

| Parametr | Wartość |
|---|---|
| Chain ID | `heya-1` |
| Denom | `uheya` (1 HEYA = 1,000,000 uheya) |
| Max supply | 10,000,000,000 HEYA |
| Initial supply | 5,000,000,000 HEYA |
| Inflation | ~13% rocznie |
| Bech32 prefix | `heya` |

## Endpoints (mainnet)

| Service | Address |
|---|---|
| RPC | `http://178.63.164.6:26657` |
| gRPC | `178.63.164.6:9090` |
| REST API | `http://178.63.164.6:1317` |
| P2P | `178.63.164.6:26656` |

## Build

```bash
go build -o build/heyad ./cmd/heyad/
```

## Genesis

Genesis available at `https://github.com/heyanetwork/heya` or on request.

## License

MIT
