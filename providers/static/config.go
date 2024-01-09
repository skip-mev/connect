package static

import (
	"github.com/spf13/viper"
)

// StaticMockProviderConfig is a map of token names to their metadata.
type StaticMockProviderConfig struct { //nolint
	// TokenPrices is a map of token names to their metadata.
	TokenPrices map[string]string `json:"tokenPrices" validate:"required"`
}

// ReadStaticMockProviderConfigFromFile reads the static mock provider config from the given file.
func ReadStaticMockProviderConfigFromFile(path string) (StaticMockProviderConfig, error) {
	var config StaticMockProviderConfig

	// read in the config file
	viper.SetConfigFile(path)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	// parse config
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
