package erc20

import (
	"context"
	"encoding/json"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"heya/x/erc20/client/cli"
	"heya/x/erc20/types"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
	_ appmodule.AppModule   = AppModule{}
)

type AppModuleBasic struct{}

func (b AppModuleBasic) Name() string { return types.ModuleName }

func (b AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return json.RawMessage("{}")
}

func (b AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterRESTRoutes(_ client.Context, _ *runtime.ServeMux) {}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}

func (b AppModuleBasic) GetTxCmd() *cobra.Command { return cli.NewTxCmd() }

func (b AppModuleBasic) GetQueryCmd() *cobra.Command { return cli.NewQueryCmd() }

func (b AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

type AppModule struct {
	AppModuleBasic
	keeper    Keeper
	cdc       codec.Codec
	authority string
}

func NewAppModule(cdc codec.Codec, keeper Keeper, authority string) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		cdc:            cdc,
		authority:      authority,
	}
}

func (am AppModule) IsOnePerModuleType() {}

func (am AppModule) IsAppModule() {}

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), NewQuerier(am.keeper))
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	if am.keeper.GetAuthority(ctx) == "" && am.authority != "" {
		_ = am.keeper.SetAuthority(ctx, am.authority)
	}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return json.RawMessage("{}")
}

func (am AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) BeginBlock(context.Context) error { return nil }

func (am AppModule) EndBlock(context.Context) error { return nil }

func (am AppModule) GenerateGenesisState(simState *module.SimulationState) {}

func (am AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}

func (am AppModule) RegisterStoreKey(_ store.KVStoreService) {}

var _ appmodule.HasBeginBlocker = AppModule{}
var _ appmodule.HasEndBlocker = AppModule{}
