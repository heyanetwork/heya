package app

import "encoding/json"

// GenesisState defines the genesis state of the application.
type GenesisState map[string]json.RawMessage
