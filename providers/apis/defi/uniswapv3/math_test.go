package uniswapv3_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/providers/apis/defi/uniswapv3"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func TestConvertSquareRootX96Price(t *testing.T) {
	t.Run("weth/usdc uniswap primer example", func(t *testing.T) {
		val, converted := big.NewInt(1).SetString("2018382873588440326581633304624437", 10)
		require.True(t, converted)

		expected := big.NewFloat(649004842.7013700766389061032587755).SetPrec(precision)
		actual := uniswapv3.ConvertSquareRootX96Price(val).SetPrec(precision)
		require.Equal(t, expected, actual)
	})

	t.Run("works with a value of 0", func(t *testing.T) {
		val := big.NewInt(0)
		expected := big.NewFloat(0).SetPrec(precision)
		actual := uniswapv3.ConvertSquareRootX96Price(val).SetPrec(precision)
		require.Equal(t, expected, actual)
	})

	t.Run("should be 1 when the value is 2^96", func(t *testing.T) {
		val := new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
		expected := big.NewFloat(1).SetPrec(precision)
		actual := uniswapv3.ConvertSquareRootX96Price(val).SetPrec(precision)
		require.Equal(t, expected, actual)
	})
}

func TestScalePrice(t *testing.T) {
	testCases := []struct {
		name     string
		price    *big.Float
		cfg      uniswapv3.PoolConfig
		ticker   mmtypes.Ticker
		expected *big.Float
	}{
		{
			name:  "uniswap primer example for weth/usdc",
			price: big.NewFloat(649004842.70137),
			cfg: uniswapv3.PoolConfig{
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        true,
			},
			ticker: mmtypes.Ticker{
				Decimals: 18,
			},
			expected: big.NewFloat(1540.82 * 1e18),
		},
		{
			name:  "uniswap primer example for eth/usdc but with lower precision",
			price: big.NewFloat(649004842.70137),
			cfg: uniswapv3.PoolConfig{
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        true,
			},
			ticker: mmtypes.Ticker{
				Decimals: 6,
			},
			expected: big.NewFloat(1540.82 * 1e6),
		},
		{
			name:  "mainnet example for weth/usdc",
			price: big.NewFloat(2.913786192888320737692333570997812e+08),
			cfg: uniswapv3.PoolConfig{
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        true,
			},
			ticker: mmtypes.Ticker{
				Decimals: 18,
			},
			expected: big.NewFloat(3431.960802205393266704 * 1e18),
		},
		{
			name:  "mainnet example for weth/usdc but with lower precision",
			price: big.NewFloat(2.913786192888320737692333570997812e+08),
			cfg: uniswapv3.PoolConfig{
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        true,
			},
			ticker: mmtypes.Ticker{
				Decimals: 6,
			},
			expected: big.NewFloat(3431.96 * 1e6),
		},
		{
			name:  "mainnet example for eth/usdc",
			price: big.NewFloat(2.926645918358364572159014271666027e+08),
			cfg: uniswapv3.PoolConfig{
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        true,
			},
			ticker: mmtypes.Ticker{
				Decimals: 18,
			},
			expected: big.NewFloat(3416.880715658719806983 * 1e18),
		},
		{
			name:  "mainnet example for mog/eth",
			price: big.NewFloat(1.63833946559934409985296037965e-10),
			cfg: uniswapv3.PoolConfig{
				BaseDecimals:  18,
				QuoteDecimals: 18,
				Invert:        false,
			},
			ticker: mmtypes.Ticker{
				Decimals: 18,
			},
			expected: big.NewFloat(163833946),
		},
		{
			name:  "mainnet example for mog/eth but with lower precision",
			price: big.NewFloat(1.63833946559934409985296037965e-10),
			cfg: uniswapv3.PoolConfig{
				BaseDecimals:  18,
				QuoteDecimals: 18,
				Invert:        false,
			},
			ticker: mmtypes.Ticker{
				Decimals: 12,
			},
			expected: big.NewFloat(163.833946),
		},
		{
			name:  "mainnet example for btc/usdt",
			price: big.NewFloat(688.8936521667327881055693350566),
			cfg: uniswapv3.PoolConfig{
				BaseDecimals:  8,
				QuoteDecimals: 6,
				Invert:        false,
			},
			ticker: mmtypes.Ticker{
				Decimals: 5,
			},
			expected: big.NewFloat(68889.36521667327881055693350566 * 1e5),
		},
		{
			name:  "mainnet example for btc/usdt where usdt now assumes 10",
			price: big.NewFloat(6888936.521667327881055693350566),
			cfg: uniswapv3.PoolConfig{
				BaseDecimals:  8,
				QuoteDecimals: 10,
				Invert:        false,
			},
			ticker: mmtypes.Ticker{
				Decimals: 5,
			},
			expected: big.NewFloat(68889.36521667327881055693350566 * 1e5),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := uniswapv3.ScalePrice(tc.ticker, tc.cfg, tc.price).SetPrec(5)
			require.Equal(t, tc.expected.SetPrec(5), actual)
		})
	}
}

func TestGetScalingFactor(t *testing.T) {
	t.Run("base and quote decimals for erc20 tokens are the same", func(t *testing.T) {
		cfg := uniswapv3.PoolConfig{
			BaseDecimals:  18,
			QuoteDecimals: 18,
		}

		actual := uniswapv3.GetScalingFactor(
			cfg.BaseDecimals,
			cfg.QuoteDecimals,
		).SetPrec(precision)

		expected := big.NewFloat(1).SetPrec(precision)
		require.Equal(t, expected, actual)
	})

	t.Run("base decimals are greater than quote decimals for erc20 tokens", func(t *testing.T) {
		cfg := uniswapv3.PoolConfig{
			BaseDecimals:  18,
			QuoteDecimals: 6,
		}

		actual := uniswapv3.GetScalingFactor(
			cfg.BaseDecimals,
			cfg.QuoteDecimals,
		).SetPrec(precision)

		expected := big.NewFloat(1e12).SetPrec(precision)
		require.Equal(t, expected, actual)
	})

	t.Run("quote decimals are greater than base decimals for erc20 tokens", func(t *testing.T) {
		cfg := uniswapv3.PoolConfig{
			BaseDecimals:  6,
			QuoteDecimals: 18,
		}

		actual := uniswapv3.GetScalingFactor(
			cfg.BaseDecimals,
			cfg.QuoteDecimals,
		).SetPrec(precision)

		expected := big.NewFloat(1e-12).SetPrec(precision)
		require.Equal(t, expected, actual)
	})
}
