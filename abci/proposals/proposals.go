package proposals

import (
	"bytes"
	"time"

	servicemetrics "github.com/skip-mev/slinky/service/metrics"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/abci/strategies/codec"
	"github.com/skip-mev/slinky/abci/strategies/currencypair"
	"github.com/skip-mev/slinky/abci/types"
	"github.com/skip-mev/slinky/abci/ve"
)

const (
	// NumInjectedTxs is the number of transactions that were injected into
	// the proposal but are not actual transactions. In this case, the oracle
	// info is injected into the proposal but should be ignored by the application.
	NumInjectedTxs = 1

	// OracleInfoIndex is the index of the oracle info in the proposal.
	OracleInfoIndex = 0
)

// The proposalhandler is responsible primarily for:
//  1. Filling a proposal with transactions.
//  2. Injecting vote extensions into the proposal (if vote extensions are enabled).
//  3. Verifying that the vote extensions injected are valid.
//
// To verify the validity of the vote extensions, the proposal handler will
// call the validateVoteExtensionsFn. This function is responsible for verifying
// that the vote extensions included in the proposal are valid and compose a
// supermajority of signatures and vote extensions for the current block.
// The given VoteExtensionCodec must be the same used by the VoteExtensionHandler,
// the extended commit is decoded in accordance with the given ExtendedCommitCodec.
type ProposalHandler struct {
	logger log.Logger

	// prepareProposalHandler fills a proposal with transactions.
	prepareProposalHandler sdk.PrepareProposalHandler

	// processProposalHandler processes transactions in a proposal.
	processProposalHandler sdk.ProcessProposalHandler

	// validateVoteExtensionsFn validates the vote extensions included in a proposal.
	validateVoteExtensionsFn ve.ValidateVoteExtensionsFn

	// voteExtensionCodec is used to decode vote extensions.
	voteExtensionCodec codec.VoteExtensionCodec

	// extendedCommitCodec is used to decode extended commit info.
	extendedCommitCodec codec.ExtendedCommitCodec

	// currencyPairStrategy is the strategy used to determine the price information
	// from a given oracle vote extension.
	currencyPairStrategy currencypair.CurrencyPairStrategy
	// metrics is responsible for reporting / aggregating consensus-specific
	// metrics for this validator.
	metrics servicemetrics.Metrics
}

// NewProposalHandler returns a new ProposalHandler.
func NewProposalHandler(
	logger log.Logger,
	prepareProposalHandler sdk.PrepareProposalHandler,
	processProposalHandler sdk.ProcessProposalHandler,
	validateVoteExtensionsFn ve.ValidateVoteExtensionsFn,
	voteExtensionCodec codec.VoteExtensionCodec,
	extendedCommitInfoCodec codec.ExtendedCommitCodec,
	currencyPairStrategy currencypair.CurrencyPairStrategy,
	metrics servicemetrics.Metrics,
) *ProposalHandler {
	return &ProposalHandler{
		logger:                   logger,
		prepareProposalHandler:   prepareProposalHandler,
		processProposalHandler:   processProposalHandler,
		validateVoteExtensionsFn: validateVoteExtensionsFn,
		voteExtensionCodec:       voteExtensionCodec,
		extendedCommitCodec:      extendedCommitInfoCodec,
		currencyPairStrategy:     currencyPairStrategy,
		metrics:                  metrics,
	}
}

