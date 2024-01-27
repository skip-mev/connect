package oracle

import (
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/abci/proposals"
	compression "github.com/skip-mev/slinky/abci/strategies/codec"
	slinkyabci "github.com/skip-mev/slinky/abci/types"
	"github.com/skip-mev/slinky/abci/ve/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	servicemetrics "github.com/skip-mev/slinky/service/metrics"
)

// Vote encapsulates the validator and oracle data contained within a vote extension.
type Vote struct {
	// ConsAddress is the validator that submitted the vote extension.
	ConsAddress sdk.ConsAddress
	// OracleVoteExtension
	OracleVoteExtension types.OracleVoteExtension
}

// WritePrices writes the oracle data to state. Note, this will only write prices
// that are already present in state.
func (h *PreBlockHandler) WritePrices(ctx sdk.Context, prices map[oracletypes.CurrencyPair]*big.Int) error {
	currencyPairs := h.keeper.GetAllCurrencyPairs(ctx)
	for _, cp := range currencyPairs {
		price, ok := prices[cp]
		if !ok || price == nil {
			h.logger.Info(
				"no price for currency pair",
				"currency_pair", cp.String(),
			)

			continue
		}

		// Convert the price to a quote price and write it to state.
		quotePrice := oracletypes.QuotePrice{
			Price:          math.NewIntFromBigInt(price),
			BlockTimestamp: ctx.BlockHeader().Time,
			BlockHeight:    uint64(ctx.BlockHeight()),
		}

		if err := h.keeper.SetPriceForCurrencyPair(ctx, cp, quotePrice); err != nil {
			h.logger.Error(
				"failed to set price for currency pair",
				"currency_pair", cp.String(),
				"quote_price", cp.String(),
				"err", err,
			)

			return err
		}
		
		// observe prices
		recordPrice(h.metrics, price, cp)

		h.logger.Info(
			"set price for currency pair",
			"currency_pair", cp.String(),
			"quote_price", quotePrice.Price.String(),
		)
	}

	return nil
}

// recordPrice takes the given metrics + big.Int value and reports the price as a float64
func recordPrice(metrics servicemetrics.Metrics, value *big.Int, ticker oracletypes.CurrencyPair) {
	float, _ := value.Float64()
	metrics.ObservePriceForTicker(ticker, float)
}

// recordMetrics reports whether the validator's vote-extension was included in the last commit, and
// the number of tickers for which the validator reported a price.
func (h *PreBlockHandler) recordMetrics(validatorVotePresent bool) {
	// determine which tickers this validator reported prices for
	validatorPrices := h.priceAggregator.GetDataByProvider(h.validatorAddress.String())

	h.logger.Info(
		"recording metrics for validator",
		"validator", h.validatorAddress.String(),
		"is_vote_present_in_commit", validatorVotePresent,
		"num_tickers", len(validatorPrices),
	)

	// determine if the validator's vote was included in the last commit
	h.metrics.AddVoteIncludedInLastCommit(validatorVotePresent)

	for ticker := range validatorPrices {
		h.metrics.AddTickerInclusionStatus(ticker.String(), true)
	}
}

// GetOracleVotes returns all of the oracle vote extensions that were injected into
// the block. Note that all of the vote extensions included are necessarily valid at this point
// because the vote extensions were validated by the vote extension and proposal handlers.
func GetOracleVotes(
	proposal [][]byte,
	veCodec compression.VoteExtensionCodec,
	extCommitCodec compression.ExtendedCommitCodec,
) ([]Vote, error) {
	if len(proposal) < proposals.NumInjectedTxs {
		return nil, slinkyabci.MissingCommitInfoError{}
	}

	extendedCommitInfo, err := extCommitCodec.Decode(proposal[proposals.OracleInfoIndex])
	if err != nil {
		return nil, slinkyabci.CodecError{
			Err: fmt.Errorf("error decoding extended-commit-info: %v", err),
		}
	}

	votes := make([]Vote, len(extendedCommitInfo.Votes))
	for i, voteInfo := range extendedCommitInfo.Votes {
		voteExtension, err := veCodec.Decode(voteInfo.VoteExtension)
		if err != nil {
			return nil, slinkyabci.CodecError{
				Err: fmt.Errorf("error decoding vote-extension: %v", err),
			}
		}

		address := sdk.ConsAddress{}
		if err := address.Unmarshal(voteInfo.Validator.Address); err != nil {
			return nil, slinkyabci.CodecError{
				Err: fmt.Errorf("error decoding validator address: %v", err),
			}
		}

		votes[i] = Vote{
			ConsAddress:         address,
			OracleVoteExtension: voteExtension,
		}
	}

	return votes, nil
}
