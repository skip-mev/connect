package orchestrator

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"go.uber.org/zap"
)

// listenForMarketMapUpdates is a goroutine that listens for market map updates and
// updates the orchestrated providers with the new market map.
func (o *ProviderOrchestrator) listenForMarketMapUpdates(ctx context.Context) {
	mmProvider := o.GetMarketMapProvider()
	ids := mmProvider.GetIDs()
	if len(ids) != 1 {
		o.logger.Error("market map provider can only be responsible for one chain", zap.Any("ids", ids))
		return
	}

	apiCfg := mmProvider.GetAPIConfig()
	ticker := time.NewTicker(apiCfg.Interval)
	chain := ids[0]
	o.logger.Info("listening for market map updates", zap.String("chain", chain.String()))
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Fetch the latest market map.
			response := mmProvider.GetData()
			if response == nil {
				o.logger.Info("market map provider returned nil response")
				continue
			}

			result, ok := response[chain]
			if !ok {
				o.logger.Debug("market map provider response missing chain", zap.Any("chain", chain))
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
				continue
			}

			// Write the market map to the configured path.
			if err := o.WriteMarketMap(); err != nil {
				o.logger.Error("failed to write market map", zap.Error(err))
			}

			o.logger.Info("updated orchestrator with new market map", zap.Any("market_map", updated))
		}
	}
}

// WriteMarketMap writes the orchestrator's market map to the configured path.
func (o *ProviderOrchestrator) WriteMarketMap() error {
	if len(o.writeTo) == 0 {
		return nil
	}

	o.mut.Lock()
	defer o.mut.Unlock()

	f, err := os.Create(o.writeTo)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(o.marketMap); err != nil {
		return err
	}

	o.logger.Debug("wrote market map to file", zap.String("path", o.writeTo))
	return nil
}
