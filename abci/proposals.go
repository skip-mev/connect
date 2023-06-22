package abci

import (
	"fmt"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/abci/types"
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
//  1. Filling a proposal with transactions
//  2. Aggregating oracle data from each validator and injecting the oracle data
//     into the proposal.
//  3. Processing & verifying transactions/oracle data in a given proposal
//  4. Updating the oracle state
func NewProposalHandler(
	logger log.Logger,
	prepareProposalHandler sdk.PrepareProposalHandler,
	processProposalHandler sdk.ProcessProposalHandler,
	aggregateFn oracleservice.AggregateFn,
	baseapp App,
	oracleKeeper OracleKeeper,
	validateVoteExtensionsFn ValidateVoteExtensionsFn,
) *ProposalHandler {
	return &ProposalHandler{
		prepareProposalHandler:   prepareProposalHandler,
		processProposalHandler:   processProposalHandler,
		logger:                   logger,
		priceAggregator:          oracleservice.NewPriceAggregator(aggregateFn),
		baseApp:                  baseapp,
		oracleKeeper:             oracleKeeper,
		validateVoteExtensionsFn: validateVoteExtensionsFn,
	}
}

// PrepareProposalHandler returns a PrepareProposalHandler that will be called
// by base app when a new block proposal is requested. The PrepareProposalHandler
// will first attempt to fill the proposal with transactions. Then, the
// handler will aggregate all of the oracle data provided by the validators
// and inject the oracle data into the proposal. Oracle data is provided by
// validators in the form of vote extensions. The final proposal will include
// the oracle data as the first transaction in the proposal, followed by the
// transactions filled by the prepareProposalHandler.
func (h *ProposalHandler) PrepareProposalHandler() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
		if req == nil {
			h.logger.Error("prepare proposal request is nil")
			return nil, fmt.Errorf("nil request")
		}

		// Create a proposal full of transactions.
		h.logger.Info("preparing proposal", "height", req.Height)
		resp, err := h.prepareProposalHandler(ctx, req)
		if err != nil {
			h.logger.Error("failed to prepare proposal", "height", req.Height, "err", err)
			return nil, err
		}

		// Aggregate all of the oracle data provided by the validators.
		oracleInfo, err := h.AggregateOracleData(ctx, req.LocalLastCommit)
		if err != nil {
			h.logger.Error(
				"failed to aggregate oracle vote extensions",
				"height", req.Height,
				"err", err,
			)

			return nil, err
		}

		// Inject the oracle data into the proposal.
		resp.Txs = append([][]byte{oracleInfo}, resp.Txs...)

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
		h.logger.Info("processing proposal for height %d", req.Height)

		// There must be at least one slot in the proposal for the oracle info.
		if len(req.Txs) < NumInjectedTxs {
			h.logger.Error("invalid number of transactions in proposal", "height", req.Height)
			return nil, fmt.Errorf("invalid number of transactions in proposal; expected at least %d txs", NumInjectedTxs)
		}

		// Retrieve the oracle info from the proposal. This cannot be empty as we have to at least
		// verify that vote extensions were included and that they are valid.
		oracleInfoBytes := req.Txs[OracleInfoIndex]
		if len(oracleInfoBytes) == 0 {
			return nil, fmt.Errorf("oracle data is nil")
		}

		proposalOracleData := types.OracleData{}
		if err := proposalOracleData.Unmarshal(oracleInfoBytes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal oracle data: %w", err)
		}

		// Unmarshal the latest commit info which contains all of the vote extensions that were utilized
		// to create the oracle info.
		extendedCommitInfo := cometabci.ExtendedCommitInfo{}
		if err := extendedCommitInfo.Unmarshal(proposalOracleData.ExtendedCommitInfo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal extended commit info: %w", err)
		}

		// Verify that the vote extensions included in the proposal are valid.
		if err := h.validateVoteExtensionsFn(extendedCommitInfo); err != nil {
			h.logger.Error("failed to validate vote extensions", "err", err)
			return nil, err
		}

		// Verify that the oracle data provided by the proposer matches the vote extensions
		// included in the proposal.
		oracleData, err := h.VerifyOracleData(ctx, proposalOracleData, extendedCommitInfo)
		if err != nil {
			h.logger.Error("failed to verify oracle data", "err", err)
			return nil, err
		}

		// Write the oracle data to state.
		if err := h.WriteOracleData(ctx, oracleData); err != nil {
			h.logger.Error("failed to write oracle data to state", "err", err)
			return nil, err
		}

		// Process the transactions in the proposal with the oracle data removed.
		req.Txs = req.Txs[NumInjectedTxs:]
		resp, err := h.processProposalHandler(ctx, req)
		if err != nil {
			h.logger.Error("failed to process proposal", "err", err)
			return nil, err
		}

		return resp, nil
	}
}
