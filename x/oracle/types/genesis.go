package types

import "fmt"

// ValidateBasic validates that the CurrencyPair is valid, and performs any necessary validation on the
// genesis QuotePrice for the CurrencyPair. This fails if the CurrencyPair is invalid, or if the QuotePrice is nil,
// but the Nonce is non-nil.
func (cpg CurrencyPairGenesis) ValidateBasic() error {
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

// NewGenesisState returns a new genesis-state from a set of CurrencyPairGeneses
func NewGenesisState(cpgs []CurrencyPairGenesis) *GenesisState {
	return &GenesisState{
		CurrencyPairGenesis: cpgs,
	}
}

// DefaultGenesisState returns a default genesis state for the oracle module
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(nil)
}

// Validate validates the currency-pair geneses that the Genesis-State is composed of
func (gs GenesisState) Validate() error {
	for _, cpg := range gs.CurrencyPairGenesis {
		if err := cpg.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}
