package oracle

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"go.uber.org/zap"
)

// listenForMarketMapUpdates is a goroutine that listens for market map updates and
// updates the orchestrated providers with the new market map. This method assumes a market map provider is present,
// so callers of this method must nil check the provider first.
func (o *OracleImpl) listenForMarketMapUpdates(ctx context.Context) {
	ticker := time.NewTicker(o.cfg.UpdateInterval)
	o.logger.Info("listening for market map updates", zap.Int("num_provider", len(o.mmProviders)))
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Fetch the latest market map.
			for _, state := range o.mmProviders {
				// Update the oracle with the latest market map iff the market map has changed.
				o.logger.Info("checking for market map updates", zap.String("provider", state.Provider.Name()))
				o.checkMarketMapUpdates(state)
			}
		}
	}
}

func (o *OracleImpl) checkMarketMapUpdates(state *MarketMapProviderState) {
	provider, marketmap := state.Provider, state.MarketMap

	// Fetch the latest market map(s) for the provider.
	for _, resp := range provider.GetData() {
		updated := resp.Value.MarketMap
		if err := updated.ValidateBasic(); err != nil {
			o.logger.Error("failed to validate market map", zap.Error(err), zap.String("provider", provider.Name()))
			return
		}

		if marketmap.Equal(updated) {
			o.logger.Debug("market map has not changed", zap.String("provider", provider.Name()))
			return
		}

		o.logger.Info("updating oracle with new market map")
		if err := o.UpdateMarketMap(updated); err != nil {
			o.logger.Error("failed to update oracle with new market map", zap.Error(err))
			return
		}

		// Write the market map to the configured path.
		if err := o.writeMarketMap(); err != nil {
			o.logger.Error("failed to write market map", zap.Error(err))
		}

		o.logger.Info("updated oracle with new market map", zap.Any("market_map", updated))
	}
}

// writeMarketMap writes the oracle's market map to the configured path.
func (o *OracleImpl) writeMarketMap() error {
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
