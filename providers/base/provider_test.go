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

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	apierrors "github.com/skip-mev/slinky/providers/base/api/errors"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	apihandlermocks "github.com/skip-mev/slinky/providers/base/api/handlers/mocks"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	metricmocks "github.com/skip-mev/slinky/providers/base/metrics/mocks"
	"github.com/skip-mev/slinky/providers/base/testutils"
	wserrors "github.com/skip-mev/slinky/providers/base/websocket/errors"
	wshandlers "github.com/skip-mev/slinky/providers/base/websocket/handlers"
	wshandlermocks "github.com/skip-mev/slinky/providers/base/websocket/handlers/mocks"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	logger = zap.NewExample()
	apiCfg = config.APIConfig{
		Enabled:    true,
		Timeout:    time.Millisecond * 250,
		Interval:   time.Millisecond * 500,
		MaxQueries: 1,
		URL:        "localhost:8080",
		Name:       "api",
	}
	wsCfg = config.WebSocketConfig{
		Enabled:                       true,
		MaxBufferSize:                 10,
		ReconnectionTimeout:           time.Millisecond * 500,
		WSS:                           "wss:localhost:8080",
		Name:                          "websocket",
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             config.DefaultEnableCompression,
		ReadTimeout:                   config.DefaultReadTimeout,
		WriteTimeout:                  config.DefaultWriteTimeout,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	}

	wsCfgMultiplex = config.WebSocketConfig{
		Enabled:                       true,
		MaxBufferSize:                 10,
		ReconnectionTimeout:           time.Millisecond * 500,
		WSS:                           "wss:localhost:8080",
		Name:                          "websocket",
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             config.DefaultEnableCompression,
		ReadTimeout:                   config.DefaultReadTimeout,
		WriteTimeout:                  config.DefaultWriteTimeout,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: 1,
	}

	pairs = []oracletypes.CurrencyPair{
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
		handler := apihandlermocks.NewQueryHandler[oracletypes.CurrencyPair, *big.Int](t)

		provider, err := base.NewProvider(
			base.WithName[oracletypes.CurrencyPair, *big.Int](apiCfg.Name),
			base.WithAPIQueryHandler[oracletypes.CurrencyPair, *big.Int](handler),
			base.WithAPIConfig[oracletypes.CurrencyPair, *big.Int](apiCfg),
			base.WithLogger[oracletypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[oracletypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err = provider.Start(ctx)
		require.Equal(t, context.Canceled, err)
	})

	t.Run("closes with deadline with api", func(t *testing.T) {
		handler := apihandlermocks.NewQueryHandler[oracletypes.CurrencyPair, *big.Int](t)
		handler.On("Query", mock.Anything, mock.Anything, mock.Anything).Return()

		provider, err := base.NewProvider(
			base.WithName[oracletypes.CurrencyPair, *big.Int](apiCfg.Name),
			base.WithAPIQueryHandler[oracletypes.CurrencyPair, *big.Int](handler),
			base.WithAPIConfig[oracletypes.CurrencyPair, *big.Int](apiCfg),
			base.WithLogger[oracletypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[oracletypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), apiCfg.Interval*2)
		defer cancel()

		err = provider.Start(ctx)
		require.Equal(t, context.DeadlineExceeded, err)
	})

	t.Run("closes on cancel with websocket", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		handler := wshandlermocks.NewWebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int](t)

		provider, err := base.NewProvider(
			base.WithName[oracletypes.CurrencyPair, *big.Int](wsCfg.Name),
			base.WithWebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int](handler),
			base.WithWebSocketConfig[oracletypes.CurrencyPair, *big.Int](wsCfg),
			base.WithLogger[oracletypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[oracletypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		err = provider.Start(ctx)
		require.Equal(t, context.Canceled, err)
	})

	t.Run("closes on cancel with websocket multiplex", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		handler := wshandlermocks.NewWebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int](t)

		provider, err := base.NewProvider(
			base.WithName[oracletypes.CurrencyPair, *big.Int](wsCfgMultiplex.Name),
			base.WithWebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int](handler),
			base.WithWebSocketConfig[oracletypes.CurrencyPair, *big.Int](wsCfgMultiplex),
			base.WithLogger[oracletypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[oracletypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		err = provider.Start(ctx)
		require.Equal(t, context.Canceled, err)
	})

	t.Run("closes with deadline with websocket", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), apiCfg.Interval*2)
		defer cancel()

		handler := wshandlermocks.NewWebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int](t)
		handler.On("Start", mock.Anything, mock.Anything, mock.Anything).Return(func() error {
			<-ctx.Done()
			return ctx.Err()
		}()).Maybe()

		provider, err := base.NewProvider(
			base.WithName[oracletypes.CurrencyPair, *big.Int](wsCfg.Name),
			base.WithWebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int](handler),
			base.WithWebSocketConfig[oracletypes.CurrencyPair, *big.Int](wsCfg),
			base.WithLogger[oracletypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[oracletypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		err = provider.Start(ctx)
		require.Equal(t, context.DeadlineExceeded, err)
	})

	t.Run("closes with deadline with websocket multiplex", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), apiCfg.Interval*2)
		defer cancel()

		handler := wshandlermocks.NewWebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int](t)
		handler.On("Start", mock.Anything, mock.Anything, mock.Anything).Return(func() error {
			<-ctx.Done()
			return ctx.Err()
		}()).Maybe()

		provider, err := base.NewProvider(
			base.WithName[oracletypes.CurrencyPair, *big.Int](wsCfgMultiplex.Name),
			base.WithWebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int](handler),
			base.WithWebSocketConfig[oracletypes.CurrencyPair, *big.Int](wsCfgMultiplex),
			base.WithLogger[oracletypes.CurrencyPair, *big.Int](logger),
			base.WithIDs[oracletypes.CurrencyPair, *big.Int](pairs),
		)
		require.NoError(t, err)

		err = provider.Start(ctx)
		require.Equal(t, context.DeadlineExceeded, err)
	})
}

