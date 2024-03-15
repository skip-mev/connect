package binance

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

// NOTE: All documentation for this file can be located on the Binance GitHub
// API documentation: https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#symbol-price-ticker. This
// API does not require a subscription to use (i.e. No API key is required).

const (
	// Name is the name of the Binance provider.
	Name = "Binance"

	// URL is the base URL of the Binance API. This includes the base and quote
	// currency pairs that need to be inserted into the URL. This URL should be utilized
	// by Non-US users.
	URL = "https://api.binance.com/api/v3/ticker/price?symbols=%s%s%s"

	// US_URL is the base URL of the Binance US API. This includes the base and quote
	// currency pairs that need to be inserted into the URL. This URL should be utilized
	// by US users. Note that the US URL does not support all the currency pairs that
	// the Non-US URL supports.
	US_URL = "https://api.binance.us/api/v3/ticker/price?symbols=%s%s%s"

	Quotation    = "%22"
	Separator    = ","
	LeftBracket  = "%5B"
	RightBracket = "%5D"
)

var (
	// DefaultUSAPIConfig is the default configuration for the Binance API.
	DefaultUSAPIConfig = config.APIConfig{
		Name:             Name,
		Atomic:           true,
		Enabled:          true,
		Timeout:          500 * time.Millisecond,
		Interval:         150 * time.Millisecond,
		ReconnectTimeout: 2000 * time.Millisecond,
		MaxQueries:       1,
		URL:              US_URL,
		Type:             types.ConfigType,
	}

	// DefaultNonUSAPIConfig is the default configuration for the Binance API.
	DefaultNonUSAPIConfig = config.APIConfig{
		Name:             Name,
		Atomic:           true,
		Enabled:          true,
		Timeout:          500 * time.Millisecond,
		Interval:         150 * time.Millisecond,
		ReconnectTimeout: 2000 * time.Millisecond,
		MaxQueries:       1,
		URL:              URL,
		Type:             types.ConfigType,
	}

	// DefaultUSMarketConfig is the default US market configuration for Binance.
	DefaultUSMarketConfig = types.TickerToProviderConfig{
		constants.APE_USDT: {
			Name:           Name,
			OffChainTicker: "APEUSDT",
		},
		constants.APTOS_USDT: {
			Name:           Name,
			OffChainTicker: "APTUSDT",
		},
		constants.ARBITRUM_USDT: {
			Name:           Name,
			OffChainTicker: "ARBUSDT",
		},
		constants.ATOM_USDT: {
			Name:           Name,
			OffChainTicker: "ATOMUSDT",
		},
		constants.AVAX_USDT: {
			Name:           Name,
			OffChainTicker: "AVAXUSDT",
		},
		constants.BCH_USDT: {
			Name:           Name,
			OffChainTicker: "BCHUSDT",
		},
		constants.BITCOIN_USDC: {
			Name:           Name,
			OffChainTicker: "BTCUSDC",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "BTCUSDT",
		},
		constants.CARDANO_USDT: {
			Name:           Name,
			OffChainTicker: "ADAUSDT",
		},
		constants.CHAINLINK_USDT: {
			Name:           Name,
			OffChainTicker: "LINKUSDT",
		},
		constants.COMPOUND_USDT: {
			Name:           Name,
			OffChainTicker: "COMPUSDT",
		},
		constants.CURVE_USDT: {
			Name:           Name,
			OffChainTicker: "CRVUSDT",
		},
		constants.DOGE_USDT: {
			Name:           Name,
			OffChainTicker: "DOGEUSDT",
		},
		constants.ETC_USDT: {
			Name:           Name,
			OffChainTicker: "ETCUSDT",
		},
		constants.ETHEREUM_BITCOIN: {
			Name:           Name,
			OffChainTicker: "ETHBTC",
		},
		constants.ETHEREUM_USDC: {
			Name:           Name,
			OffChainTicker: "ETHUSDC",
		},
		constants.ETHEREUM_USDT: {
			Name:           Name,
			OffChainTicker: "ETHUSDT",
		},
		constants.FILECOIN_USDT: {
			Name:           Name,
			OffChainTicker: "FILUSDT",
		},
		constants.LIDO_USDT: {
			Name:           Name,
			OffChainTicker: "LDOUSDT",
		},
		constants.LITECOIN_USDT: {
			Name:           Name,
			OffChainTicker: "LTCUSDT",
		},
		constants.MAKER_USDT: {
			Name:           Name,
			OffChainTicker: "MKRUSDT",
		},
		constants.NEAR_USDT: {
			Name:           Name,
			OffChainTicker: "NEARUSDT",
		},
		constants.OPTIMISM_USDT: {
			Name:           Name,
			OffChainTicker: "OPUSDT",
		},
		constants.POLKADOT_USDT: {
			Name:           Name,
			OffChainTicker: "DOTUSDT",
		},
		constants.RIPPLE_USDT: {
			Name:           Name,
			OffChainTicker: "XRPUSDT",
		},
		constants.SEI_USDT: {
			Name:           Name,
			OffChainTicker: "SEIUSDT",
		},
		constants.SHIBA_USDT: {
			Name:           Name,
			OffChainTicker: "SHIBUSDT",
		},
		constants.SOLANA_USDC: {
			Name:           Name,
			OffChainTicker: "SOLUSDC",
		},
		constants.SOLANA_USDT: {
			Name:           Name,
			OffChainTicker: "SOLUSDT",
		},
		constants.STELLAR_USDT: {
			Name:           Name,
			OffChainTicker: "XLMUSDT",
		},
		constants.SUI_USDT: {
			Name:           Name,
			OffChainTicker: "SUIUSDT",
		},
		constants.TRON_USDT: {
			Name:           Name,
			OffChainTicker: "TRXUSDT",
		},
		constants.UNISWAP_USDT: {
			Name:           Name,
			OffChainTicker: "UNIUSDT",
		},
		constants.USDC_USDT: {
			Name:           Name,
			OffChainTicker: "USDCUSDT",
		},
		constants.WORLD_USDT: {
			Name:           Name,
			OffChainTicker: "WLDUSDT",
		},
	}

	// DefaultNonUSMarketConfig is the default market configuration for Binance.
	DefaultNonUSMarketConfig = types.TickerToProviderConfig{
		constants.APE_USDT: {
			Name:           Name,
			OffChainTicker: "APEUSDT",
		},
		constants.APTOS_USDT: {
			Name:           Name,
			OffChainTicker: "APTUSDT",
		},
		constants.ARBITRUM_USDT: {
			Name:           Name,
			OffChainTicker: "ARBUSDT",
		},
		constants.ATOM_USDT: {
			Name:           Name,
			OffChainTicker: "ATOMUSDT",
		},
		constants.AVAX_USDT: {
			Name:           Name,
			OffChainTicker: "AVAXUSDT",
		},
		constants.BCH_USDT: {
			Name:           Name,
			OffChainTicker: "BCHUSDT",
		},
		constants.BITCOIN_USDC: {
			Name:           Name,
			OffChainTicker: "BTCUSDC",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "BTCUSDT",
		},
		constants.CARDANO_USDT: {
			Name:           Name,
			OffChainTicker: "ADAUSDT",
		},
		constants.CHAINLINK_USDT: {
			Name:           Name,
			OffChainTicker: "LINKUSDT",
		},
		constants.COMPOUND_USDT: {
			Name:           Name,
			OffChainTicker: "COMPUSDT",
		},
		constants.CURVE_USDT: {
			Name:           Name,
			OffChainTicker: "CRVUSDT",
		},
		constants.DOGE_USDT: {
			Name:           Name,
			OffChainTicker: "DOGEUSDT",
		},
		constants.DYDX_USDT: {
			Name:           Name,
			OffChainTicker: "DYDXUSDT",
		},
		constants.ETC_USDT: {
			Name:           Name,
			OffChainTicker: "ETCUSDT",
		},
		constants.ETHEREUM_BITCOIN: {
			Name:           Name,
			OffChainTicker: "ETHBTC",
		},
		constants.ETHEREUM_USDC: {
			Name:           Name,
			OffChainTicker: "ETHUSDC",
		},
		constants.ETHEREUM_USDT: {
			Name:           Name,
			OffChainTicker: "ETHUSDT",
		},
		constants.FILECOIN_USDT: {
			Name:           Name,
			OffChainTicker: "FILUSDT",
		},
		constants.LIDO_USDT: {
			Name:           Name,
			OffChainTicker: "LDOUSDT",
		},
		constants.LITECOIN_USDT: {
			Name:           Name,
			OffChainTicker: "LTCUSDT",
		},
		constants.MAKER_USDT: {
			Name:           Name,
			OffChainTicker: "MKRUSDT",
		},
		constants.NEAR_USDT: {
			Name:           Name,
			OffChainTicker: "NEARUSDT",
		},
		constants.OPTIMISM_USDT: {
			Name:           Name,
			OffChainTicker: "OPUSDT",
		},
		constants.PEPE_USDT: {
			Name:           Name,
			OffChainTicker: "PEPEUSDT",
		},
		constants.POLKADOT_USDT: {
			Name:           Name,
			OffChainTicker: "DOTUSDT",
		},
		constants.POLYGON_USDT: {
			Name:           Name,
			OffChainTicker: "MATICUSDT",
		},
		constants.RIPPLE_USDT: {
			Name:           Name,
			OffChainTicker: "XRPUSDT",
		},
		constants.SEI_USDT: {
			Name:           Name,
			OffChainTicker: "SEIUSDT",
		},
		constants.SHIBA_USDT: {
			Name:           Name,
			OffChainTicker: "SHIBUSDT",
		},
		constants.SOLANA_USDC: {
			Name:           Name,
			OffChainTicker: "SOLUSDC",
		},
		constants.SOLANA_USDT: {
			Name:           Name,
			OffChainTicker: "SOLUSDT",
		},
		constants.STELLAR_USDT: {
			Name:           Name,
			OffChainTicker: "XLMUSDT",
		},
		constants.SUI_USDT: {
			Name:           Name,
			OffChainTicker: "SUIUSDT",
		},
		constants.TRON_USDT: {
			Name:           Name,
			OffChainTicker: "TRXUSDT",
		},
		constants.UNISWAP_USDT: {
			Name:           Name,
			OffChainTicker: "UNIUSDT",
		},
		constants.USDC_USDT: {
			Name:           Name,
			OffChainTicker: "USDCUSDT",
		},
		constants.USDT_USD: {
			Name:           Name,
			OffChainTicker: "USDTUSD",
		},
		constants.WORLD_USDT: {
			Name:           Name,
			OffChainTicker: "WLDUSDT",
		},
	}
)

type (
	// Response is the expected response returned by the Binance API.
	// The response is json formatted.
	// Response format:
	//
	//	[
	//  {
	//    "symbol": "LTCBTC",
	//    "price": "4.00000200"
	//  },
	//  {
	//    "symbol": "ETHBTC",
	//    "price": "0.07946600"
	//  }
	// ].
	Response []Data

	// Data BinanceData is the data returned by the Binance API.
	Data struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
)

// Decode decodes the given http response into a BinanceResponse.
func Decode(resp *http.Response) (Response, error) {
	// Parse the response into a BinanceResponse.
	var result Response
	err := json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
