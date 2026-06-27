package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	"heya/x/tokenfactory/types"
)

func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Token factory transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		NewCreateDenomCmd(),
		NewMintCmd(),
		NewBurnCmd(),
		NewChangeAdminCmd(),
		NewAcceptAdminCmd(),
		NewForceTransferCmd(),
		NewUpdateParamsCmd(),
		NewUpdateSupplyCapCmd(),
	)
	return txCmd
}

func NewUpdateParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update-params [denom-creation-fee]",
		Short:   "Update token factory params via governance proposal",
		Long:    "Update token factory params. The authority (governance module) must sign. Use 'heyad tx gov submit-proposal' for standard governance flow.",
		Example: version.AppName + " tx tokenfactory update-params 500000000uheya",
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

func NewUpdateSupplyCapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update-supply-cap [denom] [new-cap]",
		Short:   "Update supply cap for a denom via governance proposal",
		Long:    "Update the supply cap for a denom. The authority (governance module) must sign. The new cap must be >= current supply. Use 'heyad tx gov submit-proposal' for standard governance flow.",
		Example: version.AppName + " tx tokenfactory update-supply-cap factory/heya1abc123/mytoken 2000000000000000",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgUpdateSupplyCap(clientCtx.GetFromAddress().String(), args[0], args[1])
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewCreateDenomCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-denom [subdenom]",
		Short:   "Create a new denom",
		Long:    "Create a new denom. The sender becomes the admin. Denom format: factory/{creator}/{subdenom}",
		Example: version.AppName + " tx tokenfactory create-denom mytoken",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgCreateDenom(clientCtx.GetFromAddress().String(), args[0])
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewMintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mint [amount]",
		Short:   "Mint tokens for a denom you are admin of",
		Long:    "Mint tokens for a denom. The amount must include the full denom (e.g. 1000factory/heya1.../mytoken). Use --mint-to to send to a different address.",
		Example: version.AppName + " tx tokenfactory mint 1000000factory/heya1abc123/mytoken",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			mintTo, _ := cmd.Flags().GetString("mint-to")
			msg := types.NewMsgMint(clientCtx.GetFromAddress().String(), args[0], mintTo)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String("mint-to", "", "address to mint tokens to (defaults to sender)")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewBurnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "burn [amount]",
		Short:   "Burn tokens for a denom you are admin of",
		Long:    "Burn tokens. The amount must include the full denom (e.g. 1000factory/heya1.../mytoken).",
		Example: version.AppName + " tx tokenfactory burn 1000000factory/heya1abc123/mytoken",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgBurn(clientCtx.GetFromAddress().String(), args[0])
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewChangeAdminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "change-admin [denom] [new-admin]",
		Short:   "Propose a new admin for a denom (two-step transfer)",
		Long:    "Propose a new admin for a denom. The new admin must accept with accept-admin. Only the current admin can execute this.",
		Example: version.AppName + " tx tokenfactory change-admin factory/heya1abc123/mytoken heya1xyz789",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgChangeAdmin(clientCtx.GetFromAddress().String(), args[0], args[1])
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewAcceptAdminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "accept-admin [denom]",
		Short:   "Accept the pending admin role for a denom",
		Long:    "Accept the pending admin role for a denom. Only the address proposed as new admin can execute this.",
		Example: version.AppName + " tx tokenfactory accept-admin factory/heya1abc123/mytoken",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgAcceptAdmin(clientCtx.GetFromAddress().String(), args[0])
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewForceTransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "force-transfer [amount] [dest-addr]",
		Short:   "Force transfer tokens to an address (admin only)",
		Long:    "Forcefully transfer tokens from any holder to a destination address. Only the denom admin can execute this. Use --from-address to specify the source (defaults to admin's address).",
		Example: version.AppName + " tx tokenfactory force-transfer 1000000factory/heya1abc123/mytoken heya1destaddr --from-address heya1source",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddr, _ := cmd.Flags().GetString("from-address")
			msg := types.NewMsgForceTransferFull(clientCtx.GetFromAddress().String(), args[0], args[1], fromAddr)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String("from-address", "", "address to transfer from (defaults to sender)")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
