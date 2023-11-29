package oracle_test

import (
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/mock"

	preblock "github.com/skip-mev/slinky/abci/preblock/oracle"
	preblockmath "github.com/skip-mev/slinky/abci/preblock/oracle/math"
	"github.com/skip-mev/slinky/abci/preblock/oracle/math/mocks"
	preblockmock "github.com/skip-mev/slinky/abci/preblock/oracle/mocks"
	strategymocks "github.com/skip-mev/slinky/abci/strategies/mocks"
	"github.com/skip-mev/slinky/abci/testutils"
	metricmock "github.com/skip-mev/slinky/service/metrics/mocks"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	btcUSD = oracletypes.CurrencyPair{
		Base:  "BTC",
		Quote: "USD",
	}
	ethUSD = oracletypes.CurrencyPair{
		Base:  "ETH",
		Quote: "USD",
	}
	ethBTC = oracletypes.CurrencyPair{
		Base:  "ETH",
		Quote: "BTC",
	}

	oneHundred   = uint256.NewInt(100)
	twoHundred   = uint256.NewInt(200)
	threeHundred = uint256.NewInt(300)
	fourHundred  = uint256.NewInt(400)
	sixHundred   = uint256.NewInt(600)
	sevenHundred = uint256.NewInt(700)
	eightHundred = uint256.NewInt(800)
	nineHundred  = uint256.NewInt(900)

	val1 = sdk.ConsAddress([]byte("val1"))
	val2 = sdk.ConsAddress([]byte("val2"))

	ongodhecappin = append([]byte("ongodhecappin"), make([]byte, 32)...)
)

