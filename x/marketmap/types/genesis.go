package types

// NewGenesisState returns an instance of GenesisState.
func NewGenesisState(
	tickers map[string]Ticker,
	paths map[string]Paths,
	providers map[string]Providers,
) GenesisState {
	return GenesisState{
		MarketMap: MarketMap{
			Tickers:   tickers,
			Paths:     paths,
			Providers: providers,
		},
	}
}

// ValidateBasic performs basic validation on the GenesisState.
func (gs *GenesisState) ValidateBasic() error {
	return gs.MarketMap.ValidateBasic()
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
