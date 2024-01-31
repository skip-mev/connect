package oracle

import (
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cometabci "github.com/cometbft/cometbft/abci/types"
	cometproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/skip-mev/slinky/abci/proposals"
	compression "github.com/skip-mev/slinky/abci/strategies/codec"
	slinkyabci "github.com/skip-mev/slinky/abci/types"
	"github.com/skip-mev/slinky/abci/ve/types"
	servicemetrics "github.com/skip-mev/slinky/service/metrics"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
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

		h.logger.Info(
			"set price for currency pair",
			"currency_pair", cp.String(),
			"quote_price", quotePrice.Price.String(),
		)
	}

	return nil
}

// recordPrice records all the given prices per ticker, and reports them as a float64
func (h *PreBlockHandler) recordPrices(prices map[oracletypes.CurrencyPair]*big.Int) {
	for ticker, price := range prices {
		floatPrice, _ := price.Float64()
		h.metrics.ObservePriceForTicker(ticker, floatPrice)
	}
}

// recordValidatorReports takes the commit decided for this block, and for each validator in the commit, records
// whether their vote was included in the commit, whether they reported a price for each currency-pair, and if so
// the price they reported.
func (h *PreBlockHandler) recordValidatorReports(ctx sdk.Context, decidedCommit cometabci.CommitInfo) {
	pricesToReport := h.keeper.GetAllCurrencyPairs(ctx)

	// iterate over each validator in the commit
	for _, vote := range decidedCommit.Votes {
		var nilVote bool
		validator := sdk.ConsAddress(vote.Validator.Address).String()
		// if the validator voted nil, record that status
		if vote.BlockIdFlag != cometproto.BlockIDFlagCommit {
			nilVote = true
		}
		// iterate over each currency-pair, and record whether the validator reported a price for it
		validatorPrices := h.priceAggregator.GetDataByProvider(validator)
		for _, cp := range pricesToReport {
			// if the validator reported a nil-vote, record that and skip
			if nilVote {
				h.metrics.AddValidatorReportForTicker(validator, cp, servicemetrics.Absent)
				continue
			}

			// otherwise, check if the validator reported a price for the currency-pair
			price, ok := validatorPrices[cp]
			if !ok || price == nil {
				h.metrics.AddValidatorReportForTicker(validator, cp, servicemetrics.MissingPrice)
				continue
			}

			// if the validator reported a price, record that price
			floatPrice, _ := price.Float64()
			h.metrics.AddValidatorReportForTicker(validator, cp, servicemetrics.WithPrice)
			h.metrics.AddValidatorPriceForTicker(validator, cp, floatPrice)
		}
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
