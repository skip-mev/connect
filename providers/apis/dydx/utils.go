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
	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
)

const (
	// Name is the name of the MarketMap provider.
	Name = "dydx_api"

	// ChainID is the chain ID for the dYdX market map provider.
	ChainID = "dydx-mainnet-1"

	// Endpoint is the endpoint for the dYdX market map API.
	Endpoint = "%s/dydxprotocol/prices/params/market?limit=10000"

	// Delimeter is the delimeter used to separate the base and quote assets in a pair.
	Delimeter = "-"
)

// DefaultAPIConfig returns the default configuration for the dYdX market map API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           true,
	Enabled:          true,
	Timeout:          20 * time.Second, // Set a high timeout to account for slow API responses in the case where many markets are queried.
	Interval:         10 * time.Second,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	URL:              "localhost:1317",
}

// ProviderMapping is referencing the different providers that are supported by the dYdX market params.
//
// ref: https://github.com/dydxprotocol/v4-chain/blob/main/protocol/daemons/pricefeed/client/constants/exchange_common/exchange_id.go
var ProviderMapping = map[string][]string{
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
func ConvertMarketParamsToMarketMap(params dydxtypes.QueryAllMarketParamsResponse) (mmtypes.MarketMapResponse, error) {
	marketMap := mmtypes.MarketMap{
		Markets: make(map[string]mmtypes.Market),
	}

	for _, market := range params.MarketParams {
		ticker, err := CreateTickerFromMarket(market)
		if err != nil {
			return mmtypes.MarketMapResponse{}, err
		}

		var exchangeConfigJSON dydxtypes.ExchangeConfigJson
		if err := json.Unmarshal([]byte(market.ExchangeConfigJson), &exchangeConfigJSON); err != nil {
			return mmtypes.MarketMapResponse{}, fmt.Errorf("failed to unmarshal exchange json config for %s: %w", ticker.String(), err)
		}

		// Convert the exchange config JSON to a set of paths and providers.
		tickerProviders, err := ConvertExchangeConfigJSON(exchangeConfigJSON)
		if err != nil {
			return mmtypes.MarketMapResponse{}, fmt.Errorf("failed to convert exchange config json for %s: %w", ticker.String(), err)
		}

		// Add the ticker, provider, and paths to the market map.
		marketMap.Markets[ticker.String()] = mmtypes.Market{
			Ticker:          ticker,
			ProviderConfigs: tickerProviders,
		}
	}

	if err := marketMap.ValidateBasic(); err != nil {
		return mmtypes.MarketMapResponse{}, fmt.Errorf("failed to validate market map: %w", err)
	}

	return mmtypes.MarketMapResponse{
		MarketMap: marketMap,
	}, nil
}

// CreateCurrencyPairFromPair creates a currency pair from a dYdX market.
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
	config dydxtypes.ExchangeConfigJson,
) ([]mmtypes.ProviderConfig, error) {
	var (
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
		exchangeNames, ok := ProviderMapping[cfg.ExchangeName]
		if !ok {
			return nil, fmt.Errorf("unsupported exchange: %s", cfg.ExchangeName)
		}

		var (
			err              error
			createdProviders []mmtypes.ProviderConfig
		)

		// Determine the relevant operations and provider configs based on the exchange config.
		switch {
		case len(cfg.AdjustByMarket) == 0 && !cfg.Invert:
			createdProviders, err = DirectConversion(cfg, exchangeNames)
		case len(cfg.AdjustByMarket) == 0 && cfg.Invert:
			createdProviders, err = InvertedConversion(cfg, exchangeNames)

		case len(cfg.AdjustByMarket) > 0 && !cfg.Invert:
			createdProviders, err = IndirectConversion(cfg, exchangeNames)

		case len(cfg.AdjustByMarket) > 0 && cfg.Invert:
			createdProviders, err = IndirectInvertedConversion(cfg, exchangeNames)
		}
		if err != nil {
			return nil, err
		}

		providers = append(providers, createdProviders...)

	}

	return providers, nil
}

// DirectConversion is a conversion from market to desired ticker i.e. BTC/USD -> BTC/USD.
func DirectConversion(
	cfg dydxtypes.ExchangeMarketConfigJson,
	exchangeNames []string,
) ([]mmtypes.ProviderConfig, error) {
	providers := make([]mmtypes.ProviderConfig, len(exchangeNames))

	offChainTicker := ConvertDenomByProvider(cfg.Ticker, cfg.ExchangeName)
	for i, name := range exchangeNames {
		providers[i] = mmtypes.ProviderConfig{
			Name:            name,
			OffChainTicker:  offChainTicker,
			Invert:          false,
			NormalizeByPair: nil,
		}
	}
	return providers, nil
}

// InvertedConversion is a conversion with an inverted price i.e. USD/BTC ^ -1 = BTC/USD.
func InvertedConversion(
	cfg dydxtypes.ExchangeMarketConfigJson,
	exchangeNames []string,
) ([]mmtypes.ProviderConfig, error) {
	providers := make([]mmtypes.ProviderConfig, len(exchangeNames))

	offChainTicker := ConvertDenomByProvider(cfg.Ticker, cfg.ExchangeName)
	for i, name := range exchangeNames {
		providers[i] = mmtypes.ProviderConfig{
			Name:            name,
			OffChainTicker:  offChainTicker,
			Invert:          true,
			NormalizeByPair: nil,
		}
	}

	return providers, nil
}

// IndirectConversion is a conversion of two markets i.e. BTC/USDT * USDT/USD = BTC/USD.
func IndirectConversion(
	cfg dydxtypes.ExchangeMarketConfigJson,
	exchangeNames []string,
) ([]mmtypes.ProviderConfig, error) {
	providers := make([]mmtypes.ProviderConfig, len(exchangeNames))

	cp, err := CreateCurrencyPairFromPair(cfg.AdjustByMarket)
	if err != nil {
		return nil, err
	}

	offChainTicker := ConvertDenomByProvider(cfg.Ticker, cfg.ExchangeName)
	for i, name := range exchangeNames {
		providers[i] = mmtypes.ProviderConfig{
			Name:            name,
			OffChainTicker:  offChainTicker,
			Invert:          false,
			NormalizeByPair: &cp,
		}
	}

	return providers, nil
}

// IndirectInvertedConversion is a conversion of two markets to a desired ticker
// where the inverted quote of the first market and quote of the second market are used.
// i.e. BTC/USDT ^ -1 * BTC/USD = USDT/USD.
func IndirectInvertedConversion(
	cfg dydxtypes.ExchangeMarketConfigJson,
	exchangeNames []string,
) ([]mmtypes.ProviderConfig, error) {
	providers := make([]mmtypes.ProviderConfig, len(exchangeNames))

	cp, err := CreateCurrencyPairFromPair(cfg.AdjustByMarket)
	if err != nil {
		return nil, err
	}

	offChainTicker := ConvertDenomByProvider(cfg.Ticker, cfg.ExchangeName)
	for i, name := range exchangeNames {
		providers[i] = mmtypes.ProviderConfig{
			Name:            name,
			OffChainTicker:  offChainTicker,
			Invert:          true,
			NormalizeByPair: &cp,
		}
	}

	return providers, nil
}

// ConvertDenomByProvider converts a given denom to a format that is compatible with a given provider.
// Specifically, this is used to convert API to WebSocket representations of denoms where necessary.
func ConvertDenomByProvider(denom string, exchange string) string {
	providers, ok := ProviderMapping[exchange]
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
