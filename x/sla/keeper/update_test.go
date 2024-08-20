package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	slakeeper "github.com/skip-mev/connect/v2/x/sla/keeper"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

func (s *KeeperTestSuite) TestUpdatePriceFeeds() {
	id := "id"
	sla := slatypes.NewPriceFeedSLA(id, 10, math.LegacyMustNewDecFromStr("1.0"), math.LegacyMustNewDecFromStr("1.0"), 5, 5)

	consAddress1 := sdk.ConsAddress("consAddress1")

	cp := slinkytypes.NewCurrencyPair("btc", "usd")

	priceFeedUpdates := slakeeper.NewPriceFeedUpdates()
	priceFeedUpdates.CurrencyPairs[cp] = struct{}{}

	validatorUpdates := slakeeper.NewValidatorUpdate(consAddress1)
	validatorUpdates.Updates[cp] = slatypes.VoteWithPrice

	priceFeedUpdates.ValidatorUpdates[consAddress1.String()] = validatorUpdates

	cps := make(map[slinkytypes.CurrencyPair]struct{})
	cps[cp] = struct{}{}

	s.Run("correctly updates price feeds with updates", func() {
		err := s.keeper.SetSLA(s.ctx, sla)
		s.Require().NoError(err)

		cps, err := s.keeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(cps, 0)

		err = s.keeper.UpdatePriceFeeds(s.ctx, priceFeedUpdates)
		s.Require().NoError(err)

		// Check that the new price feed was added.
		priceFeeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id)
		s.Require().NoError(err)
		s.Require().Len(priceFeeds, 1)

		// Check that the currency pair was added.
		cps, err = s.keeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(cps, 1)
		s.Require().Contains(cps, cp)
	})

	s.Run("correctly updates price feeds with no updates", func() {
		err := s.keeper.SetSLA(s.ctx, sla)
		s.Require().NoError(err)

		err = s.keeper.UpdatePriceFeeds(s.ctx, slakeeper.NewPriceFeedUpdates())
		s.Require().NoError(err)

		// check that there are no price feeds.
		priceFeeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id)
		s.Require().NoError(err)
		s.Require().Len(priceFeeds, 0)

		// Check that no currency pair was added.
		cps, err := s.keeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(cps, 0)
	})

	s.Run("can remove price feeds", func() {
		err := s.keeper.SetSLA(s.ctx, sla)
		s.Require().NoError(err)

		feed, err := slatypes.NewPriceFeed(uint(sla.MaximumViableWindow), consAddress1, cp, sla.ID)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, feed)
		s.Require().NoError(err)

		err = s.keeper.SetCurrencyPairs(s.ctx, cps)
		s.Require().NoError(err)

		err = s.keeper.UpdatePriceFeeds(s.ctx, slakeeper.NewPriceFeedUpdates())
		s.Require().NoError(err)

		// Check that the price feed was removed.
		priceFeeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id)
		s.Require().NoError(err)
		s.Require().Len(priceFeeds, 0)

		// Check that the currency pair was removed.
		cps, err := s.keeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(cps, 0)
	})

	s.Run("can update price feeds", func() {
		err := s.keeper.SetSLA(s.ctx, sla)
		s.Require().NoError(err)

		feed, err := slatypes.NewPriceFeed(uint(sla.MaximumViableWindow), consAddress1, cp, sla.ID)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, feed)
		s.Require().NoError(err)

		err = s.keeper.SetCurrencyPairs(s.ctx, cps)
		s.Require().NoError(err)

		validatorUpdates.Updates[cp] = slatypes.VoteWithoutPrice
		priceFeedUpdates.ValidatorUpdates[consAddress1.String()] = validatorUpdates

		err = s.keeper.UpdatePriceFeeds(s.ctx, priceFeedUpdates)
		s.Require().NoError(err)

		// Check that the price feed was updated.
		priceFeeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id)
		s.Require().NoError(err)
		s.Require().Len(priceFeeds, 1)

		feed = priceFeeds[0]
		s.Require().Equal(consAddress1, sdk.ConsAddress(feed.Validator))
		s.Require().Equal(cp, feed.CurrencyPair)
		s.Require().Equal(sla.ID, feed.ID)
		s.Require().Equal(sla.MaximumViableWindow, feed.MaximumViableWindow)

		// Check that the status was correctly set.
		numVotes, err := feed.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numVotes)

		numPriceUpdates, err := feed.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(0), numPriceUpdates)
	})
}

