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

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	"github.com/skip-mev/slinky/providers/websockets/cryptodotcom"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	logger = zap.NewExample()
	mogusd = mmtypes.NewTicker("MOG", "USD", 8, 1)
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
					constants.BITCOIN_USD: types.NewPriceResult(big.NewInt(4206900000000), time.Now()),
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
					constants.BITCOIN_USD:  types.NewPriceResult(big.NewInt(4206900000000), time.Now()),
					constants.ETHEREUM_USD: types.NewPriceResult(big.NewInt(200000000000), time.Now()),
					constants.SOLANA_USD:   types.NewPriceResult(big.NewInt(100000000000), time.Now()),
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
					constants.BITCOIN_USD: types.NewPriceResult(big.NewInt(4206900000000), time.Now()),
				},
				UnResolved: types.UnResolvedPrices{
					constants.SOLANA_USD: providertypes.UnresolvedResult{
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
			marketConfig, err := types.NewProviderMarketMap(cryptodotcom.Name, cryptodotcom.DefaultMarketConfig)
			require.NoError(t, err)

			wsHandler, err := cryptodotcom.NewWebSocketDataHandler(logger, marketConfig, cryptodotcom.DefaultWebSocketConfig)
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
		msg         cryptodotcom.InstrumentRequestMessage
		expectedErr bool
	}{
		{
			name:        "no currency pairs",
			cps:         []mmtypes.Ticker{},
			msg:         cryptodotcom.InstrumentRequestMessage{},
			expectedErr: true,
		},
		{
			name: "one currency pair",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USD,
			},
			msg: cryptodotcom.InstrumentRequestMessage{
				Method: "subscribe",
				Params: cryptodotcom.InstrumentParams{
					Channels: []string{"ticker.BTCUSD-PERP"},
				},
			},
			expectedErr: false,
		},
		{
			name: "multiple currency pairs",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USD,
				constants.ETHEREUM_USD,
				constants.SOLANA_USD,
			},
			msg: cryptodotcom.InstrumentRequestMessage{
				Method: "subscribe",
				Params: cryptodotcom.InstrumentParams{
					Channels: []string{
						"ticker.BTCUSD-PERP",
						"ticker.ETHUSD-PERP",
						"ticker.SOLUSD-PERP",
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "one found and one not found",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USD,
				mogusd,
			},
			msg: cryptodotcom.InstrumentRequestMessage{
				Method: "subscribe",
				Params: cryptodotcom.InstrumentParams{
					Channels: []string{"ticker.BTCUSD-PERP"},
				},
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			marketConfig, err := types.NewProviderMarketMap(cryptodotcom.Name, cryptodotcom.DefaultMarketConfig)
			require.NoError(t, err)

			wsHandler, err := cryptodotcom.NewWebSocketDataHandler(logger, marketConfig, cryptodotcom.DefaultWebSocketConfig)
			require.NoError(t, err)

			msgs, err := wsHandler.CreateMessages(tc.cps)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			expectedBz, err := json.Marshal(tc.msg)
			require.NoError(t, err)
			require.Equal(t, 1, len(msgs))
			require.EqualValues(t, expectedBz, msgs[0])
		})
	}
}
