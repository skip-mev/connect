package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesisState returns a default genesis state, with default params
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams(sdk.DefaultBondDenom, nil), nil)
}

// NewGenesisState creates a new GenesisState object given the params for the module.
func NewGenesisState(params Params, alerts []AlertWithStatus) GenesisState {
	gs := GenesisState{
		Params: params,
		Alerts: alerts,
	}

	return gs
}

// ValidateBasic performs stateless validation on the GenesisState, specifically it
// validates that the Params are valid, and that each of the alerts is also valid.
func (gs GenesisState) ValidateBasic() error {
	// validate the genesis-state
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// alerts
	alerts := make(map[string]struct{})

	// validate each alert
	for _, alert := range gs.Alerts {
		if err := alert.ValidateBasic(); err != nil {
			return err
		}

		// check for uniqueness on the alert.height + currency-pair
		if _, ok := alerts[alertKey(alert)]; ok {
			return fmt.Errorf("duplicate alert in genesis state: %v", alert.Alert)
		}

		alerts[alertKey(alert)] = struct{}{}
	}

	return nil
}

func alertKey(a AlertWithStatus) string {
	return fmt.Sprint(a.Alert.Height) + a.Alert.CurrencyPair.ToString()
}
