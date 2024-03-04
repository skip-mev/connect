package oracle_test

import (
	"fmt"
	"math/big"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	cometproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	preblock "github.com/skip-mev/slinky/abci/preblock/oracle"
	preblockmock "github.com/skip-mev/slinky/abci/preblock/oracle/mocks"
	compression "github.com/skip-mev/slinky/abci/strategies/codec"
	codecmock "github.com/skip-mev/slinky/abci/strategies/codec/mocks"
	currencypairmock "github.com/skip-mev/slinky/abci/strategies/currencypair/mocks"
	"github.com/skip-mev/slinky/abci/testutils"
	"github.com/skip-mev/slinky/abci/types"
	vetypes "github.com/skip-mev/slinky/abci/ve/types"
	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/pkg/math/voteweighted"
	"github.com/skip-mev/slinky/pkg/math/voteweighted/mocks"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	servicemetrics "github.com/skip-mev/slinky/service/metrics"
	metricmock "github.com/skip-mev/slinky/service/metrics/mocks"
	"github.com/skip-mev/slinky/x/oracle/keeper"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var maxUint256, _ = new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)

type PreBlockTestSuite struct {
	suite.Suite

	myVal         sdk.ConsAddress
	ctx           sdk.Context
	currencyPairs []slinkytypes.CurrencyPair
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

	s.currencyPairs = []slinkytypes.CurrencyPair{
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
	aggregationFn := voteweighted.MedianFromContext(
		log.NewTestLogger(s.T()),
		mockValidatorStore,
		voteweighted.DefaultPowerThreshold,
	)

	// Use mock metrics
	s.mockMetrics = metricmock.NewMetrics(s.T())

	// Create the oracle keeper
	s.oracleKeeper = testutils.CreateTestOracleKeeperWithGenesis(s.T(), s.ctx, s.key, s.genesis)

	s.cpID = currencypairmock.NewCurrencyPairStrategy(s.T())

	s.handler = preblock.NewOraclePreBlockHandler(
		log.NewTestLogger(s.T()),
		aggregationFn,
		&s.oracleKeeper,
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
		prices := map[slinkytypes.CurrencyPair]*big.Int{
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
		prices := map[slinkytypes.CurrencyPair]*big.Int{
			s.currencyPairs[0]: big.NewInt(1),
			s.currencyPairs[1]: big.NewInt(2),
			s.currencyPairs[2]: maxUint256,
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
		prices := map[slinkytypes.CurrencyPair]*big.Int{
			s.currencyPairs[0]: nil,
		}

		err := s.handler.WritePrices(s.ctx, prices)
		s.Require().NoError(err)

		// Check that the price was not written to state.
		_, err = s.oracleKeeper.GetPriceForCurrencyPair(s.ctx, s.currencyPairs[0])
		s.Require().Error(err)
	})

	s.Run("attempting to set price for unsupported currency pair", func() {
		unsupportedCP := slinkytypes.CurrencyPair{
			Base:  "cap",
			Quote: "on-god",
		}
		prices := map[slinkytypes.CurrencyPair]*big.Int{
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
			cometproto.ConsensusParams{
				Abci: &cometproto.ABCIParams{
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
			cometproto.ConsensusParams{
				Abci: &cometproto.ABCIParams{
					VoteExtensionsEnableHeight: 2,
				},
			},
		)
		s.mockMetrics.On("ObserveABCIMethodLatency", servicemetrics.PreBlock, mock.Anything).Return()
		s.mockMetrics.On("AddABCIRequest", servicemetrics.PreBlock, mock.Anything)
		// run preblocker
		s.handler.PreBlocker()(s.ctx, &cmtabci.RequestFinalizeBlock{
			Height: 1,
		})
	})
}

func (s *PreBlockTestSuite) TestPreBlockStatus() {
	s.Run("success", func() {
		metrics := metricmock.NewMetrics(s.T())
		handler := preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context) aggregator.AggregateFn[string, map[slinkytypes.CurrencyPair]*big.Int] {
				return func(_ aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]) map[slinkytypes.CurrencyPair]*big.Int {
					return nil
				}
			},
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
			func(_ sdk.Context) aggregator.AggregateFn[string, map[slinkytypes.CurrencyPair]*big.Int] {
				return func(_ aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]) map[slinkytypes.CurrencyPair]*big.Int {
					return nil
				}
			},
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
			func(_ sdk.Context) aggregator.AggregateFn[string, map[slinkytypes.CurrencyPair]*big.Int] {
				return func(_ aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]) map[slinkytypes.CurrencyPair]*big.Int {
					return nil
				}
			},
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

func (s *PreBlockTestSuite) TestValidatorReports() {
	// test that no reports are recorded when exec-mode is not finalize
	s.Run("test that no reports are recorded when exec-mode is not finalize", func() {
		metrics := metricmock.NewMetrics(s.T())
		handler := preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context) aggregator.AggregateFn[string, map[slinkytypes.CurrencyPair]*big.Int] {
				return func(_ aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]) map[slinkytypes.CurrencyPair]*big.Int {
					return nil
				}
			},
			nil,
			metrics,
			nil,
			nil,
			nil,
		)

		// change exec mode to not be finalize
		s.ctx = s.ctx.WithExecMode(sdk.ExecModeVoteExtension)

		// enable vote-extensions
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 2)
		_, err := handler.PreBlocker()(s.ctx, nil)
		s.Require().NoError(err)
	})

	s.Run("test that no reports are recorded if there is an error", func() {
		metrics := metricmock.NewMetrics(s.T())
		handler := preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context) aggregator.AggregateFn[string, map[slinkytypes.CurrencyPair]*big.Int] {
				return func(_ aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]) map[slinkytypes.CurrencyPair]*big.Int {
					return nil
				}
			},
			nil,
			metrics,
			nil,
			nil,
			nil,
		)

		// expect metrics calls
		metrics.On("ObserveABCIMethodLatency", servicemetrics.PreBlock, mock.Anything).Return()
		metrics.On("AddABCIRequest", servicemetrics.PreBlock, types.MissingCommitInfoError{}).Return()

		// don't expect preblocker to record per validator metrics

		// enable ves + set exec mode
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 2)
		s.ctx = s.ctx.WithBlockHeight(4)
		s.ctx = s.ctx.WithExecMode(sdk.ExecModeFinalize)
		// run preblocker
		_, err := handler.PreBlocker()(s.ctx, &cmtabci.RequestFinalizeBlock{})
		s.Require().Error(err, types.MissingCommitInfoError{})
	})

	// test that if a validator's report is absent, this is recorded
	s.Run("test that if a validator's report is absent, this is recorded", func() {
		metrics := metricmock.NewMetrics(s.T())
		val1 := sdk.ConsAddress("val1")
		val2 := sdk.ConsAddress("val2")
		val3 := sdk.ConsAddress("val3")

		mockOracleKeeper := preblockmock.NewKeeper(s.T())
		currencyPairStrategyMock := currencypairmock.NewCurrencyPairStrategy(s.T())

		btcUsd := slinkytypes.NewCurrencyPair("BTC", "USD")
		mogUsd := slinkytypes.NewCurrencyPair("MOG", "USD")

		handler := preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context) aggregator.AggregateFn[string, map[slinkytypes.CurrencyPair]*big.Int] {
				return func(_ aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]) map[slinkytypes.CurrencyPair]*big.Int {
					return map[slinkytypes.CurrencyPair]*big.Int{
						// return default values
						btcUsd: big.NewInt(1),
						mogUsd: maxUint256,
					}
				}
			},
			mockOracleKeeper,
			metrics,
			currencyPairStrategyMock,
			compression.NewDefaultVoteExtensionCodec(),
			compression.NewDefaultExtendedCommitCodec(),
		)

		// enable ves + set exec mode
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 2)
		s.ctx = s.ctx.WithBlockHeight(4)
		s.ctx = s.ctx.WithExecMode(sdk.ExecModeFinalize)

		// mock currency-pair strategy calls
		currencyPairStrategyMock.On("FromID", s.ctx, uint64(0)).Return(btcUsd, nil)
		currencyPairStrategyMock.On("FromID", s.ctx, uint64(1)).Return(mogUsd, nil)
		currencyPairStrategyMock.On("GetDecodedPrice", s.ctx, btcUsd, big.NewInt(1).Bytes()).Return(big.NewInt(1), nil)
		currencyPairStrategyMock.On("GetDecodedPrice", s.ctx, btcUsd, big.NewInt(2).Bytes()).Return(big.NewInt(2), nil)
		currencyPairStrategyMock.On("GetDecodedPrice", s.ctx, mogUsd, mock.Anything).Return(maxUint256, nil)

		// mock oracle keeper calls
		mockOracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{btcUsd, mogUsd}, nil)
		mockOracleKeeper.On("SetPriceForCurrencyPair", s.ctx, btcUsd, mock.Anything).Return(nil)
		mockOracleKeeper.On("SetPriceForCurrencyPair", s.ctx, mogUsd, mock.Anything).Return(nil)

		// create extended commit info
		val1Vote, err := testutils.CreateExtendedVoteInfo(val1, map[uint64][]byte{
			0: big.NewInt(1).Bytes(),
			1: maxUint256.Bytes(),
		}, compression.NewDefaultVoteExtensionCodec())
		s.Require().NoError(err)

		val2Vote, err := testutils.CreateExtendedVoteInfo(val2, map[uint64][]byte{
			0: big.NewInt(2).Bytes(),
		}, compression.NewDefaultVoteExtensionCodec())
		s.Require().NoError(err)

		_, extCommitBz, err := testutils.CreateExtendedCommitInfo([]cmtabci.ExtendedVoteInfo{val1Vote, val2Vote}, compression.NewDefaultExtendedCommitCodec())
		s.Require().NoError(err)

		// expect metrics calls
		metrics.On("ObserveABCIMethodLatency", servicemetrics.PreBlock, mock.Anything).Return()
		metrics.On("AddABCIRequest", servicemetrics.PreBlock, servicemetrics.Success{}).Return()
		// expect price reports
		float, _ := big.NewInt(1).Float64()
		metrics.On("ObservePriceForTicker", btcUsd, float)
		float, _ = maxUint256.Float64()
		metrics.On("ObservePriceForTicker", mogUsd, float)

		// expect per validator metrics
		// val1
		metrics.On("AddValidatorReportForTicker", val1.String(), btcUsd, servicemetrics.WithPrice)
		metrics.On("AddValidatorReportForTicker", val1.String(), mogUsd, servicemetrics.WithPrice)
		float, _ = big.NewInt(1).Float64()
		metrics.On("AddValidatorPriceForTicker", val1.String(), btcUsd, float)
		float, _ = maxUint256.Float64()
		metrics.On("AddValidatorPriceForTicker", val1.String(), mogUsd, float)
		// val2
		metrics.On("AddValidatorReportForTicker", val2.String(), btcUsd, servicemetrics.WithPrice)
		metrics.On("AddValidatorReportForTicker", val2.String(), mogUsd, servicemetrics.MissingPrice)
		float, _ = big.NewInt(2).Float64()
		metrics.On("AddValidatorPriceForTicker", val2.String(), btcUsd, float).Once()
		// val3
		metrics.On("AddValidatorReportForTicker", val3.String(), btcUsd, servicemetrics.Absent)
		metrics.On("AddValidatorReportForTicker", val3.String(), mogUsd, servicemetrics.Absent)

		// run preblocker
		_, err = handler.PreBlocker()(s.ctx, &cmtabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitBz},
			DecidedLastCommit: cmtabci.CommitInfo{
				Votes: []cmtabci.VoteInfo{
					{
						Validator: cmtabci.Validator{
							Address: val1,
						},
						BlockIdFlag: cometproto.BlockIDFlagCommit,
					},
					{
						Validator: cmtabci.Validator{
							Address: val2,
						},
						BlockIdFlag: cometproto.BlockIDFlagCommit,
					},
					{
						Validator: cmtabci.Validator{
							Address: val3,
						},
						BlockIdFlag: cometproto.BlockIDFlagAbsent,
					},
				},
			},
		})
		s.Require().NoError(err)
	})
}
