package currencypair

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
)

const (
	// DefaultCacheInitialCapacity is the initial capacity to initialize the cache map for the
	// DefaultCurrencyPairStrategy.  This will prevent reallocation when under this size of CPs.
	DefaultCacheInitialCapacity = 50
)

// DefaultCurrencyPairStrategy is a strategy that uses the currency pair ID stored in the x/oracle state as
// the unique ID for a given currency pair and utilizes raw prices stored in the x/oracle state as the price
// representation for a given currency pair.
type DefaultCurrencyPairStrategy struct {
	oracleKeeper   OracleKeeper
	idCache        map[uint64]connecttypes.CurrencyPair
	previousHeight int64
}

// NewDefaultCurrencyPairStrategy returns a new DefaultCurrencyPairStrategy instance.
func NewDefaultCurrencyPairStrategy(oracleKeeper OracleKeeper) *DefaultCurrencyPairStrategy {
	strategy := &DefaultCurrencyPairStrategy{
		oracleKeeper: oracleKeeper,
		idCache:      make(map[uint64]connecttypes.CurrencyPair, DefaultCacheInitialCapacity),
	}
	return strategy
}

// ID returns the ID of the given currency pair, by querying the x/oracle state for the ID of the given
// currency pair. This method returns an error if the given currency pair is not found in the x/oracle state.
func (s *DefaultCurrencyPairStrategy) ID(ctx sdk.Context, cp connecttypes.CurrencyPair) (uint64, error) {
	// reset cache if the block height has changed
	height := ctx.BlockHeight()
	if height != s.previousHeight {
		s.idCache = make(map[uint64]connecttypes.CurrencyPair, DefaultCacheInitialCapacity)
		s.previousHeight = height
	}

	id, found := s.oracleKeeper.GetIDForCurrencyPair(ctx, cp)
	if !found {
		return 0, fmt.Errorf("currency pair %s not found in x/oracle state", cp.String())
	}

	// cache the currency pair for future lookups
	s.idCache[id] = cp

	return id, nil
}

// FromID returns the currency pair with the given ID, by querying the x/oracle state for the currency pair
// with the given ID. this method returns an error if the given ID is not currently present for an existing currency-pair.
func (s *DefaultCurrencyPairStrategy) FromID(ctx sdk.Context, id uint64) (connecttypes.CurrencyPair, error) {
	// reset cache if the block height has changed
	height := ctx.BlockHeight()
	if height != s.previousHeight {
		s.idCache = make(map[uint64]connecttypes.CurrencyPair, DefaultCacheInitialCapacity)
		s.previousHeight = height
	}

	// check the cache first
	if cp, found := s.idCache[id]; found {
		return cp, nil
	}

	cp, found := s.oracleKeeper.GetCurrencyPairFromID(ctx, id)
	if !found {
		return connecttypes.CurrencyPair{}, fmt.Errorf("id %d not found", id)
	}

	// cache the currency pair for future lookups
	s.idCache[id] = cp

	return cp, nil
}

// GetEncodedPrice returns the encoded price for the given currency pair. The default implementation
// returns the raw price, encoded into bytes.
func (s *DefaultCurrencyPairStrategy) GetEncodedPrice(
	_ sdk.Context,
	_ connecttypes.CurrencyPair,
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
	_ connecttypes.CurrencyPair,
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

// GetMaxNumCP returns the number of pairs that the VEs should include.  This method returns an error if the size cannot
// be queried from the x/oracle state. Specifically, this method should return the maximum number of currency pairs that
// could have existed at the time at which the votes were created. As such, if the execution mode is PrepareProposal or
// ProcessProposal, the number of removed currency pairs in the previous block should be included in the total.
func (s *DefaultCurrencyPairStrategy) GetMaxNumCP(
	ctx sdk.Context,
) (uint64, error) {
	current, err := s.oracleKeeper.GetNumCurrencyPairs(ctx)
	if err != nil {
		return 0, err
	}

	if mode := ctx.ExecMode(); mode == sdk.ExecModePrepareProposal || mode == sdk.ExecModeProcessProposal {
		removed, err := s.oracleKeeper.GetNumRemovedCurrencyPairs(ctx)
		if err != nil {
			return 0, err
		}

		return current + removed, nil
	}

	return current, nil
}
