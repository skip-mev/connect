package types_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

func TestGenesisState(t *testing.T) {
	t.Run("default genesis state is valid", func(t *testing.T) {
		gs := slatypes.NewDefaultGenesisState()
		err := gs.ValidateBasic()
		require.NoError(t, err)
	})

	validSLA := slatypes.NewPriceFeedSLA(
		"test",
		10,
		math.LegacyMustNewDecFromStr("0.5"),
		math.LegacyMustNewDecFromStr("0.5"),
		5,
		5,
	)

	invalidSLA := slatypes.NewPriceFeedSLA(
		"test2",
		0,
		math.LegacyMustNewDecFromStr("0.0"),
		math.LegacyMustNewDecFromStr("0.0"),
		0,
		0,
	)

	validSLA2 := slatypes.NewPriceFeedSLA(
		"test2",
		10,
		math.LegacyMustNewDecFromStr("0.5"),
		math.LegacyMustNewDecFromStr("0.5"),
		5,
		5,
	)

	val1 := sdk.ConsAddress("val1")

	cp1 := slinkytypes.NewCurrencyPair("BTC", "USD")

	goodFeed1, err := slatypes.NewPriceFeed(10, val1, cp1, "test")
	require.NoError(t, err)

	goodFeed2, err := slatypes.NewPriceFeed(10, val1, cp1, "test2")
	require.NoError(t, err)

	badPriceFeed1, err := slatypes.NewPriceFeed(10, val1, cp1, "no match sla")
	require.NoError(t, err)

	badPriceFeed2, err := slatypes.NewPriceFeed(11, val1, cp1, "test")
	require.NoError(t, err)

	t.Run("genesis state with duplicate ids should be rejected", func(t *testing.T) {
		gs := slatypes.NewGenesisState([]slatypes.PriceFeedSLA{validSLA, validSLA}, nil, slatypes.DefaultParams())
		err := gs.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("genesis state with invalid sla should be rejected", func(t *testing.T) {
		gs := slatypes.NewGenesisState([]slatypes.PriceFeedSLA{invalidSLA}, nil, slatypes.DefaultParams())
		err := gs.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("genesis state with valid slas should be accepted", func(t *testing.T) {
		gs := slatypes.NewGenesisState([]slatypes.PriceFeedSLA{validSLA, validSLA2}, nil, slatypes.DefaultParams())
		err := gs.ValidateBasic()
		require.NoError(t, err)
	})

	t.Run("genesis state with valid price feed", func(t *testing.T) {
		gs := slatypes.NewGenesisState([]slatypes.PriceFeedSLA{validSLA}, []slatypes.PriceFeed{goodFeed1}, slatypes.DefaultParams())
		err := gs.ValidateBasic()
		require.NoError(t, err)
	})

	t.Run("genesis state with invalid price feed that has no matching SLA", func(t *testing.T) {
		gs := slatypes.NewGenesisState([]slatypes.PriceFeedSLA{validSLA}, []slatypes.PriceFeed{badPriceFeed1}, slatypes.DefaultParams())
		err := gs.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("sla has a different maximum viable window than the price feed with same id", func(t *testing.T) {
		gs := slatypes.NewGenesisState([]slatypes.PriceFeedSLA{validSLA}, []slatypes.PriceFeed{badPriceFeed2}, slatypes.DefaultParams())
		err := gs.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("duplicate price feeds", func(t *testing.T) {
		gs := slatypes.NewGenesisState([]slatypes.PriceFeedSLA{validSLA}, []slatypes.PriceFeed{goodFeed1, goodFeed1}, slatypes.DefaultParams())
		err := gs.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("same (val, cp) pair for different SLAs", func(t *testing.T) {
		gs := slatypes.NewGenesisState([]slatypes.PriceFeedSLA{validSLA, validSLA2}, []slatypes.PriceFeed{goodFeed1, goodFeed2}, slatypes.DefaultParams())
		err := gs.ValidateBasic()
		require.NoError(t, err)
	})
}
