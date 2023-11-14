package coingecko

import "github.com/spf13/viper"

// Config is the config struct for the CoinGecko provider.
type Config struct {
	// APIKey is the API key used to make requests to the CoinGecko API.
	APIKey string `mapstructure:"api_key" toml:"api_key"`
}

// ReadCoinGeckoConfigFromFile reads a config from a file and returns the config.
func ReadCoinGeckoConfigFromFile(path string) (Config, error) {
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
