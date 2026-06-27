package tokenfactory

import (
	"context"
	"testing"

	"cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/stretchr/testify/require"

	"heya/x/tokenfactory/types"
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

type mockBankKeeper struct {
	balances map[string]sdk.Coins
	supply   map[string]sdk.Coin
}

func newMockBankKeeper() *mockBankKeeper {
	return &mockBankKeeper{
		balances: make(map[string]sdk.Coins),
		supply:   make(map[string]sdk.Coin),
	}
}

func (m *mockBankKeeper) MintCoins(ctx context.Context, moduleAddr string, amt sdk.Coins) error {
	m.balances[moduleAddr] = m.balances[moduleAddr].Add(amt...)
	for _, coin := range amt {
		if existing, ok := m.supply[coin.Denom]; ok {
			m.supply[coin.Denom] = existing.Add(coin)
		} else {
			m.supply[coin.Denom] = coin
		}
	}
	return nil
}

func (m *mockBankKeeper) BurnCoins(ctx context.Context, moduleAddr string, amt sdk.Coins) error {
	m.balances[moduleAddr] = m.balances[moduleAddr].Sub(amt...)
	for _, coin := range amt {
		if existing, ok := m.supply[coin.Denom]; ok {
			m.supply[coin.Denom] = existing.Sub(coin)
		}
	}
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	m.balances[senderModule] = m.balances[senderModule].Sub(amt...)
	m.balances[recipientAddr.String()] = m.balances[recipientAddr.String()].Add(amt...)
	return nil
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	m.balances[senderAddr.String()] = m.balances[senderAddr.String()].Sub(amt...)
	m.balances[recipientModule] = m.balances[recipientModule].Add(amt...)
	return nil
}

func (m *mockBankKeeper) SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	m.balances[fromAddr.String()] = m.balances[fromAddr.String()].Sub(amt...)
	m.balances[toAddr.String()] = m.balances[toAddr.String()].Add(amt...)
	return nil
}

func (m *mockBankKeeper) GetSupply(ctx context.Context, denom string) sdk.Coin {
	if coin, ok := m.supply[denom]; ok {
		return coin
	}
	return sdk.NewCoin(denom, sdkmath.ZeroInt())
}

var validSender string

func init() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("heya", "heyapub")

	addr := sdk.AccAddress(make([]byte, 20))
	validSender, _ = bech32.ConvertAndEncode("heya", addr)
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
	bech, err := bech32.ConvertAndEncode("heya", addr)
	require.NoError(t, err)
	return bech
}

func setupTest(t *testing.T) (Keeper, *mockBankKeeper, types.MsgServer) {
	t.Helper()

	mockStoreService := newMockKVStoreService()
	mockBank := newMockBankKeeper()

	mockBank.balances[validSender] = sdk.NewCoins(sdk.NewCoin("uheya", sdkmath.NewInt(100_000_000_000)))

	keeper := NewKeeper(mockStoreService, mockBank, validSender)
	server := NewMsgServerImpl(keeper)

	return keeper, mockBank, server
}

func createDenom(t *testing.T, server types.MsgServer, ctx sdk.Context, sender, subdenom string) string {
	t.Helper()
	createMsg := types.NewMsgCreateDenom(sender, subdenom)
	resp, err := server.CreateDenom(ctx, createMsg)
	require.NoError(t, err)
	return resp.Denom
}

func TestCreateDenom(t *testing.T) {
	keeper, _, server := setupTest(t)
	ctx := testCtx()

	denom := createDenom(t, server, ctx, validSender, "mytoken")
	require.Equal(t, "factory/"+validSender+"/mytoken", denom)

	admin, exists, err := keeper.GetDenomAdmin(ctx, denom)
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, validSender, admin)

	cap, hasCap, err := keeper.GetSupplyCap(ctx, denom)
	require.NoError(t, err)
	require.True(t, hasCap)
	require.True(t, cap.Equal(sdkmath.NewInt(1_000_000_000_000_000)))

	_, err = server.CreateDenom(ctx, types.NewMsgCreateDenom(validSender, "mytoken"))
	require.ErrorIs(t, err, types.ErrDenomExists)
}

