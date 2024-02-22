package kucoin_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	"github.com/skip-mev/slinky/providers/websockets/kucoin"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	logger = zap.NewExample()
	mogusd = mmtypes.NewTicker("MOG", "USD", 8, 1)
)

func TestHandleMessage(t *testing.T) {
	testCases := []struct {
		name        string
		msg         func() []byte
		resp        types.PriceResponse
		updateMsg   func() []handlers.WebsocketEncodedMessage
		expectedErr bool
	}{
		{
			name: "invalid message",
			msg: func() []byte {
				return []byte("invalid")
			},
			resp:        types.PriceResponse{},
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
			resp:        types.PriceResponse{},
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
			resp:        types.PriceResponse{},
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
			resp:        types.PriceResponse{},
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
			resp:        types.PriceResponse{},
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
			resp:        types.PriceResponse{},
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
			resp: types.PriceResponse{
				Resolved: types.ResolvedPrices{
					constants.BITCOIN_USDT: {
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
			resp: types.PriceResponse{
				UnResolved: types.UnResolvedPrices{
					constants.BITCOIN_USDT: fmt.Errorf("err"),
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
			resp: types.PriceResponse{
				UnResolved: types.UnResolvedPrices{
					constants.BITCOIN_USDT: fmt.Errorf("received out of order ticker response message"),
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
			resp:        types.PriceResponse{},
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
			resp: types.PriceResponse{
				UnResolved: types.UnResolvedPrices{
					constants.BITCOIN_USDT: fmt.Errorf("failed to parse price %s", "failed to parse float64 string to big int: invalid"),
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
			resp: types.PriceResponse{},
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
			resp: types.PriceResponse{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
	}

	marketConfig, err := types.NewProviderMarketMap(kucoin.Name, kucoin.DefaultMarketConfig)
	require.NoError(t, err)

	handler, err := kucoin.NewWebSocketDataHandler(logger, marketConfig, kucoin.DefaultWebSocketConfig)
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
		cps         []mmtypes.Ticker
		expected    func() []handlers.WebsocketEncodedMessage
		expectedErr bool
	}{
		{
			name: "no currency pairs",
			cps:  []mmtypes.Ticker{},
			expected: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "one currency pair",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
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
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
				constants.ETHEREUM_USDT,
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
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
				constants.ETHEREUM_USDT,
				mogusd,
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			marketConfig, err := types.NewProviderMarketMap(kucoin.Name, kucoin.DefaultMarketConfig)
			require.NoError(t, err)

			handler, err := kucoin.NewWebSocketDataHandler(logger, marketConfig, kucoin.DefaultWebSocketConfig)
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
