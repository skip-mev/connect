package kraken_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/kraken"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	logger = zap.NewExample()

	cfg = config.ProviderConfig{
		Name: kraken.Name,
		WebSocket: config.WebSocketConfig{
			Enabled:             true,
			MaxBufferSize:       1024,
			ReconnectionTimeout: 10 * time.Second,
			WSS:                 kraken.URL,
			Name:                kraken.Name,
			ReadBufferSize:      config.DefaultReadBufferSize,
			WriteBufferSize:     config.DefaultWriteBufferSize,
			HandshakeTimeout:    config.DefaultHandshakeTimeout,
			EnableCompression:   config.DefaultEnableCompression,
			ReadTimeout:         config.DefaultReadTimeout,
			WriteTimeout:        config.DefaultWriteTimeout,
		},
		Market: config.MarketConfig{
			Name: kraken.Name,
			CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
				"BITCOIN/USD": {
					Ticker:       "XBT/USD",
					CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				},
				"ETHEREUM/USD": {
					Ticker:       "ETH/USD",
					CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
				},
			},
		},
	}
)

func TestHandleMessage(t *testing.T) {
	testCases := []struct {
		name        string
		msg         func() []byte
		resp        providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
		updateMsg   func() []handlers.WebsocketEncodedMessage
		expectedErr bool
	}{
		{
			name: "valid ticker response message",
			msg: func() []byte {
				return []byte(`[340,{"a":["42694.60000",31,"31.27308189"],"b":["42694.50000",1,"1.01355072"],"c":["42694.60000","0.00455773"],"v":["2068.49653432","2075.61202911"],"p":["42596.41907","42598.31137"],"t":[21771,22049],"l":["42190.20000","42190.20000"],"h":["43165.00000","43165.00000"],"o":["43134.70000","43159.20000"]},"ticker","XBT/USD"]`)
			},
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
				Resolved: map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): {
						Value: big.NewInt(4259641907000),
					},
				},
				UnResolved: map[oracletypes.CurrencyPair]error{},
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
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := kraken.NewWebSocketDataHandler(logger, cfg)
			require.NoError(t, err)

			resp, updateMsg, err := handler.HandleMessage(tc.msg())
			fmt.Println(err)
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
			expectedErr: true,
		},
		{
			name: "single currency pair",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
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
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := kraken.SubscribeRequestMessage{
					Event: string(kraken.SubscribeEvent),
					Pair:  []string{"XBT/USD", "ETH/USD"},
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
			name: "one known and one unknown currency pair",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				oracletypes.NewCurrencyPair("MOG", "USD"),
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := kraken.NewWebSocketDataHandler(logger, cfg)
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
			fmt.Println(actual, err)
			if tc.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}