// PrepareProposalHandler returns a PrepareProposalHandler that will be called
// by base app when a new block proposal is requested. The PrepareProposalHandler
// will first fill the proposal with transactions. Then, if vote extensions are
// enabled, the handler will inject the extended commit info into the proposal.
func (h *ProposalHandler) PrepareProposalHandler() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *cometabci.RequestPrepareProposal) (resp *cometabci.ResponsePrepareProposal, err error) {
		var (
			extInfoBz                     []byte
			wrappedPrepareProposalLatency time.Duration
		)
		startTime := time.Now()

		// report the slinky specific PrepareProposal latency
		defer func() {
			totalLatency := time.Since(startTime)
			h.logger.Info(
				"recording handle time metrics of prepare-proposal (seconds)",
				"total latency", totalLatency.Seconds(),
				"wrapped prepare proposal latency", wrappedPrepareProposalLatency.Seconds(),
				"slinky prepare proposal latency", (totalLatency - wrappedPrepareProposalLatency).Seconds(),
			)

			types.RecordLatencyAndStatus(h.metrics, totalLatency-wrappedPrepareProposalLatency, err, servicemetrics.PrepareProposal)
		}()

		if req == nil {
			h.logger.Error("PrepareProposalHandler received a nil request")
			err = types.NilRequestError{
				Handler: servicemetrics.PrepareProposal,
			}
			return nil, err
		}

		// If vote extensions are enabled, the current proposer must inject the extended commit
		// info into the proposal. This extended commit info contains the oracle data
		// for the current block.
		voteExtensionsEnabled := ve.VoteExtensionsEnabled(ctx)
		if voteExtensionsEnabled {
			h.logger.Info(
				"injecting oracle data into proposal",
				"height", req.Height,
				"vote_extensions_enabled", voteExtensionsEnabled,
			)

			extInfo := req.LocalLastCommit
			if err = h.ValidateExtendedCommitInfo(ctx, req.Height, extInfo); err != nil {
				h.logger.Error(
					"failed to validate vote extensions",
					"height", req.Height,
					"commit_info", extInfo,
					"err", err,
				)
				err = InvalidExtendedCommitInfoError{
					Err: err,
				}

				return &cometabci.ResponsePrepareProposal{Txs: make([][]byte, 0)}, err
			}

			// Create the vote extension injection data which will be injected into the proposal. These contain the
			// oracle data for the current block which will be committed to state in PreBlock.
			extInfoBz, err = h.extendedCommitCodec.Encode(extInfo)
			if err != nil {
				h.logger.Error(
					"failed to extended commit info",
					"commit_info", extInfo,
					"err", err,
				)
				err = types.CodecError{
					Err: err,
				}

				return &cometabci.ResponsePrepareProposal{Txs: make([][]byte, 0)}, err
			}
			// Inject our VE Tx to the req Txs, we want to do this before h.prepareProposalHandler(ctx, req) so that
			// the wrapped application can have access to the injected VE tx.
			extInfoBzSize := int64(len(extInfoBz))
			if extInfoBzSize < req.MaxTxBytes {
				req.Txs = append([][]byte{extInfoBz}, req.Txs...)
				// Reserve bytes for our VE Tx
				req.MaxTxBytes -= extInfoBzSize
			} else {
				h.logger.Error("omitting VE because size consumes entire block",
					"extInfoBzSize", extInfoBzSize,
					"MaxTxBytes", req.MaxTxBytes)
				extInfoBz = []byte{}
			}
		}

		// Build the proposal. Get the duration that the wrapped prepare proposal handler executed for.
		wrappedPrepareProposalStartTime := time.Now()
		resp, err = h.prepareProposalHandler(ctx, req)
		wrappedPrepareProposalLatency = time.Since(wrappedPrepareProposalStartTime)
		if err != nil {
			h.logger.Error("failed to prepare proposal", "err", err)
			err = types.WrappedHandlerError{
				Handler: servicemetrics.PrepareProposal,
				Err:     err,
			}

			return &cometabci.ResponsePrepareProposal{Txs: make([][]byte, 0)}, err
		}
		h.logger.Info("wrapped prepareProposalHandler produced response ", "txs", len(resp.Txs))

		// Inject our VE Tx ( if extInfoBz is non-empty), and resize our response Txs to respect req.MaxTxBytes
		resp.Txs = h.injectAndResize(resp.Txs, extInfoBz, req.MaxTxBytes+int64(len(extInfoBz)))

		h.logger.Info(
			"prepared proposal",
			"txs", len(resp.Txs),
			"vote_extensions_enabled", voteExtensionsEnabled,
		)

		return resp, nil
	}
}

