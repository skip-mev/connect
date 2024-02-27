package ve

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	compression "github.com/skip-mev/slinky/abci/strategies/codec"
	"github.com/skip-mev/slinky/abci/strategies/currencypair"
	slinkyabci "github.com/skip-mev/slinky/abci/types"
	"github.com/skip-mev/slinky/abci/ve/types"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	client "github.com/skip-mev/slinky/service/clients/oracle"
	servicemetrics "github.com/skip-mev/slinky/service/metrics"
	servicetypes "github.com/skip-mev/slinky/service/servers/oracle/types"
)

// VoteExtensionHandler is a handler that extends a vote with the oracle's
// current price feed. In the case where oracle data is unable to be fetched
// or correctly marshalled, the handler will return an empty vote extension to
// ensure liveliness.
type VoteExtensionHandler struct {
	logger log.Logger

	// oracleClient is the remote oracle client that is responsible for fetching prices
	oracleClient client.OracleClient

	// timeout is the maximum amount of time to wait for the oracle to respond
	// to a price request.
	timeout time.Duration

	// currencyPairStrategy is the strategy used to determine the price information
	// to include in the vote extension.
	currencyPairStrategy currencypair.CurrencyPairStrategy

	// voteExtensionCodec is an interface to handle the marshalling / unmarshalling of vote-extensions
	voteExtensionCodec compression.VoteExtensionCodec

	// preBlocker is utilized to update and retrieve the latest on-chain price information.
	preBlocker sdk.PreBlocker

	// metrics is the service metrics interface that the vote-extension handler will use to report metrics.
	metrics servicemetrics.Metrics
}

// NewVoteExtensionHandler returns a new VoteExtensionHandler.
func NewVoteExtensionHandler(
	logger log.Logger,
	oracleClient client.OracleClient,
	timeout time.Duration,
	strategy currencypair.CurrencyPairStrategy,
	codec compression.VoteExtensionCodec,
	preBlocker sdk.PreBlocker,
	metrics servicemetrics.Metrics,
) *VoteExtensionHandler {
	return &VoteExtensionHandler{
		logger:               logger,
		oracleClient:         oracleClient,
		timeout:              timeout,
		currencyPairStrategy: strategy,
		voteExtensionCodec:   codec,
		preBlocker:           preBlocker,
		metrics:              metrics,
	}
}

