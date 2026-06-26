package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		DenomCreationFee: "1000000000uheya",
		Denoms:           []*DenomAuthority{},
	}
}

func (gs *GenesisState) Validate() error {
	if _, err := sdk.ParseCoinNormalized(gs.DenomCreationFee); err != nil {
		return err
	}
	seen := make(map[string]bool)
	for _, d := range gs.Denoms {
		if d.Denom == "" {
			return ErrInvalidDenom
		}
		if d.Admin == "" {
			return ErrInvalidCreator
		}
		if _, err := sdk.AccAddressFromBech32(d.Admin); err != nil {
			return ErrInvalidCreator
		}
		if seen[d.Denom] {
			return ErrDenomExists
		}
		seen[d.Denom] = true
	}
	return nil
}

func (gs *GenesisState) DenomCreationFeeCoin() (sdk.Coin, error) {
	return sdk.ParseCoinNormalized(gs.DenomCreationFee)
}
