package abci

import (
	"fmt"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracleservice "github.com/skip-mev/slinky/oracle/types"
)

const (
	// NumInjectedTxs is the number of transactions that were injected into
	// the proposal but are not actual transactions. In this case, the oracle
	// info is injected into the proposal but should be ignored by the application.
	NumInjectedTxs = 1

	// OracleInfoIndex is the index of the oracle info in the proposal.
	OracleInfoIndex = 0
)

type (
	ProposalHandler struct {
		logger log.Logger

		// prepareProposalHandler fills a proposal with transactions.
		prepareProposalHandler sdk.PrepareProposalHandler

		// processProposalHandler processes transactions in a proposal.
		processProposalHandler sdk.ProcessProposalHandler

		// priceAggregator is responsible for aggregating prices from each validator
		// and computing the final oracle price for each asset.
		priceAggregator *oracleservice.PriceAggregator

		// baseApp is the base application. This is utilized to retrieve the
		// state context for writing oracle data to state.
		baseApp App

		// oraclekeeper is the keeper for the oracle module. This is utilized
		// to write oracle data to state.
		oracleKeeper OracleKeeper

		// validateVoteExtensionsFn is the function responsible for validating vote extensions.
		validateVoteExtensionsFn ValidateVoteExtensionsFn
	}
)

// NewProposalHandler returns a new ProposalHandler. The proposalhandler is responsible
// primarily for:
//  1. Filling a proposal with transactions.
//  2. Aggregating oracle data from each validator and injecting the oracle data
//     into the proposal.
//  3. Processing & verifying transactions/oracle data in a given proposal.
//  4. Updating the oracle module state.
func NewProposalHandler(
	logger log.Logger,
	prepareProposalHandler sdk.PrepareProposalHandler,
	processProposalHandler sdk.ProcessProposalHandler,
	aggregateFn oracleservice.AggregateFn,
	baseApp App,
	oracleKeeper OracleKeeper,
	validateVoteExtensionsFn ValidateVoteExtensionsFn,
) *ProposalHandler {
	return &ProposalHandler{
		prepareProposalHandler:   prepareProposalHandler,
		processProposalHandler:   processProposalHandler,
		logger:                   logger,
		priceAggregator:          oracleservice.NewPriceAggregator(aggregateFn),
		baseApp:                  baseApp,
		oracleKeeper:             oracleKeeper,
		validateVoteExtensionsFn: validateVoteExtensionsFn,
	}
}

// PrepareProposalHandler returns a PrepareProposalHandler that will be called
// by base app when a new block proposal is requested. The PrepareProposalHandler
// will first attempt to aggregate all of the oracle data provided by the validators
// in their vote extensions. Then, it will apply the oracle data to the current state
// so that transactions can be executed on top of the latest oracle data. Finally,
// the handler will fill the proposal with transactions and inject the oracle data
// into the proposal (if vote extensions are enabled).
func (h *ProposalHandler) PrepareProposalHandler() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
		if req == nil {
			h.logger.Error("prepare proposal received a nil request")
			return nil, fmt.Errorf("nil request")
		}

		oracleDataBz := []byte{}
		voteExtensionsEnabled := h.VoteExtensionsEnabled(ctx)

		h.logger.Info(
			"preparing proposal",
			"height", req.Height,
			"vote_extensions_enabled", voteExtensionsEnabled,
		)

		if voteExtensionsEnabled {
			// Aggregate all of the oracle data provided by the validators in their vote extensions.
			oracleData, err := h.AggregateOracleData(ctx, req.LocalLastCommit)
			if err != nil {
				h.logger.Error("failed to aggregate oracle vote extensions", "err", err)
				return nil, err
			}

			// Apply the oracle data to the current state so that transactions can be executed
			// on top of the latest oracle data.
			if err := h.WriteOracleData(ctx, oracleData); err != nil {
				h.logger.Error("failed to write oracle data to state", "err", err)
				return nil, err
			}

			oracleDataBz, err = oracleData.Marshal()
			if err != nil {
				h.logger.Error("failed to marshal oracle data", "err", err)
				return nil, err
			}
		}

		resp, err := h.prepareProposalHandler(ctx, req)
		if err != nil {
			h.logger.Error("failed to prepare proposal", "err", err)
			return nil, err
		}

		// If vote extensions are enabled, each validator must inject the oracle data
		// into the proposal.
		if voteExtensionsEnabled {
			resp.Txs = append([][]byte{oracleDataBz}, resp.Txs...)
		}

		h.logger.Info("prepared proposal", "txs", len(resp.Txs))

		return resp, nil
	}
}

// ProcessProposalHandler returns a ProcessProposalHandler that will be called
// by base app when a new block proposal needs to be verified. The ProcessProposalHandler
// will first verify that the vote extensions included in the proposal are valid and compose
// a supermajority of signatures and vote extensions for the current block. Then, the
// handler will verify that the oracle data provided by the proposer matches the vote extensions
// included in the proposal. Finally, the handler will write the oracle data to state and
// process the transactions in the proposal with the oracle data removed.
func (h *ProposalHandler) ProcessProposalHandler() sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
		voteExtensionsEnabled := h.VoteExtensionsEnabled(ctx)

		h.logger.Info(
			"processing proposal",
			"height", req.Height,
			"num_txs", len(req.Txs),
			"vote_extensions_enabled", voteExtensionsEnabled,
		)

		if voteExtensionsEnabled {
			// Verify that the oracle data is valid
			oracleData, err := h.CheckOracleData(ctx, req)
			if err != nil {
				h.logger.Error("failed to verify oracle data", "err", err)
				return nil, err
			}

			// Apply the oracle data to the current state so that transactions can be executed
			// on top of the latest oracle data.
			if err := h.WriteOracleData(ctx, oracleData); err != nil {
				h.logger.Error("failed to write oracle data to state", "err", err)
				return nil, err
			}

			// Process the transactions in the proposal with the oracle data removed.
			req.Txs = req.Txs[NumInjectedTxs:]
		}

		resp, err := h.processProposalHandler(ctx, req)
		if err != nil {
			h.logger.Error("failed to process proposal", "err", err)
			return nil, err
		}

		return resp, nil
	}
}

// VoteExtensionsEnabled determines if vote extensions are enabled for the current block.
func (h *ProposalHandler) VoteExtensionsEnabled(ctx sdk.Context) bool {
	cp := ctx.ConsensusParams()
	if cp.Abci == nil || cp.Abci.VoteExtensionsEnableHeight == 0 {
		return false
	}

	// Per the cosmos sdk, the first block should not utilize the latest finalize block state. This means
	// vote extensions should NOT be making state changes.
	//
	// Ref: https://github.com/cosmos/cosmos-sdk/blob/2100a73dcea634ce914977dbddb4991a020ee345/baseapp/baseapp.go#L488-L495
	if ctx.BlockHeight() <= 1 {
		return false
	}

	// We do a +1 here because the vote extensions are enabled at height h
	// but a proposer will only receive vote extensions in height h+1.
	return cp.Abci.VoteExtensionsEnableHeight+1 < ctx.BlockHeight()
}
