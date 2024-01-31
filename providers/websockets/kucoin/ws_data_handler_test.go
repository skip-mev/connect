package kucoin_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/kucoin"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var logger = zap.NewExample()

func TestHandleMessage(t *testing.T) {
	testCases := []struct {
		name        string
		msg         func() []byte
		resp        providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
		updateMsg   func() []handlers.WebsocketEncodedMessage
		expectedErr bool
	}{
		{
			name: "invalid message",
			msg: func() []byte {
				return []byte("invalid")
			},
			resp:        providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
			updateMsg:   func() []handlers.WebsocketEncodedMessage { return nil },
			expectedErr: true,
		},
		{
			name: "welcome message",
			msg: func() []byte {
				return []byte(`{
					"id": "id",
					"type": "welcome"
				}`)
			},
			resp:        providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
			updateMsg:   func() []handlers.WebsocketEncodedMessage { return nil },
			expectedErr: false,
		},
		{
			name: "pong message",
			msg: func() []byte {
				return []byte(`{
					"id": "id",
					"type": "pong"
				}`)
			},
			resp:        providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
			updateMsg:   func() []handlers.WebsocketEncodedMessage { return nil },
			expectedErr: false,
		},
		{
			name: "subscription response message",
			msg: func() []byte {
				return []byte(`{
					"id": "id",
					"type": "ack"
				}`)
			},
			resp:        providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
			updateMsg:   func() []handlers.WebsocketEncodedMessage { return nil },
			expectedErr: false,
		},
		{
			name: "unknown message type",
			msg: func() []byte {
				return []byte(`{
					"id": "id",
					"type": "unknown"
				}`)
			},
			resp:        providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
			updateMsg:   func() []handlers.WebsocketEncodedMessage { return nil },
			expectedErr: true,
		},
		{
			name: "invalid ticker response message",
			msg: func() []byte {
				return []byte(`{
					"id": "id",
					"type": "message",
					"topic": "/market/ticker:BTC-USDT",
					"data": "invalid"
				}`)
			},
			resp:        providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
			updateMsg:   func() []handlers.WebsocketEncodedMessage { return nil },
			expectedErr: true,
		},
		{
			name: "valid ticker response message",
			msg: func() []byte {
				return []byte(`{
					"id": "id",
					"type": "message",
					"topic": "/market/ticker:BTC-USDT",
					"subject": "trade.ticker",
					"data": {
						"sequence": "1",
						"price": "0.1"
					}
				}`)
			},
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
				Resolved: map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals): {
						Value: big.NewInt(10000000),
					},
				},
			},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: false,
		},
		{
			name: "duplicate valid ticker response message",
			msg: func() []byte {
				return []byte(`{
					"id": "id",
					"type": "message",
					"topic": "/market/ticker:BTC-USDT",
					"subject": "trade.ticker",
					"data": {
						"sequence": "1",
						"price": "0.1"
					}
				}`)
			},
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
				UnResolved: map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals): fmt.Errorf("err"),
				},
			},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "unable to parse sequence",
			msg: func() []byte {
				return []byte(`{
					"id": "id",
					"type": "message",
					"topic": "/market/ticker:BTC-USDT",
					"subject": "trade.ticker",
					"data": {
						"sequence": "mog",
						"price": "0.1"
					}
				}`)
			},
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
				UnResolved: map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals): fmt.Errorf("received out of order ticker response message"),
				},
			},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "missing ticker data",
			msg: func() []byte {
				return []byte(`{
					"id": "id",
					"type": "message",
					"topic": "/market/ticker",
					"subject": "trade.ticker",
					"data": {
						"price": "0.1"
					}
				}`)
			},
			resp:        providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
			updateMsg:   func() []handlers.WebsocketEncodedMessage { return nil },
			expectedErr: true,
		},
		{
			name: "invalid ticker data",
			msg: func() []byte {
				return []byte(`{
					"id": "id",
					"type": "message",
					"topic": "/market/ticker:BTC-USDT",
					"subject": "trade.ticker",
					"data": {
						"price": "invalid"
					}
				}`)
			},
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
				UnResolved: map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals): fmt.Errorf("failed to parse price %s", "failed to parse float64 string to big int: invalid"),
				},
			},
			updateMsg:   func() []handlers.WebsocketEncodedMessage { return nil },
			expectedErr: true,
		},
		{
			name: "invalid ticker subject",
			msg: func() []byte {
				return []byte(`{
					"id": "id",
					"type": "message",
					"topic": "/market/ticker:BTC-USDT",
					"subject": "trade.futures",
					"data": {
						"price": "0.1"
					}
				}`)
			},
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "invalid ticker topic",
			msg: func() []byte {
				return []byte(`{
					"id": "id",
					"type": "message",
					"topic": "/market/ticker:MOG-USDT",
					"subject": "trade.ticker",
					"data": {
						"price": "0.1"
					}
				}`)
			},
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
	}

	handler, err := kucoin.NewWebSocketDataHandler(logger, providerCfg)
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, updateMsg, err := handler.HandleMessage(tc.msg())
			if tc.expectedErr {
				require.Error(t, err)

				require.LessOrEqual(t, len(resp.UnResolved), 1)
				for cp := range tc.resp.UnResolved {
					require.Contains(t, resp.UnResolved, cp)
					require.Error(t, resp.UnResolved[cp])
				}
				return
			}

			require.NoError(t, err)

			// The response should contain a single resolved price update.
			require.LessOrEqual(t, len(resp.Resolved), 1)
			require.LessOrEqual(t, len(resp.UnResolved), 1)

			require.Equal(t, tc.updateMsg(), updateMsg)

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

