package bybit_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	"github.com/skip-mev/slinky/providers/websockets/bybit"
)

var (
	btcusdt = types.DefaultProviderTicker{
		OffChainTicker: "BTCUSDT",
	}
	ethusdt = types.DefaultProviderTicker{
		OffChainTicker: "ETHUSDT",
	}
	logger  = zap.NewExample()
)

func TestHandlerMessage(t *testing.T) {
	testCases := []struct {
		name      string
		msg       func() []byte
		resp      types.PriceResponse
		updateMsg func() []handlers.WebsocketEncodedMessage
		expErr    bool
	}{
		{
			name: "invalid message",
			msg: func() []byte {
				return []byte("invalid message")
			},
			resp:      types.NewPriceResponse(nil, nil),
			updateMsg: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:    true,
		},
		{
			name: "invalid message type",
			msg: func() []byte {
				msg := bybit.BaseResponse{
					Op: "unknown",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp:      types.NewPriceResponse(nil, nil),
			updateMsg: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:    true,
		},
		{
			name: "price update",
			msg: func() []byte {
				msg := bybit.TickerUpdateMessage{
					Topic: "tickers.BTCUSDT",
					Data: bybit.TickerUpdateData{
						Symbol:    "BTCUSDT",
						LastPrice: "1",
					},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{
					btcusdt: {
						Value: big.NewFloat(1.0),
					},
				},
				types.UnResolvedPrices{},
			),
			updateMsg: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:    false,
		},
		{
			name: "price update with unknown pair ID",
			msg: func() []byte {
				msg := bybit.TickerUpdateMessage{
					Topic: "tickers.MOGUSDT",
					Data: bybit.TickerUpdateData{
						Symbol:    "MOGUSDT",
						LastPrice: "1",
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
			updateMsg: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:    true,
		},
		{
			name: "successful subscription",
			msg: func() []byte {
				msg := bybit.SubscriptionResponse{
					BaseResponse: bybit.BaseResponse{
						Success: true,
						RetMsg:  string(bybit.OperationSubscribe),
						ConnID:  "90190u1309",
						Op:      string(bybit.OperationSubscribe),
					},
					ReqID: "1012901",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{},
			),
			updateMsg: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:    false,
		},
		{
			name: "subscription error",
			msg: func() []byte {
				msg := bybit.SubscriptionResponse{
					BaseResponse: bybit.BaseResponse{
						Success: false,
						RetMsg:  string(bybit.OperationSubscribe),
						ConnID:  "90190u1309",
						Op:      string(bybit.OperationSubscribe),
					},
					ReqID: "1012901",
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			resp: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{},
			),
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wsHandler, err := bybit.NewWebSocketDataHandler(logger, bybit.DefaultWebSocketConfig)
			require.NoError(t, err)

			// Update the cache since it is assumed that CreateMessages is executed before anything else.
			_, err = wsHandler.CreateMessages([]types.ProviderTicker{btcusdt, ethusdt})
			require.NoError(t, err)

			resp, updateMsg, err := wsHandler.HandleMessage(tc.msg())
			if tc.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tc.updateMsg(), updateMsg)

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
		expected    func() [][]byte
		expectedErr bool
	}{
		{
			name: "no currency pairs",
			cps:  []types.ProviderTicker{},
			expected: func() [][]byte {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "one currency pair",
			cps: []types.ProviderTicker{
				btcusdt,
			},
			expected: func() [][]byte {
				msg := bybit.SubscriptionRequest{
					BaseRequest: bybit.BaseRequest{
						Op: string(bybit.OperationSubscribe),
					},
					Args: []string{"tickers.BTCUSDT"},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return [][]byte{bz}
			},
			expectedErr: false,
		},
		{
			name: "two currency pairs",
			cps: []types.ProviderTicker{
				btcusdt,
				ethusdt,
			},
			expected: func() [][]byte {
				msgs := make([][]byte, 2)
				for i, ticker := range []string{"tickers.BTCUSDT", "tickers.ETHUSDT"} {
					msg := bybit.SubscriptionRequest{
						BaseRequest: bybit.BaseRequest{
							Op: string(bybit.OperationSubscribe),
						},
						Args: []string{ticker},
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
			wsHandler, err := bybit.NewWebSocketDataHandler(logger, bybit.DefaultWebSocketConfig)
			require.NoError(t, err)

			msgs, err := wsHandler.CreateMessages(tc.cps)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			expected := tc.expected()
			require.Equal(t, len(expected), len(msgs))
			for i, msg := range msgs {
				require.EqualValues(t, expected[i], msg)
			}
		})
	}
}
