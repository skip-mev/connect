package okx_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	logger = zap.NewExample()
	mogusd = mmtypes.NewTicker("MOG", "USDT", 8, 1)
)

func TestHandlerMessage(t *testing.T) {
	testCases := []struct {
		name          string
		msg           func() []byte
		resp          types.PriceResponse
		updateMessage func() []handlers.WebsocketEncodedMessage
		expErr        bool
	}{
		{
			name: "invalid message",
			msg: func() []byte {
				return []byte("invalid message")
			},
			resp:          types.NewPriceResponse(nil, nil),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
		{
			name: "invalid message type",
			msg: func() []byte {
				msg := okx.BaseMessage{
					Event: "unknown",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp:          types.NewPriceResponse(nil, nil),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
		{
			name: "instrument price update",
			msg: func() []byte {
				msg := okx.IndexTickersResponseMessage{
					Arguments: okx.SubscriptionTopic{
						Channel:      string(okx.IndexTickersChannel),
						InstrumentID: "BTC-USDT",
					},
					Data: []okx.IndexTicker{
						{
							ID:         "BTC-USDT",
							IndexPrice: "1",
						},
					},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{
					constants.BITCOIN_USDT: {
						Value: big.NewInt(100000000),
					},
				},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        false,
		},
		{
			name: "multiple instruments included in the response",
			msg: func() []byte {
				msg := okx.IndexTickersResponseMessage{
					Arguments: okx.SubscriptionTopic{
						Channel:      string(okx.IndexTickersChannel),
						InstrumentID: "BTC-USDT",
					},
					Data: []okx.IndexTicker{
						{
							ID:         "BTC-USDT",
							IndexPrice: "1",
						},
						{
							ID:         "ETH-USDT",
							IndexPrice: "2",
						},
					},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{
					constants.BITCOIN_USDT: {
						Value: big.NewInt(100000000),
					},
					constants.ETHEREUM_USDT: {
						Value: big.NewInt(200000000),
					},
				},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        false,
		},
		{
			name: "instrument price update with unknown instrument ID",
			msg: func() []byte {
				msg := okx.IndexTickersResponseMessage{
					Arguments: okx.SubscriptionTopic{
						Channel:      string(okx.IndexTickersChannel),
						InstrumentID: "MOG-USDT",
					},
					Data: []okx.IndexTicker{
						{
							ID:         "MOG-USDT",
							IndexPrice: "1",
						},
					},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        false,
		},
		{
			name: "successful subscription",
			msg: func() []byte {
				msg := okx.SubscribeResponseMessage{
					Arguments: okx.SubscriptionTopic{
						Channel:      string(okx.IndexTickersChannel),
						InstrumentID: "BTC-USDT",
					},
					Event:        string(okx.EventSubscribe),
					ConnectionID: "123",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        false,
		},
		{
			name: "subscription error",
			msg: func() []byte {
				initSubMessage := okx.SubscribeRequestMessage{
					Operation: string(okx.OperationSubscribe),
					Arguments: []okx.SubscriptionTopic{
						{
							Channel:      string(okx.IndexTickersChannel),
							InstrumentID: "BTC-USDT",
						},
					},
				}

				bz, err := json.Marshal(initSubMessage)
				require.NoError(t, err)
				errMsg := fmt.Sprintf("%s%s", okx.ExpectedErrorPrefix, string(bz))

				msg := okx.SubscribeResponseMessage{
					Event:        string(okx.EventError),
					Code:         "123",
					Message:      errMsg,
					ConnectionID: "123",
				}

				bz, err = json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				msg := okx.SubscribeRequestMessage{
					Operation: string(okx.OperationSubscribe),
					Arguments: []okx.SubscriptionTopic{
						{
							Channel:      string(okx.IndexTickersChannel),
							InstrumentID: "BTC-USDT",
						},
					},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{bz}
			},
			expErr: false,
		},
		{
			name: "subscription error with invalid message",
			msg: func() []byte {
				msg := okx.SubscribeResponseMessage{
					Event:        string(okx.EventError),
					Code:         "123",
					Message:      "invalidmessage",
					ConnectionID: "123",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
		{
			name: "subscription error with invalid message format",
			msg: func() []byte {
				initSubMessage := okx.SubscribeRequestMessage{
					Operation: string(okx.OperationSubscribe),
					Arguments: []okx.SubscriptionTopic{
						{
							Channel:      string(okx.IndexTickersChannel),
							InstrumentID: "BTC-USDT",
						},
					},
				}

				bz, err := json.Marshal(initSubMessage)
				require.NoError(t, err)
				errMsg := fmt.Sprintf("%s%s", okx.ExpectedErrorPrefix, string(bz)+"invalidmessage")

				msg := okx.SubscribeResponseMessage{
					Event:        string(okx.EventError),
					Code:         "123",
					Message:      errMsg,
					ConnectionID: "123",
				}

				bz, err = json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			marketConfig, err := types.NewProviderMarketMap(okx.Name, okx.DefaultMarketConfig)
			require.NoError(t, err)

			wsHandler, err := okx.NewWebSocketDataHandler(logger, marketConfig, okx.DefaultWebSocketConfig)
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
				msg := okx.SubscribeRequestMessage{
					Operation: string(okx.OperationSubscribe),
					Arguments: []okx.SubscriptionTopic{
						{
							Channel:      string(okx.IndexTickersChannel),
							InstrumentID: "BTC-USDT",
						},
					},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{bz}
			},
			expectedErr: false,
		},
		{
			name: "two currency pairs",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
				constants.ETHEREUM_USDT,
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := okx.SubscribeRequestMessage{
					Operation: string(okx.OperationSubscribe),
					Arguments: []okx.SubscriptionTopic{
						{
							Channel:      string(okx.IndexTickersChannel),
							InstrumentID: "BTC-USDT",
						},
						{
							Channel:      string(okx.IndexTickersChannel),
							InstrumentID: "ETH-USDT",
						},
					},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{bz}
			},
			expectedErr: false,
		},
		{
			name: "one currency pair not in config",
			cps: []mmtypes.Ticker{
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
			marketConfig, err := types.NewProviderMarketMap(okx.Name, okx.DefaultMarketConfig)
			require.NoError(t, err)

			wsHandler, err := okx.NewWebSocketDataHandler(logger, marketConfig, okx.DefaultWebSocketConfig)
			require.NoError(t, err)

			msgs, err := wsHandler.CreateMessages(tc.cps)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, tc.expected(), msgs)
		})
	}
}
