package base_test

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/base"
	apierrors "github.com/skip-mev/connect/v2/providers/base/api/errors"
	apihandlers "github.com/skip-mev/connect/v2/providers/base/api/handlers"
	apihandlermocks "github.com/skip-mev/connect/v2/providers/base/api/handlers/mocks"
	providermetrics "github.com/skip-mev/connect/v2/providers/base/metrics"
	metricmocks "github.com/skip-mev/connect/v2/providers/base/metrics/mocks"
	"github.com/skip-mev/connect/v2/providers/base/testutils"
	wserrors "github.com/skip-mev/connect/v2/providers/base/websocket/errors"
	wshandlers "github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
	wshandlermocks "github.com/skip-mev/connect/v2/providers/base/websocket/handlers/mocks"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

var (
	logger = zap.NewExample()
	apiCfg = config.APIConfig{
		Enabled:          true,
		Timeout:          time.Millisecond * 250,
		Interval:         time.Millisecond * 500,
		ReconnectTimeout: time.Millisecond * 500,
		MaxQueries:       100,
		Endpoints:        []config.Endpoint{{URL: "http://test.com"}},
		Name:             "api",
	}
	wsCfg = config.WebSocketConfig{
		Enabled:             true,
		MaxBufferSize:       10,
		ReconnectionTimeout: time.Millisecond * 500,
		Endpoints: []config.Endpoint{
			{
				URL: "ws://localhost:8080",
			},
		},
		Name:                          "websocket",
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             config.DefaultEnableCompression,
		ReadTimeout:                   config.DefaultReadTimeout,
		WriteTimeout:                  config.DefaultWriteTimeout,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
		MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
	}

	wsCfgMultiplex = config.WebSocketConfig{
		Enabled:             true,
		MaxBufferSize:       10,
		ReconnectionTimeout: time.Millisecond * 500,
		Endpoints: []config.Endpoint{
			{
				URL: "ws://localhost:8080",
			},
		},
		Name:                          "websocket",
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             config.DefaultEnableCompression,
		ReadTimeout:                   config.DefaultReadTimeout,
		WriteTimeout:                  config.DefaultWriteTimeout,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: 1,
		MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
	}

	pairs = []connecttypes.CurrencyPair{
		{
			Base:  "BTC",
			Quote: "USD",
		},
		{
			Base:  "ETH",
			Quote: "USD",
		},
	}

	respTime = time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC)
)

