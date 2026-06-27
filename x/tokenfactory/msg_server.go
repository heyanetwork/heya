package tokenfactory

import (
	"context"
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"heya/x/tokenfactory/types"
)

type msgServer struct {
	types.UnimplementedMsgServer
	keeper Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (s msgServer) checkPaused(ctx sdk.Context) error {
	if s.keeper.IsPaused(ctx) {
		return types.ErrPaused
	}
	return nil
}

func (s msgServer) CreateDenom(goCtx context.Context, msg *types.MsgCreateDenom) (*types.MsgCreateDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := s.checkPaused(ctx); err != nil {
		return nil, err
	}

	denom := types.NewDenom(msg.Sender, msg.Subdenom)
	_, exists, err := s.keeper.GetDenomAdmin(ctx, denom)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, types.ErrDenomExists
	}

	params, err := s.keeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	fee := params.DenomCreationFee
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

	// Set default supply cap: 1 billion tokens
	defaultCap := sdkmath.NewInt(1_000_000_000_000_000)
	if err := s.keeper.SetSupplyCap(ctx, denom, defaultCap); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("action", "create_denom"),
		sdk.NewAttribute("creator", msg.Sender),
		sdk.NewAttribute("denom", denom),
		sdk.NewAttribute("subdenom", msg.Subdenom),
	))
	return &types.MsgCreateDenomResponse{Denom: denom}, nil
}

func (s msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := s.checkPaused(ctx); err != nil {
		return nil, err
	}

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

	// Check supply cap
	cap, hasCap, err := s.keeper.GetSupplyCap(ctx, coin.Denom)
	if err != nil {
		return nil, err
	}
	if hasCap {
		currentSupply := s.keeper.bankKeeper.GetSupply(ctx, coin.Denom)
		newSupply := currentSupply.Amount.Add(coin.Amount)
		if newSupply.GT(cap) {
			return nil, types.ErrSupplyCapExceeded.Wrapf(
				"current supply %s + mint amount %s exceeds cap %s",
				currentSupply.Amount.String(), coin.Amount.String(), cap.String(),
			)
		}
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
		sdk.NewAttribute("action", "mint"),
		sdk.NewAttribute("minter", msg.Sender),
		sdk.NewAttribute("recipient", recipient),
		sdk.NewAttribute("denom", coin.Denom),
		sdk.NewAttribute("amount", coin.Amount.String()),
	))
	return &types.MsgMintResponse{}, nil
}

func (s msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := s.checkPaused(ctx); err != nil {
		return nil, err
	}

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
		sdk.NewAttribute("action", "burn"),
		sdk.NewAttribute("burner", msg.Sender),
		sdk.NewAttribute("denom", coin.Denom),
		sdk.NewAttribute("amount", coin.Amount.String()),
	))
	return &types.MsgBurnResponse{}, nil
}

func (s msgServer) ChangeAdmin(goCtx context.Context, msg *types.MsgChangeAdmin) (*types.MsgChangeAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := s.checkPaused(ctx); err != nil {
		return nil, err
	}

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

	// Two-step admin transfer: set pending admin instead of direct change
	if err := s.keeper.SetPendingAdmin(ctx, msg.Denom, msg.NewAdmin); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("action", "change_admin_pending"),
		sdk.NewAttribute("current_admin", msg.Sender),
		sdk.NewAttribute("pending_admin", msg.NewAdmin),
		sdk.NewAttribute("denom", msg.Denom),
	))
	return &types.MsgChangeAdminResponse{}, nil
}

func (s msgServer) AcceptAdmin(goCtx context.Context, msg *types.MsgAcceptAdmin) (*types.MsgAcceptAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := s.checkPaused(ctx); err != nil {
		return nil, err
	}

	if err := validateFactoryDenom(msg.Denom); err != nil {
		return nil, err
	}
	pendingAdmin, exists, err := s.keeper.GetPendingAdmin(ctx, msg.Denom)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, types.ErrPendingAdminNotFound
	}
	if pendingAdmin != msg.Admin {
		return nil, types.ErrNotPendingAdmin
	}

	if err := s.keeper.SetDenomAdmin(ctx, msg.Denom, msg.Admin); err != nil {
		return nil, err
	}
	if err := s.keeper.SetPendingAdmin(ctx, msg.Denom, ""); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("action", "accept_admin"),
		sdk.NewAttribute("new_admin", msg.Admin),
		sdk.NewAttribute("denom", msg.Denom),
	))
	return &types.MsgAcceptAdminResponse{}, nil
}

