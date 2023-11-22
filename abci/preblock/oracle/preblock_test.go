package oracle_test

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/suite"

	preblock "github.com/skip-mev/slinky/abci/preblock/oracle"
	preblockmath "github.com/skip-mev/slinky/abci/preblock/oracle/math"
	"github.com/skip-mev/slinky/abci/preblock/oracle/math/mocks"
	strategymock "github.com/skip-mev/slinky/abci/strategies/mocks"
	"github.com/skip-mev/slinky/abci/testutils"
	metricmock "github.com/skip-mev/slinky/service/metrics/mocks"
	"github.com/skip-mev/slinky/x/oracle/keeper"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

type PreBlockTestSuite struct {
	suite.Suite

	myVal         sdk.ConsAddress
	ctx           sdk.Context
	currencyPairs []oracletypes.CurrencyPair
	genesis       oracletypes.GenesisState
	key           *storetypes.KVStoreKey
	transientKey  *storetypes.TransientStoreKey
	oracleKeeper  keeper.Keeper
	handler       *preblock.PreBlockHandler
	cpID          *strategymock.CurrencyPairIDStrategy
}

func TestPreBlockTestSuite(t *testing.T) {
	suite.Run(t, new(PreBlockTestSuite))
}

func (s *PreBlockTestSuite) SetupTest() {
	s.myVal = sdk.ConsAddress([]byte("myVal"))

	s.currencyPairs = []oracletypes.CurrencyPair{
		{
			Base:  "BTC",
			Quote: "ETH",
		},
		{
			Base:  "BTC",
			Quote: "USD",
		},
		{
			Base:  "ETH",
			Quote: "USD",
		},
	}
	genesisCPs := []oracletypes.CurrencyPairGenesis{
		{
			CurrencyPair: s.currencyPairs[0],
			Nonce:        0,
			Id:           0,
		},
		{
			CurrencyPair: s.currencyPairs[1],
			Nonce:        0,
			Id:           1,
		},
		{
			CurrencyPair: s.currencyPairs[2],
			Nonce:        0,
			Id:           2,
		},
	}
	s.genesis = oracletypes.GenesisState{
		CurrencyPairGenesis: genesisCPs,
		NextId:              3,
	}
}

func (s *PreBlockTestSuite) SetupSubTest() {
	s.key = storetypes.NewKVStoreKey(oracletypes.StoreKey)
	s.transientKey = storetypes.NewTransientStoreKey("transient_test")
	s.ctx = testutils.CreateBaseSDKContextWithKeys(s.T(), s.key, s.transientKey)

	// Use the default aggregation function for testing
	mockValidatorStore := mocks.NewValidatorStore(s.T())
	aggregationFn := preblockmath.VoteWeightedMedianFromContext(
		log.NewTestLogger(s.T()),
		mockValidatorStore,
		preblockmath.DefaultPowerThreshold,
	)

	// Use mock metrics
	mockMetrics := metricmock.NewMetrics(s.T())

	// Create the oracle keeper
	s.oracleKeeper = testutils.CreateTestOracleKeeperWithGenesis(s.ctx, s.key, s.genesis)

	s.cpID = strategymock.NewCurrencyPairIDStrategy(s.T())

	s.handler = preblock.NewOraclePreBlockHandler(
		log.NewTestLogger(s.T()),
		aggregationFn,
		s.oracleKeeper,
		s.myVal,
		mockMetrics,
		s.cpID,
	)
}

func (s *PreBlockTestSuite) TestPreBlockHandler() {}

func (s *PreBlockTestSuite) TestWritePrices() {
	s.Run("no prices", func() {
		err := s.handler.WritePrices(s.ctx, nil)
		s.Require().NoError(err)
	})

	s.Run("single price update", func() {
		prices := map[oracletypes.CurrencyPair]*uint256.Int{
			s.currencyPairs[0]: uint256.NewInt(1),
		}

		err := s.handler.WritePrices(s.ctx, prices)
		s.Require().NoError(err)

		// Check that the price was written to state.
		oraclePrice, err := s.oracleKeeper.GetPriceForCurrencyPair(s.ctx, s.currencyPairs[0])
		s.Require().NoError(err)
		s.Require().Equal(math.NewIntFromBigInt(prices[s.currencyPairs[0]].ToBig()), oraclePrice.Price)
	})

	s.Run("multiple price updates", func() {
		prices := map[oracletypes.CurrencyPair]*uint256.Int{
			s.currencyPairs[0]: uint256.NewInt(1),
			s.currencyPairs[1]: uint256.NewInt(2),
			s.currencyPairs[2]: uint256.NewInt(3),
		}

		err := s.handler.WritePrices(s.ctx, prices)
		s.Require().NoError(err)

		// Check that the prices were written to state.
		for _, cp := range s.currencyPairs {
			oraclePrice, err := s.oracleKeeper.GetPriceForCurrencyPair(s.ctx, cp)
			s.Require().NoError(err)
			s.Require().Equal(math.NewIntFromBigInt(prices[cp].ToBig()), oraclePrice.Price)
		}
	})

	s.Run("single price update with a nil price", func() {
		prices := map[oracletypes.CurrencyPair]*uint256.Int{
			s.currencyPairs[0]: nil,
		}

		err := s.handler.WritePrices(s.ctx, prices)
		s.Require().NoError(err)

		// Check that the price was not written to state.
		_, err = s.oracleKeeper.GetPriceForCurrencyPair(s.ctx, s.currencyPairs[0])
		s.Require().Error(err)
	})

	s.Run("attempting to set price for unsupported currency pair", func() {
		unsupportedCP := oracletypes.CurrencyPair{
			Base:  "cap",
			Quote: "on-god",
		}
		prices := map[oracletypes.CurrencyPair]*uint256.Int{
			unsupportedCP: uint256.NewInt(1),
		}

		err := s.handler.WritePrices(s.ctx, prices)
		s.Require().NoError(err)

		// Check that the price was not written to state.
		_, err = s.oracleKeeper.GetPriceForCurrencyPair(s.ctx, unsupportedCP)
		s.Require().Error(err)
	})
}
