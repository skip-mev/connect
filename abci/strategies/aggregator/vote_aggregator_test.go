package aggregator_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/skip-mev/connect/v2/abci/strategies/aggregator"
	"github.com/skip-mev/connect/v2/abci/strategies/codec"
	currencypairmocks "github.com/skip-mev/connect/v2/abci/strategies/currencypair/mocks"
	"github.com/skip-mev/connect/v2/abci/testutils"
	"github.com/skip-mev/connect/v2/pkg/math/voteweighted"
	"github.com/skip-mev/connect/v2/pkg/math/voteweighted/mocks"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
)

var (
	btcUSD = connecttypes.CurrencyPair{
		Base:  "BTC",
		Quote: "USD",
	}
	ethUSD = connecttypes.CurrencyPair{
		Base:  "ETH",
		Quote: "USD",
	}
	ethBTC = connecttypes.CurrencyPair{
		Base:  "ETH",
		Quote: "BTC",
	}

	oneHundred   = big.NewInt(100)
	twoHundred   = big.NewInt(200)
	threeHundred = big.NewInt(300)
	fourHundred  = big.NewInt(400)
	sixHundred   = big.NewInt(600)
	sevenHundred = big.NewInt(700)
	eightHundred = big.NewInt(800)
	nineHundred  = big.NewInt(900)

	val1 = sdk.ConsAddress("val1")
	val2 = sdk.ConsAddress("val2")

	ongodhecappin = append([]byte("ongodhecappin"), make([]byte, 32)...)
)

type VoteAggregatorTestSuite struct {
	suite.Suite

	commitCodec codec.ExtendedCommitCodec

	ctx sdk.Context

	veCodec codec.VoteExtensionCodec

	myVal sdk.ConsAddress
}

func (s *VoteAggregatorTestSuite) SetupTest() {
	s.myVal = sdk.ConsAddress("myVal")

	s.veCodec = codec.NewCompressionVoteExtensionCodec(
		codec.NewDefaultVoteExtensionCodec(),
		codec.NewZLibCompressor(),
	)

	s.commitCodec = codec.NewCompressionExtendedCommitCodec(
		codec.NewDefaultExtendedCommitCodec(),
		codec.NewZLibCompressor(),
	)
}

func TestVoteAggregatorTestSuite(t *testing.T) {
	suite.Run(t, new(VoteAggregatorTestSuite))
}

