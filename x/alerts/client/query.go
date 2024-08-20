package client

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/skip-mev/connect/v2/x/alerts/types"
)

const (
	flagAlertStatusID = "alert-status"
)

// GetQueryCmd returns the parent command for all x/alerts cli query commands.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the alerts module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryAlerts(),
	)

	return cmd
}

// CmdQueryParams returns the command for querying the module's parameters.
func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current alerts module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Params(clientCtx.CmdContext, &types.ParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryAlerts returns the command for querying alerts.
func CmdQueryAlerts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alerts",
		Short: "Query alerts by status (concluded, open, or all). See --help for more info",
		Long: `
The query is expected to look as follows:
	query alerts alerts --status <concluded|open> -> returns all queries with the given status
	query alerts alerts -> returns all alerts
		`,
		Example: "alerts alerts --alert-status concluded",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// get the alert status from flags if it exists
			alertStatusID, err := cmd.Flags().GetString(flagAlertStatusID)
			if err != nil {
				return err
			}

			// convert the alert status to an alert status id
			status, err := stringToAlertStatusID(alertStatusID)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Alerts(clientCtx.CmdContext, &types.AlertsRequest{
				Status: status,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().String(flagAlertStatusID, "", "filter alerts by status")

	return cmd
}

func stringToAlertStatusID(status string) (types.AlertStatusID, error) {
	switch status {
	case "open":
		return types.AlertStatusID_CONCLUSION_STATUS_UNCONCLUDED, nil
	case "concluded":
		return types.AlertStatusID_CONCLUSION_STATUS_CONCLUDED, nil
	case "":
		return types.AlertStatusID_CONCLUSION_STATUS_UNSPECIFIED, nil
	default:
		return types.AlertStatusID_CONCLUSION_STATUS_UNSPECIFIED, fmt.Errorf("invalid alert status: %s", status)
	}
}
