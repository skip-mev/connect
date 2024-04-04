package uniswap_test

import (
	"math/big"
	"testing"

	"github.com/skip-mev/slinky/providers/apis/uniswap"
	"github.com/stretchr/testify/require"
)

// precision is the precision used for big.Float calculations. Specifically
// this is used to ensure that float values are the same within a certain
// precision.
const precision = 30

func TestConvertSquareRootX96Price(t *testing.T) {
	t.Run("weth/usdc uniswap primer example", func(t *testing.T) {
		val, converted := big.NewInt(1).SetString("2018382873588440326581633304624437", 10)
		require.True(t, converted)

		expected := big.NewFloat(649004842.7013700766389061032587755).SetPrec(precision)
		actual := uniswap.ConvertSquareRootX96Price(val).SetPrec(precision)
		require.Equal(t, expected, actual)
	})

	t.Run("works with a value of 0", func(t *testing.T) {
		val := big.NewInt(0)
		expected := big.NewFloat(0).SetPrec(precision)
		actual := uniswap.ConvertSquareRootX96Price(val).SetPrec(precision)
		require.Equal(t, expected, actual)
	})

	t.Run("should be 1 when the value is 2^96", func(t *testing.T) {
		val := new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
		expected := big.NewFloat(1).SetPrec(precision)
		actual := uniswap.ConvertSquareRootX96Price(val).SetPrec(precision)
		require.Equal(t, expected, actual)
	})
}

func TestScalePrice(t *testing.T) {
}

func TestInvertCurrencyPairPrice(t *testing.T) {

}
