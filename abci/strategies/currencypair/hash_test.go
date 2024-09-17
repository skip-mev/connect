package currencypair_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	strategies "github.com/skip-mev/connect/v2/abci/strategies/currencypair"
	"github.com/skip-mev/connect/v2/abci/strategies/currencypair/mocks"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
)

func TestHashCurrencyPairStrategyID(t *testing.T) {
	ok := mocks.NewOracleKeeper(t)
	ctx := sdk.Context{}
	strategy := strategies.NewHashCurrencyPairStrategy(ok)

	t.Run("test a single valid currency pair getting a hash", func(t *testing.T) {
		ok.On("GetIDForCurrencyPair", mock.Anything, btcusd).Return(uint64(0), true).Once()

		// expect the first currency-pair to have ID 0
		id, err := strategy.ID(ctx, btcusd)
		require.NoError(t, err)
		require.NotEqual(t, uint64(0), id)
	})

	t.Run("test a currency pair with no base/quote", func(t *testing.T) {
		ok.On("GetIDForCurrencyPair", mock.Anything, connecttypes.CurrencyPair{}).Return(uint64(0), false).Once()

		// expect an error when querying for a currency-pair not in module-state
		id, err := strategy.ID(ctx, connecttypes.CurrencyPair{})
		require.Error(t, err)
		require.Equal(t, uint64(0), id) // not equal to 0 because we have a delimiter.
	})

	t.Run("test equality of hashing", func(t *testing.T) {
		ok.On("GetIDForCurrencyPair", mock.Anything, btcusd).Return(uint64(0), true).Twice()

		id1, err := strategy.ID(ctx, btcusd)
		require.NoError(t, err)

		id2, err := strategy.ID(ctx, btcusd)
		require.NoError(t, err)

		require.Equal(t, id1, id2)
	})
}

func TestHashCurrencyPairStrategyFromID(t *testing.T) {
	ok := mocks.NewOracleKeeper(t)
	ctx := sdk.Context{}
	strategy := strategies.NewHashCurrencyPairStrategy(ok)

	t.Run("test getting currency pair for currency pair that does not exist in state", func(t *testing.T) {
		ok.On("GetAllCurrencyPairs", ctx).Return([]connecttypes.CurrencyPair{}).Once()

		id, err := strategies.CurrencyPairToHashID(btcusd.String())
		require.NoError(t, err)

		_, err = strategy.FromID(ctx, id)
		require.Error(t, err)
	})

	t.Run("test getting currency pair for currency pair that exists in state (no cache)", func(t *testing.T) {
		ok.On("GetAllCurrencyPairs", ctx).Return([]connecttypes.CurrencyPair{btcusd}).Once()

		id, err := strategies.CurrencyPairToHashID(btcusd.String())
		require.NoError(t, err)

		cp, err := strategy.FromID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, btcusd, cp)
	})

	t.Run("test getting currency pair for currency pair that exists in state (with cache)", func(t *testing.T) {
		id, err := strategies.CurrencyPairToHashID(btcusd.String())
		require.NoError(t, err)

		cp, err := strategy.FromID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, btcusd, cp)
	})
}
