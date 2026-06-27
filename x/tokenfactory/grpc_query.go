package tokenfactory

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"heya/x/tokenfactory/types"
)

var _ types.QueryServer = Querier{}

type Querier struct {
	types.UnimplementedQueryServer
	keeper Keeper
}

func NewQuerier(keeper Keeper) Querier {
	return Querier{keeper: keeper}
}

func (q Querier) DenomAdmin(goCtx context.Context, req *types.QueryDenomAdminRequest) (*types.QueryDenomAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	admin, exists, err := q.keeper.GetDenomAdmin(ctx, req.Denom)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, types.ErrDenomNotFound
	}
	return &types.QueryDenomAdminResponse{Admin: admin}, nil
}

func (q Querier) Params(goCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params, err := q.keeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryParamsResponse{DenomCreationFee: params.DenomCreationFee.String()}, nil
}
