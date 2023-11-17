package oracle

import (
	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/abci/ve"
	"github.com/skip-mev/slinky/aggregator"
	servicemetrics "github.com/skip-mev/slinky/service/metrics"
)

// PreBlockHandler is responsible for aggregating oracle data from each
// validator and writing the oracle data into the store before any transactions
// are executed/finalized for a given block.
type PreBlockHandler struct { //golint:ignore
	logger log.Logger

	// priceAggregator is responsible for aggregating prices from each validator
	// and computing the final oracle price for each asset.
	priceAggregator *aggregator.PriceAggregator

	// aggregateFnWithCtx is the aggregate function parametrized by the latest
	// state of the application.
	aggregateFnWithCtx aggregator.AggregateFnFromContext

	// metrics is responsible for reporting / aggregating consensus-specific
	// metrics for this validator.
	metrics servicemetrics.Metrics

	// validatorAddress is the consensus address of the validator running this
	// oracle.
	validatorAddress sdk.ConsAddress

	// keeper is the keeper for the oracle module. This is utilized to write
	// oracle data to state.
	keeper Keeper
}

// NewOraclePreBlockHandler returns a new PreBlockHandler. The handler
// is responsible for writing oracle data included in vote extensions to state.
func NewOraclePreBlockHandler(
	logger log.Logger,
	aggregateFn aggregator.AggregateFnFromContext,
	oracleKeeper Keeper,
	validatorConsAddress sdk.ConsAddress,
	metrics servicemetrics.Metrics,
) *PreBlockHandler {
	return &PreBlockHandler{
		logger:             logger,
		priceAggregator:    aggregator.NewPriceAggregator(aggregateFn(sdk.Context{})),
		aggregateFnWithCtx: aggregateFn,
		keeper:             oracleKeeper,
		validatorAddress:   validatorConsAddress,
		metrics:            metrics,
	}
}

// PreBlocker is called by the base app before the block is finalized. It
// is responsible for aggregating oracle data from each validator and writing
// the oracle data to the store.
func (h *PreBlockHandler) PreBlocker() sdk.PreBlocker {
	return func(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
		// If vote extensions are not enabled, then we don't need to do anything.
		if !ve.VoteExtensionsEnabled(ctx) || req == nil {
			h.logger.Info(
				"vote extensions are not enabled",
				"height", ctx.BlockHeight(),
			)

			return &sdk.ResponsePreBlock{}, nil
		}

		h.logger.Info(
			"executing the pre-finalize block hook",
			"height", req.Height,
		)

		// If vote extensions have been enabled, the extended commit info - which
		// contains the vote extensions - must be included in the request.
		votes, err := GetOracleVotes(req.Txs)
		if err != nil {
			h.logger.Error(
				"failed to get extended commit info from proposal",
				"height", req.Height,
				"num_txs", len(req.Txs),
				"err", err,
			)

			return &sdk.ResponsePreBlock{}, err
		}

		// Aggregate all of the oracle vote extensions into a single set of prices.
		prices, err := h.AggregateOracleVotes(ctx, votes)
		if err != nil {
			h.logger.Error(
				"failed to aggregate oracle votes",
				"height", req.Height,
				"err", err,
			)

			return &sdk.ResponsePreBlock{}, err
		}

		// Write the oracle data to the store.
		if err := h.WritePrices(ctx, prices); err != nil {
			h.logger.Error(
				"failed to write oracle data to store",
				"prices", prices,
				"err", err,
			)

			return &sdk.ResponsePreBlock{}, err
		}

		h.logger.Info("finished executing the oracle pre-block hook")

		return &sdk.ResponsePreBlock{}, nil
	}
}
