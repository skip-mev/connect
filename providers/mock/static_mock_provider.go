package mock

import (
	"context"
	"fmt"
	"math/big"

	"github.com/spf13/viper"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var _ oracle.Provider = (*StaticMockProvider)(nil)

type (
	// StaticMockProvider defines a mocked exchange rate provider using fixed exchange
	// rates.
	StaticMockProvider struct {
		exchangeRates map[oracletypes.CurrencyPair]aggregator.QuotePrice
		currencyPairs []oracletypes.CurrencyPair
	}

	// StaticMockProviderConfig is a map of token names to their metadata.
	StaticMockProviderConfig struct {
		// TokenPrices is a map of token names to their metadata.
		TokenPrices map[string]string `mapstructure:"tokens" toml:"tokens"`
	}
)

// NewStaticMockProvider returns a new mock provider. The mock provider
// will always return the same static data. Meant to be used for testing.
func NewStaticMockProvider() *StaticMockProvider {
	return &StaticMockProvider{
		exchangeRates: map[oracletypes.CurrencyPair]aggregator.QuotePrice{
			oracletypes.NewCurrencyPair("COSMOS", "USDC"):   {Price: big.NewInt(1134)},
			oracletypes.NewCurrencyPair("COSMOS", "USDT"):   {Price: big.NewInt(1135)},
			oracletypes.NewCurrencyPair("COSMOS", "USD"):    {Price: big.NewInt(1136)},
			oracletypes.NewCurrencyPair("OSMOSIS", "USDC"):  {Price: big.NewInt(1137)},
			oracletypes.NewCurrencyPair("OSMOSIS", "USDT"):  {Price: big.NewInt(1138)},
			oracletypes.NewCurrencyPair("OSMOSIS", "USD"):   {Price: big.NewInt(1139)},
			oracletypes.NewCurrencyPair("ETHEREUM", "USDC"): {Price: big.NewInt(1140)},
			oracletypes.NewCurrencyPair("ETHEREUM", "USDT"): {Price: big.NewInt(1141)},
			oracletypes.NewCurrencyPair("ETHEREUM", "USD"):  {Price: big.NewInt(1142)},
			oracletypes.NewCurrencyPair("BITCOIN", "USD"):   {Price: big.NewInt(1143)},
		},
		currencyPairs: []oracletypes.CurrencyPair{
			oracletypes.NewCurrencyPair("COSMOS", "USDC"),
			oracletypes.NewCurrencyPair("COSMOS", "USDT"),
			oracletypes.NewCurrencyPair("COSMOS", "USD"),
			oracletypes.NewCurrencyPair("OSMOSIS", "USDC"),
			oracletypes.NewCurrencyPair("OSMOSIS", "USDT"),
			oracletypes.NewCurrencyPair("OSMOSIS", "USD"),
			oracletypes.NewCurrencyPair("ETHEREUM", "USDC"),
			oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
			oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			oracletypes.NewCurrencyPair("BITCOIN", "USD"),
		},
	}
}

// NewStaticMockProviderFromConfig constructs a new static mock provider from the config
// Notice this method expects the TokenNameToSymbol map to be populated w/ entries of the form
// CurrencyPair.ToString(): big.NewInt(price)
func NewStaticMockProviderFromConfig(providerConfig config.ProviderConfig) (*StaticMockProvider, error) {
	if providerConfig.Name != "static-mock-provider" {
		return nil, fmt.Errorf("expected provider config name to be static-mock-provider, got %s", providerConfig.Name)
	}

	config, err := ReadStaticMockProviderConfigFromFile(providerConfig.Path)
	if err != nil {
		return nil, err
	}

	s := StaticMockProvider{
		exchangeRates: make(map[oracletypes.CurrencyPair]aggregator.QuotePrice),
		currencyPairs: make([]oracletypes.CurrencyPair, 0),
	}

	for cpString, price := range config.TokenPrices {
		cp, err := oracletypes.CurrencyPairFromString(cpString)
		if err != nil {
			continue
		}

		price, converted := big.NewInt(0).SetString(price, 10)
		if !converted {
			return nil, fmt.Errorf("failed to parse price %s for currency pair %s", price, cpString)
		}

		s.exchangeRates[cp] = aggregator.QuotePrice{Price: price}
		s.currencyPairs = append(s.currencyPairs, cp)
	}

	return &s, nil
}

// Name returns the name of the mock provider.
func (p StaticMockProvider) Name() string {
	return "static-mock-provider"
}

// GetPrices returns the mocked exchange rates.
func (p StaticMockProvider) GetPrices(_ context.Context) (map[oracletypes.CurrencyPair]aggregator.QuotePrice, error) {
	return p.exchangeRates, nil
}

// SetPairs is a no-op for the mock provider.
func (p StaticMockProvider) SetPairs(_ ...oracletypes.CurrencyPair) {}

// GetPairs is a no-op for the mock provider.
func (p StaticMockProvider) GetPairs() []oracletypes.CurrencyPair {
	return p.currencyPairs
}

// ReadStaticMockProviderConfigFromFile reads the static mock provider config from the given file.
func ReadStaticMockProviderConfigFromFile(path string) (StaticMockProviderConfig, error) {
	// read in the config file
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return StaticMockProviderConfig{}, err
	}

	// parse config
	var config StaticMockProviderConfig
	if err := viper.Unmarshal(&config); err != nil {
		return StaticMockProviderConfig{}, err
	}

	return config, nil
}
