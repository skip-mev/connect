package coinbase

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

// NOTE: All documentation for this file can be located on the Coinbase
// API documentation: https://docs.cloud.coinbase.com/sign-in-with-coinbase/docs/api-prices#get-spot-price. This
// API does not require a subscription to use (i.e. No API key is required).

const (
	// Name is the name of the Coinbase provider.
	Name = "coinbase_api"

	Type = types.ConfigType

	// URL is the base URL of the Coinbase API. This includes the base and quote
	// currency pairs that need to be inserted into the URL.
	URL = "https://api.coinbase.com/v2/prices/%s/spot"
)

var (
	// DefaultAPIConfig is the default configuration for the Coinbase API.
	DefaultAPIConfig = config.APIConfig{
		Name:             Name,
		Atomic:           false,
		Enabled:          true,
		Timeout:          3000 * time.Millisecond,
		Interval:         100 * time.Millisecond,
		ReconnectTimeout: 2000 * time.Millisecond,
		MaxQueries:       5,
		URL:              URL,
	}

	DefaultProviderConfig = config.ProviderConfig{
		Name: Name,
		API:  DefaultAPIConfig,
		Type: Type,
	}

	// DefaultMarketConfig is the default market configuration for Coinbase.
	DefaultMarketConfig = types.CurrencyPairsToProviderTickers{
		constants.APE_USD: {
			OffChainTicker: "APE-USD",
		},
		constants.APE_USDC: {
			OffChainTicker: "APE-USDC",
		},
		constants.APE_USDT: {
			OffChainTicker: "APE-USDT",
		},
		constants.APTOS_USD: {
			OffChainTicker: "APT-USD",
		},
		constants.ARBITRUM_USD: {
			OffChainTicker: "ARB-USD",
		},
		constants.ATOM_USD: {
			OffChainTicker: "ATOM-USD",
		},
		constants.ATOM_USDC: {
			OffChainTicker: "ATOM-USDC",
		},
		constants.ATOM_USDT: {
			OffChainTicker: "ATOM-USDT",
		},
		constants.AVAX_USD: {
			OffChainTicker: "AVAX-USD",
		},
		constants.AVAX_USDC: {
			OffChainTicker: "AVAX-USDC",
		},
		constants.AVAX_USDT: {
			OffChainTicker: "AVAX-USDT",
		},
		constants.BCH_USD: {
			OffChainTicker: "BCH-USD",
		},
		constants.BITCOIN_USD: {
			OffChainTicker: "BTC-USD",
		},
		constants.BITCOIN_USDT: {
			OffChainTicker: "BTC-USDT",
		},
		constants.BITCOIN_USDC: {
			OffChainTicker: "BTC-USDC",
		},
		constants.BITCOIN_USDT: {
			OffChainTicker: "BTC-USDT",
		},
		constants.BLUR_USD: {
			OffChainTicker: "BLUR-USD",
		},
		constants.CARDANO_USD: {
			OffChainTicker: "ADA-USD",
		},
		constants.CELESTIA_USD: {
			OffChainTicker: "TIA-USD",
		},
		constants.CELESTIA_USDC: {
			OffChainTicker: "TIA-USDC",
		},
		constants.CELESTIA_USDT: {
			OffChainTicker: "TIA-USDT",
		},
		constants.CHAINLINK_USD: {
			OffChainTicker: "LINK-USD",
		},
		constants.COMPOUND_USD: {
			OffChainTicker: "COMP-USD",
		},
		constants.CURVE_USD: {
			OffChainTicker: "CRV-USD",
		},
		constants.DOGE_USD: {
			OffChainTicker: "DOGE-USD",
		},
		constants.DYDX_USD: {
			OffChainTicker: "DYDX-USD",
		},
		constants.DYDX_USDC: {
			OffChainTicker: "DYDX-USDC",
		},
		constants.DYDX_USDT: {
			OffChainTicker: "DYDX-USDT",
		},
		constants.ETC_USD: {
			OffChainTicker: "ETC-USD",
		},
		constants.ETHEREUM_BITCOIN: {
			OffChainTicker: "ETH-BTC",
		},
		constants.ETHEREUM_USD: {
			OffChainTicker: "ETH-USD",
		},
		constants.ETHEREUM_USDC: {
			OffChainTicker: "ETH-USDC",
		},
		constants.ETHEREUM_USDT: {
			OffChainTicker: "ETH-USDT",
		},
		constants.FILECOIN_USD: {
			OffChainTicker: "FIL-USD",
		},
		constants.LIDO_USD: {
			OffChainTicker: "LDO-USD",
		},
		constants.LITECOIN_USD: {
			OffChainTicker: "LTC-USD",
		},
		constants.MAKER_USD: {
			OffChainTicker: "MKR-USD",
		},
		constants.NEAR_USD: {
			OffChainTicker: "NEAR-USD",
		},
		constants.OPTIMISM_USD: {
			OffChainTicker: "OP-USD",
		},
		constants.OSMOSIS_USD: {
			OffChainTicker: "OSMO-USD",
		},
		constants.OSMOSIS_USDC: {
			OffChainTicker: "OSMO-USDC",
		},
		constants.OSMOSIS_USDT: {
			OffChainTicker: "OSMO-USDT",
		},
		constants.POLKADOT_USD: {
			OffChainTicker: "DOT-USD",
		},
		constants.POLYGON_USD: {
			OffChainTicker: "MATIC-USD",
		},
		constants.RIPPLE_USD: {
			OffChainTicker: "XRP-USD",
		},
		constants.SEI_USD: {
			OffChainTicker: "SEI-USD",
		},
		constants.SHIBA_USD: {
			OffChainTicker: "SHIB-USD",
		},
		constants.SOLANA_USD: {
			OffChainTicker: "SOL-USD",
		},
		constants.SOLANA_USDC: {
			OffChainTicker: "SOL-USDC",
		},
		constants.SOLANA_USDT: {
			OffChainTicker: "SOL-USDT",
		},
		constants.STELLAR_USD: {
			OffChainTicker: "XLM-USD",
		},
		constants.SUI_USD: {
			OffChainTicker: "SUI-USD",
		},
		constants.UNISWAP_USD: {
			OffChainTicker: "UNI-USD",
		},
		constants.USDC_USD: {
			OffChainTicker: "USDC-USD",
		},
		constants.USDC_USDT: {
			OffChainTicker: "USDC-USDT",
		},
		constants.USDT_USD: {
			OffChainTicker: "USDT-USD",
		},
	}
)

type (
	// CoinBaseResponse is the expected response returned by the Coinbase API.
	// The response is json formatted.
	// Response format:
	//
	//	{
	//	  "data": {
	//	    "amount": "1020.25",
	//	    "currency": "USD"
	//	  }
	//	}
	CoinBaseResponse struct { //nolint
		Data CoinBaseData `json:"data"`
	}

	// CoinBaseData is the data returned by the Coinbase API.
	CoinBaseData struct { //nolint
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	}
)