// injectAndResize returns a tx array containing the injectTx at the beginning followed by appTxs.
// The returned transaction array is bounded by maxSizeBytes, and the function is idempotent meaning the
// injectTx will only appear once regardless of how many times you attempt to inject it.
// If injectTx is large enough, all originalTxs may end up being excluded from the returned tx array.
func (h *ProposalHandler) injectAndResize(appTxs [][]byte, injectTx []byte, maxSizeBytes int64) [][]byte {
	var returnedTxs [][]byte
	var consumedBytes int64
	// If VEs are enabled and our VE Tx isn't already in the appTxs, inject it here
	if len(injectTx) != 0 && (len(appTxs) < 1 || !bytes.Equal(appTxs[0], injectTx)) {
		injectBytes := int64(len(injectTx))
		// Ensure the VE Tx is in the response if we have room.
		// We may want to be more aggressive in the future about dedicating block space for application-specific Txs.
		// However, the VE Tx size should be relatively stable so MaxTxBytes should be set w/ plenty of headroom.
		if injectBytes <= maxSizeBytes {
			consumedBytes += injectBytes
			returnedTxs = append(returnedTxs, injectTx)
		}
	}
	// Add as many appTxs to the returned proposal as possible given our maxSizeBytes constraint
	for _, tx := range appTxs {
		consumedBytes += int64(len(tx))
		if consumedBytes > maxSizeBytes {
			return returnedTxs
		}
		returnedTxs = append(returnedTxs, tx)
	}
	return returnedTxs
}

// ProcessProposalHandler returns a ProcessProposalHandler that will be called
// by base app when a new block proposal needs to be verified. The ProcessProposalHandler
// will verify that the vote extensions included in the proposal are valid and compose
// a supermajority of signatures and vote extensions for the current block.
func (h *ProposalHandler) ProcessProposalHandler() sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *cometabci.RequestProcessProposal) (resp *cometabci.ResponseProcessProposal, err error) {
		start := time.Now()
		var wrappedProcessProposalLatency time.Duration

		// Defer a function to record the total time it took to process the proposal.
		defer func() {
			// record latency
			totalLatency := time.Since(start)
			h.logger.Info(
				"recording handle time metrics of process-proposal (seconds)",
				"total latency", totalLatency.Seconds(),
				"wrapped prepare proposal latency", wrappedProcessProposalLatency.Seconds(),
				"slinky prepare proposal latency", (totalLatency - wrappedProcessProposalLatency).Seconds(),
			)
			types.RecordLatencyAndStatus(h.metrics, totalLatency-wrappedProcessProposalLatency, err, servicemetrics.ProcessProposal)
		}()

		// this should never happen, but just in case
		if req == nil {
			h.logger.Error("ProcessProposalHandler received a nil request")
			err = types.NilRequestError{
				Handler: servicemetrics.ProcessProposal,
			}
			return nil, err
		}

		voteExtensionsEnabled := ve.VoteExtensionsEnabled(ctx)

		h.logger.Info(
			"processing proposal",
			"height", req.Height,
			"num_txs", len(req.Txs),
			"vote_extensions_enabled", voteExtensionsEnabled,
		)

		if voteExtensionsEnabled {
			// Ensure that the commit info was correctly injected into the proposal.
			if len(req.Txs) < NumInjectedTxs {
				h.logger.Error("failed to process proposal: missing commit info", "num_txs", len(req.Txs))
				err = MissingCommitInfoError{}
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT},
					err
			}

			// Validate the vote extensions included in the proposal.
			var extInfo cometabci.ExtendedCommitInfo
			extInfo, err = h.extendedCommitCodec.Decode(req.Txs[OracleInfoIndex])
			if err != nil {
				h.logger.Error("failed to unmarshal commit info", "err", err)
				err = types.CodecError{
					Err: err,
				}
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT},
					err
			}

			if err := h.ValidateExtendedCommitInfo(ctx, req.Height, extInfo); err != nil {
				h.logger.Error(
					"failed to validate vote extensions",
					"height", req.Height,
					"commit_info", extInfo,
					"err", err,
				)
				err = InvalidExtendedCommitInfoError{
					Err: err,
				}

				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT},
					err
			}

			// Process the transactions in the proposal with the oracle data removed.
			req.Txs = req.Txs[NumInjectedTxs:]
		}

		// call the wrapped process-proposal
		wrappedProcessProposalStartTime := time.Now()
		resp, err = h.processProposalHandler(ctx, req)
		if err != nil {
			err = types.WrappedHandlerError{
				Handler: servicemetrics.ProcessProposal,
				Err:     err,
			}
		}

		wrappedProcessProposalLatency = time.Since(wrappedProcessProposalStartTime)
		return resp, err
	}
}
