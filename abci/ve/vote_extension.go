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

	"github.com/skip-mev/connect/v2/abci/strategies/aggregator"
	compression "github.com/skip-mev/connect/v2/abci/strategies/codec"
	"github.com/skip-mev/connect/v2/abci/strategies/currencypair"
	connectabci "github.com/skip-mev/connect/v2/abci/types"
	"github.com/skip-mev/connect/v2/abci/ve/types"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	servicemetrics "github.com/skip-mev/connect/v2/service/metrics"
	servicetypes "github.com/skip-mev/connect/v2/service/servers/oracle/types"
)

// VoteExtensionHandler is a handler that extends a vote with the oracle's
// current price feed. In the case where oracle data is unable to be fetched
// or correctly marshalled, the handler will return an empty vote extension to
// ensure liveliness.
type VoteExtensionHandler struct {
	logger log.Logger

	// oracleClient is the remote oracle client that is responsible for fetching prices
	oracleClient connectabci.OracleClient

	// timeout is the maximum amount of time to wait for the oracle to respond
	// to a price request.
	timeout time.Duration

	// currencyPairStrategy is the strategy used to determine the price information
	// to include in the vote extension.
	currencyPairStrategy currencypair.CurrencyPairStrategy

	// voteExtensionCodec is an interface to handle the marshalling / unmarshalling of vote-extensions
	voteExtensionCodec compression.VoteExtensionCodec

	// PriceApplier is the price applier that is used to decode vote-extensions, aggregate price reports, and write prices to state.
	priceApplier aggregator.PriceApplier

	// metrics is the service metrics interface that the vote-extension handler will use to report metrics.
	metrics servicemetrics.Metrics
}

// NewVoteExtensionHandler returns a new VoteExtensionHandler.
func NewVoteExtensionHandler(
	logger log.Logger,
	oracleClient connectabci.OracleClient,
	timeout time.Duration,
	strategy currencypair.CurrencyPairStrategy,
	codec compression.VoteExtensionCodec,
	priceApplier aggregator.PriceApplier,
	metrics servicemetrics.Metrics,
) *VoteExtensionHandler {
	return &VoteExtensionHandler{
		logger:               logger,
		oracleClient:         oracleClient,
		timeout:              timeout,
		currencyPairStrategy: strategy,
		voteExtensionCodec:   codec,
		metrics:              metrics,
		priceApplier:         priceApplier,
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
			h.logger.Debug(
				"extend vote handler",
				"duration (seconds)", latency.Seconds(),
				"err", err,
			)
			connectabci.RecordLatencyAndStatus(h.metrics, latency, err, servicemetrics.ExtendVote)

			// ignore all non-panic errors
			var p ErrPanic
			if !errors.As(err, &p) {
				err = nil
			}
		}()

		if req == nil {
			h.logger.Error("extend vote handler received a nil request")
			err = connectabci.NilRequestError{
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
		if _, err = h.priceApplier.ApplyPricesFromVoteExtensions(ctx, reqFinalizeBlock); err != nil {
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

			err = connectabci.CodecError{
				Err: err,
			}

			return &cometabci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		h.logger.Debug(
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
			h.logger.Debug(
				"verify vote extension handler",
				"duration (seconds)", latency.Seconds(),
			)

			connectabci.RecordLatencyAndStatus(h.metrics, latency, err, servicemetrics.VerifyVoteExtension)
		}()

		if req == nil {
			err = connectabci.NilRequestError{
				Handler: servicemetrics.VerifyVoteExtension,
			}
			h.logger.Error("VerifyVoteExtensionHandler received a nil request")
			return nil, err
		}

		// By default, we accept empty vote extensions.
		if len(req.VoteExtension) == 0 {
			h.logger.Info(
				"empty vote extension",
				"height", req.Height,
			)

			return &cometabci.ResponseVerifyVoteExtension{Status: cometabci.ResponseVerifyVoteExtension_ACCEPT}, nil
		}

		// decode the vote-extension bytes
		voteExtension, err := h.voteExtensionCodec.Decode(req.VoteExtension)
		if err != nil {
			h.logger.Error(
				"failed to decode vote extension",
				"height", req.Height,
				"err", err,
			)
			err = connectabci.CodecError{
				Err: err,
			}

			return &cometabci.ResponseVerifyVoteExtension{Status: cometabci.ResponseVerifyVoteExtension_REJECT}, err
		}

		if err := ValidateOracleVoteExtension(ctx, voteExtension, h.currencyPairStrategy); err != nil {
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

		h.logger.Debug(
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
// does this by iterating over the prices submitted by the oracle service and determining the
// correct decoded price / ID based on the currency pair strategy.
func (h *VoteExtensionHandler) transformOracleServicePrices(ctx sdk.Context, prices map[string]string) (types.OracleVoteExtension, error) {
	strategyPrices := make(map[uint64][]byte)

	// Iterate over the prices and transform them into the correct format.
	for currencyPairID, priceString := range prices {
		cp, err := connecttypes.CurrencyPairFromString(currencyPairID)
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

		h.logger.Debug(
			"transformed oracle price",
			"currency_pair", cp,
			"height", ctx.BlockHeight(),
		)

		strategyPrices[cpID] = encodedPrice
	}

	h.logger.Debug("transformed oracle prices", "prices", len(strategyPrices))

	return types.OracleVoteExtension{
		Prices: strategyPrices,
	}, nil
}
