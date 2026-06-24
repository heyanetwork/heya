package supplycap

import (
	"context"
	"encoding/json"
	"math/big"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
)

const (
	ModuleName = "supplycap"
)

// MaxSupply is 10,000,000,000 HEYA in uheya (10B * 1_000_000)
var MaxSupply = new(big.Int).Mul(big.NewInt(10_000_000_000), big.NewInt(1_000_000))

var (
	_ appmodule.AppModule       = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
	_ module.AppModuleBasic     = AppModule{}
)

type AppModule struct {
	bankKeeper bankkeeper.Keeper
	mintKeeper mintkeeper.Keeper
}

func NewAppModule(bk bankkeeper.Keeper, mk mintkeeper.Keeper) AppModule {
	return AppModule{bankKeeper: bk, mintKeeper: mk}
}

func (AppModule) IsOnePerModuleType() {}
func (AppModule) IsAppModule()        {}

func (am AppModule) BeginBlock(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	total := am.bankKeeper.GetSupply(sdkCtx, "uheya")
	if total.Amount.BigInt().Cmp(MaxSupply) >= 0 {
		params, err := am.mintKeeper.Params.Get(sdkCtx)
		if err != nil {
			return err
		}
		if !params.InflationMax.IsZero() {
			params.InflationMax = math.LegacyNewDec(0)
			params.InflationMin = math.LegacyNewDec(0)
			params.InflationRateChange = math.LegacyNewDec(0)
			if err := am.mintKeeper.Params.Set(sdkCtx, params); err != nil {
				return err
			}
			sdkCtx.Logger().Info("max supply 10B HEYA reached, inflation set to 0")
		}
	}
	return nil
}

func (AppModule) Name() string { return ModuleName }

func (AppModule) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}

func (AppModule) RegisterInterfaces(codectypes.InterfaceRegistry) {}

func (AppModule) RegisterGRPCGatewayRoutes(client.Context, *runtime.ServeMux) {}

func (AppModule) DefaultGenesis(_ codec.JSONCodec) json.RawMessage { return json.RawMessage("{}") }

func (AppModule) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, _ json.RawMessage) error {
	return nil
}
