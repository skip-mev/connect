package sla_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/skip-mev/connect/v2/abci/preblock/sla"
	"github.com/skip-mev/connect/v2/abci/preblock/sla/mocks"
	voteaggregator "github.com/skip-mev/connect/v2/abci/strategies/aggregator"
	compression "github.com/skip-mev/connect/v2/abci/strategies/codec"
	currencypairmocks "github.com/skip-mev/connect/v2/abci/strategies/currencypair/mocks"
	"github.com/skip-mev/connect/v2/abci/testutils"
	oraclevetypes "github.com/skip-mev/connect/v2/abci/ve/types"
	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	slakeeper "github.com/skip-mev/connect/v2/x/sla/keeper"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
	slamocks "github.com/skip-mev/connect/v2/x/sla/types/mocks"
)

var oneHundred = big.NewInt(100)

type SLAPreBlockerHandlerTestSuite struct {
	suite.Suite

	ctx sdk.Context

	// The handler being tested.
	handler *sla.PreBlockHandler

	// Mocks for the handler's dependencies.
	oracleKeeper           *mocks.OracleKeeper
	stakingKeeper          *mocks.StakingKeeper
	slaKeeper              *slakeeper.Keeper
	currencyPairIDStrategy *currencypairmocks.CurrencyPairStrategy

	veCodec        compression.VoteExtensionCodec
	extCommitCodec compression.ExtendedCommitCodec

	val1      stakingtypes.Validator
	consAddr1 sdk.ConsAddress

	val2      stakingtypes.Validator
	consAddr2 sdk.ConsAddress

	val3      stakingtypes.Validator
	consAddr3 sdk.ConsAddress

	cp1 slinkytypes.CurrencyPair
	cp2 slinkytypes.CurrencyPair
	cp3 slinkytypes.CurrencyPair

	veEnabled bool

	sla    slatypes.PriceFeedSLA
	setSLA bool
}

func TestSLAPreBlockerHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SLAPreBlockerHandlerTestSuite))
}

func (s *SLAPreBlockerHandlerTestSuite) SetupTest() {
	pks := simtestutil.CreateTestPubKeys(4)

	var err error
	val1 := sdk.ValAddress("val1")
	s.val1, err = stakingtypes.NewValidator(val1.String(), pks[0], stakingtypes.Description{})
	s.Require().NoError(err)
	s.consAddr1, err = s.val1.GetConsAddr()
	s.Require().NoError(err)

	val2 := sdk.ValAddress("val2")
	s.val2, err = stakingtypes.NewValidator(val2.String(), pks[1], stakingtypes.Description{})
	s.Require().NoError(err)
	s.consAddr2, err = s.val2.GetConsAddr()
	s.Require().NoError(err)

	val3 := sdk.ValAddress("val3")
	s.val3, err = stakingtypes.NewValidator(val3.String(), pks[2], stakingtypes.Description{})
	s.Require().NoError(err)
	s.consAddr3, err = s.val3.GetConsAddr()
	s.Require().NoError(err)

	s.cp1 = slinkytypes.NewCurrencyPair("btc", "usd")
	s.cp2 = slinkytypes.NewCurrencyPair("eth", "usd")
	s.cp3 = slinkytypes.NewCurrencyPair("btc", "eth")

	s.currencyPairIDStrategy = currencypairmocks.NewCurrencyPairStrategy(s.T())

	s.veCodec = compression.NewCompressionVoteExtensionCodec(
		compression.NewDefaultVoteExtensionCodec(),
		compression.NewZLibCompressor(),
	)

	s.extCommitCodec = compression.NewCompressionExtendedCommitCodec(
		compression.NewDefaultExtendedCommitCodec(),
		compression.NewZLibCompressor(),
	)

	// Set a single sla in the store for subsequent testing.
	s.initHandler(s.veEnabled, s.setSLA)
}

func (s *SLAPreBlockerHandlerTestSuite) SetupSubTest() {
	s.initHandler(s.veEnabled, s.setSLA)
}

