package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
)

// ValidateBasic validates that the CurrencyPair is valid, and performs any necessary validation on the
// genesis QuotePrice for the CurrencyPair. This fails if the CurrencyPair is invalid, or if the QuotePrice is nil,
// but the Nonce is non-nil.
func (cpg *CurrencyPairGenesis) ValidateBasic() error {
	// validate the CurrencyPair
	if err := cpg.CurrencyPair.ValidateBasic(); err != nil {
		return err
	}
	// check validity of nonce, the only time a nonce will be non-zero will be if a price update has been made for the
	// CurrencyPair
	if cpg.CurrencyPairPrice == nil && cpg.Nonce != 0 {
		return fmt.Errorf("invalid nonce, no price update but non-zero nonce: %v", cpg.Nonce)
	}

	return nil
}

// NewGenesisState returns a new genesis-state from a set of CurrencyPairGeneses.
func NewGenesisState(cpgs []CurrencyPairGenesis, nextID uint64) *GenesisState {
	return &GenesisState{
		CurrencyPairGenesis: cpgs,
		NextId:              nextID,
	}
}

// DefaultGenesisState returns a default genesis state for the oracle module.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(nil, 0)
}

// Validate validates the currency-pair geneses that the Genesis-State is composed of
// valid CurrencyPairGenesis, and that no ID for a currency-pair is repeated.
func (gs *GenesisState) Validate() error {
	ids := make(map[uint64]struct{})
	cps := make(map[string]struct{})
	for _, cpg := range gs.CurrencyPairGenesis {
		// validate the currency-pair genesis
		if err := cpg.ValidateBasic(); err != nil {
			return err
		}

		// check if the ID > gs.NextID
		if cpg.Id >= gs.NextId {
			return fmt.Errorf("invalid id: %v, must be less than next id: %v", cpg.Id, gs.NextId)
		}

		// check for a repeated ID
		if _, ok := ids[cpg.Id]; ok {
			return fmt.Errorf("repeated id: %v", cpg.Id)
		}

		// check for repeated currency-pairs
		if _, ok := cps[cpg.CurrencyPair.String()]; ok {
			return fmt.Errorf("repeated currency-pair: %v", cpg.CurrencyPair.String())
		}

		// add the ID to the set of IDs
		ids[cpg.Id] = struct{}{}

		// add the currency-pair to the set of currency-pairs
		cps[cpg.CurrencyPair.String()] = struct{}{}
	}

	return nil
}

// GetGenesisStateFromAppState returns x/oracle GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return genesisState
}
