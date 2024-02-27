package manager_test

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/providers/apis/binance"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"go.uber.org/zap"
)

var (
	logger = zap.NewExample()

	oracleCfg = config.OracleConfig{
		Production: true,
		Metrics: config.MetricsConfig{
			Enabled: false,
		},
		UpdateInterval: 1500 * time.Millisecond,
		MaxPriceAge:    2 * time.Minute,
		Providers: []config.ProviderConfig{
			{
				Name: binance.Name,
				API:  binance.DefaultUSAPIConfig,
			},
			{
				Name: coinbase.Name,
				API:  coinbase.DefaultAPIConfig,
			},
			{
				Name:      okx.Name,
				WebSocket: okx.DefaultWebSocketConfig,
			},
		},
	}

	// Coinbase and OKX are supported by the marketmap
	marketMap = mmtypes.MarketMap{
		Tickers: map[string]mmtypes.Ticker{
			constants.BITCOIN_USD.String():  constants.BITCOIN_USD,
			constants.ETHEREUM_USD.String(): constants.ETHEREUM_USD,
		},
		Providers: map[string]mmtypes.Providers{
			constants.BITCOIN_USD.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.BITCOIN_USD],
					okx.DefaultMarketConfig[constants.BITCOIN_USD],
				},
			},
			constants.ETHEREUM_USD.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.ETHEREUM_USD],
					okx.DefaultMarketConfig[constants.ETHEREUM_USD],
				},
			},
		},
	}
)