func (s *SLAPreBlockerHandlerTestSuite) TestPreBlocker() {
	s.Run("returns error if req is nil", func() {
		_, err := s.handler.PreBlocker()(s.ctx, nil)
		s.Require().Error(err)
	})

	s.Run("returns if vote extensions have not been enabled", func() {
		req := &cometabci.RequestFinalizeBlock{}

		_, err := s.handler.PreBlocker()(s.ctx, req)
		s.Require().NoError(err)
	})

	// Enable vote extensions.
	s.veEnabled = true

	s.Run("returns an error if the vote extensions are not included", func() {
		req := &cometabci.RequestFinalizeBlock{}
		_, err := s.handler.PreBlocker()(s.ctx, req)
		s.Require().Error(err)
	})

	s.Run("returns with no vote extensions in the block", func() {
		// create the vote extensions
		_, bz, err := testutils.CreateExtendedCommitInfo(nil, s.extCommitCodec)
		s.Require().NoError(err)

		// create the request
		req := &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{bz},
		}

		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{})

		_, err = s.handler.PreBlocker()(s.ctx, req)
		s.Require().NoError(err)

		cps, err := s.slaKeeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(0, len(cps))
	})

	s.setSLA = true

	s.Run("returns with no vote extensions in the block with a single sla set", func() {
		// create the vote extensions
		_, bz, err := testutils.CreateExtendedCommitInfo(nil, s.extCommitCodec)
		s.Require().NoError(err)

		// create the request
		req := &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{bz},
		}

		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{})

		_, err = s.handler.PreBlocker()(s.ctx, req)
		s.Require().NoError(err)

		// Check that no currency pairs were added.
		cps, err := s.slaKeeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(0, len(cps))

		// Check that no new price feeds were created.
		feeds, err := s.slaKeeper.GetAllPriceFeeds(s.ctx, s.sla.ID)
		s.Require().NoError(err)
		s.Require().Equal(0, len(feeds))
	})

	s.Run("returns with single validator, single cp, and vote with price", func() {
		ve1, err := testutils.CreateExtendedVoteInfo(
			s.consAddr1,
			map[uint64][]byte{
				0: oneHundred.Bytes(),
			},
			s.veCodec,
		)
		s.Require().NoError(err)

		// create the vote extensions
		_, bz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{ve1}, s.extCommitCodec)
		s.Require().NoError(err)

		// create the request
		req := &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{bz},
		}

		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{s.val1}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{s.cp1})

		s.currencyPairIDStrategy.On("ID", mock.Anything, s.cp1).Return(uint64(0), nil)

		_, err = s.handler.PreBlocker()(s.ctx, req)
		s.Require().NoError(err)

		// Check that the currency pair was added.
		cps, err := s.slaKeeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(1, len(cps))

		// Check that the price feed was created.
		feeds, err := s.slaKeeper.GetAllPriceFeeds(s.ctx, s.sla.ID)
		s.Require().NoError(err)
		s.Require().Equal(1, len(feeds))

		// Check that the price feed was created with the correct values.
		feed := feeds[0]
		s.Require().Equal(s.cp1, feed.CurrencyPair)
		s.Require().Equal(s.consAddr1, sdk.ConsAddress(feed.Validator))

		updates, err := feed.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), updates)

		votes, err := feed.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), votes)
	})

	s.Run("correctly updates with single validator, single cp, and validator with no vote extension", func() {
		// create the vote extensions
		_, bz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{}, s.extCommitCodec)
		s.Require().NoError(err)

		// create the request
		req := &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{bz},
		}

		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{s.val1}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{s.cp1})

		_, err = s.handler.PreBlocker()(s.ctx, req)
		s.Require().NoError(err)

		// Check that the currency pair was added.
		cps, err := s.slaKeeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(1, len(cps))

		// Check that the price feed was created.
		feeds, err := s.slaKeeper.GetAllPriceFeeds(s.ctx, s.sla.ID)
		s.Require().NoError(err)
		s.Require().Equal(1, len(feeds))

		// Check that the price feed was created with the correct values.
		feed := feeds[0]
		s.Require().Equal(s.cp1, feed.CurrencyPair)
		s.Require().Equal(s.consAddr1, sdk.ConsAddress(feed.Validator))

		updates, err := feed.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(0), updates)

		votes, err := feed.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(0), votes)
	})

	s.Run("correctly updates with single validator, single cp, and vote extension without the price", func() {
		ve1, err := testutils.CreateExtendedVoteInfo(
			s.consAddr1,
			map[uint64][]byte{},
			s.veCodec,
		)
		s.Require().NoError(err)

		// create the vote extensions
		_, bz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{ve1}, s.extCommitCodec)
		s.Require().NoError(err)

		// create the request
		req := &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{bz},
		}

		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{s.val1}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{s.cp1})

		_, err = s.handler.PreBlocker()(s.ctx, req)
		s.Require().NoError(err)

		// Check that the currency pair was added.
		cps, err := s.slaKeeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(1, len(cps))

		// Check that the price feed was created.
		feeds, err := s.slaKeeper.GetAllPriceFeeds(s.ctx, s.sla.ID)
		s.Require().NoError(err)
		s.Require().Equal(1, len(feeds))

		// Check that the price feed was created with the correct values.
		feed := feeds[0]
		s.Require().Equal(s.cp1, feed.CurrencyPair)
		s.Require().Equal(s.consAddr1, sdk.ConsAddress(feed.Validator))

		updates, err := feed.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(0), updates)

		votes, err := feed.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), votes)
	})

	s.Run("correctly updates an existing price feed with single validator, single cp, and vote extension with price", func() {
		feed, err := slatypes.NewPriceFeed(
			uint(s.sla.MaximumViableWindow),
			s.consAddr1,
			s.cp1,
			s.sla.ID,
		)
		s.Require().NoError(err)

		err = feed.SetUpdate(slatypes.VoteWithPrice)
		s.Require().NoError(err)

		// Add the feed to the store.
		err = s.slaKeeper.SetPriceFeed(s.ctx, feed)
		s.Require().NoError(err)

		ve1, err := testutils.CreateExtendedVoteInfo(
			s.consAddr1,
			map[uint64][]byte{
				0: oneHundred.Bytes(),
			},
			s.veCodec,
		)
		s.Require().NoError(err)

		// create the vote extensions
		_, bz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{ve1}, s.extCommitCodec)
		s.Require().NoError(err)

		// create the request
		req := &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{bz},
		}

		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{s.val1}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{s.cp1})

		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp1).Return(uint64(0), nil)

		_, err = s.handler.PreBlocker()(s.ctx, req)
		s.Require().NoError(err)

		// Check that the currency pair was added.
		cps, err := s.slaKeeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(1, len(cps))

		// Check that the price feed was updated.
		feeds, err := s.slaKeeper.GetAllPriceFeeds(s.ctx, s.sla.ID)
		s.Require().NoError(err)
		s.Require().Equal(1, len(feeds))

		// Check that the price feed was updated with the correct values.
		feed = feeds[0]
		s.Require().Equal(s.cp1, feed.CurrencyPair)
		s.Require().Equal(s.consAddr1, sdk.ConsAddress(feed.Validator))

		updates, err := feed.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), updates)

		votes, err := feed.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), votes)

		updates, err = feed.GetNumPriceUpdatesWithWindow(2)
		s.Require().NoError(err)
		s.Require().Equal(uint(2), updates)

		votes, err = feed.GetNumVotesWithWindow(2)
		s.Require().NoError(err)
		s.Require().Equal(uint(2), votes)
	})

	s.Run("correctly updates an existing price feed with single validator, single cp, and no vote extension", func() {
		feed, err := slatypes.NewPriceFeed(
			uint(s.sla.MaximumViableWindow),
			s.consAddr1,
			s.cp1,
			s.sla.ID,
		)
		s.Require().NoError(err)

		err = feed.SetUpdate(slatypes.VoteWithPrice)
		s.Require().NoError(err)

		// Add the feed to the store.
		err = s.slaKeeper.SetPriceFeed(s.ctx, feed)
		s.Require().NoError(err)

		// create the vote extensions
		_, bz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{}, s.extCommitCodec)
		s.Require().NoError(err)

		// create the request
		req := &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{bz},
		}

		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{s.val1}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{s.cp1})

		_, err = s.handler.PreBlocker()(s.ctx, req)
		s.Require().NoError(err)

		// Check that the currency pair was added.
		cps, err := s.slaKeeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(1, len(cps))

		// Check that the price feed was updated.
		feeds, err := s.slaKeeper.GetAllPriceFeeds(s.ctx, s.sla.ID)
		s.Require().NoError(err)
		s.Require().Equal(1, len(feeds))

		// Check that the price feed was updated with the correct values.
		feed = feeds[0]
		s.Require().Equal(s.cp1, feed.CurrencyPair)
		s.Require().Equal(s.consAddr1, sdk.ConsAddress(feed.Validator))

		updates, err := feed.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(0), updates)

		votes, err := feed.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(0), votes)

		updates, err = feed.GetNumPriceUpdatesWithWindow(2)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), updates)

		votes, err = feed.GetNumVotesWithWindow(2)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), votes)
	})

	s.Run("correctly updates an existing price feed with single validator, single cp, and vote extension without price", func() {
		feed, err := slatypes.NewPriceFeed(
			uint(s.sla.MaximumViableWindow),
			s.consAddr1,
			s.cp1,
			s.sla.ID,
		)
		s.Require().NoError(err)

		err = feed.SetUpdate(slatypes.VoteWithPrice)
		s.Require().NoError(err)

		// Add the feed to the store.
		err = s.slaKeeper.SetPriceFeed(s.ctx, feed)
		s.Require().NoError(err)

		ve1, err := testutils.CreateExtendedVoteInfo(
			s.consAddr1,
			map[uint64][]byte{},
			s.veCodec,
		)
		s.Require().NoError(err)

		// create the vote extensions
		_, bz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{ve1}, s.extCommitCodec)
		s.Require().NoError(err)

		// create the request
		req := &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{bz},
		}

		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{s.val1}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{s.cp1})

		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp1).Return(uint64(0), nil)

		_, err = s.handler.PreBlocker()(s.ctx, req)
		s.Require().NoError(err)

		// Check that the currency pair was added.
		cps, err := s.slaKeeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(1, len(cps))

		// Check that the price feed was updated.
		feeds, err := s.slaKeeper.GetAllPriceFeeds(s.ctx, s.sla.ID)
		s.Require().NoError(err)
		s.Require().Equal(1, len(feeds))

		// Check that the price feed was updated with the correct values.
		feed = feeds[0]
		s.Require().Equal(s.cp1, feed.CurrencyPair)
		s.Require().Equal(s.consAddr1, sdk.ConsAddress(feed.Validator))

		updates, err := feed.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(0), updates)

		votes, err := feed.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), votes)

		updates, err = feed.GetNumPriceUpdatesWithWindow(2)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), updates)

		votes, err = feed.GetNumVotesWithWindow(2)
		s.Require().NoError(err)
		s.Require().Equal(uint(2), votes)
	})
}

