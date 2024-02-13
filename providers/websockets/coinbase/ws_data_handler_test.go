package coinbase_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/coinbase"
)

var (
	providerCfg = config.ProviderConfig{
		Name:      coinbase.Name,
		WebSocket: coinbase.DefaultWebSocketConfig,
		Market: config.MarketConfig{
			Name: coinbase.Name,
			CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
				"BITCOIN/USD": {
					Ticker:       "BTC-USD",
					CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
				},
				"ETHEREUM/USD": {
					Ticker:       "ETH-USD",
					CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
				},
			},
		},
	}

	logger = zap.NewExample()
)

func TestHandleMessage(t *testing.T) {
	testCases := []struct {
		name          string
		msg           func() []byte
		resp          providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]
		updateMessage func() []handlers.WebsocketEncodedMessage
		expErr        bool
	}{
		{
			name: "unknown message",
			msg: func() []byte {
				return []byte(`{"type":"unknown"}`)
			},
			resp:          providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{},
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
		{
			name: "ticker message",
			msg: func() []byte {
				msg := coinbase.TickerResponseMessage{
					Type:     string(coinbase.TickerMessage),
					Ticker:   "BTC-USD",
					Price:    "10000.00",
					Sequence: 1,
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				Resolved: map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					slinkytypes.NewCurrencyPair("BITCOIN", "USD"): {
						Value: big.NewInt(1000000000000),
					},
				},
			},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: false,
		},
		{
			name: "ticker message with invalid ticker",
			msg: func() []byte {
				msg := coinbase.TickerResponseMessage{
					Type:     string(coinbase.TickerMessage),
					Ticker:   "MOG-USD",
					Price:    "10000.00",
					Sequence: 1,
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: true,
		},
		{
			name: "ticker message with bad price",
			msg: func() []byte {
				msg := coinbase.TickerResponseMessage{
					Type:     string(coinbase.TickerMessage),
					Ticker:   "BTC-USD",
					Price:    "$10000.00.00",
					Sequence: 1,
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				UnResolved: map[slinkytypes.CurrencyPair]error{
					slinkytypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("failed to convert price to big int"),
				},
			},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: true,
		},
		{
			name: "ticker message with out of order sequence number",
			msg: func() []byte {
				msg := coinbase.TickerResponseMessage{
					Type:     string(coinbase.TickerMessage),
					Ticker:   "BTC-USD", // We have already received a message with sequence number 1.
					Price:    "10000.00",
					Sequence: 1,
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				UnResolved: map[slinkytypes.CurrencyPair]error{
					slinkytypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("received out of order ticker response message"),
				},
			},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: true,
		},
		{
			name: "subscriptions message",
			msg: func() []byte {
				msg := coinbase.SubscribeResponseMessage{
					Type: string(coinbase.SubscriptionsMessage),
					Channels: []coinbase.Channel{
						{
							Name: string(coinbase.TickerChannel),
							Instruments: []string{
								"BTC-USD",
							},
						},
					},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: false,
		},
	}

	wsHandler, err := coinbase.NewWebSocketDataHandler(logger, providerCfg)
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, updateMsg, err := wsHandler.HandleMessage(tc.msg())
			if tc.expErr {
				require.Error(t, err)

				require.Equal(t, len(tc.resp.UnResolved), len(resp.UnResolved))
				for cp := range tc.resp.UnResolved {
					require.Contains(t, resp.UnResolved, cp)
					require.Error(t, resp.UnResolved[cp])
				}
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

func TestCreateMessages(t *testing.T) {
	testCases := []struct {
		name        string
		cps         []slinkytypes.CurrencyPair
		expected    func() []handlers.WebsocketEncodedMessage
		expectedErr bool
	}{
		{
			name: "no currency pairs to subscribe to",
			cps:  []slinkytypes.CurrencyPair{},
			expected: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "one currency pair to subscribe to",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := coinbase.SubscribeRequestMessage{
					Type: string(coinbase.SubscribeMessage),
					ProductIDs: []string{
						"BTC-USD",
					},
					Channels: []string{
						string(coinbase.TickerChannel),
					},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{bz}
			},
			expectedErr: false,
		},
		{
			name: "multiple currency pairs to subscribe to",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
				slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := coinbase.SubscribeRequestMessage{
					Type: string(coinbase.SubscribeMessage),
					ProductIDs: []string{
						"BTC-USD",
						"ETH-USD",
					},
					Channels: []string{
						string(coinbase.TickerChannel),
					},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{bz}
			},
			expectedErr: false,
		},
		{
			name: "multiple currency pairs to subscribe to with one not supported",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
				slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
				slinkytypes.NewCurrencyPair("MOG", "USD"),
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := coinbase.SubscribeRequestMessage{
					Type: string(coinbase.SubscribeMessage),
					ProductIDs: []string{
						"BTC-USD",
						"ETH-USD",
					},
					Channels: []string{
						string(coinbase.TickerChannel),
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
			handler, err := coinbase.NewWebSocketDataHandler(logger, providerCfg)
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
