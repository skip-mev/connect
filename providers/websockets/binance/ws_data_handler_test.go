package binance_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	"github.com/skip-mev/connect/v2/providers/websockets/binance"
)

var (
	logger  = zap.NewExample()
	btcusdt = types.NewProviderTicker("BTCUSDT", "")
	ethusdt = types.NewProviderTicker("ETHUSDT", "")
	mogusdt = types.NewProviderTicker("MOGUSDT", "")
)

func TestHandleMessage(t *testing.T) {
	cases := []struct {
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
			name: "subscription to unknown instruments",
			msg: func() []byte {
				msg := binance.SubscribeMessageResponse{
					Result: nil,
					ID:     1,
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
			name: "subscription to known instruments success",
			msg: func() []byte {
				msg := binance.SubscribeMessageResponse{
					Result: nil,
					ID:     2,
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)
				return bz
			},
			resp:          types.NewPriceResponse(nil, nil),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        false,
		},
		{
			name: "subscription to known instruments failure",
			msg: func() []byte {
				msg := binance.SubscribeMessageResponse{
					Result: "error",
					ID:     3,
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)
				return bz
			},
			resp: types.NewPriceResponse(nil, nil),
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				msg := binance.SubscribeMessageRequest{
					Method: string(binance.SubscribeMethod),
					Params: []string{
						"ethusdt@aggTrade",
						"ethusdt@ticker",
					},
					ID: 3,
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return []handlers.WebsocketEncodedMessage{bz}
			},
			expErr: false,
		},
		{
			name: "ticker stream message with good price",
			msg: func() []byte {
				msg := `
				{
					"stream": "btcusdt@ticker",
					"data": {
						"s": "btcusdt",
						"c": "10000.00000000",
						"C": 1600000000000
						}
				}`

				return []byte(msg)
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{
					btcusdt: {
						Value: big.NewFloat(10000.0),
					},
				},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        false,
		},
		{
			name: "ticker stream message with bad price",
			msg: func() []byte {
				msg := `
				{
					"stream": "btcusdt@ticker",
					"data": {
						"s": "btcusdt",
						"c": "bad_price",
						"C": 1600000000000
						}
				}`

				return []byte(msg)
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btcusdt: {
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("failed to parse price: strconv.ParseFloat: parsing \"bad_price\": invalid syntax"), providertypes.ErrorFailedToParsePrice),
					},
				},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
		{
			name: "ticker stream message with unknown instrument",
			msg: func() []byte {
				msg := `
				{
					"stream": "unknown@ticker",
					"data": {
						"s": "mogmeusdt",
						"c": "10000.00000000",
						"C": 1600000000000
						}
				}`

				return []byte(msg)
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
		{
			name: "agg trade stream message with good price",
			msg: func() []byte {
				msg := `
				{
					"stream": "btcusdt@aggTrade",
					"data": {
						"s": "btcusdt",
						"p": "10000.00000000"
						}
				}`

				return []byte(msg)
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{
					btcusdt: {
						Value: big.NewFloat(10000.0),
					},
				},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        false,
		},
		{
			name: "agg trade stream message with bad price",
			msg: func() []byte {
				msg := `
				{
					"stream": "btcusdt@aggTrade",
					"data": {
						"s": "btcusdt",
						"p": "bad_price"
						}
				}`

				return []byte(msg)
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btcusdt: {
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("failed to parse price: strconv.ParseFloat: parsing \"bad_price\": invalid syntax"), providertypes.ErrorFailedToParsePrice),
					},
				},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
		{
			name: "agg trade stream message with unknown instrument",
			msg: func() []byte {
				msg := `
				{
					"stream": "unknown@aggTrade",
					"data": {
						"s": "mogmeusdt",
						"p": "10000.00000000"
						}
				}`

				return []byte(msg)
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{},
			),
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wsHandlerI, err := binance.NewWebSocketDataHandler(logger, binance.DefaultWebSocketConfig)
			require.NoError(t, err)

			// Set the instruments for the handler.
			wsHandler := wsHandlerI.(*binance.WebSocketHandler)
			wsHandler.SetIDForInstruments(2, []string{btcusdt.GetOffChainTicker()})
			wsHandler.SetIDForInstruments(3, []string{ethusdt.GetOffChainTicker()})

			_, err = wsHandler.CreateMessages([]types.ProviderTicker{btcusdt, ethusdt})
			require.NoError(t, err)

			resp, updateMsgs, err := wsHandler.HandleMessage(tc.msg())
			if tc.expErr {
				require.Error(t, err)
				require.Equal(t, tc.updateMessage(), updateMsgs)
			} else {
				require.NoError(t, err)
			}

			require.Len(t, updateMsgs, len(tc.updateMessage()))
			seenIDs := make(map[int64]struct{})
			for i, m := range updateMsgs {
				var msg binance.SubscribeMessageRequest
				require.NoError(t, json.Unmarshal(m, &msg))

				var expected binance.SubscribeMessageRequest
				require.NoError(t, json.Unmarshal(tc.updateMessage()[i], &expected))

				require.Equal(t, expected.Method, msg.Method)
				require.Equal(t, expected.Params, msg.Params)
				require.NotZero(t, msg.ID)
				require.NotContains(t, seenIDs, msg.ID)
				seenIDs[msg.ID] = struct{}{}
			}

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
	batchCfg := binance.DefaultWebSocketConfig
	batchCfg.MaxSubscriptionsPerBatch = 2

	cases := []struct {
		name        string
		ticker      []types.ProviderTicker
		cfg         config.WebSocketConfig
		expected    func() []binance.SubscribeMessageRequest
		expectedErr bool
	}{
		{
			name:   "no tickers",
			ticker: []types.ProviderTicker{},
			cfg:    binance.DefaultWebSocketConfig,
			expected: func() []binance.SubscribeMessageRequest {
				return []binance.SubscribeMessageRequest{}
			},
			expectedErr: true,
		},
		{
			name: "single ticker",
			ticker: []types.ProviderTicker{
				btcusdt,
			},
			cfg: binance.DefaultWebSocketConfig,
			expected: func() []binance.SubscribeMessageRequest {
				return []binance.SubscribeMessageRequest{
					{
						Method: string(binance.SubscribeMethod),
						Params: []string{
							"btcusdt@aggTrade",
							"btcusdt@ticker",
						},
						ID: 1,
					},
				}
			},
			expectedErr: false,
		},
		{
			name: "multiple tickers",
			ticker: []types.ProviderTicker{
				btcusdt,
				ethusdt,
			},
			cfg: binance.DefaultWebSocketConfig,
			expected: func() []binance.SubscribeMessageRequest {
				return []binance.SubscribeMessageRequest{
					{
						Method: string(binance.SubscribeMethod),
						Params: []string{
							"btcusdt@aggTrade",
							"btcusdt@ticker",
						},
						ID: 1,
					},
					{
						Method: string(binance.SubscribeMethod),
						Params: []string{
							"ethusdt@aggTrade",
							"ethusdt@ticker",
						},
						ID: 1,
					},
				}
			},
			expectedErr: false,
		},
		{
			name: "multiple tickers with batch config",
			ticker: []types.ProviderTicker{
				btcusdt,
				ethusdt,
			},
			cfg: batchCfg,
			expected: func() []binance.SubscribeMessageRequest {
				return []binance.SubscribeMessageRequest{
					{
						Method: string(binance.SubscribeMethod),
						Params: []string{
							"btcusdt@aggTrade",
							"btcusdt@ticker",
							"ethusdt@aggTrade",
							"ethusdt@ticker",
						},
						ID: 1,
					},
				}
			},
			expectedErr: false,
		},
		{
			name: "multiple tickers with batch config and multiple batches",
			ticker: []types.ProviderTicker{
				btcusdt,
				ethusdt,
				mogusdt,
			},
			cfg: batchCfg,
			expected: func() []binance.SubscribeMessageRequest {
				return []binance.SubscribeMessageRequest{
					{
						Method: string(binance.SubscribeMethod),
						Params: []string{
							"btcusdt@aggTrade",
							"btcusdt@ticker",
							"ethusdt@aggTrade",
							"ethusdt@ticker",
						},
						ID: 1,
					},
					{
						Method: string(binance.SubscribeMethod),
						Params: []string{
							"mogusdt@aggTrade",
							"mogusdt@ticker",
						},
						ID: 2,
					},
				}
			},
			expectedErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := binance.NewWebSocketDataHandler(logger, tc.cfg)
			require.NoError(t, err)

			actual, err := handler.CreateMessages(tc.ticker)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			expected := tc.expected()
			require.Equal(t, len(expected), len(actual))

			seenIDs := make(map[int64]struct{})
			for i, m := range actual {
				var msg binance.SubscribeMessageRequest
				require.NoError(t, json.Unmarshal(m, &msg))

				require.Equal(t, expected[i].Method, msg.Method)
				require.Equal(t, expected[i].Params, msg.Params)
				require.NotZero(t, msg.ID)
				require.NotContains(t, seenIDs, msg.ID)
				seenIDs[msg.ID] = struct{}{}
			}
		})
	}
}
