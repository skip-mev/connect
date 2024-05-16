package bitstamp_test

import (
	"fmt"
	"math/big"
	"testing"

	providertypes "github.com/skip-mev/slinky/providers/types"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	"github.com/skip-mev/slinky/providers/websockets/bitstamp"
)

var (
	btcusd = types.DefaultProviderTicker{
		OffChainTicker: "BTCUSD",
	}
	ethusd = types.DefaultProviderTicker{
		OffChainTicker: "ETHUSD",
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
			name: "unknown message",
			msg: func() []byte {
				return []byte(`{"event":"unknown"}`)
			},
			resp:          types.PriceResponse{},
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        true,
		},
		{
			name: "heartbeat message",
			msg: func() []byte {
				return []byte(`{"event":"bts:heartbeat"}`)
			},
			resp:          types.PriceResponse{},
			updateMessage: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:        false,
		},
		{
			name: "reconnect message",
			msg: func() []byte {
				return []byte(`{"event":"bts:request_reconnect"}`)
			},
			resp: types.PriceResponse{},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return []handlers.WebsocketEncodedMessage{
					[]byte(`{"event":"bts:request_reconnect"}`),
				}
			},
			expErr: false,
		},
		{
			name: "successful subscription",
			msg: func() []byte {
				return []byte(`{"event":"bts:subscription_succeeded","channel":"live_trades_btcusd"}`)
			},
			resp: types.PriceResponse{},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: false,
		},
		{
			name: "good trade update",
			msg: func() []byte {
				return []byte(`{"event":"trade","channel":"live_trades_btcusd","data":{"microtimestamp":"1612185600000000","amount":"0.00000001","buy_order_id":0,"sell_order_id":0,"amount_str":"0.00000001","price_str":"100000.00","timestamp":"1612185600","price":"100000.00","type":1,"id":123456789}}`)
			},
			resp: types.PriceResponse{
				Resolved: types.ResolvedPrices{
					btcusd: {
						Value: big.NewFloat(100000.00),
					},
				},
			},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: false,
		},
		{
			name: "wrong channel trade update",
			msg: func() []byte {
				return []byte(`{"event":"trade","channel":"futures_ethusd","data":{"microtimestamp":"1612185600000000","amount":"0.00000001","buy_order_id":0,"sell_order_id":0,"amount_str":"0.00000001","price_str":"100000.00","timestamp":"1612185600","price":"100000.00","type":1,"id":123456789}}`)
			},
			resp: types.PriceResponse{},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: true,
		},
		{
			name: "missing ticker data in channel",
			msg: func() []byte {
				return []byte(`{"event":"trade","channel":"live_trades_","data":{"microtimestamp":"1612185600000000","amount":"0.00000001","buy_order_id":0,"sell_order_id":0,"amount_str":"0.00000001","price_str":"100000.00","timestamp":"1612185600","price":"100000.00","type":1,"id":123456789}}`)
			},
			resp: types.PriceResponse{},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: true,
		},
		{
			name: "unsupported ticker",
			msg: func() []byte {
				return []byte(`{"event":"trade","channel":"live_trades_mogusd","data":{"microtimestamp":"1612185600000000","amount":"0.00000001","buy_order_id":0,"sell_order_id":0,"amount_str":"0.00000001","price_str":"100000.00","timestamp":"1612185600","price":"100000.00","type":1,"id":123456789}}`)
			},
			resp: types.PriceResponse{},
			updateMessage: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: true,
		},
		{
			name: "bad price",
			msg: func() []byte {
				return []byte(`{"event":"trade","channel":"live_trades_btcusd","data":{"microtimestamp":"1612185600000000","amount":"0.00000001","buy_order_id":0,"sell_order_id":0,"amount_str":"0.00000001","price_str":"bad","timestamp":"1612185600","price":"bad","type":1,"id":123456789}}`)
			},
			resp: types.PriceResponse{
				UnResolved: types.UnResolvedPrices{
					btcusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("error"), providertypes.ErrorWebSocketGeneral),
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
			wsHandler, err := bitstamp.NewWebSocketDataHandler(logger, bitstamp.DefaultWebSocketConfig)
			require.NoError(t, err)

			// Update the cache since it is assumed that CreateMessages is executed before anything else.
			_, err = wsHandler.CreateMessages([]types.ProviderTicker{btcusd, ethusd})
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
	testCases := []struct {
		name        string
		cps         []types.ProviderTicker
		expected    func() []handlers.WebsocketEncodedMessage
		expectedErr bool
	}{
		{
			name: "no currency pairs",
			cps:  []types.ProviderTicker{},
			expected: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "one currency pair",
			cps: []types.ProviderTicker{
				btcusd,
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				return []handlers.WebsocketEncodedMessage{
					[]byte(`{"event":"bts:subscribe","data":{"channel":"live_trades_btcusd"}}`),
				}
			},
			expectedErr: false,
		},
		{
			name: "two currency pairs",
			cps: []types.ProviderTicker{
				btcusd,
				ethusd,
			},
			expected: func() []handlers.WebsocketEncodedMessage {
				return []handlers.WebsocketEncodedMessage{
					[]byte(`{"event":"bts:subscribe","data":{"channel":"live_trades_btcusd"}}`),
					[]byte(`{"event":"bts:subscribe","data":{"channel":"live_trades_ethusd"}}`),
				}
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wsHandler, err := bitstamp.NewWebSocketDataHandler(logger, bitstamp.DefaultWebSocketConfig)
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

func TestHeartBeat(t *testing.T) {
	wsHandler, err := bitstamp.NewWebSocketDataHandler(logger, bitstamp.DefaultWebSocketConfig)
	require.NoError(t, err)

	msgs, err := wsHandler.HeartBeatMessages()
	require.NoError(t, err)

	msg := handlers.WebsocketEncodedMessage([]byte(`{"event":"bts:heartbeat"}`))
	require.Len(t, msgs, 1)
	require.Equal(t, msg, msgs[0])
}