func (s *SLAPreBlockerHandlerTestSuite) TestGetUpdates() {
	s.Run("returns with no voting updates", func() {
		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{})

		updates, err := s.handler.GetUpdates(s.ctx, nil)
		s.Require().NoError(err)

		s.Require().Equal(0, len(updates.CurrencyPairs))
		s.Require().Equal(0, len(updates.ValidatorUpdates))
	})

	s.Run("returns single validator, with single currency pair, and no updates", func() {
		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{s.val1}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{s.cp1})

		updates, err := s.handler.GetUpdates(s.ctx, nil)
		s.Require().NoError(err)

		s.Require().Equal(1, len(updates.CurrencyPairs))
		s.Require().Equal(1, len(updates.ValidatorUpdates))
		s.Require().Contains(updates.CurrencyPairs, s.cp1)
		s.Require().Contains(updates.ValidatorUpdates, s.consAddr1.String())

		// Ensure the values are correct.
		validator := updates.ValidatorUpdates[s.consAddr1.String()]
		s.Require().Equal(validator.ConsAddress, s.consAddr1)
		s.Require().Equal(1, len(validator.Updates))

		s.Require().Contains(validator.Updates, s.cp1)
		s.Require().Equal(slatypes.NoVote, validator.Updates[s.cp1])
	})

	s.Run("returns with single validator, single cp, and vote with price", func() {
		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{s.val1}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{s.cp1})

		votes := []voteaggregator.Vote{
			{
				ConsAddress: s.consAddr1,
				OracleVoteExtension: oraclevetypes.OracleVoteExtension{
					Prices: map[uint64][]byte{
						0: oneHundred.Bytes(),
					},
				},
			},
		}

		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp1).Return(uint64(0), nil)

		updates, err := s.handler.GetUpdates(s.ctx, votes)
		s.Require().NoError(err)

		s.Require().Equal(1, len(updates.CurrencyPairs))
		s.Require().Equal(1, len(updates.ValidatorUpdates))
		s.Require().Contains(updates.CurrencyPairs, s.cp1)
		s.Require().Contains(updates.ValidatorUpdates, s.consAddr1.String())

		// Ensure the values are correct.
		validator := updates.ValidatorUpdates[s.consAddr1.String()]
		s.Require().Equal(validator.ConsAddress, s.consAddr1)

		s.Require().Contains(validator.Updates, s.cp1)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp1])
	})

	s.Run("returns with single validator, single cp, and vote without price", func() {
		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{s.val1}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{s.cp1})

		votes := []voteaggregator.Vote{
			{
				ConsAddress: s.consAddr1,
				OracleVoteExtension: oraclevetypes.OracleVoteExtension{
					Prices: map[uint64][]byte{},
				},
			},
		}

		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp1).Return(uint64(0), nil)

		updates, err := s.handler.GetUpdates(s.ctx, votes)
		s.Require().NoError(err)

		s.Require().Equal(1, len(updates.CurrencyPairs))
		s.Require().Equal(1, len(updates.ValidatorUpdates))
		s.Require().Contains(updates.CurrencyPairs, s.cp1)
		s.Require().Contains(updates.ValidatorUpdates, s.consAddr1.String())

		// Ensure the values are correct.
		validator := updates.ValidatorUpdates[s.consAddr1.String()]
		s.Require().Equal(validator.ConsAddress, s.consAddr1)

		s.Require().Contains(validator.Updates, s.cp1)
		s.Require().Equal(slatypes.VoteWithoutPrice, validator.Updates[s.cp1])
	})

	s.Run("returns with single validator, multiple cps, and votes with Prices", func() {
		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{s.val1}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{s.cp1, s.cp2, s.cp3})

		votes := []voteaggregator.Vote{
			{
				ConsAddress: s.consAddr1,
				OracleVoteExtension: oraclevetypes.OracleVoteExtension{
					Prices: map[uint64][]byte{
						0: oneHundred.Bytes(),
						1: oneHundred.Bytes(),
						2: oneHundred.Bytes(),
					},
				},
			},
		}

		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp1).Return(uint64(0), nil)
		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp2).Return(uint64(1), nil)
		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp3).Return(uint64(2), nil)

		updates, err := s.handler.GetUpdates(s.ctx, votes)
		s.Require().NoError(err)

		s.Require().Equal(3, len(updates.CurrencyPairs))
		s.Require().Equal(1, len(updates.ValidatorUpdates))
		s.Require().Contains(updates.CurrencyPairs, s.cp1)
		s.Require().Contains(updates.CurrencyPairs, s.cp2)
		s.Require().Contains(updates.CurrencyPairs, s.cp3)
		s.Require().Contains(updates.ValidatorUpdates, s.consAddr1.String())

		// Ensure the values are correct.
		validator := updates.ValidatorUpdates[s.consAddr1.String()]
		s.Require().Equal(validator.ConsAddress, s.consAddr1)
		s.Require().Equal(3, len(validator.Updates))

		s.Require().Contains(validator.Updates, s.cp1)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp1])

		s.Require().Contains(validator.Updates, s.cp2)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp2])

		s.Require().Contains(validator.Updates, s.cp3)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp3])
	})

	s.Run("returns with single validator, multiple cps, and some Prices", func() {
		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{s.val1}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{s.cp1, s.cp2, s.cp3})

		votes := []voteaggregator.Vote{
			{
				ConsAddress: s.consAddr1,
				OracleVoteExtension: oraclevetypes.OracleVoteExtension{
					Prices: map[uint64][]byte{
						0: oneHundred.Bytes(),
						1: oneHundred.Bytes(),
					},
				},
			},
		}

		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp1).Return(uint64(0), nil)
		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp2).Return(uint64(1), nil)
		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp3).Return(uint64(2), nil)

		updates, err := s.handler.GetUpdates(s.ctx, votes)
		s.Require().NoError(err)

		s.Require().Equal(3, len(updates.CurrencyPairs))
		s.Require().Equal(1, len(updates.ValidatorUpdates))
		s.Require().Contains(updates.CurrencyPairs, s.cp1)
		s.Require().Contains(updates.CurrencyPairs, s.cp2)
		s.Require().Contains(updates.CurrencyPairs, s.cp3)
		s.Require().Contains(updates.ValidatorUpdates, s.consAddr1.String())

		// Ensure the values are correct.
		validator := updates.ValidatorUpdates[s.consAddr1.String()]
		s.Require().Equal(validator.ConsAddress, s.consAddr1)
		s.Require().Equal(3, len(validator.Updates))

		s.Require().Contains(validator.Updates, s.cp1)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp1])

		s.Require().Contains(validator.Updates, s.cp2)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp2])

		s.Require().Contains(validator.Updates, s.cp3)
		s.Require().Equal(slatypes.VoteWithoutPrice, validator.Updates[s.cp3])
	})

	s.Run("returns with 2 validators, single cp, and votes with Prices", func() {
		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{s.val1, s.val2}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{s.cp1})

		votes := []voteaggregator.Vote{
			{
				ConsAddress: s.consAddr1,
				OracleVoteExtension: oraclevetypes.OracleVoteExtension{
					Prices: map[uint64][]byte{
						0: oneHundred.Bytes(),
					},
				},
			},
			{
				ConsAddress: s.consAddr2,
				OracleVoteExtension: oraclevetypes.OracleVoteExtension{
					Prices: map[uint64][]byte{
						0: oneHundred.Bytes(),
					},
				},
			},
		}

		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp1).Return(uint64(0), nil).Twice()

		updates, err := s.handler.GetUpdates(s.ctx, votes)
		s.Require().NoError(err)

		s.Require().Equal(1, len(updates.CurrencyPairs))
		s.Require().Equal(2, len(updates.ValidatorUpdates))
		s.Require().Contains(updates.CurrencyPairs, s.cp1)
		s.Require().Contains(updates.ValidatorUpdates, s.consAddr1.String())
		s.Require().Contains(updates.ValidatorUpdates, s.consAddr2.String())

		// Ensure the values are correct.
		validator := updates.ValidatorUpdates[s.consAddr1.String()]
		s.Require().Equal(validator.ConsAddress, s.consAddr1)
		s.Require().Equal(1, len(validator.Updates))

		s.Require().Contains(validator.Updates, s.cp1)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp1])

		validator = updates.ValidatorUpdates[s.consAddr2.String()]
		s.Require().Equal(validator.ConsAddress, s.consAddr2)
		s.Require().Equal(1, len(validator.Updates))

		s.Require().Contains(validator.Updates, s.cp1)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp1])
	})

	s.Run("multiple validators, single cp, and one validator did not post any price updates", func() {
		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{s.val1, s.val2}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{s.cp1})

		votes := []voteaggregator.Vote{
			{
				ConsAddress: s.consAddr1,
				OracleVoteExtension: oraclevetypes.OracleVoteExtension{
					Prices: map[uint64][]byte{
						0: oneHundred.Bytes(),
					},
				},
			},
		}

		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp1).Return(uint64(0), nil)

		updates, err := s.handler.GetUpdates(s.ctx, votes)
		s.Require().NoError(err)

		s.Require().Equal(1, len(updates.CurrencyPairs))
		s.Require().Equal(2, len(updates.ValidatorUpdates))
		s.Require().Contains(updates.CurrencyPairs, s.cp1)
		s.Require().Contains(updates.ValidatorUpdates, s.consAddr1.String())
		s.Require().Contains(updates.ValidatorUpdates, s.consAddr2.String())

		// Ensure the values are correct.
		validator := updates.ValidatorUpdates[s.consAddr1.String()]
		s.Require().Equal(1, len(validator.Updates))
		s.Require().Equal(validator.ConsAddress, s.consAddr1)

		s.Require().Contains(validator.Updates, s.cp1)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp1])

		validator = updates.ValidatorUpdates[s.consAddr2.String()]
		s.Require().Equal(1, len(validator.Updates))
		s.Require().Equal(validator.ConsAddress, s.consAddr2)

		s.Require().Contains(validator.Updates, s.cp1)
		s.Require().Equal(slatypes.NoVote, validator.Updates[s.cp1])
	})

	s.Run("multiple validators, multiple cps, and all validators posted Prices", func() {
		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{s.val1, s.val2}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{s.cp1, s.cp2, s.cp3})

		votes := []voteaggregator.Vote{
			{
				ConsAddress: s.consAddr1,
				OracleVoteExtension: oraclevetypes.OracleVoteExtension{
					Prices: map[uint64][]byte{
						0: oneHundred.Bytes(),
						1: oneHundred.Bytes(),
						2: oneHundred.Bytes(),
					},
				},
			},
			{
				ConsAddress: s.consAddr2,
				OracleVoteExtension: oraclevetypes.OracleVoteExtension{
					Prices: map[uint64][]byte{
						0: oneHundred.Bytes(),
						1: oneHundred.Bytes(),
						2: oneHundred.Bytes(),
					},
				},
			},
		}

		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp1).Return(uint64(0), nil).Twice()
		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp2).Return(uint64(1), nil).Twice()
		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp3).Return(uint64(2), nil).Twice()

		updates, err := s.handler.GetUpdates(s.ctx, votes)
		s.Require().NoError(err)

		s.Require().Equal(3, len(updates.CurrencyPairs))
		s.Require().Equal(2, len(updates.ValidatorUpdates))
		s.Require().Contains(updates.CurrencyPairs, s.cp1)
		s.Require().Contains(updates.CurrencyPairs, s.cp2)
		s.Require().Contains(updates.CurrencyPairs, s.cp3)
		s.Require().Contains(updates.ValidatorUpdates, s.consAddr1.String())
		s.Require().Contains(updates.ValidatorUpdates, s.consAddr2.String())

		// Ensure the values are correct.
		validator := updates.ValidatorUpdates[s.consAddr1.String()]
		s.Require().Equal(validator.ConsAddress, s.consAddr1)
		s.Require().Equal(3, len(validator.Updates))

		s.Require().Contains(validator.Updates, s.cp1)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp1])

		s.Require().Contains(validator.Updates, s.cp2)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp2])

		s.Require().Contains(validator.Updates, s.cp3)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp3])

		validator = updates.ValidatorUpdates[s.consAddr2.String()]
		s.Require().Equal(validator.ConsAddress, s.consAddr2)
		s.Require().Equal(3, len(validator.Updates))

		s.Require().Contains(validator.Updates, s.cp1)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp1])

		s.Require().Contains(validator.Updates, s.cp2)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp2])

		s.Require().Contains(validator.Updates, s.cp3)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp3])
	})

	s.Run("3 validators, 3 cps, 1 validator did not vote, 1 validator posted Prices for some, 1 posted for all", func() {
		s.stakingKeeper.On("GetBondedValidatorsByPower", s.ctx).Return([]stakingtypes.Validator{s.val1, s.val2, s.val3}, nil)
		s.oracleKeeper.On("GetAllCurrencyPairs", s.ctx).Return([]slinkytypes.CurrencyPair{s.cp1, s.cp2, s.cp3})

		votes := []voteaggregator.Vote{
			{
				ConsAddress: s.consAddr1,
				OracleVoteExtension: oraclevetypes.OracleVoteExtension{
					Prices: map[uint64][]byte{
						0: oneHundred.Bytes(),
						1: oneHundred.Bytes(),
						2: oneHundred.Bytes(),
					},
				},
			},
			{
				ConsAddress: s.consAddr2,
				OracleVoteExtension: oraclevetypes.OracleVoteExtension{
					Prices: map[uint64][]byte{
						0: oneHundred.Bytes(),
					},
				},
			},
		}

		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp1).Return(uint64(0), nil).Twice()
		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp2).Return(uint64(1), nil)
		s.currencyPairIDStrategy.On("ID", s.ctx, s.cp3).Return(uint64(2), nil)

		updates, err := s.handler.GetUpdates(s.ctx, votes)
		s.Require().NoError(err)

		s.Require().Equal(3, len(updates.CurrencyPairs))
		s.Require().Equal(3, len(updates.ValidatorUpdates))
		s.Require().Contains(updates.CurrencyPairs, s.cp1)
		s.Require().Contains(updates.CurrencyPairs, s.cp2)
		s.Require().Contains(updates.CurrencyPairs, s.cp3)
		s.Require().Contains(updates.ValidatorUpdates, s.consAddr1.String())
		s.Require().Contains(updates.ValidatorUpdates, s.consAddr2.String())
		s.Require().Contains(updates.ValidatorUpdates, s.consAddr3.String())

		// Ensure the values are correct.
		validator := updates.ValidatorUpdates[s.consAddr1.String()]
		s.Require().Equal(validator.ConsAddress, s.consAddr1)
		s.Require().Equal(3, len(validator.Updates))

		s.Require().Contains(validator.Updates, s.cp1)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp1])

		s.Require().Contains(validator.Updates, s.cp2)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp2])

		s.Require().Contains(validator.Updates, s.cp3)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp3])

		validator = updates.ValidatorUpdates[s.consAddr2.String()]
		s.Require().Equal(validator.ConsAddress, s.consAddr2)
		s.Require().Equal(3, len(validator.Updates))

		s.Require().Contains(validator.Updates, s.cp1)
		s.Require().Equal(slatypes.VoteWithPrice, validator.Updates[s.cp1])

		s.Require().Contains(validator.Updates, s.cp2)
		s.Require().Equal(slatypes.VoteWithoutPrice, validator.Updates[s.cp2])

		s.Require().Contains(validator.Updates, s.cp3)
		s.Require().Equal(slatypes.VoteWithoutPrice, validator.Updates[s.cp3])

		validator = updates.ValidatorUpdates[s.consAddr3.String()]
		s.Require().Equal(validator.ConsAddress, s.consAddr3)
		s.Require().Equal(3, len(validator.Updates))

		s.Require().Contains(validator.Updates, s.cp1)
		s.Require().Equal(slatypes.NoVote, validator.Updates[s.cp1])

		s.Require().Contains(validator.Updates, s.cp2)
		s.Require().Equal(slatypes.NoVote, validator.Updates[s.cp2])

		s.Require().Contains(validator.Updates, s.cp3)
		s.Require().Equal(slatypes.NoVote, validator.Updates[s.cp3])
	})
}