func (s *KeeperTestSuite) TestUpdatePriceFeedsForSLA() {
	id := "id"
	sla := slatypes.NewPriceFeedSLA(id, 10, math.LegacyMustNewDecFromStr("1.0"), math.LegacyMustNewDecFromStr("1.0"), 5, 5)

	consAddress1 := sdk.ConsAddress("consAddress1")
	consAddress2 := sdk.ConsAddress("consAddress2")

	cp1 := slinkytypes.NewCurrencyPair("btc", "usd")
	cp2 := slinkytypes.NewCurrencyPair("eth", "usd")

	s.Run("correctly updates price feeds with no updates", func() {
		updates := slakeeper.NewPriceFeedUpdates()
		err := s.keeper.UpdatePriceFeedsForSLA(s.ctx, sla, updates)
		s.Require().NoError(err)

		// Check that no price feeds were added or removed.
		priceFeeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id)
		s.Require().NoError(err)
		s.Require().Len(priceFeeds, 0)
	})

	priceFeedUpdates := slakeeper.NewPriceFeedUpdates()
	priceFeedUpdates.CurrencyPairs[cp1] = struct{}{}

	valUpdates := slakeeper.NewValidatorUpdate(consAddress1)
	valUpdates.Updates[cp1] = slatypes.VoteWithPrice

	priceFeedUpdates.ValidatorUpdates[consAddress1.String()] = valUpdates

	s.Run("correctly can create a new price feed with price update", func() {
		err := s.keeper.UpdatePriceFeedsForSLA(s.ctx, sla, priceFeedUpdates)
		s.Require().NoError(err)

		// Check that the price feed was added.
		priceFeeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id)
		s.Require().NoError(err)
		s.Require().Len(priceFeeds, 1)

		feed := priceFeeds[0]
		s.Require().Equal(consAddress1, sdk.ConsAddress(feed.Validator))
		s.Require().Equal(cp1, feed.CurrencyPair)
		s.Require().Equal(sla.ID, feed.ID)
		s.Require().Equal(sla.MaximumViableWindow, feed.MaximumViableWindow)

		// Check that the status was correctly set.
		numVotes, err := feed.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numVotes)

		numPriceUpdates, err := feed.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numPriceUpdates)
	})

	s.Run("correctly can create a new price feed with vote but no price update", func() {
		valUpdates.Updates[cp1] = slatypes.VoteWithoutPrice
		priceFeedUpdates.ValidatorUpdates[consAddress1.String()] = valUpdates

		err := s.keeper.UpdatePriceFeedsForSLA(s.ctx, sla, priceFeedUpdates)
		s.Require().NoError(err)

		// Check that the price feed was added.
		priceFeeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id)
		s.Require().NoError(err)
		s.Require().Len(priceFeeds, 1)

		feed := priceFeeds[0]
		s.Require().Equal(consAddress1, sdk.ConsAddress(feed.Validator))
		s.Require().Equal(cp1, feed.CurrencyPair)
		s.Require().Equal(sla.ID, feed.ID)
		s.Require().Equal(sla.MaximumViableWindow, feed.MaximumViableWindow)

		// Check that the status was correctly set.
		numVotes, err := feed.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numVotes)

		numPriceUpdates, err := feed.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(0), numPriceUpdates)
	})

	s.Run("correctly can create a new price feed with no vote", func() {
		valUpdates.Updates[cp1] = slatypes.NoVote
		priceFeedUpdates.ValidatorUpdates[consAddress1.String()] = valUpdates

		err := s.keeper.UpdatePriceFeedsForSLA(s.ctx, sla, priceFeedUpdates)
		s.Require().NoError(err)

		// Check that the price feed was added.
		priceFeeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id)
		s.Require().NoError(err)
		s.Require().Len(priceFeeds, 1)

		feed := priceFeeds[0]
		s.Require().Equal(consAddress1, sdk.ConsAddress(feed.Validator))
		s.Require().Equal(cp1, feed.CurrencyPair)
		s.Require().Equal(sla.ID, feed.ID)
		s.Require().Equal(sla.MaximumViableWindow, feed.MaximumViableWindow)

		// Check that the status was correctly set.
		numVotes, err := feed.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(0), numVotes)

		numPriceUpdates, err := feed.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(0), numPriceUpdates)
	})

	s.Run("correctly can update an existing price feed with price update", func() {
		feed, err := slatypes.NewPriceFeed(uint(sla.MaximumViableWindow), consAddress1, cp1, sla.ID)
		s.Require().NoError(err)

		err = feed.SetUpdate(slatypes.VoteWithoutPrice)
		s.Require().NoError(err)

		err = feed.SetUpdate(slatypes.VoteWithPrice)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, feed)
		s.Require().NoError(err)

		valUpdates.Updates[cp1] = slatypes.VoteWithPrice
		priceFeedUpdates.ValidatorUpdates[consAddress1.String()] = valUpdates

		err = s.keeper.UpdatePriceFeedsForSLA(s.ctx, sla, priceFeedUpdates)
		s.Require().NoError(err)

		// Check that the price feed was updated.
		priceFeeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id)
		s.Require().NoError(err)
		s.Require().Len(priceFeeds, 1)

		feed = priceFeeds[0]
		s.Require().Equal(consAddress1, sdk.ConsAddress(feed.Validator))
		s.Require().Equal(cp1, feed.CurrencyPair)
		s.Require().Equal(sla.ID, feed.ID)
		s.Require().Equal(sla.MaximumViableWindow, feed.MaximumViableWindow)

		// Check that the status was correctly set.
		numVotes, err := feed.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numVotes)

		numPriceUpdates, err := feed.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numPriceUpdates)

		numVotes, err = feed.GetNumVotesWithWindow(2)
		s.Require().NoError(err)
		s.Require().Equal(uint(2), numVotes)

		numPriceUpdates, err = feed.GetNumPriceUpdatesWithWindow(2)
		s.Require().NoError(err)
		s.Require().Equal(uint(2), numPriceUpdates)

		numVotes, err = feed.GetNumVotesWithWindow(3)
		s.Require().NoError(err)
		s.Require().Equal(uint(3), numVotes)

		numPriceUpdates, err = feed.GetNumPriceUpdatesWithWindow(3)
		s.Require().NoError(err)
		s.Require().Equal(uint(2), numPriceUpdates)
	})

	s.Run("correctly can update an existing price feed with only a vote and no price update", func() {
		feed, err := slatypes.NewPriceFeed(uint(sla.MaximumViableWindow), consAddress1, cp1, sla.ID)
		s.Require().NoError(err)

		err = feed.SetUpdate(slatypes.VoteWithoutPrice)
		s.Require().NoError(err)

		err = feed.SetUpdate(slatypes.VoteWithPrice)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, feed)
		s.Require().NoError(err)

		valUpdates.Updates[cp1] = slatypes.VoteWithoutPrice
		priceFeedUpdates.ValidatorUpdates[consAddress1.String()] = valUpdates

		err = s.keeper.UpdatePriceFeedsForSLA(s.ctx, sla, priceFeedUpdates)
		s.Require().NoError(err)

		// Check that the price feed was updated.
		priceFeeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id)
		s.Require().NoError(err)
		s.Require().Len(priceFeeds, 1)

		feed = priceFeeds[0]
		s.Require().Equal(consAddress1, sdk.ConsAddress(feed.Validator))
		s.Require().Equal(cp1, feed.CurrencyPair)
		s.Require().Equal(sla.ID, feed.ID)
		s.Require().Equal(sla.MaximumViableWindow, feed.MaximumViableWindow)

		// Check that the status was correctly set.
		numVotes, err := feed.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numVotes)

		numPriceUpdates, err := feed.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(0), numPriceUpdates)

		numVotes, err = feed.GetNumVotesWithWindow(2)
		s.Require().NoError(err)
		s.Require().Equal(uint(2), numVotes)

		numPriceUpdates, err = feed.GetNumPriceUpdatesWithWindow(2)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numPriceUpdates)

		numVotes, err = feed.GetNumVotesWithWindow(3)
		s.Require().NoError(err)
		s.Require().Equal(uint(3), numVotes)

		numPriceUpdates, err = feed.GetNumPriceUpdatesWithWindow(3)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numPriceUpdates)
	})

	s.Run("correctly can update an existing price feed with no vote", func() {
		feed, err := slatypes.NewPriceFeed(uint(sla.MaximumViableWindow), consAddress1, cp1, sla.ID)
		s.Require().NoError(err)

		err = feed.SetUpdate(slatypes.VoteWithoutPrice)
		s.Require().NoError(err)

		err = feed.SetUpdate(slatypes.VoteWithPrice)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, feed)
		s.Require().NoError(err)

		valUpdates.Updates[cp1] = slatypes.NoVote
		priceFeedUpdates.ValidatorUpdates[consAddress1.String()] = valUpdates

		err = s.keeper.UpdatePriceFeedsForSLA(s.ctx, sla, priceFeedUpdates)
		s.Require().NoError(err)

		// Check that the price feed was updated.
		priceFeeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id)
		s.Require().NoError(err)
		s.Require().Len(priceFeeds, 1)

		feed = priceFeeds[0]
		s.Require().Equal(consAddress1, sdk.ConsAddress(feed.Validator))
		s.Require().Equal(cp1, feed.CurrencyPair)
		s.Require().Equal(sla.ID, feed.ID)
		s.Require().Equal(sla.MaximumViableWindow, feed.MaximumViableWindow)

		// Check that the status was correctly set.
		numVotes, err := feed.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(0), numVotes)

		numPriceUpdates, err := feed.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(0), numPriceUpdates)

		numVotes, err = feed.GetNumVotesWithWindow(2)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numVotes)

		numPriceUpdates, err = feed.GetNumPriceUpdatesWithWindow(2)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numPriceUpdates)

		numVotes, err = feed.GetNumVotesWithWindow(3)
		s.Require().NoError(err)
		s.Require().Equal(uint(2), numVotes)

		numPriceUpdates, err = feed.GetNumPriceUpdatesWithWindow(3)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numPriceUpdates)
	})

	s.Run("can correctly create a new price feed with different currency pairs", func() {
		priceFeedUpdates.CurrencyPairs[cp2] = struct{}{}

		valUpdates.Updates[cp1] = slatypes.VoteWithoutPrice
		valUpdates.Updates[cp2] = slatypes.VoteWithPrice
		priceFeedUpdates.ValidatorUpdates[consAddress1.String()] = valUpdates

		err := s.keeper.UpdatePriceFeedsForSLA(s.ctx, sla, priceFeedUpdates)
		s.Require().NoError(err)

		// Check that the price feed was added.
		priceFeeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id)
		s.Require().NoError(err)
		s.Require().Len(priceFeeds, 2)

		feedCP1 := priceFeeds[0]
		feedCP2 := priceFeeds[1]
		if feedCP1.CurrencyPair.String() != cp1.String() {
			feedCP1 = priceFeeds[1]
			feedCP2 = priceFeeds[0]
		}

		// Check that the statuses were correctly set.
		numVotes, err := feedCP1.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numVotes)

		numPriceUpdates, err := feedCP1.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(0), numPriceUpdates)

		numVotes, err = feedCP2.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numVotes)

		numPriceUpdates, err = feedCP2.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numPriceUpdates)
	})

	s.Run("can correctly create a new price feed and update an existing one", func() {
		priceFeedUpdates := slakeeper.NewPriceFeedUpdates()
		priceFeedUpdates.CurrencyPairs[cp1] = struct{}{}

		validatorUpdates := slakeeper.NewValidatorUpdate(consAddress1)
		validatorUpdates.Updates[cp1] = slatypes.VoteWithPrice

		feed, err := slatypes.NewPriceFeed(uint(sla.MaximumViableWindow), consAddress1, cp1, sla.ID)
		s.Require().NoError(err)

		err = feed.SetUpdate(slatypes.VoteWithPrice)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, feed)
		s.Require().NoError(err)

		priceFeedUpdates = slakeeper.NewPriceFeedUpdates()
		priceFeedUpdates.CurrencyPairs[cp1] = struct{}{}

		validatorUpdates = slakeeper.NewValidatorUpdate(consAddress1)
		validatorUpdates.Updates[cp1] = slatypes.VoteWithoutPrice
		priceFeedUpdates.ValidatorUpdates[consAddress1.String()] = validatorUpdates

		validatorUpdates2 := slakeeper.NewValidatorUpdate(consAddress2)
		validatorUpdates2.Updates[cp1] = slatypes.VoteWithPrice
		priceFeedUpdates.ValidatorUpdates[consAddress2.String()] = validatorUpdates2

		err = s.keeper.UpdatePriceFeedsForSLA(s.ctx, sla, priceFeedUpdates)
		s.Require().NoError(err)

		// Check that the price feed was updated.
		priceFeeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id)
		s.Require().NoError(err)
		s.Require().Len(priceFeeds, 2)

		val1Feed := priceFeeds[0]
		val2Feed := priceFeeds[1]
		if !sdk.ConsAddress(feed.Validator).Equals(sdk.ConsAddress(val1Feed.Validator)) {
			val1Feed = priceFeeds[1]
			val2Feed = priceFeeds[0]
		}

		// Check that the statuses were correctly set.
		numVotes, err := val1Feed.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numVotes)

		numPriceUpdates, err := val1Feed.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(0), numPriceUpdates)

		numVotes, err = val2Feed.GetNumVotesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numVotes)

		numPriceUpdates, err = val2Feed.GetNumPriceUpdatesWithWindow(1)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numPriceUpdates)

		numVotes, err = val1Feed.GetNumVotesWithWindow(2)
		s.Require().NoError(err)
		s.Require().Equal(uint(2), numVotes)

		numPriceUpdates, err = val1Feed.GetNumPriceUpdatesWithWindow(2)
		s.Require().NoError(err)
		s.Require().Equal(uint(1), numPriceUpdates)
	})
}
