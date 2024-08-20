package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

func TestSLAValidateBasic(t *testing.T) {
	t.Run("empty id should be rejected", func(t *testing.T) {
		sla := slatypes.NewPriceFeedSLA(
			"",
			1,
			math.LegacyMustNewDecFromStr("0.5"),
			math.LegacyMustNewDecFromStr("0.5"),
			1,
			1,
		)
		err := sla.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("maximum viable window of 0 should be rejected", func(t *testing.T) {
		sla := slatypes.NewPriceFeedSLA(
			"test",
			0,
			math.LegacyMustNewDecFromStr("0.5"),
			math.LegacyMustNewDecFromStr("0.5"),
			1,
			1,
		)
		err := sla.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("equal viable windows should be rejected", func(t *testing.T) {
		sla := slatypes.NewPriceFeedSLA(
			"test",
			1,
			math.LegacyMustNewDecFromStr("0.5"),
			math.LegacyMustNewDecFromStr("0.5"),
			1,
			1,
		)
		err := sla.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("expected uptime of 0 or less should be rejected", func(t *testing.T) {
		sla := slatypes.NewPriceFeedSLA(
			"test",
			1,
			math.LegacyMustNewDecFromStr("0"),
			math.LegacyMustNewDecFromStr("0.5"),
			1,
			1,
		)
		err := sla.ValidateBasic()
		require.Error(t, err)

		sla = slatypes.NewPriceFeedSLA(
			"test",
			1,
			math.LegacyMustNewDecFromStr("-1"),
			math.LegacyMustNewDecFromStr("0.5"),
			1,
			1,
		)
		err = sla.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("slashing constant of 0 or less should be rejected", func(t *testing.T) {
		sla := slatypes.NewPriceFeedSLA(
			"test",
			1,
			math.LegacyMustNewDecFromStr("0.5"),
			math.LegacyMustNewDecFromStr("0"),
			1,
			1,
		)
		err := sla.ValidateBasic()
		require.Error(t, err)

		sla = slatypes.NewPriceFeedSLA(
			"test",
			1,
			math.LegacyMustNewDecFromStr("0.5"),
			math.LegacyMustNewDecFromStr("-1"),
			1,
			1,
		)
		err = sla.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("minimum block updates of 0 should be rejected", func(t *testing.T) {
		sla := slatypes.NewPriceFeedSLA(
			"test",
			1,
			math.LegacyMustNewDecFromStr("0.5"),
			math.LegacyMustNewDecFromStr("0.5"),
			0,
			1,
		)
		err := sla.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("minimum block updates greater than maximum viable window should be rejected", func(t *testing.T) {
		sla := slatypes.NewPriceFeedSLA(
			"test",
			1,
			math.LegacyMustNewDecFromStr("0.5"),
			math.LegacyMustNewDecFromStr("0.5"),
			2,
			1,
		)
		err := sla.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("frequency of 0 should be rejected", func(t *testing.T) {
		sla := slatypes.NewPriceFeedSLA(
			"test",
			2,
			math.LegacyMustNewDecFromStr("0.5"),
			math.LegacyMustNewDecFromStr("0.5"),
			1,
			0,
		)
		err := sla.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("frequency that is greater than the maximum viable window must be rejected", func(t *testing.T) {
		sla := slatypes.NewPriceFeedSLA(
			"test",
			2,
			math.LegacyMustNewDecFromStr("0.5"),
			math.LegacyMustNewDecFromStr("0.5"),
			1,
			3,
		)
		err := sla.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("valid sla should be accepted", func(t *testing.T) {
		sla := slatypes.NewPriceFeedSLA(
			"test",
			10,
			math.LegacyMustNewDecFromStr("0.5"),
			math.LegacyMustNewDecFromStr("0.5"),
			5,
			5,
		)
		err := sla.ValidateBasic()
		require.NoError(t, err)
	})
}

func TestQualifies(t *testing.T) {
	sla := slatypes.NewPriceFeedSLA(
		"test",
		10,
		math.LegacyMustNewDecFromStr("0.5"),
		math.LegacyMustNewDecFromStr("0.5"),
		5,
		5,
	)

	t.Run("does not qualify when the SLA ID is different from price feed ID", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, val, cp, "differentID")
		require.NoError(t, err)

		qualifies, err := sla.Qualifies(priceFeed)
		require.NoError(t, err)
		require.False(t, qualifies)
	})

	t.Run("errors when the IDs are the same but the maximum viable window is different", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(20, val, cp, "test")
		require.NoError(t, err)

		qualifies, err := sla.Qualifies(priceFeed)
		require.Error(t, err)
		require.False(t, qualifies)
	})

	t.Run("does not qualify when price feed has not seen enough votes", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, val, cp, "test")
		require.NoError(t, err)

		qualifies, err := sla.Qualifies(priceFeed)
		require.NoError(t, err)
		require.False(t, qualifies)
	})

	t.Run("qualifies when price feed has seen enough votes", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, val, cp, "test")
		require.NoError(t, err)

		// Vote with price for 5 blocks
		for i := 0; i < 5; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		}

		qualifies, err := sla.Qualifies(priceFeed)
		require.NoError(t, err)
		require.True(t, qualifies)
	})

	t.Run("qualifies when price feed has seen enough votes with wraparound", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, val, cp, "test")
		require.NoError(t, err)

		// Vote with price for 5 blocks
		for i := 0; i < 5; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		}

		// Vote with price for 5 blocks
		for i := 0; i < 5; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.NoVote))
		}

		qualifies, err := sla.Qualifies(priceFeed)
		require.NoError(t, err)
		require.True(t, qualifies)
	})

	t.Run("does not qualify when price feed has not seen enough votes with wraparound", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, val, cp, "test")
		require.NoError(t, err)

		// Vote with price for 5 blocks
		for i := 0; i < 5; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		}

		// Vote with price for 5 blocks
		for i := 0; i < 6; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.NoVote))
		}

		qualifies, err := sla.Qualifies(priceFeed)
		require.NoError(t, err)
		require.False(t, qualifies)
	})

	t.Run("qualifies when there is a large wrap around", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(10, val, cp, "test")
		require.NoError(t, err)

		// Vote with price for 5 blocks
		for i := 0; i < 5; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		}

		// Vote with price for 5 blocks
		for i := 0; i < 15; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.NoVote))
		}

		// Vote with price for 5 blocks
		for i := 0; i < 5; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		}

		qualifies, err := sla.Qualifies(priceFeed)
		require.NoError(t, err)
		require.True(t, qualifies)
	})
}

