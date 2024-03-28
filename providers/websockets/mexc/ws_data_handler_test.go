package mexc_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/mexc"
	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
)

var (
	logger = zap.NewExample()
	mogusd = mmtypes.NewTicker("MOG", "USD", 8, 1)
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
					constants.BITCOIN_USDT: {
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
					constants.BITCOIN_USDT: providertypes.UnresolvedResult{
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
					constants.BITCOIN_USDT: providertypes.UnresolvedResult{
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
			marketConfig, err := types.NewProviderMarketMap(mexc.Name, mexc.DefaultProviderConfig)
			require.NoError(t, err)

			wsHandler, err := mexc.NewWebSocketDataHandler(logger, marketConfig, mexc.DefaultWebSocketConfig)
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
		cps         []mmtypes.Ticker
		expected    func() []handlers.WebsocketEncodedMessage
		expectedErr bool
	}{
		{
			name: "single currency pair",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := `{"method":"SUBSCRIPTION","params":["spot@public.miniTicker.v3.api@BTCUSDT@UTC+8"]}`
				return []handlers.WebsocketEncodedMessage{[]byte(msg)}
			},
			expectedErr: false,
		},
		{
			name: "multiple currency pairs",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
				constants.ETHEREUM_USDT,
				constants.ATOM_USDC,
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				msg := `{"method":"SUBSCRIPTION","params":["spot@public.miniTicker.v3.api@BTCUSDT@UTC+8","spot@public.miniTicker.v3.api@ETHUSDT@UTC+8","spot@public.miniTicker.v3.api@ATOMUSDC@UTC+8"]}`
				return []handlers.WebsocketEncodedMessage{[]byte(msg)}
			},
			expectedErr: false,
		},
		{
			name: "unsupported currency pair",
			cps: []mmtypes.Ticker{
				mogusd,
			},
			expected:    func() []handlers.WebsocketEncodedMessage { return nil },
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			marketConfig, err := types.NewProviderMarketMap(mexc.Name, mexc.DefaultProviderConfig)
			require.NoError(t, err)

			wsHandler, err := mexc.NewWebSocketDataHandler(logger, marketConfig, mexc.DefaultWebSocketConfig)
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
	marketConfig, err := types.NewProviderMarketMap(mexc.Name, mexc.DefaultProviderConfig)
	require.NoError(t, err)

	wsHandler, err := mexc.NewWebSocketDataHandler(logger, marketConfig, mexc.DefaultWebSocketConfig)
	require.NoError(t, err)

	expected := []handlers.WebsocketEncodedMessage{
		[]byte(`{"id":0,"code":0,"msg":"PING"}`),
	}

	msgs, err := wsHandler.HeartBeatMessages()
	require.NoError(t, err)
	require.Equal(t, expected, msgs)
}
