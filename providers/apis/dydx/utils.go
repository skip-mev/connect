package dydx

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/apis/binance"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/coinbase"
	dydxtypes "github.com/skip-mev/slinky/providers/apis/dydx/types"
	"github.com/skip-mev/slinky/providers/apis/kraken"
	"github.com/skip-mev/slinky/providers/websockets/bitfinex"
	"github.com/skip-mev/slinky/providers/websockets/bitstamp"
	"github.com/skip-mev/slinky/providers/websockets/bybit"
	coinbasews "github.com/skip-mev/slinky/providers/websockets/coinbase"
	"github.com/skip-mev/slinky/providers/websockets/cryptodotcom"
	"github.com/skip-mev/slinky/providers/websockets/gate"
	"github.com/skip-mev/slinky/providers/websockets/huobi"
	"github.com/skip-mev/slinky/providers/websockets/kucoin"
	"github.com/skip-mev/slinky/providers/websockets/mexc"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

const (
	// Name is the name of the MarketMap provider.
	Name = "dydx_api"

	// ChainID is the chain ID for the dYdX market map provider.
	ChainID = "dydx-mainnet-1"

	// Endpoint is the endpoint for the dYdX market map API.
	Endpoint = "%s/dydxprotocol/prices/params/market"

	// Delimeter is the delimeter used to separate the base and quote assets in a pair.
	Delimeter = "-"
)

// DefaultAPIConfig returns the default configuration for the dYdX market map API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           true,
	Enabled:          true,
	Timeout:          2000 * time.Millisecond,
	Interval:         10 * time.Second,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	URL:              "https://dydx-api.polkachu.com",
}

// These are references to the different providers that are supported by the dYdX market map.
//
// ref: https://github.com/dydxprotocol/v4-chain/blob/main/protocol/daemons/pricefeed/client/constants/exchange_common/exchange_id.go
var providerMapping = map[string][]string{
	"Binance":     {binance.Name},
	"BinanceUS":   {binance.Name},
	"Bitfinex":    {bitfinex.Name},
	"Kraken":      {kraken.Name}, // We only support the API since the WebSocket has different pairs.
	"Gate":        {gate.Name},
	"Bitstamp":    {bitstamp.Name},
	"Bybit":       {bybit.Name},
	"CryptoCom":   {cryptodotcom.Name},
	"Huobi":       {huobi.Name},
	"Kucoin":      {kucoin.Name},
	"Okx":         {okx.Name},
	"Mexc":        {mexc.Name},
	"CoinbasePro": {coinbaseapi.Name, coinbasews.Name}, // We support both the API and WebSocket.
}

// ConvertMarketParamsToMarketMap converts a dYdX market params response to a slinky market map response.
func ConvertMarketParamsToMarketMap(params dydxtypes.QueryAllMarketParamsResponse) (mmtypes.GetMarketMapResponse, error) {
	marketMap := mmtypes.MarketMap{
		Tickers:         make(map[string]mmtypes.Ticker),
		Providers:       make(map[string]mmtypes.Providers),
		Paths:           make(map[string]mmtypes.Paths),
		AggregationType: mmtypes.AggregationType_INDEX_PRICE_AGGREGATION,
	}

	for _, market := range params.MarketParams {
		ticker, err := CreateTickerFromMarket(market)
		if err != nil {
			return mmtypes.GetMarketMapResponse{}, err
		}

		var exchangeConfigJSON dydxtypes.ExchangeConfigJson
		if err := json.Unmarshal([]byte(market.ExchangeConfigJson), &exchangeConfigJSON); err != nil {
			return mmtypes.GetMarketMapResponse{}, fmt.Errorf("failed to unmarshal exchange json config for %s: %w", ticker.String(), err)
		}

		// Convert the exchange config JSON to a set of paths and providers.
		tickerPaths, tickerProviders := ConvertExchangeConfigJSON(ticker, exchangeConfigJSON)
		if len(tickerPaths.Paths) == 0 || len(tickerProviders.Providers) == 0 {
			continue
		}

		// Add the ticker, provider, and paths to the market map.
		marketMap.Tickers[ticker.String()] = ticker
		marketMap.Paths[ticker.String()] = tickerPaths
		marketMap.Providers[ticker.String()] = tickerProviders
	}

	if err := marketMap.ValidateBasic(); err != nil {
		return mmtypes.GetMarketMapResponse{}, fmt.Errorf("failed to validate market map: %w", err)
	}

	return mmtypes.GetMarketMapResponse{
		MarketMap: marketMap,
	}, nil
}