func (s *VoteAggregatorTestSuite) TestAggregateOracleVotes() {
	// Use the default aggregation function for testing
	mockValidatorStore := mocks.NewValidatorStore(s.T())
	aggregationFn := voteweighted.MedianFromContext(
		log.NewTestLogger(s.T()),
		mockValidatorStore,
		voteweighted.DefaultPowerThreshold,
	)
	mockValidatorStore.On("TotalBondedTokens", mock.Anything).Return(math.NewInt(100), nil)

	cpID := currencypairmocks.NewCurrencyPairStrategy(s.T())

	handler := aggregator.NewDefaultVoteAggregator(
		log.NewTestLogger(s.T()),
		aggregationFn,
		cpID,
	)

	s.Run("no oracle data", func() {
		_, commitBz, err := testutils.CreateExtendedCommitInfo(nil, s.commitCodec)
		s.Require().NoError(err)

		proposal := [][]byte{commitBz}
		votes, err := aggregator.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

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
		votes, err := aggregator.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

		cpID.On("FromID", s.ctx, uint64(0)).Return(btcUSD, nil).Once()
		cpID.On("GetDecodedPrice", s.ctx, btcUSD, oneHundred.Bytes()).Return(oneHundred, nil).Once()

		// Assume the validator takes up all voting power
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
		votes, err := aggregator.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

		// Assume the validator takes up all voting power
		mockValidatorStore.On("ValidatorByConsAddr", mock.Anything, s.myVal).Return(
			stakingtypes.Validator{
				Tokens: math.NewInt(100),
				Status: stakingtypes.Bonded,
			},
			nil,
		).Once()

		cpID.On("FromID", s.ctx, uint64(0)).Return(btcUSD, nil).Once()
		cpID.On("GetDecodedPrice", s.ctx, btcUSD, oneHundred.Bytes()).Return(oneHundred, nil).Once()

		cpID.On("FromID", s.ctx, uint64(1)).Return(ethUSD, nil).Once()
		cpID.On("GetDecodedPrice", s.ctx, ethUSD, twoHundred.Bytes()).Return(twoHundred, nil).Once()

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
		votes, err := aggregator.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

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
		cpID.On("GetDecodedPrice", s.ctx, btcUSD, oneHundred.Bytes()).Return(oneHundred, nil).Once()
		cpID.On("GetDecodedPrice", s.ctx, btcUSD, twoHundred.Bytes()).Return(twoHundred, nil).Once()

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
		votes, err := aggregator.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

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
		cpID.On("GetDecodedPrice", s.ctx, btcUSD, oneHundred.Bytes()).Return(oneHundred, nil).Once()
		cpID.On("GetDecodedPrice", s.ctx, btcUSD, twoHundred.Bytes()).Return(twoHundred, nil).Once()

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
		votes, err := aggregator.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

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

		cpID.On("GetDecodedPrice", s.ctx, btcUSD, oneHundred.Bytes()).Return(oneHundred, nil).Once()
		cpID.On("GetDecodedPrice", s.ctx, btcUSD, threeHundred.Bytes()).Return(threeHundred, nil).Once()

		cpID.On("GetDecodedPrice", s.ctx, ethUSD, twoHundred.Bytes()).Return(twoHundred, nil).Once()
		cpID.On("GetDecodedPrice", s.ctx, ethUSD, fourHundred.Bytes()).Return(fourHundred, nil).Once()

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
		votes, err := aggregator.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

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

		cpID.On("FromID", s.ctx, uint64(0)).Return(btcUSD, nil).Times(3)
		cpID.On("FromID", s.ctx, uint64(1)).Return(ethUSD, nil).Twice()
		cpID.On("FromID", s.ctx, uint64(2)).Return(ethBTC, nil).Times(3)

		cpID.On("GetDecodedPrice", s.ctx, btcUSD, oneHundred.Bytes()).Return(oneHundred, nil).Once()
		cpID.On("GetDecodedPrice", s.ctx, btcUSD, fourHundred.Bytes()).Return(fourHundred, nil).Once()
		cpID.On("GetDecodedPrice", s.ctx, btcUSD, sevenHundred.Bytes()).Return(sevenHundred, nil).Once()

		cpID.On("GetDecodedPrice", s.ctx, ethUSD, twoHundred.Bytes()).Return(twoHundred, nil).Once()
		cpID.On("GetDecodedPrice", s.ctx, ethUSD, eightHundred.Bytes()).Return(eightHundred, nil).Once()

		cpID.On("GetDecodedPrice", s.ctx, ethBTC, threeHundred.Bytes()).Return(threeHundred, nil).Once()
		cpID.On("GetDecodedPrice", s.ctx, ethBTC, sixHundred.Bytes()).Return(sixHundred, nil).Once()
		cpID.On("GetDecodedPrice", s.ctx, ethBTC, nineHundred.Bytes()).Return(nineHundred, nil).Once()

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().NoError(err)
		s.Require().Len(prices, 2)

		// Check that the prices are correct
		s.Require().Equal(fourHundred.String(), prices[btcUSD].String())
		s.Require().Equal(sixHundred.String(), prices[ethBTC].String())
	})

	s.Run("continues when the validator's prices are malformed", func() {
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
		votes, err := aggregator.GetOracleVotes(proposal, s.veCodec, s.commitCodec)
		s.Require().NoError(err)

		mockValidatorStore.On("ValidatorByConsAddr", mock.Anything, s.myVal).Return(
			stakingtypes.Validator{
				Tokens: math.NewInt(25),
				Status: stakingtypes.Bonded,
			},
			nil,
		).Once()

		// Aggregate oracle data
		prices, err := handler.AggregateOracleVotes(s.ctx, votes)
		s.Require().NoError(err)
		s.Require().Len(prices, 0)
	})
}
