package tokenfactory

import (
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"heya/x/tokenfactory/types"
)

type Keeper struct {
	storeService store.KVStoreService
	bankKeeper   types.BankKeeper
	authority    string
}

func NewKeeper(
	storeService store.KVStoreService,
	bankKeeper types.BankKeeper,
	authority string,
) Keeper {
	return Keeper{
		storeService: storeService,
		bankKeeper:   bankKeeper,
		authority:    authority,
	}
}

func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) GetDenomAdmin(ctx sdk.Context, denom string) (string, bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.DenomKey(denom))
	if err != nil {
		return "", false, err
	}
	if bz == nil {
		return "", false, nil
	}
	return string(bz), true, nil
}

func (k Keeper) SetDenomAdmin(ctx sdk.Context, denom, admin string) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Set(types.DenomKey(denom), []byte(admin))
}

func (k Keeper) AllDenomAdmins(ctx sdk.Context) ([]*types.DenomAuthority, error) {
	store := k.storeService.OpenKVStore(ctx)
	iter, err := store.Iterator(types.DenomKeyPrefix, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var denoms []*types.DenomAuthority
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		value := iter.Value()
		denoms = append(denoms, &types.DenomAuthority{
			Denom: string(key[len(types.DenomKeyPrefix):]),
			Admin: string(value),
		})
	}
	return denoms, nil
}

func (k Keeper) GetSupplyCap(ctx sdk.Context, denom string) (sdkmath.Int, bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.SupplyCapKey(denom))
	if err != nil {
		return sdkmath.Int{}, false, err
	}
	if bz == nil {
		return sdkmath.Int{}, false, nil
	}
	cap, ok := sdkmath.NewIntFromString(string(bz))
	if !ok {
		return sdkmath.Int{}, false, types.ErrInvalidDenom.Wrap("invalid supply cap value")
	}
	return cap, true, nil
}

func (k Keeper) SetSupplyCap(ctx sdk.Context, denom string, cap sdkmath.Int) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Set(types.SupplyCapKey(denom), []byte(cap.String()))
}

func (k Keeper) IsPaused(ctx sdk.Context) bool {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.PausedKey)
	if err != nil || bz == nil {
		return false
	}
	return string(bz) == "true"
}

func (k Keeper) SetPaused(ctx sdk.Context, paused bool) error {
	store := k.storeService.OpenKVStore(ctx)
	if paused {
		return store.Set(types.PausedKey, []byte("true"))
	}
	return store.Set(types.PausedKey, []byte("false"))
}

func (k Keeper) GetPendingAdmin(ctx sdk.Context, denom string) (string, bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.PendingAdminKey(denom))
	if err != nil {
		return "", false, err
	}
	if bz == nil {
		return "", false, nil
	}
	return string(bz), true, nil
}

func (k Keeper) SetPendingAdmin(ctx sdk.Context, denom, admin string) error {
	store := k.storeService.OpenKVStore(ctx)
	if admin == "" {
		return store.Delete(types.PendingAdminKey(denom))
	}
	return store.Set(types.PendingAdminKey(denom), []byte(admin))
}

func (k Keeper) GetParams(ctx sdk.Context) (types.Params, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ParamsKey)
	if err != nil || bz == nil {
		return types.DefaultParams(), nil
	}
	coin, err := sdk.ParseCoinNormalized(string(bz))
	if err != nil {
		return types.DefaultParams(), nil
	}
	return types.Params{DenomCreationFee: coin}, nil
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Set(types.ParamsKey, []byte(params.DenomCreationFee.String()))
}
