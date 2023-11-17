package types_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/x/sla/types"
)

func TestMsgAddSLAs(t *testing.T) {
	t.Run("should reject a msg with an invalid authority address", func(t *testing.T) {
		msg := types.NewMsgAddSLAs("invalid", []types.PriceFeedSLA{})
		err := msg.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("should accept an empty message with a valid authority address", func(t *testing.T) {
		msg := types.NewMsgAddSLAs(sdk.AccAddress([]byte("test")).String(), []types.PriceFeedSLA{})
		err := msg.ValidateBasic()
		require.NoError(t, err)
	})

	invalidSLA := types.NewPriceFeedSLA(
		"test",
		0,
		math.LegacyMustNewDecFromStr("0.0"),
		math.LegacyMustNewDecFromStr("0.0"),
		0,
		0,
	)

	validSLA := types.NewPriceFeedSLA(
		"test",
		10,
		math.LegacyMustNewDecFromStr("0.5"),
		math.LegacyMustNewDecFromStr("0.5"),
		5,
		5,
	)

	validSLA2 := types.NewPriceFeedSLA(
		"test2",
		10,
		math.LegacyMustNewDecFromStr("0.5"),
		math.LegacyMustNewDecFromStr("0.5"),
		5,
		5,
	)

	t.Run("should reject a message with an invalid sla", func(t *testing.T) {
		msg := types.NewMsgAddSLAs(sdk.AccAddress([]byte("test")).String(), []types.PriceFeedSLA{invalidSLA})
		err := msg.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("should reject a message with duplicate slas", func(t *testing.T) {
		msg := types.NewMsgAddSLAs(sdk.AccAddress([]byte("test")).String(), []types.PriceFeedSLA{validSLA, validSLA})
		err := msg.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("should accept a message with valid slas", func(t *testing.T) {
		msg := types.NewMsgAddSLAs(sdk.AccAddress([]byte("test")).String(), []types.PriceFeedSLA{validSLA, validSLA2})
		err := msg.ValidateBasic()
		require.NoError(t, err)
	})
}

func TestMsgRemoveSLAs(t *testing.T) {
	t.Run("should reject a msg with an invalid authority address", func(t *testing.T) {
		msg := types.NewMsgRemoveSLAs("invalid", []string{})
		err := msg.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("should accept an empty message with a valid authority address", func(t *testing.T) {
		msg := types.NewMsgRemoveSLAs(sdk.AccAddress([]byte("test")).String(), []string{})
		err := msg.ValidateBasic()
		require.NoError(t, err)
	})

	t.Run("should reject a message with duplicate ids", func(t *testing.T) {
		msg := types.NewMsgRemoveSLAs(sdk.AccAddress([]byte("test")).String(), []string{"test", "test"})
		err := msg.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("should accept a message with valid ids", func(t *testing.T) {
		msg := types.NewMsgRemoveSLAs(sdk.AccAddress([]byte("test")).String(), []string{"test", "test2"})
		err := msg.ValidateBasic()
		require.NoError(t, err)
	})
}

func TestMsgParams(t *testing.T) {
	t.Run("should reject a message with an invalid authority address", func(t *testing.T) {
		msg := types.NewMsgParams("invalid", types.DefaultParams())
		err := msg.ValidateBasic()
		require.Error(t, err)
	})

	t.Run("should accept an empty message with a valid authority address", func(t *testing.T) {
		msg := types.NewMsgParams(sdk.AccAddress([]byte("test")).String(), types.DefaultParams())
		err := msg.ValidateBasic()
		require.NoError(t, err)
	})
}
