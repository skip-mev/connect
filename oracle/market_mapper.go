package oracle

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

// listenForMarketMapUpdates is a goroutine that listens for market map updates and
// updates the orchestrated providers with the new market map. This method assumes a market map provider is present,
// so callers of this method must nil check the provider first.
func (o *OracleImpl) listenForMarketMapUpdates(ctx context.Context) {
	mmProvider := o.mmProvider
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

			if o.lastUpdated != 0 && o.lastUpdated == result.Value.LastUpdated {
				o.logger.Debug("skipping market map update on no lastUpdated change", zap.Uint64("lastUpdated", o.lastUpdated))
				continue
			}

			validSubset, err := result.Value.MarketMap.GetValidSubset()
			if err != nil {
				o.logger.Error("failed to validate market map", zap.Error(err))
				continue
			}

			// Detect removed markets and surface info about the removals
			var removedMarkets []string
			for t := range result.Value.MarketMap.Markets {
				if _, in := validSubset.Markets[t]; !in {
					removedMarkets = append(removedMarkets, t)
				}
			}
			if len(validSubset.Markets) == 0 || len(validSubset.Markets) != len(result.Value.MarketMap.Markets) {
				o.logger.Warn("invalid market map update has caused some markets to be removed")
				o.logger.Info("markets removed from invalid market map", zap.String("markets", strings.Join(removedMarkets, " ")))
			}

			// Update the oracle with the latest market map iff the market map has changed.
			updated := validSubset
			if o.marketMap.Equal(updated) {
				o.logger.Debug("market map has not changed")
				continue
			}

			o.logger.Info("updating oracle with new market map")
			if err := o.UpdateMarketMap(updated); err != nil {
				o.logger.Error("failed to update oracle with new market map", zap.Error(err))
				continue
			}

			o.lastUpdated = result.Value.GetLastUpdated()

			// Write the market map to the configured path.
			if err := o.WriteMarketMap(); err != nil {
				o.logger.Error("failed to write market map", zap.Error(err))
			}

			o.logger.Info("updated oracle with new market map")
			o.logger.Debug("updated oracle with new market map", zap.Any("market_map", updated))
		}
	}
}

// WriteMarketMap writes the oracle's market map to the configured path.
func (o *OracleImpl) WriteMarketMap() error {
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
