package base_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/base"
	"github.com/skip-mev/connect/v2/providers/base/testutils"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

var (
	btcusd = connecttypes.NewCurrencyPair("BITCOIN", "USD")
	ethusd = connecttypes.NewCurrencyPair("ETHEREUM", "USD")
	solusd = connecttypes.NewCurrencyPair("SOLANA", "USD")
)

func TestConfigUpdater(t *testing.T) {
	t.Run("restart on IDs update with an API provider", func(t *testing.T) {
		pairs := []connecttypes.CurrencyPair{btcusd}
		apiHandler := testutils.CreateAPIQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
			t,
			logger,
			nil,
			200*time.Millisecond,
		)

		provider, err := base.NewProvider[connecttypes.CurrencyPair, *big.Int](
			base.WithName[connecttypes.CurrencyPair, *big.Int](apiCfg.Name),
			base.WithAPIQueryHandler[connecttypes.CurrencyPair, *big.Int](apiHandler),
			base.WithAPIConfig[connecttypes.CurrencyPair, *big.Int](apiCfg),
			base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[connecttypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		// Start the provider and run it for a few seconds.
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		go func() {
			provider.Start(ctx)
		}()

		// The initial IDs should be the same as the provider's IDs.
		ids := provider.GetIDs()
		require.Equal(t, pairs, ids)

		// Wait for a few seconds and update the IDs.
		time.Sleep(2 * time.Second)
		updated := []connecttypes.CurrencyPair{ethusd, solusd, btcusd}
		logger.Debug("test case updating ids")
		provider.Update(base.WithNewIDs[connecttypes.CurrencyPair, *big.Int](updated))

		// Wait for the provider to restart.
		time.Sleep(2 * time.Second)

		// The IDs should be updated.
		ids = provider.GetIDs()
		require.Equal(t, updated, ids)

		// Check that the provider exited without error.
		provider.Stop()
		require.Eventually(t, func() bool { return !provider.IsRunning() }, 2*time.Second, 100*time.Millisecond)
	})

	t.Run("restart on IDs update with a websocket provider", func(t *testing.T) {
		pairs := []connecttypes.CurrencyPair{btcusd}
		wsHandler := testutils.CreateWebSocketQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
			t,
			time.Second,
			logger,
			nil,
		)

		provider, err := base.NewProvider[connecttypes.CurrencyPair, *big.Int](
			base.WithName[connecttypes.CurrencyPair, *big.Int](wsCfg.Name),
			base.WithWebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int](wsHandler),
			base.WithWebSocketConfig[connecttypes.CurrencyPair, *big.Int](wsCfg),
			base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[connecttypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		// Start the provider and run it for a few seconds.
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		go func() {
			provider.Start(ctx)
		}()

		// The initial IDs should be the same as the provider's IDs.
		ids := provider.GetIDs()
		require.Equal(t, pairs, ids)

		// Wait for a few seconds and update the IDs.
		time.Sleep(2 * time.Second)
		updated := []connecttypes.CurrencyPair{ethusd, solusd, btcusd}
		logger.Debug("test case updating ids")
		provider.Update(base.WithNewIDs[connecttypes.CurrencyPair, *big.Int](updated))

		// Wait for the provider to restart.
		time.Sleep(2 * time.Second)

		// The IDs should be updated.
		ids = provider.GetIDs()
		require.Equal(t, updated, ids)

		// Check that the provider exited without error.
		provider.Stop()
		require.Eventually(t, func() bool {
			return !provider.IsRunning()
		}, 2*time.Second, 100*time.Millisecond)
	})

	t.Run("restart on API handler update", func(t *testing.T) {
		pairs := []connecttypes.CurrencyPair{btcusd}
		apiHandler := testutils.CreateAPIQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
			t,
			logger,
			nil,
			200*time.Millisecond,
		)

		provider, err := base.NewProvider[connecttypes.CurrencyPair, *big.Int](
			base.WithName[connecttypes.CurrencyPair, *big.Int](apiCfg.Name),
			base.WithAPIQueryHandler[connecttypes.CurrencyPair, *big.Int](apiHandler),
			base.WithAPIConfig[connecttypes.CurrencyPair, *big.Int](apiCfg),
			base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[connecttypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		// Start the provider and run it for a few seconds.
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		go func() {
			provider.Start(ctx)
		}()

		// The initial API handler should be the same as the provider's API handler.
		handler := provider.GetAPIHandler()
		require.Equal(t, apiHandler, handler)

		// Wait for a few seconds and update the API handler with a handler that returns some data.
		time.Sleep(4 * time.Second)

		resolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
			pairs[0]: {
				Value:     big.NewInt(100),
				Timestamp: respTime,
			},
		}
		responses := []providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
			providertypes.NewGetResponse(resolved, nil),
		}
		updatedAPIHandler := testutils.CreateAPIQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
			t,
			logger,
			responses,
			200*time.Millisecond,
		)
		logger.Debug("test case updating api handler")
		provider.Update(base.WithNewAPIHandler[connecttypes.CurrencyPair, *big.Int](updatedAPIHandler))

		// Wait for the provider to restart.
		time.Sleep(2 * time.Second)

		// The API handler should be updated.
		handler = provider.GetAPIHandler()
		require.Equal(t, updatedAPIHandler, handler)

		// Check that the provider exited without error.
		provider.Stop()
		require.Eventually(t, func() bool { return !provider.IsRunning() }, 2*time.Second, 100*time.Millisecond)
	})

	t.Run("restart on WebSocket handler update", func(t *testing.T) {
		pairs := []connecttypes.CurrencyPair{btcusd}
		wsHandler := testutils.CreateWebSocketQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
			t,
			time.Second,
			logger,
			nil,
		)

		provider, err := base.NewProvider[connecttypes.CurrencyPair, *big.Int](
			base.WithName[connecttypes.CurrencyPair, *big.Int](wsCfg.Name),
			base.WithWebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int](wsHandler),
			base.WithWebSocketConfig[connecttypes.CurrencyPair, *big.Int](wsCfg),
			base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[connecttypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		// Start the provider and run it for a few seconds.
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		go func() {
			provider.Start(ctx)
		}()

		// The initial WebSocket handler should be the same as the provider's WebSocket handler.
		handler := provider.GetWebSocketHandler()
		require.Equal(t, wsHandler, handler)

		// Wait for a few seconds and update the WebSocket handler with a handler that returns some data.
		time.Sleep(4 * time.Second)

		resolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
			pairs[0]: {
				Value:     big.NewInt(100),
				Timestamp: respTime,
			},
		}
		responses := []providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
			providertypes.NewGetResponse(resolved, nil),
		}
		updatedWsHandler := testutils.CreateWebSocketQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](t, time.Second, logger, responses)
		logger.Debug("test case updating websocket handler")
		provider.Update(base.WithNewWebSocketHandler[connecttypes.CurrencyPair, *big.Int](updatedWsHandler))

		// Wait for the provider to restart.
		time.Sleep(2 * time.Second)

		// The WebSocket handler should be updated.
		handler = provider.GetWebSocketHandler()
		require.Equal(t, updatedWsHandler, handler)

		// Check that the provider exited without error.
		provider.Stop()
		require.Eventually(t, func() bool { return !provider.IsRunning() }, 2*time.Second, 100*time.Millisecond)
	})
}
