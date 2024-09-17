package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/oracle/types"
)

// GetQueryCmd returns the parent command for all x/oracle cli query commands. The
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
		GetPriceCmd(),
		GetAllCurrencyPairsCmd(),
	)

	return cmd
}

// GetPriceCmd returns the cli-command that queries the price information for a given CurrencyPair. This is essentially a wrapper around the module's
// QueryClient, as under-the-hood it constructs a request to a query-client served over a grpc-conn embedded in the clientCtx.
func GetPriceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "price [base] [quote]",
		Short: "Query for the price of a specified currency-pair",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// get context
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// retrieve CurrencyPair from arguments
			cp := connecttypes.NewCurrencyPair(args[0], args[1])

			// create client
			qc := types.NewQueryClient(clientCtx)

			// query for prices
			res, err := qc.GetPrice(cmd.Context(), &types.GetPriceRequest{
				CurrencyPair: cp.String(),
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

// GetAllCurrencyPairsCmd returns the cli-command that queries for all CurrencyPairs in the module. This is essentially a wrapper around the module's
// QueryClient, as under-the-hood it constructs a request to a query-client served over a grpc-conn embedded in the clientCtx.
func GetAllCurrencyPairsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "currency-pairs",
		Short: "Query for all the currency-pairs being tracked by the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// get the context
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// create a new query client
			qc := types.NewQueryClient(clientCtx)

			// query for all CurrencyPairs
			res, err := qc.GetAllCurrencyPairs(cmd.Context(), &types.GetAllCurrencyPairsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
