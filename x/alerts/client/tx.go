package client

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/alerts/types"
)

// GetTxCmd returns the parent command for all x/alerts cli transaction commands.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "alerts",
		Short:                      "Alerts transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	return cmd
}

func AlertTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alert",
		Short: "Create a new alert",
		Long: `
Create a new alert with the specified height, sender, and currency-pair.
	Example: "slinkyd tx alerts alert cosmos... 1 BTC/USD"
	Structure: "slinkyd tx alerts alert <sender> <height> <currency-pair>
`,
		Example: "slinkyd tx alerts alert cosmos... 1 BTC/USD",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if err := cmd.Flags().Set(flags.FlagFrom, args[0]); err != nil {
				return err
			}

			// get the height
			height, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			// get the currency-pair
			cp, err := slinkytypes.CurrencyPairFromString(args[2])
			if err != nil {
				return err
			}

			alert := types.NewAlert(height, clientCtx.FromAddress, cp)
			alertMsg := types.NewMsgAlert(alert)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), alertMsg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
