package abci

import (
	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PreBlockHandler is responsible for aggregating oracle data from each validator and writing the
// oracle data into the store before any transactions are finalized for a given block.
type PreBlockHandler struct {
	logger log.Logger

	// oracle tracks, updates, and verifies oracle data.
	oracle *Oracle
}

// NewPreBlockHandler returns a new PreBlockHandler. The handler is responsible for writing oracle
// data included in vote extensions to state.
func NewPreBlockHandler(
	logger log.Logger,
	oracle *Oracle,
) *PreBlockHandler {
	return &PreBlockHandler{
		logger: logger.With("module", "oracle"),
		oracle: oracle,
	}
}

// PreBlockHook is called by the base app before the block is finalized. It is responsible
// for aggregating oracle data from each validator and writing the oracle data to the store.
//
// NOTE: The results of the oracle verification between PrepareProposal, ProcessProposal and
// PreBlock SHOULD be the same.
//
// TODO: Figure out what (if any) are the consequences of committing prices in PreBlock instead of
// ProcessProposal.
func (hook *PreBlockHandler) PreBlockHook() sdk.PreBlocker {
	return func(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
		// If vote extensions are not enabled, then we don't need to do anything.
		if !VoteExtensionsEnabled(ctx) || req == nil {
			return &sdk.ResponsePreBlock{}, nil
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
			return &sdk.ResponsePreBlock{}, err
		}

		// Write the oracle data to the store. This should never fail because the
		// oracle data was already verified in ProcessProposal.
		if err := hook.oracle.WriteOracleData(ctx, oracleData); err != nil {
			hook.logger.Error(
				"failed to write oracle data",
				"num_prices", len(oracleData.Prices),
				"err", err,
			)
			return &sdk.ResponsePreBlock{}, err
		}

		return &sdk.ResponsePreBlock{}, nil
	}
}
