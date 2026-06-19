# Heya Chain - Lokalna sieć Cosmos SDK

Nazwa: **Heya** (HEYA)
Denom: `uheya` (1 HEYA = 1,000,000 uheya)
Chain ID: `heya-1`
Adresy: `heya1...`

## Struktura projektu

```
/root/heya/
  app/              # Główna aplikacja Cosmos SDK
  cmd/heyad/      # Binary węzła (heyad)
  config.yml        # Konfiguracja Ignite CLI
  docs/             # Dokumentacja API
  proto/            # Pliki protobuf
  testutil/         # Test helpers
```

## Uruchomienie

```bash
# Inicjalizacja (jeśli potrzeba od nowa)
rm -rf ~/.heya
heyad init my-moniker --chain-id heya-1

# Dodanie kont
heyad keys add alice --keyring-backend test
heyad keys add bob --keyring-backend test

# Dodanie do genesis
heyad genesis add-genesis-account alice 100000000000uheya --keyring-backend test
heyad genesis add-genesis-account bob 50000000000uheya --keyring-backend test

# Stworzenie walidatora
heyad genesis gentx alice 50000000000uheya --keyring-backend test --chain-id heya-1
heyad genesis collect-gentxs

# Konfiguracja denom in genesis (stake -> uheya)
jq '.app_state.crisis.constant_fee.denom = "uheya" |
    .app_state.gov.params.min_deposit[0].denom = "uheya" |
    .app_state.mint.params.mint_denom = "uheya" |
    .app_state.staking.params.bond_denom = "uheya"' \
    ~/.heya/config/genesis.json > tmp.json && mv tmp.json ~/.heya/config/genesis.json

# Uruchomienie
heyad start
```

## Przydatne komendy

```bash
# Status węzła
heyad status

# Balans konta
heyad query bank balances heya1...

# Wysyłanie tokenów
heyad tx bank send alice heya1... 1000000uheya --keyring-backend test --chain-id heya-1

# Staking (delegacja)
heyad tx staking delegate heyavaloper1... 1000000uheya --from alice --keyring-backend test --chain-id heya-1

# Nagrody walidatora
heyad query distribution rewards heya1... --chain-id heya-1

# Propozycja governance
heyad tx gov submit-proposal --title "Test" --description "Test" --deposit 10000000uheya --from alice --chain-id heya-1 --keyring-backend test
```

## Porty (domyślnie)
- RPC: 26657
- P2P: 26656
- gRPC: 9090
- REST API: 1317
