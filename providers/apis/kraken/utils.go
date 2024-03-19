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
	Name = "Kraken"

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
		Timeout:          500 * time.Millisecond,
		Interval:         400 * time.Millisecond,
		ReconnectTimeout: 2000 * time.Millisecond,
		MaxQueries:       1,
		URL:              URL,
	}

	// DefaultMarketConfig is the default market configuration for Kraken.
	DefaultMarketConfig = types.TickerToProviderConfig{
		constants.APE_USDT: {
			Name:           Name,
			OffChainTicker: "APEUSDT",
		},
		constants.APTOS_USD: {
			Name:           Name,
			OffChainTicker: "APTUSD",
		},
		constants.ARBITRUM_USD: {
			Name:           Name,
			OffChainTicker: "ARBUSD",
		},
		constants.ATOM_USDT: {
			Name:           Name,
			OffChainTicker: "ATOMUSDT",
		},
		constants.ATOM_USD: {
			Name:           Name,
			OffChainTicker: "ATOMUSD",
		},
		constants.AVAX_USDT: {
			Name:           Name,
			OffChainTicker: "AVAXUSDT",
		},
		constants.AVAX_USD: {
			Name:           Name,
			OffChainTicker: "AVAXUSD",
		},
		constants.BCH_USDT: {
			Name:           Name,
			OffChainTicker: "BCHUSDT",
		},
		constants.BCH_USD: {
			Name:           Name,
			OffChainTicker: "BCHUSD",
		},
		constants.BITCOIN_USDC: {
			Name:           Name,
			OffChainTicker: "XBTUSDC",
		},
		constants.BITCOIN_USD: {
			Name:           Name,
			OffChainTicker: "XXBTZUSD",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "XBTUSDT",
		},
		constants.CARDANO_USDT: {
			Name:           Name,
			OffChainTicker: "ADAUSDT",
		},
		constants.CARDANO_USD: {
			Name:           Name,
			OffChainTicker: "ADAUSD",
		},
		constants.CHAINLINK_USDT: {
			Name:           Name,
			OffChainTicker: "LINKUSDT",
		},
		constants.CHAINLINK_USD: {
			Name:           Name,
			OffChainTicker: "LINKUSD",
		},
		constants.COMPOUND_USD: {
			Name:           Name,
			OffChainTicker: "COMPUSD",
		},
		constants.CURVE_USD: {
			Name:           Name,
			OffChainTicker: "CRVUSD",
		},
		constants.DOGE_USDT: {
			Name:           Name,
			OffChainTicker: "XDGUSDT",
		},
		constants.DOGE_USD: {
			Name:           Name,
			OffChainTicker: "XDGUSD",
		},
		constants.DYDX_USD: {
			Name:           Name,
			OffChainTicker: "DYDXUSD",
		},
		constants.ETC_USD: {
			Name:           Name,
			OffChainTicker: "ETCUSD",
		},
		constants.ETHEREUM_BITCOIN: {
			Name:           Name,
			OffChainTicker: "XETHXXBT",
		},
		constants.ETHEREUM_USDC: {
			Name:           Name,
			OffChainTicker: "ETHUSDC",
		},
		constants.ETHEREUM_USDT: {
			Name:           Name,
			OffChainTicker: "ETHUSDT",
		},
		constants.ETHEREUM_USD: {
			Name:           Name,
			OffChainTicker: "XETHZUSD",
		},
		constants.FILECOIN_USD: {
			Name:           Name,
			OffChainTicker: "FILUSD",
		},
		constants.LIDO_USD: {
			Name:           Name,
			OffChainTicker: "LDOUSD",
		},
		constants.LITECOIN_USDT: {
			Name:           Name,
			OffChainTicker: "LTCUSDT",
		},
		constants.LITECOIN_USD: {
			Name:           Name,
			OffChainTicker: "XLTCZUSD",
		},
		constants.MAKER_USD: {
			Name:           Name,
			OffChainTicker: "MKRUSD",
		},
		constants.NEAR_USD: {
			Name:           Name,
			OffChainTicker: "NEARUSD",
		},
		constants.OPTIMISM_USD: {
			Name:           Name,
			OffChainTicker: "OPUSD",
		},
		constants.PEPE_USD: {
			Name:           Name,
			OffChainTicker: "PEPEUSD",
		},
		constants.POLKADOT_USDT: {
			Name:           Name,
			OffChainTicker: "DOTUSDT",
		},
		constants.POLKADOT_USD: {
			Name:           Name,
			OffChainTicker: "DOTUSD",
		},
		constants.POLYGON_USDT: {
			Name:           Name,
			OffChainTicker: "MATICUSDT",
		},
		constants.POLYGON_USD: {
			Name:           Name,
			OffChainTicker: "MATICUSD",
		},
		constants.RIPPLE_USDT: {
			Name:           Name,
			OffChainTicker: "XRPUSDT",
		},
		constants.RIPPLE_USD: {
			Name:           Name,
			OffChainTicker: "XXRPZUSD",
		},
		constants.SEI_USD: {
			Name:           Name,
			OffChainTicker: "SEIUSD",
		},
		constants.SHIBA_USDT: {
			Name:           Name,
			OffChainTicker: "SHIBUSDT",
		},
		constants.SHIBA_USD: {
			Name:           Name,
			OffChainTicker: "SHIBUSD",
		},
		constants.SOLANA_USDT: {
			Name:           Name,
			OffChainTicker: "SOLUSDT",
		},
		constants.SOLANA_USD: {
			Name:           Name,
			OffChainTicker: "SOLUSD",
		},
		constants.STELLAR_USD: {
			Name:           Name,
			OffChainTicker: "XXLMZUSD",
		},
		constants.SUI_USD: {
			Name:           Name,
			OffChainTicker: "SUIUSD",
		},
		constants.TRON_USD: {
			Name:           Name,
			OffChainTicker: "TRXUSD",
		},
		constants.UNISWAP_USD: {
			Name:           Name,
			OffChainTicker: "UNIUSD",
		},
		constants.USDC_USDT: {
			Name:           Name,
			OffChainTicker: "USDCUSDT",
		},
		constants.USDT_USD: {
			Name:           Name,
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
