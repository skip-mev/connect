package provider

import (
	"github.com/skip-mev/slinky/oracle/types"
)

type (
	// Provider defines an interface an exchange price provider must implement.
	Provider interface {
		// Name returns the name of the provider.
		Name() string

		// GetTickerPrices returns the tickerPrices based on the provided pairs.
		GetTickerPrices(...types.CurrencyPair) (map[string]types.TickerPrice, error)

		// GetCandlePrices returns the candlePrices based on the provided pairs.
		GetCandlePrices(...types.CurrencyPair) (map[string][]types.Candle, error)

		// GetAvailablePairs return all available pairs symbol to subscribe.
		GetAvailablePairs() (map[string]struct{}, error)
	}

	// Endpoint defines an override setting in our config for the
	// hardcoded rest and websocket api endpoints.
	Endpoint struct {
		// Name of the provider, ex. "binance"
		Name string `toml:"name"`

		// Rest endpoint for the provider, ex. "https://api1.binance.com"
		Rest string `toml:"rest"`

		// Websocket endpoint for the provider, ex. "stream.binance.com:9443"
		Websocket string `toml:"websocket"`

		// APIKey for API Key protected endpoints
		APIKey string `toml:"apikey"`
	}
)
