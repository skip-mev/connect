package oracle

import (
	"fmt"
	"math/big"
	"time"

	"github.com/cosmos/cosmos-sdk/types/module"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	abciaggregator "github.com/skip-mev/connect/v2/abci/strategies/aggregator"
	"github.com/skip-mev/connect/v2/abci/strategies/codec"
	"github.com/skip-mev/connect/v2/abci/strategies/currencypair"
	connectabcitypes "github.com/skip-mev/connect/v2/abci/types"
	"github.com/skip-mev/connect/v2/abci/ve"
	"github.com/skip-mev/connect/v2/aggregator"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	servicemetrics "github.com/skip-mev/connect/v2/service/metrics"
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
	keeper connectabcitypes.OracleKeeper

	// pa is the price applier that is used to decode vote-extensions, aggregate price reports, and write prices to state.
	pa abciaggregator.PriceApplier
}

// NewOraclePreBlockHandler returns a new PreBlockHandler. The handler
// is responsible for writing oracle data included in vote extensions to state.
func NewOraclePreBlockHandler(
	logger log.Logger,
	aggregateFn aggregator.AggregateFnFromContext[string, map[connecttypes.CurrencyPair]*big.Int],
	oracleKeeper connectabcitypes.OracleKeeper,
	metrics servicemetrics.Metrics,
	strategy currencypair.CurrencyPairStrategy,
	veCodec codec.VoteExtensionCodec,
	ecCodec codec.ExtendedCommitCodec,
) *PreBlockHandler {
	va := abciaggregator.NewDefaultVoteAggregator(
		logger,
		aggregateFn,
		strategy,
	)
	pa := abciaggregator.NewOraclePriceApplier(
		va,
		oracleKeeper,
		veCodec,
		ecCodec,
		logger,
	)

	return &PreBlockHandler{
		logger:  logger,
		keeper:  oracleKeeper,
		metrics: metrics,
		pa:      pa,
	}
}

// WrappedPreBlocker is called by the base app before the block is finalized. It
// is responsible for calling the module manager's PreBlock method, aggregating oracle data from each validator and
// writing the oracle data to the store.
func (h *PreBlockHandler) WrappedPreBlocker(mm *module.Manager) sdk.PreBlocker {
	return func(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) (response *sdk.ResponsePreBlock, err error) {
		if req == nil {
			ctx.Logger().Error(
				"received nil RequestFinalizeBlock in oracle preblocker",
				"height", ctx.BlockHeight(),
			)

			return &sdk.ResponsePreBlock{}, fmt.Errorf("received nil RequestFinalizeBlock in oracle preblocker: height %d", ctx.BlockHeight())
		}

		// call module manager's PreBlocker first in case there is changes made on upgrades
		// that can modify state and lead to serialization/deserialization issues
		response, err = mm.PreBlock(ctx)
		if err != nil {
			return response, err
		}

		start := time.Now()
		var prices map[connecttypes.CurrencyPair]*big.Int
		defer func() {
			// only measure latency in Finalize
			if ctx.ExecMode() == sdk.ExecModeFinalize {
				latency := time.Since(start)
				h.logger.Debug(
					"finished executing the pre-block hook",
					"height", ctx.BlockHeight(),
					"latency (seconds)", latency.Seconds(),
				)
				connectabcitypes.RecordLatencyAndStatus(h.metrics, latency, err, servicemetrics.PreBlock)

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
		if !ve.VoteExtensionsEnabled(ctx) {
			h.logger.Info(
				"vote extensions are not enabled",
				"height", ctx.BlockHeight(),
			)

			return response, nil
		}

		h.logger.Debug(
			"executing the pre-finalize block hook",
			"height", req.Height,
		)

		// decode vote-extensions + apply prices to state
		prices, err = h.pa.ApplyPricesFromVoteExtensions(ctx, req)
		if err != nil {
			h.logger.Error(
				"failed to apply prices from vote extensions",
				"height", req.Height,
				"error", err,
			)

			return response, err
		}

		return response, nil
	}
}

// PreBlocker is called by the base app before the block is finalized. It
// is responsible for aggregating oracle data from each validator and writing
// the oracle data to the store.
//
// Deprecated: using PreBlocker requires wrapping module Manager's PreBlock call. This method should no longer be used.
// Use WrappedPreBlocker instead to handle this functionality automatically.
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
		var prices map[connecttypes.CurrencyPair]*big.Int
		defer func() {
			// only measure latency in Finalize
			if ctx.ExecMode() == sdk.ExecModeFinalize {
				latency := time.Since(start)
				h.logger.Debug(
					"finished executing the pre-block hook",
					"height", ctx.BlockHeight(),
					"latency (seconds)", latency.Seconds(),
				)
				connectabcitypes.RecordLatencyAndStatus(h.metrics, latency, err, servicemetrics.PreBlock)

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
		if !ve.VoteExtensionsEnabled(ctx) {
			h.logger.Info(
				"vote extensions are not enabled",
				"height", ctx.BlockHeight(),
			)

			return &sdk.ResponsePreBlock{}, nil
		}

		h.logger.Debug(
			"executing the pre-finalize block hook",
			"height", req.Height,
		)

		// decode vote-extensions + apply prices to state
		prices, err = h.pa.ApplyPricesFromVoteExtensions(ctx, req)
		if err != nil {
			h.logger.Error(
				"failed to apply prices from vote extensions",
				"height", req.Height,
				"error", err,
			)

			return &sdk.ResponsePreBlock{}, err
		}

		return &sdk.ResponsePreBlock{}, nil
	}
}
