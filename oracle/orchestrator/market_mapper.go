package orchestrator

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// listenForMarketMapUpdates is a goroutine that listens for market map updates and
// updates the orchestrated providers with the new market map.
func (o *ProviderOrchestrator) listenForMarketMapUpdates(ctx context.Context) func() error {
	return func() error {
		mapper := o.GetMarketMapper()
		ids := mapper.GetIDs()
		if len(ids) != 1 {
			return fmt.Errorf("expected 1 id, got %d", len(ids))
		}

		apiCfg := mapper.GetAPIConfig()
		ticker := time.NewTicker(apiCfg.Interval)
		chain := ids[0]
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				// Fetch the latest market map.
				response := mapper.GetData()
				if response == nil {
					o.logger.Debug("market mapper returned nil response")
					continue
				}

				result, ok := response[chain]
				if !ok {
					o.logger.Debug("market mapper response missing chain", zap.Any("chain", chain))
					continue
				}

				// Update the orchestrator with the latest market map iff the market map has changed.
				updated := result.Value.MarketMap
				if o.marketMap.Equal(updated) {
					o.logger.Debug("market map has not changed")
					continue
				}

				o.logger.Info("updating orchestrator with new market map")
				if err := o.UpdateWithMarketMap(updated); err != nil {
					o.logger.Error("failed to update orchestrator with new market map", zap.Error(err))
					return err
				}
			}
		}
	}
}
