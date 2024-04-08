package gate_test

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	"github.com/skip-mev/slinky/providers/websockets/gate"
)

var (
	btc_usdt = gate.DefaultMarketConfig.MustGetProviderTicker(constants.BITCOIN_USDT)
	eth_usdt = gate.DefaultMarketConfig.MustGetProviderTicker(constants.ETHEREUM_USDT)
	logger   = zap.NewExample()
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
			name: "invalid base message type",
			msg: func() []byte {
				msg := gate.BaseMessage{
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
			name: "unsupported ticker stream channel",
			msg: func() []byte {
				msg := gate.TickerStream{
					BaseMessage: gate.BaseMessage{
						Time:    0,
						Channel: "unknown",
						Event:   string(gate.EventUpdate),
					},
					Result: gate.TickerResult{
						CurrencyPair: "BTC_USDT",
						Last:         "1",
					},
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
			name: "ticker stream price update",
			msg: func() []byte {
				msg := gate.TickerStream{
					BaseMessage: gate.BaseMessage{
						Time:    0,
						Channel: string(gate.ChannelTickers),
						Event:   string(gate.EventUpdate),
					},
					Result: gate.TickerResult{
						CurrencyPair: "BTC_USDT",
						Last:         "1",
					},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{
					btc_usdt: {
						Value: big.NewFloat(1e18),
					},
				},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        false,
		},
		{
			name: "price update with unknown currency pair",
			msg: func() []byte {
				msg := gate.TickerStream{
					BaseMessage: gate.BaseMessage{
						Time:    0,
						Channel: string(gate.ChannelTickers),
						Event:   string(gate.EventUpdate),
					},
					Result: gate.TickerResult{
						CurrencyPair: "MOG_USDT",
						Last:         "1",
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
			expErr:        true,
		},
		{
			name: "successful subscription",
			msg: func() []byte {
				msg := gate.SubscribeResponse{
					BaseMessage: gate.BaseMessage{
						Time:    0,
						Channel: string(gate.ChannelTickers),
						Event:   string(gate.EventSubscribe),
					},
					ID: 0,
					Error: gate.ErrorMessage{
						Code:    0,
						Message: "",
					},
					Result: gate.RequestResult{Status: string(gate.StatusSuccess)},
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
				msg := gate.SubscribeResponse{
					BaseMessage: gate.BaseMessage{
						Time:    0,
						Channel: string(gate.ChannelTickers),
						Event:   string(gate.EventSubscribe),
					},
					ID: 0,
					Error: gate.ErrorMessage{
						Code:    int(gate.ErrorInvalidRequestBody),
						Message: "invalid body",
					},
					Result: gate.RequestResult{Status: "error"},
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wsHandler, err := gate.NewWebSocketDataHandler(logger, gate.DefaultWebSocketConfig)
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
				require.Equal(t, result.Value.SetPrec(types.DefaultTickerDecimals), resp.Resolved[cp].Value.SetPrec(types.DefaultTickerDecimals))
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
				msg := gate.SubscribeRequest{
					BaseMessage: gate.BaseMessage{
						Time:    time.Now().Second(),
						Channel: string(gate.ChannelTickers),
						Event:   string(gate.EventSubscribe),
					},
					ID:      time.Now().Second(),
					Payload: []string{"BTC_USDT"},
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
				msg := gate.SubscribeRequest{
					BaseMessage: gate.BaseMessage{
						Time:    time.Now().Second(),
						Channel: string(gate.ChannelTickers),
						Event:   string(gate.EventSubscribe),
					},
					ID:      time.Now().Second(),
					Payload: []string{"BTC_USDT", "ETH_USDT"},
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
			handler, err := gate.NewWebSocketDataHandler(logger, gate.DefaultWebSocketConfig)
			require.NoError(t, err)

			msgs, err := handler.CreateMessages(tc.cps)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			var (
				gotMsg      gate.SubscribeRequest
				expectedMsg gate.SubscribeRequest
			)

			// need to check the non-time based fields
			err = json.Unmarshal(msgs[0], &gotMsg)
			require.NoError(t, err)
			err = json.Unmarshal(tc.expected()[0], &expectedMsg)
			require.NoError(t, err)

			require.Equal(t, expectedMsg.Event, gotMsg.Event)
			require.Equal(t, expectedMsg.Channel, gotMsg.Channel)
			require.Equal(t, expectedMsg.Payload, gotMsg.Payload)
		})
	}
}
