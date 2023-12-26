package base

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// loop is the main loop for the provider. It continuously attempts to request data
// from the APIDataHandler until the context is cancelled.
func (p *BaseProvider[K, V]) loop(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ticker := time.NewTicker(p.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("provider stopped via context")
			return ctx.Err()

		case <-ticker.C:
			p.logger.Debug("attempting to fetch new data")
			p.attemptDataUpdate(ctx)
		}
	}
}

// attemptDataUpdate tries to update data by fetching and parsing API data.
// It logs any errors encountered during the process.
func (p *BaseProvider[K, V]) attemptDataUpdate(ctx context.Context) {
	fetchCtx, cancel := context.WithTimeout(ctx, p.cfg.Timeout)
	defer cancel()

	// Retrieve API Data.
	data, err := p.handler.Get(fetchCtx)
	if err != nil {
		p.logger.Debug("failed to fetch data from API", zap.Error(err))
		return
	}

	if len(data) == 0 {
		p.logger.Debug("no data returned from API")
		return
	}

	// Update the data.
	p.setData(data)
	p.setLastUpdate(time.Now().UTC())
	p.logger.Debug("data updated successfully", zap.Int("num_data_points", len(data)))
}

// setData sets the latest data for the provider.
func (p *BaseProvider[K, V]) setData(data map[K]V) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.data = data
}

// setLastUpdate sets the time at which the data was last updated.
func (p *BaseProvider[K, V]) setLastUpdate(t time.Time) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.lastUpdate = t
}
