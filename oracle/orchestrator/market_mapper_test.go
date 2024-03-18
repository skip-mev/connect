package orchestrator_test

import (
	"context"
	"testing"
	"time"

	"github.com/skip-mev/slinky/oracle/orchestrator"
	"github.com/stretchr/testify/require"
)

func TestListenForMarketMapUpdates(t *testing.T) {
	t.Run("mapper has no chain IDs to fetch should not update the orchestrator", func(t *testing.T) {
		_, provider := createTestMarketMapProvider(t, nil, time.Second, nil)

		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithMarketMapper(provider),
		)
		require.NoError(t, err)

		marketMap := o.GetMarketMap()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			o.Start(ctx)
		}()

		// Wait for the orchestrator to start.
		time.Sleep(500 * time.Millisecond)

		// The orchestrator should not have been updated.
		require.Equal(t, marketMap, o.GetMarketMap())

		// Stop the orchestrator.
		cancel()
		o.Stop()
	})
}
