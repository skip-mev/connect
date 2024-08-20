package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

// GetQueryCmd returns the parent command for all x/sla cli query commands.
func GetQueryCmd() *cobra.Command {
	// create base command
	cmd := &cobra.Command{
		Use:   slatypes.ModuleName,
		Short: fmt.Sprintf("Querying commands for the %s module", slatypes.ModuleName),
		RunE:  client.ValidateCmd,
	}

	// add sub-commands
	cmd.AddCommand(
		GetAllSLAsCmd(),
		GetParamsCmd(),
	)

	return cmd
}

// GetAllSLAsCmd returns the cli-command that queries all SLAs in the store.
func GetAllSLAsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slas",
		Short: "Query for all SLAs in the store",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := slatypes.NewQueryClient(clientCtx)
			resp, err := queryClient.GetAllSLAs(clientCtx.CmdContext, &slatypes.GetAllSLAsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetParamsCmd returns the cli-command that queries the current SLA parameters.
func GetParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query for the current SLA parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := slatypes.NewQueryClient(clientCtx)
			resp, err := queryClient.Params(clientCtx.CmdContext, &slatypes.ParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
