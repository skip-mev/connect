package oracle

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"

	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// IsMarketMapValidUpdated checks if the given MarketMapResponse is an update to the existing MarketMap.
// - returns an error if the market map is fully invalid or the response is invalid
// - returns false if the market map is not updated
// - returns true and the new market map if the new market map is updated and valid.
func (o *OracleImpl) IsMarketMapValidUpdated(resp *mmtypes.MarketMapResponse) (mmtypes.MarketMap, bool, error) {
	if resp == nil {
		return mmtypes.MarketMap{}, false, fmt.Errorf("nil response")
	}

	// TODO: restore LastUpdated check when on-chain logic is fixed

	// check equality of the response and our current market map
	if o.marketMap.Equal(resp.MarketMap) {
		o.logger.Info("market map has not changed")
		return mmtypes.MarketMap{}, false, nil
	}

	// if the value has changed, check for a Valid subset
	validSubset, err := resp.MarketMap.GetValidSubset()
	if err != nil {
		o.logger.Error("failed to validate market map", zap.Error(err))
		return mmtypes.MarketMap{}, false, fmt.Errorf("failed to validate market map: %w", err)
	}

	// Detect removed markets and surface info about the removals
	var removedMarkets []string
	for t := range resp.MarketMap.Markets {
		if _, in := validSubset.Markets[t]; !in {
			removedMarkets = append(removedMarkets, t)
		}
	}
	if len(validSubset.Markets) == 0 || len(validSubset.Markets) != len(resp.MarketMap.Markets) {
		o.logger.Warn("invalid market map update has caused some markets to be removed")
		o.logger.Info("markets removed from invalid market map", zap.String("markets", strings.Join(removedMarkets, " ")))
	}

	// Update the oracle with the latest market map iff the market map has changed.
	if o.marketMap.Equal(validSubset) {
		o.logger.Debug("market map has not changed")
		return mmtypes.MarketMap{}, false, nil
	}

	return validSubset, true, nil
}

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

			newMarketMap, isUpdated, err := o.IsMarketMapValidUpdated(result.Value)
			if err != nil {
				o.logger.Error("failed to check new market map", zap.Error(err))
				continue
			}

			if !isUpdated {
				continue
			}

			o.logger.Info("updating oracle with new market map")
			if err := o.UpdateMarketMap(newMarketMap); err != nil {
				o.logger.Error("failed to update oracle with new market map", zap.Error(err))
				continue
			}

			o.lastUpdated = result.Value.GetLastUpdated()

			// Write the market map to the configured path.
			if err := o.WriteMarketMap(); err != nil {
				o.logger.Error("failed to write market map", zap.Error(err))
			}

			o.logger.Info("updated oracle with new market map")
			o.logger.Debug("updated oracle with new market map", zap.Any("market_map", newMarketMap))
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
