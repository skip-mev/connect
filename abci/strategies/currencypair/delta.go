package currencypair

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
)

// DeltaCurrencyPairStrategy is a strategy that inherits from the DefaultCurrencyPairStrategy but
// encodes/decodes the price as the difference between the current price and the previous price.
type DeltaCurrencyPairStrategy struct {
	*DefaultCurrencyPairStrategy
	cache          map[connecttypes.CurrencyPair]*big.Int
	previousHeight int64
}

// NewDeltaCurrencyPairStrategy returns a new DeltaCurrencyPairStrategy instance.
func NewDeltaCurrencyPairStrategy(oracleKeeper OracleKeeper) *DeltaCurrencyPairStrategy {
	return &DeltaCurrencyPairStrategy{
		DefaultCurrencyPairStrategy: NewDefaultCurrencyPairStrategy(oracleKeeper),
		cache:                       make(map[connecttypes.CurrencyPair]*big.Int, DefaultCacheInitialCapacity),
	}
}

// GetEncodedPrice returns the encoded price for the given currency pair. Before a price is encoded,
// it is first converted to a delta price by subtracting the current on-chain price with the given
// price. The delta price is then encoded into bytes.
func (s *DeltaCurrencyPairStrategy) GetEncodedPrice(
	ctx sdk.Context,
	cp connecttypes.CurrencyPair,
	price *big.Int,
) ([]byte, error) {
	if price.Sign() < 0 {
		return nil, fmt.Errorf("price cannot be negative: %s", price.String())
	}

	onChainPrice, err := s.getOnChainPrice(ctx, cp)
	if err != nil {
		return nil, err
	}

	deltaPrice := new(big.Int).Sub(price, onChainPrice)

	ctx.Logger().Debug(
		"encoded oracle price",
		"currency_pair", cp.String(),
		"price", deltaPrice.String(),
	)

	return deltaPrice.GobEncode()
}

// GetDecodedPrice returns the decoded price for the given currency pair. The inputted price will
// be decoded into a delta price, which is then added to the current on-chain price to get the
// final price. If the price for the currency pair is not currently present in state, the delta
// is strictly positive and is the price.
func (s *DeltaCurrencyPairStrategy) GetDecodedPrice(
	ctx sdk.Context,
	cp connecttypes.CurrencyPair,
	priceBytes []byte,
) (*big.Int, error) {
	onChainPrice, err := s.getOnChainPrice(ctx, cp)
	if err != nil {
		return nil, err
	}

	// Decode the price bytes into a delta price.
	var delta big.Int
	if err := delta.GobDecode(priceBytes); err != nil {
		return nil, err
	}

	updatedPrice := new(big.Int).Add(&delta, onChainPrice)
	if updatedPrice.Sign() < 0 {
		return nil, fmt.Errorf("price cannot be negative: %s", updatedPrice.String())
	}

	return updatedPrice, nil
}

// getOnChainPrice returns the on-chain price, if it exists, for the given currency pair. If the
// price is successfully fetched, the price is cached for future calls. The price cache is cleared
// when the height changes.
func (s *DeltaCurrencyPairStrategy) getOnChainPrice(ctx sdk.Context, cp connecttypes.CurrencyPair) (*big.Int, error) {
	height := ctx.BlockHeight()
	if height != s.previousHeight {
		s.cache = make(map[connecttypes.CurrencyPair]*big.Int, DefaultCacheInitialCapacity)
		s.previousHeight = height
	}

	// If the price is already cached, return it.
	if price, ok := s.cache[cp]; ok {
		return price, nil
	}

	// Fetch the current price for the currency pair.
	currentPrice := big.NewInt(0)
	quote, err := s.oracleKeeper.GetPriceForCurrencyPair(ctx, cp)
	if err != nil {
		var quotePriceNotExistError oracletypes.QuotePriceNotExistError
		noPriceErr := errors.As(err, &quotePriceNotExistError)
		if !noPriceErr {
			return nil, fmt.Errorf(
				"error getting price for currency pair (%s): %w",
				cp.String(),
				err,
			)
		}

	} else {
		currentPrice = quote.Price.BigInt()
	}

	// Cache the price and return it.
	s.cache[cp] = currentPrice
	return currentPrice, nil
}
