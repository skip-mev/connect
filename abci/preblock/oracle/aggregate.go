package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/holiman/uint256"
	"github.com/skip-mev/slinky/abci/ve/types"
	"github.com/skip-mev/slinky/aggregator"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// AggregateOracleVotes ingresses vote information which contains all of the
// vote extensions each validator extended in the previous block. it is important
// to note that
//  1. The vote extension may be nil, in which case the validator is not providing
//     any oracle data for the current block. This could have occurred because the
//     validator was offline, or its local oracle service was down.
//  2. The vote extension may contain prices updates for only a subset of currency pairs.
//     This could have occurred because the price providers for the validator were
//     offline, or the price providers did not provide a price update for a given
//     currency pair.
//
// In order for a currency pair to be included in the final oracle price, the currency
// pair must be provided by a supermajority (2/3+) of validators. This is enforced by the
// price aggregator but can be replaced by the application.
func (h *PreBlockHandler) AggregateOracleVotes(
	ctx sdk.Context,
	votes []Vote,
) (map[oracletypes.CurrencyPair]*uint256.Int, error) {
	// Reset the price aggregator and set the aggregationFn to use the latest application-state.
	h.priceAggregator.SetAggregationFn(h.aggregateFnWithCtx(ctx))
	h.priceAggregator.ResetProviderPrices()

	// Iterate through all vote extensions and consolidate all price info before
	// aggregating.
	isVotePresentInCommit := false
	for _, vote := range votes {
		if vote.ConsAddress.String() == h.validatorAddress.String() {
			isVotePresentInCommit = true
		}

		if err := h.addVoteToAggregator(vote.ConsAddress.String(), vote.OracleVoteExtension); err != nil {
			h.logger.Error(
				"failed to add vote to aggregator",
				"validator_address", vote.ConsAddress.String(),
				"err", err,
			)

			return nil, err
		}
	}

	// Compute the final prices for each currency pair.
	h.priceAggregator.UpdatePrices()
	prices := h.priceAggregator.GetPrices()

	// Record metrics for this validator.
	h.recordMetrics(isVotePresentInCommit)

	h.logger.Info(
		"aggregated oracle data",
		"num_prices", len(prices),
	)

	return prices, nil
}

// addVoteToAggregator consolidates the oracle data from a single validator
// into the price aggregator. The oracle data is provided in the form of a vote
// extension. The vote extension contains the prices for each currency pair that
// the validator is providing for the current block.
func (h *PreBlockHandler) addVoteToAggregator(address string, oracleData types.OracleVoteExtension) error {
	if len(oracleData.Prices) == 0 {
		return nil
	}

	// Format all of the prices into a map of currency pair -> price.
	prices := make(map[oracletypes.CurrencyPair]aggregator.QuotePrice)
	for pair, priceString := range oracleData.Prices {
		// Convert the price to a uint256.Int. All price feeds are expected to be
		// in the form of a string hex before conversion.
		price, err := uint256.FromHex(priceString)
		if err != nil {
			return err
		}

		// Convert the asset into a currency pair.
		currencyPair, err := oracletypes.CurrencyPairFromString(pair)
		if err != nil {
			return err
		}

		prices[currencyPair] = aggregator.QuotePrice{
			Price:     price,
			Timestamp: oracleData.Timestamp,
		}
	}

	h.logger.Debug(
		"adding oracle prices to aggregator",
		"num_prices", len(prices),
		"validator_address", address,
	)

	h.priceAggregator.SetProviderPrices(address, prices)

	return nil
}
