package tokenfactory

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"heya/x/tokenfactory/types"
)

type msgServer struct {
	keeper Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (s msgServer) CreateDenom(goCtx context.Context, msg *types.MsgCreateDenom) (*types.MsgCreateDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	denom := types.NewDenom(msg.Sender, msg.Subdenom)
	_, exists, err := s.keeper.GetDenomAdmin(ctx, denom)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, types.ErrDenomExists
	}

	fee := types.DefaultParams().DenomCreationFee
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	if err := s.keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, sdk.NewCoins(fee)); err != nil {
		return nil, err
	}
	if err := s.keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(fee)); err != nil {
		return nil, err
	}

	if err := s.keeper.SetDenomAdmin(ctx, denom, msg.Sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("creator", msg.Sender),
		sdk.NewAttribute("denom", denom),
		sdk.NewAttribute("subdenom", msg.Subdenom),
	))
	return &types.MsgCreateDenomResponse{Denom: denom}, nil
}

func (s msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	coin, err := sdk.ParseCoinNormalized(msg.Amount)
	if err != nil {
		return nil, err
	}
	if err := validateFactoryDenom(coin.Denom); err != nil {
		return nil, err
	}
	admin, exists, err := s.keeper.GetDenomAdmin(ctx, coin.Denom)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, types.ErrDenomNotFound
	}
	if admin != msg.Sender {
		return nil, types.ErrUnauthorized
	}
	if err := s.keeper.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}
	recipient := msg.MintTo
	if recipient == "" {
		recipient = msg.Sender
	}
	recipientAddr, err := sdk.AccAddressFromBech32(recipient)
	if err != nil {
		return nil, err
	}
	if err := s.keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddr, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("minter", msg.Sender),
		sdk.NewAttribute("recipient", recipient),
		sdk.NewAttribute("denom", coin.Denom),
		sdk.NewAttribute("amount", coin.Amount.String()),
	))
	return &types.MsgMintResponse{}, nil
}

func (s msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	coin, err := sdk.ParseCoinNormalized(msg.Amount)
	if err != nil {
		return nil, err
	}
	if err := validateFactoryDenom(coin.Denom); err != nil {
		return nil, err
	}
	admin, exists, err := s.keeper.GetDenomAdmin(ctx, coin.Denom)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, types.ErrDenomNotFound
	}
	if admin != msg.Sender {
		return nil, types.ErrUnauthorized
	}
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	if err := s.keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}
	if err := s.keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("burner", msg.Sender),
		sdk.NewAttribute("denom", coin.Denom),
		sdk.NewAttribute("amount", coin.Amount.String()),
	))
	return &types.MsgBurnResponse{}, nil
}

func (s msgServer) ChangeAdmin(goCtx context.Context, msg *types.MsgChangeAdmin) (*types.MsgChangeAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := validateFactoryDenom(msg.Denom); err != nil {
		return nil, err
	}
	admin, exists, err := s.keeper.GetDenomAdmin(ctx, msg.Denom)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, types.ErrDenomNotFound
	}
	if admin != msg.Sender {
		return nil, types.ErrUnauthorized
	}
	if err := s.keeper.SetDenomAdmin(ctx, msg.Denom, msg.NewAdmin); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("old_admin", msg.Sender),
		sdk.NewAttribute("new_admin", msg.NewAdmin),
		sdk.NewAttribute("denom", msg.Denom),
	))
	return &types.MsgChangeAdminResponse{}, nil
}

func (s msgServer) ForceTransfer(goCtx context.Context, msg *types.MsgForceTransfer) (*types.MsgForceTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	coin, err := sdk.ParseCoinNormalized(msg.Amount)
	if err != nil {
		return nil, err
	}
	if err := validateFactoryDenom(coin.Denom); err != nil {
		return nil, err
	}
	admin, exists, err := s.keeper.GetDenomAdmin(ctx, coin.Denom)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, types.ErrDenomNotFound
	}
	if admin != msg.Sender {
		return nil, types.ErrUnauthorized
	}

	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	destAddr, err := sdk.AccAddressFromBech32(msg.DestAddr)
	if err != nil {
		return nil, err
	}

	if err := s.keeper.bankKeeper.SendCoins(ctx, senderAddr, destAddr, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("force_transfer_from", msg.Sender),
		sdk.NewAttribute("force_transfer_to", msg.DestAddr),
		sdk.NewAttribute("denom", coin.Denom),
		sdk.NewAttribute("amount", coin.Amount.String()),
	))
	return &types.MsgForceTransferResponse{}, nil
}

func validateFactoryDenom(denom string) error {
	if !strings.HasPrefix(denom, types.DenomPrefix+"/") {
		return types.ErrInvalidDenom
	}
	return nil
}
