package mexc_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/mexc"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	cfg = config.ProviderConfig{
		Name:      mexc.Name,
		WebSocket: mexc.DefaultWebSocketConfig,
		Market:    mexc.DefaultMarketConfig,
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
			name: "pong message",
			msg: func() []byte {
				return []byte(`{"id":0,"code":0,"msg":"PONG"}`)
			},
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				Resolved: map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					oracletypes.NewCurrencyPair("BITCOIN", "USDT"): {
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
			name: "unsupported market price update",
			msg: func() []byte {
				msg := `{"c":"spot@public.miniTicker.v3.api@MOGUSDT@UTC+8","d":{"s":"MOGUSDT","p":"10000.00"}}`
				return []byte(msg)
			},
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				UnResolved: map[slinkytypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USDT"): fmt.Errorf("invalid channel"),
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
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				UnResolved: map[slinkytypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USDT"): fmt.Errorf("invalid price"),
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
			wsHandler, err := mexc.NewWebSocketDataHandler(logger, cfg)
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
			name: "single currency pair",
			cps: []slinkytypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := `{"method":"SUBSCRIPTION","params":["spot@public.miniTicker.v3.api@BTCUSDT@UTC+8"]}`
				return []handlers.WebsocketEncodedMessage{[]byte(msg)}
			},
			expectedErr: false,
		},
		{
			name: "multiple currency pairs",
			cps: []slinkytypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
				oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
				oracletypes.NewCurrencyPair("ATOM", "USDC"),
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := `{"method":"SUBSCRIPTION","params":["spot@public.miniTicker.v3.api@BTCUSDT@UTC+8","spot@public.miniTicker.v3.api@ETHUSDT@UTC+8","spot@public.miniTicker.v3.api@ATOMUSDC@UTC+8"]}`
				return []handlers.WebsocketEncodedMessage{[]byte(msg)}
			},
			expectedErr: false,
		},
		{
			name: "unsupported currency pair",
			cps: []slinkytypes.CurrencyPair{
				oracletypes.NewCurrencyPair("MOG", "USD"),
			},
			expected:    func() []handlers.WebsocketEncodedMessage { return nil },
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := mexc.NewWebSocketDataHandler(logger, cfg)
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

func TestHeartBeatMessages(t *testing.T) {
	handler, err := mexc.NewWebSocketDataHandler(logger, cfg)
	require.NoError(t, err)

	expected := []handlers.WebsocketEncodedMessage{
		[]byte(`{"id":0,"code":0,"msg":"PING"}`),
	}

	msgs, err := handler.HeartBeatMessages()
	require.NoError(t, err)
	require.Equal(t, expected, msgs)
}