// CreateCurrencyPairFromMarket creates a currency pair from a dYdX market.
func CreateCurrencyPairFromPair(pair string) (slinkytypes.CurrencyPair, error) {
	split := strings.Split(pair, Delimeter)
	if len(split) != 2 {
		return slinkytypes.CurrencyPair{}, fmt.Errorf("expected pair (%s) to have 2 elements, got %d", pair, len(split))
	}

	cp := slinkytypes.NewCurrencyPair(
		strings.ToUpper(split[0]), // Base
		strings.ToUpper(split[1]), // Quote
	)

	return cp, cp.ValidateBasic()
}

// CreateTickerFromMarket creates a ticker from a dYdX market.
func CreateTickerFromMarket(market dydxtypes.MarketParam) (mmtypes.Ticker, error) {
	cp, err := CreateCurrencyPairFromPair(market.Pair)
	if err != nil {
		return mmtypes.Ticker{}, err
	}

	t := mmtypes.Ticker{
		CurrencyPair:     cp,
		Decimals:         uint64(market.Exponent * -1),
		MinProviderCount: uint64(market.MinExchanges),
	}

	return t, t.ValidateBasic()
}

// ConvertExchangeConfigJSON creates a set of paths and providers for a given ticker
// from a dYdX market. These paths represent the different ways to convert a currency
// pair using the dYdX market.
func ConvertExchangeConfigJSON(
	ticker mmtypes.Ticker,
	config dydxtypes.ExchangeConfigJson,
) (mmtypes.Paths, mmtypes.Providers) {
	var (
		paths     []mmtypes.Path
		providers []mmtypes.ProviderConfig
		seen      = make(map[dydxtypes.ExchangeMarketConfigJson]struct{})
	)

	for _, cfg := range config.Exchanges {
		// Ignore duplicates.
		if _, ok := seen[cfg]; ok {
			continue
		}
		seen[cfg] = struct{}{}

		// This means we have seen an exchange that slinky cannot support.
		exchangeNames, ok := providerMapping[cfg.ExchangeName]
		if !ok {
			continue
		}

		var (
			exchangePaths []mmtypes.Path
			err           error
			addProviders  = true
		)
		// Determine the relevant operations and provider configs based on the exchange config.
		switch {
		case len(cfg.AdjustByMarket) == 0 && !cfg.Invert:
			exchangePaths = DirectConversion(ticker, exchangeNames)
		case len(cfg.AdjustByMarket) == 0 && cfg.Invert:
			exchangePaths = InvertedConversion(ticker, exchangeNames)
		case len(cfg.AdjustByMarket) > 0 && !cfg.Invert:
			exchangePaths, err = IndirectConversion(ticker, cfg, exchangeNames)
		case len(cfg.AdjustByMarket) > 0 && cfg.Invert:
			exchangePaths, err = IndirectInvertedConversion(cfg, exchangeNames)
			addProviders = false
		}

		// We passively ignore errors here, as we don't want to fail the entire conversion.
		if err != nil {
			continue
		}
		paths = append(paths, exchangePaths...)

		// We only update the providers for a given ticker if the conversion includes the exchanges
		// off-chain representation i.e. Case 1,2,3.
		if addProviders {
			offChainTicker := ConvertDenomByProvider(cfg.Ticker, cfg.ExchangeName)
			for _, name := range exchangeNames {
				providers = append(providers, mmtypes.ProviderConfig{
					Name:           name,
					OffChainTicker: offChainTicker,
				})
			}
		}

	}

	return mmtypes.Paths{Paths: paths}, mmtypes.Providers{Providers: providers}
}

