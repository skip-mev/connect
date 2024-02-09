package types_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/alerts/types"
)

func TestGenesisValidation(t *testing.T) {
	type testCase struct {
		name    string
		genesis types.GenesisState
		valid   bool
	}

	cases := []testCase{
		{
			"genesis with invalid params - fail",
			types.GenesisState{
				Params: types.NewParams(types.AlertParams{false, sdk.NewCoin("test", math.NewInt(1000000)), 0}, nil, types.PruningParams{}),
			},
			false,
		},
		{
			"genesis with valid params - pass",
			types.GenesisState{
				Params: types.NewParams(types.AlertParams{false, sdk.NewCoin("test", math.NewInt(0)), 0}, nil, types.PruningParams{}),
			},
			true,
		},
		{
			"genesis with an invalid alert - fail",
			types.NewGenesisState(types.NewParams(types.AlertParams{true, sdk.NewCoin("test", math.NewInt(1000000)), 1}, nil, types.PruningParams{}), []types.AlertWithStatus{
				types.NewAlertWithStatus(
					types.Alert{
						Height: 1,
						Signer: "",
					},
					types.NewAlertStatus(1, 2, time.Now(), 1),
				),
				types.NewAlertWithStatus(
					types.NewAlert(1, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BASE", "QUOTE")),
					types.NewAlertStatus(1, 2, time.Now(), 1),
				),
			}),
			false,
		},
		{
			"genesis with duplicate alerts - fail",
			types.NewGenesisState(types.NewParams(types.AlertParams{true, sdk.NewCoin("test", math.NewInt(1000000)), 1}, nil, types.PruningParams{}), []types.AlertWithStatus{
				types.NewAlertWithStatus(
					types.NewAlert(1, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BASE", "QUOTE")),
					types.NewAlertStatus(1, 2, time.Now(), 1),
				),
				types.NewAlertWithStatus(
					types.NewAlert(1, sdk.AccAddress("test1"), slinkytypes.NewCurrencyPair("BASE", "QUOTE")),
					types.NewAlertStatus(1, 2, time.Now(), 0),
				),
			}),
			false,
		},
		{
			"genesis with valid non-duplicatge alerts - pass",
			types.NewGenesisState(types.NewParams(types.AlertParams{true, sdk.NewCoin("test", math.NewInt(1000000)), 1}, nil, types.PruningParams{}), []types.AlertWithStatus{
				types.NewAlertWithStatus(
					types.NewAlert(1, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BASE", "QUOTE")),
					types.NewAlertStatus(1, 2, time.Now(), 1),
				),
				types.NewAlertWithStatus(
					types.NewAlert(0, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BASE2", "QUOTE2")),
					types.NewAlertStatus(1, 2, time.Now(), 0),
				),
			}),
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.genesis.ValidateBasic()
			if tc.valid && err != nil {
				t.Errorf("expected genesis to be valid, got error: %v", err)
			}
			if !tc.valid && err == nil {
				t.Errorf("expected genesis to be invalid, got nil error")
			}
		})
	}
}

func TestDefaultGenesisValidation(t *testing.T) {
	// test that the default genesis is valid
	genesis := types.DefaultGenesisState()
	assert.NoError(t, genesis.ValidateBasic())
}
