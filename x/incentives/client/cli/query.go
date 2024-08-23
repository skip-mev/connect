package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/skip-mev/connect/v2/x/incentives/types"
)

// GetQueryCmd returns the parent command for all x/incentives cli query commands. The
// provided clientCtx should have, at a minimum, a verifier, CometBFT RPC client,
// and marshaler set.
func GetQueryCmd() *cobra.Command {
	// create base-command
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		RunE:  client.ValidateCmd,
	}

	// add sub-commands
	cmd.AddCommand(
		GetIncentivesByTypeCmd(),
		GetAllIncentivesCmd(),
	)

	return cmd
}

// GetIncentivesByTypeCmd returns the cli-command that queries the incentives of a given type.
// This is essentially a wrapper around the module's QueryClient, as under-the-hood it constructs
// a request to a query-client served over a grpc-conn embedded in the clientCtx.
func GetIncentivesByTypeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "type [type]",
		Short: "Query for all incentives of a specified type",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// get context
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// retrieve incentive type from arguments
			incentiveType := args[0]

			// create client
			qc := types.NewQueryClient(clientCtx)

			// query for incentives
			res, err := qc.GetIncentivesByType(clientCtx.CmdContext, &types.GetIncentivesByTypeRequest{
				IncentiveType: incentiveType,
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

// GetAllIncentivesCmd returns the cli-command that queries all incentives currently stored in the
// incentives module. This is essentially a wrapper around the module's QueryClient, as under-the-hood
// it constructs a request to a query-client served over a grpc-conn embedded in the clientCtx.
func GetAllIncentivesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-incentives",
		Short: "Query for all incentives currently stored in the module",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			// get context
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// create client
			qc := types.NewQueryClient(clientCtx)

			// query for incentives
			res, err := qc.GetAllIncentives(clientCtx.CmdContext, &types.GetAllIncentivesRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
