package oracle

import (
	"math/big"
	"time"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/abci/strategies/codec"
	"github.com/skip-mev/slinky/abci/strategies/currencypair"
	"github.com/skip-mev/slinky/abci/ve"
	"github.com/skip-mev/slinky/aggregator"
	servicemetrics "github.com/skip-mev/slinky/service/metrics"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// PreBlockHandler is responsible for aggregating oracle data from each
// validator and writing the oracle data into the store before any transactions
// are executed/finalized for a given block.
type PreBlockHandler struct { //golint:ignore
	logger log.Logger

	// priceAggregator is responsible for aggregating prices from each validator
	// and computing the final oracle price for each asset.
	priceAggregator *aggregator.DataAggregator[string, map[oracletypes.CurrencyPair]*big.Int]

	// aggregateFnWithCtx is the aggregate function parametrized by the latest
	// state of the application.
	aggregateFnWithCtx aggregator.AggregateFnFromContext[string, map[oracletypes.CurrencyPair]*big.Int]

	// metrics is responsible for reporting / aggregating consensus-specific
	// metrics for this validator.
	metrics servicemetrics.Metrics

	// validatorAddress is the consensus address of the validator running this
	// oracle.
	validatorAddress sdk.ConsAddress

	// keeper is the keeper for the oracle module. This is utilized to write
	// oracle data to state.
	keeper Keeper

	// currencyPairStrategy is the strategy used for generating / retrieving
	// price information for currency pairs.
	currencyPairStrategy currencypair.CurrencyPairStrategy

	// voteExtensionCodec is the codec used for encoding / decoding vote extensions.
	// This is used to decode vote extensions included in transactions.
	voteExtensionCodec codec.VoteExtensionCodec

	// extendedCommitCodec is the codec used for encoding / decoding extended
	// commit messages. This is used to decode extended commit messages included
	// in transactions.
	extendedCommitCodec codec.ExtendedCommitCodec
}

// NewOraclePreBlockHandler returns a new PreBlockHandler. The handler
// is responsible for writing oracle data included in vote extensions to state.
func NewOraclePreBlockHandler(
	logger log.Logger,
	aggregateFn aggregator.AggregateFnFromContext[string, map[oracletypes.CurrencyPair]*big.Int],
	oracleKeeper Keeper,
	validatorConsAddress sdk.ConsAddress,
	metrics servicemetrics.Metrics,
	strategy currencypair.CurrencyPairStrategy,
	veCodec codec.VoteExtensionCodec,
	ecCodec codec.ExtendedCommitCodec,
) *PreBlockHandler {
	priceAggregator := aggregator.NewDataAggregator[string, map[oracletypes.CurrencyPair]*big.Int](
		aggregator.WithAggregateFnFromContext(aggregateFn),
	)

	return &PreBlockHandler{
		logger:               logger,
		priceAggregator:      priceAggregator,
		aggregateFnWithCtx:   aggregateFn,
		keeper:               oracleKeeper,
		validatorAddress:     validatorConsAddress,
		metrics:              metrics,
		currencyPairStrategy: strategy,
		voteExtensionCodec:   veCodec,
		extendedCommitCodec:  ecCodec,
	}
}

// PreBlocker is called by the base app before the block is finalized. It
// is responsible for aggregating oracle data from each validator and writing
// the oracle data to the store.
func (h *PreBlockHandler) PreBlocker() sdk.PreBlocker {
	return func(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
		start := time.Now()
		defer func() {
			// only measure latency in Finalize
			if ctx.ExecMode() == sdk.ExecModeFinalize {
				latency := time.Since(start)
				h.logger.Info(
					"finished executing the pre-block hook",
					"height", ctx.BlockHeight(),
					"latency (seconds)", latency.Seconds(),
				)
				h.metrics.ObserveABCIMethodLatency(servicemetrics.PreBlock, latency)
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
		votes, err := GetOracleVotes(req.Txs, h.voteExtensionCodec, h.extendedCommitCodec)
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
