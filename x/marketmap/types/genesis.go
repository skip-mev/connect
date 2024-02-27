package types

// NewGenesisState returns an instance of GenesisState.
func NewGenesisState(
	marketMap MarketMap,
	lastUpdated uint64,
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
		LastUpdated: 0,
		Params:      DefaultParams(),
	}
}
