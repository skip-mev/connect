package uniswapv3

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/constants"
)

const (
	// BaseName is the name of the Uniswap V3 API.
	BaseName = "uniswapv3_api"

	// NameSeparator is the character used to separate elements of dynamic naming for the provider.
	NameSeparator = "-"

	// ContractMethod is the contract method to call for the Uniswap V3 API.
	ContractMethod = "slot0"

	// ETH_URL is the URL for the Uniswap V3 API. This uses a free public RPC provider on Ethereum Mainnet.
	ETH_URL = "https://eth.public-rpc.com/"

	// BASE_URL is the URL for the Uniswap V3 API. This uses a free public RPC provider on Base Mainnet.
	BASE_URL = "https://mainnet.base.org"
)

// ProviderNames is the set of all supported "dynamic" names mapped by chain.
var ProviderNames = map[string]string{
	constants.ETHEREUM: strings.Join([]string{BaseName, constants.ETHEREUM}, NameSeparator),
	constants.BASE:     strings.Join([]string{BaseName, constants.BASE}, NameSeparator),
}

// IsValidProviderName returns a bool based on the validity of the passed in name.
// Dynamic provider naming is supported via `BaseName“NameSeparator“SupportedChain`.
func IsValidProviderName(name string) bool {
	for _, providerName := range ProviderNames {
		if name == providerName {
			return true
		}
	}
	return false
}

// PoolConfig is the configuration for a Uniswap V3 pool. This is specific to each pair of tokens.
type PoolConfig struct {
	// Address is the Uniswap V3 pool address.
	Address string `json:"address"`
	// BaseDecimals is the number of decimals for the base token. This should be derived from the
	// token contract.
	BaseDecimals int64 `json:"base_decimals"`
	// QuoteDecimals is the number of decimals for the quote token. This should be derived from the
	// token contract.
	QuoteDecimals int64 `json:"quote_decimals"`
	// Invert is utilized to invert the price of a pool's reserves. This may be required for certain
	// pools as the price is derived based on the sorted order of the ERC20 addresses of the tokens
	// in the pool.
	Invert bool `json:"invert"`
}

// ValidateBasic validates the pool configuration.
func (pc *PoolConfig) ValidateBasic() error {
	if !common.IsHexAddress(pc.Address) {
		return fmt.Errorf("pool address is not a valid ethereum address")
	}

	if pc.BaseDecimals < 0 {
		return fmt.Errorf("base decimals must be non-negative")
	}

	if pc.QuoteDecimals < 0 {
		return fmt.Errorf("quote decimals must be non-negative")
	}

	return nil
}

// MustToJSON converts the pool configuration to JSON.
func (pc PoolConfig) MustToJSON() string {
	b, err := json.Marshal(pc)
	if err != nil {
		panic(err)
	}
	return string(b)
}

var (
	// DefaultETHAPIConfig is the default configuration for the Uniswap API. Specifically this is for
	// Ethereum mainnet.
	DefaultETHAPIConfig = config.APIConfig{
		Name:              fmt.Sprintf("%s%s%s", BaseName, NameSeparator, constants.ETHEREUM),
		Atomic:            true,
		Enabled:           true,
		Timeout:           1000 * time.Millisecond,
		Interval:          2000 * time.Millisecond,
		ReconnectTimeout:  2000 * time.Millisecond,
		MaxQueries:        1,
		Endpoints:         []config.Endpoint{{URL: ETH_URL}},
		MaxBlockHeightAge: 30 * time.Second,
	}

	// DefaultBaseAPIConfig is the default configuration for the Uniswap API. Specifically this is for
	// Base mainnet.
	DefaultBaseAPIConfig = config.APIConfig{
		Name:              fmt.Sprintf("%s%s%s", BaseName, NameSeparator, constants.BASE),
		Atomic:            true,
		Enabled:           true,
		Timeout:           1000 * time.Millisecond,
		Interval:          2000 * time.Millisecond,
		ReconnectTimeout:  2000 * time.Millisecond,
		MaxQueries:        1,
		Endpoints:         []config.Endpoint{{URL: BASE_URL}},
		MaxBlockHeightAge: 30 * time.Second,
	}
)
