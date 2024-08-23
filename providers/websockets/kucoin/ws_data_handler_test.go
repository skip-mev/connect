package kucoin_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	providertypes "github.com/skip-mev/connect/v2/providers/types"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
	"github.com/skip-mev/connect/v2/providers/websockets/kucoin"
)

var (
	btcusdt = types.DefaultProviderTicker{
		OffChainTicker: "BTC-USDT",
	}
	ethusdt = types.DefaultProviderTicker{
		OffChainTicker: "ETH-USDT",
	}
	mogusdt = types.DefaultProviderTicker{
		OffChainTicker: "MOG-USDT",
	}
	logger = zap.NewExample()
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
			resp: types.PriceResponse{
				Resolved: types.ResolvedPrices{
					btcusdt: {
						Value:        big.NewFloat(0),
						ResponseCode: providertypes.ResponseCodeUnchanged,
					},
					ethusdt: {
						Value:        big.NewFloat(0),
						ResponseCode: providertypes.ResponseCodeUnchanged,
					},
				},
			},
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
					btcusdt: {
						Value: big.NewFloat(0.1),
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
					btcusdt: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("error"), providertypes.ErrorWebSocketGeneral),
					},
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
					btcusdt: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("received out of order ticker response message"), providertypes.ErrorWebSocketGeneral),
					},
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
					btcusdt: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("failed to parse price %s", "failed to parse float64 string to big int: invalid"), providertypes.ErrorWebSocketGeneral),
					},
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

	handler, err := kucoin.NewWebSocketDataHandler(logger, kucoin.DefaultWebSocketConfig)
	require.NoError(t, err)

	// Update the cache since it is assumed that CreateMessages is executed before anything else.
	_, err = handler.CreateMessages([]types.ProviderTicker{btcusdt, ethusdt})
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
			require.LessOrEqual(t, len(resp.Resolved), 2)
			require.LessOrEqual(t, len(resp.UnResolved), 2)

			require.Equal(t, tc.updateMsg(), updateMsg)

			for cp, result := range tc.resp.Resolved {
				require.Contains(t, resp.Resolved, cp)
				require.Equal(t, result.Value.SetPrec(18), resp.Resolved[cp].Value.SetPrec(18))
				require.Equal(t, result.ResponseCode, resp.Resolved[cp].ResponseCode)
			}

			for cp := range tc.resp.UnResolved {
				require.Contains(t, resp.UnResolved, cp)
				require.Error(t, resp.UnResolved[cp])
			}
		})
	}
}

func TestCreateMessages(t *testing.T) {
	batchCfg := kucoin.DefaultWebSocketConfig
	batchCfg.MaxSubscriptionsPerBatch = 2

	testCases := []struct {
		name        string
		cps         []types.ProviderTicker
		cfg         config.WebSocketConfig
		expected    func() []handlers.WebsocketEncodedMessage
		expectedErr bool
	}{
		{
			name: "no currency pairs",
			cps:  []types.ProviderTicker{},
			cfg:  kucoin.DefaultWebSocketConfig,
			expected: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "one currency pair",
			cps: []types.ProviderTicker{
				btcusdt,
			},
			cfg: kucoin.DefaultWebSocketConfig,
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
			cps: []types.ProviderTicker{
				btcusdt,
				ethusdt,
			},
			cfg: kucoin.DefaultWebSocketConfig,
			expected: func() []handlers.WebsocketEncodedMessage {
				msgs := make([]handlers.WebsocketEncodedMessage, 2)
				for i, ticker := range []string{"BTC-USDT", "ETH-USDT"} {
					msg := kucoin.SubscribeRequestMessage{
						Type: string(kucoin.SubscribeMessage),
						Topic: fmt.Sprintf(
							"%s%s",
							kucoin.TickerTopic,
							ticker,
						),
						PrivateChannel: false,
						Response:       false,
					}

					bz, err := json.Marshal(msg)
					require.NoError(t, err)
					msgs[i] = bz
				}

				return msgs
			},
			expectedErr: false,
		},
		{
			name: "multiple currency pairs with batch config",
			cps: []types.ProviderTicker{
				btcusdt,
				ethusdt,
			},
			cfg: batchCfg,
			expected: func() []handlers.WebsocketEncodedMessage {
				msgs := make([]handlers.WebsocketEncodedMessage, 1)
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
				msgs[0] = bz

				return msgs
			},
			expectedErr: false,
		},
		{
			name: "multiple currency pairs with batch config and remainder",
			cps: []types.ProviderTicker{
				btcusdt,
				ethusdt,
				mogusdt,
			},
			cfg: batchCfg,
			expected: func() []handlers.WebsocketEncodedMessage {
				msgs := make([]handlers.WebsocketEncodedMessage, 2)
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
				msgs[0] = bz

				msg = kucoin.SubscribeRequestMessage{
					Type: string(kucoin.SubscribeMessage),
					Topic: fmt.Sprintf(
						"%s%s",
						kucoin.TickerTopic,
						"MOG-USDT",
					),
					PrivateChannel: false,
					Response:       false,
				}

				bz, err = json.Marshal(msg)
				require.NoError(t, err)
				msgs[1] = bz

				return msgs
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := kucoin.NewWebSocketDataHandler(logger, tc.cfg)
			require.NoError(t, err)

			actual, err := handler.CreateMessages(tc.cps)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			expected := tc.expected()
			require.Len(t, actual, len(expected))

			for i := range expected {
				var expectedMsg kucoin.SubscribeRequestMessage
				require.NoError(t, json.Unmarshal(expected[i], &expectedMsg))
				var actualMsg kucoin.SubscribeRequestMessage
				require.NoError(t, json.Unmarshal(actual[i], &actualMsg))

				require.Equal(t, expectedMsg.Type, actualMsg.Type)
				require.Equal(t, expectedMsg.Topic, actualMsg.Topic)
				require.Equal(t, expectedMsg.PrivateChannel, actualMsg.PrivateChannel)
				require.Equal(t, expectedMsg.Response, actualMsg.Response)
			}
		})
	}
}
