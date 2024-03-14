package currencypair

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
)

// DefaultCurrencyPairStrategy is a strategy that uses the currency pair ID stored in the x/oracle state as
// the unique ID for a given currency pair and utilizes raw prices stored in the x/oracle state as the price
// representation for a given currency pair.
type DefaultCurrencyPairStrategy struct {
	oracleKeeper     OracleKeeper
	idToCurrencyPair map[uint64]slinkytypes.CurrencyPair
	useCache         bool
}

// Option is a function that enables optional configuration of the DefaultCurrencyPairStrategy.
type Option func(*DefaultCurrencyPairStrategy)

// WithCache enables the cache for DefaultCurrencyPairStrategy.
func WithCache() Option {
	return func(s *DefaultCurrencyPairStrategy) {
		s.useCache = true
	}
}

// NewDefaultCurrencyPairStrategy returns a new DefaultCurrencyPairStrategy instance.
func NewDefaultCurrencyPairStrategy(oracleKeeper OracleKeeper, opts ...Option) *DefaultCurrencyPairStrategy {
	strategy := &DefaultCurrencyPairStrategy{
		oracleKeeper:     oracleKeeper,
		idToCurrencyPair: make(map[uint64]slinkytypes.CurrencyPair),
		useCache:         false,
	}

	// run options
	for _, opt := range opts {
		opt(strategy)
	}

	return strategy
}

// ID returns the ID of the given currency pair, by querying the x/oracle state for the ID of the given
// currency pair. This method returns an error if the given currency pair is not found in the x/oracle state.
func (s *DefaultCurrencyPairStrategy) ID(ctx sdk.Context, cp slinkytypes.CurrencyPair) (uint64, error) {
	id, found := s.oracleKeeper.GetIDForCurrencyPair(ctx, cp)
	if !found {
		return 0, fmt.Errorf("currency pair %s not found in x/oracle state", cp.String())
	}

	// cache the currency pair for future lookups
	s.idToCurrencyPair[id] = cp

	return id, nil
}

// FromID returns the currency pair with the given ID, by querying the x/oracle state for the currency pair
// with the given ID. this method returns an error if the given ID is not currently present for an existing currency-pair.
func (s *DefaultCurrencyPairStrategy) FromID(ctx sdk.Context, id uint64) (slinkytypes.CurrencyPair, error) {
	// check the cache first
	if cp, found := s.idToCurrencyPair[id]; found {
		return cp, nil
	}

	cp, found := s.oracleKeeper.GetCurrencyPairFromID(ctx, id)
	if !found {
		return slinkytypes.CurrencyPair{}, fmt.Errorf("id %d out of bounds", id)
	}

	// cache the currency pair for future lookups
	s.idToCurrencyPair[id] = cp

	return cp, nil
}

// GetEncodedPrice returns the encoded price for the given currency pair. The default implementation
// returns the raw price, encoded into bytes.
func (s *DefaultCurrencyPairStrategy) GetEncodedPrice(
	_ sdk.Context,
	_ slinkytypes.CurrencyPair,
	price *big.Int,
) ([]byte, error) {
	if price.Sign() < 0 {
		return nil, fmt.Errorf("price cannot be negative: %s", price.String())
	}

	return price.GobEncode()
}

// GetDecodedPrice returns the decoded price for the given currency pair. The default implementation
// returns the raw price, decoded from bytes.
func (s *DefaultCurrencyPairStrategy) GetDecodedPrice(
	_ sdk.Context,
	_ slinkytypes.CurrencyPair,
	priceBytes []byte,
) (*big.Int, error) {
	var price big.Int
	if err := price.GobDecode(priceBytes); err != nil {
		return nil, err
	}

	if price.Sign() < 0 {
		return nil, fmt.Errorf("price cannot be negative: %s", price.String())
	}

	return &price, nil
}
