package static

import "github.com/spf13/viper"

// StaticMockProviderConfig is a map of token names to their metadata.
type StaticMockProviderConfig struct { //nolint
	// TokenPrices is a map of token names to their metadata.
	TokenPrices map[string]string `mapstructure:"tokens" toml:"tokens"`
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
