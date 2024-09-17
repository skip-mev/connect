package currencypair

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
)

// HashCurrencyPairStrategy is a strategy that uses the sha256 hash of the currency
// pair as the unique ID for a given currency pair and utilizes raw prices as the
// price representation for a given currency pair.
type HashCurrencyPairStrategy struct {
	*DefaultCurrencyPairStrategy
}

// NewHashCurrencyPairStrategy returns a new HashCurrencyPairStrategy instance.
func NewHashCurrencyPairStrategy(oracleKeeper OracleKeeper) *HashCurrencyPairStrategy {
	return &HashCurrencyPairStrategy{
		DefaultCurrencyPairStrategy: NewDefaultCurrencyPairStrategy(oracleKeeper),
	}
}

// ID returns the ID of the given currency pair, by taking the hash of the currency
// pair and using that as the ID.
func (s *HashCurrencyPairStrategy) ID(ctx sdk.Context, cp connecttypes.CurrencyPair) (uint64, error) {
	// reset cache if the block height has changed
	height := ctx.BlockHeight()
	if height != s.previousHeight {
		s.idCache = make(map[uint64]connecttypes.CurrencyPair, DefaultCacheInitialCapacity)
		s.previousHeight = height
	}

	// Check that the currency pair exists in state.
	_, found := s.oracleKeeper.GetIDForCurrencyPair(ctx, cp)
	if !found {
		return 0, fmt.Errorf("currency pair %s not found in x/oracle state", cp.String())
	}

	hash, err := CurrencyPairToHashID(cp.String())
	if err != nil {
		return 0, fmt.Errorf("failed to hash currency pair %s: %w", cp.String(), err)
	}

	s.idCache[hash] = cp
	return hash, nil
}

// FromID returns the currency pair with the given ID, it first checks the cache
// for the currency pair and if it is not found in the cache, it will attempt to
// retrieve the currency pair from the x/oracle state.
func (s *HashCurrencyPairStrategy) FromID(ctx sdk.Context, id uint64) (connecttypes.CurrencyPair, error) {
	// reset cache if the block height has changed
	height := ctx.BlockHeight()
	if height != s.previousHeight {
		s.idCache = make(map[uint64]connecttypes.CurrencyPair, DefaultCacheInitialCapacity)
		s.previousHeight = height
	}

	cp, found := s.idCache[id]
	if found {
		return cp, nil
	}

	// if the currency pair is not found in the cache, attempt to retrieve it from
	// the x/oracle state by populating the cache with all currency pairs. This
	// should only be executed once per block height.
	allCPs := s.oracleKeeper.GetAllCurrencyPairs(ctx)
	for _, cp := range allCPs {
		hash, err := CurrencyPairToHashID(cp.String())
		if err != nil {
			return connecttypes.CurrencyPair{}, fmt.Errorf("failed to hash currency pair %s: %w", cp.String(), err)
		}

		s.idCache[hash] = cp
	}

	cp, found = s.idCache[id]
	if !found {
		return connecttypes.CurrencyPair{}, fmt.Errorf("currency pair with sha256 hashed ID %d not found in x/oracle state", id)
	}

	return cp, nil
}

// CurrencyPairToHashID returns the ID of the given currency pair, by taking the
// sha256 hash of the currency pair.
func CurrencyPairToHashID(currencyPair string) (uint64, error) { //nolint
	hash := sha256.New()
	if _, err := hash.Write([]byte(currencyPair)); err != nil {
		return 0, err
	}

	md := hash.Sum(nil)
	return binary.LittleEndian.Uint64(md), nil
}
