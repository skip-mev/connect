package oracle_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/pkg/math/oracle"
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

func TestScaleUpCurrencyPairPrice(t *testing.T) {
	t.Run("can scale up a price of 1", func(t *testing.T) {
		price := big.NewInt(1)
		scaledPrice, err := oracle.ScaleUpCurrencyPairPrice(0, price)
		require.NoError(t, err)

		one := oracle.ScaledOne(oracle.ScaledDecimals)
		require.Equal(t, one, scaledPrice)
	})

	t.Run("can scale up a price of 2000", func(t *testing.T) {
		price := big.NewInt(2000)
		exp := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(oracle.ScaledDecimals-8), nil)
		expectedPrice := new(big.Int).Mul(price, exp)

		scaledPrice, err := oracle.ScaleUpCurrencyPairPrice(8, price)
		require.NoError(t, err)
		require.Equal(t, expectedPrice, scaledPrice)
	})

	t.Run("errors when scaling up a price with more decimals than the standard", func(t *testing.T) {
		price := big.NewInt(2000)
		_, err := oracle.ScaleUpCurrencyPairPrice(oracle.ScaledDecimals+1, price)
		require.Error(t, err)
	})

	t.Run("equal number of decimal points", func(t *testing.T) {
		price := big.NewInt(2000)
		exp := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(oracle.ScaledDecimals), nil)
		expectedPrice := new(big.Int).Mul(price, exp)

		scaledPrice, err := oracle.ScaleUpCurrencyPairPrice(oracle.ScaledDecimals, expectedPrice)
		require.NoError(t, err)
		require.Equal(t, expectedPrice, scaledPrice)
	})
}

func TestScaleDownCurrencyPairPrice(t *testing.T) {
	t.Run("can scale down a price of 1", func(t *testing.T) {
		one := oracle.ScaledOne(oracle.ScaledDecimals)
		scaledPrice, err := oracle.ScaleDownCurrencyPairPrice(0, one)
		require.NoError(t, err)

		require.Equal(t, big.NewInt(1), scaledPrice)
	})

	t.Run("can scale down a price of 2000", func(t *testing.T) {
		price := big.NewInt(2000)
		exp := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(8), nil)
		price = new(big.Int).Mul(price, exp)

		scaledPrice, err := oracle.ScaleUpCurrencyPairPrice(8, price)
		require.NoError(t, err)

		unscaledPrice, err := oracle.ScaleDownCurrencyPairPrice(8, scaledPrice)
		require.NoError(t, err)
		require.Equal(t, price, unscaledPrice)
	})

	t.Run("errors when scaling down a price with more decimals than the standard", func(t *testing.T) {
		price := big.NewInt(2000)
		_, err := oracle.ScaleDownCurrencyPairPrice(oracle.ScaledDecimals+1, price)
		require.Error(t, err)
	})

	t.Run("equal number of decimal points", func(t *testing.T) {
		price := big.NewInt(2000)
		exp := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(oracle.ScaledDecimals), nil)
		price = new(big.Int).Mul(price, exp)

		scaledPrice, err := oracle.ScaleDownCurrencyPairPrice(oracle.ScaledDecimals, price)
		require.NoError(t, err)
		require.Equal(t, price, scaledPrice)
	})
}
