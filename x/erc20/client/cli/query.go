package cli

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"heya/x/erc20/types"
)

func NewQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the ERC20 module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(
		NewBalanceOfCmd(),
		NewAllowanceCmd(),
		NewTotalSupplyCmd(),
		NewParamsCmd(),
	)
	return queryCmd
}

func NewBalanceOfCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "balance-of [owner] [denom]",
		Short:   "Get the balance of an address for a denom",
		Example: version.AppName + " query erc20 balance-of heya1abc factory/heya1abc/mytoken",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.BalanceOf(cmd.Context(), &types.QueryBalanceOfRequest{
				Owner: args[0],
				Denom: args[1],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewAllowanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "allowance [owner] [spender] [denom]",
		Short:   "Get the allowance of a spender for an owner",
		Example: version.AppName + " query erc20 allowance heya1owner heya1spender factory/heya1abc/mytoken",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Allowance(cmd.Context(), &types.QueryAllowanceRequest{
				Owner:   args[0],
				Spender: args[1],
				Denom:   args[2],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewTotalSupplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "total-supply [denom]",
		Short:   "Get the total supply of a denom",
		Example: version.AppName + " query erc20 total-supply factory/heya1abc/mytoken",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.TotalSupply(cmd.Context(), &types.QueryTotalSupplyRequest{
				Denom: args[0],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Short:   "Get erc20 module parameters",
		Args:    cobra.NoArgs,
		Example: "heyad query erc20 params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}
			bz, err := json.Marshal(resp)
			if err != nil {
				return err
			}
			return clientCtx.PrintString(string(bz) + "\n")
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
