package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/alerts/types"
)

type ConclusionStatus uint64

const (
	Positive ConclusionStatus = iota
	Negative
)

// ConcludeAlert takes an Alert and status. This method returns the Alert's bond to the Alert's owner if the Alert
// is concluded with positive status, if the alert is concluded with negative status, the bond is burned.
// Finally, the alert's status is set to Concluded, and it's purge height is set to alert.Height + AlertParams.MaxBlockAge.
func (k *Keeper) ConcludeAlert(ctx sdk.Context, alertToConclude types.Alert, status ConclusionStatus) error {
	// check that the alert is valid
	if err := alertToConclude.ValidateBasic(); err != nil {
		return err
	}

	// get the alert from state
	alert, ok := k.GetAlert(ctx, alertToConclude)
	if !ok {
		return fmt.Errorf("alert not found: %v", alertToConclude)
	}

	// check if the alert is already concluded
	if alert.Status.ConclusionStatus != uint64(types.Unconcluded) {
		return fmt.Errorf("alert already concluded")
	}

	params := k.GetParams(ctx)

	switch status {
	case Positive:
		if err := k.unescrowBond(ctx, alert.Alert, params.AlertParams.BondAmount); err != nil {
			return err
		}
	case Negative:
		if err := k.burnBond(ctx, params.AlertParams.BondAmount); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid status: %v", status)
	}

	// update the status of the alert
	alert.Status.ConclusionStatus = uint64(types.Concluded)

	// set the purge height of the alert
	alert.Status.PurgeHeight = alert.Alert.Height + params.AlertParams.MaxBlockAge

	// set the alert
	return k.SetAlert(ctx, alert)
}

// unescrowBond sends the bond at the module account back to the alert's signer.
func (k *Keeper) unescrowBond(ctx sdk.Context, a types.Alert, bond sdk.Coin) error {
	alertSigner, err := sdk.AccAddressFromBech32(a.Signer)
	if err != nil {
		return err
	}

	// send the coins from the module account to the signer
	return k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		alertSigner,
		sdk.NewCoins(bond),
	)
}

// burnBond burns the bond stored at the module account's address.
func (k *Keeper) burnBond(ctx sdk.Context, bond sdk.Coin) error {
	// burn the coins
	return k.bankKeeper.BurnCoins(
		ctx,
		types.ModuleName,
		sdk.NewCoins(bond),
	)
}
