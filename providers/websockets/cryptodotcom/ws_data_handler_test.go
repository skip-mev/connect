package cryptodotcom_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	providertypes "github.com/skip-mev/slinky/providers/types"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	"github.com/skip-mev/slinky/providers/websockets/cryptodotcom"
)

var (
	btcusd = types.DefaultProviderTicker{
		OffChainTicker: "BTCUSD",
	}
	ethusd = types.DefaultProviderTicker{
		OffChainTicker: "ETHUSD",
	}
	solusd = types.DefaultProviderTicker{
		OffChainTicker: "SOLUSD",
	}
	logger = zap.NewExample()
)

func TestHandleMessage(t *testing.T) {
	testCases := []struct {
		name         string
		msg          func() []byte
		resp         types.PriceResponse
		expUpdateMsg func() []handlers.WebsocketEncodedMessage
		expErr       bool
	}{
		{
			name: "cannot unmarshal to base message",
			msg: func() []byte {
				return []byte(`no rizz message`)
			},
			resp:         types.PriceResponse{},
			expUpdateMsg: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:       true,
		},
		{
			name: "unknown method type",
			msg: func() []byte {
				msg := cryptodotcom.InstrumentResponseMessage{
					Method: "unknown",
				}
				bz, err := json.Marshal(msg)
				require.NoError(t, err)
				return bz
			},
			resp:         types.PriceResponse{},
			expUpdateMsg: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:       true,
		},
		{
			name: "unknown status code",
			msg: func() []byte {
				msg := cryptodotcom.InstrumentResponseMessage{
					Method: string(cryptodotcom.InstrumentMethod),
					Code:   1,
				}
				bz, err := json.Marshal(msg)
				require.NoError(t, err)
				return bz
			},
			resp:         types.PriceResponse{},
			expUpdateMsg: func() []handlers.WebsocketEncodedMessage { return nil },
			expErr:       true,
		},
		{
			name: "heartbeat",
			msg: func() []byte {
				msg := cryptodotcom.InstrumentResponseMessage{
					ID:     42069,
					Method: string(cryptodotcom.HeartBeatRequestMethod),
					Code:   0,
				}
				bz, err := json.Marshal(msg)
				require.NoError(t, err)
				return bz
			},
			resp: types.PriceResponse{},
			expUpdateMsg: func() []handlers.WebsocketEncodedMessage {
				msg := cryptodotcom.HeartBeatResponseMessage{
					ID:     42069,
					Method: string(cryptodotcom.HeartBeatResponseMethod),
				}
				bz, err := json.Marshal(msg)
				require.NoError(t, err)
				return []handlers.WebsocketEncodedMessage{bz}
			},
			expErr: false,
		},
		{
			name: "instrument response with no data",
			msg: func() []byte {
				msg := cryptodotcom.InstrumentResponseMessage{
					ID:     42069,
					Method: string(cryptodotcom.InstrumentMethod),
					Code:   0,
					Result: cryptodotcom.InstrumentResult{
						Data: []cryptodotcom.InstrumentData{},
					},
				}
				bz, err := json.Marshal(msg)
				require.NoError(t, err)
				return bz
			},
			resp: types.PriceResponse{},
			expUpdateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: true,
		},
		{
			name: "instrument response with one instrument",
			msg: func() []byte {
				msg := cryptodotcom.InstrumentResponseMessage{
					ID:     42069,
					Method: string(cryptodotcom.InstrumentMethod),
					Code:   0,
					Result: cryptodotcom.InstrumentResult{
						Data: []cryptodotcom.InstrumentData{
							{
								Name:             "BTCUSD-PERP",
								LatestTradePrice: "42069",
							},
						},
					},
				}
				bz, err := json.Marshal(msg)
				require.NoError(t, err)
				return bz
			},
			resp: types.PriceResponse{
				Resolved: types.ResolvedPrices{
					btcusd: types.NewPriceResult(big.NewFloat(42069.00), time.Now()),
				},
				UnResolved: types.UnResolvedPrices{},
			},
			expUpdateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: false,
		},
		{
			name: "unknown instrument",
			msg: func() []byte {
				msg := cryptodotcom.InstrumentResponseMessage{
					ID:     42069,
					Method: string(cryptodotcom.InstrumentMethod),
					Code:   0,
					Result: cryptodotcom.InstrumentResult{
						Data: []cryptodotcom.InstrumentData{
							{
								Name:             "MOGUSD-PERP",
								LatestTradePrice: "42069",
							},
						},
					},
				}
				bz, err := json.Marshal(msg)
				require.NoError(t, err)
				return bz
			},
			resp: types.PriceResponse{},
			expUpdateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: false,
		},
		{
			name: "instrument response with multiple instruments",
			msg: func() []byte {
				msg := cryptodotcom.InstrumentResponseMessage{
					ID:     42069,
					Method: string(cryptodotcom.InstrumentMethod),
					Code:   0,
					Result: cryptodotcom.InstrumentResult{
						Data: []cryptodotcom.InstrumentData{
							{
								Name:             "BTCUSD-PERP",
								LatestTradePrice: "42069",
							},
							{
								Name:             "ETHUSD-PERP",
								LatestTradePrice: "2000",
							},
							{
								Name:             "SOLUSD-PERP",
								LatestTradePrice: "1000",
							},
						},
					},
				}
				bz, err := json.Marshal(msg)
				require.NoError(t, err)
				return bz
			},
			resp: types.PriceResponse{
				Resolved: types.ResolvedPrices{
					btcusd: types.NewPriceResult(big.NewFloat(42069.00), time.Now()),
					ethusd: types.NewPriceResult(big.NewFloat(2000.00), time.Now()),
					solusd: types.NewPriceResult(big.NewFloat(1000.00), time.Now()),
				},
				UnResolved: types.UnResolvedPrices{},
			},
			expUpdateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: false,
		},
		{
			name: "instrument response with one instrument and one bad price instrument",
			msg: func() []byte {
				msg := cryptodotcom.InstrumentResponseMessage{
					ID:     42069,
					Method: string(cryptodotcom.InstrumentMethod),
					Code:   0,
					Result: cryptodotcom.InstrumentResult{
						Data: []cryptodotcom.InstrumentData{
							{
								Name:             "BTCUSD-PERP",
								LatestTradePrice: "42069",
							},
							{
								Name:             "SOLUSD-PERP",
								LatestTradePrice: "$42,069.00",
							},
						},
					},
				}
				bz, err := json.Marshal(msg)
				require.NoError(t, err)
				return bz
			},
			resp: types.PriceResponse{
				Resolved: types.ResolvedPrices{
					btcusd: types.NewPriceResult(big.NewFloat(42069.00), time.Now()),
				},
				UnResolved: types.UnResolvedPrices{
					solusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("failed to parse price $42,069.00: invalid syntax"), providertypes.ErrorWebSocketGeneral),
					},
				},
			},
			expUpdateMsg: func() []handlers.WebsocketEncodedMessage {
				return nil
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wsHandler, err := cryptodotcom.NewWebSocketDataHandler(logger, cryptodotcom.DefaultWebSocketConfig)
			require.NoError(t, err)

			// Update the cache since it is assumed that CreateMessages is executed before anything else.
			_, err = wsHandler.CreateMessages([]types.ProviderTicker{btcusd, ethusd, solusd})
			require.NoError(t, err)

			resp, updateMsg, err := wsHandler.HandleMessage(tc.msg())
			if tc.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expUpdateMsg(), updateMsg)

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
		msgs        []cryptodotcom.InstrumentRequestMessage
		expectedErr bool
	}{
		{
			name:        "no currency pairs",
			cps:         []types.ProviderTicker{},
			msgs:        []cryptodotcom.InstrumentRequestMessage{},
			expectedErr: true,
		},
		{
			name: "one currency pair",
			cps: []types.ProviderTicker{
				btcusd,
			},
			msgs: []cryptodotcom.InstrumentRequestMessage{
				{
					Method: "subscribe",
					Params: cryptodotcom.InstrumentParams{
						Channels: []string{"ticker.BTCUSD-PERP"},
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "multiple currency pairs",
			cps: []types.ProviderTicker{
				btcusd,
				ethusd,
				solusd,
			},
			msgs: []cryptodotcom.InstrumentRequestMessage{
				{
					Method: "subscribe",
					Params: cryptodotcom.InstrumentParams{
						Channels: []string{
							"ticker.BTCUSD-PERP",
						},
					},
				},
				{
					Method: "subscribe",
					Params: cryptodotcom.InstrumentParams{
						Channels: []string{
							"ticker.ETHUSD-PERP",
						},
					},
				},
				{
					Method: "subscribe",
					Params: cryptodotcom.InstrumentParams{
						Channels: []string{
							"ticker.SOLUSD-PERP",
						},
					},
				},
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wsHandler, err := cryptodotcom.NewWebSocketDataHandler(logger, cryptodotcom.DefaultWebSocketConfig)
			require.NoError(t, err)

			msgs, err := wsHandler.CreateMessages(tc.cps)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, len(tc.msgs), len(msgs))
			for i := range tc.msgs {
				expectedBz, err := json.Marshal(tc.msgs[i])
				require.NoError(t, err)
				require.EqualValues(t, expectedBz, msgs[i])
			}
		})
	}
}
