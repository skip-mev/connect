package types

// NewGenesisState returns an instance of GenesisState.
func NewGenesisState(tickers ...Ticker) GenesisState {
	return GenesisState{
		Tickers: tickers,
	}
}

// ValidateBasic performs basic validation on the GenesisState.
func (gs *GenesisState) ValidateBasic() error {
	return Tickers(gs.Tickers).ValidateBasic()
}
