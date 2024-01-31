package base

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/pkg/math"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// fetch is the main blocker for the provider. It is responsible for fetching data from the
// data provider and updating the data.
func (p *Provider[K, V]) fetch(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// responseCh is used to receive the response(s) from the query handler.
	var responseCh chan providertypes.GetResponse[K, V]
	switch {
	case p.api != nil:
		// The buffer size is set to the minimum of the number of IDs and the max number of queries.
		// This is to ensure that the response channel does not block the query handler and that the
		// query handler does not exceed the rate limit parameters of the provider.
		responseCh = make(chan providertypes.GetResponse[K, V], math.Min(len(p.ids), p.apiCfg.MaxQueries))
	case p.ws != nil:
		// Otherwise, the buffer size is set to the max buffer size configured for the websocket.
		responseCh = make(chan providertypes.GetResponse[K, V], p.wsCfg.MaxBufferSize)
	default:
		return fmt.Errorf("no api or websocket configured")
	}

	// Start the receive loop.
	go p.recv(ctx, responseCh)

	// Determine which loop to use based on whether the provider is an API or webSocket provider.
	switch {
	case p.api != nil:
		return p.startAPI(ctx, responseCh)
	case p.ws != nil:
		return p.startWebSocket(ctx, responseCh)
	default:
		return fmt.Errorf("no api or websocket configured")
	}
}

// startAPI is the main loop for the provider. It is responsible for fetching data from the API
// and updating the data.
func (p *Provider[K, V]) startAPI(ctx context.Context, responseCh chan<- providertypes.GetResponse[K, V]) error {
	p.logger.Info("starting api query handler")

	ticker := time.NewTicker(p.apiCfg.Interval)
	defer ticker.Stop()

	// Start the data update loop.
	for {
		select {
		case <-ctx.Done():
			p.logger.Info("provider stopped via context")
			return ctx.Err()

		case <-ticker.C:
			p.logger.Debug(
				"attempting to fetch new data",
				zap.Int("num_ids", len(p.ids)),
				zap.Int("buffer_size", len(responseCh)),
			)

			p.attemptAPIDataUpdate(ctx, responseCh)
		}
	}
}

// attemptAPIDataUpdate tries to update data by fetching and parsing API data.
// It logs any errors encountered during the process.
func (p *Provider[K, V]) attemptAPIDataUpdate(ctx context.Context, responseCh chan<- providertypes.GetResponse[K, V]) {
	if len(p.ids) == 0 {
		p.logger.Debug("no ids to fetch")
		return
	}

	ctx, cancel := context.WithTimeout(ctx, p.apiCfg.Timeout)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				p.logger.Error("panic in query handler", zap.Any("panic", r))
			}
			cancel()
			p.logger.Debug("finished query handler")
		}()

		// Start the query handler. The handler must respect the context timeout.
		p.logger.Debug("starting query handler")
		p.api.Query(ctx, p.ids, responseCh)
	}()
}

// startWebSocket starts a connection to the websocket and handles the incoming messages.
func (p *Provider[K, V]) startWebSocket(ctx context.Context, responseCh chan<- providertypes.GetResponse[K, V]) error {
	// Start the websocket query handler. If the connection fails to start, then the query handler
	// will be restarted after a timeout.
	for {
		select {
		case <-ctx.Done():
			p.logger.Info("provider stopped via context")
			return ctx.Err()
		default:
			p.logger.Debug("starting websocket query handler")
			// create sub handlers
			// if len(ids) == 30 and MaxSubscriptionsPerConnection == 45
			// 30 / 45 = 0 -> need one sub handler
			maxSubsPerConn := p.wsCfg.MaxSubscriptionsPerConnection
			if maxSubsPerConn > 0 {
				// case where we will split ID's across sub handlers
				numSubHandlers := (len(p.ids) / maxSubsPerConn) + 1
				// split ids
				var subIDs []K
				for i := 0; i < numSubHandlers; i++ {
					start := i
					end := maxSubsPerConn * (i + 1)
					if i+1 == numSubHandlers {
						subIDs = p.ids[start:]

					} else {
						subIDs = p.ids[start:end]
					}

					// spin up a goroutine for parallel handlers
					go func(ids []K) {
						if err := p.ws.Start(ctx, ids, responseCh); err != nil {
							p.logger.Error("websocket query handler returned error", zap.Error(err))
						}
					}(subIDs)
				}
			} else {
				// case where there is 1 sub handler
				if err := p.ws.Start(ctx, p.ids, responseCh); err != nil {
					p.logger.Error("websocket query handler returned error", zap.Error(err))
				}
			}

			// If the websocket query handler returns, then the connection was closed. Wait for
			// a bit before trying to reconnect.
			time.Sleep(p.wsCfg.ReconnectionTimeout)
		}
	}
}

// recv receives responses from the response channel and updates the data.
func (p *Provider[K, V]) recv(ctx context.Context, responseCh <-chan providertypes.GetResponse[K, V]) {
	// Wait for the data to be retrieved until the context is cancelled.
	for {
		select {
		case <-ctx.Done():
			p.logger.Debug("finishing recv and closing with request context err", zap.Error(ctx.Err()))
			return
		case r := <-responseCh:
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
