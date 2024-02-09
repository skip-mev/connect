package types

// NewGenesisState returns a instance of GenesisState.
func NewGenesisState(config AggregateMarketConfig) GenesisState {
	return GenesisState{
		Config: config,
	}
}

// ValidateBasic performs basic validation on the GenesisState.
func (gs GenesisState) ValidateBasic() error {
	return gs.Config.ValidateBasic()
}