// ExtendVoteHandler returns a handler that extends a vote with the oracle's
// current price feed. In the case where oracle data is unable to be fetched
// or correctly marshalled, the handler will return an empty vote extension to
// ensure liveness.
func (h *VoteExtensionHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *cometabci.RequestExtendVote) (resp *cometabci.ResponseExtendVote, err error) {
		start := time.Now()

		// measure latencies from invocation to return, catch panics first
		defer func() {
			// catch panics if possible
			if r := recover(); r != nil {
				h.logger.Error(
					"recovered from panic in ExtendVoteHandler",
					"err", r,
				)

				resp, err = &cometabci.ResponseExtendVote{VoteExtension: []byte{}}, ErrPanic{fmt.Errorf("%v", r)}
			}

			// measure latency
			latency := time.Since(start)
			h.logger.Info(
				"extend vote handler",
				"duration (seconds)", latency.Seconds(),
				"err", err,
			)
			slinkyabci.RecordLatencyAndStatus(h.metrics, latency, err, servicemetrics.ExtendVote)

			// ignore all non-panic errors
			var p ErrPanic
			if !errors.As(err, &p) {
				err = nil
			}
		}()

		if req == nil {
			h.logger.Error("extend vote handler received a nil request")
			err = slinkyabci.NilRequestError{
				Handler: servicemetrics.ExtendVote,
			}
			return nil, err
		}

		// Update the latest on-chain prices with the vote extensions included in the current
		// block proposal.
		reqFinalizeBlock := &cometabci.RequestFinalizeBlock{
			Txs:    req.Txs,
			Height: req.Height,
		}
		if _, err = h.preBlocker(ctx, reqFinalizeBlock); err != nil {
			h.logger.Error(
				"failed to aggregate oracle votes",
				"height", req.Height,
				"err", err,
			)
			err = PreBlockError{err}

			return &cometabci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		// Create a context with a timeout to ensure we do not wait forever for the oracle
		// to respond.
		reqCtx, cancel := context.WithTimeout(ctx.Context(), h.timeout)
		defer cancel()

		// To ensure liveness, we return a vote even if the oracle is not running
		// or if the oracle returns a bad response.
		oracleResp, err := h.oracleClient.Prices(ctx.WithContext(reqCtx), &servicetypes.QueryPricesRequest{})
		if err != nil {
			h.logger.Error(
				"failed to retrieve oracle prices for vote extension; returning empty vote extension",
				"height", req.Height,
				"ctx_err", reqCtx.Err(),
				"err", err,
			)

			err = OracleClientError{
				Err: err,
			}

			return &cometabci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		// If we get no response, we return an empty vote extension.
		if oracleResp == nil {
			h.logger.Error(
				"oracle returned nil prices for vote extension; returning empty vote extension",
				"height", req.Height,
			)

			err = OracleClientError{fmt.Errorf("oracle returned nil prices")}

			return &cometabci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		// Transform the response prices into a vote extension.
		voteExt, err := h.transformOracleServicePrices(ctx, oracleResp.Prices)
		if err != nil {
			h.logger.Error(
				"failed to transform oracle prices for vote extension; returning empty vote extension",
				"height", req.Height,
				"err", err,
			)

			err = TransformPricesError{
				Err: err,
			}

			return &cometabci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		bz, err := h.voteExtensionCodec.Encode(voteExt)
		if err != nil {
			h.logger.Error(
				"failed to marshal vote extension; returning empty vote extension",
				"height", req.Height,
				"err", err,
			)

			err = slinkyabci.CodecError{
				Err: err,
			}

			return &cometabci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		h.logger.Info(
			"extending vote with oracle prices",
			"req_height", req.Height,
		)

		return &cometabci.ResponseExtendVote{VoteExtension: bz}, nil
	}
}

// VerifyVoteExtensionHandler returns a handler that verifies the vote extension provided by
// a validator is valid. In the case when the vote extension is empty, we return ACCEPT. This means
// that the validator may have been unable to fetch prices from the oracle and is voting an empty vote extension.
// We reject any vote extensions that are not empty and fail to unmarshal or contain invalid prices.
func (h *VoteExtensionHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(ctx sdk.Context, req *cometabci.RequestVerifyVoteExtension) (_ *cometabci.ResponseVerifyVoteExtension, err error) {
		start := time.Now()

		// measure latencies from invocation to return
		defer func() {
			latency := time.Since(start)
			h.logger.Info(
				"verify vote extension handler",
				"duration (seconds)", latency.Seconds(),
			)

			slinkyabci.RecordLatencyAndStatus(h.metrics, latency, err, servicemetrics.VerifyVoteExtension)
		}()

		if req == nil {
			err = slinkyabci.NilRequestError{
				Handler: servicemetrics.VerifyVoteExtension,
			}
			h.logger.Error("VerifyVoteExtensionHandler received a nil request")
			return nil, err
		}
		// decode the vote-extension bytes
		voteExtension, err := h.voteExtensionCodec.Decode(req.VoteExtension)
		if err != nil {
			h.logger.Error(
				"failed to decode vote extension",
				"height", req.Height,
				"err", err,
			)
			err = slinkyabci.CodecError{
				Err: err,
			}

			return &cometabci.ResponseVerifyVoteExtension{Status: cometabci.ResponseVerifyVoteExtension_REJECT}, err
		}

		if err := ValidateOracleVoteExtension(voteExtension, h.currencyPairStrategy); err != nil {
			h.logger.Error(
				"failed to validate vote extension",
				"height", req.Height,
				"err", err,
			)
			err = ValidateVoteExtensionError{
				Err: err,
			}

			return &cometabci.ResponseVerifyVoteExtension{Status: cometabci.ResponseVerifyVoteExtension_REJECT}, err
		}

		h.logger.Info(
			"validated vote extension",
			"height", req.Height,
			"size (bytes)", len(req.VoteExtension),
		)

		// observe message size
		h.metrics.ObserveMessageSize(servicemetrics.VoteExtension, len(req.VoteExtension))

		return &cometabci.ResponseVerifyVoteExtension{Status: cometabci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}

// transformOracleServicePrices transforms the oracle service prices into a vote extension. It
// does this by iterating over the the prices submitted by the oracle service and determining the
// correct decoded price / ID based on the currency pair strategy.
func (h *VoteExtensionHandler) transformOracleServicePrices(ctx sdk.Context, prices map[string]string) (types.OracleVoteExtension, error) {
	strategyPrices := make(map[uint64][]byte)

	// Iterate over the prices and transform them into the correct format.
	for currencyPairID, priceString := range prices {
		cp, err := slinkytypes.CurrencyPairFromString(currencyPairID)
		if err != nil {
			return types.OracleVoteExtension{}, err
		}

		rawPrice, converted := new(big.Int).SetString(priceString, 10)
		if !converted {
			return types.OracleVoteExtension{}, fmt.Errorf("failed to convert price string to big.Int: %s", priceString)
		}

		// Determine if the currency pair is supported by the network.
		cpID, err := h.currencyPairStrategy.ID(ctx, cp)
		if err != nil {
			h.logger.Debug(
				"failed to get currency pair ID",
				"currency_pair", cp,
				"err", err,
			)

			continue
		}

		// Determine the encoded price for the currency pair based on the strategy.
		encodedPrice, err := h.currencyPairStrategy.GetEncodedPrice(ctx, cp, rawPrice)
		if err != nil {
			h.logger.Debug(
				"failed to get current price for currency pair",
				"currency_pair", cp,
				"err", err,
			)

			continue
		}

		h.logger.Info(
			"transformed oracle price",
			"currency_pair", cp,
			"height", ctx.BlockHeight(),
		)

		strategyPrices[cpID] = encodedPrice
	}

	h.logger.Info("transformed oracle prices", "prices", len(strategyPrices))

	return types.OracleVoteExtension{
		Prices: strategyPrices,
	}, nil
}
