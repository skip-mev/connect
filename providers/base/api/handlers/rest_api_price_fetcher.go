package handlers

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/api/errors"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// RestAPIFetcher handles the logic of fetching prices from a REST API. This implementation
// depends on an APIDataHandler to handle the creation of URLs / parsing the API response.
type RestAPIFetcher[K providertypes.ResponseKey, V providertypes.ResponseValue] struct {
	// requestHandler is responsible for making outgoing HTTP requests with a given URL.
	requestHandler RequestHandler

	// apiDataHandler is responsible for creating URLs and parsing the API response.
	apiDataHandler APIDataHandler[K, V]

	// metrics is responsible for tracking metrics related to the API.
	metrics metrics.APIMetrics

	// config is the configuration for the API. Specifically configuring the timeouts
	// for outgoing requests
	config config.APIConfig

	// logger
	logger *zap.Logger
}

// NewRestAPIFetcher creates a new RestAPIFetcher.
func NewRestAPIFetcher[K providertypes.ResponseKey, V providertypes.ResponseValue](
	requestHandler RequestHandler,
	apiDataHandler APIDataHandler[K, V],
	metrics metrics.APIMetrics,
	config config.APIConfig,
	logger *zap.Logger,
) (*RestAPIFetcher[K, V], error) {
	if err := config.ValidateBasic(); err != nil {
		return nil, err
	}

	if !config.Enabled {
		return nil, fmt.Errorf("api is disabled")
	}

	if requestHandler == nil {
		return nil, fmt.Errorf("request handler is nil")
	}

	if apiDataHandler == nil {
		return nil, fmt.Errorf("api data handler is nil")
	}

	if metrics == nil {
		return nil, fmt.Errorf("metrics is nil")
	}

	return &RestAPIFetcher[K, V]{
		requestHandler: requestHandler,
		apiDataHandler: apiDataHandler,
		metrics:        metrics,
		config:         config,
		logger:         logger,
	}, nil
}

// Fetch is used to fetch the corresponding IDs from the API. This method blocks until the
// response is received from the API, parsed, and returned.
func (pf *RestAPIFetcher[K, V]) Fetch(
	ctx context.Context,
	ids []K,
) providertypes.GetResponse[K, V] {
	// Create the URL for the request.
	url, err := pf.apiDataHandler.CreateURL(ids)
	if err != nil {
		return providertypes.NewGetResponseWithErr[K, V](
			ids,
			providertypes.NewErrorWithCode(
				errors.ErrCreateURLWithErr(err),
				providertypes.ErrorUnableToCreateURL,
			),
		)
	}

	pf.logger.Debug("created url", zap.String("url", url))

	// Make the request.
	apiCtx, cancel := context.WithTimeout(ctx, pf.config.Timeout)
	defer cancel()

	pf.logger.Debug("making request", zap.String("url", url))

	// Record the status code in the metrics.
	resp, err := pf.requestHandler.Do(apiCtx, url)
	pf.metrics.AddHTTPStatusCode(pf.config.Name, resp)
	if err != nil {
		status := providertypes.ErrorUnknown
		if resp != nil {
			status = providertypes.ErrorCode(resp.StatusCode)
		}

		pf.logger.Error(
			"failed to make request",
			zap.Error(err),
			zap.String("url", url),
		)

		return providertypes.NewGetResponseWithErr[K, V](
			ids,
			providertypes.NewErrorWithCode( // TODO(nikhil): coordinate api-errors w/ correct metric codes
				errors.ErrDoRequestWithErr(err),
				status,
			),
		)
	}
	defer resp.Body.Close()

	pf.logger.Debug("received response", zap.Int("status_code", resp.StatusCode))
	// TODO: add more error handling here.
	// TODO(nikhil): move this logic to a shared HTTPClient
	var response providertypes.GetResponse[K, V]
	switch {
	case resp.StatusCode == http.StatusTooManyRequests:
		response = providertypes.NewGetResponseWithErr[K, V](
			ids,
			providertypes.NewErrorWithCode(
				errors.ErrRateLimit,
				providertypes.ErrorRateLimitExceeded,
			),
		)
	case resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices:
		response = providertypes.NewGetResponseWithErr[K, V](
			ids,
			providertypes.NewErrorWithCode(
				errors.ErrUnexpectedStatusCodeWithCode(resp.StatusCode),
				providertypes.ErrorCode(resp.StatusCode),
			),
		)
	default:
		response = pf.apiDataHandler.ParseResponse(ids, resp)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		pf.logger.Error(
			"failed to make and parse response",
			zap.Int("status_code", resp.StatusCode),
			zap.String("url", url),
		)
	}

	return response
}
