package types

type (
	// Provider defines an interface an exchange price provider must implement.
	Provider interface {
		// Name returns the name of the provider.
		Name() string

		// GetPrices returns the tickerPrices based on the provided pairs.
		GetPrices() (map[string]TickerPrice, error)

		// SetPairs sets the pairs that the provider should fetch prices for.
		SetPairs(...CurrencyPair)

		// GetPairs returns the pairs that the provider is fetching prices for.
		GetPairs() []CurrencyPair
	}
)
