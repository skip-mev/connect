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
	Name = "binance_api"

	Type = types.ConfigType

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
		Timeout:          3000 * time.Millisecond,
		Interval:         750 * time.Millisecond,
		ReconnectTimeout: 2000 * time.Millisecond,
		MaxQueries:       1,
		URL:              US_URL,
	}

	// DefaultNonUSAPIConfig is the default configuration for the Binance API.
	DefaultNonUSAPIConfig = config.APIConfig{
		Name:             Name,
		Atomic:           true,
		Enabled:          true,
		Timeout:          3000 * time.Millisecond,
		Interval:         750 * time.Millisecond,
		ReconnectTimeout: 2000 * time.Millisecond,
		MaxQueries:       1,
		URL:              URL,
	}

	DefaultProviderConfig = config.ProviderConfig{
		Name: Name,
		API:  DefaultNonUSAPIConfig,
		Type: Type,
	}

	// DefaultUSMarketConfig is the default US market configuration for Binance.
	DefaultUSMarketConfig = types.CurrencyPairsToProviderTickers{
		constants.APE_USDT: {
			OffChainTicker: "APEUSDT",
		},
		constants.APTOS_USDT: {
			OffChainTicker: "APTUSDT",
		},
		constants.ARBITRUM_USDT: {
			OffChainTicker: "ARBUSDT",
		},
		constants.ATOM_USDT: {
			OffChainTicker: "ATOMUSDT",
		},
		constants.AVAX_USDT: {
			OffChainTicker: "AVAXUSDT",
		},
		constants.BCH_USDT: {
			OffChainTicker: "BCHUSDT",
		},
		constants.BITCOIN_USDC: {
			OffChainTicker: "BTCUSDC",
		},
		constants.BITCOIN_USDT: {
			OffChainTicker: "BTCUSDT",
		},
		constants.CARDANO_USDT: {
			OffChainTicker: "ADAUSDT",
		},
		constants.CHAINLINK_USDT: {
			OffChainTicker: "LINKUSDT",
		},
		constants.COMPOUND_USDT: {
			OffChainTicker: "COMPUSDT",
		},
		constants.CURVE_USDT: {
			OffChainTicker: "CRVUSDT",
		},
		constants.DOGE_USDT: {
			OffChainTicker: "DOGEUSDT",
		},
		constants.ETC_USDT: {
			OffChainTicker: "ETCUSDT",
		},
		constants.ETHEREUM_BITCOIN: {
			OffChainTicker: "ETHBTC",
		},
		constants.ETHEREUM_USDC: {
			OffChainTicker: "ETHUSDC",
		},
		constants.ETHEREUM_USDT: {
			OffChainTicker: "ETHUSDT",
		},
		constants.FILECOIN_USDT: {
			OffChainTicker: "FILUSDT",
		},
		constants.LIDO_USDT: {
			OffChainTicker: "LDOUSDT",
		},
		constants.LITECOIN_USDT: {
			OffChainTicker: "LTCUSDT",
		},
		constants.MAKER_USDT: {
			OffChainTicker: "MKRUSDT",
		},
		constants.NEAR_USDT: {
			OffChainTicker: "NEARUSDT",
		},
		constants.OPTIMISM_USDT: {
			OffChainTicker: "OPUSDT",
		},
		constants.POLKADOT_USDT: {
			OffChainTicker: "DOTUSDT",
		},
		constants.RIPPLE_USDT: {
			OffChainTicker: "XRPUSDT",
		},
		constants.SEI_USDT: {
			OffChainTicker: "SEIUSDT",
		},
		constants.SHIBA_USDT: {
			OffChainTicker: "SHIBUSDT",
		},
		constants.SOLANA_USDC: {
			OffChainTicker: "SOLUSDC",
		},
		constants.SOLANA_USDT: {
			OffChainTicker: "SOLUSDT",
		},
		constants.STELLAR_USDT: {
			OffChainTicker: "XLMUSDT",
		},
		constants.SUI_USDT: {
			OffChainTicker: "SUIUSDT",
		},
		constants.TRON_USDT: {
			OffChainTicker: "TRXUSDT",
		},
		constants.UNISWAP_USDT: {
			OffChainTicker: "UNIUSDT",
		},
		constants.USDC_USDT: {
			OffChainTicker: "USDCUSDT",
		},
		constants.WORLD_USDT: {
			OffChainTicker: "WLDUSDT",
		},
	}

	// DefaultNonUSMarketConfig is the default market configuration for Binance.
	DefaultNonUSMarketConfig = types.CurrencyPairsToProviderTickers{
		constants.APE_USDT: {
			OffChainTicker: "APEUSDT",
		},
		constants.APTOS_USDT: {
			OffChainTicker: "APTUSDT",
		},
		constants.ARBITRUM_USDT: {
			OffChainTicker: "ARBUSDT",
		},
		constants.ATOM_USDT: {
			OffChainTicker: "ATOMUSDT",
		},
		constants.AVAX_USDT: {
			OffChainTicker: "AVAXUSDT",
		},
		constants.BCH_USDT: {
			OffChainTicker: "BCHUSDT",
		},
		constants.BITCOIN_USDC: {
			OffChainTicker: "BTCUSDC",
		},
		constants.BITCOIN_USDT: {
			OffChainTicker: "BTCUSDT",
		},
		constants.CARDANO_USDT: {
			OffChainTicker: "ADAUSDT",
		},
		constants.CHAINLINK_USDT: {
			OffChainTicker: "LINKUSDT",
		},
		constants.COMPOUND_USDT: {
			OffChainTicker: "COMPUSDT",
		},
		constants.CURVE_USDT: {
			OffChainTicker: "CRVUSDT",
		},
		constants.DOGE_USDT: {
			OffChainTicker: "DOGEUSDT",
		},
		constants.DYDX_USDT: {
			OffChainTicker: "DYDXUSDT",
		},
		constants.ETC_USDT: {
			OffChainTicker: "ETCUSDT",
		},
		constants.ETHEREUM_BITCOIN: {
			OffChainTicker: "ETHBTC",
		},
		constants.ETHEREUM_USDC: {
			OffChainTicker: "ETHUSDC",
		},
		constants.ETHEREUM_USDT: {
			OffChainTicker: "ETHUSDT",
		},
		constants.FILECOIN_USDT: {
			OffChainTicker: "FILUSDT",
		},
		constants.LIDO_USDT: {
			OffChainTicker: "LDOUSDT",
		},
		constants.LITECOIN_USDT: {
			OffChainTicker: "LTCUSDT",
		},
		constants.MAKER_USDT: {
			OffChainTicker: "MKRUSDT",
		},
		constants.NEAR_USDT: {
			OffChainTicker: "NEARUSDT",
		},
		constants.OPTIMISM_USDT: {
			OffChainTicker: "OPUSDT",
		},
		constants.PEPE_USDT: {
			OffChainTicker: "PEPEUSDT",
		},
		constants.POLKADOT_USDT: {
			OffChainTicker: "DOTUSDT",
		},
		constants.POLYGON_USDT: {
			OffChainTicker: "MATICUSDT",
		},
		constants.RIPPLE_USDT: {
			OffChainTicker: "XRPUSDT",
		},
		constants.SEI_USDT: {
			OffChainTicker: "SEIUSDT",
		},
		constants.SHIBA_USDT: {
			OffChainTicker: "SHIBUSDT",
		},
		constants.SOLANA_USDC: {
			OffChainTicker: "SOLUSDC",
		},
		constants.SOLANA_USDT: {
			OffChainTicker: "SOLUSDT",
		},
		constants.STELLAR_USDT: {
			OffChainTicker: "XLMUSDT",
		},
		constants.SUI_USDT: {
			OffChainTicker: "SUIUSDT",
		},
		constants.TRON_USDT: {
			OffChainTicker: "TRXUSDT",
		},
		constants.UNISWAP_USDT: {
			OffChainTicker: "UNIUSDT",
		},
		constants.USDC_USDT: {
			OffChainTicker: "USDCUSDT",
		},
		constants.USDT_USD: {
			OffChainTicker: "USDTUSD",
		},
		constants.WORLD_USDT: {
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
