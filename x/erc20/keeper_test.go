package erc20

import (
	"context"
	"testing"

	"cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"heya/x/erc20/types"
)

type mockKVStore struct {
	data map[string][]byte
}

func newMockKVStore() *mockKVStore {
	return &mockKVStore{data: make(map[string][]byte)}
}

func (m *mockKVStore) Get(key []byte) ([]byte, error) {
	val, ok := m.data[string(key)]
	if !ok {
		return nil, nil
	}
	return val, nil
}

func (m *mockKVStore) Has(key []byte) (bool, error) {
	_, ok := m.data[string(key)]
	return ok, nil
}

func (m *mockKVStore) Set(key, value []byte) error {
	m.data[string(key)] = value
	return nil
}

func (m *mockKVStore) Delete(key []byte) error {
	delete(m.data, string(key))
	return nil
}

func (m *mockKVStore) Iterator(start, end []byte) (store.Iterator, error) {
	return nil, nil
}

func (m *mockKVStore) ReverseIterator(start, end []byte) (store.Iterator, error) {
	return nil, nil
}

type mockKVStoreService struct {
	store *mockKVStore
}

func newMockKVStoreService() *mockKVStoreService {
	return &mockKVStoreService{store: newMockKVStore()}
}

func (m *mockKVStoreService) OpenKVStore(ctx context.Context) store.KVStore {
	return m.store
}

type mockBankKeeper struct{}

func newMockBankKeeper() *mockBankKeeper {
	return &mockBankKeeper{}
}

func (m *mockBankKeeper) SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

func (m *mockBankKeeper) GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return sdk.NewCoin(denom, sdkmath.NewInt(1000))
}

func (m *mockBankKeeper) GetSupply(ctx context.Context, denom string) sdk.Coin {
	return sdk.NewCoin(denom, sdkmath.NewInt(1000))
}

func testCtx() sdk.Context {
	return sdk.Context{}.
		WithContext(context.Background()).
		WithEventManager(sdk.NewEventManager())
}

func makeAddr(t *testing.T, seed byte) string {
	t.Helper()
	addr := sdk.AccAddress(make([]byte, 20))
	addr[0] = seed
	return addr.String()
}

func setupKeeper(t *testing.T) (Keeper, sdk.Context) {
	t.Helper()
	storeSvc := newMockKVStoreService()
	bankKeeper := newMockBankKeeper()
	keeper := NewKeeper(storeSvc, bankKeeper)
	return keeper, testCtx()
}

func setupMsgServer(t *testing.T) (types.MsgServer, Keeper, sdk.Context) {
	t.Helper()
	keeper, ctx := setupKeeper(t)
	server := NewMsgServerImpl(keeper)
	return server, keeper, ctx
}

func TestGetSetAllowance(t *testing.T) {
	keeper, ctx := setupKeeper(t)

	amt := keeper.GetAllowance(ctx, "owner", "spender", "denom")
	require.True(t, amt.IsZero())

	err := keeper.SetAllowance(ctx, "owner", "spender", "denom", sdkmath.NewInt(100))
	require.NoError(t, err)

	amt = keeper.GetAllowance(ctx, "owner", "spender", "denom")
	require.Equal(t, sdkmath.NewInt(100), amt)

	err = keeper.SetAllowance(ctx, "owner", "spender", "denom", sdkmath.NewInt(50))
	require.NoError(t, err)

	amt = keeper.GetAllowance(ctx, "owner", "spender", "denom")
	require.Equal(t, sdkmath.NewInt(50), amt)
}

func TestDeleteAllowanceOnZero(t *testing.T) {
	keeper, ctx := setupKeeper(t)

	err := keeper.SetAllowance(ctx, "owner", "spender", "denom", sdkmath.NewInt(100))
	require.NoError(t, err)

	err = keeper.SetAllowance(ctx, "owner", "spender", "denom", sdkmath.ZeroInt())
	require.NoError(t, err)

	amt := keeper.GetAllowance(ctx, "owner", "spender", "denom")
	require.True(t, amt.IsZero())
}