func TestStart(t *testing.T) {
	t.Parallel()

	t.Run("closes on cancel with api", func(t *testing.T) {
		t.Parallel()

		handler := apihandlermocks.NewQueryHandler[connecttypes.CurrencyPair, *big.Int](t)

		handler.On("Query", mock.Anything, mock.Anything, mock.Anything).Return().Maybe().After(200 * time.Millisecond)

		provider, err := base.NewProvider(
			base.WithName[connecttypes.CurrencyPair, *big.Int](apiCfg.Name),
			base.WithAPIQueryHandler[connecttypes.CurrencyPair, *big.Int](handler),
			base.WithAPIConfig[connecttypes.CurrencyPair, *big.Int](apiCfg),
			base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[connecttypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err = provider.Start(ctx)
		require.Equal(t, context.Canceled, err)
	})

	t.Run("closes with deadline with api", func(t *testing.T) {
		t.Parallel()

		handler := apihandlermocks.NewQueryHandler[connecttypes.CurrencyPair, *big.Int](t)
		handler.On("Query", mock.Anything, mock.Anything, mock.Anything).Return().Maybe().After(200 * time.Millisecond)

		provider, err := base.NewProvider(
			base.WithName[connecttypes.CurrencyPair, *big.Int](apiCfg.Name),
			base.WithAPIQueryHandler[connecttypes.CurrencyPair, *big.Int](handler),
			base.WithAPIConfig[connecttypes.CurrencyPair, *big.Int](apiCfg),
			base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[connecttypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), apiCfg.Interval*2)
		defer cancel()

		err = provider.Start(ctx)
		require.Equal(t, context.DeadlineExceeded, err)
	})

	t.Run("closes on cancel with websocket", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		handler := wshandlermocks.NewWebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int](t)
		handler.On("Copy").Return(handler).Maybe()

		provider, err := base.NewProvider(
			base.WithName[connecttypes.CurrencyPair, *big.Int](wsCfg.Name),
			base.WithWebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int](handler),
			base.WithWebSocketConfig[connecttypes.CurrencyPair, *big.Int](wsCfg),
			base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[connecttypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		err = provider.Start(ctx)
		require.Equal(t, context.Canceled, err)
	})

	t.Run("closes on cancel with websocket multiplex", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		handler := wshandlermocks.NewWebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int](t)
		handler.On("Copy").Return(handler).Maybe()

		provider, err := base.NewProvider(
			base.WithName[connecttypes.CurrencyPair, *big.Int](wsCfgMultiplex.Name),
			base.WithWebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int](handler),
			base.WithWebSocketConfig[connecttypes.CurrencyPair, *big.Int](wsCfgMultiplex),
			base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[connecttypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		err = provider.Start(ctx)
		require.Equal(t, context.Canceled, err)
	})

	t.Run("closes with deadline with websocket", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), apiCfg.Interval*2)
		defer cancel()

		handler := wshandlermocks.NewWebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int](t)
		handler.On("Copy").Return(handler).Maybe()
		handler.On("Start", mock.Anything, mock.Anything, mock.Anything).Return(func() error {
			<-ctx.Done()
			return ctx.Err()
		}()).Maybe()

		provider, err := base.NewProvider(
			base.WithName[connecttypes.CurrencyPair, *big.Int](wsCfg.Name),
			base.WithWebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int](handler),
			base.WithWebSocketConfig[connecttypes.CurrencyPair, *big.Int](wsCfg),
			base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[connecttypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		err = provider.Start(ctx)
		require.Equal(t, context.DeadlineExceeded, err)
	})

	t.Run("closes with deadline with websocket multiplex", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), apiCfg.Interval*2)
		defer cancel()

		handler := wshandlermocks.NewWebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int](t)
		handler.On("Copy").Return(handler).Maybe()
		handler.On("Start", mock.Anything, mock.Anything, mock.Anything).Return(func() error {
			<-ctx.Done()
			return ctx.Err()
		}()).Maybe()

		provider, err := base.NewProvider(
			base.WithName[connecttypes.CurrencyPair, *big.Int](wsCfgMultiplex.Name),
			base.WithWebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int](handler),
			base.WithWebSocketConfig[connecttypes.CurrencyPair, *big.Int](wsCfgMultiplex),
			base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[connecttypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		err = provider.Start(ctx)
		require.Equal(t, context.DeadlineExceeded, err)
	})
}

