package oracle

import (
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	voteaggregator "github.com/skip-mev/slinky/abci/strategies/aggregator"
	"github.com/skip-mev/slinky/abci/strategies/codec"
	"github.com/skip-mev/slinky/abci/strategies/currencypair"
	"github.com/skip-mev/slinky/abci/types"
	"github.com/skip-mev/slinky/abci/ve"
	"github.com/skip-mev/slinky/aggregator"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	servicemetrics "github.com/skip-mev/slinky/service/metrics"
)

// PreBlockHandler is responsible for aggregating oracle data from each
// validator and writing the oracle data into the store before any transactions
// are executed/finalized for a given block.
type PreBlockHandler struct { //golint:ignore
	logger log.Logger

	// metrics is responsible for reporting / aggregating consensus-specific
	// metrics for this validator.
	metrics servicemetrics.Metrics

	// keeper is the keeper for the oracle module. This is utilized to write
	// oracle data to state.
	keeper Keeper

	// voteExtensionCodec is the codec used for encoding / decoding vote extensions.
	// This is used to decode vote extensions included in transactions.
	voteExtensionCodec codec.VoteExtensionCodec

	// extendedCommitCodec is the codec used for encoding / decoding extended
	// commit messages. This is used to decode extended commit messages included
	// in transactions.
	extendedCommitCodec codec.ExtendedCommitCodec

	// voteAggregator is responsible for aggregating votes from an extended commit into the canonical prices
	voteAggregator voteaggregator.VoteAggregator
}

// NewOraclePreBlockHandler returns a new PreBlockHandler. The handler
// is responsible for writing oracle data included in vote extensions to state.
func NewOraclePreBlockHandler(
	logger log.Logger,
	aggregateFn aggregator.AggregateFnFromContext[string, map[slinkytypes.CurrencyPair]*big.Int],
	oracleKeeper Keeper,
	metrics servicemetrics.Metrics,
	strategy currencypair.CurrencyPairStrategy,
	veCodec codec.VoteExtensionCodec,
	ecCodec codec.ExtendedCommitCodec,
) *PreBlockHandler {
	va := voteaggregator.NewDefaultVoteAggregator(
		logger,
		aggregateFn,
		strategy,
	)

	return &PreBlockHandler{
		logger:              logger,
		keeper:              oracleKeeper,
		metrics:             metrics,
		voteExtensionCodec:  veCodec,
		extendedCommitCodec: ecCodec,
		voteAggregator:      va,
	}
}

// PreBlocker is called by the base app before the block is finalized. It
// is responsible for aggregating oracle data from each validator and writing
// the oracle data to the store.
func (h *PreBlockHandler) PreBlocker() sdk.PreBlocker {
	return func(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) (_ *sdk.ResponsePreBlock, err error) {
		if req == nil {
			ctx.Logger().Error(
				"received nil RequestFinalizeBlock in oracle preblocker",
				"height", ctx.BlockHeight(),
			)

			return &sdk.ResponsePreBlock{}, fmt.Errorf("received nil RequestFinalizeBlock in oracle preblocker: height %d", ctx.BlockHeight())
		}

		start := time.Now()
		var prices map[slinkytypes.CurrencyPair]*big.Int
		defer func() {
			// only measure latency in Finalize
			if ctx.ExecMode() == sdk.ExecModeFinalize {
				latency := time.Since(start)
				h.logger.Info(
					"finished executing the pre-block hook",
					"height", ctx.BlockHeight(),
					"latency (seconds)", latency.Seconds(),
				)
				types.RecordLatencyAndStatus(h.metrics, latency, err, servicemetrics.PreBlock)

				// record prices + ticker metrics per validator (only do so if there was no error writing the prices)
				if err == nil && prices != nil {
					// record price metrics
					h.recordPrices(prices)

					// record validator report metrics
					h.recordValidatorReports(ctx, req.DecidedLastCommit)
				}
			}
		}()

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
		votes, err := voteaggregator.GetOracleVotes(req.Txs, h.voteExtensionCodec, h.extendedCommitCodec)
		if err != nil {
			h.logger.Error(
				"failed to get extended commit info from proposal",
				"height", req.Height,
				"num_txs", len(req.Txs),
				"err", err,
			)

			return &sdk.ResponsePreBlock{}, err
		}

		// Aggregate all oracle vote extensions into a single set of prices.
		prices, err = h.voteAggregator.AggregateOracleVotes(ctx, votes)
		if err != nil {
			h.logger.Error(
				"failed to aggregate oracle votes",
				"height", req.Height,
				"err", err,
			)

			err = PriceAggregationError{
				Err: err,
			}
			return &sdk.ResponsePreBlock{}, err
		}

		// Write the oracle data to the store.
		if err := h.WritePrices(ctx, prices); err != nil {
			h.logger.Error(
				"failed to write oracle data to store",
				"prices", prices,
				"err", err,
			)

			err = CommitPricesError{
				Err: err,
			}

			return &sdk.ResponsePreBlock{}, err
		}

		h.logger.Info("finished executing the oracle pre-block hook")

		return &sdk.ResponsePreBlock{}, nil
	}
}