func TestGetSetAuthority(t *testing.T) {
	keeper, ctx := setupKeeper(t)

	auth := keeper.GetAuthority(ctx)
	require.Empty(t, auth)

	err := keeper.SetAuthority(ctx, "heya1authority")
	require.NoError(t, err)

	auth = keeper.GetAuthority(ctx)
	require.Equal(t, "heya1authority", auth)
}

func TestTransfer(t *testing.T) {
	server, _, ctx := setupMsgServer(t)

	addr1 := makeAddr(t, 1)
	addr2 := makeAddr(t, 2)

	_, err := server.Transfer(ctx, types.NewMsgTransfer(addr1, addr2, "100factory/creator/mytoken"))
	require.NoError(t, err)
}

func TestTransferInvalidAmount(t *testing.T) {
	server, _, ctx := setupMsgServer(t)

	addr1 := makeAddr(t, 1)
	addr2 := makeAddr(t, 2)

	_, err := server.Transfer(ctx, types.NewMsgTransfer(addr1, addr2, "notanumber"))
	require.Error(t, err)
}

func TestApprove(t *testing.T) {
	server, keeper, ctx := setupMsgServer(t)

	addr1 := makeAddr(t, 1)
	addr2 := makeAddr(t, 2)

	_, err := server.Approve(ctx, types.NewMsgApprove(addr1, addr2, "100factory/creator/mytoken"))
	require.NoError(t, err)

	amt := keeper.GetAllowance(ctx, addr1, addr2, "factory/creator/mytoken")
	require.Equal(t, sdkmath.NewInt(100), amt)
}

func TestTransferFromInsufficientAllowance(t *testing.T) {
	server, _, ctx := setupMsgServer(t)

	addr1 := makeAddr(t, 1)
	addr2 := makeAddr(t, 2)
	addr3 := makeAddr(t, 3)

	_, err := server.TransferFrom(ctx, types.NewMsgTransferFrom(addr2, addr1, addr3, "100factory/creator/mytoken"))
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient allowance")
}

func TestTransferFromWithAllowance(t *testing.T) {
	server, keeper, ctx := setupMsgServer(t)

	addr1 := makeAddr(t, 1)
	addr2 := makeAddr(t, 2)
	addr3 := makeAddr(t, 3)

	err := keeper.SetAllowance(ctx, addr1, addr2, "factory/creator/mytoken", sdkmath.NewInt(200))
	require.NoError(t, err)

	_, err = server.TransferFrom(ctx, types.NewMsgTransferFrom(addr2, addr1, addr3, "100factory/creator/mytoken"))
	require.NoError(t, err)

	amt := keeper.GetAllowance(ctx, addr1, addr2, "factory/creator/mytoken")
	require.Equal(t, sdkmath.NewInt(100), amt)
}

func TestIncreaseAllowance(t *testing.T) {
	server, keeper, ctx := setupMsgServer(t)

	addr1 := makeAddr(t, 1)
	addr2 := makeAddr(t, 2)

	err := keeper.SetAllowance(ctx, addr1, addr2, "denom", sdkmath.NewInt(100))
	require.NoError(t, err)

	_, err = server.IncreaseAllowance(ctx, types.NewMsgIncreaseAllowance(addr1, addr2, "denom", "50"))
	require.NoError(t, err)

	amt := keeper.GetAllowance(ctx, addr1, addr2, "denom")
	require.Equal(t, sdkmath.NewInt(150), amt)
}

func TestIncreaseAllowanceFromZero(t *testing.T) {
	server, keeper, ctx := setupMsgServer(t)

	addr1 := makeAddr(t, 1)
	addr2 := makeAddr(t, 2)

	_, err := server.IncreaseAllowance(ctx, types.NewMsgIncreaseAllowance(addr1, addr2, "denom", "75"))
	require.NoError(t, err)

	amt := keeper.GetAllowance(ctx, addr1, addr2, "denom")
	require.Equal(t, sdkmath.NewInt(75), amt)
}

