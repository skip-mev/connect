package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
)

// NewDefaultGenesisState returns a default genesis state for the module.
func NewDefaultGenesisState() *GenesisState {
	return &GenesisState{}
}

// NewGenesisState returns a new genesis state for the module.
func NewGenesisState(incentives []IncentivesByType) *GenesisState {
	return &GenesisState{
		Registry: incentives,
	}
}

// ValidateBasic performs basic validation of the genesis state data returning an
// error for any failed validation criteria.
func (gs *GenesisState) ValidateBasic() error {
	seen := make(map[string]struct{})
	for _, entry := range gs.Registry {
		if _, ok := seen[entry.IncentiveType]; ok {
			return fmt.Errorf("duplicate incentive name %s", entry.IncentiveType)
		}

		if err := entry.ValidateBasic(); err != nil {
			return err
		}

		seen[entry.IncentiveType] = struct{}{}
	}

	return nil
}

// NewIncentives returns a new Incentives instance.
func NewIncentives(name string, incentives [][]byte) IncentivesByType {
	return IncentivesByType{
		IncentiveType: name,
		Entries:       incentives,
	}
}

// ValidateBasic performs basic validation of the Incentives data returning an
// error for any failed validation criteria.
func (i *IncentivesByType) ValidateBasic() error {
	if len(i.IncentiveType) == 0 {
		return fmt.Errorf("incentive name cannot be empty")
	}

	if len(i.Entries) == 0 {
		return fmt.Errorf("incentive %s must have at least one incentive", i.IncentiveType)
	}

	for _, incentive := range i.Entries {
		if len(incentive) == 0 {
			return fmt.Errorf("incentive %s cannot be empty", i.IncentiveType)
		}
	}

	return nil
}

// GetGenesisStateFromAppState returns x/incentives GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return genesisState
}