func TestCreateMessages(t *testing.T) {
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
			expectedErr: true,
		},
		{
			name: "one currency pair",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := kucoin.SubscribeRequestMessage{
					Type: string(kucoin.SubscribeMessage),
					Topic: fmt.Sprintf(
						"%s%s",
						kucoin.TickerTopic,
						"BTC-USDT",
					),
					PrivateChannel: false,
					Response:       false,
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{bz}
			},
			expectedErr: false,
		},
		{
			name: "multiple currency pairs",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
				oracletypes.NewCurrencyPair("ETHEREUM", "USD", oracletypes.DefaultDecimals),
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := kucoin.SubscribeRequestMessage{
					Type: string(kucoin.SubscribeMessage),
					Topic: fmt.Sprintf(
						"%s%s,%s",
						kucoin.TickerTopic,
						"BTC-USDT",
						"ETH-USDT",
					),
					PrivateChannel: false,
					Response:       false,
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{bz}
			},
			expectedErr: false,
		},
		{
			name: "multiple currency pairs with one not found in market configs",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
				oracletypes.NewCurrencyPair("ETHEREUM", "USD", oracletypes.DefaultDecimals),
				oracletypes.NewCurrencyPair("MOG", "USD", oracletypes.DefaultDecimals),
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := kucoin.SubscribeRequestMessage{
					Type: string(kucoin.SubscribeMessage),
					Topic: fmt.Sprintf(
						"%s%s,%s",
						kucoin.TickerTopic,
						"BTC-USDT",
						"ETH-USDT",
					),
					PrivateChannel: false,
					Response:       false,
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{bz}
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := kucoin.NewWebSocketDataHandler(logger, providerCfg)
			require.NoError(t, err)

			actual, err := handler.CreateMessages(tc.cps)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, actual, 1)

			var expected kucoin.SubscribeRequestMessage
			require.NoError(t, json.Unmarshal(tc.expected()[0], &expected))

			var actualMsg kucoin.SubscribeRequestMessage
			require.NoError(t, json.Unmarshal(actual[0], &actualMsg))

			require.Equal(t, expected.Type, actualMsg.Type)
			require.Equal(t, expected.Topic, actualMsg.Topic)
			require.Equal(t, expected.PrivateChannel, actualMsg.PrivateChannel)
			require.Equal(t, expected.Response, actualMsg.Response)
		})
	}
}
