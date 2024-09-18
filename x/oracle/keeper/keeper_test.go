package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/oracle/keeper"
	"github.com/skip-mev/connect/v2/x/oracle/types"
	"github.com/skip-mev/connect/v2/x/oracle/types/mocks"
)

const (
	moduleAuth = "authority"
)

var moduleAuthAddr = sdk.AccAddress(moduleAuth)

type KeeperTestSuite struct {
	suite.Suite

	oracleKeeper        keeper.Keeper
	mockMarketMapKeeper *mocks.MarketMapKeeper
	ctx                 sdk.Context
}

func (s *KeeperTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	ss := runtime.NewKVStoreService(key)
	encCfg := moduletestutil.MakeTestEncodingConfig()
	s.mockMarketMapKeeper = mocks.NewMarketMapKeeper(s.T())
	s.oracleKeeper = keeper.NewKeeper(ss, encCfg.Codec, s.mockMarketMapKeeper, moduleAuthAddr)
	s.ctx = testutil.DefaultContext(key, storetypes.NewTransientStoreKey("transient_key"))

	s.Require().NotPanics(func() {
		s.oracleKeeper.InitGenesis(s.ctx, *types.DefaultGenesisState())
	})
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupWithNoMMKeeper() {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	ss := runtime.NewKVStoreService(key)
	encCfg := moduletestutil.MakeTestEncodingConfig()
	s.oracleKeeper = keeper.NewKeeper(ss, encCfg.Codec, nil, moduleAuthAddr)
	s.ctx = testutil.DefaultContext(key, storetypes.NewTransientStoreKey("transient_key"))
}

func (s *KeeperTestSuite) TestSetPriceForCurrencyPair() {
	tcs := []struct {
		name       string
		cp         connecttypes.CurrencyPair
		price      types.QuotePrice
		expectPass bool
	}{
		{
			"if the currency pair is correctly formatted - pass",
			connecttypes.CurrencyPair{
				Base:  "AA",
				Quote: "BB",
			},
			types.QuotePrice{
				BlockTimestamp: time.Now(),
				BlockHeight:    100,
				Price:          sdkmath.NewInt(100),
			},
			true,
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			// set the price to state
			err := s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, tc.cp, tc.price)

			switch tc.expectPass {
			case true:
				// expect the quote price to be written to state for the currency pair
				s.Require().Nil(err)
				// expect that we can retrieve the QuotePrice for the currency pair
				qp, err := s.oracleKeeper.GetPriceForCurrencyPair(s.ctx, tc.cp)
				s.Require().Nil(err)
				// check equality
				checkQuotePriceEqual(s.T(), qp, tc.price)
			default:
				// check that there was a failure setting the currency pair
				s.Require().NotNil(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestSetPriceIncrementNonce() {
	// insert a cp + qp pair, and check that the nonce is zero
	cp := connecttypes.CurrencyPair{
		Base:  "AA",
		Quote: "BB",
	}
	qp := types.QuotePrice{
		Price: sdkmath.NewInt(100),
	}
	// attempt to get the qp for cp (should fail)
	_, err := s.oracleKeeper.GetPriceWithNonceForCurrencyPair(s.ctx, cp)
	s.Require().NotNil(err)

	// attempt to get the nonce for the cp (should fail)
	_, err = s.oracleKeeper.GetNonceForCurrencyPair(s.ctx, cp)
	s.Require().NotNil(err)

	// set the qp
	err = s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp, qp)
	s.Require().Nil(err)

	// check that the nonce is zero for the cp
	qpn, err := s.oracleKeeper.GetPriceWithNonceForCurrencyPair(s.ctx, cp)
	s.Require().Nil(err)
	require.Equal(s.T(), qpn.Nonce(), uint64(0))

	// update the qp
	qp.Price = sdkmath.NewInt(101)
	err = s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp, qp)
	s.Require().Nil(err)

	// get the nonce again, and expect it to have incremented
	qpn, err = s.oracleKeeper.GetPriceWithNonceForCurrencyPair(s.ctx, cp)
	s.Require().Nil(err)
	s.Require().Equal(qpn.Nonce(), uint64(1))
}

func checkQuotePriceEqual(t *testing.T, qp1, qp2 types.QuotePrice) {
	t.Helper()

	require.Equal(t, qp1.BlockHeight, qp2.BlockHeight)
	require.Equal(t, qp1.BlockTimestamp.UnixMilli(), qp2.BlockTimestamp.UnixMilli())
	require.Equal(t, qp1.Price.Int64(), qp2.Price.Int64())
}

func (s *KeeperTestSuite) TestGetAllCPs() {
	// insert multiple currency pairs
	cp1 := connecttypes.CurrencyPair{
		Base:  "AA",
		Quote: "BB",
	}
	qp1 := types.QuotePrice{
		Price: sdkmath.NewInt(100),
	}
	cp2 := connecttypes.CurrencyPair{
		Base:  "CC",
		Quote: "DD",
	}
	qp2 := types.QuotePrice{
		Price: sdkmath.NewInt(120),
	}

	// insert
	s.Require().Nil(s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp1, qp1))
	s.Require().Nil(s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp2, qp2))

	// get all cps
	expectedCurrencyPairs := map[string]struct{}{"AA/BB": {}, "CC/DD": {}}
	tickers := s.oracleKeeper.GetAllCurrencyPairs(s.ctx)

	s.Require().Equal(len(tickers), 2)

	// check for inclusion
	for _, ticker := range tickers {
		ts := ticker.String()
		_, ok := expectedCurrencyPairs[ts]
		s.Require().True(ok)
	}
}

func (s *KeeperTestSuite) TestCreateCurrencyPair() {
	cp := connecttypes.CurrencyPair{
		Base:  "NEW",
		Quote: "PAIR",
	}
	s.Run("creating a new currency-pair initializes correctly", func() {
		// create the currency pair
		err := s.oracleKeeper.CreateCurrencyPair(s.ctx, cp)
		s.Require().Nil(err)

		// check that the currency pair exists
		nonce, err := s.oracleKeeper.GetNonceForCurrencyPair(s.ctx, cp)
		s.Require().Nil(err)
		s.Require().Equal(nonce, uint64(0))

		// check that the set of all cps includes the currency-pair
		cps := s.oracleKeeper.GetAllCurrencyPairs(s.ctx)

		var found bool
		for _, cp := range cps {
			if cp.String() == "NEW/PAIR" {
				found = true
				break
			}
		}
		s.Require().True(found)
	})

	s.Run("creating a currency-pair twice fails", func() {
		err := s.oracleKeeper.CreateCurrencyPair(s.ctx, cp)
		s.Require().Equal(err.Error(), types.NewCurrencyPairAlreadyExistsError(cp).Error())
	})
}

func (s *KeeperTestSuite) TestIDForCurrencyPair() {
	cp1 := connecttypes.CurrencyPair{
		Base:  "PAIR",
		Quote: "1",
	}

	cp2 := connecttypes.CurrencyPair{
		Base:  "PAIR",
		Quote: "2",
	}

	s.Run("test setting ids for currency pairs", func() {
		s.Require().Nil(s.oracleKeeper.CreateCurrencyPair(s.ctx, cp1))

		// get the id for the currency-pair
		id, ok := s.oracleKeeper.GetIDForCurrencyPair(s.ctx, cp1)
		s.Require().True(ok)

		// set the next currency-pair
		s.Require().Nil(s.oracleKeeper.CreateCurrencyPair(s.ctx, cp2))

		// get the id for the currency-pair
		id2, ok := s.oracleKeeper.GetIDForCurrencyPair(s.ctx, cp2)
		s.Require().True(ok)

		// check that the ids are different
		s.Require().Equal(id+1, id2)
	})

	s.Run("test getting ids for currency-pairs", func() {
		// get the id for the currency-pair
		id, ok := s.oracleKeeper.GetIDForCurrencyPair(s.ctx, cp1)
		s.Require().True(ok)

		// get the currency-pair for the id
		cp, ok := s.oracleKeeper.GetCurrencyPairFromID(s.ctx, id)
		s.Require().True(ok)

		// check that the currency-pair is the same
		s.Require().Equal(cp1, cp)

		// get the id for the currency-pair
		id2, ok := s.oracleKeeper.GetIDForCurrencyPair(s.ctx, cp2)
		s.Require().True(ok)

		// get the currency-pair for the id
		cp, ok = s.oracleKeeper.GetCurrencyPairFromID(s.ctx, id2)
		s.Require().True(ok)

		// check that the currency-pair is the same
		s.Require().Equal(cp2, cp)
	})

	var unusedID uint64
	s.Run("test that removing a currency-pair removes the ID for that currency-pair", func() {
		var ok bool
		unusedID, ok = s.oracleKeeper.GetIDForCurrencyPair(s.ctx, cp2)
		s.Require().True(ok)

		// remove the currency-pair
		s.oracleKeeper.RemoveCurrencyPair(s.ctx, cp2)

		// check that the id is no longer in use
		_, ok = s.oracleKeeper.GetCurrencyPairFromID(s.ctx, unusedID)
		s.Require().False(ok)
	})

	s.Run("insert another currency-pair, and expect that unusedID + 1 is used", func() {
		cp3 := connecttypes.CurrencyPair{
			Base:  "PAIR",
			Quote: "3",
		}

		s.Require().Nil(s.oracleKeeper.CreateCurrencyPair(s.ctx, cp3))

		// get the id for the currency-pair
		id, ok := s.oracleKeeper.GetIDForCurrencyPair(s.ctx, cp3)
		s.Require().True(ok)

		// check that the id is unusedID + 1
		s.Require().Equal(unusedID+1, id)
	})
}

func (s *KeeperTestSuite) TestGetNumRemovedCurrencyPairs() {
	s.Run("get 0 with no state", func() {
		s.SetupTest()

		removes, err := s.oracleKeeper.GetNumRemovedCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(removes, uint64(0))
	})

	s.Run("get 1 with 1 remove", func() {
		s.SetupTest()

		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, connecttypes.CurrencyPair{Base: "test", Quote: "coin1"}))
		s.Require().NoError(s.oracleKeeper.RemoveCurrencyPair(s.ctx, connecttypes.CurrencyPair{Base: "test", Quote: "coin1"}))

		removes, err := s.oracleKeeper.GetNumRemovedCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(removes, uint64(1))
	})

	s.Run("get 2 with 2 removes", func() {
		s.SetupTest()

		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, connecttypes.CurrencyPair{Base: "test", Quote: "coin1"}))
		s.Require().NoError(s.oracleKeeper.RemoveCurrencyPair(s.ctx, connecttypes.CurrencyPair{Base: "test", Quote: "coin1"}))
		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, connecttypes.CurrencyPair{Base: "test", Quote: "coin2"}))
		s.Require().NoError(s.oracleKeeper.RemoveCurrencyPair(s.ctx, connecttypes.CurrencyPair{Base: "test", Quote: "coin2"}))

		removes, err := s.oracleKeeper.GetNumRemovedCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(removes, uint64(2))
	})
}

func (s *KeeperTestSuite) TestGetNumCurrencyPairs() {
	s.Run("get 0 with no state", func() {
		s.SetupTest()

		num, err := s.oracleKeeper.GetNumCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(num, uint64(0))
	})

	s.Run("get 1 with 1 cp", func() {
		s.SetupTest()

		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, connecttypes.CurrencyPair{Base: "test", Quote: "coin1"}))

		cps, err := s.oracleKeeper.GetNumCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(cps, uint64(1))
	})

	s.Run("get 2 with 2 cp", func() {
		s.SetupTest()

		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, connecttypes.CurrencyPair{Base: "test", Quote: "coin1"}))
		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, connecttypes.CurrencyPair{Base: "test", Quote: "coin2"}))

		cps, err := s.oracleKeeper.GetNumCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(cps, uint64(2))
	})
}
