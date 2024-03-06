package base

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"go.uber.org/zap"

	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// fetch is the main blocker for the provider. It is responsible for fetching data from
// the data provider and updating the data. Note that the context passed here is valid
// until either the parent context (provider's main context) is cancelled, the fetch routine
// encounters an error, or the provider is manually stopped.
func (p *Provider[K, V]) fetch(ctx context.Context) error {
	// Determine which loop to use based on whether the provider is an API or webSocket provider.
	switch {
	case p.Type() == providertypes.API:
		return p.startAPI(ctx)
	case p.Type() == providertypes.WebSockets:
		return p.startMultiplexWebsocket(ctx)
	default:
		return fmt.Errorf("no api or websocket configured")
	}
}

// startAPI is the main loop for the provider. It is responsible for fetching data from the API
// and updating the data.
func (p *Provider[K, V]) startAPI(ctx context.Context) error {
	p.logger.Info("starting api query handler")

	// Start the data update loop.
	handler := p.GetAPIHandler()
	ids := p.GetIDs()
	for {
		select {
		case <-ctx.Done():
			p.logger.Info("api stopped via context")
			return ctx.Err()

		default:
			p.logger.Debug(
				"attempting to fetch new data",
				zap.Int("buffer_size", len(p.responseCh)),
				zap.Int("num_ids", len(ids)),
			)

			handler.Query(ctx, ids, p.responseCh)
		}
	}
}

// startMultiplexWebsocket is the main loop for web socket providers. It is responsible for
// creating a connection to the websocket and handling the incoming messages. In the case
// where multiple connections (multiplexing) are used, this function will start multiple
// connections.
func (p *Provider[K, V]) startMultiplexWebsocket(ctx context.Context) error {
	var (
		maxSubsPerConn = p.wsCfg.MaxSubscriptionsPerConnection
		subTasks       = make([][]K, 0)
		wg             = errgroup.Group{}
	)

	// create sub handlers
	// if len(ids) == 30 and MaxSubscriptionsPerConnection == 45
	// 30 / 45 = 0 -> need one sub handler
	ids := p.GetIDs()
	if maxSubsPerConn > 0 {
		// case where we will split ID's across sub handlers
		numSubHandlers := int(math.Ceil(float64(len(ids)) / float64(maxSubsPerConn)))
		p.logger.Info("setting number of web socket handlers for provider", zap.Int("sub_handlers", numSubHandlers))
		wg.SetLimit(numSubHandlers)

		// split ids
		for i := 0; i < numSubHandlers; i++ {
			start := i * maxSubsPerConn

			// Copy the IDs over.
			subIDs := make([]K, 0)
			if end := start + maxSubsPerConn; end >= len(ids) {
				subIDs = append(subIDs, ids[start:]...)
			} else {
				subIDs = append(subIDs, ids[start:end]...)
			}

			subTasks = append(subTasks, subIDs)
		}
	} else {
		// case where there is 1 sub handler
		subTasks = append(subTasks, ids)
		wg.SetLimit(1)
	}

	for _, subIDs := range subTasks {
		wg.Go(p.startWebSocket(ctx, subIDs))
	}

	// Wait for all the sub handlers to finish.
	return wg.Wait()
}

// startWebSocket starts a connection to the websocket and handles the incoming messages.
func (p *Provider[K, V]) startWebSocket(ctx context.Context, subIDs []K) func() error {
	return func() error {
		// Start the websocket query handler. If the connection fails to start, then the query handler
		// will be restarted after a timeout.
		restarts := 0
		handler := p.GetWebSocketHandler()
		handler = handler.Copy()
    
		for {
			select {
			case <-ctx.Done():
				p.logger.Info("web socket stopped via context")
				return ctx.Err()
			default:
				if restarts > 0 {
					p.logger.Info("rgestarting websocket query handler", zap.Int("num_restarts", restarts))

					// If the websocket query handler returns, then the connection was closed. Wait for
					// a bit before trying to reconnect.
					time.Sleep(p.wsCfg.ReconnectionTimeout)
				}

				p.logger.Debug("starting websocket query handler", zap.Int("num_ids", len(subIDs)), zap.Any("ids", subIDs))
				if err := handler.Start(ctx, subIDs, responseCh); err != nil {
					p.logger.Error("websocket query handler returned error", zap.Error(err))
				}
				restarts++
			}
		}
	}
}

// recv receives responses from the response channel and updates the data.
func (p *Provider[K, V]) recv(ctx context.Context) {
	// Wait for the data to be retrieved until the context is cancelled.
	for {
		select {
		case <-ctx.Done():
			p.logger.Debug("finishing recv and closing with request context err", zap.Error(ctx.Err()))
			return
		case r := <-p.responseCh:
			resolved, unResolved := r.Resolved, r.UnResolved

			// Update all the resolved data.
			for id, result := range resolved {
				p.logger.Debug(
					"successfully fetched data",
					zap.String("id", id.String()),
					zap.String("result", result.String()),
				)

				p.updateData(id, result)

				// Update the metrics.
				strID := strings.ToLower(id.String())
				p.metrics.AddProviderResponseByID(p.name, strID, providermetrics.Success, p.Type())
				p.metrics.AddProviderResponse(p.name, providermetrics.Success, p.Type())
				p.metrics.LastUpdated(p.name, strID, p.Type())
			}

			// Log and record all the unresolved data.
			for id, err := range unResolved {
				p.logger.Debug(
					"failed to fetch data",
					zap.Any("id", id),
					zap.Error(err),
				)

				// Update the metrics.
				strID := strings.ToLower(id.String())
				p.metrics.AddProviderResponseByID(p.name, strID, providermetrics.Failure, p.Type())
				p.metrics.AddProviderResponse(p.name, providermetrics.Failure, p.Type())
			}
		}
	}
}

// updateData sets the latest data for the provider. This will only update the data if the timestamp
// of the data is greater than the current data.
func (p *Provider[K, V]) updateData(id K, result providertypes.Result[V]) {
	p.mu.Lock()
	defer p.mu.Unlock()

	current, ok := p.data[id]
	if !ok {
		p.data[id] = result
		return
	}

	// If the timestamp of the result is less than the current timestamp, then we do not update the data.
	if result.Timestamp.Before(current.Timestamp) {
		p.logger.Debug(
			"result timestamp is before current timestamp",
			zap.Time("result_timestamp", result.Timestamp),
			zap.Time("current_timestamp", p.data[id].Timestamp),
			zap.String("id", fmt.Sprint(id)),
		)
		return
	}

	p.logger.Debug(
		"updating base provider data",
		zap.String("id", fmt.Sprint(id)),
		zap.String("result", result.String()),
	)
	p.data[id] = result
}
