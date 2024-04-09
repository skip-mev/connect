package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// APIQueryHandler is an interface that encapsulates querying a data provider for info.
// The handler must respect the context timeout and cancel the request if the context
// is cancelled. All responses must be sent to the response channel. These are processed
// asynchronously by the provider.
//
//go:generate mockery --name APIQueryHandler --output ./mocks/ --case underscore
type APIQueryHandler[K providertypes.ResponseKey, V providertypes.ResponseValue] interface {
	Query(
		ctx context.Context,
		ids []K,
		responseCh chan<- providertypes.GetResponse[K, V],
	)
}

// APIFetcher is an interface that encapsulates fetching data from a provider. This interface
// is meant to abstract over the various processes of interacting w/ GRPC, JSON-RPC, REST, etc. APIs.
type APIFetcher[K providertypes.ResponseKey, V providertypes.ResponseValue] interface {
	// Fetch fetches data from the API for the given IDs. The response is returned as a map of IDs to
	// their respective responses. The request should respect the context timeout and cancel the request
	// if the context is cancelled.
	Fetch(
		ctx context.Context,
		ids []K,
	) providertypes.GetResponse[K, V]
}

// APIQueryHandlerImpl is the default API implementation of the QueryHandler interface.
// This is used to query using http requests. It manages querying the data provider
// by using the APIDataHandler and RequestHandler. All responses are sent to the
// response channel. In the case where the APIQueryHandler is atomic, the handler
// will make a single request for all IDs. If the APIQueryHandler is not
// atomic, the handler will make a request for each ID in a separate go routine.
type APIQueryHandlerImpl[K providertypes.ResponseKey, V providertypes.ResponseValue] struct {
	logger  *zap.Logger
	metrics metrics.APIMetrics
	config  config.APIConfig

	// fetcher is responsible for fetching data from the API.
	fetcher APIFetcher[K, V]
}

// NewAPIQueryHandler creates a new APIQueryHandler. It manages querying the data
// provider by using the APIDataHandler and RequestHandler.
func NewAPIQueryHandler[K providertypes.ResponseKey, V providertypes.ResponseValue](
	logger *zap.Logger,
	cfg config.APIConfig,
	requestHandler RequestHandler,
	apiHandler APIDataHandler[K, V],
	metrics metrics.APIMetrics,
) (APIQueryHandler[K, V], error) {
	fetcher, err := NewRestAPIFetcher(requestHandler, apiHandler, metrics, cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create api fetcher: %w", err)
	}

	return &APIQueryHandlerImpl[K, V]{
		logger:  logger.With(zap.String("api_query_handler", cfg.Name)),
		config:  cfg,
		metrics: metrics,
		fetcher: fetcher,
	}, nil
}

// NewAPIQueryHandlerWithFetcher creates a new APIQueryHandler with a custom api fetcher.
func NewAPIQueryHandlerWithFetcher[K providertypes.ResponseKey, V providertypes.ResponseValue](
	logger *zap.Logger,
	cfg config.APIConfig,
	fetcher APIFetcher[K, V],
	metrics metrics.APIMetrics,
) (APIQueryHandler[K, V], error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid provider config: %w", err)
	}

	if !cfg.Enabled {
		return nil, fmt.Errorf("api query handler is not enabled for the provider")
	}

	if logger == nil {
		return nil, fmt.Errorf("no logger specified for api query handler")
	}

	if metrics == nil {
		return nil, fmt.Errorf("no metrics specified for api query handler")
	}

	if fetcher == nil {
		return nil, fmt.Errorf("no fetcher specified for api query handler")
	}

	return &APIQueryHandlerImpl[K, V]{
		logger:  logger.With(zap.String("api_query_handler", cfg.Name)),
		config:  cfg,
		metrics: metrics,
		fetcher: fetcher,
	}, nil
}

