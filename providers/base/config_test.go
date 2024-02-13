package base_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/skip-mev/slinky/providers/base"
	"github.com/skip-mev/slinky/providers/base/testutils"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/require"
)

var (
	btcusd = oracletypes.NewCurrencyPair("BITCOIN", "USD")
	ethusd = oracletypes.NewCurrencyPair("ETHEREUM", "USD")
	solusd = oracletypes.NewCurrencyPair("SOLANA", "USD")
)

func TestConfigUpdater(t *testing.T) {
	t.Run("restart on IDs update with an API provider", func(t *testing.T) {
		pairs := []oracletypes.CurrencyPair{btcusd}
		updater := base.NewConfigUpdater[oracletypes.CurrencyPair]()
		apiHandler := testutils.CreateAPIQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
			t,
			logger,
			nil,
		)

		provider, err := base.NewProvider[oracletypes.CurrencyPair, *big.Int](
			base.WithName[oracletypes.CurrencyPair, *big.Int](apiCfg.Name),
			base.WithAPIQueryHandler[oracletypes.CurrencyPair, *big.Int](apiHandler),
			base.WithAPIConfig[oracletypes.CurrencyPair, *big.Int](apiCfg),
			base.WithLogger[oracletypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[oracletypes.CurrencyPair, *big.Int](pairs),
			base.WithConfigUpdater[oracletypes.CurrencyPair, *big.Int](updater),
		)
		require.NoError(t, err)

		// Start the provider and run it for a few seconds.
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		errCh := make(chan error)
		go func() {
			errCh <- provider.Start(ctx)
		}()

		// The initial IDs should be the same as the provider's IDs.
		ids := provider.GetIDs()
		require.Equal(t, pairs, ids)

		// Wait for a few seconds and update the IDs.
		time.Sleep(2 * time.Second)
		updated := []oracletypes.CurrencyPair{ethusd, solusd, btcusd}
		updater.UpdateIDs(updated)

		// Wait for the provider to restart.
		time.Sleep(2 * time.Second)

		// The IDs should be updated.
		ids = provider.GetIDs()
		require.Equal(t, updated, ids)

		// Check that the provider exited without error.
		require.Equal(t, context.DeadlineExceeded, <-errCh)
	})

	t.Run("restart on IDs update with a websocket provider", func(t *testing.T) {
		pairs := []oracletypes.CurrencyPair{btcusd}
		updater := base.NewConfigUpdater[oracletypes.CurrencyPair]()
		wsHandler := testutils.CreateWebSocketQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
			t,
			time.Second,
			logger,
			nil,
		)

		provider, err := base.NewProvider[oracletypes.CurrencyPair, *big.Int](
			base.WithName[oracletypes.CurrencyPair, *big.Int](wsCfg.Name),
			base.WithWebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int](wsHandler),
			base.WithWebSocketConfig[oracletypes.CurrencyPair, *big.Int](wsCfg),
			base.WithLogger[oracletypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[oracletypes.CurrencyPair, *big.Int](pairs),
			base.WithConfigUpdater[oracletypes.CurrencyPair, *big.Int](updater),
		)
		require.NoError(t, err)

		// Start the provider and run it for a few seconds.
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		errCh := make(chan error)
		go func() {
			errCh <- provider.Start(ctx)
		}()

		// The initial IDs should be the same as the provider's IDs.
		ids := provider.GetIDs()
		require.Equal(t, pairs, ids)

		// Wait for a few seconds and update the IDs.
		time.Sleep(2 * time.Second)
		updated := []oracletypes.CurrencyPair{ethusd, solusd, btcusd}
		updater.UpdateIDs(updated)

		// Wait for the provider to restart.
		time.Sleep(2 * time.Second)

		// The IDs should be updated.
		ids = provider.GetIDs()
		require.Equal(t, updated, ids)

		// Check that the provider exited without error.
		require.Equal(t, context.DeadlineExceeded, <-errCh)
	})
}
