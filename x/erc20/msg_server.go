package erc20

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"heya/x/erc20/types"
)

type msgServer struct {
	types.UnimplementedMsgServer
	keeper Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (s msgServer) Transfer(goCtx context.Context, msg *types.MsgTransfer) (*types.MsgTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	coin, err := sdk.ParseCoinNormalized(msg.Amount)
	if err != nil {
		return nil, err
	}

	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	to, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return nil, err
	}

	if err := s.keeper.bankKeeper.SendCoins(ctx, from, to, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("action", "transfer"),
		sdk.NewAttribute("from", msg.Sender),
		sdk.NewAttribute("to", msg.To),
		sdk.NewAttribute("denom", coin.Denom),
		sdk.NewAttribute("amount", coin.Amount.String()),
	))
	return &types.MsgTransferResponse{}, nil
}

func (s msgServer) Approve(goCtx context.Context, msg *types.MsgApprove) (*types.MsgApproveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	coin, err := sdk.ParseCoinNormalized(msg.Amount)
	if err != nil {
		return nil, err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return nil, err
	}
	if _, err := sdk.AccAddressFromBech32(msg.Spender); err != nil {
		return nil, err
	}

	if err := s.keeper.SetAllowance(ctx, msg.Owner, msg.Spender, coin.Denom, coin.Amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("action", "approve"),
		sdk.NewAttribute("owner", msg.Owner),
		sdk.NewAttribute("spender", msg.Spender),
		sdk.NewAttribute("denom", coin.Denom),
		sdk.NewAttribute("amount", coin.Amount.String()),
	))
	return &types.MsgApproveResponse{}, nil
}

func (s msgServer) TransferFrom(goCtx context.Context, msg *types.MsgTransferFrom) (*types.MsgTransferFromResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	coin, err := sdk.ParseCoinNormalized(msg.Amount)
	if err != nil {
		return nil, err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Caller); err != nil {
		return nil, err
	}
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}
	to, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return nil, err
	}

	allowance := s.keeper.GetAllowance(ctx, msg.From, msg.Caller, coin.Denom)
	if allowance.LT(coin.Amount) {
		return nil, types.ErrInsufficientAllowance.Wrapf(
			"allowance %s < transfer amount %s", allowance.String(), coin.Amount.String(),
		)
	}

	newAllowance := allowance.Sub(coin.Amount)
	if err := s.keeper.SetAllowance(ctx, msg.From, msg.Caller, coin.Denom, newAllowance); err != nil {
		return nil, err
	}

	if err := s.keeper.bankKeeper.SendCoins(ctx, from, to, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("action", "transfer_from"),
		sdk.NewAttribute("caller", msg.Caller),
		sdk.NewAttribute("from", msg.From),
		sdk.NewAttribute("to", msg.To),
		sdk.NewAttribute("denom", coin.Denom),
		sdk.NewAttribute("amount", coin.Amount.String()),
	))
	return &types.MsgTransferFromResponse{}, nil
}

func (s msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Authority != s.keeper.GetAuthority(ctx) {
		return nil, types.ErrUnauthorized.Wrap("only the module authority can update params")
	}

	if err := s.keeper.SetAuthority(ctx, msg.NewAuthority); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("action", "update_params"),
		sdk.NewAttribute("authority", msg.Authority),
		sdk.NewAttribute("new_authority", msg.NewAuthority),
	))
	return &types.MsgUpdateParamsResponse{}, nil
}

func (s msgServer) IncreaseAllowance(goCtx context.Context, msg *types.MsgIncreaseAllowance) (*types.MsgIncreaseAllowanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	amt, ok := sdkmath.NewIntFromString(msg.Amount)
	if !ok || !amt.IsPositive() {
		return nil, fmt.Errorf("amount must be positive integer")
	}

	current := s.keeper.GetAllowance(ctx, msg.Owner, msg.Spender, msg.Denom)
	newAllowance := current.Add(amt)

	if err := s.keeper.SetAllowance(ctx, msg.Owner, msg.Spender, msg.Denom, newAllowance); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("action", "increase_allowance"),
		sdk.NewAttribute("owner", msg.Owner),
		sdk.NewAttribute("spender", msg.Spender),
		sdk.NewAttribute("denom", msg.Denom),
		sdk.NewAttribute("amount", msg.Amount),
	))
	return &types.MsgIncreaseAllowanceResponse{}, nil
}

func (s msgServer) DecreaseAllowance(goCtx context.Context, msg *types.MsgDecreaseAllowance) (*types.MsgDecreaseAllowanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	amt, ok := sdkmath.NewIntFromString(msg.Amount)
	if !ok || !amt.IsPositive() {
		return nil, fmt.Errorf("amount must be positive integer")
	}

	current := s.keeper.GetAllowance(ctx, msg.Owner, msg.Spender, msg.Denom)
	if current.LT(amt) {
		return nil, types.ErrInsufficientAllowance.Wrapf(
			"allowance %s < decrease amount %s", current.String(), msg.Amount,
		)
	}

	newAllowance := current.Sub(amt)
	if err := s.keeper.SetAllowance(ctx, msg.Owner, msg.Spender, msg.Denom, newAllowance); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("action", "decrease_allowance"),
		sdk.NewAttribute("owner", msg.Owner),
		sdk.NewAttribute("spender", msg.Spender),
		sdk.NewAttribute("denom", msg.Denom),
		sdk.NewAttribute("amount", msg.Amount),
	))
	return &types.MsgDecreaseAllowanceResponse{}, nil
}
