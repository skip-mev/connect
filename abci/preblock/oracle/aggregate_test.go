package oracle_test

import (
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/mock"

	preblock "github.com/skip-mev/slinky/abci/preblock/oracle"
	preblockmath "github.com/skip-mev/slinky/abci/preblock/oracle/math"
	"github.com/skip-mev/slinky/abci/preblock/oracle/math/mocks"
	preblockmock "github.com/skip-mev/slinky/abci/preblock/oracle/mocks"
	"github.com/skip-mev/slinky/abci/testutils"
	merticmock "github.com/skip-mev/slinky/service/metrics/mocks"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	cp1 = oracletypes.CurrencyPair{
		Base:  "BTC",
		Quote: "USD",
	}
	cp2 = oracletypes.CurrencyPair{
		Base:  "ETH",
		Quote: "USD",
	}
	cp3 = oracletypes.CurrencyPair{
		Base:  "ETH",
		Quote: "BTC",
	}

	val1 = sdk.ConsAddress([]byte("val1"))
	val2 = sdk.ConsAddress([]byte("val2"))
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
	mockMetrics := merticmock.NewMetrics(s.T())

	handler := preblock.NewOraclePreBlockHandler(
		log.NewTestLogger(s.T()),
		aggregationFn,
		mockOracleKeeper,
		s.myVal,
		mockMetrics,
	)

	s.Run("no oracle data", func() {
		_, commitBz, err := testutils.CreateExtendedCommitInfo(nil)
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal)
		s.Require().NoError(err)

		mockMetrics.On("AddVoteIncludedInLastCommit", false).Once()

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().NoError(err)
		s.Require().Len(prices, 0)
	})

	s.Run("single oracle data", func() {
		// Create a single vote extension from my validator
		myValPrices := map[string]string{
			cp1.ToString(): "0x100",
		}
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(s.myVal, myValPrices, time.Now(), 2)
		s.Require().NoError(err)

		// Create the extended commit info
		_, commitBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo})
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal)
		s.Require().NoError(err)

		// The validator is included in the commit and the price should be included
		mockMetrics.On("AddVoteIncludedInLastCommit", true).Once()
		mockMetrics.On("AddTickerInclusionStatus", cp1.ToString(), true).Once()

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
		s.Require().Equal(myValPrices[cp1.ToString()], prices[cp1].String())
	})

	s.Run("multiple prices from a single validator", func() {
		// Create a single vote extension from my validator
		myValPrices := map[string]string{
			cp1.ToString(): "0x100",
			cp2.ToString(): "0x200",
		}
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(s.myVal, myValPrices, time.Now(), 2)
		s.Require().NoError(err)

		// Create the extended commit info
		_, commitBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo})
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal)
		s.Require().NoError(err)

		// The validator is included in the commit and the price should be included
		mockMetrics.On("AddVoteIncludedInLastCommit", true).Once()
		mockMetrics.On("AddTickerInclusionStatus", cp1.ToString(), true).Once()
		mockMetrics.On("AddTickerInclusionStatus", cp2.ToString(), true).Once()

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
		s.Require().Len(prices, 2)

		// Check that the prices are correct
		s.Require().Equal(myValPrices[cp1.ToString()], prices[cp1].String())
		s.Require().Equal(myValPrices[cp2.ToString()], prices[cp2].String())
	})

	s.Run("single price from two different validators", func() {
		// Create a single vote extension from my validator
		myValPrices := map[string]string{
			cp1.ToString(): "0x100",
		}
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(s.myVal, myValPrices, time.Now(), 2)
		s.Require().NoError(err)

		// Create a single vote extension from another validator
		otherValPrices := map[string]string{
			cp1.ToString(): "0x200",
		}
		otherValVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, otherValPrices, time.Now(), 2)
		s.Require().NoError(err)

		// Create the extended commit info
		_, commitBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo, otherValVoteInfo})
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal)
		s.Require().NoError(err)

		// The validator is included in the commit and the price should be included
		mockMetrics.On("AddVoteIncludedInLastCommit", true).Once()
		mockMetrics.On("AddTickerInclusionStatus", cp1.ToString(), true).Once()

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

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().NoError(err)
		s.Require().Len(prices, 1)

		// Check that the prices are correct
		s.Require().Equal("0x100", prices[cp1].String())
	})

	s.Run("single price update from multiple validators but not enough voting power", func() {
		// Create a single vote extension from my validator
		myValPrices := map[string]string{
			cp1.ToString(): "0x100",
		}
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(s.myVal, myValPrices, time.Now(), 2)
		s.Require().NoError(err)

		// Create a single vote extension from another validator
		otherValPrices := map[string]string{
			cp1.ToString(): "0x200",
		}
		otherValVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, otherValPrices, time.Now(), 2)
		s.Require().NoError(err)

		// Create the extended commit info
		_, commitBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo, otherValVoteInfo})
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal)
		s.Require().NoError(err)

		// The validator is included in the commit and the price should be included
		mockMetrics.On("AddVoteIncludedInLastCommit", true).Once()
		mockMetrics.On("AddTickerInclusionStatus", cp1.ToString(), true).Once()

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

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().NoError(err)
		s.Require().Len(prices, 0)
	})

	s.Run("multiple prices from multiple validators", func() {
		// Create a vote extension with multiple prices from my validator
		myValPrices := map[string]string{
			cp1.ToString(): "0x100",
			cp2.ToString(): "0x200",
		}
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(s.myVal, myValPrices, time.Now(), 2)
		s.Require().NoError(err)

		// Create a vote extension with multiple prices from another validator
		otherValPrices := map[string]string{
			cp1.ToString(): "0x300",
			cp2.ToString(): "0x400",
		}
		otherValVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, otherValPrices, time.Now(), 2)
		s.Require().NoError(err)

		// Create the extended commit info
		_, commitBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo, otherValVoteInfo})
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal)
		s.Require().NoError(err)

		// The validator is included in the commit and the price should be included
		mockMetrics.On("AddVoteIncludedInLastCommit", true).Once()
		mockMetrics.On("AddTickerInclusionStatus", cp1.ToString(), true).Once()
		mockMetrics.On("AddTickerInclusionStatus", cp2.ToString(), true).Once()

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

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().NoError(err)
		s.Require().Len(prices, 2)

		// Check that the prices are correct
		s.Require().Equal("0x300", prices[cp1].String())
		s.Require().Equal("0x400", prices[cp2].String())
	})

	s.Run("multiple prices from multiple validators but not enough voting power for some", func() {
		myValPrices := map[string]string{
			cp1.ToString(): "0x100",
			cp2.ToString(): "0x200",
			cp3.ToString(): "0x300",
		}
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(s.myVal, myValPrices, time.Now(), 2)
		s.Require().NoError(err)

		val1Prices := map[string]string{
			cp1.ToString(): "0x400",
			cp3.ToString(): "0x600",
		}
		val1VoteInfo, err := testutils.CreateExtendedVoteInfo(val1, val1Prices, time.Now(), 2)
		s.Require().NoError(err)

		val2Prices := map[string]string{
			cp1.ToString(): "0x700",
			cp2.ToString(): "0x800",
			cp3.ToString(): "0x900",
		}
		val2VoteInfo, err := testutils.CreateExtendedVoteInfo(val2, val2Prices, time.Now(), 2)
		s.Require().NoError(err)

		// Create the extended commit info
		_, commitBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo, val1VoteInfo, val2VoteInfo})
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal)
		s.Require().NoError(err)

		// The validator is included in the commit and the price should be included
		mockMetrics.On("AddVoteIncludedInLastCommit", true).Once()
		mockMetrics.On("AddTickerInclusionStatus", cp1.ToString(), true).Once()
		mockMetrics.On("AddTickerInclusionStatus", cp2.ToString(), true).Once()
		mockMetrics.On("AddTickerInclusionStatus", cp3.ToString(), true).Once()

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

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().NoError(err)
		s.Require().Len(prices, 2)

		// Check that the prices are correct
		s.Require().Equal("0x400", prices[cp1].String())
		s.Require().Equal("0x600", prices[cp3].String())
	})

	s.Run("errors when the validator's prices are malformed", func() {
		// Create a single vote extension from my validator
		myValPrices := map[string]string{
			cp1.ToString(): "ongodhecappin",
		}
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(s.myVal, myValPrices, time.Now(), 2)
		s.Require().NoError(err)

		// Create the extended commit info
		_, commitBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo})
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := preblock.GetOracleVotes(proposal)
		s.Require().NoError(err)

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().Error(err)
		s.Require().Len(prices, 0)
	})
}
