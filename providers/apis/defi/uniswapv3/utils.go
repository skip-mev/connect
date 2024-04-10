package uniswapv3

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// Name is the name of the Uniswap V3 API.
	Name = "uniswapv3_api"

	// ContractMethod is the contract method to call for the Uniswap V3 API.
	ContractMethod = "slot0"
)

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

	if pc.BaseDecimals <= 0 {
		return fmt.Errorf("base decimals must be positive")
	}

	if pc.QuoteDecimals <= 0 {
		return fmt.Errorf("quote decimals must be positive")
	}

	return nil
}

// MustToJSON converts the pool configuration to JSON.
func (pc *PoolConfig) MustToJSON() string {
	b, err := json.Marshal(pc)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// DefaultAPIConfig is the default configuration for the Uniswap API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           true,
	Enabled:          true,
	Timeout:          1000 * time.Millisecond,
	Interval:         2000 * time.Millisecond,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	URL:              "https://eth.public-rpc.com/",
}
