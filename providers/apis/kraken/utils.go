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
	// Name is the name of the Kraken provider.
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
		Interval:         150 * time.Millisecond,
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

// TickerResult is the result of a Kraken API call for a single ticker. .
//
// https://api.kraken.com/0/public/Ticker
// https://docs.kraken.com/rest/#tag/Market-Data/operation/getTickerInformation
type TickerResult struct {
	pair            string
	AskPriceStats   []string `json:"a" validate:"len=3,dive,positive-float-string"`
	BidPriceStats   []string `json:"b" validate:"len=3,dive,positive-float-string"`
	ClosePriceStats []string `json:"c" validate:"len=2,dive,positive-float-string"`
}

func (ktr *TickerResult) GetAskPrice() string {
	return ktr.AskPriceStats[0]
}

func (ktr *TickerResult) GetBidPrice() string {
	return ktr.BidPriceStats[0]
}

func (ktr *TickerResult) GetLastPrice() string {
	return ktr.ClosePriceStats[0]
}

type ResponseBody struct {
	// As of this time, the Kraken API response is all-or-nothing - either valid ticker data, or one or more errors,
	// but not both. We enforce this expectation by defining mutual exclusivity in the validation tags of the Errors
	// field so that any validated API result always meets our expectation in the response parsing logic.
	Errors  []string                `json:"error" validate:"omitempty"`
	Tickers map[string]TickerResult `validate:"required_without=Errors,excluded_with=Errors,dive" json:"result"`
}

// Decode decodes the given http response into a TickerResult.
func Decode(resp *http.Response) (ResponseBody, error) {
	// Parse the response into a ResponseBody.
	var result ResponseBody
	err := json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
