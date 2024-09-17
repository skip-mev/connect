package currencypair_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/abci/strategies/currencypair"
	mocks "github.com/skip-mev/connect/v2/abci/strategies/currencypair/mocks"
	"github.com/skip-mev/connect/v2/abci/testutils"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
)

func TestDeltaCurrencyPairStrategyGetEncodedPrice(t *testing.T) {
	cp := connecttypes.NewCurrencyPair("BTC", "USD")

	t.Run("price does not exist in state, delta is final price", func(t *testing.T) {
		ok := mocks.NewOracleKeeper(t)
		ctx := testutils.CreateBaseSDKContext(t)
		strategy := currencypair.NewDeltaCurrencyPairStrategy(ok)

		ok.On(
			"GetPriceForCurrencyPair",
			mock.Anything,
			cp,
		).Return(oracletypes.QuotePrice{}, oracletypes.QuotePriceNotExistError{})

		price := big.NewInt(100)
		encodedPrice, err := strategy.GetEncodedPrice(ctx, cp, price)
		require.NoError(t, err)

		// the encoded price should be the delta price
		var deltaPrice big.Int
		err = deltaPrice.GobDecode(encodedPrice)
		require.NoError(t, err)
		require.Equal(t, price, &deltaPrice)
	})

	t.Run("price exists in state, inputted price is smaller so delta should be negative", func(t *testing.T) {
		ok := mocks.NewOracleKeeper(t)
		ctx := testutils.CreateBaseSDKContext(t)
		strategy := currencypair.NewDeltaCurrencyPairStrategy(ok)

		onChainPrice := math.NewInt(100)
		ok.On(
			"GetPriceForCurrencyPair",
			mock.Anything,
			cp,
		).Return(oracletypes.QuotePrice{Price: onChainPrice}, nil)

		price := big.NewInt(80)
		encodedPrice, err := strategy.GetEncodedPrice(ctx, cp, price)
		require.NoError(t, err)

		// the encoded price should be the delta price
		var deltaPrice big.Int
		err = deltaPrice.GobDecode(encodedPrice)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(-20), &deltaPrice)
	})

	t.Run("price exists in state, inputted price is larger so delta should be positive", func(t *testing.T) {
		ok := mocks.NewOracleKeeper(t)
		ctx := testutils.CreateBaseSDKContext(t)
		strategy := currencypair.NewDeltaCurrencyPairStrategy(ok)

		onChainPrice := math.NewInt(100)
		ok.On(
			"GetPriceForCurrencyPair",
			mock.Anything,
			cp,
		).Return(oracletypes.QuotePrice{Price: onChainPrice}, nil)

		price := big.NewInt(120)
		encodedPrice, err := strategy.GetEncodedPrice(ctx, cp, price)
		require.NoError(t, err)

		// the encoded price should be the delta price
		var deltaPrice big.Int
		err = deltaPrice.GobDecode(encodedPrice)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(20), &deltaPrice)
	})

	t.Run("price exists in state, inputted price is equal so delta should be zero", func(t *testing.T) {
		ok := mocks.NewOracleKeeper(t)
		ctx := testutils.CreateBaseSDKContext(t)
		strategy := currencypair.NewDeltaCurrencyPairStrategy(ok)

		onChainPrice := math.NewInt(100)
		ok.On(
			"GetPriceForCurrencyPair",
			mock.Anything,
			cp,
		).Return(oracletypes.QuotePrice{Price: onChainPrice}, nil)

		price := big.NewInt(100)
		encodedPrice, err := strategy.GetEncodedPrice(ctx, cp, price)
		require.NoError(t, err)

		// the encoded price should be the delta price
		var deltaPrice big.Int
		err = deltaPrice.GobDecode(encodedPrice)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(0), &deltaPrice)
	})

	t.Run("price cache works for several calls to the same height", func(t *testing.T) {
		ok := mocks.NewOracleKeeper(t)
		ctx := testutils.CreateBaseSDKContext(t)
		strategy := currencypair.NewDeltaCurrencyPairStrategy(ok)

		onChainPrice := math.NewInt(100)
		ok.On(
			"GetPriceForCurrencyPair",
			mock.Anything,
			cp,
		).Return(oracletypes.QuotePrice{Price: onChainPrice}, nil)

		price := big.NewInt(120)
		encodedPrice, err := strategy.GetEncodedPrice(ctx, cp, price)
		require.NoError(t, err)

		// the encoded price should be the delta price
		var deltaPrice big.Int
		err = deltaPrice.GobDecode(encodedPrice)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(20), &deltaPrice)

		// the second call should return the same encoded price
		encodedPrice, err = strategy.GetEncodedPrice(ctx, cp, price)
		require.NoError(t, err)
		err = deltaPrice.GobDecode(encodedPrice)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(20), &deltaPrice)
	})

	t.Run("price cache is cleared when height changes", func(t *testing.T) {
		ok := mocks.NewOracleKeeper(t)
		ctx := testutils.CreateBaseSDKContext(t)
		strategy := currencypair.NewDeltaCurrencyPairStrategy(ok)

		onChainPrice := math.NewInt(100)
		ok.On(
			"GetPriceForCurrencyPair",
			ctx,
			cp,
		).Return(oracletypes.QuotePrice{Price: onChainPrice}, nil).Once()

		price := big.NewInt(120)
		encodedPrice, err := strategy.GetEncodedPrice(ctx, cp, price)
		require.NoError(t, err)

		// the encoded price should be the delta price
		var deltaPrice big.Int
		err = deltaPrice.GobDecode(encodedPrice)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(20), &deltaPrice)

		// update the height to clear the cache
		ctx = ctx.WithBlockHeight(5)
		updatedOnChainPrice := math.NewInt(200)
		ok.On(
			"GetPriceForCurrencyPair",
			ctx,
			cp,
		).Return(oracletypes.QuotePrice{Price: updatedOnChainPrice}, nil).Once()

		price = big.NewInt(50)
		encodedPrice, err = strategy.GetEncodedPrice(ctx, cp, price)
		require.NoError(t, err)

		// the encoded price should be the delta price
		err = deltaPrice.GobDecode(encodedPrice)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(-150), &deltaPrice)
	})

	t.Run("error when the price is negative", func(t *testing.T) {
		ok := mocks.NewOracleKeeper(t)
		ctx := testutils.CreateBaseSDKContext(t)
		strategy := currencypair.NewDeltaCurrencyPairStrategy(ok)

		price := big.NewInt(-120)
		_, err := strategy.GetEncodedPrice(ctx, cp, price)
		require.Error(t, err)
	})
}

