package currencypair_test

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	strategies "github.com/skip-mev/slinky/abci/strategies/currencypair"
	"github.com/skip-mev/slinky/abci/strategies/currencypair/mocks"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
)

var (
	btcusd = slinkytypes.NewCurrencyPair("BTC", "USD")
	usdeth = slinkytypes.NewCurrencyPair("USD", "ETH")
	ethbtc = slinkytypes.NewCurrencyPair("ETH", "BTC")
)

func TestDefaultCurrencyPairStrategyID(t *testing.T) {
	ok := mocks.NewOracleKeeper(t)
	ctx := sdk.Context{}
	strategy := strategies.NewDefaultCurrencyPairStrategy(ok)

	// test that the strategy returns IDs from the oracle module
	t.Run("test getting ids with two currency-pairs in module-state", func(t *testing.T) {
		// expect the first currency-pair to have ID 0
		ok.On("GetIDForCurrencyPair", ctx, btcusd).Return(uint64(0), true).Once()
		id, err := strategy.ID(ctx, btcusd)
		require.NoError(t, err)
		require.Equal(t, uint64(0), id)

		// expect the second currency-pair to have ID 1
		ok.On("GetIDForCurrencyPair", ctx, usdeth).Return(uint64(1), true).Once()
		id, err = strategy.ID(ctx, usdeth)
		require.NoError(t, err)
		require.Equal(t, uint64(1), id)
	})

	// test that if a currency-pair does not have an ID w/ x/oracle, a failure is returned
	t.Run("expect error when currency-pair not found in module-state", func(t *testing.T) {
		// expect an error when querying for a currency-pair not in module-state
		ok.On("GetIDForCurrencyPair", ctx, ethbtc).Return(uint64(0), false).Once()
		_, err := strategy.ID(ctx, ethbtc)
		require.Error(t, err)
	})
}

func TestDefaultCurrencyPairStrategyFromID(t *testing.T) {
	ok := mocks.NewOracleKeeper(t)
	ctx := sdk.Context{}
	strategy := strategies.NewDefaultCurrencyPairStrategy(ok)

	// test that the strategy returns IDs from the oracle module and then will come from the cache
	t.Run("test getting ids with two currency-pairs in module-state", func(t *testing.T) {
		// expect the first currency-pair to have ID 0
		ok.On("GetCurrencyPairFromID", ctx, uint64(0)).Return(btcusd, true).Once()
		cp, err := strategy.FromID(ctx, uint64(0))
		require.NoError(t, err)
		require.Equal(t, btcusd, cp)

		// expect the second currency-pair to have ID 1
		ok.On("GetCurrencyPairFromID", ctx, uint64(1)).Return(usdeth, true).Once()
		cp, err = strategy.FromID(ctx, uint64(1))
		require.NoError(t, err)
		require.Equal(t, usdeth, cp)

		// call ID to populate the cache
		// expect the first currency-pair to have ID 0
		ok.On("GetIDForCurrencyPair", ctx, btcusd).Return(uint64(0), true).Once()
		id, err := strategy.ID(ctx, btcusd)
		require.NoError(t, err)
		require.Equal(t, uint64(0), id)

		// expect the second currency-pair to have ID 1
		ok.On("GetIDForCurrencyPair", ctx, usdeth).Return(uint64(1), true).Once()
		id, err = strategy.ID(ctx, usdeth)
		require.NoError(t, err)
		require.Equal(t, uint64(1), id)

		// expect the first currency-pair to have ID 0
		cp, err = strategy.FromID(ctx, uint64(0))
		require.NoError(t, err)
		require.Equal(t, btcusd, cp)

		// expect the second currency-pair to have ID 1
		cp, err = strategy.FromID(ctx, uint64(1))
		require.NoError(t, err)
		require.Equal(t, usdeth, cp)
	})

	// test that if a currency-pair does not have an ID w/ x/oracle, a failure is returned
	t.Run("expect error when currency-pair not found in module-state", func(t *testing.T) {
		// expect an error when querying for a currency-pair not in module-state
		ok.On("GetCurrencyPairFromID", ctx, uint64(2)).Return(slinkytypes.CurrencyPair{}, false).Once()
		_, err := strategy.FromID(ctx, uint64(2))
		require.Error(t, err)
	})
}

func TestDefaultCurrencyPairStrategyGetEncodedPrice(t *testing.T) {
	ok := mocks.NewOracleKeeper(t)
	ctx := sdk.Context{}
	strategy := strategies.NewDefaultCurrencyPairStrategy(ok)

	cp := btcusd
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