func TestDecreaseAllowance(t *testing.T) {
	server, keeper, ctx := setupMsgServer(t)

	addr1 := makeAddr(t, 1)
	addr2 := makeAddr(t, 2)

	err := keeper.SetAllowance(ctx, addr1, addr2, "denom", sdkmath.NewInt(100))
	require.NoError(t, err)

	_, err = server.DecreaseAllowance(ctx, types.NewMsgDecreaseAllowance(addr1, addr2, "denom", "30"))
	require.NoError(t, err)

	amt := keeper.GetAllowance(ctx, addr1, addr2, "denom")
	require.Equal(t, sdkmath.NewInt(70), amt)
}

func TestDecreaseAllowanceInsufficient(t *testing.T) {
	server, _, ctx := setupMsgServer(t)

	addr1 := makeAddr(t, 1)
	addr2 := makeAddr(t, 2)

	_, err := server.DecreaseAllowance(ctx, types.NewMsgDecreaseAllowance(addr1, addr2, "denom", "10"))
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient allowance")
}

func TestDecreaseAllowanceMoreThanBalance(t *testing.T) {
	server, keeper, ctx := setupMsgServer(t)

	addr1 := makeAddr(t, 1)
	addr2 := makeAddr(t, 2)

	err := keeper.SetAllowance(ctx, addr1, addr2, "denom", sdkmath.NewInt(5))
	require.NoError(t, err)

	_, err = server.DecreaseAllowance(ctx, types.NewMsgDecreaseAllowance(addr1, addr2, "denom", "10"))
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient allowance")
}

func TestUpdateParams(t *testing.T) {
	server, keeper, ctx := setupMsgServer(t)

	addr1 := makeAddr(t, 1)
	addr2 := makeAddr(t, 2)

	err := keeper.SetAuthority(ctx, addr1)
	require.NoError(t, err)

	_, err = server.UpdateParams(ctx, types.NewMsgUpdateParams(addr1, addr2))
	require.NoError(t, err)

	auth := keeper.GetAuthority(ctx)
	require.Equal(t, addr2, auth)
}

func TestUpdateParamsUnauthorized(t *testing.T) {
	server, keeper, ctx := setupMsgServer(t)

	addr1 := makeAddr(t, 1)
	addr2 := makeAddr(t, 2)
	addr3 := makeAddr(t, 3)

	err := keeper.SetAuthority(ctx, addr1)
	require.NoError(t, err)

	_, err = server.UpdateParams(ctx, types.NewMsgUpdateParams(addr2, addr3))
	require.Error(t, err)
	require.ErrorContains(t, err, "unauthorized")
}

func TestParamsQuery(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	querier := NewQuerier(keeper)
	resp, err := querier.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Empty(t, resp.Authority)
}

func TestBalanceOfQuery(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	querier := NewQuerier(keeper)
	resp, err := querier.BalanceOf(ctx, &types.QueryBalanceOfRequest{
		Owner: makeAddr(t, 1),
		Denom: "factory/creator/mytoken",
	})
	require.NoError(t, err)
	require.Equal(t, "1000", resp.Amount)
}

func TestAllowanceQuery(t *testing.T) {
	keeper, ctx := setupKeeper(t)

	err := keeper.SetAllowance(ctx, "owner", "spender", "denom", sdkmath.NewInt(50))
	require.NoError(t, err)

	querier := NewQuerier(keeper)
	resp, err := querier.Allowance(ctx, &types.QueryAllowanceRequest{
		Owner:   "owner",
		Spender: "spender",
		Denom:   "denom",
	})
	require.NoError(t, err)
	require.Equal(t, "50", resp.Amount)
}

func TestTotalSupplyQuery(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	querier := NewQuerier(keeper)
	resp, err := querier.TotalSupply(ctx, &types.QueryTotalSupplyRequest{
		Denom: "factory/creator/mytoken",
	})
	require.NoError(t, err)
	require.Equal(t, "1000", resp.Amount)
}
