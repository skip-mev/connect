package base

import (
	"context"

	"go.uber.org/zap"
)

// listenOnConfigUpdater listens for updates from the config updater and updates the
// provider's internal configurations. This will trigger the provider to restart
// and is blocking until the context is cancelled.
func (p *Provider[K, V]) listenOnConfigUpdater(ctx context.Context) {
	if p.updater == nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("stopping config client listener")
			return
		case ids := <-p.updater.GetIDs():
			p.logger.Debug("received new ids", zap.Any("ids", ids))
			p.SetIDs(ids)

			// Signal the provider to restart.
			p.restartCh <- struct{}{}
		}
	}
}