func TestDeltaCurrencyPairStrategyGetDecodedPrice(t *testing.T) {
	t.Run("price does not exist in state, delta is final price", func(t *testing.T) {
		ok := mocks.NewOracleKeeper(t)
		ctx := testutils.CreateBaseSDKContext(t)
		strategy := currencypair.NewDeltaCurrencyPairStrategy(ok)

		cp := connecttypes.NewCurrencyPair("BTC", "USD")

		ok.On(
			"GetPriceForCurrencyPair",
			mock.Anything,
			cp,
		).Return(oracletypes.QuotePrice{}, oracletypes.QuotePriceNotExistError{})

		price := big.NewInt(100)
		priceBytes, err := price.GobEncode()
		require.NoError(t, err)

		decodedPrice, err := strategy.GetDecodedPrice(ctx, cp, priceBytes)
		require.NoError(t, err)
		require.Equal(t, price, decodedPrice)
	})

	t.Run("price exists in state, negative delta with increase price returned", func(t *testing.T) {
		ok := mocks.NewOracleKeeper(t)
		ctx := testutils.CreateBaseSDKContext(t)
		strategy := currencypair.NewDeltaCurrencyPairStrategy(ok)

		cp := connecttypes.NewCurrencyPair("BTC", "USD")

		onChainPrice := math.NewInt(100)
		ok.On(
			"GetPriceForCurrencyPair",
			mock.Anything,
			cp,
		).Return(oracletypes.QuotePrice{Price: onChainPrice}, nil)

		delta := big.NewInt(-20)
		priceBytes, err := delta.GobEncode()
		require.NoError(t, err)

		decodedPrice, err := strategy.GetDecodedPrice(ctx, cp, priceBytes)
		require.NoError(t, err)

		expectedPrice := big.NewInt(80)
		require.Equal(t, expectedPrice, decodedPrice)
	})

	t.Run("price exists in state, positive delta with increase price returned", func(t *testing.T) {
		ok := mocks.NewOracleKeeper(t)
		ctx := testutils.CreateBaseSDKContext(t)
		strategy := currencypair.NewDeltaCurrencyPairStrategy(ok)

		cp := connecttypes.NewCurrencyPair("BTC", "USD")

		onChainPrice := math.NewInt(100)
		ok.On(
			"GetPriceForCurrencyPair",
			mock.Anything,
			cp,
		).Return(oracletypes.QuotePrice{Price: onChainPrice}, nil)

		delta := big.NewInt(20)
		priceBytes, err := delta.GobEncode()
		require.NoError(t, err)

		decodedPrice, err := strategy.GetDecodedPrice(ctx, cp, priceBytes)
		require.NoError(t, err)

		expectedPrice := big.NewInt(120)
		require.Equal(t, expectedPrice, decodedPrice)
	})

	t.Run("price exists in state, zero delta with no increase price returned", func(t *testing.T) {
		ok := mocks.NewOracleKeeper(t)
		ctx := testutils.CreateBaseSDKContext(t)
		strategy := currencypair.NewDeltaCurrencyPairStrategy(ok)

		cp := connecttypes.NewCurrencyPair("BTC", "USD")

		onChainPrice := math.NewInt(100)
		ok.On(
			"GetPriceForCurrencyPair",
			mock.Anything,
			cp,
		).Return(oracletypes.QuotePrice{Price: onChainPrice}, nil)

		delta := big.NewInt(0)
		priceBytes, err := delta.GobEncode()
		require.NoError(t, err)

		decodedPrice, err := strategy.GetDecodedPrice(ctx, cp, priceBytes)
		require.NoError(t, err)

		expectedPrice := big.NewInt(100)
		require.Equal(t, expectedPrice, decodedPrice)
	})

	t.Run("negative delta with negative price returns error", func(t *testing.T) {
		ok := mocks.NewOracleKeeper(t)
		ctx := testutils.CreateBaseSDKContext(t)
		strategy := currencypair.NewDeltaCurrencyPairStrategy(ok)

		cp := connecttypes.NewCurrencyPair("BTC", "USD")

		onChainPrice := math.NewInt(100)
		ok.On(
			"GetPriceForCurrencyPair",
			mock.Anything,
			cp,
		).Return(oracletypes.QuotePrice{Price: onChainPrice}, nil)

		delta := big.NewInt(-120)
		priceBytes, err := delta.GobEncode()
		require.NoError(t, err)

		_, err = strategy.GetDecodedPrice(ctx, cp, priceBytes)
		require.Error(t, err)
	})

	t.Run("price cache works for several calls to the same height", func(t *testing.T) {
		ok := mocks.NewOracleKeeper(t)
		ctx := testutils.CreateBaseSDKContext(t)
		strategy := currencypair.NewDeltaCurrencyPairStrategy(ok)
		cp := connecttypes.NewCurrencyPair("BTC", "USD")

		onChainPrice := math.NewInt(100)
		ok.On(
			"GetPriceForCurrencyPair",
			ctx,
			cp,
		).Return(oracletypes.QuotePrice{Price: onChainPrice}, nil).Once()

		delta := big.NewInt(20)
		priceBytes, err := delta.GobEncode()
		require.NoError(t, err)

		decodedPrice, err := strategy.GetDecodedPrice(ctx, cp, priceBytes)
		require.NoError(t, err)

		expectedPrice := big.NewInt(120)
		require.Equal(t, expectedPrice, decodedPrice)

		// the second call should return the same decoded price
		decodedPrice, err = strategy.GetDecodedPrice(ctx, cp, priceBytes)
		require.NoError(t, err)
		require.Equal(t, expectedPrice, decodedPrice)
	})

	t.Run("price cache is cleared when height changes", func(t *testing.T) {
		ok := mocks.NewOracleKeeper(t)
		ctx := testutils.CreateBaseSDKContext(t)
		strategy := currencypair.NewDeltaCurrencyPairStrategy(ok)
		cp := connecttypes.NewCurrencyPair("BTC", "USD")

		onChainPrice := math.NewInt(100)
		ok.On(
			"GetPriceForCurrencyPair",
			ctx,
			cp,
		).Return(oracletypes.QuotePrice{Price: onChainPrice}, nil).Once()

		delta := big.NewInt(20)
		priceBytes, err := delta.GobEncode()
		require.NoError(t, err)

		decodedPrice, err := strategy.GetDecodedPrice(ctx, cp, priceBytes)
		require.NoError(t, err)

		expectedPrice := big.NewInt(120)
		require.Equal(t, expectedPrice, decodedPrice)

		// update the height to clear the cache
		ctx = ctx.WithBlockHeight(5)
		updatedOnChainPrice := math.NewInt(200)
		ok.On(
			"GetPriceForCurrencyPair",
			ctx,
			cp,
		).Return(oracletypes.QuotePrice{Price: updatedOnChainPrice}, nil).Once()

		delta = big.NewInt(50)
		priceBytes, err = delta.GobEncode()
		require.NoError(t, err)

		decodedPrice, err = strategy.GetDecodedPrice(ctx, cp, priceBytes)
		require.NoError(t, err)

		expectedPrice = big.NewInt(250)
		require.Equal(t, expectedPrice, decodedPrice)
	})
}
