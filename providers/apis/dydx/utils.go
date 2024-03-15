package dydx

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	dydxtypes "github.com/skip-mev/slinky/providers/apis/dydx/types"
	mmclienttypes "github.com/skip-mev/slinky/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

const (
	// Name is the name of the MarketMap provider.
	Name = "dYdX"

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
	Timeout:          1000 * time.Millisecond,
	Interval:         10 * time.Second,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	URL:              "replace-me",
	Type:             mmclienttypes.ConfigType,
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

		paths, providers, err := ConvertExchangeConfigJSON(ticker, exchangeConfigJSON)
		if err != nil {
			return mmtypes.GetMarketMapResponse{}, err
		}

		// Add the ticker, provider, and paths to the market map.
		marketMap.Tickers[ticker.String()] = ticker
		marketMap.Paths[ticker.String()] = paths
		marketMap.Providers[ticker.String()] = providers
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
// pair using the dYdX market. This involves either a direct or indirect conversion
// (via an index price).
func ConvertExchangeConfigJSON(
	ticker mmtypes.Ticker,
	config dydxtypes.ExchangeConfigJson,
) (mmtypes.Paths, mmtypes.Providers, error) {
	paths := make([]mmtypes.Path, 0)
	providers := make([]mmtypes.ProviderConfig, 0)

	for _, cfg := range config.Exchanges {
		var operations []mmtypes.Operation
		switch {
		case len(cfg.AdjustByMarket) == 0 && !cfg.Invert:
			// Direct conversion.
			operation := mmtypes.Operation{
				CurrencyPair: ticker.CurrencyPair,
				Provider:     cfg.ExchangeName,
				Invert:       false,
			}
			providerCfg := mmtypes.ProviderConfig{
				Name:           cfg.ExchangeName,
				OffChainTicker: cfg.Ticker,
			}

			operations = append(operations, operation)
			providers = append(providers, providerCfg)
		case len(cfg.AdjustByMarket) == 0 && cfg.Invert:
			// Direct conversion with inverted price.
			operation := mmtypes.Operation{
				CurrencyPair: ticker.CurrencyPair,
				Provider:     cfg.ExchangeName,
				Invert:       true,
			}

			providerCfg := mmtypes.ProviderConfig{
				Name:           cfg.ExchangeName,
				OffChainTicker: cfg.Ticker,
			}

			operations = append(operations, operation)
			providers = append(providers, providerCfg)
		case len(cfg.AdjustByMarket) > 0 && !cfg.Invert:
			// Indirect conversion with index price. This is effectively a conversion
			// from the base currency to the quote currency via the index currency.
			// Ex. BTC/USD via BTC/USDT and USDT/USD.
			first := mmtypes.Operation{
				CurrencyPair: ticker.CurrencyPair,
				Provider:     cfg.ExchangeName,
				Invert:       false,
			}

			cp, err := CreateCurrencyPairFromPair(cfg.AdjustByMarket)
			if err != nil {
				return mmtypes.Paths{}, mmtypes.Providers{}, err
			}
			second := mmtypes.Operation{
				CurrencyPair: cp,
				Provider:     mmtypes.IndexPrice,
				Invert:       false,
			}

			providerCfg := mmtypes.ProviderConfig{
				Name:           cfg.ExchangeName,
				OffChainTicker: cfg.Ticker,
			}

			operations = append(operations, first, second)
			providers = append(providers, providerCfg)
		case len(cfg.AdjustByMarket) > 0 && cfg.Invert:
			// Indirect conversion with inverted index price. This is effectively a conversion
			// from the base currency to the quote currency via the index currency with an inverted
			// price. Ex. USDT/USD via BTC/USDT and BTC/USD. In this case, we are not defining
			// a new market but are instead using an existing one. The existing market must match
			// to the market used in the index price.
			cp, err := CreateCurrencyPairFromPair(cfg.AdjustByMarket)
			if err != nil {
				return mmtypes.Paths{}, mmtypes.Providers{}, err
			}
			first := mmtypes.Operation{
				CurrencyPair: cp,
				Provider:     cfg.ExchangeName,
				Invert:       true,
			}

			second := mmtypes.Operation{
				CurrencyPair: cp,
				Provider:     mmtypes.IndexPrice,
				Invert:       false,
			}

			operations = append(operations, first, second)
		}

		// Add the provider config and operations to the paths and providers.
		paths = append(paths, mmtypes.Path{Operations: operations})
	}

	allPaths := mmtypes.Paths{
		Paths: paths,
	}
	allProviders := mmtypes.Providers{
		Providers: providers,
	}

	return allPaths, allProviders, nil
}
