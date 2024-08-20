package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

var (
	val = sdk.ConsAddress("validator1")
	cp  = slinkytypes.NewCurrencyPair("BTC", "ETH")
	id  = "testID"
)

func TestSetUpdate(t *testing.T) {
	t.Run("vote with price", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, val, cp, id)
		require.NoError(t, err)
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))

		count, err := priceFeed.GetUpdateCount()
		require.NoError(t, err)
		require.Equal(t, uint(1), count)

		count, err = priceFeed.GetInclusionCount()
		require.NoError(t, err)
		require.Equal(t, uint(1), count)

		bit, err := priceFeed.GetInclusionBit(0)
		require.NoError(t, err)
		require.Equal(t, true, bit)

		bit, err = priceFeed.GetUpdateBit(0)
		require.NoError(t, err)
		require.Equal(t, true, bit)

		require.Equal(t, uint(1), uint(priceFeed.Index))
	})

	t.Run("vote without price", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, val, cp, id)
		require.NoError(t, err)
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))

		count, err := priceFeed.GetUpdateCount()
		require.NoError(t, err)
		require.Equal(t, uint(0), count)

		count, err = priceFeed.GetInclusionCount()
		require.NoError(t, err)
		require.Equal(t, uint(1), count)

		bit, err := priceFeed.GetInclusionBit(0)
		require.NoError(t, err)
		require.Equal(t, true, bit)

		bit, err = priceFeed.GetUpdateBit(0)
		require.NoError(t, err)
		require.Equal(t, false, bit)

		require.Equal(t, uint(1), uint(priceFeed.Index))
	})

	t.Run("multiple votes", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, val, cp, id)
		require.NoError(t, err)
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))

		count, err := priceFeed.GetUpdateCount()
		require.NoError(t, err)
		require.Equal(t, uint(2), count)

		count, err = priceFeed.GetInclusionCount()
		require.NoError(t, err)
		require.Equal(t, uint(4), count)

		bit, err := priceFeed.GetInclusionBit(0)
		require.NoError(t, err)
		require.Equal(t, true, bit)

		bit, err = priceFeed.GetInclusionBit(1)
		require.NoError(t, err)
		require.Equal(t, true, bit)

		bit, err = priceFeed.GetInclusionBit(2)
		require.NoError(t, err)
		require.Equal(t, true, bit)

		bit, err = priceFeed.GetInclusionBit(3)
		require.NoError(t, err)
		require.Equal(t, true, bit)

		bit, err = priceFeed.GetUpdateBit(0)
		require.NoError(t, err)
		require.Equal(t, true, bit)

		bit, err = priceFeed.GetUpdateBit(1)
		require.NoError(t, err)
		require.Equal(t, false, bit)

		bit, err = priceFeed.GetUpdateBit(2)
		require.NoError(t, err)
		require.Equal(t, true, bit)

		bit, err = priceFeed.GetUpdateBit(3)
		require.NoError(t, err)
		require.Equal(t, false, bit)

		require.Equal(t, uint(4), uint(priceFeed.Index))
	})

	t.Run("multiple votes with wraparound", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(4, val, cp, id)
		require.NoError(t, err)

		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))

		count, err := priceFeed.GetUpdateCount()
		require.NoError(t, err)
		require.Equal(t, uint(2), count)

		count, err = priceFeed.GetInclusionCount()
		require.NoError(t, err)
		require.Equal(t, uint(4), count)

		require.Equal(t, uint(0), uint(priceFeed.Index))
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))

		count, err = priceFeed.GetUpdateCount()
		require.NoError(t, err)
		require.Equal(t, uint(1), count)
	})

	t.Run("no vote", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, val, cp, id)
		require.NoError(t, err)
		require.NoError(t, priceFeed.SetUpdate(slatypes.NoVote))

		count, err := priceFeed.GetUpdateCount()
		require.NoError(t, err)
		require.Equal(t, uint(0), count)

		count, err = priceFeed.GetInclusionCount()
		require.NoError(t, err)
		require.Equal(t, uint(0), count)

		require.Equal(t, uint(1), uint(priceFeed.Index))
	})
}

