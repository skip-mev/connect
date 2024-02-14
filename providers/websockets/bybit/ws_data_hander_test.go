package bybit_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	"github.com/skip-mev/slinky/providers/constants"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/bybit"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	logger = zap.NewExample()
	mogusd = mmtypes.NewTicker("MOG", "USD", 8, 1)
)

func TestHandlerMessage(t *testing.T) {
	testCases := []struct {
		name      string
		msg       func() []byte
		resp      providertypes.GetResponse[mmtypes.Ticker, *big.Int]
		updateMsg func() []handlers.WebsocketEncodedMessage
		expErr    bool
	}{
		{
			name: "invalid message",
			msg: func() []byte {
				return []byte("invalid message")
			},
			resp:      providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](nil, nil),
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
			resp:      providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](nil, nil),
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
			resp: providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{
					constants.BITCOIN_USDT: {
						Value: big.NewInt(100000000),
					},
				},
				map[mmtypes.Ticker]error{},
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
			resp: providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{},
				map[mmtypes.Ticker]error{},
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
			resp: providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{},
				map[mmtypes.Ticker]error{},
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
			resp: providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{},
				map[mmtypes.Ticker]error{},
			),
			updateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wsHandler, err := bybit.NewWebSocketDataHandler(logger, bybit.DefaultMarketConfig, bybit.DefaultWebSocketConfig)
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
		expected    func() []byte
		expectedErr bool
	}{
		{
			name: "no currency pairs",
			cps:  []mmtypes.Ticker{},
			expected: func() []byte {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "one currency pair",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
			},
			expected: func() []byte {
				msg := bybit.SubscriptionRequest{
					BaseRequest: bybit.BaseRequest{
						Op: string(bybit.OperationSubscribe),
					},
					Args: []string{"tickers.BTCUSDT"},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			expectedErr: false,
		},
		{
			name: "two currency pairs",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
				constants.ETHEREUM_USDT,
			},
			expected: func() []byte {
				msg := bybit.SubscriptionRequest{
					BaseRequest: bybit.BaseRequest{
						Op: string(bybit.OperationSubscribe),
					},
					Args: []string{"tickers.BTCUSDT", "tickers.ETHUSDT"},
				}

				bz, err := json.Marshal(msg)
				require.NoError(t, err)

				return bz
			},
			expectedErr: false,
		},
		{
			name: "one currency pair not in config",
			cps: []mmtypes.Ticker{
				mogusd,
			},
			expected: func() []byte {
				return nil
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := bybit.NewWebSocketDataHandler(logger, bybit.DefaultMarketConfig, bybit.DefaultWebSocketConfig)
			require.NoError(t, err)

			msgs, err := handler.CreateMessages(tc.cps)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, 1, len(msgs))
			require.EqualValues(t, tc.expected(), []byte(msgs[0]))
		})
	}
}
