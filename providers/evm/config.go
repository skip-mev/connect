package evm

import (
	"github.com/spf13/viper"
)

type (
	// Config is the config struct for EVM providers.
	Config struct {
		// APIKey is the API key used to make requests to the EVM API.
		APIKey string `mapstructure:"api_key" toml:"api_key"`
		// TokenNameToMetadata is a map of token names to their metadata.
		TokenNameToMetadata map[string]TokenMetadata `mapstructure:"token_name_to_metadata" toml:"token_name_to_metadata"`
		// RPCEndpoint is the endpoint of the ethereum rpc node to use for querying
		RPCEndpoint string `mapstructure:"rpc_endpoint" toml:"rpc_endpoint"`
	}

	// TokenMetadata is the metadata for a token.
	TokenMetadata struct {
		// Symbol is the provider-specific token identifier. This can be a name, ticker, contract address, etc.
		Symbol string `mapstructure:"symbol" toml:"symbol"`
		// Decimals is the number of decimal places the token has on chain.
		Decimals uint64 `mapstructure:"decimals" toml:"decimals"`
		// IsTWAP indicates whether this token's price is a time weighted average.
		IsTWAP bool `mapstructure:"is_twap" toml:"is_twap"`
	}
)

// ReadEVMConfigFromFile reads a config from a file and returns the config.
func ReadEVMConfigFromFile(path string) (Config, error) {
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
