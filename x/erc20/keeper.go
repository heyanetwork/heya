package erc20

import (
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"heya/x/erc20/types"
)

type Keeper struct {
	storeService store.KVStoreService
	bankKeeper   types.BankKeeper
}

func NewKeeper(storeService store.KVStoreService, bankKeeper types.BankKeeper) Keeper {
	return Keeper{
		storeService: storeService,
		bankKeeper:   bankKeeper,
	}
}

func (k Keeper) GetAuthority(ctx sdk.Context) string {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ParamsKey)
	if err != nil || bz == nil {
		return ""
	}
	return string(bz)
}

func (k Keeper) SetAuthority(ctx sdk.Context, authority string) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Set(types.ParamsKey, []byte(authority))
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) GetAllowance(ctx sdk.Context, owner, spender, denom string) sdkmath.Int {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.AllowanceKey(owner, spender, denom))
	if err != nil || bz == nil {
		return sdkmath.ZeroInt()
	}
	amt, ok := sdkmath.NewIntFromString(string(bz))
	if !ok {
		return sdkmath.ZeroInt()
	}
	return amt
}

func (k Keeper) SetAllowance(ctx sdk.Context, owner, spender, denom string, amount sdkmath.Int) error {
	store := k.storeService.OpenKVStore(ctx)
	if amount.IsZero() {
		return store.Delete(types.AllowanceKey(owner, spender, denom))
	}
	return store.Set(types.AllowanceKey(owner, spender, denom), []byte(amount.String()))
}
