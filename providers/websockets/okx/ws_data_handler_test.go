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
)

var (
	btc_usdt = okx.DefaultMarketConfig.MustGetProviderTicker(constants.BITCOIN_USDT)
	eth_usdt = okx.DefaultMarketConfig.MustGetProviderTicker(constants.ETHEREUM_USDT)
	logger   = zap.NewExample()
)

func TestHandleMessage(t *testing.T) {
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
					btc_usdt: {
						Value: big.NewFloat(1.0),
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
					btc_usdt: {
						Value: big.NewFloat(1.0),
					},
					eth_usdt: {
						Value: big.NewFloat(2.0),
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
			wsHandler, err := okx.NewWebSocketDataHandler(logger, okx.DefaultWebSocketConfig)
			require.NoError(t, err)

			// Update the cache since it is assumed that CreateMessages is executed before anything else.
			_, err = wsHandler.CreateMessages([]types.ProviderTicker{btc_usdt, eth_usdt})
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
				require.Equal(t, result.Value.SetPrec(18), resp.Resolved[cp].Value.SetPrec(18))
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
		cps         []types.ProviderTicker
		expected    func() []handlers.WebsocketEncodedMessage
		expectedErr bool
	}{
		{
			name: "no currency pairs",
			cps:  []types.ProviderTicker{},
			expected: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "one currency pair",
			cps: []types.ProviderTicker{
				btc_usdt,
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
			cps: []types.ProviderTicker{
				btc_usdt,
				eth_usdt,
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wsHandler, err := okx.NewWebSocketDataHandler(logger, okx.DefaultWebSocketConfig)
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
