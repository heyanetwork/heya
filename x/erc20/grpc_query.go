package erc20

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"heya/x/erc20/types"
)

var _ types.QueryServer = Querier{}

type Querier struct {
	types.UnimplementedQueryServer
	keeper Keeper
}

func NewQuerier(keeper Keeper) Querier {
	return Querier{keeper: keeper}
}

func (q Querier) BalanceOf(goCtx context.Context, req *types.QueryBalanceOfRequest) (*types.QueryBalanceOfResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, err
	}
	coin := q.keeper.bankKeeper.GetBalance(ctx, addr, req.Denom)
	return &types.QueryBalanceOfResponse{Amount: coin.Amount.String()}, nil
}

func (q Querier) Allowance(goCtx context.Context, req *types.QueryAllowanceRequest) (*types.QueryAllowanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	allowance := q.keeper.GetAllowance(ctx, req.Owner, req.Spender, req.Denom)
	return &types.QueryAllowanceResponse{Amount: allowance.String()}, nil
}

func (q Querier) TotalSupply(goCtx context.Context, req *types.QueryTotalSupplyRequest) (*types.QueryTotalSupplyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	coin := q.keeper.bankKeeper.GetSupply(ctx, req.Denom)
	return &types.QueryTotalSupplyResponse{Amount: coin.Amount.String()}, nil
}

func (q Querier) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	authority := q.keeper.GetAuthority(ctx)
	return &types.QueryParamsResponse{Authority: authority}, nil
}
