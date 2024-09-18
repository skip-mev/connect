package aggregator

import (
	"math/big"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/abci/strategies/codec"
	connectabcitypes "github.com/skip-mev/connect/v2/abci/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
)

// PriceApplier is an interface used in `ExtendVote` and `PreBlock` to apply the prices
// derived from the latest votes to state.
//
//go:generate mockery --name PriceApplier --filename mock_price_applier.go
type PriceApplier interface {
	// ApplyPricesFromVoteExtensions derives the aggregate prices per asset in accordance with the given
	// vote extensions + VoteAggregator. If a price exists for an asset, it is written to state. The
	// prices aggregated from vote-extensions are returned if no errors are encountered in execution,
	// otherwise an error is returned + nil prices.
	ApplyPricesFromVoteExtensions(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) (map[connecttypes.CurrencyPair]*big.Int, error)

	// GetPriceForValidator gets the prices reported by a given validator. This method depends
	// on the prices from the latest set of aggregated votes.
	GetPricesForValidator(validator sdk.ConsAddress) map[connecttypes.CurrencyPair]*big.Int
}

// oraclePriceApplier is an implementation of PriceApplier that applies prices to the oracle module.
type oraclePriceApplier struct {
	// va is a VoteAggregator that is used to aggregate votes into prices.
	va VoteAggregator

	// ok is the oracle keeper that is used to write prices to state.
	ok connectabcitypes.OracleKeeper

	// logger
	logger log.Logger

	// codecs
	voteExtensionCodec  codec.VoteExtensionCodec
	extendedCommitCodec codec.ExtendedCommitCodec
}

// NewOraclePriceApplier returns a new oraclePriceApplier.
func NewOraclePriceApplier(
	va VoteAggregator,
	ok connectabcitypes.OracleKeeper,
	voteExtensionCodec codec.VoteExtensionCodec,
	extendedCommitCodec codec.ExtendedCommitCodec,
	logger log.Logger,
) PriceApplier {
	return &oraclePriceApplier{
		va:                  va,
		ok:                  ok,
		logger:              logger,
		voteExtensionCodec:  voteExtensionCodec,
		extendedCommitCodec: extendedCommitCodec,
	}
}

func (opa *oraclePriceApplier) ApplyPricesFromVoteExtensions(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) (map[connecttypes.CurrencyPair]*big.Int, error) {
	// If vote extensions have been enabled, the extended commit info - which
	// contains the vote extensions - must be included in the request.
	votes, err := GetOracleVotes(req.Txs, opa.voteExtensionCodec, opa.extendedCommitCodec)
	if err != nil {
		opa.logger.Error(
			"failed to get extended commit info from proposal",
			"height", req.Height,
			"num_txs", len(req.Txs),
			"err", err,
		)

		return nil, err
	}

	opa.logger.Debug(
		"got oracle vote extensions",
		"height", req.Height,
		"num_votes", len(votes),
	)

	// Aggregate all oracle vote extensions into a single set of prices.
	prices, err := opa.va.AggregateOracleVotes(ctx, votes)
	if err != nil {
		opa.logger.Error(
			"failed to aggregate oracle votes",
			"height", req.Height,
			"err", err,
		)

		err = PriceAggregationError{
			Err: err,
		}
		return nil, err
	}

	currencyPairs := opa.ok.GetAllCurrencyPairs(ctx)
	for _, cp := range currencyPairs {
		price, ok := prices[cp]
		if !ok || price == nil {
			opa.logger.Debug(
				"no price for currency pair",
				"currency_pair", cp.String(),
			)

			continue
		}

		if price.Sign() == -1 {
			opa.logger.Error(
				"price is negative",
				"currency_pair", cp.String(),
				"price", price.String(),
			)

			continue
		}

		// Convert the price to a quote price and write it to state.
		quotePrice := oracletypes.QuotePrice{
			Price:          math.NewIntFromBigInt(price),
			BlockTimestamp: ctx.BlockHeader().Time,
			BlockHeight:    uint64(ctx.BlockHeight()), //nolint:gosec
		}

		if err := opa.ok.SetPriceForCurrencyPair(ctx, cp, quotePrice); err != nil {
			opa.logger.Error(
				"failed to set price for currency pair",
				"currency_pair", cp.String(),
				"quote_price", cp.String(),
				"err", err,
			)

			return nil, err
		}

		opa.logger.Debug(
			"set price for currency pair",
			"currency_pair", cp.String(),
			"quote_price", quotePrice.Price.String(),
		)
	}

	return prices, nil
}

func (opa *oraclePriceApplier) GetPricesForValidator(validator sdk.ConsAddress) map[connecttypes.CurrencyPair]*big.Int {
	return opa.va.GetPriceForValidator(validator)
}
