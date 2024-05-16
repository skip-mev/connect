package kraken_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	"github.com/skip-mev/slinky/providers/websockets/kraken"
)

var (
	btcusd = types.DefaultProviderTicker{
		OffChainTicker: "BTCUSD",
	}
	ethusd = types.DefaultProviderTicker{
		OffChainTicker: "ETHUSD",
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
			name: "valid ticker response message",
			msg: func() []byte {
				return []byte(`[340,{"a":["42694.60000",31,"31.27308189"],"b":["42694.50000",1,"1.01355072"],"c":["42694.60000","0.00455773"],"v":["2068.49653432","2075.61202911"],"p":["42596.41907","42598.31137"],"t":[21771,22049],"l":["42190.20000","42190.20000"],"h":["43165.00000","43165.00000"],"o":["43134.70000","43159.20000"]},"ticker","XBT/USD"]`)
			},
			resp: types.PriceResponse{
				Resolved: types.ResolvedPrices{
					btcusd: {
						Value: big.NewFloat(42596.41907000),
					},
				},
				UnResolved: types.UnResolvedPrices{},
			},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: false,
		},
		{
			name: "invalid ticker response message",
			msg: func() []byte {
				return []byte(`[340,{"a":["42694.60000",31,"31.27308189"],"b":["42694.50000",1,"1.01355072"],"c":["42694.60000","0.00455773"],"v":["2068.49653432","2075.61202911"],"p":["42596.41907","42598.31137"],"t":[21771,22049],"l":["42190.20000","42190.20000"],"h":["43165.00000","43165.00000"],"o":["43134.70000","43159.20000"]},"ticker"]`)
			},
			resp: types.PriceResponse{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "invalid response channel for ticker response message",
			msg: func() []byte {
				return []byte(`[340,{"a":["42694.60000",31,"31.27308189"],"b":["42694.50000",1,"1.01355072"],"c":["42694.60000","0.00455773"],"v":["2068.49653432","2075.61202911"],"p":["42596.41907","42598.31137"],"t":[21771,22049],"l":["42190.20000","42190.20000"],"h":["43165.00000","43165.00000"],"o":["43134.70000","43159.20000"]},"book", "XBT/USD"]`)
			},
			resp: types.PriceResponse{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "received ticker response message with unknown currency pair",
			msg: func() []byte {
				return []byte(`[340,{"a":["42694.60000",31,"31.27308189"],"b":["42694.50000",1,"1.01355072"],"c":["42694.60000","0.00455773"],"v":["2068.49653432","2075.61202911"],"p":["42596.41907","42598.31137"],"t":[21771,22049],"l":["42190.20000","42190.20000"],"h":["43165.00000","43165.00000"],"o":["43134.70000","43159.20000"]},"ticker","MOG/USD"]`)
			},
			resp: types.PriceResponse{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "missing price update in ticker response message",
			msg: func() []byte {
				return []byte(`[340,{"a":["42694.60000",31,"31.27308189"],"b":["42694.50000",1,"1.01355072"],"v":["2068.49653432","2075.61202911"],"p":["42596.41907"],"t":[21771,22049],"l":["42190.20000","42190.20000"],"h":["43165.00000","43165.00000"],"o":["43134.70000","43159.20000"]},"ticker","XBT/USD"]`)
			},
			resp: types.PriceResponse{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "cannot convert price update in ticker response message to big int",
			msg: func() []byte {
				return []byte(`[340,{"a":["42694.60000",31,"31.27308189"],"b":["42694.50000",1,"1.01355072"],"v":["2068.49653432","2075.61202911"],"p":["$42,596.41907","42598.31137"],"t":[21771,22049],"l":["42190.20000","42190.20000"],"h":["43165.00000","43165.00000"],"o":["43134.70000","43159.20000"]},"ticker","XBT/USD"]`)
			},
			resp: types.PriceResponse{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "system status response message (online)",
			msg: func() []byte {
				return []byte(`{"connectionID": 1234, "event": "systemStatus", "status": "online", "version": "1.0.0"}`)
			},
			resp: types.PriceResponse{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: false,
		},
		{
			name: "system status response message (offline)",
			msg: func() []byte {
				return []byte(`{"connectionID": 1234, "event": "systemStatus", "status": "offline", "version": "1.0.0"}`)
			},
			resp: types.PriceResponse{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "heartbeat response message",
			msg: func() []byte {
				return []byte(`{"connectionID": 1234, "event": "heartbeat"}`)
			},
			resp: types.PriceResponse{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: false,
		},
		{
			name: "subscription status response message (subscribed)",
			msg: func() []byte {
				return []byte(`{"channelID": 1234, "event": "subscriptionStatus", "pair": "XBT/USD", "status": "subscribed", "subscription": {"name": "ticker"}, "channelName": "ticker"}`)
			},
			resp: types.PriceResponse{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: false,
		},
		{
			name: "subscription status response message (error)",
			msg: func() []byte {
				return []byte(`{"errorMessage": "Subscription depth not supported", "event": "subscriptionStatus", "pair": "XBT/USD", "status": "error", "subscription": {"name": "ticker"}, "channelName": "ticker"}`)
			},
			resp: types.PriceResponse{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				msg := kraken.SubscribeRequestMessage{
					Event: string(kraken.SubscribeEvent),
					Pair:  []string{"XBT/USD"},
					Subscription: kraken.Subscription{
						Name: string(kraken.TickerChannel),
					},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{bz}
			},
			expectedErr: false,
		},
		{
			name: "unknown subscription status response message",
			msg: func() []byte {
				return []byte(`{"channelID": 1234, "event": "subscriptionStatus", "pair": "XBT/USD", "status": "unknown", "subscription": {"name": "ticker"}, "channelName": "ticker"}`)
			},
			resp: types.PriceResponse{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "unknown event type",
			msg: func() []byte {
				return []byte(`{"event": "unknown"}`)
			},
			resp: types.PriceResponse{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := kraken.NewWebSocketDataHandler(logger, kraken.DefaultWebSocketConfig)
			require.NoError(t, err)

			// Update the cache since it is assumed that CreateMessages is executed before anything else.
			_, err = handler.CreateMessages([]types.ProviderTicker{btcusd, ethusd})
			require.NoError(t, err)

			resp, updateMsg, err := handler.HandleMessage(tc.msg())
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// The response should contain a single resolved price update.
			require.LessOrEqual(t, len(resp.Resolved), 1)
			require.LessOrEqual(t, len(resp.UnResolved), 1)

			require.Equal(t, tc.updateMsg(), updateMsg)

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
			name: "single currency pair",
			cps: []types.ProviderTicker{
				btcusd,
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := kraken.SubscribeRequestMessage{
					Event: string(kraken.SubscribeEvent),
					Pair:  []string{"XBT/USD"},
					Subscription: kraken.Subscription{
						Name: string(kraken.TickerChannel),
					},
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
				btcusd,
				ethusd,
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msgs := make([]handlers.WebsocketEncodedMessage, 2)
				for i, ticker := range []string{"XBT/USD", "ETH/USD"} {
					msg := kraken.SubscribeRequestMessage{
						Event: string(kraken.SubscribeEvent),
						Pair:  []string{ticker},
						Subscription: kraken.Subscription{
							Name: string(kraken.TickerChannel),
						},
					}

					bz, err := json.Marshal(msg)
					require.NoError(t, err)
					msgs[i] = bz
				}

				return msgs
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := kraken.NewWebSocketDataHandler(logger, kraken.DefaultWebSocketConfig)
			require.NoError(t, err)

			actual, err := handler.CreateMessages(tc.cps)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected(), actual)
		})
	}
}

func TestDecodeTickerResponseMessage(t *testing.T) {
	testCases := []struct {
		name     string
		response string
		expected kraken.TickerResponseMessage
		expErr   bool
	}{
		{
			name:     "valid response",
			response: `[340,{"a":["42694.60000",31,"31.27308189"],"b":["42694.50000",1,"1.01355072"],"c":["42694.60000","0.00455773"],"v":["2068.49653432","2075.61202911"],"p":["42596.41907","42598.31137"],"t":[21771,22049],"l":["42190.20000","42190.20000"],"h":["43165.00000","43165.00000"],"o":["43134.70000","43159.20000"]},"ticker","XBT/USD"]`,
			expected: kraken.TickerResponseMessage{
				ChannelID: 340,
				TickerData: kraken.TickerData{
					VolumeWeightedAveragePrice: []string{"42596.41907", "42598.31137"},
				},
				ChannelName: "ticker",
				Pair:        "XBT/USD",
			},
		},
		{
			name:     "invalid response with missing channel ID",
			response: `[{"a":["42694.60000",31,"31.27308189"],"b":["42694.50000",1,"1.01355072"],"c":["42694.60000","0.00455773"],"v":["2068.49653432","2075.61202911"],"p":["42596.41907","42598.31137"],"t":[21771,22049],"l":["42190.20000","42190.20000"],"h":["43165.00000","43165.00000"],"o":["43134.70000","43159.20000"]},"ticker","XBT/USD"]`,
			expected: kraken.TickerResponseMessage{},
			expErr:   true,
		},
		{
			name:     "invalid response with missing ticker data",
			response: `[340,"ticker","XBT/USD"]`,
			expected: kraken.TickerResponseMessage{},
			expErr:   true,
		},
		{
			name:     "invalid response with missing channel name",
			response: `[340,{"a":["42694.60000",31,"31.27308189"],"b":["42694.50000",1,"1.01355072"],"c":["42694.60000","0.00455773"],"v":["2068.49653432","2075.61202911"],"p":["42596.41907","42598.31137"],"t":[21771,22049],"l":["42190.20000","42190.20000"],"h":["43165.00000","43165.00000"],"o":["43134.70000","43159.20000"]},"XBT/USD"]`,
			expected: kraken.TickerResponseMessage{},
			expErr:   true,
		},
		{
			name:     "invalid response with missing pair",
			response: `[340,{"a":["42694.60000",31,"31.27308189"],"b":["42694.50000",1,"1.01355072"],"c":["42694.60000","0.00455773"],"v":["2068.49653432","2075.61202911"],"p":["42596.41907","42598.31137"],"t":[21771,22049],"l":["42190.20000","42190.20000"],"h":["43165.00000","43165.00000"],"o":["43134.70000","43159.20000"]},"ticker"]`,
			expected: kraken.TickerResponseMessage{},
			expErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := kraken.DecodeTickerResponseMessage([]byte(tc.response))
			if tc.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}
