package tokenfactory

import (
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"heya/x/tokenfactory/types"
)

type Keeper struct {
	storeService store.KVStoreService
	bankKeeper   types.BankKeeper
}

func NewKeeper(
	storeService store.KVStoreService,
	bankKeeper types.BankKeeper,
) Keeper {
	return Keeper{
		storeService: storeService,
		bankKeeper:   bankKeeper,
	}
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
