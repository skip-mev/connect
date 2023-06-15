package keeper_test

import (
	"testing"
	"time"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/x/oracle/keeper"
	"github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	oracleKeeper keeper.Keeper
	key          storetypes.StoreKey
	ctx          sdk.Context
}

func (s *KeeperTestSuite) SetupTest() {
	s.key = storetypes.NewKVStoreKey(types.StoreKey)
	s.oracleKeeper = keeper.NewKeeper(s.key)
	s.ctx = testutil.DefaultContext(s.key, storetypes.NewTransientStoreKey("transient_key"))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) TestSetPriceForCurrencyPair() {
	tcs := []struct {
		name       string
		cp         types.CurrencyPair
		price      types.QuotePrice
		expectPass bool
	}{
		{
			"if the currency pair is correctly formatted - pass",
			types.CurrencyPair{
				Base:  "AA",
				Quote: "BB",
			},
			types.QuotePrice{
				BlockTimestamp: time.Now(),
				BlockHeight:    100,
				Price:          sdk.NewInt(100),
			},
			true,
		},
	}

	for _, tc := range tcs {
		// set the price to state
		err := s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, tc.cp, tc.price)

		switch tc.expectPass {
		case true:
			// expect the quote price to be written to state for the currency pair
			assert.Nil(s.T(), err)
			// expect that we can retrieve the QuotePrice for the currency pair
			qp, err := s.oracleKeeper.GetPriceForCurrencyPair(s.ctx, tc.cp)
			assert.Nil(s.T(), err)
			// check equality
			checkQuotePriceEqual(s.T(), qp, tc.price)
		default:
			// check that there was a failure setting the currency pair
			assert.NotNil(s.T(), err)
		}
	}
}

func (s *KeeperTestSuite) TestSetPriceIncrementNonce() {
	// insert a cp + qp pair, and check that the nonce is zero
	cp := types.CurrencyPair{
		Base:  "AA",
		Quote: "BB",
	}
	qp := types.QuotePrice{
		Price: sdk.NewInt(100),
	}
	// attempt to get the qp for cp (should fail)
	_, err := s.oracleKeeper.GetPriceWithNonceForCurrencyPair(s.ctx, cp)
	assert.NotNil(s.T(), err)

	// set the qp
	err = s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp, qp)
	assert.Nil(s.T(), err)

	// check that the nonce is zero for the cp
	qpn, err := s.oracleKeeper.GetPriceWithNonceForCurrencyPair(s.ctx, cp)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), qpn.Nonce(), uint64(0))

	// update the qp
	qp.Price = sdk.NewInt(101)
	err = s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp, qp)
	assert.Nil(s.T(), err)

	// get the nonce again, and expect it to have incremented
	qpn, err = s.oracleKeeper.GetPriceWithNonceForCurrencyPair(s.ctx, cp)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), qpn.Nonce(), uint64(1))
}

func checkQuotePriceEqual(t *testing.T, qp1, qp2 types.QuotePrice) {
	assert.Equal(t, qp1.BlockHeight, qp2.BlockHeight)
	assert.Equal(t, qp1.BlockTimestamp.UnixMilli(), qp2.BlockTimestamp.UnixMilli())
	assert.Equal(t, qp1.Price.Int64(), qp2.Price.Int64())
}

func (s *KeeperTestSuite) TestGetAllTickers() {
	// insert multiple currency pairs
	cp1 := types.CurrencyPair{
		Base:  "AA",
		Quote: "BB",
	}
	qp1 := types.QuotePrice{
		Price: sdk.NewInt(100),
	}
	cp2 := types.CurrencyPair{
		Base:  "CC",
		Quote: "DD",
	}
	qp2 := types.QuotePrice{
		Price: sdk.NewInt(120),
	}

	// insert
	assert.Nil(s.T(), s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp1, qp1))
	assert.Nil(s.T(), s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp2, qp2))

	// get all tickers
	expectedTickers := map[string]struct{}{"AA/BB": {}, "CC/DD": {}}
	tickers := s.oracleKeeper.GetAllTickers(s.ctx)

	assert.Equal(s.T(), len(tickers), 2)

	// check for inclusion
	for _, ticker := range tickers {
		ts := ticker.ToString()
		_, ok := expectedTickers[ts]
		assert.True(s.T(), ok)
	}
}
