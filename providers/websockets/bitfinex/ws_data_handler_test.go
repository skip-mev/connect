package bitfinex_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/bitfinex"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	providerCfg = config.ProviderConfig{
		Name:      bitfinex.Name,
		WebSocket: bitfinex.DefaultWebSocketConfig,
		Market: config.MarketConfig{
			Name: bitfinex.Name,
			CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
				"BITCOIN/USDT": {
					Ticker:       "BTC-USDT",
					CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
				},
				"ETHEREUM/USDT": {
					Ticker:       "ETH-USDT",
					CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
				},
			},
		},
	}

	channelMap = map[string]config.CurrencyPairMarketConfig{
		"btc-channel": providerCfg.Market.CurrencyPairToMarketConfigs["BITCOIN/USDT"],
		"eth-channel": providerCfg.Market.CurrencyPairToMarketConfigs["ETHEREUM/USDT"],
	}

	logger = zap.NewExample()
)

func TestHandlerMessage(t *testing.T) {
	testCases := []struct {
		name          string
		preRun        func() []byte
		msg           func() []byte
		resp          providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
		updateMessage func() []handlers.WebsocketEncodedMessage
		expErr        bool
	}{
		{
			name:   "invalid message",
			preRun: func() []byte { return nil },
			msg: func() []byte {
				return []byte("invalid message")
			},
			resp:          providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](nil, nil),
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
			resp:          providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](nil, nil),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
		{
			name: "ticker price update",
			preRun: func() []byte {
				msg := bitfinex.SubscribedMessage{
					BaseMessage: bitfinex.BaseMessage{Event: string(bitfinex.EventSubscribed)},
					Channel:     string(bitfinex.ChannelTicker),
					ChannelID:   "btc-channel",
					Pair:        "BTC-USDT",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			msg: func() []byte {
				msg := bitfinex.TickerStream{
					BaseStreamMessage: bitfinex.BaseStreamMessage{ChannelID: "btc-channel"},
					LastPrice:         "1",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					oracletypes.NewCurrencyPair("BITCOIN", "USDT"): {
						Value: big.NewInt(100000000),
					},
				},
				map[oracletypes.CurrencyPair]error{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        false,
		},
		{
			name:   "ticker price update with unknown channel ID",
			preRun: func() []byte { return nil },
			msg: func() []byte {
				msg := bitfinex.TickerStream{
					BaseStreamMessage: bitfinex.BaseStreamMessage{ChannelID: "unknown"},
					LastPrice:         "",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        false,
		},
		{
			name:   "successful subscription",
			preRun: func() []byte { return nil },
			msg: func() []byte {
				msg := bitfinex.SubscribedMessage{
					BaseMessage: bitfinex.BaseMessage{Event: string(bitfinex.EventSubscribed)},
					Channel:     string(bitfinex.ChannelTicker),
					ChannelID:   "btc-channel",
					Pair:        "BTC-USDT",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{},
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
			resp: providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wsHandler, err := bitfinex.NewWebSocketDataHandler(logger, providerCfg)
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
			expectedErr: false,
		},
		{
			name: "one currency pair",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := bitfinex.SubscribeMessage{
					BaseMessage: bitfinex.BaseMessage{Event: string(bitfinex.EventSubscribe)},
					Channel:     string(bitfinex.ChannelTicker),
					Symbol:      "BTC-USDT",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{bz}
			},
			expectedErr: false,
		},
		{
			name: "two currency pairs",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
				oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := bitfinex.SubscribeMessage{
					BaseMessage: bitfinex.BaseMessage{Event: string(bitfinex.EventSubscribe)},
					Channel:     string(bitfinex.ChannelTicker),
					Symbol:      "BTC-USDT",
				}
				bz1, err := json.Marshal(msg)
				require.NoError(t, err)

				msg = bitfinex.SubscribeMessage{
					BaseMessage: bitfinex.BaseMessage{Event: string(bitfinex.EventSubscribe)},
					Channel:     string(bitfinex.ChannelTicker),
					Symbol:      "ETH-USDT",
				}
				bz2, err := json.Marshal(msg)
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{bz1, bz2}
			},
			expectedErr: false,
		},
		{
			name: "one currency pair not in config",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("MOG", "USDT"),
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := bitfinex.NewWebSocketDataHandler(logger, providerCfg)
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