func TestStop(t *testing.T) {
	t.Run("no error when not running", func(t *testing.T) {
		handler := apihandlermocks.NewQueryHandler[connecttypes.CurrencyPair, *big.Int](t)

		provider, err := base.NewProvider(
			base.WithName[connecttypes.CurrencyPair, *big.Int](apiCfg.Name),
			base.WithAPIQueryHandler[connecttypes.CurrencyPair, *big.Int](handler),
			base.WithAPIConfig[connecttypes.CurrencyPair, *big.Int](apiCfg),
			base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[connecttypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)
		provider.Stop()

		require.Eventually(t, func() bool { return !provider.IsRunning() }, time.Second*3, time.Millisecond*100)
	})

	t.Run("no error when running an API provider", func(t *testing.T) {
		handler := testutils.CreateAPIQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
			t,
			logger,
			nil,
			200*time.Millisecond,
		)

		provider, err := base.NewProvider(
			base.WithName[connecttypes.CurrencyPair, *big.Int](apiCfg.Name),
			base.WithAPIQueryHandler[connecttypes.CurrencyPair, *big.Int](handler),
			base.WithAPIConfig[connecttypes.CurrencyPair, *big.Int](apiCfg),
			base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[connecttypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		// Use a timeout greater than the interval to ensure that the provider is running.
		now := time.Now()
		timeout := apiCfg.Interval * 10
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		go func() {
			err = provider.Start(ctx)
			require.Error(t, err)
		}()

		time.Sleep(time.Second * 3)
		provider.Stop()
		require.True(t, time.Since(now) < timeout)

		require.Eventually(t, func() bool { return !provider.IsRunning() }, time.Second*3, time.Millisecond*100)
	})

	t.Run("no error when running a WebSocket provider", func(t *testing.T) {
		handler := testutils.CreateWebSocketQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
			t,
			time.Second,
			logger,
			nil,
		)

		provider, err := base.NewProvider(
			base.WithName[connecttypes.CurrencyPair, *big.Int](wsCfg.Name),
			base.WithWebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int](handler),
			base.WithWebSocketConfig[connecttypes.CurrencyPair, *big.Int](wsCfg),
			base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[connecttypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		// Use a timeout greater than the interval to ensure that the provider is running.
		now := time.Now()
		timeout := wsCfg.ReconnectionTimeout * 10
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		go func() {
			err = provider.Start(ctx)
			require.Error(t, err)
		}()

		time.Sleep(time.Second * 3)
		provider.Stop()
		require.True(t, time.Since(now) < timeout)

		require.Eventually(t, func() bool { return !provider.IsRunning() }, time.Second*3, time.Millisecond*100)
	})
}

func TestWebSocketProvider(t *testing.T) {
	testCases := []struct {
		name           string
		handler        func() wshandlers.WebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int]
		pairs          []connecttypes.CurrencyPair
		cfg            config.WebSocketConfig
		expectedPrices map[connecttypes.CurrencyPair]*big.Int
	}{
		{
			name: "no prices to fetch",
			handler: func() wshandlers.WebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				return testutils.CreateWebSocketQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
					t,
					time.Second,
					logger,
					nil,
				)
			},
			pairs:          []connecttypes.CurrencyPair{},
			cfg:            wsCfg,
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{},
		},
		{
			name: "can fetch a single price",
			handler: func() wshandlers.WebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				resolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: respTime,
					},
				}
				responses := []providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse(resolved, nil),
				}

				return testutils.CreateWebSocketQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
					t,
					time.Second,
					logger,
					responses,
				)
			},
			pairs: []connecttypes.CurrencyPair{
				pairs[0],
			},
			cfg: wsCfg,
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
			},
		},
		{
			name: "can fetch prices and only updates if the timestamp is greater than the current data",
			handler: func() wshandlers.WebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				fn := func(ctx context.Context, responseCh chan<- providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]) {
					resolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
						pairs[0]: {
							Value:     big.NewInt(100),
							Timestamp: respTime,
						},
					}
					resp := providertypes.NewGetResponse[connecttypes.CurrencyPair, *big.Int](resolved, nil)

					select {
					case <-ctx.Done():
						return
					case responseCh <- resp:
						logger.Debug("sending response", zap.String("response", resp.String()))
						time.Sleep(time.Second)
					}

					resolved = map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
						pairs[0]: {
							Value:     big.NewInt(200),
							Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					}
					resp = providertypes.NewGetResponse[connecttypes.CurrencyPair, *big.Int](resolved, nil)

					select {
					case <-ctx.Done():
						return
					case responseCh <- resp:
						logger.Debug("sending response", zap.String("response", resp.String()))
						time.Sleep(time.Second)
					}
				}

				return testutils.CreateWebSocketQueryHandlerWithResponseFn[connecttypes.CurrencyPair, *big.Int](
					t,
					fn,
				)
			},
			pairs: []connecttypes.CurrencyPair{
				pairs[0],
			},
			cfg: wsCfg,
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
			},
		},
		{
			name: "can fetch multiple prices",
			handler: func() wshandlers.WebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				resolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: respTime,
					},
					pairs[1]: {
						Value:     big.NewInt(200),
						Timestamp: respTime,
					},
				}

				responses := []providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse(resolved, nil),
				}

				return testutils.CreateWebSocketQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
					t,
					time.Second,
					logger,
					responses,
				)
			},
			pairs: []connecttypes.CurrencyPair{
				pairs[0],
				pairs[1],
			},
			cfg: wsCfg,
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
				pairs[1]: big.NewInt(200),
			},
		},
		{
			name: "can fetch multiple prices multiplexed",
			handler: func() wshandlers.WebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				resolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: respTime,
					},
					pairs[1]: {
						Value:     big.NewInt(200),
						Timestamp: respTime,
					},
				}

				responses := []providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse(resolved, nil),
				}

				return testutils.CreateWebSocketQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
					t,
					time.Second,
					logger,
					responses,
				)
			},
			pairs: []connecttypes.CurrencyPair{
				pairs[0],
				pairs[1],
			},
			cfg: wsCfgMultiplex,
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
				pairs[1]: big.NewInt(200),
			},
		},
		{
			name: "does not update if the response included an error",
			handler: func() wshandlers.WebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				unResolved := map[connecttypes.CurrencyPair]providertypes.UnresolvedResult{
					pairs[0]: {
						ErrorWithCode: providertypes.NewErrorWithCode(wserrors.ErrHandleMessage, providertypes.ErrorWebSocketGeneral),
					},
				}

				responses := []providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse[connecttypes.CurrencyPair, *big.Int](nil, unResolved),
				}

				return testutils.CreateWebSocketQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
					t,
					time.Second,
					logger,
					responses,
				)
			},
			pairs: []connecttypes.CurrencyPair{
				pairs[0],
			},
			cfg:            wsCfg,
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{},
		},
		{
			name: "continues restarting if the query handler returns",
			handler: func() wshandlers.WebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				handler := wshandlermocks.NewWebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int](t)

				handler.On("Copy").Return(handler).Maybe()
				handler.On("Start", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("no gib price updates")).Maybe()

				return handler
			},
			pairs: []connecttypes.CurrencyPair{
				pairs[0],
			},
			cfg:            wsCfg,
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{},
		},
		{
			name: "continues restarting if the query handler returns multiplexed",
			handler: func() wshandlers.WebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				handler := wshandlermocks.NewWebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int](t)

				handler.On("Copy").Return(handler).Maybe()

				handler.On("Start", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("no gib price updates")).Maybe()

				return handler
			},
			pairs: []connecttypes.CurrencyPair{
				pairs[0],
			},
			cfg:            wsCfgMultiplex,
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{},
		},
		{
			name: "updates the timestamp associated with a result if the the data is unchanged and still valid",
			handler: func() wshandlers.WebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				// First response is valid and sets the data.
				resolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}

				unchangedResolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					pairs[0]: {
						Value:        big.NewInt(100),
						Timestamp:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						ResponseCode: providertypes.ResponseCodeUnchanged,
					},
				}

				responses := []providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse(resolved, nil),
					providertypes.NewGetResponse(unchangedResolved, nil),
				}

				return testutils.CreateWebSocketQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
					t,
					time.Second,
					logger,
					responses,
				)
			},
			pairs: []connecttypes.CurrencyPair{
				pairs[0],
			},
			cfg: wsCfg,
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
			},
		},
		{
			name: "does not update the base provider if the result is unchanged but the cache has no entry for the id",
			handler: func() wshandlers.WebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				// First response is valid and sets the data.
				resolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					pairs[0]: {
						Value:        big.NewInt(100),
						Timestamp:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						ResponseCode: providertypes.ResponseCodeUnchanged,
					},
				}

				responses := []providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse(resolved, nil),
				}

				return testutils.CreateWebSocketQueryHandlerWithGetResponses(
					t,
					time.Second,
					logger,
					responses,
				)
			},
			pairs: []connecttypes.CurrencyPair{
				pairs[0],
			},
			cfg:            wsCfg,
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider, err := base.NewProvider[connecttypes.CurrencyPair, *big.Int](
				base.WithName[connecttypes.CurrencyPair, *big.Int](tc.cfg.Name),
				base.WithWebSocketQueryHandler[connecttypes.CurrencyPair, *big.Int](tc.handler()),
				base.WithWebSocketConfig[connecttypes.CurrencyPair, *big.Int](tc.cfg),
				base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
				base.WithIDs[connecttypes.CurrencyPair, *big.Int](tc.pairs),
			)
			require.NoError(t, err)

			now := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
			defer cancel()

			provider.Start(ctx)

			data := provider.GetData()
			for cp, price := range tc.expectedPrices {
				require.Contains(t, data, cp)
				result := data[cp]
				require.Equal(t, price, result.Value)
				require.True(t, result.Timestamp.After(now))
			}
		})
	}
}