// Query is used to query the API data provider for the given IDs. This method blocks
// until all responses have been sent to the response channel. Query will only
// make N concurrent requests at a time, where N is the capacity of the response channel.
func (h *APIQueryHandlerImpl[K, V]) Query(
	ctx context.Context,
	ids []K,
	responseCh chan<- providertypes.GetResponse[K, V],
) {
	if len(ids) == 0 {
		h.logger.Debug("no ids to query")
		return
	}

	// Observe the total amount of time it takes to fulfill the request(s).
	h.logger.Debug("starting api query handler")
	defer func() {
		if r := recover(); r != nil {
			h.logger.Error("panic in api query handler", zap.Any("panic", r))
		}

		h.logger.Debug("finished api query handler")
	}()

	// Set the concurrency limit based on the maximum number of queries allowed for a single
	// interval.
	wg := errgroup.Group{}
	limit := math.Min(h.config.MaxQueries, len(ids))
	wg.SetLimit(limit)
	h.logger.Debug("setting concurrency limit", zap.Int("limit", limit))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// If our task is atomic, we can make a single request for all the IDs. Otherwise,
	// we need to make a request for each ID.
	var tasks []func() error
	if h.config.Atomic {
		tasks = append(tasks, h.subTask(ctx, ids, responseCh))
	} else {
		for i := 0; i < len(ids); i++ {
			id := ids[i]
			tasks = append(tasks, h.subTask(ctx, []K{id}, responseCh))
		}
	}

	// Block each task until the wait group has capacity to accept a new response.
	index := 0
MainLoop:
	for {
		select {
		case <-ctx.Done():
			h.logger.Debug("context cancelled, stopping queries")
			break MainLoop
		default:
			wg.Go(tasks[index])
			index++
			index %= len(tasks)

			// Sleep for a bit to prevent the loop from spinning too fast.
			h.logger.Debug("sleeping", zap.Duration("interval", h.config.Interval), zap.Int("index", index))
			time.Sleep(h.config.Interval)
		}
	}

	// Wait for all tasks to complete.
	h.logger.Debug("waiting for api sub-tasks to complete")
	if err := wg.Wait(); err != nil {
		h.logger.Error("error querying ids", zap.Error(err))
	}
	h.logger.Debug("all api sub-tasks completed")
}

// subTask is the subtask that is used to query the data provider for the given IDs,
// parse the response, and write the response to the response channel.
func (h *APIQueryHandlerImpl[K, V]) subTask(
	ctx context.Context,
	ids []K,
	responseCh chan<- providertypes.GetResponse[K, V],
) func() error {
	return func() error {
		start := time.Now().UTC()

		defer func() {
			// Recover from any panics that occur.
			if r := recover(); r != nil {
				h.logger.Error("panic occurred in subtask", zap.Any("panic", r), zap.Any("ids", ids))
			}

			h.metrics.ObserveProviderResponseLatency(h.config.Name, time.Since(start))
			h.logger.Debug("finished subtask", zap.Any("ids", ids))
		}()

		h.logger.Debug("starting subtask", zap.Any("ids", ids))

		h.writeResponse(ctx, responseCh, h.fetcher.Fetch(ctx, ids))
		return nil
	}
}

// writeResponse is used to write the response to the response channel.
func (h *APIQueryHandlerImpl[K, V]) writeResponse(
	ctx context.Context,
	responseCh chan<- providertypes.GetResponse[K, V],
	response providertypes.GetResponse[K, V],
) {
	// Write the response to the response channel. We only do so if the
	// context has not been cancelled. Otherwise, we risk writing to a
	// channel that is not being read from.
	select {
	case <-ctx.Done():
		h.logger.Debug("context cancelled, stopping write response")
		return
	case responseCh <- response:
		h.logger.Debug("wrote response", zap.String("response", response.String()))
	}

	// Update the metrics.
	for id := range response.Resolved {
		h.metrics.AddProviderResponse(h.config.Name, strings.ToLower(id.String()), providertypes.OK)
	}
	for id, unresolvedResult := range response.UnResolved {
		h.metrics.AddProviderResponse(h.config.Name, strings.ToLower(id.String()), unresolvedResult.Code())
	}
}