func TestMint(t *testing.T) {
	_, _, server := setupTest(t)
	ctx := testCtx()

	denom := createDenom(t, server, ctx, validSender, "mytoken")

	mintMsg := types.NewMsgMint(validSender, "1000"+denom, "")
	_, err := server.Mint(ctx, mintMsg)
	require.NoError(t, err)
}

func TestMintExceedsSupplyCap(t *testing.T) {
	keeper, _, server := setupTest(t)
	ctx := testCtx()

	denom := createDenom(t, server, ctx, validSender, "limited")

	err := keeper.SetSupplyCap(ctx, denom, sdkmath.NewInt(5000))
	require.NoError(t, err)

	_, err = server.Mint(ctx, types.NewMsgMint(validSender, "3000"+denom, ""))
	require.NoError(t, err)

	_, err = server.Mint(ctx, types.NewMsgMint(validSender, "3000"+denom, ""))
	require.ErrorIs(t, err, types.ErrSupplyCapExceeded)
}

func TestMintUnauthorized(t *testing.T) {
	_, _, server := setupTest(t)
	ctx := testCtx()

	otherSender := makeAddr(t, 1)
	denom := createDenom(t, server, ctx, validSender, "mytoken")

	_, err := server.Mint(ctx, types.NewMsgMint(otherSender, "1000"+denom, ""))
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestBurn(t *testing.T) {
	_, _, server := setupTest(t)
	ctx := testCtx()

	denom := createDenom(t, server, ctx, validSender, "mytoken")
	_, err := server.Mint(ctx, types.NewMsgMint(validSender, "1000"+denom, ""))
	require.NoError(t, err)

	_, err = server.Burn(ctx, types.NewMsgBurn(validSender, "500"+denom))
	require.NoError(t, err)
}

func TestTwoStepChangeAdmin(t *testing.T) {
	keeper, _, server := setupTest(t)
	ctx := testCtx()

	newAdmin := makeAddr(t, 2)
	denom := createDenom(t, server, ctx, validSender, "mytoken")

	// Step 1: Current admin proposes new admin
	_, err := server.ChangeAdmin(ctx, types.NewMsgChangeAdmin(validSender, denom, newAdmin))
	require.NoError(t, err)

	// Admin should NOT have changed yet
	admin, _, err := keeper.GetDenomAdmin(ctx, denom)
	require.NoError(t, err)
	require.Equal(t, validSender, admin)

	// Pending admin should be set
	pending, hasPending, err := keeper.GetPendingAdmin(ctx, denom)
	require.NoError(t, err)
	require.True(t, hasPending)
	require.Equal(t, newAdmin, pending)

	// Step 2: Wrong address tries to accept
	wrongAddr := makeAddr(t, 3)
	_, err = server.AcceptAdmin(ctx, types.NewMsgAcceptAdmin(wrongAddr, denom))
	require.ErrorIs(t, err, types.ErrNotPendingAdmin)

	// Step 3: Correct pending admin accepts
	_, err = server.AcceptAdmin(ctx, types.NewMsgAcceptAdmin(newAdmin, denom))
	require.NoError(t, err)

	// Admin should now be the new admin
	admin, _, err = keeper.GetDenomAdmin(ctx, denom)
	require.NoError(t, err)
	require.Equal(t, newAdmin, admin)

	// Pending should be cleared
	_, hasPending, err = keeper.GetPendingAdmin(ctx, denom)
	require.NoError(t, err)
	require.False(t, hasPending)
}

func TestAcceptAdminNoPending(t *testing.T) {
	_, _, server := setupTest(t)
	ctx := testCtx()

	denom := createDenom(t, server, ctx, validSender, "mytoken")
	newAdmin := makeAddr(t, 2)

	_, err := server.AcceptAdmin(ctx, types.NewMsgAcceptAdmin(newAdmin, denom))
	require.ErrorIs(t, err, types.ErrPendingAdminNotFound)
}

func TestChangeAdminUnauthorized(t *testing.T) {
	_, _, server := setupTest(t)
	ctx := testCtx()

	otherSender := makeAddr(t, 1)
	denom := createDenom(t, server, ctx, validSender, "mytoken")

	_, err := server.ChangeAdmin(ctx, types.NewMsgChangeAdmin(otherSender, denom, otherSender))
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestForceTransfer(t *testing.T) {
	_, bank, server := setupTest(t)
	ctx := testCtx()

	destAddr := makeAddr(t, 2)
	denom := createDenom(t, server, ctx, validSender, "mytoken")

	_, err := server.Mint(ctx, types.NewMsgMint(validSender, "1000"+denom, ""))
	require.NoError(t, err)

	// Force transfer from admin's own account (backward compat)
	_, err = server.ForceTransfer(ctx, types.NewMsgForceTransferFull(validSender, "300"+denom, destAddr, ""))
	require.NoError(t, err)

	// Verify balances
	senderAddr, _ := sdk.AccAddressFromBech32(validSender)
	destAcc, _ := sdk.AccAddressFromBech32(destAddr)
	require.True(t, bank.balances[destAcc.String()].AmountOf(denom).Equal(sdkmath.NewInt(300)))
	require.True(t, bank.balances[senderAddr.String()].AmountOf(denom).Equal(sdkmath.NewInt(700)))
}

func TestForceTransferFromAnyAddress(t *testing.T) {
	_, bank, server := setupTest(t)
	ctx := testCtx()

	holder := makeAddr(t, 2)
	destAddr := makeAddr(t, 3)
	denom := createDenom(t, server, ctx, validSender, "mytoken")

	// Mint tokens to holder
	_, err := server.Mint(ctx, types.NewMsgMint(validSender, "1000"+denom, holder))
	require.NoError(t, err)

	// Force transfer from holder to dest (admin authorizes, tokens come from holder)
	_, err = server.ForceTransfer(ctx, types.NewMsgForceTransferFull(validSender, "400"+denom, destAddr, holder))
	require.NoError(t, err)

	holderAcc, _ := sdk.AccAddressFromBech32(holder)
	destAcc, _ := sdk.AccAddressFromBech32(destAddr)
	require.True(t, bank.balances[holderAcc.String()].AmountOf(denom).Equal(sdkmath.NewInt(600)))
	require.True(t, bank.balances[destAcc.String()].AmountOf(denom).Equal(sdkmath.NewInt(400)))
}

func TestPausedBlocksOperations(t *testing.T) {
	keeper, _, server := setupTest(t)
	ctx := testCtx()

	err := keeper.SetPaused(ctx, true)
	require.NoError(t, err)

	_, err = server.CreateDenom(ctx, types.NewMsgCreateDenom(validSender, "mytoken"))
	require.ErrorIs(t, err, types.ErrPaused)

	_, err = server.Mint(ctx, types.NewMsgMint(validSender, "1000factory/heya1abc/mytoken", ""))
	require.ErrorIs(t, err, types.ErrPaused)

	_, err = server.Burn(ctx, types.NewMsgBurn(validSender, "1000factory/heya1abc/mytoken"))
	require.ErrorIs(t, err, types.ErrPaused)

	_, err = server.ChangeAdmin(ctx, types.NewMsgChangeAdmin(validSender, "factory/heya1abc/mytoken", validSender))
	require.ErrorIs(t, err, types.ErrPaused)

	_, err = server.AcceptAdmin(ctx, types.NewMsgAcceptAdmin(validSender, "factory/heya1abc/mytoken"))
	require.ErrorIs(t, err, types.ErrPaused)

	_, err = server.ForceTransfer(ctx, types.NewMsgForceTransfer(validSender, "1000factory/heya1abc/mytoken", validSender))
	require.ErrorIs(t, err, types.ErrPaused)

	err = keeper.SetPaused(ctx, false)
	require.NoError(t, err)

	_, err = server.CreateDenom(ctx, types.NewMsgCreateDenom(validSender, "mytoken2"))
	require.NoError(t, err)
}

func TestPausedState(t *testing.T) {
	keeper, _, _ := setupTest(t)
	ctx := testCtx()

	require.False(t, keeper.IsPaused(ctx))

	err := keeper.SetPaused(ctx, true)
	require.NoError(t, err)
	require.True(t, keeper.IsPaused(ctx))

	err = keeper.SetPaused(ctx, false)
	require.NoError(t, err)
	require.False(t, keeper.IsPaused(ctx))
}

func TestOnChainParams(t *testing.T) {
	keeper, _, server := setupTest(t)
	ctx := testCtx()

	// Default params should work
	params, err := keeper.GetParams(ctx)
	require.NoError(t, err)
	require.True(t, params.DenomCreationFee.Equal(sdk.NewCoin("uheya", sdkmath.NewInt(1_000_000_000))))

	// Change creation fee
	newParams := types.Params{
		DenomCreationFee: sdk.NewCoin("uheya", sdkmath.NewInt(100_000_000)),
	}
	err = keeper.SetParams(ctx, newParams)
	require.NoError(t, err)

	// Verify params changed
	params, err = keeper.GetParams(ctx)
	require.NoError(t, err)
	require.True(t, params.DenomCreationFee.Equal(sdk.NewCoin("uheya", sdkmath.NewInt(100_000_000))))

	// Create denom should use new fee (mock bank only has 100B uheya, should have enough)
	denom := createDenom(t, server, ctx, validSender, "newfee")
	require.Contains(t, denom, "newfee")
}

func TestParamsQuery(t *testing.T) {
	keeper, _, _ := setupTest(t)
	ctx := testCtx()

	expectedFee := sdk.NewCoin("uheya", sdkmath.NewInt(1_000_000_000))
	err := keeper.SetParams(ctx, types.Params{DenomCreationFee: expectedFee})
	require.NoError(t, err)

	params, err := keeper.GetParams(ctx)
	require.NoError(t, err)
	require.True(t, params.DenomCreationFee.Equal(expectedFee))
}

func TestUpdateParams(t *testing.T) {
	keeper, _, server := setupTest(t)
	ctx := testCtx()

	// Valid update by authority should succeed
	msg := types.NewMsgUpdateParams(validSender, "500000000uheya")
	_, err := server.UpdateParams(ctx, msg)
	require.NoError(t, err)

	// Verify params changed
	params, err := keeper.GetParams(ctx)
	require.NoError(t, err)
	require.True(t, params.DenomCreationFee.Equal(sdk.NewCoin("uheya", sdkmath.NewInt(500_000_000))))
}

func TestUpdateParamsUnauthorized(t *testing.T) {
	_, _, server := setupTest(t)
	ctx := testCtx()

	// Update by non-authority should fail
	attacker := makeAddr(t, 99)
	msg := types.NewMsgUpdateParams(attacker, "500000000uheya")
	_, err := server.UpdateParams(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unauthorized")
}

func TestUpdateParamsInvalidFee(t *testing.T) {
	_, _, server := setupTest(t)
	ctx := testCtx()

	// Invalid fee should fail validation
	msg := types.NewMsgUpdateParams(validSender, "invalid")
	_, err := server.UpdateParams(ctx, msg)
	require.Error(t, err)
}

func TestSupplyCapKeeper(t *testing.T) {
	keeper, _, _ := setupTest(t)
	ctx := testCtx()

	denom := "factory/heya1abc/test"

	// No cap initially
	_, has, err := keeper.GetSupplyCap(ctx, denom)
	require.NoError(t, err)
	require.False(t, has)

	// Set cap
	err = keeper.SetSupplyCap(ctx, denom, sdkmath.NewInt(1000000))
	require.NoError(t, err)

	cap, has, err := keeper.GetSupplyCap(ctx, denom)
	require.NoError(t, err)
	require.True(t, has)
	require.True(t, cap.Equal(sdkmath.NewInt(1000000)))
}
