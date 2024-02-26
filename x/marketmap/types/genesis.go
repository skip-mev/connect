package types

import "fmt"

// NewGenesisState returns an instance of GenesisState.
func NewGenesisState(
	marketMap MarketMap,
	lastUpdated int64,
	params Params,
) GenesisState {
	return GenesisState{
		MarketMap:   marketMap,
		LastUpdated: lastUpdated,
		Params:      params,
	}
}

// ValidateBasic performs basic validation on the GenesisState.
func (gs *GenesisState) ValidateBasic() error {
	if err := gs.MarketMap.ValidateBasic(); err != nil {
		return err
	}

	if gs.LastUpdated < 0 {
		return fmt.Errorf("LastUpdated height cannot be less than 0, got %d", gs.LastUpdated)
	}

	return gs.Params.ValidateBasic()
}

// DefaultGenesisState returns the default genesis of the marketmap module.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		MarketMap: MarketMap{
			Tickers:   make(map[string]Ticker),
			Paths:     make(map[string]Paths),
			Providers: make(map[string]Providers),
		},
		LastUpdated: 1,
		Params:      DefaultParams(),
	}
}
