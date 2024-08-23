package mexc_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	"github.com/skip-mev/connect/v2/providers/websockets/mexc"
)

var (
	btcusdt = types.DefaultProviderTicker{
		OffChainTicker: "BTCUSDT",
	}
	ethusdt = types.DefaultProviderTicker{
		OffChainTicker: "ETHUSDT",
	}
	atomusdc = types.DefaultProviderTicker{
		OffChainTicker: "ATOMUSDC",
	}
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
			name: "pong message",
			msg: func() []byte {
				return []byte(`{"id":0,"code":0,"msg":"PONG"}`)
			},
			resp: types.PriceResponse{},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: false,
		},
		{
			name: "subscription message",
			msg: func() []byte {
				return []byte(`{"id":0,"code":0,"msg":"spot@public.miniTicker.v3.api@BTCUSDT@UTC+8"}`)
			},
			resp: types.PriceResponse{},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: false,
		},
		{
			name: "unknown message",
			msg: func() []byte {
				return []byte(`{"id":0,"code":0,"msg":"UNKNOWN"}`)
			},
			resp: types.PriceResponse{},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: true,
		},
		{
			name: "price update message",
			msg: func() []byte {
				msg := `{"c":"spot@public.miniTicker.v3.api@BTCUSDT@UTC+8","d":{"s":"BTCUSDT","p":"10000.00"}}`
				return []byte(msg)
			},
			resp: types.PriceResponse{
				Resolved: types.ResolvedPrices{
					btcusdt: {
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
			name: "unsupported market price update",
			msg: func() []byte {
				msg := `{"c":"spot@public.miniTicker.v3.api@MOGUSDT@UTC+8","d":{"s":"MOGUSDT","p":"10000.00"}}`
				return []byte(msg)
			},
			resp: types.PriceResponse{},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: true,
		},
		{
			name: "price update from incorrect channel",
			msg: func() []byte {
				msg := `{"c":"futures@public.miniTicker.v3.api@BTCUSDT@UTC+8","d":{"s":"BTCUSDT","p":"10000.00"}}`
				return []byte(msg)
			},
			resp: types.PriceResponse{
				UnResolved: types.UnResolvedPrices{
					btcusdt: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("invalid channel"), providertypes.ErrorWebSocketGeneral),
					},
				},
			},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: true,
		},
		{
			name: "price update with invalid price",
			msg: func() []byte {
				msg := `{"c":"spot@public.miniTicker.v3.api@BTCUSDT@UTC+8","d":{"s":"BTCUSDT","p":"$10,000.00"}}`
				return []byte(msg)
			},
			resp: types.PriceResponse{
				UnResolved: types.UnResolvedPrices{
					btcusdt: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("invalid price"), providertypes.ErrorWebSocketGeneral),
					},
				},
			},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wsHandler, err := mexc.NewWebSocketDataHandler(logger, mexc.DefaultWebSocketConfig)
			require.NoError(t, err)

			// Update the cache since it is assumed that CreateMessages is executed before anything else.
			_, err = wsHandler.CreateMessages([]types.ProviderTicker{btcusdt, ethusdt, atomusdc})
			require.NoError(t, err)

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
	batchCfg := mexc.DefaultWebSocketConfig
	batchCfg.MaxSubscriptionsPerBatch = 2

	testCases := []struct {
		name        string
		cps         []types.ProviderTicker
		cfg         config.WebSocketConfig
		expected    func() []handlers.WebsocketEncodedMessage
		expectedErr bool
	}{
		{
			name: "single currency pair",
			cps: []types.ProviderTicker{
				btcusdt,
			},
			cfg: mexc.DefaultWebSocketConfig,
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := `{"method":"SUBSCRIPTION","params":["spot@public.miniTicker.v3.api@BTCUSDT@UTC+8"]}`
				return []handlers.WebsocketEncodedMessage{[]byte(msg)}
			},
			expectedErr: false,
		},
		{
			name: "multiple currency pairs",
			cps: []types.ProviderTicker{
				btcusdt,
				ethusdt,
				atomusdc,
			},
			cfg: mexc.DefaultWebSocketConfig,
			expected: func() []handlers.WebsocketEncodedMessage {
				msg1 := `{"method":"SUBSCRIPTION","params":["spot@public.miniTicker.v3.api@BTCUSDT@UTC+8"]}`
				msg2 := `{"method":"SUBSCRIPTION","params":["spot@public.miniTicker.v3.api@ETHUSDT@UTC+8"]}`
				msg3 := `{"method":"SUBSCRIPTION","params":["spot@public.miniTicker.v3.api@ATOMUSDC@UTC+8"]}`
				return []handlers.WebsocketEncodedMessage{[]byte(msg1), []byte(msg2), []byte(msg3)}
			},
			expectedErr: false,
		},
		{
			name: "multiple currency pairs with batch",
			cps: []types.ProviderTicker{
				btcusdt,
				ethusdt,
			},
			cfg: batchCfg,
			expected: func() []handlers.WebsocketEncodedMessage {
				msg1 := `{"method":"SUBSCRIPTION","params":["spot@public.miniTicker.v3.api@BTCUSDT@UTC+8","spot@public.miniTicker.v3.api@ETHUSDT@UTC+8"]}`
				return []handlers.WebsocketEncodedMessage{[]byte(msg1)}
			},
			expectedErr: false,
		},
		{
			name: "multiple currency pairs with batch and remainder",
			cps: []types.ProviderTicker{
				btcusdt,
				ethusdt,
				atomusdc,
			},
			cfg: batchCfg,
			expected: func() []handlers.WebsocketEncodedMessage {
				msg1 := `{"method":"SUBSCRIPTION","params":["spot@public.miniTicker.v3.api@BTCUSDT@UTC+8","spot@public.miniTicker.v3.api@ETHUSDT@UTC+8"]}`
				msg2 := `{"method":"SUBSCRIPTION","params":["spot@public.miniTicker.v3.api@ATOMUSDC@UTC+8"]}`
				return []handlers.WebsocketEncodedMessage{[]byte(msg1), []byte(msg2)}
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wsHandler, err := mexc.NewWebSocketDataHandler(logger, tc.cfg)
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

func TestHeartBeatMessages(t *testing.T) {
	wsHandler, err := mexc.NewWebSocketDataHandler(logger, mexc.DefaultWebSocketConfig)
	require.NoError(t, err)

	expected := []handlers.WebsocketEncodedMessage{
		[]byte(`{"id":0,"code":0,"msg":"PING"}`),
	}

	msgs, err := wsHandler.HeartBeatMessages()
	require.NoError(t, err)
	require.Equal(t, expected, msgs)
}
