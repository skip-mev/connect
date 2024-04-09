package dydx

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/apis/binance"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	dydxtypes "github.com/skip-mev/slinky/providers/apis/dydx/types"
	"github.com/skip-mev/slinky/providers/apis/kraken"
	"github.com/skip-mev/slinky/providers/volatile"
	"github.com/skip-mev/slinky/providers/websockets/bitfinex"
	"github.com/skip-mev/slinky/providers/websockets/bitstamp"
	"github.com/skip-mev/slinky/providers/websockets/bybit"
	"github.com/skip-mev/slinky/providers/websockets/cryptodotcom"
	"github.com/skip-mev/slinky/providers/websockets/gate"
	"github.com/skip-mev/slinky/providers/websockets/huobi"
	"github.com/skip-mev/slinky/providers/websockets/kucoin"
	"github.com/skip-mev/slinky/providers/websockets/mexc"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
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
}

// ConvertMarketParamsToMarketMap converts a dYdX market params response to a slinky market map response.
func (h *APIHandler) ConvertMarketParamsToMarketMap(
	params dydxtypes.QueryAllMarketParamsResponse,
) (mmtypes.MarketMapResponse, error) {
	marketMap := mmtypes.MarketMap{
		Markets: make(map[string]mmtypes.Market),
	}

	for _, market := range params.MarketParams {
		ticker, err := h.CreateTickerFromMarket(market)
		if err != nil {
			h.logger.Error(
				"failed to create ticker from market",
				zap.String("market", market.Pair),
				zap.Error(err),
			)

			return mmtypes.MarketMapResponse{}, fmt.Errorf("failed to create ticker from market: %w", err)
		}

		var exchangeConfigJSON dydxtypes.ExchangeConfigJson
		if err := json.Unmarshal([]byte(market.ExchangeConfigJson), &exchangeConfigJSON); err != nil {
			h.logger.Error(
				"failed to unmarshal exchange json config",
				zap.String("ticker", ticker.String()),
				zap.Error(err),
			)

			return mmtypes.MarketMapResponse{}, fmt.Errorf("failed to unmarshal exchange json config: %w", err)
		}

		// Convert the exchange config JSON to a set of paths and providers.
		providers, err := h.ConvertExchangeConfigJSON(exchangeConfigJSON)
		if err != nil {
			h.logger.Error(
				"failed to convert exchange config json",
				zap.String("ticker", ticker.String()),
				zap.Error(err),
			)

			return mmtypes.MarketMapResponse{}, fmt.Errorf("failed to convert exchange config json: %w", err)
		}

		marketMap.Markets[ticker.String()] = mmtypes.Market{
			Ticker:          ticker,
			ProviderConfigs: providers,
		}
	}

	if err := marketMap.ValidateBasic(); err != nil {
		return mmtypes.MarketMapResponse{}, fmt.Errorf("failed to validate market map: %w", err)
	}

	return mmtypes.MarketMapResponse{
		MarketMap: marketMap,
	}, nil
}

// CreateTickerFromMarket creates a ticker from a dYdX market.
func (h *APIHandler) CreateTickerFromMarket(market dydxtypes.MarketParam) (mmtypes.Ticker, error) {
	cp, err := h.CreateCurrencyPairFromPair(market.Pair)
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

// CreateCurrencyPairFromMarket creates a currency pair from a dYdX market.
func (h *APIHandler) CreateCurrencyPairFromPair(pair string) (slinkytypes.CurrencyPair, error) {
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

// ConvertExchangeConfigJSON creates a set of paths and providers for a given ticker
// from a dYdX market. These paths represent the different ways to convert a currency
// pair using the dYdX market.
func (h *APIHandler) ConvertExchangeConfigJSON(
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

		// This means we have seen an exchange that slinky cannot support.
		exchange, ok := ProviderMapping[cfg.ExchangeName]
		if !ok {
			// ignore unsupported exchanges
			h.logger.Error(
				"skipping unsupported exchange",
				zap.String("exchange", cfg.ExchangeName),
				zap.String("ticker", cfg.Ticker),
			)

			continue
		}

		// Determine if the exchange needs to have an normalizeByPair.
		var normalizeByPair *slinkytypes.CurrencyPair
		if len(cfg.AdjustByMarket) > 0 {
			temp, err := h.CreateCurrencyPairFromPair(cfg.AdjustByMarket)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to create normalize by pair for %s: %w",
					cfg.AdjustByMarket,
					err,
				)
			}

			normalizeByPair = &temp
		}

		// Convert to a provider config.
		providers = append(providers, mmtypes.ProviderConfig{
			Name:            exchange,
			OffChainTicker:  ConvertDenomByProvider(cfg.Ticker, exchange), // Convert the ticker to the provider's format.
			Invert:          cfg.Invert,
			NormalizeByPair: normalizeByPair,
		})
	}

	return providers, nil
}

// ConvertDenomByProvider converts a given denom to a format that is compatible with a given provider.
// Specifically, this is used to convert API to WebSocket representations of denoms where necessary.
func ConvertDenomByProvider(denom string, exchange string) string {
	switch {
	case exchange == mexc.Name:
		if strings.Contains(denom, "_") {
			return strings.ReplaceAll(denom, "_", "")
		}

		return denom
	case exchange == bitstamp.Name:
		if strings.Contains(denom, "/") {
			return strings.ToLower(strings.ReplaceAll(denom, "/", ""))
		}

		return strings.ToLower(denom)
	default:
		return denom
	}
}
