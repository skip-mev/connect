package uniswapv3_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/providers/apis/defi/uniswapv3"
)

func TestConvertSquareRootX96Price(t *testing.T) {
	t.Run("weth/usdc uniswap primer example", func(t *testing.T) {
		val, converted := big.NewInt(1).SetString("2018382873588440326581633304624437", 10)
		require.True(t, converted)

		expected := big.NewFloat(649004842.7013700766389061032587755).SetPrec(40)
		actual := uniswapv3.ConvertSquareRootX96Price(val).SetPrec(40)
		require.Equal(t, expected, actual)
	})

	t.Run("works with a value of 0", func(t *testing.T) {
		val := big.NewInt(0)
		expected := big.NewFloat(0).SetPrec(40)
		actual := uniswapv3.ConvertSquareRootX96Price(val).SetPrec(40)
		require.Equal(t, expected, actual)
	})

	t.Run("should be 1 when the value is 2^96", func(t *testing.T) {
		val := new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
		expected := big.NewFloat(1).SetPrec(40)
		actual := uniswapv3.ConvertSquareRootX96Price(val).SetPrec(40)
		require.Equal(t, expected, actual)
	})
}

func TestScalePrice(t *testing.T) {
	testCases := []struct {
		name     string
		price    *big.Float
		cfg      uniswapv3.PoolConfig
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
			expected: big.NewFloat(1540.820552028),
		},
		{
			name:  "mainnet example for weth/usdc",
			price: big.NewFloat(2.913786192888320737692333570997812e+08),
			cfg: uniswapv3.PoolConfig{
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        true,
			},
			expected: big.NewFloat(3431.960802205393266704),
		},
		{
			name:  "mainnet example for eth/usdc",
			price: big.NewFloat(2.926645918358364572159014271666027e+08),
			cfg: uniswapv3.PoolConfig{
				BaseDecimals:  18,
				QuoteDecimals: 6,
				Invert:        true,
			},
			expected: big.NewFloat(3416.880715658719806983),
		},
		{
			name:  "mainnet example for mog/eth",
			price: big.NewFloat(1.63833946559934409985296037965e-10),
			cfg: uniswapv3.PoolConfig{
				BaseDecimals:  18,
				QuoteDecimals: 18,
				Invert:        false,
			},
			expected: big.NewFloat(1.63833946559934409985296037965e-10),
		},
		{
			name:  "mainnet example for btc/usdt",
			price: big.NewFloat(688.8936521667327881055693350566),
			cfg: uniswapv3.PoolConfig{
				BaseDecimals:  8,
				QuoteDecimals: 6,
				Invert:        false,
			},
			expected: big.NewFloat(68889.36521667327881055693350566),
		},
		{
			name:  "mainnet example for btc/usdt where usdt now assumes 10",
			price: big.NewFloat(6888936.521667327881055693350566),
			cfg: uniswapv3.PoolConfig{
				BaseDecimals:  8,
				QuoteDecimals: 10,
				Invert:        false,
			},
			expected: big.NewFloat(68889.36521667327881055693350566),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := uniswapv3.ScalePrice(tc.cfg, tc.price).SetPrec(40)
			require.Equal(t, tc.expected.SetPrec(40), actual)
		})
	}
}
