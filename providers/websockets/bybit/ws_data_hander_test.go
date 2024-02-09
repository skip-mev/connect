package bybit_test

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"

	"github.com/skip-mev/slinky/providers/websockets/bybit"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	logger = zap.NewExample()

	cfg = config.ProviderConfig{
		Name: bybit.Name,
		WebSocket: config.WebSocketConfig{
			Enabled:             true,
			MaxBufferSize:       1024,
			ReconnectionTimeout: 10 * time.Second,
			WSS:                 bybit.URLProd,
			Name:                bybit.Name,
			ReadBufferSize:      config.DefaultReadBufferSize,
			WriteBufferSize:     config.DefaultWriteBufferSize,
			HandshakeTimeout:    config.DefaultHandshakeTimeout,
			EnableCompression:   config.DefaultEnableCompression,
			ReadTimeout:         config.DefaultReadTimeout,
			WriteTimeout:        config.DefaultWriteTimeout,
		},
		Market: config.MarketConfig{
			Name: bybit.Name,
			CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
				"BITCOIN/USD": {
					Ticker:       "BTCUSD",
					CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
				},
				"ETHEREUM/USD": {
					Ticker:       "ETHUSD",
					CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
				},
			},
		},
	}
)

func TestHandlerMessage(t *testing.T) {
	testCases := []struct {
		name      string
		msg       func() []byte
		resp      providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]
		updateMsg func() []handlers.WebsocketEncodedMessage
		expErr    bool
	}{
		{
			name: "invalid message",
			msg: func() []byte {
				return []byte("invalid message")
			},
			resp:      providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](nil, nil),
			updateMsg: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:    true,
		},
		{
			name: "invalid message type",
			msg: func() []byte {
				msg := bybit.BaseResponse{
					Op: "unknown",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp:      providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](nil, nil),
			updateMsg: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:    true,
		},
		{
			name: "price update",
			msg: func() []byte {
				msg := bybit.TickerUpdateMessage{
					Topic: "tickers.BTCUSD",
					Data: bybit.TickerUpdateData{
						Symbol:    "BTCUSD",
						LastPrice: "1",
					},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): {
						Value: big.NewInt(100000000),
					},
				},
				map[slinkytypes.CurrencyPair]error{},
			),
			updateMsg: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:    false,
		},
		{
			name: "price update with unknown pair ID",
			msg: func() []byte {
				msg := bybit.TickerUpdateMessage{
					Topic: "tickers.MOGUSDT",
					Data: bybit.TickerUpdateData{
						Symbol:    "MOGUSDT",
						LastPrice: "1",
					},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[slinkytypes.CurrencyPair]error{},
			),
			updateMsg: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:    false,
		},
		{
			name: "successful subscription",
			msg: func() []byte {
				msg := bybit.SubscriptionResponse{
					BaseResponse: bybit.BaseResponse{
						Success: true,
						RetMsg:  string(bybit.OperationSubscribe),
						ConnID:  "90190u1309",
						Op:      string(bybit.OperationSubscribe),
					},
					ReqID: "1012901",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[slinkytypes.CurrencyPair]error{},
			),
			updateMsg: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:    false,
		},
		{
			name: "subscription error",
			msg: func() []byte {
				msg := bybit.SubscriptionResponse{
					BaseResponse: bybit.BaseResponse{
						Success: false,
						RetMsg:  string(bybit.OperationSubscribe),
						ConnID:  "90190u1309",
						Op:      string(bybit.OperationSubscribe),
					},
					ReqID: "1012901",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[slinkytypes.CurrencyPair]error{},
			),
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wsHandler, err := bybit.NewWebSocketDataHandler(logger, cfg)
			require.NoError(t, err)

			resp, updateMsg, err := wsHandler.HandleMessage(tc.msg())
			if tc.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tc.updateMsg(), updateMsg)

			require.Equal(t, len(tc.resp.Resolved), len(resp.Resolved))
			require.Equal(t, len(tc.resp.UnResolved), len(resp.UnResolved))

			for cp, result := range tc.resp.Resolved {
				require.Contains(t, resp.Resolved, cp)
				require.Equal(t, result.Value, resp.Resolved[cp].Value)
			}

			for cp := range tc.resp.UnResolved {
				require.Contains(t, resp.UnResolved, cp)
				require.Error(t, resp.UnResolved[cp])
			}
		})
	}
}

func TestCreateMessage(t *testing.T) {
	testCases := []struct {
		name        string
		cps         []slinkytypes.CurrencyPair
		expected    func() []byte
		expectedErr bool
	}{
		{
			name: "no currency pairs",
			cps:  []slinkytypes.CurrencyPair{},
			expected: func() []byte {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "one currency pair",
			cps: []slinkytypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			expected: func() []byte {
				msg := bybit.SubscriptionRequest{
					BaseRequest: bybit.BaseRequest{
						Op: string(bybit.OperationSubscribe),
					},
					Args: []string{"tickers.BTCUSD"},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			expectedErr: false,
		},
		{
			name: "two currency pairs",
			cps: []slinkytypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			expected: func() []byte {
				msg := bybit.SubscriptionRequest{
					BaseRequest: bybit.BaseRequest{
						Op: string(bybit.OperationSubscribe),
					},
					Args: []string{"tickers.BTCUSD", "tickers.ETHUSD"},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			expectedErr: false,
		},
		{
			name: "one currency pair not in config",
			cps: []slinkytypes.CurrencyPair{
				oracletypes.NewCurrencyPair("MOG", "USD"),
			},
			expected: func() []byte {
				return nil
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := bybit.NewWebSocketDataHandler(logger, cfg)
			require.NoError(t, err)

			msgs, err := handler.CreateMessages(tc.cps)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, 1, len(msgs))
			require.EqualValues(t, tc.expected(), []byte(msgs[0]))
		})
	}
}
