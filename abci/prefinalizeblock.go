package abci

import (
	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PreFinalizeBlockHook is called by the base app before the block is finalized. It is
// responsible for aggregating oracle data from each validator and writing the oracle data
// into the store before any transactions are finalized for a given block.
type PreFinalizeBlockHandler struct {
	logger log.Logger

	// oracle tracks, updates, and verifies oracle data.
	oracle *Oracle
}

// NewPreFinalizeBlockHandler returns a new PreFinalizeBlockHandler. The handler is
// responsible for writing oracle data included in vote extensions to state.
func NewPreFinalizeBlockHandler(
	logger log.Logger,
	oracle *Oracle,
) *PreFinalizeBlockHandler {
	return &PreFinalizeBlockHandler{
		logger: logger.With("module", "oracle"),
		oracle: oracle,
	}
}

// PreFinalizeBlock is called by the base app before the block is finalized.
// It is responsible for aggregating oracle data from each validator and writing the oracle data
// to the store.
//
// NOTE: The results of the oracle verification between PrepareProposal, ProcessProposal and
// PreFinalizeBlock SHOULD be the same.
//
// TODO: Figure out what (if any) are the consequences of committing prices in
// PreFinalizeBlock instead of ProcessProposal.
func (hook *PreFinalizeBlockHandler) PreFinalizeBlockHook() sdk.PreFinalizeBlockHook {
	return func(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) error {
		// If vote extensions are not enabled, then we don't need to do anything.
		if !VoteExtensionsEnabled(ctx) || req == nil {
			return nil
		}

		hook.logger.Info(
			"executing the pre-finalize block hook",
			"height", req.Height,
		)

		// Check the oracle data included in the vote extensions. This should never
		// fail because the oracle data was already verified in ProcessProposal.
		oracleData, err := hook.oracle.CheckOracleData(ctx, req.Txs, req.Height)
		if err != nil {
			hook.logger.Error(
				"failed to check oracle data",
				"num_txs", len(req.Txs),
				"err", err,
			)

			return err
		}

		// Write the oracle data to the store. This should never fail because the
		// oracle data was already verified in ProcessProposal.
		if err := hook.oracle.WriteOracleData(ctx, oracleData); err != nil {
			hook.logger.Error(
				"failed to write oracle data",
				"num_prices", len(oracleData.Prices),
				"err", err,
			)

			return err
		}

		return nil
	}
}
