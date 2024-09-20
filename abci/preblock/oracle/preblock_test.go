package oracle_test

import (
	"context"
	"math/big"
	"testing"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cometproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	preblock "github.com/skip-mev/connect/v2/abci/preblock/oracle"
	compression "github.com/skip-mev/connect/v2/abci/strategies/codec"
	codecmock "github.com/skip-mev/connect/v2/abci/strategies/codec/mocks"
	"github.com/skip-mev/connect/v2/abci/strategies/currencypair"
	currencypairmock "github.com/skip-mev/connect/v2/abci/strategies/currencypair/mocks"
	"github.com/skip-mev/connect/v2/abci/testutils"
	"github.com/skip-mev/connect/v2/abci/types"
	connectabcimocks "github.com/skip-mev/connect/v2/abci/types/mocks"
	vetypes "github.com/skip-mev/connect/v2/abci/ve/types"
	"github.com/skip-mev/connect/v2/aggregator"
	"github.com/skip-mev/connect/v2/pkg/math/voteweighted"
	voteweightedmocks "github.com/skip-mev/connect/v2/pkg/math/voteweighted/mocks"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	servicemetrics "github.com/skip-mev/connect/v2/service/metrics"
	metricmock "github.com/skip-mev/connect/v2/service/metrics/mocks"
	"github.com/skip-mev/connect/v2/x/oracle/keeper"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
)

var maxUint256, _ = new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)

type PreBlockTestSuite struct {
	suite.Suite

	myVal         sdk.ConsAddress
	ctx           sdk.Context
	currencyPairs []connecttypes.CurrencyPair
	genesis       oracletypes.GenesisState
	key           *storetypes.KVStoreKey
	transientKey  *storetypes.TransientStoreKey
	oracleKeeper  keeper.Keeper
	handler       *preblock.PreBlockHandler
	cpID          *currencypairmock.CurrencyPairStrategy
	veCodec       compression.VoteExtensionCodec
	commitCodec   compression.ExtendedCommitCodec
	mockMetrics   *metricmock.Metrics
	mm            *module.Manager
}

func TestPreBlockTestSuite(t *testing.T) {
	suite.Run(t, new(PreBlockTestSuite))
}