func TestGetNumberOfPriceUpdates(t *testing.T) {
	t.Run("correctly set all bits and returns correct number of updates", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(4, val, cp, id)
		require.NoError(t, err)

		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))

		numUpdates, err := priceFeed.GetNumPriceUpdatesWithWindow(1)
		require.NoError(t, err)
		require.Equal(t, uint(1), numUpdates)

		numUpdates, err = priceFeed.GetNumPriceUpdatesWithWindow(2)
		require.NoError(t, err)
		require.Equal(t, uint(1), numUpdates)

		numUpdates, err = priceFeed.GetNumPriceUpdatesWithWindow(3)
		require.NoError(t, err)
		require.Equal(t, uint(2), numUpdates)

		numUpdates, err = priceFeed.GetNumPriceUpdatesWithWindow(4)
		require.NoError(t, err)
		require.Equal(t, uint(2), numUpdates)

		_, err = priceFeed.GetNumPriceUpdatesWithWindow(5)
		require.Error(t, err)
	})

	t.Run("correctly set all bits and returns correct number of updates with wrap around", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(4, val, cp, id)
		require.NoError(t, err)

		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))

		numUpdates, err := priceFeed.GetNumPriceUpdatesWithWindow(2)
		require.NoError(t, err)
		require.Equal(t, uint(1), numUpdates)

		numUpdates, err = priceFeed.GetNumPriceUpdatesWithWindow(4)
		require.NoError(t, err)
		require.Equal(t, uint(2), numUpdates)

		require.NoError(t, priceFeed.SetUpdate(slatypes.NoVote))

		numUpdates, err = priceFeed.GetNumPriceUpdatesWithWindow(4)
		require.NoError(t, err)
		require.Equal(t, uint(1), numUpdates)

		require.NoError(t, priceFeed.SetUpdate(slatypes.NoVote))

		numUpdates, err = priceFeed.GetNumPriceUpdatesWithWindow(4)
		require.NoError(t, err)
		require.Equal(t, uint(1), numUpdates)
	})
}

