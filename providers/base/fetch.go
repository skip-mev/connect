package base

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// loop is the main loop for the provider. It continuously attempts to request data
// from the APIDataHandler until the context is cancelled.
func (p *BaseProvider[K, V]) loop(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ticker := time.NewTicker(p.cfg.Interval)
	defer ticker.Stop()

	// responseCh is used to receive the response(s) from the query handler. The buffer size is set
	// to the minimum of the number of IDs and the max number of queries. This is to ensure that
	// the response channel does not block the query handler and that the query handler does not
	// exceed the rate limit parameters of the provider.
	responseCh := make(chan providertypes.GetResponse[K, V], math.Min(len(p.ids), p.cfg.MaxQueries))

	// Start the receive loop.
	go p.recv(ctx, responseCh)

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

			p.attemptDataUpdate(ctx, responseCh)
		}
	}
}

// attemptDataUpdate tries to update data by fetching and parsing API data.
// It logs any errors encountered during the process.
func (p *BaseProvider[K, V]) attemptDataUpdate(ctx context.Context, responseCh chan<- providertypes.GetResponse[K, V]) {
	if len(p.ids) == 0 {
		p.logger.Debug("no ids to fetch")
		return
	}

	ctx, cancel := context.WithTimeout(ctx, p.cfg.Timeout)
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
		p.handler.Query(ctx, p.ids, responseCh)
	}()
}

// recv receives responses from the response channel and updates the data.
func (p *BaseProvider[K, V]) recv(ctx context.Context, responseCh <-chan providertypes.GetResponse[K, V]) {
	// Wait for the data to be retrieved until the context is cancelled.
	for {
		select {
		case <-ctx.Done():
			p.logger.Debug("finishing recv and closing with request context err", zap.Error(ctx.Err()))
			return
		case r := <-responseCh:
			resolved, unResolved := r.Resolved, r.UnResolved

			// Update all of the resolved data.
			for id, result := range resolved {
				p.logger.Debug(
					"successfully fetched data",
					zap.Any("id", id),
					zap.String("result", result.String()),
				)

				p.updateData(id, result)
			}

			// Log all of the unresolved data.
			for id, err := range unResolved {
				p.logger.Debug(
					"failed to fetch data",
					zap.Any("id", id),
					zap.Error(err),
				)
			}
		}
	}
}

// updateData sets the latest data for the provider. This will only update the data if the timestamp
// of the data is greater than the current data.
func (p *BaseProvider[K, V]) updateData(id K, result providertypes.Result[V]) {
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
