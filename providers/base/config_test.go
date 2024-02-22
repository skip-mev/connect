package base_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/base"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var (
	btcusd = slinkytypes.NewCurrencyPair("BITCOIN", "USD")
	ethusd = slinkytypes.NewCurrencyPair("ETHEREUM", "USD")
	solusd = slinkytypes.NewCurrencyPair("SOLANA", "USD")
)

func TestConfigUpdater(t *testing.T) {
	t.Run("restart on IDs update with an API provider", func(t *testing.T) {
		pairs := []slinkytypes.CurrencyPair{btcusd}
		updater := base.NewConfigUpdater[slinkytypes.CurrencyPair, *big.Int]()
		apiHandler := testutils.CreateAPIQueryHandlerWithGetResponses[slinkytypes.CurrencyPair, *big.Int](
			t,
			logger,
			nil,
		)

		provider, err := base.NewProvider[slinkytypes.CurrencyPair, *big.Int](
			base.WithName[slinkytypes.CurrencyPair, *big.Int](apiCfg.Name),
			base.WithAPIQueryHandler[slinkytypes.CurrencyPair, *big.Int](apiHandler),
			base.WithAPIConfig[slinkytypes.CurrencyPair, *big.Int](apiCfg),
			base.WithLogger[slinkytypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[slinkytypes.CurrencyPair, *big.Int](pairs),
			base.WithConfigUpdater[slinkytypes.CurrencyPair, *big.Int](updater),
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
		updated := []slinkytypes.CurrencyPair{ethusd, solusd, btcusd}
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
		pairs := []slinkytypes.CurrencyPair{btcusd}
		updater := base.NewConfigUpdater[slinkytypes.CurrencyPair, *big.Int]()
		wsHandler := testutils.CreateWebSocketQueryHandlerWithGetResponses[slinkytypes.CurrencyPair, *big.Int](
			t,
			time.Second,
			logger,
			nil,
		)

		provider, err := base.NewProvider[slinkytypes.CurrencyPair, *big.Int](
			base.WithName[slinkytypes.CurrencyPair, *big.Int](wsCfg.Name),
			base.WithWebSocketQueryHandler[slinkytypes.CurrencyPair, *big.Int](wsHandler),
			base.WithWebSocketConfig[slinkytypes.CurrencyPair, *big.Int](wsCfg),
			base.WithLogger[slinkytypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[slinkytypes.CurrencyPair, *big.Int](pairs),
			base.WithConfigUpdater[slinkytypes.CurrencyPair, *big.Int](updater),
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
		updated := []slinkytypes.CurrencyPair{ethusd, solusd, btcusd}
		updater.UpdateIDs(updated)

		// Wait for the provider to restart.
		time.Sleep(2 * time.Second)

		// The IDs should be updated.
		ids = provider.GetIDs()
		require.Equal(t, updated, ids)

		// Check that the provider exited without error.
		require.Equal(t, context.DeadlineExceeded, <-errCh)
	})

	t.Run("restart on API handler update", func(t *testing.T) {
		pairs := []slinkytypes.CurrencyPair{btcusd}
		updater := base.NewConfigUpdater[slinkytypes.CurrencyPair, *big.Int]()
		apiHandler := testutils.CreateAPIQueryHandlerWithGetResponses[slinkytypes.CurrencyPair, *big.Int](
			t,
			logger,
			nil,
		)

		provider, err := base.NewProvider[slinkytypes.CurrencyPair, *big.Int](
			base.WithName[slinkytypes.CurrencyPair, *big.Int](apiCfg.Name),
			base.WithAPIQueryHandler[slinkytypes.CurrencyPair, *big.Int](apiHandler),
			base.WithAPIConfig[slinkytypes.CurrencyPair, *big.Int](apiCfg),
			base.WithLogger[slinkytypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[slinkytypes.CurrencyPair, *big.Int](pairs),
			base.WithConfigUpdater[slinkytypes.CurrencyPair, *big.Int](updater),
		)
		require.NoError(t, err)

		// Start the provider and run it for a few seconds.
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		errCh := make(chan error)
		go func() {
			errCh <- provider.Start(ctx)
		}()

		// The initial API handler should be the same as the provider's API handler.
		handler := provider.GetAPIHandler()
		require.Equal(t, apiHandler, handler)

		// Wait for a few seconds and update the API handler with a handler that returns some data.
		time.Sleep(2 * time.Second)

		resolved := map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
			pairs[0]: {
				Value:     big.NewInt(100),
				Timestamp: respTime,
			},
		}
		responses := []providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
			providertypes.NewGetResponse(resolved, nil),
		}
		updatedAPIHandler := testutils.CreateAPIQueryHandlerWithGetResponses[slinkytypes.CurrencyPair, *big.Int](t, logger, responses)
		updater.UpdateAPIHandler(updatedAPIHandler)

		// Wait for the provider to restart.
		time.Sleep(2 * time.Second)

		// The API handler should be updated.
		handler = provider.GetAPIHandler()
		require.Equal(t, updatedAPIHandler, handler)

		// Check that the provider exited without error.
		require.Equal(t, context.DeadlineExceeded, <-errCh)
	})

	t.Run("restart on WebSocket handler update", func(t *testing.T) {
		pairs := []slinkytypes.CurrencyPair{btcusd}
		updater := base.NewConfigUpdater[slinkytypes.CurrencyPair, *big.Int]()
		wsHandler := testutils.CreateWebSocketQueryHandlerWithGetResponses[slinkytypes.CurrencyPair, *big.Int](
			t,
			time.Second,
			logger,
			nil,
		)

		provider, err := base.NewProvider[slinkytypes.CurrencyPair, *big.Int](
			base.WithName[slinkytypes.CurrencyPair, *big.Int](wsCfg.Name),
			base.WithWebSocketQueryHandler[slinkytypes.CurrencyPair, *big.Int](wsHandler),
			base.WithWebSocketConfig[slinkytypes.CurrencyPair, *big.Int](wsCfg),
			base.WithLogger[slinkytypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[slinkytypes.CurrencyPair, *big.Int](pairs),
			base.WithConfigUpdater[slinkytypes.CurrencyPair, *big.Int](updater),
		)
		require.NoError(t, err)

		// Start the provider and run it for a few seconds.
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		errCh := make(chan error)
		go func() {
			errCh <- provider.Start(ctx)
		}()

		// The initial WebSocket handler should be the same as the provider's WebSocket handler.
		handler := provider.GetWebSocketHandler()
		require.Equal(t, wsHandler, handler)

		// Wait for a few seconds and update the WebSocket handler with a handler that returns some data.
		time.Sleep(2 * time.Second)

		resolved := map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
			pairs[0]: {
				Value:     big.NewInt(100),
				Timestamp: respTime,
			},
		}
		responses := []providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
			providertypes.NewGetResponse(resolved, nil),
		}
		updatedWsHandler := testutils.CreateWebSocketQueryHandlerWithGetResponses[slinkytypes.CurrencyPair, *big.Int](t, time.Second, logger, responses)
		updater.UpdateWebSocketHandler(updatedWsHandler)

		// Wait for the provider to restart.
		time.Sleep(2 * time.Second)

		// The WebSocket handler should be updated.
		handler = provider.GetWebSocketHandler()
		require.Equal(t, updatedWsHandler, handler)

		// Check that the provider exited without error.
		require.Equal(t, context.DeadlineExceeded, <-errCh)
	})
}
