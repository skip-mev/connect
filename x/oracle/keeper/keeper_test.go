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

	"github.com/skip-mev/slinky/x/oracle/keeper"
	"github.com/skip-mev/slinky/x/oracle/types"
)

const (
	moduleAuth = "authority"
)

var moduleAuthAddr = sdk.AccAddress([]byte(moduleAuth))

type KeeperTestSuite struct {
	suite.Suite

	oracleKeeper keeper.Keeper
	ctx          sdk.Context
}

func (s *KeeperTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	ss := runtime.NewKVStoreService(key)
	encCfg := moduletestutil.MakeTestEncodingConfig()
	s.oracleKeeper = keeper.NewKeeper(ss, encCfg.Codec, moduleAuthAddr)
	s.ctx = testutil.DefaultContext(key, storetypes.NewTransientStoreKey("transient_key"))
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
	cp := types.CurrencyPair{
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
	cp1 := types.CurrencyPair{
		Base:  "AA",
		Quote: "BB",
	}
	qp1 := types.QuotePrice{
		Price: sdkmath.NewInt(100),
	}
	cp2 := types.CurrencyPair{
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
	cp := types.CurrencyPair{
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
	cp1 := types.CurrencyPair{
		Base:  "PAIR",
		Quote: "1",
	}

	cp2 := types.CurrencyPair{
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
		cp3 := types.CurrencyPair{
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
