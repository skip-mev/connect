package currencypair

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// DefaultCurrencyPairStrategy is a strategy that uses the currency pair ID stored in the x/oracle state as
// the unique ID for a given currency pair and utilizes raw prices stored in the x/oracle state as the price
// representation for a given currency pair.
type DefaultCurrencyPairStrategy struct {
	oracleKeeper OracleKeeper
}

// NewDefaultCurrencyPairStrategy returns a new DefaultCurrencyPairStrategy instance.
func NewDefaultCurrencyPairStrategy(oracleKeeper OracleKeeper) *DefaultCurrencyPairStrategy {
	return &DefaultCurrencyPairStrategy{
		oracleKeeper: oracleKeeper,
	}
}

// ID returns the ID of the given currency pair, by querying the x/oracle state for the ID of the given
// currency pair. this method returns an error if the given currency pair is not found in the x/oracle state.
func (s *DefaultCurrencyPairStrategy) ID(ctx sdk.Context, cp oracletypes.CurrencyPair) (uint64, error) {
	id, found := s.oracleKeeper.GetIDForCurrencyPair(ctx, cp)
	if !found {
		return 0, fmt.Errorf("currency pair %s not found in x/oracle state", cp.ToString())
	}

	return id, nil
}

// FromID returns the currency pair with the given ID, by querying the x/oracle state for the currency pair
// with the given ID. this method returns an error if the given ID is not currently present for an existing currency-pair.
func (s *DefaultCurrencyPairStrategy) FromID(ctx sdk.Context, id uint64) (oracletypes.CurrencyPair, error) {
	cp, found := s.oracleKeeper.GetCurrencyPairFromID(ctx, id)
	if !found {
		return oracletypes.CurrencyPair{}, fmt.Errorf("id %d out of bounds", id)
	}

	return cp, nil
}
