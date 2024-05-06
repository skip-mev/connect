package coinbase_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	providertypes "github.com/skip-mev/slinky/providers/types"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	"github.com/skip-mev/slinky/providers/websockets/coinbase"
)

var (
	btcusd = coinbase.DefaultMarketConfig.MustGetProviderTicker(constants.BITCOIN_USD)
	ethusd = coinbase.DefaultMarketConfig.MustGetProviderTicker(constants.ETHEREUM_USD)
	logger = zap.NewExample()
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
			name: "unknown message",
			msg: func() []byte {
				return []byte(`{"type":"unknown"}`)
			},
			resp:          types.PriceResponse{},
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
			resp: types.PriceResponse{
				Resolved: types.ResolvedPrices{
					btcusd: {
						Value: big.NewFloat(10000.00),
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
			resp: types.PriceResponse{},
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
			resp: types.PriceResponse{
				UnResolved: types.UnResolvedPrices{
					btcusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("failed to convert price to big int"), providertypes.ErrorUnknown),
					},
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
			resp: types.PriceResponse{
				UnResolved: types.UnResolvedPrices{
					btcusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("received out of order ticker response message"), providertypes.ErrorUnknown),
					},
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
			resp: types.PriceResponse{},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: false,
		},
	}

	wsHandler, err := coinbase.NewWebSocketDataHandler(logger, coinbase.DefaultWebSocketConfig)
	require.NoError(t, err)

	// Update the cache since it is assumed that CreateMessages is executed before anything else.
	_, err = wsHandler.CreateMessages([]types.ProviderTicker{btcusd, ethusd})
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
				require.Equal(t, result.Value.SetPrec(18), resp.Resolved[cp].Value.SetPrec(18))
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
		cps         []types.ProviderTicker
		expected    func() []handlers.WebsocketEncodedMessage
		expectedErr bool
	}{
		{
			name: "no currency pairs to subscribe to",
			cps:  []types.ProviderTicker{},
			expected: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "one currency pair to subscribe to",
			cps: []types.ProviderTicker{
				btcusd,
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
			cps: []types.ProviderTicker{
				btcusd,
				ethusd,
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msgs := make([]handlers.WebsocketEncodedMessage, 2)
				for i, ticker := range []string{"BTC-USD", "ETH-USD"} {
					msg := coinbase.SubscribeRequestMessage{
						Type: string(coinbase.SubscribeMessage),
						ProductIDs: []string{
							ticker,
						},
						Channels: []string{
							string(coinbase.TickerChannel),
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
			wsHandler, err := coinbase.NewWebSocketDataHandler(logger, coinbase.DefaultWebSocketConfig)
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
