package math

import (
	"math/big"
	"sort"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/x/oracle/types"
)

// DefaultPowerThreshold defines the total voting power % that must be
// submitted in order for a currency pair to be considered for the
// final oracle price. We provide a default supermajority threshold
// of 2/3+.
var DefaultPowerThreshold = math.LegacyNewDecWithPrec(667, 3)

type (
	// VoteWeightPriceInfo tracks the stake weight(s) + price(s) for a given currency pair.
	VoteWeightedPriceInfo struct {
		Prices      []VoteWeightedPricePerValidator
		TotalWeight math.Int
	}

	// VoteWeightPrice defines a price update that includes the stake weight of the validator.
	VoteWeightedPricePerValidator struct {
		VoteWeight math.Int
		Price      *big.Int
	}
)

// VoteWeightedMedianFromContext returns a new VoteWeightedMedian aggregate function that is parametrized by the
// latest state of the application.
func VoteWeightedMedianFromContext(
	logger log.Logger,
	validatorStore ValidatorStore,
	threshold math.LegacyDec,
) aggregator.AggregateFnFromContext {
	return func(ctx sdk.Context) aggregator.AggregateFn {
		return VoteWeightedMedian(ctx, logger, validatorStore, threshold)
	}
}

// VoteWeightedMedian returns an aggregation function that computes the stake weighted median price as the
// final deterministic oracle price for any qualifying currency pair (base, quote). There are a few things to
// note about the implementation:
//
//  1. Price updates for a given currency pair will only be written to state if the power % threshold
//     is met. The threshold is determined by the total voting power of all validators that
//     submitted a price update for a given currency pair divided by the total network voting power. The threshold
//     to meet is configurable by developers.
//  2. In the case where there are not enough price updates for a given currency pair, the
//     price will not be included in the final set of oracle prices.
//  3. Given the threshold is met, the final oracle price for a given currency pair is the
//     median price weighted by the stake of each validator that submitted a price.
func VoteWeightedMedian(
	ctx sdk.Context,
	logger log.Logger,
	validatorStore ValidatorStore,
	threshold math.LegacyDec,
) aggregator.AggregateFn {
	return func(providers aggregator.AggregatedProviderPrices) map[types.CurrencyPair]*big.Int {
		priceInfo := make(map[types.CurrencyPair]VoteWeightedPriceInfo)

		// Iterate through all providers and store stake weight + price for each currency pair.
		for valAddress, validatorPrices := range providers {
			// Retrieve the validator from the validator store and get its vote weight.
			address, err := sdk.ConsAddressFromBech32(valAddress)
			if err != nil {
				logger.Info(
					"failed to parse validator address; skipping validator prices",
					"validator_address", valAddress,
					"err", err,
				)

				continue
			}

			validator, err := validatorStore.ValidatorByConsAddr(ctx, address)
			if err != nil {
				logger.Info(
					"failed to retrieve validator from store; skipping validator prices",
					"validator_address", valAddress,
					"err", err,
				)

				continue
			}

			voteWeight := validator.GetBondedTokens()

			// Iterate through all prices and store the price + vote weight for each currency pair.
			for currencyPair, quotePrice := range validatorPrices {
				// Only include prices that are not nil.
				if quotePrice.Price == nil {
					logger.Info(
						"price is nil",
						"currency_pair", currencyPair.ToString(),
						"validator_address", valAddress,
					)

					continue
				}

				// Initialize the price info if it does not exist for the given currency pair.
				if _, ok := priceInfo[currencyPair]; !ok {
					priceInfo[currencyPair] = VoteWeightedPriceInfo{
						Prices:      make([]VoteWeightedPricePerValidator, 0),
						TotalWeight: math.ZeroInt(),
					}
				}

				// Update the price info.
				cpInfo := priceInfo[currencyPair]
				priceInfo[currencyPair] = VoteWeightedPriceInfo{
					Prices: append(cpInfo.Prices, VoteWeightedPricePerValidator{
						VoteWeight: voteWeight,
						Price:      quotePrice.Price,
					}),
					TotalWeight: cpInfo.TotalWeight.Add(voteWeight),
				}
			}
		}

		// Iterate through all prices and compute the median price for each asset.
		prices := make(map[types.CurrencyPair]*big.Int)
		totalBondedTokens, err := validatorStore.TotalBondedTokens(ctx)
		if err != nil {
			// This should never error.
			panic(err)
		}

		for currencyPair, info := range priceInfo {
			// The total voting power % that submitted a price update for the given currency pair must be
			// greater than the threshold to be included in the final oracle price.
			if percentSubmitted := math.LegacyNewDecFromInt(info.TotalWeight).Quo(math.LegacyNewDecFromInt(totalBondedTokens)); percentSubmitted.GTE(threshold) {
				prices[currencyPair] = ComputeVoteWeightedMedian(info)

				logger.Info(
					"computed stake-weighted median price for currency pair",
					"currency_pair", currencyPair.ToString(),
					"percent_submitted", percentSubmitted.String(),
					"threshold", threshold.String(),
					"final_price", prices[currencyPair].String(),
				)
			} else {
				logger.Info(
					"not enough voting power to compute stake-weighted median price price for currency pair",
					"currency_pair", currencyPair.ToString(),
					"threshold", threshold.String(),
					"percent_submitted", percentSubmitted.String(),
				)
			}
		}

		return prices
	}
}

// ComputeVoteWeightedMedian computes the stake-weighted median price for a given asset.
func ComputeVoteWeightedMedian(priceInfo VoteWeightedPriceInfo) *big.Int {
	// Sort the prices by price.
	sort.SliceStable(priceInfo.Prices, func(i, j int) bool {
		switch priceInfo.Prices[i].Price.Cmp(priceInfo.Prices[j].Price) {
		case -1:
			return true
		case 1:
			return false
		default:
			return true
		}
	})

	// Compute the median weight.
	middle := priceInfo.TotalWeight.QuoRaw(2)

	// Iterate through the prices and compute the median price.
	sum := math.ZeroInt()
	for index, price := range priceInfo.Prices {
		sum = sum.Add(price.VoteWeight)

		if sum.GTE(middle) {
			return price.Price
		}

		// If we reached the end of the list, return the last price.
		if index == len(priceInfo.Prices)-1 {
			return price.Price
		}
	}

	return nil
}