func (s *SLAPreBlockerHandlerTestSuite) initHandler(veEnabled, setSLA bool) {
	// Set up context
	key := storetypes.NewKVStoreKey(slatypes.StoreKey)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx

	if veEnabled {
		s.ctx = s.ctx.WithConsensusParams(cmtproto.ConsensusParams{
			Abci: &cmtproto.ABCIParams{VoteExtensionsEnableHeight: 1},
		})

		s.ctx = s.ctx.WithBlockHeight(3)
	}

	s.ctx = s.ctx.WithLogger(log.NewTestLogger(s.T()))

	// Set up for sla keeper
	// Set up store and encoding configs
	storeService := runtime.NewKVStoreService(key)
	encodingConfig := moduletestutil.MakeTestEncodingConfig()

	slatypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	// Set up keeper
	s.slaKeeper = slakeeper.NewKeeper(
		storeService,
		encodingConfig.Codec,
		sdk.AccAddress("authority"),
		slamocks.NewStakingKeeper(s.T()),
		slamocks.NewSlashingKeeper(s.T()),
	)

	s.Require().NoError(s.slaKeeper.SetParams(s.ctx, slatypes.DefaultParams()))

	if setSLA {
		id := "slaID1"
		s.sla = slatypes.NewPriceFeedSLA(
			id,
			10,
			math.LegacyMustNewDecFromStr("0.1"),
			math.LegacyMustNewDecFromStr("0.1"),
			5,
			uint64(s.ctx.BlockHeight()),
		)
		err := s.slaKeeper.SetSLA(s.ctx, s.sla)
		s.Require().NoError(err)

		slas, err := s.slaKeeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(1, len(slas))
		s.Require().Equal(s.sla, slas[0])
	}

	s.oracleKeeper = mocks.NewOracleKeeper(s.T())
	s.stakingKeeper = mocks.NewStakingKeeper(s.T())

	s.handler = sla.NewSLAPreBlockHandler(
		s.oracleKeeper,
		s.stakingKeeper,
		s.slaKeeper,
		s.currencyPairIDStrategy,
		s.veCodec,
		s.extCommitCodec,
	)
}
