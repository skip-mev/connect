package bitfinex_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/bitfinex"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	mogusd     = mmtypes.NewTicker("MOG", "USD", 8, 1)
	channelBTC = 111
	logger     = zap.NewExample()
)

func rawStringToBz(raw string) []byte {
	return []byte(raw)
}

func TestHandlerMessage(t *testing.T) {
	testCases := []struct {
		name          string
		preRun        func() []byte
		msg           func() []byte
		resp          types.PriceResponse
		updateMessage func() []handlers.WebsocketEncodedMessage
		expErr        bool
	}{
		{
			name:   "invalid message",
			preRun: func() []byte { return nil },
			msg: func() []byte {
				return []byte("invalid message")
			},
			resp:          providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](nil, nil),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
		{
			name:   "invalid message type",
			preRun: func() []byte { return nil },
			msg: func() []byte {
				msg := bitfinex.BaseMessage{
					Event: "unknown",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp:          providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](nil, nil),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
		{
			name: "ticker price update",
			preRun: func() []byte {
				msg := bitfinex.SubscribedMessage{
					BaseMessage: bitfinex.BaseMessage{Event: string(bitfinex.EventSubscribed)},
					Channel:     string(bitfinex.ChannelTicker),
					ChannelID:   channelBTC,
					Pair:        "BTCUSD",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			msg: func() []byte {
				return rawStringToBz(`[111,[14957,68.17328796,14958,55.29588132,-659,-0.0422,1.0,53723.08813995,16494,14454]]`)
			},
			resp: providertypes.NewGetResponse(
				types.ResolvedPrices{
					constants.BITCOIN_USD: {
						Value: big.NewInt(100000000),
					},
				},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        false,
		},
		{
			name:   "ticker price update with unknown channel ID",
			preRun: func() []byte { return nil },
			msg: func() []byte {
				return rawStringToBz(`[0,[14957,68.17328796,14958,55.29588132,-659,-0.0422,1.0,53723.08813995,16494,14454]]`)
			},
			resp: providertypes.NewGetResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
		{
			name:   "successful subscription",
			preRun: func() []byte { return nil },
			msg: func() []byte {
				msg := bitfinex.SubscribedMessage{
					BaseMessage: bitfinex.BaseMessage{Event: string(bitfinex.EventSubscribed)},
					Channel:     string(bitfinex.ChannelTicker),
					ChannelID:   channelBTC,
					Pair:        "BTCUSD",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: providertypes.NewGetResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        false,
		},
		{
			name:   "subscription error",
			preRun: func() []byte { return nil },
			msg: func() []byte {
				msg := bitfinex.ErrorMessage{
					BaseMessage: bitfinex.BaseMessage{Event: string(bitfinex.EventError)},
					Msg:         "error subscribing",
					Code:        int64(bitfinex.ErrorSubscriptionFailed),
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: providertypes.NewGetResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wsHandler, err := bitfinex.NewWebSocketDataHandler(logger, bitfinex.DefaultMarketConfig, bitfinex.DefaultWebSocketConfig)
			require.NoError(t, err)

			if tc.preRun() != nil {
				_, _, err = wsHandler.HandleMessage(tc.preRun())
				require.NoError(t, err)
			}

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
			expectedErr: false,
		},
		{
			name: "one currency pair",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USD,
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := bitfinex.SubscribeMessage{
					BaseMessage: bitfinex.BaseMessage{Event: string(bitfinex.EventSubscribe)},
					Channel:     string(bitfinex.ChannelTicker),
					Symbol:      "BTCUSD",
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
				constants.BITCOIN_USD,
				constants.ETHEREUM_USD,
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := bitfinex.SubscribeMessage{
					BaseMessage: bitfinex.BaseMessage{Event: string(bitfinex.EventSubscribe)},
					Channel:     string(bitfinex.ChannelTicker),
					Symbol:      "BTCUSD",
				}
				bz1, err := json.Marshal(msg)
				require.NoError(t, err)

				msg = bitfinex.SubscribeMessage{
					BaseMessage: bitfinex.BaseMessage{Event: string(bitfinex.EventSubscribe)},
					Channel:     string(bitfinex.ChannelTicker),
					Symbol:      "ETHUSD",
				}
				bz2, err := json.Marshal(msg)
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{bz1, bz2}
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
			handler, err := bitfinex.NewWebSocketDataHandler(logger, bitfinex.DefaultMarketConfig, bitfinex.DefaultWebSocketConfig)
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
