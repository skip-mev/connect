package oracle_test

import (
	"fmt"
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
	codecmock "github.com/skip-mev/slinky/abci/strategies/codec/mocks"
	currencypairmock "github.com/skip-mev/slinky/abci/strategies/currencypair/mocks"
	"github.com/skip-mev/slinky/abci/testutils"
	"github.com/skip-mev/slinky/abci/types"
	vetypes "github.com/skip-mev/slinky/abci/ve/types"
	"github.com/skip-mev/slinky/aggregator"
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

		// expect metrics call
		s.mockMetrics.On("ObservePriceForTicker", s.currencyPairs[0], float64(1))

		err := s.handler.WritePrices(s.ctx, prices)
		s.Require().NoError(err)

		// Check that the price was written to state.
		oraclePrice, err := s.oracleKeeper.GetPriceForCurrencyPair(s.ctx, s.currencyPairs[0])
		s.Require().NoError(err)
		s.Require().Equal(math.NewIntFromBigInt(prices[s.currencyPairs[0]]), oraclePrice.Price)
	})

	s.Run("multiple price updates", func() {
		bigIntPrice, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10) // use a non-uint64 price

		prices := map[oracletypes.CurrencyPair]*big.Int{
			s.currencyPairs[0]: big.NewInt(1),
			s.currencyPairs[1]: big.NewInt(2),
			s.currencyPairs[2]: bigIntPrice,
		}
		s.mockMetrics.On("ObservePriceForTicker", s.currencyPairs[0], float64(1))
		s.mockMetrics.On("ObservePriceForTicker", s.currencyPairs[1], float64(2))
		s.mockMetrics.On("ObservePriceForTicker", s.currencyPairs[2], float64(1.157920892373162e+77)) // we can represent 256 bit ints

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
		s.mockMetrics.On("AddABCIRequest", servicemetrics.PreBlock, mock.Anything)
		// run preblocker
		s.handler.PreBlocker()(s.ctx, &cmtabci.RequestFinalizeBlock{})
	})
}

func (s *PreBlockTestSuite) TestPreBlockStatus() {
	s.Run("success", func() {
		metrics := metricmock.NewMetrics(s.T())
		handler := preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			func(ctx sdk.Context) aggregator.AggregateFn[string, map[oracletypes.CurrencyPair]*big.Int] {
				return func(providers aggregator.AggregatedProviderData[string, map[oracletypes.CurrencyPair]*big.Int]) map[oracletypes.CurrencyPair]*big.Int {
					return nil
				}
			},
			nil,
			nil,
			metrics,
			nil,
			nil,
			nil,
		)

		metrics.On("ObserveABCIMethodLatency", servicemetrics.PreBlock, mock.Anything).Return()
		metrics.On("AddABCIRequest", servicemetrics.PreBlock, servicemetrics.Success{}).Return()
		// run preblocker
		_, err := handler.PreBlocker()(s.ctx, &cmtabci.RequestFinalizeBlock{})
		s.Require().NoError(err)
	})

	s.Run("error in getting oracle votes", func() {
		metrics := metricmock.NewMetrics(s.T())
		handler := preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			func(ctx sdk.Context) aggregator.AggregateFn[string, map[oracletypes.CurrencyPair]*big.Int] {
				return func(providers aggregator.AggregatedProviderData[string, map[oracletypes.CurrencyPair]*big.Int]) map[oracletypes.CurrencyPair]*big.Int {
					return nil
				}
			},
			nil,
			nil,
			metrics,
			nil,
			nil,
			nil,
		)

		expErr := types.MissingCommitInfoError{}
		metrics.On("ObserveABCIMethodLatency", servicemetrics.PreBlock, mock.Anything).Return()
		metrics.On("AddABCIRequest", servicemetrics.PreBlock, expErr).Return()

		// make ves enabled
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 2)
		s.ctx = s.ctx.WithBlockHeight(4)
		// run preblocker
		_, err := handler.PreBlocker()(s.ctx, &cmtabci.RequestFinalizeBlock{
			Txs: [][]byte{},
		})
		s.Require().Error(err, expErr)
	})

	s.Run("error in aggregating oracle votes", func() {
		metrics := metricmock.NewMetrics(s.T())
		extCodec := codecmock.NewExtendedCommitCodec(s.T())
		veCodec := codecmock.NewVoteExtensionCodec(s.T())
		handler := preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			func(ctx sdk.Context) aggregator.AggregateFn[string, map[oracletypes.CurrencyPair]*big.Int] {
				return func(providers aggregator.AggregatedProviderData[string, map[oracletypes.CurrencyPair]*big.Int]) map[oracletypes.CurrencyPair]*big.Int {
					return nil
				}
			},
			nil,
			nil,
			metrics,
			nil,
			veCodec,
			extCodec,
		)

		ca := sdk.ConsAddress([]byte("val"))
		extCodec.On("Decode", mock.Anything).Return(cmtabci.ExtendedCommitInfo{
			Votes: []cmtabci.ExtendedVoteInfo{
				{
					VoteExtension: []byte("ve"),
					Validator: cmtabci.Validator{
						Address: ca,
					},
				},
			},
		}, nil)
		veCodec.On("Decode", []byte("ve")).Return(vetypes.OracleVoteExtension{
			Prices: map[uint64][]byte{
				1: make([]byte, 34),
			},
		}, nil)
		expErr := preblock.PriceAggregationError{
			Err: fmt.Errorf("price bytes are too long: %d", 34),
		}

		metrics.On("ObserveABCIMethodLatency", servicemetrics.PreBlock, mock.Anything).Return()
		metrics.On("AddABCIRequest", servicemetrics.PreBlock, expErr).Return()
		// make ves enabled
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 2)
		s.ctx = s.ctx.WithBlockHeight(4)
		// run preblocker
		_, err := handler.PreBlocker()(s.ctx, &cmtabci.RequestFinalizeBlock{
			Txs: [][]byte{
				[]byte("abc"),
			},
		})
		s.Require().Error(err, expErr)
	})
}
