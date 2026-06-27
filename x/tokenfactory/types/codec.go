package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateDenom{}, "heya/tokenfactory/MsgCreateDenom", nil)
	cdc.RegisterConcrete(&MsgMint{}, "heya/tokenfactory/MsgMint", nil)
	cdc.RegisterConcrete(&MsgBurn{}, "heya/tokenfactory/MsgBurn", nil)
	cdc.RegisterConcrete(&MsgChangeAdmin{}, "heya/tokenfactory/MsgChangeAdmin", nil)
	cdc.RegisterConcrete(&MsgForceTransfer{}, "heya/tokenfactory/MsgForceTransfer", nil)
	cdc.RegisterConcrete(&MsgAcceptAdmin{}, "heya/tokenfactory/MsgAcceptAdmin", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "heya/tokenfactory/MsgUpdateParams", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateDenom{},
		&MsgMint{},
		&MsgBurn{},
		&MsgChangeAdmin{},
		&MsgForceTransfer{},
		&MsgAcceptAdmin{},
		&MsgUpdateParams{},
	)
}
