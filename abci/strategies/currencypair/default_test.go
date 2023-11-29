package currencypair_test

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	strategies "github.com/skip-mev/slinky/abci/strategies/currencypair"
	mocks "github.com/skip-mev/slinky/abci/strategies/currencypair/mocks"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

func TestDefaultCurrencyPairStrategyID(t *testing.T) {
	ok := mocks.NewOracleKeeper(t)

	ctx := sdk.Context{}

	strategy := strategies.NewDefaultCurrencyPairStrategy(ok)

	// test that the strategy returns IDs from the oracle module
	t.Run("test getting ids with two currency-pairs in module-state", func(t *testing.T) {
		// expect the first currency-pair to have ID 0
		ok.On("GetIDForCurrencyPair", ctx, oracletypes.NewCurrencyPair("BTC", "USD")).Return(uint64(0), true)
		id, err := strategy.ID(ctx, oracletypes.NewCurrencyPair("BTC", "USD"))
		require.NoError(t, err)
		require.Equal(t, uint64(0), id)

		// expect the second currency-pair to have ID 1
		ok.On("GetIDForCurrencyPair", ctx, oracletypes.NewCurrencyPair("USD", "ETH")).Return(uint64(1), true)
		id, err = strategy.ID(ctx, oracletypes.NewCurrencyPair("USD", "ETH"))
		require.NoError(t, err)
		require.Equal(t, uint64(1), id)
	})

	// test that if a currency-pair does not have an ID w/ x/oracle, a failure is returned
	t.Run("expect error when currency-pair not found in module-state", func(t *testing.T) {
		// expect an error when querying for a currency-pair not in module-state
		ok.On("GetIDForCurrencyPair", ctx, oracletypes.NewCurrencyPair("ETH", "BTC")).Return(uint64(0), false)
		_, err := strategy.ID(ctx, oracletypes.NewCurrencyPair("ETH", "BTC"))
		require.Error(t, err)
	})
}

func TestDefaultCurrencyPairStrategyFromID(t *testing.T) {
	ok := mocks.NewOracleKeeper(t)

	ctx := sdk.Context{}

	strategy := strategies.NewDefaultCurrencyPairStrategy(ok)

	// test that the strategy returns IDs from the oracle module
	t.Run("test getting ids with two currency-pairs in module-state", func(t *testing.T) {
		// expect the first currency-pair to have ID 0
		ok.On("GetCurrencyPairFromID", ctx, uint64(0)).Return(oracletypes.NewCurrencyPair("BTC", "USD"), true)
		cp, err := strategy.FromID(ctx, uint64(0))
		require.NoError(t, err)
		require.Equal(t, oracletypes.NewCurrencyPair("BTC", "USD"), cp)

		// expect the second currency-pair to have ID 1
		ok.On("GetCurrencyPairFromID", ctx, uint64(1)).Return(oracletypes.NewCurrencyPair("USD", "ETH"), true)
		cp, err = strategy.FromID(ctx, uint64(1))
		require.NoError(t, err)
		require.Equal(t, oracletypes.NewCurrencyPair("USD", "ETH"), cp)
	})

	// test that if a currency-pair does not have an ID w/ x/oracle, a failure is returned
	t.Run("expect error when currency-pair not found in module-state", func(t *testing.T) {
		// expect an error when querying for a currency-pair not in module-state
		ok.On("GetCurrencyPairFromID", ctx, uint64(2)).Return(oracletypes.CurrencyPair{}, false)
		_, err := strategy.FromID(ctx, uint64(2))
		require.Error(t, err)
	})
}

func TestDefaultCurrencyPairStrategyGetEncodedPrice(t *testing.T) {
	ok := mocks.NewOracleKeeper(t)

	ctx := sdk.Context{}

	strategy := strategies.NewDefaultCurrencyPairStrategy(ok)

	cp := oracletypes.NewCurrencyPair("BTC", "USD")
	t.Run("can encode a positive price", func(t *testing.T) {
		price := big.NewInt(100)
		encodedPrice, err := strategy.GetEncodedPrice(ctx, cp, price)
		require.NoError(t, err)

		decodedPrice, err := strategy.GetDecodedPrice(ctx, cp, encodedPrice)
		require.NoError(t, err)
		require.Equal(t, price, decodedPrice)
	})

	t.Run("cannot encode a negative price", func(t *testing.T) {
		price := big.NewInt(-100)
		_, err := strategy.GetEncodedPrice(ctx, cp, price)
		require.Error(t, err)

		bz, err := price.GobEncode()
		require.NoError(t, err)

		_, err = strategy.GetDecodedPrice(ctx, cp, bz)
		require.Error(t, err)
	})

	t.Run("errors when decoding a negative price", func(t *testing.T) {
		price := big.NewInt(-100)
		bz, err := price.GobEncode()
		require.NoError(t, err)

		_, err = strategy.GetDecodedPrice(ctx, cp, bz)
		require.Error(t, err)
	})
}
