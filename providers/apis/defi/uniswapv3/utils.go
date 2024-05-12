package uniswapv3

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// BaseName is the name of the Uniswap V3 API.
	BaseName = "uniswapv3_api"

	Type = types.ConfigType

	// NameSeparator is the character used to separate elements of dynamic naming for the provider.
	NameSeparator = "-"

	// ContractMethod is the contract method to call for the Uniswap V3 API.
	ContractMethod = "slot0"
)

// ProviderNames is the set of all supported "dynamic" names mapped by chain.
var ProviderNames = map[string]string{
	constants.ETHEREUM: strings.Join([]string{BaseName, constants.ETHEREUM}, NameSeparator),
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

	if pc.BaseDecimals <= 0 {
		return fmt.Errorf("base decimals must be positive")
	}

	if pc.QuoteDecimals <= 0 {
		return fmt.Errorf("quote decimals must be positive")
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
		Name:             fmt.Sprintf("%s%s%s", BaseName, NameSeparator, constants.ETHEREUM),
		Atomic:           true,
		Enabled:          true,
		Timeout:          1000 * time.Millisecond,
		Interval:         2000 * time.Millisecond,
		ReconnectTimeout: 2000 * time.Millisecond,
		MaxQueries:       1,
		URL:              "https://eth.public-rpc.com/",
	}

	DefaultETHProviderConfig = config.ProviderConfig{
		Name: ProviderNames[constants.ETHEREUM],
		API:  DefaultETHAPIConfig,
		Type: Type,
	}

	// DefaultETHMarketConfig is the default market configuration for Uniswap V3. Specifically
	// this is for Ethereum mainnet.
	DefaultETHMarketConfig = types.CurrencyPairsToProviderTickers{
		constants.AAVE_ETH: {
			OffChainTicker: constants.AAVE_ETH.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x5aB53EE1d50eeF2C1DD3d5402789cd27bB52c1bB
				Address:       "0x5aB53EE1d50eeF2C1DD3d5402789cd27bB52c1bB",
				BaseDecimals:  18,
				QuoteDecimals: 18,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.AXELAR_ETH: {
			OffChainTicker: constants.AXELAR_ETH.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0xE7F6720C1F546217081667A5ab7fEbB688036856
				Address:       "0xE7F6720C1F546217081667A5ab7fEbB688036856",
				BaseDecimals:  6,
				QuoteDecimals: 18,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.AXELAR_USDC: {
			OffChainTicker: constants.AXELAR_USDC.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0xAe2A25CBDb19d0dC0DDDD1D2f6b08A6E48c4a9a9
				Address:       "0xAe2A25CBDb19d0dC0DDDD1D2f6b08A6E48c4a9a9",
				BaseDecimals:  6,
				QuoteDecimals: 6,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.CHAINLINK_ETH: {
			OffChainTicker: constants.CHAINLINK_ETH.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0xa6Cc3C2531FdaA6Ae1A3CA84c2855806728693e8
				Address:       "0xa6Cc3C2531FdaA6Ae1A3CA84c2855806728693e8",
				BaseDecimals:  18,
				QuoteDecimals: 18,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.CHAINLINK_USDC: {
			OffChainTicker: constants.CHAINLINK_USDC.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0xFAD57d2039C21811C8F2B5D5B65308aa99D31559
				Address:       "0xFAD57d2039C21811C8F2B5D5B65308aa99D31559",
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.DAI_ETH: {
			OffChainTicker: constants.DAI_ETH.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0xC2e9F25Be6257c210d7Adf0D4Cd6E3E881ba25f8
				Address:       "0xC2e9F25Be6257c210d7Adf0D4Cd6E3E881ba25f8",
				BaseDecimals:  18,
				QuoteDecimals: 18,
				Invert:        true,
			}.MustToJSON(),
		},
		constants.DAI_USDC: {
			OffChainTicker: constants.DAI_USDC.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x5777d92f208679DB4b9778590Fa3CAB3aC9e2168
				Address:       "0x5777d92f208679DB4b9778590Fa3CAB3aC9e2168",
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.ETHEREUM_USDC: {
			OffChainTicker: constants.ETHEREUM_USDC.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640
				Address:       "0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640",
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        true,
			}.MustToJSON(),
		},
		constants.ETHEREUM_USDT: {
			OffChainTicker: constants.ETHEREUM_USDT.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x4e68Ccd3E89f51C3074ca5072bbAC773960dFa36
				Address:       "0x4e68Ccd3E89f51C3074ca5072bbAC773960dFa36",
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.ETHENA_ETH: {
			OffChainTicker: constants.ETHENA_ETH.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0xc3Db44ADC1fCdFd5671f555236eae49f4A8EEa18
				Address:       "0xc3Db44ADC1fCdFd5671f555236eae49f4A8EEa18",
				BaseDecimals:  18,
				QuoteDecimals: 18,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.ETHENA_USDC: {
			OffChainTicker: constants.ETHENA_USDC.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x408A625596f47314e1FD4a6cBCE84C4A8695bA3f
				Address:       "0x408A625596f47314e1FD4a6cBCE84C4A8695bA3f",
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.ETHENA_USDT: {
			OffChainTicker: constants.ETHENA_USDT.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x4185D2952eb74A28EF550a410BA9b8e210Ee9391
				Address:       "0x4185D2952eb74A28EF550a410BA9b8e210Ee9391",
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.ETHERFI_ETH: {
			OffChainTicker: constants.ETHERFI_ETH.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0xDeFDAC77A9A767a2c4eEd826E1AEaD2dAcE53e1C
				Address:       "0xDeFDAC77A9A767a2c4eEd826E1AEaD2dAcE53e1C",
				BaseDecimals:  18,
				QuoteDecimals: 18,
				Invert:        true,
			}.MustToJSON(),
		},
		constants.ETHERFI_USDT: {
			OffChainTicker: constants.ETHERFI_USDT.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x80fa4C1fd0fbB9A4f071999aF69531dee1016644
				Address:       "0x80fa4C1fd0fbB9A4f071999aF69531dee1016644",
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        true,
			}.MustToJSON(),
		},
		constants.HARRY_POTTER_OBAMA_SONIC_10_INU_ETH: {
			OffChainTicker: constants.HARRY_POTTER_OBAMA_SONIC_10_INU_ETH.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x0c30062368eEfB96bF3AdE1218E685306b8E89Fa
				Address:       "0x0c30062368eEfB96bF3AdE1218E685306b8E89Fa",
				BaseDecimals:  8,
				QuoteDecimals: 18,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.LIDO_ETH: {
			OffChainTicker: constants.LIDO_ETH.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0xa3f558aebAecAf0e11cA4b2199cC5Ed341edfd74
				Address:       "0xa3f558aebAecAf0e11cA4b2199cC5Ed341edfd74",
				BaseDecimals:  18,
				QuoteDecimals: 18,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.LIDO_USDC: {
			OffChainTicker: constants.LIDO_USDC.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x78235D08B2aE7a3E00184329212a4d7AcD2F9985
				Address:       "0x78235D08B2aE7a3E00184329212a4d7AcD2F9985",
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.PEPE_ETH: {
			OffChainTicker: constants.PEPE_ETH.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x11950d141EcB863F01007AdD7D1A342041227b58
				Address:       "0x11950d141EcB863F01007AdD7D1A342041227b58",
				BaseDecimals:  18,
				QuoteDecimals: 18,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.PEPE_USDC: {
			OffChainTicker: constants.PEPE_USDC.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0xcEE31C846CbF003F4cEB5Bbd234cBA03C6e940C7
				Address:       "0xcEE31C846CbF003F4cEB5Bbd234cBA03C6e940C7",
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.MAKER_ETH: {
			OffChainTicker: constants.MAKER_ETH.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0xe8c6c9227491C0a8156A0106A0204d881BB7E531
				Address:       "0xe8c6c9227491C0a8156A0106A0204d881BB7E531",
				BaseDecimals:  18,
				QuoteDecimals: 18,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.MAKER_USDC: {
			OffChainTicker: constants.MAKER_USDC.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0xC486Ad2764D55C7dc033487D634195d6e4A6917E
				Address:       "0xC486Ad2764D55C7dc033487D634195d6e4A6917E",
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.MOG_ETH: {
			OffChainTicker: constants.MOG_ETH.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x7832310Cd0de39c4cE0A635F34d9a4B5b47fd434
				Address:       "0x7832310Cd0de39c4cE0A635F34d9a4B5b47fd434",
				BaseDecimals:  18,
				QuoteDecimals: 18,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.UNISWAP_ETH: {
			OffChainTicker: constants.UNISWAP_ETH.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x1d42064Fc4Beb5F8aAF85F4617AE8b3b5B8Bd801
				Address:       "0x1d42064Fc4Beb5F8aAF85F4617AE8b3b5B8Bd801",
				BaseDecimals:  18,
				QuoteDecimals: 18,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.UNISWAP_USDC: {
			OffChainTicker: constants.UNISWAP_USDC.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0xD0fC8bA7E267f2bc56044A7715A489d851dC6D78
				Address:       "0xD0fC8bA7E267f2bc56044A7715A489d851dC6D78",
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.UNISWAP_USDT: {
			OffChainTicker: constants.UNISWAP_USDT.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x3470447f3CecfFAc709D3e783A307790b0208d60
				Address:       "0x3470447f3CecfFAc709D3e783A307790b0208d60",
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.WBITCOIN_ETH: {
			OffChainTicker: constants.WBITCOIN_ETH.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0xCBCdF9626bC03E24f779434178A73a0B4bad62eD
				Address:       "0xCBCdF9626bC03E24f779434178A73a0B4bad62eD",
				BaseDecimals:  8,
				QuoteDecimals: 18,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.WBITCOIN_USDC: {
			OffChainTicker: constants.WBITCOIN_USDC.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x99ac8cA7087fA4A2A1FB6357269965A2014ABc35
				Address:       "0x99ac8cA7087fA4A2A1FB6357269965A2014ABc35",
				BaseDecimals:  8,
				QuoteDecimals: 6,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.WSTETH_ETH: {
			OffChainTicker: constants.WSTETH_ETH.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x109830a1AAaD605BbF02a9dFA7B0B92EC2FB7dAa
				Address:       "0x109830a1AAaD605BbF02a9dFA7B0B92EC2FB7dAa",
				BaseDecimals:  18,
				QuoteDecimals: 18,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.WSTETH_USDC: {
			OffChainTicker: constants.WSTETH_USDC.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x4622Df6fB2d9Bee0DCDaCF545aCDB6a2b2f4f863
				Address:       "0x4622Df6fB2d9Bee0DCDaCF545aCDB6a2b2f4f863",
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.WTAO_ETH: {
			OffChainTicker: constants.WTAO_ETH.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0x433a00819C771b33FA7223a5B3499b24FBCd1bBC
				Address:       "0x433a00819C771b33FA7223a5B3499b24FBCd1bBC",
				BaseDecimals:  9,
				QuoteDecimals: 18,
				Invert:        false,
			}.MustToJSON(),
		},
		constants.WTAO_USDC: {
			OffChainTicker: constants.WTAO_USDC.String(),
			JSON: PoolConfig{
				// REF: https://app.uniswap.org/explore/pools/ethereum/0xf763Bb342eB3d23C02ccB86312422fe0c1c17E94
				Address:       "0xf763Bb342eB3d23C02ccB86312422fe0c1c17E94",
				BaseDecimals:  9,
				QuoteDecimals: 6,
				Invert:        false,
			}.MustToJSON(),
		},
	}
)
