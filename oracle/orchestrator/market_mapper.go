package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"go.uber.org/zap"
)

func (o *ProviderOrchestrator) listenForMarketMapUpdates(ctx context.Context) func() error {
	return func() error {
		state := o.GetMarketMapperState()
		ticker := time.NewTicker(state.Interval)
		mapper := state.Mapper
		ids := mapper.GetIDs()

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				resp := mapper.GetData()
				result, ok := resp[ids[0]]
				if !ok {
					o.logger.Error("market map update does not contain the expected id", zap.Any("ids", ids))
					continue
				}

				if mm := result.Value.MarketMap; !o.marketMap.Equal(mm) {
					o.UpdateWithMarketMap(mm)
					o.writeMarketMap(mm)
				}
			}
		}
	}
}

func (o *ProviderOrchestrator) writeMarketMap(marketMap mmtypes.MarketMap) error {
	// Open the local market config file. This will overwrite any changes made to the
	// local market config file.
	f, err := os.Create("local_market_config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating local market config file: %v\n", err)
		return err
	}
	defer f.Close()

	// Encode the local market config file.
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(marketMap); err != nil {
		return err
	}

	return nil
}