func TestAPIProviderLoop(t *testing.T) {
	testCases := []struct {
		name           string
		handler        func() apihandlers.APIQueryHandler[connecttypes.CurrencyPair, *big.Int]
		pairs          []connecttypes.CurrencyPair
		expectedPrices map[connecttypes.CurrencyPair]*big.Int
	}{
		{
			name: "no prices to fetch",
			handler: func() apihandlers.APIQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				return testutils.CreateAPIQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
					t,
					logger,
					nil,
					200*time.Millisecond,
				)
			},
			pairs:          []connecttypes.CurrencyPair{},
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{},
		},
		{
			name: "can fetch a single price",
			handler: func() apihandlers.APIQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				resolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: respTime,
					},
				}
				responses := []providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse(resolved, nil),
				}

				return testutils.CreateAPIQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
					t,
					logger,
					responses,
					200*time.Millisecond,
				)
			},
			pairs: []connecttypes.CurrencyPair{
				pairs[0],
			},
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
			},
		},
		{
			name: "can fetch prices and only updates if the timestamp is greater than the current data",
			handler: func() apihandlers.APIQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				resolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: respTime,
					},
				}
				resp := providertypes.NewGetResponse[connecttypes.CurrencyPair, *big.Int](resolved, nil)

				resolved2 := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(200),
						Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				resp2 := providertypes.NewGetResponse[connecttypes.CurrencyPair, *big.Int](resolved2, nil)

				return testutils.CreateAPIQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
					t,
					logger,
					[]providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{resp, resp2},
					200*time.Millisecond,
				)
			},
			pairs: []connecttypes.CurrencyPair{
				pairs[0],
			},
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
			},
		},
		{
			name: "can fetch multiple prices",
			handler: func() apihandlers.APIQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				resolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: respTime,
					},
					pairs[1]: {
						Value:     big.NewInt(200),
						Timestamp: respTime,
					},
				}

				responses := []providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse(resolved, nil),
				}

				return testutils.CreateAPIQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
					t,
					logger,
					responses,
					200*time.Millisecond,
				)
			},
			pairs: []connecttypes.CurrencyPair{
				pairs[0],
				pairs[1],
			},
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
				pairs[1]: big.NewInt(200),
			},
		},
		{
			name: "does not update if the response included an error",
			handler: func() apihandlers.APIQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				unResolved := map[connecttypes.CurrencyPair]providertypes.UnresolvedResult{
					pairs[0]: {
						ErrorWithCode: providertypes.NewErrorWithCode(apierrors.ErrRateLimit, providertypes.ErrorAPIGeneral),
					},
				}

				responses := []providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse[connecttypes.CurrencyPair, *big.Int](nil, unResolved),
				}

				return testutils.CreateAPIQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
					t,
					logger,
					responses,
					200*time.Millisecond,
				)
			},
			pairs: []connecttypes.CurrencyPair{
				pairs[0],
			},
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider, err := base.NewProvider[connecttypes.CurrencyPair, *big.Int](
				base.WithName[connecttypes.CurrencyPair, *big.Int](apiCfg.Name),
				base.WithAPIQueryHandler[connecttypes.CurrencyPair, *big.Int](tc.handler()),
				base.WithAPIConfig[connecttypes.CurrencyPair, *big.Int](apiCfg),
				base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
				base.WithIDs[connecttypes.CurrencyPair, *big.Int](tc.pairs),
			)
			require.NoError(t, err)

			now := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), apiCfg.Interval*5)
			defer cancel()

			err = provider.Start(ctx)
			if len(tc.pairs) > 0 {
				require.Equal(t, context.DeadlineExceeded, err)
			}

			data := provider.GetData()
			for cp, price := range tc.expectedPrices {
				require.Contains(t, data, cp)
				result := data[cp]
				require.Equal(t, price, result.Value)
				require.True(t, result.Timestamp.After(now))
			}
		})
	}
}

