package abci

import (
	"sort"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/holiman/uint256"
	oracletypes "github.com/skip-mev/slinky/oracle/types"
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
		Price      *uint256.Int
	}
)

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
	validatorStore baseapp.ValidatorStore,
	threshold math.LegacyDec,
) oracletypes.AggregateFn {
	return func(providers oracletypes.AggregatedProviderPrices) map[types.CurrencyPair]*uint256.Int {
		// Iterate through all providers and store stake weight + price for each currency pair.
		priceInfo := make(map[types.CurrencyPair]VoteWeightedPriceInfo)

		for valAddress, validatorPrices := range providers {
			for currencyPair, quotePrice := range validatorPrices {
				// Only include prices that are not nil.
				if quotePrice.Price == nil {
					continue
				}

				// Retrieve the validator from the validator store.
				address, err := sdk.ConsAddressFromBech32(valAddress)
				if err != nil {
					continue
				}

				voteWeight, _, err := validatorStore.BondedTokensAndPubKeyByConsAddr(ctx, address)
				if err != nil {
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
		prices := make(map[types.CurrencyPair]*uint256.Int)
		totalBondedTokens, err := validatorStore.TotalBondedTokens(ctx) // TODO: determine if total bonded tokens should be the staking metric that is used.
		if err != nil {
			// This should never error.
			panic(err)
		}

		for currencyPair, info := range priceInfo {
			// The total voting power % that submitted a price update for the given currency pair must be
			// greater than the threshold to be included in the final oracle price.
			if percentSubmitted := math.LegacyNewDecFromInt(info.TotalWeight).Quo(math.LegacyNewDecFromInt(totalBondedTokens)); percentSubmitted.GTE(threshold) {
				prices[currencyPair] = ComputeVoteWeightedMedian(info)
			}
		}

		return prices
	}
}

// VoteWeightedMedianFromContext returns a new VoteWeightedMedian aggregate function that is parametrized by the
// latest state of the application.
func VoteWeightedMedianFromContext(
	validatorStore baseapp.ValidatorStore,
	threshold math.LegacyDec,
) oracletypes.AggregateFnFromContext {
	return func(ctx sdk.Context) oracletypes.AggregateFn {
		return VoteWeightedMedian(ctx, validatorStore, threshold)
	}
}

// ComputeVoteWeightedMedian computes the stake-weighted median price for a given asset.
func ComputeVoteWeightedMedian(priceInfo VoteWeightedPriceInfo) *uint256.Int {
	// Sort the prices by price.
	sort.SliceStable(priceInfo.Prices, func(i, j int) bool {
		return priceInfo.Prices[i].Price.Lt(priceInfo.Prices[j].Price)
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
