package abci

import (
	"fmt"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/skip-mev/slinky/oracle/types"
)

type ProposalHandler struct {
	logger log.Logger

	// prepareProposalHandler fills a proposal with transactions.
	prepareProposalHandler sdk.PrepareProposalHandler

	// processProposalHandler processes transactions in a proposal.
	processProposalHandler sdk.ProcessProposalHandler

	// priceAggregator is responsible for aggregating prices from each validator
	// and computing the final oracle price for each asset.
	priceAggregator *oracletypes.PriceAggregator
}

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
	aggregateFn oracletypes.AggregateFn,
) *ProposalHandler {
	return &ProposalHandler{
		prepareProposalHandler: prepareProposalHandler,
		processProposalHandler: processProposalHandler,
		logger:                 logger,
		priceAggregator:        oracletypes.NewPriceAggregator(aggregateFn),
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

func (h *ProposalHandler) ProcessProposalHandler() sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
		h.logger.Info("processing proposal for height %d", req.Height)
		return h.processProposalHandler(ctx, req)
	}
}
