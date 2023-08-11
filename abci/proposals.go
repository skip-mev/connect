package abci

import (
	"fmt"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

		// oracle tracks, updates, and verifies oracle data.
		oracle *Oracle
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
	oracle *Oracle,
) *ProposalHandler {
	return &ProposalHandler{
		logger:                 logger,
		prepareProposalHandler: prepareProposalHandler,
		processProposalHandler: processProposalHandler,
		oracle:                 oracle,
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
		voteExtensionsEnabled := VoteExtensionsEnabled(ctx)

		h.logger.Info(
			"preparing proposal",
			"height", req.Height,
			"vote_extensions_enabled", voteExtensionsEnabled,
		)

		if voteExtensionsEnabled {
			// Validate the vote extensions provided by the validators.
			if err := h.oracle.validateVoteExtensionsFn(
				ctx,
				h.oracle.validatorStore,
				req.Height,
				ctx.ChainID(),
				req.LocalLastCommit,
			); err != nil {
				err = fmt.Errorf("%w; extended_commit_info: %+v", err, req.LocalLastCommit)
				h.logger.Error("failed to validate vote extensions", "err", err)
				return nil, err
			}

			// Aggregate all of the oracle data provided by the validators in their vote extensions.
			oracleData, err := h.oracle.AggregateOracleData(ctx, req.LocalLastCommit)
			if err != nil {
				h.logger.Error("failed to aggregate oracle vote extensions", "err", err)
				return nil, err
			}

			// Apply the oracle data to the current state so that transactions can be executed
			// on top of the latest oracle data.
			if err := h.oracle.WriteOracleData(ctx, oracleData); err != nil {
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
		voteExtensionsEnabled := VoteExtensionsEnabled(ctx)

		h.logger.Info(
			"processing proposal",
			"height", req.Height,
			"num_txs", len(req.Txs),
			"vote_extensions_enabled", voteExtensionsEnabled,
		)

		if voteExtensionsEnabled {
			// Verify that the oracle data is valid
			oracleData, err := h.oracle.CheckOracleData(ctx, req.Txs, req.Height)
			if err != nil {
				h.logger.Error("failed to verify oracle data", "err", err)
				return nil, err
			}

			// Apply the oracle data to the current state so that transactions can be executed
			// on top of the latest oracle data.
			if err := h.oracle.WriteOracleData(ctx, oracleData); err != nil {
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
