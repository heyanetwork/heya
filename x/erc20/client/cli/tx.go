package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	"heya/x/erc20/types"
)

func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "ERC20 token transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		NewTransferCmd(),
		NewApproveCmd(),
		NewTransferFromCmd(),
		NewIncreaseAllowanceCmd(),
		NewDecreaseAllowanceCmd(),
		NewUpdateParamsCmd(),
	)
	return txCmd
}

func NewTransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transfer [to] [amount]",
		Short:   "Transfer tokens (ERC20 transfer)",
		Long:    "Transfer tokens to an address. Amount must include the full denom (e.g. 1000factory/heya1.../mytoken).",
		Example: version.AppName + " tx erc20 transfer heya1dest 1000factory/heya1abc/mytoken",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgTransfer(clientCtx.GetFromAddress().String(), args[0], args[1])
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewApproveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "approve [spender] [amount]",
		Short:   "Approve spender to transfer tokens (ERC20 approve)",
		Long:    "Approve a spender to transfer up to the given amount of tokens on your behalf. Amount must include the full denom.",
		Example: version.AppName + " tx erc20 approve heya1spender 1000factory/heya1abc/mytoken",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgApprove(clientCtx.GetFromAddress().String(), args[0], args[1])
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewTransferFromCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transfer-from [from] [to] [amount]",
		Short:   "Transfer tokens from owner to recipient using allowance (ERC20 transferFrom)",
		Long:    "Transfer tokens from an owner to a recipient using the caller's allowance. Caller must have sufficient allowance.",
		Example: version.AppName + " tx erc20 transfer-from heya1owner heya1dest 1000factory/heya1abc/mytoken",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgTransferFrom(clientCtx.GetFromAddress().String(), args[0], args[1], args[2])
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewIncreaseAllowanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "increase-allowance [spender] [denom] [amount]",
		Short:   "Increase allowance for a spender by a given amount",
		Long:    "Safely increase the allowance for a spender by the given amount. Unlike approve, this does not have the race condition.",
		Example: version.AppName + " tx erc20 increase-allowance heya1spender factory/heya1abc/mytoken 100",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgIncreaseAllowance(clientCtx.GetFromAddress().String(), args[0], args[1], args[2])
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewDecreaseAllowanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "decrease-allowance [spender] [denom] [amount]",
		Short:   "Decrease allowance for a spender by a given amount",
		Long:    "Safely decrease the allowance for a spender by the given amount. Reverts if allowance would go negative.",
		Example: version.AppName + " tx erc20 decrease-allowance heya1spender factory/heya1abc/mytoken 50",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgDecreaseAllowance(clientCtx.GetFromAddress().String(), args[0], args[1], args[2])
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewUpdateParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update-params [new-authority]",
		Short:   "Update erc20 module authority via governance proposal",
		Long:    "Update erc20 module authority. The current authority (governance module) must sign. Use 'heyad tx gov submit-proposal' for standard governance flow.",
		Example: "heyad tx erc20 update-params heya1abc...",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgUpdateParams(clientCtx.GetFromAddress().String(), args[0])
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
