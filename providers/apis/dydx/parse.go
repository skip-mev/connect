package dydx

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/skip-mev/connect/v2/oracle/constants"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/apis/bitstamp"
	"github.com/skip-mev/connect/v2/providers/apis/coinmarketcap"
	"github.com/skip-mev/connect/v2/providers/apis/defi/raydium"
	"github.com/skip-mev/connect/v2/providers/apis/defi/uniswapv3"
	dydxtypes "github.com/skip-mev/connect/v2/providers/apis/dydx/types"
	"github.com/skip-mev/connect/v2/providers/apis/kraken"
	"github.com/skip-mev/connect/v2/providers/volatile"
	"github.com/skip-mev/connect/v2/providers/websockets/binance"
	"github.com/skip-mev/connect/v2/providers/websockets/bitfinex"
	"github.com/skip-mev/connect/v2/providers/websockets/bybit"
	"github.com/skip-mev/connect/v2/providers/websockets/coinbase"
	"github.com/skip-mev/connect/v2/providers/websockets/cryptodotcom"
	"github.com/skip-mev/connect/v2/providers/websockets/gate"
	"github.com/skip-mev/connect/v2/providers/websockets/huobi"
	"github.com/skip-mev/connect/v2/providers/websockets/kucoin"
	"github.com/skip-mev/connect/v2/providers/websockets/mexc"
	"github.com/skip-mev/connect/v2/providers/websockets/okx"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

// ProviderMapping is referencing the different providers that are supported by the dYdX market params.
//
// ref: https://github.com/dydxprotocol/v4-chain/blob/main/protocol/daemons/pricefeed/client/constants/exchange_common/exchange_id.go
var ProviderMapping = map[string]string{
	"Binance":              binance.Name,
	"BinanceUS":            binance.Name,
	"Bitfinex":             bitfinex.Name,
	"Kraken":               kraken.Name, // We only support the API since the WebSocket has different pairs.
	"Gate":                 gate.Name,
	"Bitstamp":             bitstamp.Name,
	"Bybit":                bybit.Name,
	"CryptoCom":            cryptodotcom.Name,
	"Huobi":                huobi.Name,
	"Kucoin":               kucoin.Name,
	"Okx":                  okx.Name,
	"Mexc":                 mexc.Name,
	"CoinbasePro":          coinbase.Name,
	"TestVolatileExchange": volatile.Name,
	"Raydium":              raydium.Name,
	"UniswapV3-Ethereum":   uniswapv3.ProviderNames[constants.ETHEREUM],
	"UniswapV3-Base":       uniswapv3.ProviderNames[constants.BASE],
	coinmarketcap.Name:     coinmarketcap.Name,
}

// ConvertMarketParamsToMarketMap converts a dYdX market params response to a connect market map response.
func ConvertMarketParamsToMarketMap(
	params dydxtypes.QueryAllMarketParamsResponse,
) (mmtypes.MarketMapResponse, error) {
	marketMap := mmtypes.MarketMap{
		Markets: make(map[string]mmtypes.Market),
	}

	for _, market := range params.MarketParams {
		ticker, err := CreateTickerFromMarket(market)
		if err != nil {
			return mmtypes.MarketMapResponse{}, fmt.Errorf("failed to create ticker from market %s: %w", market.Pair, err)
		}

		var exchangeConfigJSON dydxtypes.ExchangeConfigJson
		if err := json.Unmarshal([]byte(market.ExchangeConfigJson), &exchangeConfigJSON); err != nil {
			return mmtypes.MarketMapResponse{}, fmt.Errorf("failed to unmarshal exchange json config for %s: %w", ticker.String(), err)
		}

		// Convert the exchange config JSON to a set of paths and providers.
		providers, err := ConvertExchangeConfigJSON(exchangeConfigJSON)
		if err != nil {
			return mmtypes.MarketMapResponse{}, fmt.Errorf("failed to convert exchange config json for %s: %w", ticker.String(), err)
		}

		marketMap.Markets[ticker.String()] = mmtypes.Market{
			Ticker:          ticker,
			ProviderConfigs: providers,
		}
	}

	return mmtypes.MarketMapResponse{
		MarketMap: marketMap,
	}, nil
}

