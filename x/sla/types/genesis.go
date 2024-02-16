package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewDefaultGenesisState returns a default genesis state for the module.
func NewDefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:     DefaultParams(),
		SLAs:       make([]PriceFeedSLA, 0),
		PriceFeeds: make([]PriceFeed, 0),
	}
}

// NewGenesisState returns a new genesis state for the module.
func NewGenesisState(slas []PriceFeedSLA, priceFeeds []PriceFeed, params Params) *GenesisState {
	return &GenesisState{
		SLAs:       slas,
		PriceFeeds: priceFeeds,
		Params:     params,
	}
}

// ValidateBasic performs basic validation of the genesis state data returning an
// error for any failed validation criteria.
func (gs *GenesisState) ValidateBasic() error {
	seen := make(map[string]struct{})
	slaLength := make(map[string]uint64)
	for _, sla := range gs.SLAs {
		if _, ok := seen[sla.ID]; ok {
			return fmt.Errorf("duplicate price feed sla id %s", sla.ID)
		}

		if err := sla.ValidateBasic(); err != nil {
			return err
		}

		seen[sla.ID] = struct{}{}
		slaLength[sla.ID] = sla.MaximumViableWindow
	}

	seenFeeds := make(map[string]struct{})
	for _, priceFeed := range gs.PriceFeeds {
		// The SLA must exist for the given price feed.
		if _, ok := seen[priceFeed.ID]; !ok {
			return fmt.Errorf("sla %s does not exist for the given price feed", priceFeed.ID)
		}

		// The SLA must have the same maximum viable window as the price feed.
		if slaLength[priceFeed.ID] != priceFeed.MaximumViableWindow {
			return fmt.Errorf("sla %s has a different maximum viable window than the price feed with same id", priceFeed.ID)
		}

		// There cannot be duplicate feeds.
		feedTuple := priceFeed.ID + sdk.ConsAddress(priceFeed.Validator).String() + priceFeed.CurrencyPair.String()
		if _, ok := seenFeeds[feedTuple]; ok {
			return fmt.Errorf(
				"duplicate sla id %s validator address %s and currency pair %s",
				priceFeed.ID,
				priceFeed.Validator,
				priceFeed.CurrencyPair,
			)
		}

		if err := priceFeed.ValidateBasic(); err != nil {
			return err
		}

		seenFeeds[feedTuple] = struct{}{}
	}

	return nil
}

// GetGenesisStateFromAppState returns x/sla GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var gs GenesisState
	cdc.MustUnmarshalJSON(appState[ModuleName], &gs)
	return gs
}
