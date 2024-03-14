package dydx

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/pkg/math/oracle"
	pkgtypes "github.com/skip-mev/slinky/pkg/types"
	dydxtypes "github.com/skip-mev/slinky/providers/apis/dydx/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

const (
	// Name is the name of the MarketMap provider.
	Name = "dydx"

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
	Timeout:          500 * time.Millisecond,
	Interval:         5 * time.Second,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	URL:              "https://dydx-api.polkachu.com", // TEMP PLACEHOLDER
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

		var exchangeConfigJson dydxtypes.ExchangeConfigJson
		if err := json.Unmarshal([]byte(market.ExchangeConfigJson), &exchangeConfigJson); err != nil {
			return mmtypes.GetMarketMapResponse{}, fmt.Errorf("failed to unmarshal exchange json config for %s: %w", ticker.String(), err)
		}

		paths, err := CreatePathsFromMarket(ticker, exchangeConfigJson)
		if err != nil {
			return mmtypes.GetMarketMapResponse{}, err
		}

		providers, err := CreateProvidersFromMarket(exchangeConfigJson)
		if err != nil {
			return mmtypes.GetMarketMapResponse{}, err
		}

		// Add the ticker, provider, and paths to the market map.
		marketMap.Tickers[ticker.String()] = ticker
		marketMap.Providers[ticker.String()] = providers
		marketMap.Paths[ticker.String()] = paths
	}

	return mmtypes.GetMarketMapResponse{
		MarketMap: marketMap,
	}, nil
}

// CreateCurrencyPairFromMarket creates a currency pair from a dYdX market.
func CreateCurrencyPairFromMarket(pair string) (pkgtypes.CurrencyPair, error) {
	split := strings.Split(pair, Delimeter)
	if len(split) != 2 {
		return pkgtypes.CurrencyPair{}, fmt.Errorf("expected pair (%s) to have 2 elements, got %d", pair, len(split))
	}

	base := split[0]
	quote := split[1]
	return pkgtypes.NewCurrencyPair(base, quote), nil
}

// CreateTickerFromMarket creates a ticker from a dYdX market.
func CreateTickerFromMarket(market dydxtypes.MarketParam) (mmtypes.Ticker, error) {
	cp, err := CreateCurrencyPairFromMarket(market.Pair)
	if err != nil {
		return mmtypes.Ticker{}, err
	}

	return mmtypes.Ticker{
		CurrencyPair:     cp,
		Decimals:         uint64(market.Exponent * -1),
		MinProviderCount: uint64(market.MinExchanges),
	}, nil
}

// CreatePathsFromMarket creates a set of paths for a given ticker from a dYdX market.
func CreatePathsFromMarket(
	ticker mmtypes.Ticker,
	config dydxtypes.ExchangeConfigJson,
) (mmtypes.Paths, error) {
	paths := make([]mmtypes.Path, 0)

	for _, cfg := range config.Exchanges {
		path := mmtypes.Path{
			Operations: make([]mmtypes.Operation, 0),
		}

		// Add the first operation which is going to be the ticker.
		first := mmtypes.Operation{
			Provider:     cfg.ExchangeName,
			CurrencyPair: ticker.CurrencyPair,
			Invert:       cfg.Invert,
		}
		path.Operations = append(path.Operations, first)

		// Check if we need to convert via a index price.
		if len(cfg.AdjustByMarket) > 0 {
			cp, err := CreateCurrencyPairFromMarket(cfg.AdjustByMarket)
			if err != nil {
				return mmtypes.Paths{}, err
			}

			second := mmtypes.Operation{
				Provider:     oracle.IndexPrice,
				CurrencyPair: cp,
				Invert:       false,
			}
			path.Operations = append(path.Operations, second)
		}
	}

	return mmtypes.Paths{
		Paths: paths,
	}, nil
}

// CreateProvidersFromMarket creates a set of providers for a given ticker from a dYdX market.
func CreateProvidersFromMarket(config dydxtypes.ExchangeConfigJson) (mmtypes.Providers, error) {
	providerConfigs := make([]mmtypes.ProviderConfig, 0)

	for _, cfg := range config.Exchanges {
		// TODO: Is there any additional validation we need to do on the raw exchange config?
		name := cfg.ExchangeName
		offChainTicker := cfg.Ticker

		// Add the provider to the list of provider configs.
		providerConfig := mmtypes.ProviderConfig{
			Name:           name,
			OffChainTicker: offChainTicker,
		}
		providerConfigs = append(providerConfigs, providerConfig)
	}

	return mmtypes.Providers{
		Providers: providerConfigs,
	}, nil
}