// DirectConversion is a conversion from market to desired ticker i.e. BTC/USD -> BTC/USD.
func DirectConversion(
	ticker mmtypes.Ticker,
	exchangeNames []string,
) []mmtypes.Path {
	paths := make([]mmtypes.Path, len(exchangeNames))
	for i, name := range exchangeNames {
		path := mmtypes.Path{
			Operations: []mmtypes.Operation{
				{
					CurrencyPair: ticker.CurrencyPair,
					Provider:     name,
					Invert:       false,
				},
			},
		}
		paths[i] = path
	}
	return paths
}

// InvertedConversion is a conversion with an inverted price i.e. USD/BTC ^ -1 = BTC/USD.
func InvertedConversion(
	ticker mmtypes.Ticker,
	exchangeNames []string,
) []mmtypes.Path {
	paths := make([]mmtypes.Path, len(exchangeNames))
	for i, name := range exchangeNames {
		path := mmtypes.Path{
			Operations: []mmtypes.Operation{
				{
					CurrencyPair: ticker.CurrencyPair,
					Provider:     name,
					Invert:       true,
				},
			},
		}
		paths[i] = path
	}
	return paths
}

// IndirectConversion is a conversion of two markets i.e. BTC/USDT * USDT/USD = BTC/USD.
func IndirectConversion(
	ticker mmtypes.Ticker,
	cfg dydxtypes.ExchangeMarketConfigJson,
	exchangeNames []string,
) ([]mmtypes.Path, error) {
	cp, err := CreateCurrencyPairFromPair(cfg.AdjustByMarket)
	if err != nil {
		return nil, err
	}

	paths := make([]mmtypes.Path, len(exchangeNames))
	for i, name := range exchangeNames {
		path := mmtypes.Path{
			Operations: []mmtypes.Operation{
				{
					CurrencyPair: ticker.CurrencyPair,
					Provider:     name,
					Invert:       false,
				},
				{
					CurrencyPair: cp,
					Provider:     mmtypes.IndexPrice,
					Invert:       false,
				},
			},
		}
		paths[i] = path
	}

	return paths, nil
}

// IndirectInvertedConversion is a conversion of two markets to a desired ticker
// where the inverted quote of the first market and quote of the second market are used.
// i.e. BTC/USDT ^ -1 * BTC/USD = USDT/USD.
func IndirectInvertedConversion(
	cfg dydxtypes.ExchangeMarketConfigJson,
	exchangeNames []string,
) ([]mmtypes.Path, error) {
	cp, err := CreateCurrencyPairFromPair(cfg.AdjustByMarket)
	if err != nil {
		return nil, err
	}

	paths := make([]mmtypes.Path, len(exchangeNames))
	for i, name := range exchangeNames {
		path := mmtypes.Path{
			Operations: []mmtypes.Operation{
				{
					CurrencyPair: cp,
					Provider:     name,
					Invert:       true,
				},
				{
					CurrencyPair: cp,
					Provider:     mmtypes.IndexPrice,
					Invert:       false,
				},
			},
		}
		paths[i] = path
	}

	return paths, nil
}

// ConvertDenomByProvider converts a given denom to a format that is compatible with a given provider.
// Specifically, this is used to convert API to WebSocket representations of denoms where necessary.
func ConvertDenomByProvider(denom string, exchange string) string {
	providers, ok := providerMapping[exchange]
	if !ok {
		return denom
	}

	switch {
	case len(providers) == 1 && providers[0] == mexc.Name:
		if strings.Contains(denom, "_") {
			return strings.ReplaceAll(denom, "_", "")
		}

		return denom
	case len(providers) == 1 && providers[0] == bitstamp.Name:
		if strings.Contains(denom, "/") {
			return strings.ToLower(strings.ReplaceAll(denom, "/", ""))
		}

		return strings.ToLower(denom)
	default:
		return denom
	}
}
