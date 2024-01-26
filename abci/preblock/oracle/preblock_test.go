package oracle_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	cmtabci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	preblock "github.com/skip-mev/slinky/abci/preblock/oracle"
	preblockmath "github.com/skip-mev/slinky/abci/preblock/oracle/math"
	"github.com/skip-mev/slinky/abci/preblock/oracle/math/mocks"
	compression "github.com/skip-mev/slinky/abci/strategies/codec"
	currencypairmock "github.com/skip-mev/slinky/abci/strategies/currencypair/mocks"
	"github.com/skip-mev/slinky/abci/testutils"
	servicemetrics "github.com/skip-mev/slinky/service/metrics"
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
	cpID          *currencypairmock.CurrencyPairStrategy
	veCodec       compression.VoteExtensionCodec
	commitCodec   compression.ExtendedCommitCodec
	mockMetrics   *metricmock.Metrics
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

	s.veCodec = compression.NewCompressionVoteExtensionCodec(
		compression.NewDefaultVoteExtensionCodec(),
		compression.NewZLibCompressor(),
	)

	s.commitCodec = compression.NewCompressionExtendedCommitCodec(
		compression.NewDefaultExtendedCommitCodec(),
		compression.NewZLibCompressor(),
	)
}

func (s *PreBlockTestSuite) SetupSubTest() {
	s.key = storetypes.NewKVStoreKey(oracletypes.StoreKey)
	s.transientKey = storetypes.NewTransientStoreKey("transient_test")
	s.ctx = testutils.CreateBaseSDKContextWithKeys(s.T(), s.key, s.transientKey).WithExecMode(sdk.ExecModeFinalize)

	// Use the default aggregation function for testing
	mockValidatorStore := mocks.NewValidatorStore(s.T())
	aggregationFn := preblockmath.VoteWeightedMedianFromContext(
		log.NewTestLogger(s.T()),
		mockValidatorStore,
		preblockmath.DefaultPowerThreshold,
	)

	// Use mock metrics
	s.mockMetrics = metricmock.NewMetrics(s.T())

	// Create the oracle keeper
	s.oracleKeeper = testutils.CreateTestOracleKeeperWithGenesis(s.ctx, s.key, s.genesis)

	s.cpID = currencypairmock.NewCurrencyPairStrategy(s.T())

	s.handler = preblock.NewOraclePreBlockHandler(
		log.NewTestLogger(s.T()),
		aggregationFn,
		s.oracleKeeper,
		s.myVal,
		s.mockMetrics,
		s.cpID,
		s.veCodec,
		s.commitCodec,
	)
}

func (s *PreBlockTestSuite) TestPreBlockHandler() {}

func (s *PreBlockTestSuite) TestWritePrices() {
	s.Run("no prices", func() {
		err := s.handler.WritePrices(s.ctx, nil)
		s.Require().NoError(err)
	})

	s.Run("single price update", func() {
		prices := map[oracletypes.CurrencyPair]*big.Int{
			s.currencyPairs[0]: big.NewInt(1),
		}

		err := s.handler.WritePrices(s.ctx, prices)
		s.Require().NoError(err)

		// Check that the price was written to state.
		oraclePrice, err := s.oracleKeeper.GetPriceForCurrencyPair(s.ctx, s.currencyPairs[0])
		s.Require().NoError(err)
		s.Require().Equal(math.NewIntFromBigInt(prices[s.currencyPairs[0]]), oraclePrice.Price)
	})

	s.Run("multiple price updates", func() {
		prices := map[oracletypes.CurrencyPair]*big.Int{
			s.currencyPairs[0]: big.NewInt(1),
			s.currencyPairs[1]: big.NewInt(2),
			s.currencyPairs[2]: big.NewInt(3),
		}

		err := s.handler.WritePrices(s.ctx, prices)
		s.Require().NoError(err)

		// Check that the prices were written to state.
		for _, cp := range s.currencyPairs {
			oraclePrice, err := s.oracleKeeper.GetPriceForCurrencyPair(s.ctx, cp)
			s.Require().NoError(err)
			s.Require().Equal(math.NewIntFromBigInt(prices[cp]), oraclePrice.Price)
		}
	})

	s.Run("single price update with a nil price", func() {
		prices := map[oracletypes.CurrencyPair]*big.Int{
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
		prices := map[oracletypes.CurrencyPair]*big.Int{
			unsupportedCP: big.NewInt(1),
		}

		err := s.handler.WritePrices(s.ctx, prices)
		s.Require().NoError(err)

		// Check that the price was not written to state.
		_, err = s.oracleKeeper.GetPriceForCurrencyPair(s.ctx, unsupportedCP)
		s.Require().Error(err)
	})
}

func (s *PreBlockTestSuite) TestPreblockLatency() {
	s.Run("expect no metric invocation in non-Finalize Exec mode", func() {
		s.ctx = s.ctx.WithBlockHeight(1)
		s.ctx = s.ctx.WithExecMode(sdk.ExecModePrepareProposal)
		// ves are not enabled
		s.ctx = s.ctx.WithConsensusParams(
			cmtproto.ConsensusParams{
				Abci: &cmtproto.ABCIParams{
					VoteExtensionsEnableHeight: 2,
				},
			},
		)

		// run preblocker
		s.handler.PreBlocker()(s.ctx, &cmtabci.RequestFinalizeBlock{})
	})

	s.Run("expect metric invocation in Finalize Exec mode", func() {
		s.ctx = s.ctx.WithBlockHeight(1)
		s.ctx = s.ctx.WithExecMode(sdk.ExecModeFinalize)
		// ves are not enabled
		s.ctx = s.ctx.WithConsensusParams(
			cmtproto.ConsensusParams{
				Abci: &cmtproto.ABCIParams{
					VoteExtensionsEnableHeight: 2,
				},
			},
		)
		s.mockMetrics.On("ObserveABCIMethodLatency", servicemetrics.PreBlock, mock.Anything).Return()
		// run preblocker
		s.handler.PreBlocker()(s.ctx, &cmtabci.RequestFinalizeBlock{})
	})
}