func (s *PreBlockTestSuite) TestAggregateOracleData() {
	// Use the default aggregation function for testing
	mockValidatorStore := mocks.NewValidatorStore(s.T())
	aggregationFn := preblockmath.VoteWeightedMedianFromContext(
		log.NewTestLogger(s.T()),
		mockValidatorStore,
		preblockmath.DefaultPowerThreshold,
	)
	mockValidatorStore.On("TotalBondedTokens", mock.Anything).Return(math.NewInt(100), nil)

	mockOracleKeeper := preblockmock.NewKeeper(s.T())
	mockMetrics := metricmock.NewMetrics(s.T())

	cpID := strategymocks.NewCurrencyPairIDStrategy(s.T())

	handler := preblock.NewOraclePreBlockHandler(
		log.NewTestLogger(s.T()),
		aggregationFn,
		mockOracleKeeper,
		s.myVal,
		mockMetrics,
		cpID,
		s.veCodec,
		s.commitCodec,
	)

	s.Run("no oracle data", func() {
		_, commitBz, err := testutils.CreateExtendedCommitInfo(nil, s.commitCodec)
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

		mockMetrics.On("AddVoteIncludedInLastCommit", false).Once()

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().NoError(err)
		s.Require().Len(prices, 0)
	})

	s.Run("single oracle data", func() {
		// Create a single vote extension from my validator
		myValPrices := map[uint64][]byte{
			0: oneHundred.Bytes(),
		}
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(s.myVal, myValPrices, s.veCodec)
		s.Require().NoError(err)

		// Create the extended commit info
		_, commitBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo}, s.commitCodec)
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

		// The validator is included in the commit and the price should be included
		mockMetrics.On("AddVoteIncludedInLastCommit", true).Once()
		mockMetrics.On("AddTickerInclusionStatus", btcUSD.ToString(), true).Once()

		cpID.On("FromID", s.ctx, uint64(0)).Return(btcUSD, nil).Once()

		// Assume the validator takes up all of the voting power
		mockValidatorStore.On("ValidatorByConsAddr", mock.Anything, s.myVal).Return(
			stakingtypes.Validator{
				Tokens: math.NewInt(100),
				Status: stakingtypes.Bonded,
			},
			nil,
		).Once()

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().NoError(err)
		s.Require().Len(prices, 1)

		// Check that the prices are correct
		s.Require().Equal(oneHundred.String(), prices[btcUSD].String())
	})

	s.Run("multiple prices from a single validator", func() {
		// Create a single vote extension from my validator
		myValPrices := map[uint64][]byte{
			0: oneHundred.Bytes(),
			1: twoHundred.Bytes(),
		}
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(s.myVal, myValPrices, s.veCodec)
		s.Require().NoError(err)

		// Create the extended commit info
		_, commitBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo}, s.commitCodec)
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

		// The validator is included in the commit and the price should be included
		mockMetrics.On("AddVoteIncludedInLastCommit", true).Once()
		mockMetrics.On("AddTickerInclusionStatus", btcUSD.ToString(), true).Once()
		mockMetrics.On("AddTickerInclusionStatus", ethUSD.ToString(), true).Once()

		// Assume the validator takes up all of the voting power
		mockValidatorStore.On("ValidatorByConsAddr", mock.Anything, s.myVal).Return(
			stakingtypes.Validator{
				Tokens: math.NewInt(100),
				Status: stakingtypes.Bonded,
			},
			nil,
		).Once()

		cpID.On("FromID", s.ctx, uint64(0)).Return(btcUSD, nil).Once()
		cpID.On("FromID", s.ctx, uint64(1)).Return(ethUSD, nil).Once()

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().NoError(err)
		s.Require().Len(prices, 2)

		// Check that the prices are correct
		s.Require().Equal(oneHundred.String(), prices[btcUSD].String())
		s.Require().Equal(twoHundred.String(), prices[ethUSD].String())
	})

	s.Run("single price from two different validators", func() {
		// Create a single vote extension from my validator
		myValPrices := map[uint64][]byte{
			0: oneHundred.Bytes(),
		}
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(s.myVal, myValPrices, s.veCodec)
		s.Require().NoError(err)

		// Create a single vote extension from another validator
		otherValPrices := map[uint64][]byte{
			0: twoHundred.Bytes(),
		}
		otherValVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, otherValPrices, s.veCodec)
		s.Require().NoError(err)

		// Create the extended commit info
		_, commitBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo, otherValVoteInfo}, s.commitCodec)
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

		// The validator is included in the commit and the price should be included
		mockMetrics.On("AddVoteIncludedInLastCommit", true).Once()
		mockMetrics.On("AddTickerInclusionStatus", btcUSD.ToString(), true).Once()

		// Assume the validators have an equal stake
		mockValidatorStore.On("ValidatorByConsAddr", mock.Anything, s.myVal).Return(
			stakingtypes.Validator{
				Tokens: math.NewInt(50),
				Status: stakingtypes.Bonded,
			},
			nil,
		).Once()
		mockValidatorStore.On("ValidatorByConsAddr", mock.Anything, val1).Return(
			stakingtypes.Validator{
				Tokens: math.NewInt(50),
				Status: stakingtypes.Bonded,
			},
			nil,
		).Once()

		cpID.On("FromID", s.ctx, uint64(0)).Return(btcUSD, nil).Twice()

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().NoError(err)
		s.Require().Len(prices, 1)

		// Check that the prices are correct
		s.Require().Equal(oneHundred.String(), prices[btcUSD].String())
	})

	s.Run("single price update from multiple validators but not enough voting power", func() {
		// Create a single vote extension from my validator
		myValPrices := map[uint64][]byte{
			0: oneHundred.Bytes(),
		}
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(s.myVal, myValPrices, s.veCodec)
		s.Require().NoError(err)

		// Create a single vote extension from another validator
		otherValPrices := map[uint64][]byte{
			0: twoHundred.Bytes(),
		}
		otherValVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, otherValPrices, s.veCodec)
		s.Require().NoError(err)

		// Create the extended commit info
		_, commitBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo, otherValVoteInfo}, s.commitCodec)
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

		// The validator is included in the commit and the price should be included
		mockMetrics.On("AddVoteIncludedInLastCommit", true).Once()
		mockMetrics.On("AddTickerInclusionStatus", btcUSD.ToString(), true).Once()

		// Assume the validators have an equal stake
		mockValidatorStore.On("ValidatorByConsAddr", mock.Anything, s.myVal).Return(
			stakingtypes.Validator{
				Tokens: math.NewInt(10),
				Status: stakingtypes.Bonded,
			},
			nil,
		).Once()
		mockValidatorStore.On("ValidatorByConsAddr", mock.Anything, val1).Return(
			stakingtypes.Validator{
				Tokens: math.NewInt(10),
				Status: stakingtypes.Bonded,
			},
			nil,
		).Once()
		cpID.On("FromID", s.ctx, uint64(0)).Return(btcUSD, nil).Twice()

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().NoError(err)
		s.Require().Len(prices, 0)
	})

	s.Run("multiple prices from multiple validators", func() {
		// Create a vote extension with multiple prices from my validator
		myValPrices := map[uint64][]byte{
			0: oneHundred.Bytes(),
			1: twoHundred.Bytes(),
		}
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(s.myVal, myValPrices, s.veCodec)
		s.Require().NoError(err)

		// Create a vote extension with multiple prices from another validator
		otherValPrices := map[uint64][]byte{
			0: threeHundred.Bytes(),
			1: fourHundred.Bytes(),
		}
		otherValVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, otherValPrices, s.veCodec)
		s.Require().NoError(err)

		// Create the extended commit info
		_, commitBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo, otherValVoteInfo}, s.commitCodec)
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

		// The validator is included in the commit and the price should be included
		mockMetrics.On("AddVoteIncludedInLastCommit", true).Once()
		mockMetrics.On("AddTickerInclusionStatus", btcUSD.ToString(), true).Once()
		mockMetrics.On("AddTickerInclusionStatus", ethUSD.ToString(), true).Once()

		// Assume the validators have an unequal stake
		mockValidatorStore.On("ValidatorByConsAddr", mock.Anything, s.myVal).Return(
			stakingtypes.Validator{
				Tokens: math.NewInt(10),
				Status: stakingtypes.Bonded,
			},
			nil,
		).Once()
		mockValidatorStore.On("ValidatorByConsAddr", mock.Anything, val1).Return(
			stakingtypes.Validator{
				Tokens: math.NewInt(90),
				Status: stakingtypes.Bonded,
			},
			nil,
		).Once()

		cpID.On("FromID", s.ctx, uint64(0)).Return(btcUSD, nil).Twice()
		cpID.On("FromID", s.ctx, uint64(1)).Return(ethUSD, nil).Twice()

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().NoError(err)
		s.Require().Len(prices, 2)

		// Check that the prices are correct
		s.Require().Equal(threeHundred.String(), prices[btcUSD].String())
		s.Require().Equal(fourHundred.String(), prices[ethUSD].String())
	})

	s.Run("multiple prices from multiple validators but not enough voting power for some", func() {
		myValPrices := map[uint64][]byte{
			0: oneHundred.Bytes(),
			1: twoHundred.Bytes(),
			2: threeHundred.Bytes(),
		}
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(s.myVal, myValPrices, s.veCodec)
		s.Require().NoError(err)

		val1Prices := map[uint64][]byte{
			0: fourHundred.Bytes(),
			2: sixHundred.Bytes(),
		}
		val1VoteInfo, err := testutils.CreateExtendedVoteInfo(val1, val1Prices, s.veCodec)
		s.Require().NoError(err)

		val2Prices := map[uint64][]byte{
			0: sevenHundred.Bytes(),
			1: eightHundred.Bytes(),
			2: nineHundred.Bytes(),
		}
		val2VoteInfo, err := testutils.CreateExtendedVoteInfo(val2, val2Prices, s.veCodec)
		s.Require().NoError(err)

		// Create the extended commit info
		_, commitBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo, val1VoteInfo, val2VoteInfo}, s.commitCodec)
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

		// The validator is included in the commit and the price should be included
		mockMetrics.On("AddVoteIncludedInLastCommit", true).Once()
		mockMetrics.On("AddTickerInclusionStatus", btcUSD.ToString(), true).Once()
		mockMetrics.On("AddTickerInclusionStatus", ethBTC.ToString(), true).Once()
		mockMetrics.On("AddTickerInclusionStatus", ethUSD.ToString(), true).Once()

		// Assume the validators have an unequal stake
		mockValidatorStore.On("ValidatorByConsAddr", mock.Anything, s.myVal).Return(
			stakingtypes.Validator{
				Tokens: math.NewInt(25),
				Status: stakingtypes.Bonded,
			},
			nil,
		).Once()
		mockValidatorStore.On("ValidatorByConsAddr", mock.Anything, val1).Return(
			stakingtypes.Validator{
				Tokens: math.NewInt(50),
				Status: stakingtypes.Bonded,
			},
			nil,
		).Once()
		mockValidatorStore.On("ValidatorByConsAddr", mock.Anything, val2).Return(
			stakingtypes.Validator{
				Tokens: math.NewInt(25),
				Status: stakingtypes.Bonded,
			},
			nil,
		).Once()

		cpID.On("FromID", s.ctx, uint64(0)).Return(btcUSD, nil).Twice()
		cpID.On("FromID", s.ctx, uint64(0)).Return(ethBTC, nil).Once()
		cpID.On("FromID", s.ctx, uint64(1)).Return(ethBTC, nil).Twice()
		cpID.On("FromID", s.ctx, uint64(2)).Return(ethUSD, nil).Twice()
		cpID.On("FromID", s.ctx, uint64(2)).Return(ethBTC, nil).Once()

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().NoError(err)
		s.Require().Len(prices, 2)

		// Check that the prices are correct
		s.Require().Equal(fourHundred.String(), prices[btcUSD].String())
		s.Require().Equal(sixHundred.String(), prices[ethUSD].String())
	})

	s.Run("errors when the validator's prices are malformed", func() {
		// Create a single vote extension from my validator
		myValPrices := map[uint64][]byte{
			0: ongodhecappin,
		}
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(s.myVal, myValPrices, s.veCodec)
		s.Require().NoError(err)

		// Create the extended commit info
		_, commitBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo}, s.commitCodec)
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().Error(err)
		s.Require().Len(prices, 0)
	})
}