func TestGetUptimeFromPriceFeed(t *testing.T) {
	// Strategy inputs
	id := "testID"
	maximumViableWindow := uint(20)
	expectedUptime := math.LegacyMustNewDecFromStr("0.8")
	k := math.LegacyMustNewDecFromStr("1.0")
	frequency := 10
	minimumBlockUpdates := uint(10)

	sla := slatypes.NewPriceFeedSLA(
		id,
		uint64(maximumViableWindow),
		expectedUptime,
		k,
		uint64(minimumBlockUpdates),
		uint64(frequency),
	)

	t.Run("returns valid uptime", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(maximumViableWindow, val, cp, id)
		require.NoError(t, err)

		// Vote with price for 5 blocks
		for i := 0; i < 5; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		}

		// Vote without price for 5 blocks
		for i := 0; i < 5; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))
		}

		uptime, err := sla.GetUptimeFromPriceFeed(priceFeed)
		require.NoError(t, err)
		require.Equal(t, math.LegacyMustNewDecFromStr("0.5"), uptime)
	})

	t.Run("returns valid uptime with no votes in between", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(maximumViableWindow, val, cp, id)
		require.NoError(t, err)

		// Vote with price for 5 blocks
		for i := 0; i < 5; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		}

		// Vote without price for 15 blocks
		for i := 0; i < 10; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.NoVote))
		}

		// Vote without price for 5 blocks
		for i := 0; i < 5; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))
		}

		uptime, err := sla.GetUptimeFromPriceFeed(priceFeed)
		require.NoError(t, err)
		require.Equal(t, math.LegacyMustNewDecFromStr("0.5"), uptime)
	})

	t.Run("returns valid uptime with wraparound", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(maximumViableWindow, val, cp, id)
		require.NoError(t, err)

		// Vote with price for 5 blocks
		for i := 0; i < 5; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		}

		for i := 0; i < 15; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.NoVote))
		}

		// Vote without price for 5 blocks
		for i := 0; i < 5; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))
		}

		uptime, err := sla.GetUptimeFromPriceFeed(priceFeed)
		require.NoError(t, err)
		require.Equal(t, math.LegacyMustNewDecFromStr("0.0"), uptime)
	})

	t.Run("returns valid uptime where every other block is a vote", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(maximumViableWindow, val, cp, id)
		require.NoError(t, err)

		for i := 0; i < 10; i++ {
			if i%2 == 0 {
				require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
			} else {
				require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))
			}
		}

		uptime, err := sla.GetUptimeFromPriceFeed(priceFeed)
		require.NoError(t, err)
		require.Equal(t, math.LegacyMustNewDecFromStr("0.5"), uptime)
	})

	t.Run("returns a valid uptime when price feed has longer window", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(40, val, cp, id)
		require.NoError(t, err)

		// Vote with price for 5 blocks
		for i := 0; i < 10; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		}

		// Vote without price for 5 blocks
		for i := 0; i < 10; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))
		}

		// vote yes for 10 blocks
		for i := 0; i < 20; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		}

		uptime, err := sla.GetUptimeFromPriceFeed(priceFeed)
		require.NoError(t, err)
		require.Equal(t, math.LegacyMustNewDecFromStr("1.0"), uptime)

		// create different strategy that has 40 block window
		strategy := slatypes.NewPriceFeedSLA(
			id,
			40,
			expectedUptime,
			k,
			uint64(minimumBlockUpdates),
			uint64(frequency),
		)

		uptime, err = strategy.GetUptimeFromPriceFeed(priceFeed)
		require.NoError(t, err)
		require.Equal(t, math.LegacyMustNewDecFromStr("0.75"), uptime)
	})

	t.Run("same price feed for different strategies", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(100, val, cp, id)
		require.NoError(t, err)

		// Vote with price for 5 blocks
		for i := 0; i < 10; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		}

		// Vote without price for 5 blocks
		for i := 0; i < 10; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))
		}

		// vote yes for 10 blocks
		for i := 0; i < 20; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		}

		// create different strategy that has 40 block window
		sla := slatypes.NewPriceFeedSLA(
			id,
			10,
			expectedUptime,
			k,
			uint64(minimumBlockUpdates),
			uint64(frequency),
		)

		uptime, err := sla.GetUptimeFromPriceFeed(priceFeed)
		require.NoError(t, err)
		require.Equal(t, math.LegacyMustNewDecFromStr("1.0"), uptime)

		// create different strategy that has 40 block window
		sla = slatypes.NewPriceFeedSLA(
			id,
			20,
			expectedUptime,
			k,
			uint64(minimumBlockUpdates),
			uint64(frequency),
		)

		uptime, err = sla.GetUptimeFromPriceFeed(priceFeed)
		require.NoError(t, err)
		require.Equal(t, math.LegacyMustNewDecFromStr("1.0"), uptime)

		// create different strategy that has 40 block window
		sla = slatypes.NewPriceFeedSLA(
			id,
			30,
			expectedUptime,
			k,
			uint64(minimumBlockUpdates),
			uint64(frequency),
		)

		uptime, err = sla.GetUptimeFromPriceFeed(priceFeed)
		require.NoError(t, err)
		require.Equal(t, math.LegacyNewDec(2).Quo(math.LegacyNewDec(3)), uptime)
	})

	t.Run("wrap around updates with multiple slas", func(t *testing.T) {
		priceFeed, err := slatypes.NewPriceFeed(20, val, cp, id)
		require.NoError(t, err)

		// Vote with price for 5 blocks
		for i := 0; i < 10; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		}

		// Vote without price for 5 blocks
		for i := 0; i < 10; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithoutPrice))
		}

		// vote yes for 10 blocks
		for i := 0; i < 10; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.VoteWithPrice))
		}

		// no vote for 10 blocks
		for i := 0; i < 5; i++ {
			require.NoError(t, priceFeed.SetUpdate(slatypes.NoVote))
		}

		sla := slatypes.NewPriceFeedSLA(
			id,
			10,
			expectedUptime,
			k,
			uint64(minimumBlockUpdates),
			uint64(frequency),
		)

		uptime, err := sla.GetUptimeFromPriceFeed(priceFeed)
		require.NoError(t, err)
		require.Equal(t, math.LegacyMustNewDecFromStr("1.0"), uptime)

		sla = slatypes.NewPriceFeedSLA(
			id,
			20,
			expectedUptime,
			k,
			uint64(minimumBlockUpdates),
			uint64(frequency),
		)

		uptime, err = sla.GetUptimeFromPriceFeed(priceFeed)
		require.NoError(t, err)
		require.Equal(t, math.LegacyNewDec(2).Quo(math.LegacyNewDec(3)), uptime)

		sla = slatypes.NewPriceFeedSLA(
			id,
			6,
			expectedUptime,
			k,
			uint64(minimumBlockUpdates),
			uint64(frequency),
		)

		uptime, err = sla.GetUptimeFromPriceFeed(priceFeed)
		require.NoError(t, err)
		require.Equal(t, math.LegacyMustNewDecFromStr("1.0"), uptime)
	})
}
