package huobi_test

import (
	"bytes"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/klauspost/compress/gzip"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/huobi"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	providerCfg = config.ProviderConfig{
		Name:      huobi.Name,
		WebSocket: huobi.DefaultWebSocketConfig,
		Market: config.MarketConfig{
			Name: huobi.Name,
			CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
				"BITCOIN/USDT/8": {
					Ticker:       "btcusdt",
					CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT", oracletypes.DefaultDecimals),
				},
				"ETHEREUM/USDT/8": {
					Ticker:       "ethusdt",
					CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT", oracletypes.DefaultDecimals),
				},
			},
		},
	}

	logger = zap.NewExample()
)

func TestHandlerMessage(t *testing.T) {
	testCases := []struct {
		name          string
		msg           func() []byte
		resp          providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
		updateMessage func() []handlers.WebsocketEncodedMessage
		expErr        bool
	}{
		{
			name: "invalid message",
			msg: func() []byte {
				return []byte("invalid message")
			},
			resp:          providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](nil, nil),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
		{
			name: "invalid message subbed type",
			msg: func() []byte {
				msg := huobi.SubscriptionResponse{
					ID:     "test",
					Status: "ok",
					Subbed: "invalid",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				var buf bytes.Buffer
				zw := gzip.NewWriter(&buf)

				_, err = zw.Write(bz)
				require.NoError(t, err)
				require.NoError(t, zw.Close())

				return buf.Bytes()
			},
			resp:          providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](nil, nil),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
		{
			name: "ticker price update",
			msg: func() []byte {
				msg := huobi.TickerStream{
					Channel: "market.btcusdt.ticker",
					Tick:    huobi.Tick{LastPrice: 1},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				var buf bytes.Buffer
				zw := gzip.NewWriter(&buf)

				_, err = zw.Write(bz)
				require.NoError(t, err)
				require.NoError(t, zw.Close())

				return buf.Bytes()
			},
			resp: providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					oracletypes.NewCurrencyPair("BITCOIN", "USDT", oracletypes.DefaultDecimals): {
						Value: big.NewInt(100000000),
					},
				},
				map[oracletypes.CurrencyPair]error{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        false,
		},
		{
			name: "ticker price update with unknown ticker Ticker",
			msg: func() []byte {
				msg := huobi.TickerStream{
					Channel: "unknown",
					Tick:    huobi.Tick{LastPrice: 1},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				var buf bytes.Buffer
				zw := gzip.NewWriter(&buf)

				_, err = zw.Write(bz)
				require.NoError(t, err)
				require.NoError(t, zw.Close())

				return buf.Bytes()
			},
			resp: providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
		{
			name: "successful subscription",
			msg: func() []byte {
				msg := huobi.SubscriptionResponse{
					ID:     "test",
					Status: "ok",
					Subbed: "market.btcusdt.ticker",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				var buf bytes.Buffer
				zw := gzip.NewWriter(&buf)

				_, err = zw.Write(bz)
				require.NoError(t, err)
				require.NoError(t, zw.Close())

				return buf.Bytes()
			},
			resp: providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        false,
		},
		{
			name: "subscription error",
			msg: func() []byte {
				msg := huobi.SubscriptionResponse{
					ID:     "test",
					Status: "notok",
					Subbed: "market.btcusdt.ticker",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				var buf bytes.Buffer
				zw := gzip.NewWriter(&buf)

				_, err = zw.Write(bz)
				require.NoError(t, err)
				require.NoError(t, zw.Close())

				return buf.Bytes()
			},
			resp: providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				msg, err := huobi.NewSubscriptionRequest("btcusdt")
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{msg}
			},
			expErr: false,
		},
		{
			name: "valid heartbeat",
			msg: func() []byte {
				msg := huobi.PingMessage{Ping: 123}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				var buf bytes.Buffer
				zw := gzip.NewWriter(&buf)

				_, err = zw.Write(bz)
				require.NoError(t, err)
				require.NoError(t, zw.Close())

				return buf.Bytes()
			},
			resp: providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				msg := huobi.PongMessage{Pong: 123}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{bz}
			},
			expErr: false,
		},
		{
			name: "invalid empty heartbeat",
			msg: func() []byte {
				msg := huobi.PingMessage{}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				var buf bytes.Buffer
				zw := gzip.NewWriter(&buf)

				_, err = zw.Write(bz)
				require.NoError(t, err)
				require.NoError(t, zw.Close())

				return buf.Bytes()
			},
			resp: providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wsHandler, err := huobi.NewWebSocketDataHandler(logger, providerCfg)
			require.NoError(t, err)

			resp, updateMsg, err := wsHandler.HandleMessage(tc.msg())
			if tc.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.updateMessage(), updateMsg)

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
		cps         []oracletypes.CurrencyPair
		expected    func() []handlers.WebsocketEncodedMessage
		expectedErr bool
	}{
		{
			name: "no currency pairs",
			cps:  []oracletypes.CurrencyPair{},
			expected: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: false,
		},
		{
			name: "one currency pair",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT", oracletypes.DefaultDecimals),
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg, err := huobi.NewSubscriptionRequest("btcusdt")
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{msg}
			},
			expectedErr: false,
		},
		{
			name: "two currency pairs",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT", oracletypes.DefaultDecimals),
				oracletypes.NewCurrencyPair("ETHEREUM", "USDT", oracletypes.DefaultDecimals),
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				bz1, err := huobi.NewSubscriptionRequest("btcusdt")
				require.NoError(t, err)
				bz2, err := huobi.NewSubscriptionRequest("ethusdt")
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{bz1, bz2}
			},
			expectedErr: false,
		},
		{
			name: "one currency pair not in config",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("MOG", "USDT", oracletypes.DefaultDecimals),
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := huobi.NewWebSocketDataHandler(logger, providerCfg)
			require.NoError(t, err)

			msgs, err := handler.CreateMessages(tc.cps)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, tc.expected(), msgs)
		})
	}
}