func TestWebSocketProvider(t *testing.T) {
	testCases := []struct {
		name           string
		handler        func() wshandlers.WebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int]
		pairs          []oracletypes.CurrencyPair
		cfg            config.WebSocketConfig
		expectedPrices map[oracletypes.CurrencyPair]*big.Int
	}{
		{
			name: "no prices to fetch",
			handler: func() wshandlers.WebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				return testutils.CreateWebSocketQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					t,
					time.Second,
					logger,
					nil,
				)
			},
			pairs:          []oracletypes.CurrencyPair{},
			cfg:            wsCfg,
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{},
		},
		{
			name: "can fetch a single price",
			handler: func() wshandlers.WebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: respTime,
					},
				}
				responses := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse(resolved, nil),
				}

				return testutils.CreateWebSocketQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					t,
					time.Second,
					logger,
					responses,
				)
			},
			pairs: []oracletypes.CurrencyPair{
				pairs[0],
			},
			cfg: wsCfg,
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
			},
		},
		{
			name: "can fetch prices and only updates if the timestamp is greater than the current data",
			handler: func() wshandlers.WebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				fn := func(responseCh chan<- providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]) {
					resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
						pairs[0]: {
							Value:     big.NewInt(100),
							Timestamp: respTime,
						},
					}
					resp := providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, nil)

					logger.Debug("sending response", zap.String("response", resp.String()))
					time.Sleep(time.Second)
					responseCh <- resp

					resolved = map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
						pairs[0]: {
							Value:     big.NewInt(200),
							Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					}
					resp = providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, nil)

					logger.Debug("sending response", zap.String("response", resp.String()))
					time.Sleep(time.Second)
					responseCh <- resp
				}

				return testutils.CreateWebSocketQueryHandlerWithResponseFn[oracletypes.CurrencyPair, *big.Int](
					t,
					fn,
				)
			},
			pairs: []oracletypes.CurrencyPair{
				pairs[0],
			},
			cfg: wsCfg,
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
			},
		},
		{
			name: "can fetch multiple prices",
			handler: func() wshandlers.WebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: respTime,
					},
					pairs[1]: {
						Value:     big.NewInt(200),
						Timestamp: respTime,
					},
				}

				responses := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse(resolved, nil),
				}

				return testutils.CreateWebSocketQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					t,
					time.Second,
					logger,
					responses,
				)
			},
			pairs: []oracletypes.CurrencyPair{
				pairs[0],
				pairs[1],
			},
			cfg: wsCfg,
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
				pairs[1]: big.NewInt(200),
			},
		},
		{
			name: "can fetch multiple prices multiplexed",
			handler: func() wshandlers.WebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: respTime,
					},
					pairs[1]: {
						Value:     big.NewInt(200),
						Timestamp: respTime,
					},
				}

				responses := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse(resolved, nil),
				}

				return testutils.CreateWebSocketQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					t,
					time.Second,
					logger,
					responses,
				)
			},
			pairs: []oracletypes.CurrencyPair{
				pairs[0],
				pairs[1],
			},
			cfg: wsCfgMultiplex,
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
				pairs[1]: big.NewInt(200),
			},
		},
		{
			name: "does not update if the response included an error",
			handler: func() wshandlers.WebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				unResolved := map[oracletypes.CurrencyPair]error{
					pairs[0]: wserrors.ErrHandleMessage,
				}

				responses := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](nil, unResolved),
				}

				return testutils.CreateWebSocketQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					t,
					time.Second,
					logger,
					responses,
				)
			},
			pairs: []oracletypes.CurrencyPair{
				pairs[0],
			},
			cfg:            wsCfg,
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{},
		},
		{
			name: "continues restarting if the query handler returns",
			handler: func() wshandlers.WebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				handler := wshandlermocks.NewWebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int](t)

				handler.On("Start", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("no gib price updates")).Maybe()

				return handler
			},
			pairs: []oracletypes.CurrencyPair{
				pairs[0],
			},
			cfg:            wsCfg,
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{},
		},
		{
			name: "continues restarting if the query handler returns multiplexed",
			handler: func() wshandlers.WebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				handler := wshandlermocks.NewWebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int](t)

				handler.On("Start", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("no gib price updates")).Maybe()

				return handler
			},
			pairs: []oracletypes.CurrencyPair{
				pairs[0],
			},
			cfg:            wsCfgMultiplex,
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider, err := base.NewProvider[oracletypes.CurrencyPair, *big.Int](
				base.WithName[oracletypes.CurrencyPair, *big.Int](tc.cfg.Name),
				base.WithWebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int](tc.handler()),
				base.WithWebSocketConfig[oracletypes.CurrencyPair, *big.Int](tc.cfg),
				base.WithLogger[oracletypes.CurrencyPair, *big.Int](logger),
				base.WithIDs[oracletypes.CurrencyPair, *big.Int](tc.pairs),
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
		handler        func() apihandlers.APIQueryHandler[oracletypes.CurrencyPair, *big.Int]
		pairs          []oracletypes.CurrencyPair
		expectedPrices map[oracletypes.CurrencyPair]*big.Int
	}{
		{
			name: "no prices to fetch",
			handler: func() apihandlers.APIQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				return testutils.CreateAPIQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					t,
					logger,
					nil,
				)
			},
			pairs:          []oracletypes.CurrencyPair{},
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{},
		},
		{
			name: "can fetch a single price",
			handler: func() apihandlers.APIQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: respTime,
					},
				}
				responses := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse(resolved, nil),
				}

				return testutils.CreateAPIQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					t,
					logger,
					responses,
				)
			},
			pairs: []oracletypes.CurrencyPair{
				pairs[0],
			},
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
			},
		},
		{
			name: "can fetch prices and only updates if the timestamp is greater than the current data",
			handler: func() apihandlers.APIQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				fn := func(responseCh chan<- providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]) {
					resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
						pairs[0]: {
							Value:     big.NewInt(100),
							Timestamp: respTime,
						},
					}
					resp := providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, nil)

					logger.Debug("sending response", zap.String("response", resp.String()))
					responseCh <- resp

					resolved = map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
						pairs[0]: {
							Value:     big.NewInt(200),
							Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					}
					resp = providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, nil)

					logger.Debug("sending response", zap.String("response", resp.String()))
					responseCh <- resp
				}

				return testutils.CreateAPIQueryHandlerWithResponseFn[oracletypes.CurrencyPair, *big.Int](
					t,
					fn,
				)
			},
			pairs: []oracletypes.CurrencyPair{
				pairs[0],
			},
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
			},
		},
		{
			name: "can fetch multiple prices",
			handler: func() apihandlers.APIQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: respTime,
					},
					pairs[1]: {
						Value:     big.NewInt(200),
						Timestamp: respTime,
					},
				}

				responses := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse(resolved, nil),
				}

				return testutils.CreateAPIQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					t,
					logger,
					responses,
				)
			},
			pairs: []oracletypes.CurrencyPair{
				pairs[0],
				pairs[1],
			},
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
				pairs[1]: big.NewInt(200),
			},
		},
		{
			name: "does not update if the response included an error",
			handler: func() apihandlers.APIQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				unResolved := map[oracletypes.CurrencyPair]error{
					pairs[0]: apierrors.ErrRateLimit,
				}

				responses := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](nil, unResolved),
				}

				return testutils.CreateAPIQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					t,
					logger,
					responses,
				)
			},
			pairs: []oracletypes.CurrencyPair{
				pairs[0],
			},
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{},
		},
		{
			name: "continues updating even with timeouts on the query handler",
			handler: func() apihandlers.APIQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				handler := apihandlermocks.NewQueryHandler[oracletypes.CurrencyPair, *big.Int](t)

				handler.On("Query", mock.Anything, mock.Anything, mock.Anything).Return().Run(func(args mock.Arguments) {
					responseCh := args.Get(2).(chan<- providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int])

					resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
						pairs[0]: {
							Value:     big.NewInt(100),
							Timestamp: respTime,
						},
					}
					resp := providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, nil)

					logger.Debug("sending response", zap.String("response", resp.String()))
					responseCh <- resp
				}).After(apiCfg.Interval * 2)

				return handler
			},
			pairs: []oracletypes.CurrencyPair{
				pairs[0],
			},
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider, err := base.NewProvider[oracletypes.CurrencyPair, *big.Int](
				base.WithName[oracletypes.CurrencyPair, *big.Int](apiCfg.Name),
				base.WithAPIQueryHandler[oracletypes.CurrencyPair, *big.Int](tc.handler()),
				base.WithAPIConfig[oracletypes.CurrencyPair, *big.Int](apiCfg),
				base.WithLogger[oracletypes.CurrencyPair, *big.Int](logger),
				base.WithIDs[oracletypes.CurrencyPair, *big.Int](tc.pairs),
			)
			require.NoError(t, err)

			now := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), apiCfg.Interval*5)
			defer cancel()

			err = provider.Start(ctx)
			require.Equal(t, context.DeadlineExceeded, err)

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
		handler func() apihandlers.APIQueryHandler[oracletypes.CurrencyPair, *big.Int]
		metrics func() providermetrics.ProviderMetrics
		pairs   []oracletypes.CurrencyPair
	}{
		{
			name: "can fetch a single price",
			handler: func() apihandlers.APIQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: respTime,
					},
				}
				responses := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse(resolved, nil),
				}

				return testutils.CreateAPIQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					t,
					logger,
					responses,
				)
			},
			metrics: func() providermetrics.ProviderMetrics {
				m := metricmocks.NewProviderMetrics(t)
				p1 := strings.ToLower(fmt.Sprint(pairs[0]))

				m.On("AddProviderResponseByID", apiCfg.Name, p1, providermetrics.Success, providertypes.API).Maybe()
				m.On("AddProviderResponse", apiCfg.Name, providermetrics.Success, providertypes.API).Maybe()
				m.On("LastUpdated", apiCfg.Name, p1, providertypes.API).Maybe()

				return m
			},
			pairs: []oracletypes.CurrencyPair{
				pairs[0],
			},
		},
		{
			name: "updates correctly with bad responses",
			handler: func() apihandlers.APIQueryHandler[oracletypes.CurrencyPair, *big.Int] {
				unResolved := map[oracletypes.CurrencyPair]error{
					pairs[0]: apierrors.ErrRateLimit,
				}

				responses := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](nil, unResolved),
				}

				return testutils.CreateAPIQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					t,
					logger,
					responses,
				)
			},
			metrics: func() providermetrics.ProviderMetrics {
				m := metricmocks.NewProviderMetrics(t)
				p1 := strings.ToLower(fmt.Sprint(pairs[0]))

				m.On("AddProviderResponseByID", apiCfg.Name, p1, providermetrics.Failure, providertypes.API).Maybe()
				m.On("AddProviderResponse", apiCfg.Name, providermetrics.Failure, providertypes.API).Maybe()
				m.On("LastUpdated", apiCfg.Name, p1, providertypes.API).Maybe()

				return m
			},
			pairs: []oracletypes.CurrencyPair{
				pairs[0],
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider, err := base.NewProvider[oracletypes.CurrencyPair, *big.Int](
				base.WithName[oracletypes.CurrencyPair, *big.Int](apiCfg.Name),
				base.WithAPIQueryHandler[oracletypes.CurrencyPair, *big.Int](tc.handler()),
				base.WithAPIConfig[oracletypes.CurrencyPair, *big.Int](apiCfg),
				base.WithLogger[oracletypes.CurrencyPair, *big.Int](logger),
				base.WithIDs[oracletypes.CurrencyPair, *big.Int](tc.pairs),
				base.WithMetrics[oracletypes.CurrencyPair, *big.Int](tc.metrics()),
			)
			require.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), apiCfg.Interval*5)
			defer cancel()

			err = provider.Start(ctx)
			require.Equal(t, context.DeadlineExceeded, err)
		})
	}
}