func TestGetNumVotesWithWindow(t *testing.T) {
	t.Run("correctly returns no votes with no sets", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(4, val, cp, id)
		require.NoError(t, err)

		numVotes, err := priceFeed.GetNumVotesWithWindow(1)
		require.NoError(t, err)
		require.Equal(t, uint(0), numVotes)

		numVotes, err = priceFeed.GetNumVotesWithWindow(2)
		require.NoError(t, err)
		require.Equal(t, uint(0), numVotes)

		numVotes, err = priceFeed.GetNumVotesWithWindow(3)
		require.NoError(t, err)
		require.Equal(t, uint(0), numVotes)

		numVotes, err = priceFeed.GetNumVotesWithWindow(4)
		require.NoError(t, err)
		require.Equal(t, uint(0), numVotes)

		_, err = priceFeed.GetNumVotesWithWindow(5)
		require.Error(t, err)
	})

	t.Run("correctly returns with one set", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(4, val, cp, id)
		require.NoError(t, err)

		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))

		numVotes, err := priceFeed.GetNumVotesWithWindow(1)
		require.NoError(t, err)
		require.Equal(t, uint(1), numVotes)

		numVotes, err = priceFeed.GetNumVotesWithWindow(2)
		require.NoError(t, err)
		require.Equal(t, uint(1), numVotes)

		numVotes, err = priceFeed.GetNumVotesWithWindow(3)
		require.NoError(t, err)
		require.Equal(t, uint(1), numVotes)

		numVotes, err = priceFeed.GetNumVotesWithWindow(4)
		require.NoError(t, err)
		require.Equal(t, uint(1), numVotes)

		_, err = priceFeed.GetNumVotesWithWindow(5)
		require.Error(t, err)
	})

	t.Run("correctly returns with two sets", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(4, val, cp, id)
		require.NoError(t, err)

		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))

		numVotes, err := priceFeed.GetNumVotesWithWindow(1)
		require.NoError(t, err)
		require.Equal(t, uint(1), numVotes)

		numVotes, err = priceFeed.GetNumVotesWithWindow(2)
		require.NoError(t, err)
		require.Equal(t, uint(2), numVotes)

		numVotes, err = priceFeed.GetNumVotesWithWindow(3)
		require.NoError(t, err)
		require.Equal(t, uint(2), numVotes)

		numVotes, err = priceFeed.GetNumVotesWithWindow(4)
		require.NoError(t, err)
		require.Equal(t, uint(2), numVotes)

		_, err = priceFeed.GetNumVotesWithWindow(5)
		require.Error(t, err)
	})

	t.Run("correctly returns after some sets with wrap around", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(4, val, cp, id)
		require.NoError(t, err)

		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))

		numVotes, err := priceFeed.GetNumVotesWithWindow(1)
		require.NoError(t, err)
		require.Equal(t, uint(1), numVotes)

		numVotes, err = priceFeed.GetNumVotesWithWindow(2)
		require.NoError(t, err)
		require.Equal(t, uint(2), numVotes)

		require.NoError(t, priceFeed.SetUpdate(slatypes.NoVote))

		numVotes, err = priceFeed.GetNumVotesWithWindow(2)
		require.NoError(t, err)
		require.Equal(t, uint(1), numVotes)

		numVotes, err = priceFeed.GetNumVotesWithWindow(1)
		require.NoError(t, err)
		require.Equal(t, uint(0), numVotes)

		require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))

		numVotes, err = priceFeed.GetNumVotesWithWindow(2)
		require.NoError(t, err)
		require.Equal(t, uint(1), numVotes)

		numVotes, err = priceFeed.GetNumVotesWithWindow(3)
		require.NoError(t, err)
		require.Equal(t, uint(2), numVotes)

		numVotes, err = priceFeed.GetNumVotesWithWindow(4)
		require.NoError(t, err)
		require.Equal(t, uint(3), numVotes)

		// Now we have a wrap around
		require.NoError(t, priceFeed.SetUpdate(slatypes.NoVote))

		numVotes, err = priceFeed.GetNumVotesWithWindow(4)
		require.NoError(t, err)
		require.Equal(t, uint(2), numVotes)
	})
}

func TestPriceFeedValidateBasic(t *testing.T) {
	t.Run("valid price feed", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, val, cp, id)
		require.NoError(t, err)
		require.NoError(t, priceFeed.ValidateBasic())
	})

	t.Run("missing a price feed id", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, val, cp, "")
		require.NoError(t, err)
		require.Error(t, priceFeed.ValidateBasic())
	})

	t.Run("invalid max window", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(0, val, cp, id)
		require.NoError(t, err)
		require.Error(t, priceFeed.ValidateBasic())
	})

	t.Run("invalid inclusion map", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, val, cp, id)
		require.NoError(t, err)
		priceFeed.InclusionMap = []byte("invalid")
		require.Error(t, priceFeed.ValidateBasic())
	})

	t.Run("invalid update map", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, val, cp, id)
		require.NoError(t, err)
		priceFeed.UpdateMap = []byte("invalid")
		require.Error(t, priceFeed.ValidateBasic())
	})

	t.Run("invalid validator address", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, nil, cp, id)
		require.NoError(t, err)
		require.Error(t, priceFeed.ValidateBasic())
	})

	t.Run("invalid currency pair", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, val, slinkytypes.CurrencyPair{}, id)
		require.NoError(t, err)
		require.Error(t, priceFeed.ValidateBasic())
	})
}