func (s msgServer) ForceTransfer(goCtx context.Context, msg *types.MsgForceTransfer) (*types.MsgForceTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := s.checkPaused(ctx); err != nil {
		return nil, err
	}

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

	// Determine the source address: from_address if set, otherwise admin's own address
	fromAddrStr := msg.FromAddress
	if fromAddrStr == "" {
		fromAddrStr = msg.Sender
	}
	fromAddr, err := sdk.AccAddressFromBech32(fromAddrStr)
	if err != nil {
		return nil, err
	}
	destAddr, err := sdk.AccAddressFromBech32(msg.DestAddr)
	if err != nil {
		return nil, err
	}

	if err := s.keeper.bankKeeper.SendCoins(ctx, fromAddr, destAddr, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("action", "force_transfer"),
		sdk.NewAttribute("admin", msg.Sender),
		sdk.NewAttribute("from_address", fromAddrStr),
		sdk.NewAttribute("to_address", msg.DestAddr),
		sdk.NewAttribute("denom", coin.Denom),
		sdk.NewAttribute("amount", coin.Amount.String()),
	))
	return &types.MsgForceTransferResponse{}, nil
}

func (s msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Only governance module can update params
	if msg.Authority != s.keeper.GetAuthority() {
		return nil, types.ErrUnauthorized.Wrap("only governance can update params")
	}

	coin, err := sdk.ParseCoinNormalized(msg.DenomCreationFee)
	if err != nil {
		return nil, err
	}
	params := types.Params{DenomCreationFee: coin}
	if err := params.Validate(); err != nil {
		return nil, err
	}
	if err := s.keeper.SetParams(ctx, params); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("action", "update_params"),
		sdk.NewAttribute("denom_creation_fee", msg.DenomCreationFee),
	))
	return &types.MsgUpdateParamsResponse{}, nil
}

func (s msgServer) UpdateSupplyCap(goCtx context.Context, msg *types.MsgUpdateSupplyCap) (*types.MsgUpdateSupplyCapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Authority != s.keeper.GetAuthority() {
		return nil, types.ErrUnauthorized.Wrap("only governance can update supply cap")
	}

	if err := validateFactoryDenom(msg.Denom); err != nil {
		return nil, err
	}

	// Verify denom exists
	_, exists, err := s.keeper.GetDenomAdmin(ctx, msg.Denom)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, types.ErrDenomNotFound
	}

	cap, ok := sdkmath.NewIntFromString(msg.NewCap)
	if !ok || !cap.IsPositive() {
		return nil, fmt.Errorf("new_cap must be a positive integer")
	}

	// Verify new cap is not lower than current supply
	currentSupply := s.keeper.bankKeeper.GetSupply(ctx, msg.Denom)
	if currentSupply.Amount.GT(cap) {
		return nil, types.ErrSupplyCapExceeded.Wrapf(
			"current supply %s > new cap %s", currentSupply.Amount.String(), cap.String(),
		)
	}

	if err := s.keeper.SetSupplyCap(ctx, msg.Denom, cap); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ModuleName,
		sdk.NewAttribute("action", "update_supply_cap"),
		sdk.NewAttribute("authority", msg.Authority),
		sdk.NewAttribute("denom", msg.Denom),
		sdk.NewAttribute("new_cap", msg.NewCap),
	))
	return &types.MsgUpdateSupplyCapResponse{}, nil
}

func validateFactoryDenom(denom string) error {
	if !strings.HasPrefix(denom, types.DenomPrefix+"/") {
		return types.ErrInvalidDenom
	}
	return nil
}
