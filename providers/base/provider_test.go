package base_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	"github.com/skip-mev/slinky/providers/base/errors"
	"github.com/skip-mev/slinky/providers/base/handlers"
	handlermocks "github.com/skip-mev/slinky/providers/base/handlers/mocks"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewDevelopment()
	cfg       = config.ProviderConfig{
		Name:       "test",
		Path:       "test",
		Timeout:    time.Millisecond * 250,
		Interval:   time.Millisecond * 500,
		MaxQueries: 1,
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

	t.Run("closes on cancel", func(t *testing.T) {
		handler := handlermocks.NewQueryHandler[oracletypes.CurrencyPair, *big.Int](t)

		provider, err := base.NewProvider(logger, cfg, handler, pairs)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err = provider.Start(ctx)
		require.Equal(t, context.Canceled, err)
	})

	t.Run("closes with deadline", func(t *testing.T) {
		handler := handlermocks.NewQueryHandler[oracletypes.CurrencyPair, *big.Int](t)
		handler.On("Query", mock.Anything, mock.Anything, mock.Anything).Return()

		provider, err := base.NewProvider(logger, cfg, handler, pairs)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Interval*2)
		defer cancel()

		err = provider.Start(ctx)
		require.Equal(t, context.DeadlineExceeded, err)
	})
}

func TestProviderLoop(t *testing.T) {
	testCases := []struct {
		name           string
		handler        func() handlers.QueryHandler[oracletypes.CurrencyPair, *big.Int]
		pairs          []oracletypes.CurrencyPair
		expectedPrices map[oracletypes.CurrencyPair]*big.Int
	}{
		{
			name: "no prices to fetch",
			handler: func() handlers.QueryHandler[oracletypes.CurrencyPair, *big.Int] {
				return testutils.CreateQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
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
			handler: func() handlers.QueryHandler[oracletypes.CurrencyPair, *big.Int] {
				resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					pairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: respTime,
					},
				}
				responses := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse(resolved, nil),
				}

				return testutils.CreateQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
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
			handler: func() handlers.QueryHandler[oracletypes.CurrencyPair, *big.Int] {
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

				return testutils.CreateQueryHandlerWithResponseFn[oracletypes.CurrencyPair, *big.Int](
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
			handler: func() handlers.QueryHandler[oracletypes.CurrencyPair, *big.Int] {
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

				return testutils.CreateQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
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
			handler: func() handlers.QueryHandler[oracletypes.CurrencyPair, *big.Int] {
				unResolved := map[oracletypes.CurrencyPair]error{
					pairs[0]: errors.ErrRateLimit,
				}

				responses := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
					providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](nil, unResolved),
				}

				return testutils.CreateQueryHandlerWithGetResponses[oracletypes.CurrencyPair, *big.Int](
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
			handler: func() handlers.QueryHandler[oracletypes.CurrencyPair, *big.Int] {
				handler := handlermocks.NewQueryHandler[oracletypes.CurrencyPair, *big.Int](t)

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
				}).After(cfg.Interval * 2)

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
				logger,
				cfg,
				tc.handler(),
				tc.pairs,
			)
			require.NoError(t, err)

			now := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), cfg.Interval*5)
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
