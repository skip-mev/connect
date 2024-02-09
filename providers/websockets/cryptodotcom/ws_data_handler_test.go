package cryptodotcom_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/cryptodotcom"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	providerCfg = config.ProviderConfig{
		Name:      cryptodotcom.Name,
		WebSocket: cryptodotcom.DefaultWebSocketConfig,
		Market: config.MarketConfig{
			Name: cryptodotcom.Name,
			CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
				"BITCOIN/USD": {
					Ticker:       "BTCUSD-PERP",
					CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
				},
				"ETHEREUM/USD": {
					Ticker:       "ETHUSD-PERP",
					CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
				},
				"SOLANA/USD": {
					Ticker:       "SOLUSD-PERP",
					CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USD"),
				},
			},
		},
	}

	btcusd = slinkytypes.NewCurrencyPair("BITCOIN", "USD")
	ethusd = slinkytypes.NewCurrencyPair("ETHEREUM", "USD")
	solusd = oracletypes.NewCurrencyPair("SOLANA", "USD")

	logger = zap.NewExample()
)

func TestHandleMessage(t *testing.T) {
	testCases := []struct {
		name         string
		msg          func() []byte
		resp         providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]
		expUpdateMsg func() []handlers.WebsocketEncodedMessage
		expErr       bool
	}{
		{
			name: "cannot unmarshal to base message",
			msg: func() []byte {
				return []byte(`no rizz message`)
			},
			resp:         providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{},
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
			resp:         providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{},
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
			resp:         providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				Resolved: map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: providertypes.NewResult[*big.Int](big.NewInt(4206900000000), time.Now()),
				},
				UnResolved: map[slinkytypes.CurrencyPair]error{},
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
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{},
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
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				Resolved: map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: providertypes.NewResult[*big.Int](big.NewInt(4206900000000), time.Now()),
					ethusd: providertypes.NewResult[*big.Int](big.NewInt(200000000000), time.Now()),
					solusd: providertypes.NewResult[*big.Int](big.NewInt(100000000000), time.Now()),
				},
				UnResolved: map[slinkytypes.CurrencyPair]error{},
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
			resp: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				Resolved: map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: providertypes.NewResult[*big.Int](big.NewInt(4206900000000), time.Now()),
				},
				UnResolved: map[slinkytypes.CurrencyPair]error{
					solusd: fmt.Errorf("failed to parse price $42,069.00: invalid syntax"),
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
			wsHandler, err := cryptodotcom.NewWebSocketDataHandler(logger, providerCfg)
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
		cps         []slinkytypes.CurrencyPair
		msg         cryptodotcom.InstrumentRequestMessage
		expectedErr bool
	}{
		{
			name:        "no currency pairs",
			cps:         []slinkytypes.CurrencyPair{},
			msg:         cryptodotcom.InstrumentRequestMessage{},
			expectedErr: true,
		},
		{
			name: "one currency pair",
			cps:  []slinkytypes.CurrencyPair{btcusd},
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
			cps:  []slinkytypes.CurrencyPair{btcusd, ethusd, solusd},
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
			cps:  []slinkytypes.CurrencyPair{btcusd, oracletypes.NewCurrencyPair("MOG", "USD")},
			msg: cryptodotcom.InstrumentRequestMessage{
				Method: "subscribe",
				Params: cryptodotcom.InstrumentParams{
					Channels: []string{"ticker.BTCUSD-PERP"},
				},
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wsHandler, err := cryptodotcom.NewWebSocketDataHandler(logger, providerCfg)
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
