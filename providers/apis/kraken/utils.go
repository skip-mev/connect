package kraken

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

// NOTE: All documentation for this file can be located on the Kraken docs.
// API documentation: https://docs.kraken.com/rest/. This
// API does not require a subscription to use (i.e. No API key is required).

const (
	// Name is the name of the Kraken API provider.
	Name = "kraken_api"

	Type = types.ConfigType

	// URL is the base URL of the Kraken API. This includes the base and quote
	// currency pairs that need to be inserted into the URL.
	URL = "https://api.kraken.com/0/public/Ticker?pair=%s"

	// Separator is the character that separates tickers in the query URL.
	Separator = ","
)

var (
	// DefaultAPIConfig is the default configuration for the Kraken API.
	DefaultAPIConfig = config.APIConfig{
		Name:             Name,
		Atomic:           true,
		Enabled:          true,
		Timeout:          3000 * time.Millisecond,
		Interval:         600 * time.Millisecond,
		ReconnectTimeout: 2000 * time.Millisecond,
		MaxQueries:       1,
		URL:              URL,
	}

	DefaultProviderConfig = config.ProviderConfig{
		Name: Name,
		API:  DefaultAPIConfig,
		Type: Type,
	}

	// DefaultMarketConfig is the default market configuration for Kraken.
	DefaultMarketConfig = types.CurrencyPairsToProviderTickers{
		constants.APE_USDT: {
			OffChainTicker: "APEUSDT",
		},
		constants.APTOS_USD: {
			OffChainTicker: "APTUSD",
		},
		constants.ARBITRUM_USD: {
			OffChainTicker: "ARBUSD",
		},
		constants.ATOM_USDT: {
			OffChainTicker: "ATOMUSDT",
		},
		constants.ATOM_USD: {
			OffChainTicker: "ATOMUSD",
		},
		constants.AVAX_USDT: {
			OffChainTicker: "AVAXUSDT",
		},
		constants.AVAX_USD: {
			OffChainTicker: "AVAXUSD",
		},
		constants.BCH_USDT: {
			OffChainTicker: "BCHUSDT",
		},
		constants.BCH_USD: {
			OffChainTicker: "BCHUSD",
		},
		constants.BITCOIN_USDC: {
			OffChainTicker: "XBTUSDC",
		},
		constants.BITCOIN_USD: {
			OffChainTicker: "XXBTZUSD",
		},
		constants.BITCOIN_USDT: {
			OffChainTicker: "XBTUSDT",
		},
		constants.CARDANO_USDT: {
			OffChainTicker: "ADAUSDT",
		},
		constants.CARDANO_USD: {
			OffChainTicker: "ADAUSD",
		},
		constants.CHAINLINK_USDT: {
			OffChainTicker: "LINKUSDT",
		},
		constants.CHAINLINK_USD: {
			OffChainTicker: "LINKUSD",
		},
		constants.COMPOUND_USD: {
			OffChainTicker: "COMPUSD",
		},
		constants.CURVE_USD: {
			OffChainTicker: "CRVUSD",
		},
		constants.DOGE_USDT: {
			OffChainTicker: "XDGUSDT",
		},
		constants.DOGE_USD: {
			OffChainTicker: "XDGUSD",
		},
		constants.DYDX_USD: {
			OffChainTicker: "DYDXUSD",
		},
		constants.ETC_USD: {
			OffChainTicker: "ETCUSD",
		},
		constants.ETHEREUM_BITCOIN: {
			OffChainTicker: "XETHXXBT",
		},
		constants.ETHEREUM_USDC: {
			OffChainTicker: "ETHUSDC",
		},
		constants.ETHEREUM_USDT: {
			OffChainTicker: "ETHUSDT",
		},
		constants.ETHEREUM_USD: {
			OffChainTicker: "XETHZUSD",
		},
		constants.FILECOIN_USD: {
			OffChainTicker: "FILUSD",
		},
		constants.LIDO_USD: {
			OffChainTicker: "LDOUSD",
		},
		constants.LITECOIN_USDT: {
			OffChainTicker: "LTCUSDT",
		},
		constants.LITECOIN_USD: {
			OffChainTicker: "XLTCZUSD",
		},
		constants.MAKER_USD: {
			OffChainTicker: "MKRUSD",
		},
		constants.NEAR_USD: {
			OffChainTicker: "NEARUSD",
		},
		constants.OPTIMISM_USD: {
			OffChainTicker: "OPUSD",
		},
		constants.PEPE_USD: {
			OffChainTicker: "PEPEUSD",
		},
		constants.POLKADOT_USDT: {
			OffChainTicker: "DOTUSDT",
		},
		constants.POLKADOT_USD: {
			OffChainTicker: "DOTUSD",
		},
		constants.POLYGON_USDT: {
			OffChainTicker: "MATICUSDT",
		},
		constants.POLYGON_USD: {
			OffChainTicker: "MATICUSD",
		},
		constants.RIPPLE_USDT: {
			OffChainTicker: "XRPUSDT",
		},
		constants.RIPPLE_USD: {
			OffChainTicker: "XXRPZUSD",
		},
		constants.SEI_USD: {
			OffChainTicker: "SEIUSD",
		},
		constants.SHIBA_USDT: {
			OffChainTicker: "SHIBUSDT",
		},
		constants.SHIBA_USD: {
			OffChainTicker: "SHIBUSD",
		},
		constants.SOLANA_USDT: {
			OffChainTicker: "SOLUSDT",
		},
		constants.SOLANA_USD: {
			OffChainTicker: "SOLUSD",
		},
		constants.STELLAR_USD: {
			OffChainTicker: "XXLMZUSD",
		},
		constants.SUI_USD: {
			OffChainTicker: "SUIUSD",
		},
		constants.TRON_USD: {
			OffChainTicker: "TRXUSD",
		},
		constants.UNISWAP_USD: {
			OffChainTicker: "UNIUSD",
		},
		constants.USDC_USDT: {
			OffChainTicker: "USDCUSDT",
		},
		constants.USDT_USD: {
			OffChainTicker: "USDTZUSD",
		},
	}
)

// TickerResult is the result of a Kraken API call for a single ticker.
//
// https://api.kraken.com/0/public/Ticker
type TickerResult struct {
	pair            string
	ClosePriceStats []string `json:"c"`
}

func (ktr *TickerResult) LastPrice() string {
	return ktr.ClosePriceStats[0]
}

// ResponseBody returns a list of tickers for the response.  If there is an error, it will be included,
// and all Tickers will be undefined.
type ResponseBody struct {
	Errors  []string                `json:"error" validate:"omitempty"`
	Tickers map[string]TickerResult `json:"result"`
}

// Decode decodes the given http response into a TickerResult.
func Decode(resp *http.Response) (ResponseBody, error) {
	// Parse the response into a ResponseBody.
	var result ResponseBody
	err := json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
