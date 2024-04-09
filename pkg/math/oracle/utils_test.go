package oracle_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math/oracle"
	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
)

func TestGetProviderPrice(t *testing.T) {
	t.Run("provider does not exist in the cache", func(t *testing.T) {
		agg, err := oracle.NewMedianAggregator(logger, marketmap, nil)
		require.NoError(t, err)

		cfg := mmtypes.ProviderConfig{
			Name:           "test",
			OffChainTicker: "BTC/USD",
		}
		_, err = agg.GetProviderPrice(cfg)
		require.Error(t, err)
	})

	t.Run("provider exists in the cache but does not have the desired CP", func(t *testing.T) {
		agg, err := oracle.NewMedianAggregator(logger, marketmap, nil)
		require.NoError(t, err)

		cfg := mmtypes.ProviderConfig{
			Name:           "test",
			OffChainTicker: "BTC/USD",
		}
		prices := types.AggregatorPrices{
			"BTC/USDT": big.NewFloat(100),
		}
		agg.SetProviderData("test", prices)

		_, err = agg.GetProviderPrice(cfg)
		require.Error(t, err)
	})

	t.Run("provider exists in the cache and has the desired CP", func(t *testing.T) {
		agg, err := oracle.NewMedianAggregator(logger, marketmap, nil)
		require.NoError(t, err)

		cfg := mmtypes.ProviderConfig{
			Name:           "test",
			OffChainTicker: "BTC/USD",
		}
		prices := types.AggregatorPrices{
			"BTC/USD": big.NewFloat(100),
		}
		agg.SetProviderData("test", prices)

		price, err := agg.GetProviderPrice(cfg)
		require.NoError(t, err)
		require.Equal(t, big.NewFloat(100), price)
	})

	t.Run("provider exists in the cache and has the desired CP, invert is true", func(t *testing.T) {
		agg, err := oracle.NewMedianAggregator(logger, marketmap, nil)
		require.NoError(t, err)

		cfg := mmtypes.ProviderConfig{
			Name:           "test",
			OffChainTicker: "BTC/USD",
			Invert:         true,
		}
		prices := types.AggregatorPrices{
			"BTC/USD": big.NewFloat(100),
		}
		agg.SetProviderData("test", prices)

		price, err := agg.GetProviderPrice(cfg)
		require.NoError(t, err)
		require.Equal(t, big.NewFloat(0.01).SetPrec(18), price.SetPrec(18))
	})
}

func TestGetIndexPrice(t *testing.T) {
	t.Run("has no index prices", func(t *testing.T) {
		agg, err := oracle.NewMedianAggregator(logger, marketmap, nil)
		require.NoError(t, err)

		_, err = agg.GetIndexPrice(constants.ETHEREUM_USD)
		require.Error(t, err)
	})

	t.Run("has index prices", func(t *testing.T) {
		agg, err := oracle.NewMedianAggregator(logger, marketmap, nil)
		require.NoError(t, err)

		prices := types.AggregatorPrices{
			constants.BITCOIN_USD.String(): big.NewFloat(100),
		}
		agg.SetIndexPrices(prices)

		price, err := agg.GetIndexPrice(constants.BITCOIN_USD)
		require.NoError(t, err)
		require.Equal(t, big.NewFloat(100), price)
	})
}