func TestMetrics(t *testing.T) {
	testCases := []struct {
		name    string
		handler func() apihandlers.APIQueryHandler[connecttypes.CurrencyPair, *big.Int]
		metrics func() providermetrics.ProviderMetrics
		pairs   []connecttypes.CurrencyPair
	}{
		{
			name: "can fetch a single price",
			handler: func() apihandlers.APIQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				resolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: respTime,
					},
				}
				responses := []providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse(resolved, nil),
				}

				return testutils.CreateAPIQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
					t,
					logger,
					responses,
					200*time.Millisecond,
				)
			},
			metrics: func() providermetrics.ProviderMetrics {
				m := metricmocks.NewProviderMetrics(t)
				p1 := strings.ToLower(fmt.Sprint(pairs[0]))

				m.On("AddProviderResponseByID", apiCfg.Name, p1, providermetrics.Success, providertypes.OK, providertypes.API).Maybe()
				m.On("AddProviderResponse", apiCfg.Name, providermetrics.Success, providertypes.OK, providertypes.API).Maybe()
				m.On("LastUpdated", apiCfg.Name, p1, providertypes.API).Maybe()

				return m
			},
			pairs: []connecttypes.CurrencyPair{
				pairs[0],
			},
		},
		{
			name: "updates correctly with bad responses",
			handler: func() apihandlers.APIQueryHandler[connecttypes.CurrencyPair, *big.Int] {
				unResolved := map[connecttypes.CurrencyPair]providertypes.UnresolvedResult{
					pairs[0]: {
						ErrorWithCode: providertypes.NewErrorWithCode(apierrors.ErrRateLimit, providertypes.ErrorAPIGeneral),
					},
				}

				responses := []providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse[connecttypes.CurrencyPair, *big.Int](nil, unResolved),
				}

				return testutils.CreateAPIQueryHandlerWithGetResponses[connecttypes.CurrencyPair, *big.Int](
					t,
					logger,
					responses,
					200*time.Millisecond,
				)
			},
			metrics: func() providermetrics.ProviderMetrics {
				m := metricmocks.NewProviderMetrics(t)
				p1 := strings.ToLower(fmt.Sprint(pairs[0]))

				code := providertypes.ErrorAPIGeneral
				m.On("AddProviderResponseByID", apiCfg.Name, p1, providermetrics.Failure, code, providertypes.API).Maybe()
				m.On("AddProviderResponse", apiCfg.Name, providermetrics.Failure, code, providertypes.API).Maybe()
				m.On("LastUpdated", apiCfg.Name, p1, providertypes.API).Maybe()

				return m
			},
			pairs: []connecttypes.CurrencyPair{
				pairs[0],
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider, err := base.NewProvider[connecttypes.CurrencyPair, *big.Int](
				base.WithName[connecttypes.CurrencyPair, *big.Int](apiCfg.Name),
				base.WithAPIQueryHandler[connecttypes.CurrencyPair, *big.Int](tc.handler()),
				base.WithAPIConfig[connecttypes.CurrencyPair, *big.Int](apiCfg),
				base.WithLogger[connecttypes.CurrencyPair, *big.Int](logger),
				base.WithIDs[connecttypes.CurrencyPair, *big.Int](tc.pairs),
				base.WithMetrics[connecttypes.CurrencyPair, *big.Int](tc.metrics()),
			)
			require.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), apiCfg.Interval*5)
			defer cancel()

			err = provider.Start(ctx)
			require.Equal(t, context.DeadlineExceeded, err)
		})
	}
}
