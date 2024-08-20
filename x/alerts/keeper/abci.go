package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/alerts/types"
)

// EndBlocker is called at the end of every block. This function is a no-op if pruning is disabled
//
//	It is used to determine which Alerts are to be purged, and if they should be purged, the alerts will be removed from state.
//
// If the AlertStatus of the Alert is concluded, nothing will be done. If the AlertStatus
// is Unconcluded, the alert will be Concluded positively (i.e, the bond will be returned to the bond-address).
func (k *Keeper) EndBlocker(goCtx context.Context) error {
	// unwrap the context
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if Pruning is enabled, if not this is a no-op
	if params := k.GetParams(ctx); !params.PruningParams.Enabled {
		return nil
	}

	// get the current block height
	height := uint64(ctx.BlockHeight())

	// get all alerts
	alerts, err := k.GetAllAlerts(ctx)
	if err != nil {
		return err
	}

	// iterate through all alerts
	for _, alert := range alerts {
		// check what the pruning height of the alert is, only prune if
		// the current block height is greater than or equal to the pruning height
		if alert.Status.PurgeHeight > height {
			continue
		}

		// check status of the alert, if it is to be pruned
		if status := alert.Status.ConclusionStatus; status != uint64(types.Concluded) {
			// conclude the alert positively
			if err := k.ConcludeAlert(ctx, alert.Alert, Positive); err != nil {
				return err
			}
		}

		// finally delete the alert
		if err := k.RemoveAlert(ctx, alert.Alert); err != nil {
			return err
		}
	}

	return nil
}
