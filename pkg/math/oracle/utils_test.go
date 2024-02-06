package oracle_test

import (
	"math/big"
	"testing"

	"github.com/skip-mev/slinky/pkg/math/oracle"
	"github.com/stretchr/testify/require"
)

func TestInvertCurrencyPairPrice(t *testing.T) {
	t.Run("can invert a price of 1", func(t *testing.T) {
		one := oracle.ScaledOne(oracle.ScaledDecimals)
		inverted := oracle.InvertCurrencyPairPrice(one, oracle.ScaledDecimals)
		require.Equal(t, one, inverted)
	})

	t.Run("can invert a price of 2000", func(t *testing.T) {
		price := big.NewInt(2000)
		scaledPrice := new(big.Int).Mul(price, oracle.ScaledOne(oracle.ScaledDecimals))
		inverted := oracle.InvertCurrencyPairPrice(scaledPrice, oracle.ScaledDecimals)

		expectedExp := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(oracle.ScaledDecimals-4), nil)
		expectedPrice := big.NewInt(5)
		expectedScaledPrice := new(big.Int).Mul(expectedPrice, expectedExp)
		require.Equal(t, expectedScaledPrice, inverted)
	})

	t.Run("can invert a price of 2", func(t *testing.T) {
		price := big.NewInt(2)
		scaledPrice := new(big.Int).Mul(price, oracle.ScaledOne(oracle.ScaledDecimals))
		inverted := oracle.InvertCurrencyPairPrice(scaledPrice, oracle.ScaledDecimals)

		expectedExp := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(oracle.ScaledDecimals-1), nil)
		expectedPrice := big.NewInt(5)
		expectedScaledPrice := new(big.Int).Mul(expectedPrice, expectedExp)
		require.Equal(t, expectedScaledPrice, inverted)
	})

	t.Run("can invert a price of 0.5", func(t *testing.T) {
		price := big.NewInt(5)
		exp := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(oracle.ScaledDecimals-1), nil)
		scaledPrice := new(big.Int).Mul(price, exp)
		inverted := oracle.InvertCurrencyPairPrice(scaledPrice, oracle.ScaledDecimals)

		expectedExp := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(oracle.ScaledDecimals), nil)
		expectedPrice := big.NewInt(2)
		expectedScaledPrice := new(big.Int).Mul(expectedPrice, expectedExp)
		require.Equal(t, expectedScaledPrice, inverted)
	})
}

func TestScaleUpCurrencyPairPrice(t *testing.T) {}

func TestScaleDownCurrencyPairPrice(t *testing.T) {}

func TestScaledOne(t *testing.T) {}