// CreateTickerFromMarket creates a ticker from a dYdX market.
func CreateTickerFromMarket(market dydxtypes.MarketParam) (mmtypes.Ticker, error) {
	cp, err := CreateCurrencyPairFromPair(market.Pair)
	if err != nil {
		return mmtypes.Ticker{}, err
	}

	t := mmtypes.Ticker{
		CurrencyPair:     cp,
		Decimals:         uint64(market.Exponent * -1), //nolint:gosec
		MinProviderCount: uint64(market.MinExchanges),
		Enabled:          true,
	}

	return t, t.ValidateBasic()
}

// CreateCurrencyPairFromPair creates a currency pair from a dYdX market.
func CreateCurrencyPairFromPair(pair string) (connecttypes.CurrencyPair, error) {
	split := strings.Split(pair, Delimiter)
	if len(split) != 2 {
		return connecttypes.CurrencyPair{}, fmt.Errorf("expected pair (%s) to have 2 elements, got %d", pair, len(split))
	}

	cp := connecttypes.NewCurrencyPair(
		strings.ToUpper(split[0]), // Base
		strings.ToUpper(split[1]), // Quote
	)

	return cp, cp.ValidateBasic()
}

// ConvertExchangeConfigJSON creates a set of paths and providers for a given ticker
// from a dYdX market. These paths represent the different ways to convert a currency
// pair using the dYdX market.
func ConvertExchangeConfigJSON(
	config dydxtypes.ExchangeConfigJson,
) ([]mmtypes.ProviderConfig, error) {
	var (
		providers = make([]mmtypes.ProviderConfig, 0, len(config.Exchanges))
		seen      = make(map[dydxtypes.ExchangeMarketConfigJson]struct{})
	)

	for _, cfg := range config.Exchanges {
		// Ignore duplicates.
		if _, ok := seen[cfg]; ok {
			continue
		}
		seen[cfg] = struct{}{}

		// This means we have seen an exchange that connect cannot support.
		exchange, ok := ProviderMapping[cfg.ExchangeName]
		if !ok {
			continue
		}

		// Determine if the exchange needs to have an normalizeByPair.
		var normalizeByPair *connecttypes.CurrencyPair
		if len(cfg.AdjustByMarket) > 0 {
			temp, err := CreateCurrencyPairFromPair(cfg.AdjustByMarket)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to create normalize by pair for %s: %w",
					cfg.AdjustByMarket,
					err,
				)
			}

			normalizeByPair = &temp
		}

		// Convert the ticker to the provider's format.
		denom, err := ConvertDenomByProvider(cfg.Ticker, exchange)
		if err != nil {
			return nil, fmt.Errorf("failed to convert denom by provider: %w", err)
		}

		metaData, err := ExtractMetadata(exchange, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to extract metadata: %w", err)
		}

		// Convert to a provider config.
		providers = append(providers, mmtypes.ProviderConfig{
			Name:            exchange,
			OffChainTicker:  denom,
			Invert:          cfg.Invert,
			NormalizeByPair: normalizeByPair,
			Metadata_JSON:   metaData,
		})
	}

	return providers, nil
}

// ExtractMetadata extracts Metadata_JSON from ExchangeMarketConfigJson, based on the converted provider name.
func ExtractMetadata(providerName string, cfg dydxtypes.ExchangeMarketConfigJson) (string, error) {
	// Exchange-specific logic for converting a ticker to provider-specific metadata json
	switch {
	case strings.HasPrefix(providerName, uniswapv3.BaseName):
		return UniswapV3MetadataFromTicker(cfg.Ticker, cfg.Invert)
	case providerName == raydium.Name:
		return RaydiumMetadataFromTicker(cfg.Ticker)
	}
	return "", nil
}

// ConvertDenomByProvider converts a given denom to a format that is compatible with a given provider.
// Specifically, this is used to convert API to WebSocket representations of denoms where necessary.
func ConvertDenomByProvider(denom string, exchange string) (string, error) {
	switch {
	case exchange == mexc.Name:
		if strings.Contains(denom, "_") {
			return strings.ReplaceAll(denom, "_", ""), nil
		}

		return denom, nil
	case exchange == raydium.Name:
		// split the ticker by /, and expect there to at least be two values
		fields := strings.Split(denom, RaydiumTickerSeparator)
		if len(fields) < 2 {
			return "", fmt.Errorf("expected denom to have at least 2 fields, got %d for %s ticker: %s", len(fields), exchange, denom)
		}

		return connecttypes.NewCurrencyPair(fields[0], fields[1]).String(), nil
	default:
		return denom, nil
	}
}
