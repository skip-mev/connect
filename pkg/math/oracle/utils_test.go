package oracle_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math/oracle"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func TestGetTickerFromOperation(t *testing.T) {
	t.Run("has ticker included in the market config", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		operation := mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
		}
		ticker, err := m.GetTickerFromOperation(operation)
		require.NoError(t, err)
		require.Equal(t, BTC_USD, ticker)
	})

	t.Run("has ticker not included in the market config", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		operation := mmtypes.Operation{
			CurrencyPair: constants.MOG_USD.CurrencyPair,
		}
		ticker, err := m.GetTickerFromOperation(operation)
		require.Error(t, err)
		require.Empty(t, ticker)
	})
}

func TestGetProviderPrice(t *testing.T) {
	t.Run("does not have a ticker in the config", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		operation := mmtypes.Operation{
			CurrencyPair: constants.MOG_USD.CurrencyPair,
		}
		_, err = m.GetProviderPrice(operation)
		require.Error(t, err)
	})

	t.Run("has no provider prices or index prices", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		// Attempt to retrieve the provider.
		operation := mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     coinbase.Name,
		}
		_, err = m.GetProviderPrice(operation)
		require.Error(t, err)

		// Attempt to retrieve the index price.
		operation = mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     mmtypes.IndexPrice,
		}
		_, err = m.GetProviderPrice(operation)
		require.Error(t, err)
	})

	t.Run("has provider prices but no index prices", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		// Set the provider price.
		prices := types.TickerPrices{
			BTC_USD: createPrice(100, BTC_USD.Decimals),
		}
		m.PriceAggregator.SetProviderData(coinbase.Name, prices)

		// Attempt to retrieve the provider.
		operation := mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     coinbase.Name,
		}
		price, err := m.GetProviderPrice(operation)
		require.NoError(t, err)
		require.Equal(t, createPrice(100, oracle.ScaledDecimals), price)

		// Attempt to retrieve the index price.
		operation = mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     mmtypes.IndexPrice,
		}
		_, err = m.GetProviderPrice(operation)
		require.Error(t, err)
	})

	t.Run("has provider prices and index prices", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		// Set the provider price.
		prices := types.TickerPrices{
			BTC_USD: createPrice(100, BTC_USD.Decimals),
		}
		m.PriceAggregator.SetProviderData(coinbase.Name, prices)

		// Set the index price.
		m.PriceAggregator.SetAggregatedData(prices)

		// Attempt to retrieve the provider.
		operation := mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     coinbase.Name,
		}
		price, err := m.GetProviderPrice(operation)
		require.NoError(t, err)
		require.Equal(t, createPrice(100, oracle.ScaledDecimals), price)

		// Attempt to retrieve the index price.
		operation = mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     mmtypes.IndexPrice,
		}
		price, err = m.GetProviderPrice(operation)
		require.NoError(t, err)
		require.Equal(t, createPrice(100, oracle.ScaledDecimals), price)
	})

	t.Run("has provider prices and can correctly scale up", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		// Set the provider price.
		prices := types.TickerPrices{
			BTC_USD: createPrice(40_000, BTC_USD.Decimals),
		}
		m.PriceAggregator.SetProviderData(coinbase.Name, prices)

		// Attempt to retrieve the provider.
		operation := mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     coinbase.Name,
		}
		price, err := m.GetProviderPrice(operation)
		require.NoError(t, err)
		require.Equal(t, createPrice(40_000, oracle.ScaledDecimals), price)
	})

	t.Run("has provider prices and can correctly invert", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		// Set the provider price.
		prices := types.TickerPrices{
			BTC_USD: createPrice(40_000, BTC_USD.Decimals),
		}
		m.PriceAggregator.SetProviderData(coinbase.Name, prices)

		// Attempt to retrieve the provider.
		operation := mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     coinbase.Name,
			Invert:       true,
		}
		price, err := m.GetProviderPrice(operation)
		require.NoError(t, err)
		expectedPrice := createPrice(0.000025, oracle.ScaledDecimals)
		verifyPrice(t, expectedPrice, price)
	})
}
