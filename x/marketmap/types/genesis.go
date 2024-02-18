package types

// NewGenesisState returns an instance of GenesisState.
func NewGenesisState(tickersConfig TickersConfig) GenesisState {
	return GenesisState{
		Tickers: tickersConfig,
	}
}

// ValidateBasic performs basic validation on the GenesisState.
func (gs *GenesisState) ValidateBasic() error {
	return gs.Tickers.ValidateBasic()
}
