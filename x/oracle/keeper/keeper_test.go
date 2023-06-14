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
	key storetypes.StoreKey
	ctx sdk.Context
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
	tcs := []struct{
		name string
		cp types.CurrencyPair
		price types.QuotePrice
		expectPass bool
	}{
		{
			"if the currency pair is incorrectly formatted - fail",
			types.CurrencyPair{
				Base: "AA",
				Quote: "aB",
			},
			types.QuotePrice{},
			false,
		},
		{
			"if the currency pair is correctly formatted - pass",
			types.CurrencyPair{
				Base: "AA",
				Quote: "BB",
			},
			types.QuotePrice{
				BlockTimestamp: time.Now(),
				BlockHeight: 100,
				Price: sdk.NewInt(100),
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

func checkQuotePriceEqual(t *testing.T, qp1, qp2 types.QuotePrice) {
	assert.Equal(t, qp1.BlockHeight, qp2.BlockHeight)
	assert.Equal(t, qp1.BlockTimestamp.UnixMilli(), qp2.BlockTimestamp.UnixMilli())
	assert.Equal(t, qp1.Price.Int64(), qp2.Price.Int64())
}

func (s *KeeperTestSuite) TestGetAllTickers() {
	// insert multiple currency pairs
	cp1 := types.CurrencyPair{
		Base: "AA",
		Quote: "BB",
	}
	qp1 := types.QuotePrice{
		Price: sdk.NewInt(100),
	}
	cp2 := types.CurrencyPair{
		Base: "CC",
		Quote: "DD",
	}
	qp2 := types.QuotePrice{
		Price: sdk.NewInt(120),
	}

	// insert
	assert.Nil(s.T(), s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp1, qp1))
	assert.Nil(s.T(), s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp2, qp2))
	
	// get all tickers
	expectedTickers := map[string]struct{}{"AA/BB":{}, "CC/DD":{}}
	tickers, err := s.oracleKeeper.GetAllTickers(s.ctx)
	assert.Nil(s.T(), err)
	
	// check for inclusion
	for _, ticker := range tickers {
		ts := ticker.ToString()
		_, ok := expectedTickers[ts]
		assert.True(s.T(), ok)
	}
}