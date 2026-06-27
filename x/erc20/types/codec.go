package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgTransfer{}, "heya/erc20/MsgTransfer", nil)
	cdc.RegisterConcrete(&MsgApprove{}, "heya/erc20/MsgApprove", nil)
	cdc.RegisterConcrete(&MsgTransferFrom{}, "heya/erc20/MsgTransferFrom", nil)
	cdc.RegisterConcrete(&MsgIncreaseAllowance{}, "heya/erc20/MsgIncreaseAllowance", nil)
	cdc.RegisterConcrete(&MsgDecreaseAllowance{}, "heya/erc20/MsgDecreaseAllowance", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "heya/erc20/MsgUpdateParams", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTransfer{},
		&MsgApprove{},
		&MsgTransferFrom{},
		&MsgIncreaseAllowance{},
		&MsgDecreaseAllowance{},
		&MsgUpdateParams{},
	)
}
