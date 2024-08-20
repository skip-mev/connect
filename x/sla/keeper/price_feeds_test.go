package keeper_test

import (
	"testing"

	"github.com/bits-and-blooms/bitset"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

const (
	id1 = "testId"
)

func (s *KeeperTestSuite) TestSetPriceFeed() {
	cp1 := slinkytypes.NewCurrencyPair("btc", "usd")

	consAddress1 := sdk.ConsAddress("consAddress1")
	consAddress2 := sdk.ConsAddress("consAddress2")

	priceFeed1, err := slatypes.NewPriceFeed(
		10,
		consAddress1,
		cp1,
		id1,
	)
	s.Require().NoError(err)
	priceFeed2, _ := slatypes.NewPriceFeed(
		10,
		consAddress2,
		cp1,
		id1,
	)
	s.Require().NoError(err)

	s.Run("returns error when feed does not exist", func() {
		_, err := s.keeper.GetPriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().Error(err)

		contains, err := s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)
		s.Require().False(contains)
	})

	s.Run("can set and get a single feed", func() {
		err := s.keeper.SetPriceFeed(s.ctx, priceFeed1)
		s.Require().NoError(err)

		feed, err := s.keeper.GetPriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)
		checkEquality(s.T(), priceFeed1, feed)

		contains, err := s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)
		s.Require().True(contains)
	})

	s.Run("can set and get multiple feeds", func() {
		err := s.keeper.SetPriceFeed(s.ctx, priceFeed1)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, priceFeed2)
		s.Require().NoError(err)

		feed1, err := s.keeper.GetPriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)
		checkEquality(s.T(), priceFeed1, feed1)

		feed2, err := s.keeper.GetPriceFeed(s.ctx, id1, cp1, consAddress2)
		s.Require().NoError(err)
		checkEquality(s.T(), priceFeed2, feed2)

		contains, err := s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)
		s.Require().True(contains)

		contains, err = s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress2)
		s.Require().NoError(err)
		s.Require().True(contains)

		feed, err := s.keeper.GetAllPriceFeeds(s.ctx, id1)
		s.Require().NoError(err)
		s.Require().Len(feed, 2)
		checkEquality(s.T(), priceFeed1, feed[0])
		checkEquality(s.T(), priceFeed2, feed[1])
	})
}

