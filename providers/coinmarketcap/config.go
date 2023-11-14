package coinmarketcap

import "github.com/spf13/viper"

// Config is the config struct for the coinmarketcap provider.
type Config struct {
	// APIKey is the API key used to make requests to the coinmarketcap API.
	APIKey string `mapstructure:"api_key" toml:"api_key"`
	// TokenNameToMetadata is a map of token names to their metadata.
	TokenNameToSymbol map[string]string `mapstructure:"token_name_to_symbol" toml:"token_name_to_symbol"`
}

// ReadCoinMarketCapConfigFromFile reads a config from a file and returns the config.
func ReadCoinMarketCapConfigFromFile(path string) (Config, error) {
	// read in config file
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	// unmarshal config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