func (s *PreBlockTestSuite) SetupTest() {
	s.mm = &module.Manager{Modules: map[string]interface{}{}, OrderPreBlockers: make([]string, 0)}
	s.myVal = sdk.ConsAddress("myVal")
	s.currencyPairs = []connecttypes.CurrencyPair{
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
	mockValidatorStore := voteweightedmocks.NewValidatorStore(s.T())
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

type fakeModule struct{ called int }

func (f *fakeModule) IsOnePerModuleType() {}

func (f *fakeModule) IsAppModule() {}

func (f *fakeModule) Name() string { return "fake" }

type response struct{}

func (r response) IsConsensusParamsChanged() bool { return true }

func (f *fakeModule) PreBlock(_ context.Context) (appmodule.ResponsePreBlock, error) {
	f.called++
	return &response{}, nil
}

func (s *PreBlockTestSuite) TestPreBlocker() {
	mockValidatorStore := voteweightedmocks.NewValidatorStore(s.T())
	aggregationFn := voteweighted.MedianFromContext(
		log.NewTestLogger(s.T()),
		mockValidatorStore,
		voteweighted.DefaultPowerThreshold,
	)

	s.Run("fail on nil requests", func() {
		s.handler = preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			aggregationFn,
			&s.oracleKeeper,
			servicemetrics.NewNopMetrics(),
			s.cpID,
			s.veCodec,
			s.commitCodec,
		)

		_, err := s.handler.WrappedPreBlocker(s.mm)(s.ctx, nil)
		s.Require().Error(err)

		// require no updates
		cps := s.oracleKeeper.GetAllCurrencyPairs(s.ctx)
		for _, cp := range cps {
			nonce, err := s.oracleKeeper.GetNonceForCurrencyPair(s.ctx, cp)
			s.Require().NoError(err)

			require.Equal(s.T(), uint64(0), nonce)
		}
	})

	s.Run("return when ves aren't enabled", func() {
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 2).WithBlockHeight(1)

		s.handler = preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			aggregationFn,
			&s.oracleKeeper,
			servicemetrics.NewNopMetrics(),
			s.cpID,
			s.veCodec,
			s.commitCodec,
		)

		_, err := s.handler.WrappedPreBlocker(s.mm)(s.ctx, &cometabci.RequestFinalizeBlock{})
		s.Require().NoError(err)

		// require no updates
		cps := s.oracleKeeper.GetAllCurrencyPairs(s.ctx)
		for _, cp := range cps {
			nonce, err := s.oracleKeeper.GetNonceForCurrencyPair(s.ctx, cp)
			s.Require().NoError(err)

			require.Equal(s.T(), uint64(0), nonce)
		}
	})

	s.Run("manager PreBlock is called", func() {
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 2).WithBlockHeight(1)

		s.handler = preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			aggregationFn,
			&s.oracleKeeper,
			servicemetrics.NewNopMetrics(),
			s.cpID,
			s.veCodec,
			s.commitCodec,
		)
		// setup fake module.
		fake := &fakeModule{}
		s.mm.OrderPreBlockers = []string{fake.Name()}
		s.mm.Modules[fake.Name()] = fake
		res, err := s.handler.WrappedPreBlocker(s.mm)(s.ctx, &cometabci.RequestFinalizeBlock{})
		s.Require().NoError(err)
		s.Require().Equal(fake.called, 1)
		s.Require().True(res.IsConsensusParamsChanged()) // should return the response set above the fake module def.
	})

	// update ctx to enable ves
	s.ctx = s.ctx.WithBlockHeight(3)

	s.Run("ignore vote-extensions w/ prices for non-existent pairs", func() {
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 2).WithBlockHeight(3)
		s.handler = preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			aggregationFn,
			&s.oracleKeeper,
			servicemetrics.NewNopMetrics(),
			currencypair.NewDefaultCurrencyPairStrategy(
				&s.oracleKeeper,
			),
			s.veCodec,
			s.commitCodec,
		)

		priceBz, err := big.NewInt(1).GobEncode()
		s.Require().NoError(err)

		prices1 := map[uint64][]byte{
			3: priceBz,
		}
		ca1 := sdk.ConsAddress([]byte("ca1"))

		ve1, err := testutils.CreateExtendedVoteInfo(
			ca1,
			prices1,
			s.veCodec,
		)
		s.Require().NoError(err)

		prices2 := map[uint64][]byte{
			3: priceBz,
		}
		ca2 := sdk.ConsAddress([]byte("ca2"))

		ve2, err := testutils.CreateExtendedVoteInfo(
			ca2,
			prices2,
			s.veCodec,
		)
		s.Require().NoError(err)

		_, extCommitBz, err := testutils.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{
				ve1, ve2,
			},
			s.commitCodec,
		)
		s.Require().NoError(err)

		validator1 := voteweightedmocks.NewValidatorI(s.T())
		validator1.EXPECT().GetBondedTokens().Return(math.NewInt(1))

		validator2 := voteweightedmocks.NewValidatorI(s.T())
		validator2.On("GetBondedTokens").Return(math.NewInt(1))

		mockValidatorStore.On("ValidatorByConsAddr", s.ctx, ca1).Return(
			validator1, nil,
		)

		mockValidatorStore.On("ValidatorByConsAddr", s.ctx, ca2).Return(
			validator2, nil,
		)

		mockValidatorStore.On("TotalBondedTokens", s.ctx).Return(math.NewInt(2), nil)

		_, err = s.handler.WrappedPreBlocker(s.mm)(s.ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitBz},
		})
		s.Require().NoError(err)

		// require no updates
		cps := s.oracleKeeper.GetAllCurrencyPairs(s.ctx)
		for _, cp := range cps {
			nonce, err := s.oracleKeeper.GetNonceForCurrencyPair(s.ctx, cp)
			s.Require().NoError(err)

			require.Equal(s.T(), uint64(0), nonce)
		}
	})

	s.Run("multiple assets to write a price for", func() {
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 2).WithBlockHeight(3)
		s.handler = preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			aggregationFn,
			&s.oracleKeeper,
			servicemetrics.NewNopMetrics(),
			currencypair.NewDefaultCurrencyPairStrategy(
				&s.oracleKeeper,
			),
			s.veCodec,
			s.commitCodec,
		)

		price1Bz, err := big.NewInt(1).GobEncode()
		s.Require().NoError(err)
		price2Bz, err := big.NewInt(2).GobEncode()
		s.Require().NoError(err)
		price3Bz, err := big.NewInt(3).GobEncode()
		s.Require().NoError(err)

		prices := map[uint64][]byte{
			0: price1Bz,
			1: price2Bz,
			2: price3Bz,
		}
		ca1 := sdk.ConsAddress([]byte("ca1"))

		ve1, err := testutils.CreateExtendedVoteInfo(
			ca1,
			prices,
			s.veCodec,
		)
		s.Require().NoError(err)

		ca2 := sdk.ConsAddress([]byte("ca2"))

		ve2, err := testutils.CreateExtendedVoteInfo(
			ca2,
			prices,
			s.veCodec,
		)
		s.Require().NoError(err)

		_, extCommitBz, err := testutils.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{
				ve1, ve2,
			},
			s.commitCodec,
		)
		s.Require().NoError(err)

		validator1 := voteweightedmocks.NewValidatorI(s.T())
		validator1.On("GetBondedTokens").Return(math.NewInt(1))

		validator2 := voteweightedmocks.NewValidatorI(s.T())
		validator2.On("GetBondedTokens").Return(math.NewInt(1))

		mockValidatorStore.On("ValidatorByConsAddr", s.ctx, ca1).Return(
			validator1, nil,
		)

		mockValidatorStore.On("ValidatorByConsAddr", s.ctx, ca2).Return(
			validator2, nil,
		)

		mockValidatorStore.On("TotalBondedTokens", s.ctx).Return(math.NewInt(2), nil)

		_, err = s.handler.WrappedPreBlocker(s.mm)(s.ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitBz},
		})
		s.Require().NoError(err)

		cps := s.oracleKeeper.GetAllCurrencyPairs(s.ctx)
		for _, cp := range cps {
			nonce, err := s.oracleKeeper.GetNonceForCurrencyPair(s.ctx, cp)
			s.Require().NoError(err)

			require.Equal(s.T(), uint64(1), nonce)
		}
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
		_, err := s.handler.WrappedPreBlocker(s.mm)(s.ctx, &cometabci.RequestFinalizeBlock{})
		s.Require().NoError(err)
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
		_, err := s.handler.WrappedPreBlocker(s.mm)(s.ctx, &cometabci.RequestFinalizeBlock{
			Height: 1,
		})
		s.Require().NoError(err)
	})
}

