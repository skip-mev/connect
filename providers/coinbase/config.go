package coinbase

import "github.com/spf13/viper"

// Config is the configuration for the Coinbase provider.
type Config struct {
	// NameToSymbol is a map of currency names to their symbols.
	NameToSymbol map[string]string `mapstructure:"name_to_symbol" toml:"name_to_symbol"`
}

// ReadCoinbaseConfigFromFile reads a config from a file and returns the config.
func ReadCoinbaseConfigFromFile(path string) (Config, error) {
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
