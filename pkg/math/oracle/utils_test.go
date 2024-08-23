package oracle_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math/oracle"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

func TestGetProviderPrice(t *testing.T) {
	t.Run("provider does not exist in the cache", func(t *testing.T) {
		agg, err := oracle.NewIndexPriceAggregator(logger, marketmap, nil)
		require.NoError(t, err)

		cfg := mmtypes.ProviderConfig{
			Name:           "test",
			OffChainTicker: "BTC/USD",
		}
		_, err = agg.GetProviderPrice(cfg)
		require.Error(t, err)
	})

	t.Run("provider exists in the cache but does not have the desired CP", func(t *testing.T) {
		agg, err := oracle.NewIndexPriceAggregator(logger, marketmap, nil)
		require.NoError(t, err)

		cfg := mmtypes.ProviderConfig{
			Name:           "test",
			OffChainTicker: "BTC/USD",
		}
		prices := types.Prices{
			"BTC/USDT": big.NewFloat(100),
		}
		agg.SetProviderPrices("test", prices)

		_, err = agg.GetProviderPrice(cfg)
		require.Error(t, err)
	})

	t.Run("provider exists in the cache and has the desired CP", func(t *testing.T) {
		agg, err := oracle.NewIndexPriceAggregator(logger, marketmap, nil)
		require.NoError(t, err)

		cfg := mmtypes.ProviderConfig{
			Name:           "test",
			OffChainTicker: "BTC/USD",
		}
		prices := types.Prices{
			"BTC/USD": big.NewFloat(100),
		}
		agg.SetProviderPrices("test", prices)

		price, err := agg.GetProviderPrice(cfg)
		require.NoError(t, err)
		require.Equal(t, big.NewFloat(100), price)
	})

	t.Run("provider exists in the cache and has the desired CP, invert is true", func(t *testing.T) {
		agg, err := oracle.NewIndexPriceAggregator(logger, marketmap, nil)
		require.NoError(t, err)

		cfg := mmtypes.ProviderConfig{
			Name:           "test",
			OffChainTicker: "BTC/USD",
			Invert:         true,
		}
		prices := types.Prices{
			"BTC/USD": big.NewFloat(100),
		}
		agg.SetProviderPrices("test", prices)

		price, err := agg.GetProviderPrice(cfg)
		require.NoError(t, err)
		require.Equal(t, big.NewFloat(0.01).SetPrec(18), price.SetPrec(18))
	})

	t.Run("provider price is nil", func(t *testing.T) {
		agg, err := oracle.NewIndexPriceAggregator(logger, marketmap, nil)
		require.NoError(t, err)

		cfg := mmtypes.ProviderConfig{
			Name:           "test",
			OffChainTicker: "BTC/USD",
		}
		prices := types.Prices{
			"BTC/USD": nil,
		}
		agg.SetProviderPrices("test", prices)

		_, err = agg.GetProviderPrice(cfg)
		require.Error(t, err)
	})
}

func TestGetIndexPrice(t *testing.T) {
	t.Run("has no index prices", func(t *testing.T) {
		agg, err := oracle.NewIndexPriceAggregator(logger, marketmap, nil)
		require.NoError(t, err)

		_, err = agg.GetIndexPrice(ethusdCP)
		require.Error(t, err)
	})

	t.Run("has index prices", func(t *testing.T) {
		agg, err := oracle.NewIndexPriceAggregator(logger, marketmap, nil)
		require.NoError(t, err)

		prices := types.Prices{
			btcusdCP.String(): big.NewFloat(100),
		}
		agg.SetIndexPrices(prices)

		price, err := agg.GetIndexPrice(btcusdCP)
		require.NoError(t, err)
		require.Equal(t, big.NewFloat(100), price)
	})

	t.Run("index price is nil", func(t *testing.T) {
		agg, err := oracle.NewIndexPriceAggregator(logger, marketmap, nil)
		require.NoError(t, err)

		prices := types.Prices{
			btcusdCP.String(): nil,
		}
		agg.SetIndexPrices(prices)

		_, err = agg.GetIndexPrice(btcusdCP)
		require.Error(t, err)
	})
}