func (s *PreBlockTestSuite) TestPreBlockStatus() {
	s.Run("failure - nil request", func() {
		metrics := metricmock.NewMetrics(s.T())
		handler := preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context) aggregator.AggregateFn[string, map[connecttypes.CurrencyPair]*big.Int] {
				return func(_ aggregator.AggregatedProviderData[string, map[connecttypes.CurrencyPair]*big.Int]) map[connecttypes.CurrencyPair]*big.Int {
					return nil
				}
			},
			nil,
			metrics,
			nil,
			nil,
			nil,
		)

		// run preblocker
		_, err := handler.WrappedPreBlocker(s.mm)(s.ctx, nil)
		s.Require().Error(err)
	})

	s.Run("success", func() {
		metrics := metricmock.NewMetrics(s.T())
		handler := preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context) aggregator.AggregateFn[string, map[connecttypes.CurrencyPair]*big.Int] {
				return func(_ aggregator.AggregatedProviderData[string, map[connecttypes.CurrencyPair]*big.Int]) map[connecttypes.CurrencyPair]*big.Int {
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
		_, err := handler.WrappedPreBlocker(s.mm)(s.ctx, &cometabci.RequestFinalizeBlock{})
		s.Require().NoError(err)
	})

	s.Run("error in getting oracle votes", func() {
		metrics := metricmock.NewMetrics(s.T())
		handler := preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context) aggregator.AggregateFn[string, map[connecttypes.CurrencyPair]*big.Int] {
				return func(_ aggregator.AggregatedProviderData[string, map[connecttypes.CurrencyPair]*big.Int]) map[connecttypes.CurrencyPair]*big.Int {
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
		_, err := handler.WrappedPreBlocker(s.mm)(s.ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{},
		})
		s.Require().Error(err, expErr)
	})

	s.Run("too many price bytes causes no errors in aggregating currency pairs", func() {
		metrics := metricmock.NewMetrics(s.T())
		extCodec := codecmock.NewExtendedCommitCodec(s.T())
		veCodec := codecmock.NewVoteExtensionCodec(s.T())
		mockOracleKeeper := connectabcimocks.NewOracleKeeper(s.T())
		handler := preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context) aggregator.AggregateFn[string, map[connecttypes.CurrencyPair]*big.Int] {
				return func(_ aggregator.AggregatedProviderData[string, map[connecttypes.CurrencyPair]*big.Int]) map[connecttypes.CurrencyPair]*big.Int {
					return nil
				}
			},
			mockOracleKeeper,
			metrics,
			nil,
			veCodec,
			extCodec,
		)

		ca := sdk.ConsAddress("val")
		extCodec.On("Decode", mock.Anything).Return(cometabci.ExtendedCommitInfo{
			Votes: []cometabci.ExtendedVoteInfo{
				{
					VoteExtension: []byte("ve"),
					Validator: cometabci.Validator{
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

		metrics.On("ObserveABCIMethodLatency", servicemetrics.PreBlock, mock.Anything).Return()
		metrics.On("AddABCIRequest", servicemetrics.PreBlock, servicemetrics.Success{}).Return()
		// make ves enabled
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 2)
		s.ctx = s.ctx.WithBlockHeight(4)
		mockOracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return(nil)

		// run preblocker
		_, err := handler.WrappedPreBlocker(s.mm)(s.ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{
				[]byte("abc"),
			},
		})
		s.Require().NoError(err)
	})
}

func (s *PreBlockTestSuite) TestValidatorReports() {
	// test that no reports are recorded when exec-mode is not finalize
	s.Run("test that no reports are recorded when exec-mode is not finalize", func() {
		metrics := metricmock.NewMetrics(s.T())
		handler := preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context) aggregator.AggregateFn[string, map[connecttypes.CurrencyPair]*big.Int] {
				return func(_ aggregator.AggregatedProviderData[string, map[connecttypes.CurrencyPair]*big.Int]) map[connecttypes.CurrencyPair]*big.Int {
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
		req := &cometabci.RequestFinalizeBlock{}
		_, err := handler.WrappedPreBlocker(s.mm)(s.ctx, req)
		s.Require().NoError(err)
	})

	s.Run("test that no reports are recorded if there is an error", func() {
		metrics := metricmock.NewMetrics(s.T())
		handler := preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context) aggregator.AggregateFn[string, map[connecttypes.CurrencyPair]*big.Int] {
				return func(_ aggregator.AggregatedProviderData[string, map[connecttypes.CurrencyPair]*big.Int]) map[connecttypes.CurrencyPair]*big.Int {
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
		_, err := handler.WrappedPreBlocker(s.mm)(s.ctx, &cometabci.RequestFinalizeBlock{})
		s.Require().Error(err, types.MissingCommitInfoError{})
	})

	// test that if a validator's report is absent, this is recorded
	s.Run("test that if a validator's report is absent, this is recorded", func() {
		metrics := metricmock.NewMetrics(s.T())
		val1 := sdk.ConsAddress("val1")
		val2 := sdk.ConsAddress("val2")
		val3 := sdk.ConsAddress("val3")

		mockOracleKeeper := connectabcimocks.NewOracleKeeper(s.T())
		currencyPairStrategyMock := currencypairmock.NewCurrencyPairStrategy(s.T())

		btcUsd := connecttypes.NewCurrencyPair("BTC", "USD")
		mogUsd := connecttypes.NewCurrencyPair("MOG", "USD")

		handler := preblock.NewOraclePreBlockHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context) aggregator.AggregateFn[string, map[connecttypes.CurrencyPair]*big.Int] {
				return func(_ aggregator.AggregatedProviderData[string, map[connecttypes.CurrencyPair]*big.Int]) map[connecttypes.CurrencyPair]*big.Int {
					return map[connecttypes.CurrencyPair]*big.Int{
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
		mockOracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]connecttypes.CurrencyPair{btcUsd, mogUsd}, nil)
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

		_, extCommitBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{val1Vote, val2Vote}, compression.NewDefaultExtendedCommitCodec())
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
		_, err = handler.WrappedPreBlocker(s.mm)(s.ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitBz},
			DecidedLastCommit: cometabci.CommitInfo{
				Votes: []cometabci.VoteInfo{
					{
						Validator: cometabci.Validator{
							Address: val1,
						},
						BlockIdFlag: cometproto.BlockIDFlagCommit,
					},
					{
						Validator: cometabci.Validator{
							Address: val2,
						},
						BlockIdFlag: cometproto.BlockIDFlagCommit,
					},
					{
						Validator: cometabci.Validator{
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