func (s *KeeperTestSuite) TestRemovePriceFeeds() {
	cp1 := slinkytypes.NewCurrencyPair("btc", "usd")

	consAddress1 := sdk.ConsAddress("consAddress1")
	consAddress2 := sdk.ConsAddress("consAddress2")

	priceFeed1, err := slatypes.NewPriceFeed(
		10,
		consAddress1,
		cp1,
		id1,
	)
	s.Require().NoError(err)
	priceFeed2, _ := slatypes.NewPriceFeed(
		10,
		consAddress2,
		cp1,
		id1,
	)
	s.Require().NoError(err)

	s.Run("no error when removing a feed that does not exist", func() {
		contains, err := s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)
		s.Require().False(contains)

		err = s.keeper.RemovePriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)

		contains, err = s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)
		s.Require().False(contains)
	})

	s.Run("can remove a feed", func() {
		err := s.keeper.SetPriceFeed(s.ctx, priceFeed1)
		s.Require().NoError(err)

		contains, err := s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)
		s.Require().True(contains)

		err = s.keeper.RemovePriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)

		contains, err = s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)
		s.Require().False(contains)

		feeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id1)
		s.Require().NoError(err)
		s.Require().Empty(feeds)
	})

	s.Run("can remove multiple feeds", func() {
		err := s.keeper.SetPriceFeed(s.ctx, priceFeed1)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, priceFeed2)
		s.Require().NoError(err)

		contains, err := s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)
		s.Require().True(contains)

		contains, err = s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress2)
		s.Require().NoError(err)
		s.Require().True(contains)

		err = s.keeper.RemovePriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)

		contains, err = s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)
		s.Require().False(contains)

		contains, err = s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress2)
		s.Require().NoError(err)
		s.Require().True(contains)

		feeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id1)
		s.Require().NoError(err)
		s.Require().Len(feeds, 1)
		checkEquality(s.T(), priceFeed2, feeds[0])

		err = s.keeper.RemovePriceFeed(s.ctx, id1, cp1, consAddress2)
		s.Require().NoError(err)

		contains, err = s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)
		s.Require().False(contains)

		contains, err = s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress2)
		s.Require().NoError(err)
		s.Require().False(contains)
	})

	cp2 := slinkytypes.NewCurrencyPair("mog", "usd")
	priceFeed3, err := slatypes.NewPriceFeed(
		10,
		consAddress1,
		cp2, // different currency pair
		id1,
	)
	s.Require().NoError(err)

	s.Run("can remove all feeds for a given currency pair", func() {
		err = s.keeper.SetPriceFeed(s.ctx, priceFeed1)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, priceFeed2)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, priceFeed3)
		s.Require().NoError(err)

		contains, err := s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)
		s.Require().True(contains)

		contains, err = s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress2)
		s.Require().NoError(err)
		s.Require().True(contains)

		contains, err = s.keeper.ContainsPriceFeed(s.ctx, id1, cp2, consAddress1)
		s.Require().NoError(err)
		s.Require().True(contains)

		err = s.keeper.RemovePriceFeedByCurrencyPair(s.ctx, id1, cp1)
		s.Require().NoError(err)

		contains, err = s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress1)
		s.Require().NoError(err)
		s.Require().False(contains)

		contains, err = s.keeper.ContainsPriceFeed(s.ctx, id1, cp1, consAddress2)
		s.Require().NoError(err)
		s.Require().False(contains)

		contains, err = s.keeper.ContainsPriceFeed(s.ctx, id1, cp2, consAddress1)
		s.Require().NoError(err)
		s.Require().True(contains)

		feeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id1)
		s.Require().NoError(err)
		s.Require().Len(feeds, 1)
		checkEquality(s.T(), priceFeed3, feeds[0])
	})

	cp3 := slinkytypes.NewCurrencyPair("gds", "4L")

	s.Run("different currency pairs are not affected", func() {
		priceFeed1 := priceFeed1
		priceFeed2 := priceFeed2
		priceFeed3 := priceFeed3

		priceFeed1.CurrencyPair = cp1
		priceFeed2.CurrencyPair = cp2
		priceFeed3.CurrencyPair = cp3

		err = s.keeper.SetPriceFeed(s.ctx, priceFeed1)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, priceFeed2)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, priceFeed3)
		s.Require().NoError(err)

		feeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id1)
		s.Require().NoError(err)
		s.Require().Len(feeds, 3)
		checkEquality(s.T(), priceFeed1, feeds[0])
		checkEquality(s.T(), priceFeed2, feeds[1])
		checkEquality(s.T(), priceFeed3, feeds[2])

		err = s.keeper.RemovePriceFeedByCurrencyPair(s.ctx, id1, cp1)
		s.Require().NoError(err)

		feeds, err = s.keeper.GetAllPriceFeeds(s.ctx, id1)
		s.Require().NoError(err)
		s.Require().Len(feeds, 2)
		checkEquality(s.T(), priceFeed2, feeds[0])
		checkEquality(s.T(), priceFeed3, feeds[1])

		err = s.keeper.RemovePriceFeedByCurrencyPair(s.ctx, id1, cp2)
		s.Require().NoError(err)

		feeds, err = s.keeper.GetAllPriceFeeds(s.ctx, id1)
		s.Require().NoError(err)
		s.Require().Len(feeds, 1)

		err = s.keeper.RemovePriceFeedByCurrencyPair(s.ctx, id1, cp3)
		s.Require().NoError(err)

		feeds, err = s.keeper.GetAllPriceFeeds(s.ctx, id1)
		s.Require().NoError(err)
		s.Require().Empty(feeds)
	})

	s.Run("can remove all feeds for a given sla", func() {
		err := s.keeper.SetPriceFeed(s.ctx, priceFeed1)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, priceFeed2)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, priceFeed3)
		s.Require().NoError(err)

		feeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id1)
		s.Require().NoError(err)
		s.Require().Len(feeds, 3)

		err = s.keeper.RemovePriceFeedsBySLA(s.ctx, id1)
		s.Require().NoError(err)

		feeds, err = s.keeper.GetAllPriceFeeds(s.ctx, id1)
		s.Require().NoError(err)
		s.Require().Empty(feeds)
	})

	s.Run("can remove all price feeds for a given sla with other sla ids present", func() {
		err := s.keeper.SetPriceFeed(s.ctx, priceFeed1)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, priceFeed2)
		s.Require().NoError(err)

		priceFeed3.ID = "testId2"
		err = s.keeper.SetPriceFeed(s.ctx, priceFeed3)
		s.Require().NoError(err)

		feeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id1)
		s.Require().NoError(err)
		s.Require().Len(feeds, 2)

		feeds, err = s.keeper.GetAllPriceFeeds(s.ctx, "testId2")
		s.Require().NoError(err)
		s.Require().Len(feeds, 1)

		err = s.keeper.RemovePriceFeedsBySLA(s.ctx, id1)
		s.Require().NoError(err)

		feeds, err = s.keeper.GetAllPriceFeeds(s.ctx, id1)
		s.Require().NoError(err)
		s.Require().Empty(feeds)

		feeds, err = s.keeper.GetAllPriceFeeds(s.ctx, "testId2")
		s.Require().NoError(err)
		s.Require().Len(feeds, 1)
	})
}

func checkEquality(t *testing.T, sla1, sla2 slatypes.PriceFeed) {
	t.Helper()

	require.Equal(t, sla1.ID, sla2.ID)
	require.Equal(t, sla1.MaximumViableWindow, sla2.MaximumViableWindow)
	require.Equal(t, sla1.Index, sla2.Index)

	updateMap1 := bitset.New(uint(sla1.MaximumViableWindow))
	require.NoError(t, updateMap1.UnmarshalBinary(sla1.UpdateMap))

	updateMap2 := bitset.New(uint(sla2.MaximumViableWindow))
	require.NoError(t, updateMap2.UnmarshalBinary(sla2.UpdateMap))

	inclusionMap1 := bitset.New(uint(sla1.MaximumViableWindow))
	require.NoError(t, inclusionMap1.UnmarshalBinary(sla1.InclusionMap))

	inclusionMap2 := bitset.New(uint(sla2.MaximumViableWindow))
	require.NoError(t, inclusionMap2.UnmarshalBinary(sla2.InclusionMap))

	require.True(t, updateMap1.Equal(updateMap2))
	require.True(t, inclusionMap1.Equal(inclusionMap2))
}
