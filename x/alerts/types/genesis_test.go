package types_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/skip-mev/slinky/x/alerts/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

func TestGenesisValidation(t *testing.T) {
	type testCase struct {
		name    string
		genesis types.GenesisState
		valid   bool
	}

	cases := []testCase{
		{
			name: "genesis with invalid params - fail",
			genesis: types.GenesisState{
				Params: types.NewParams(types.AlertParams{BondAmount: sdk.NewCoin("test", math.NewInt(1000000))}, nil, types.PruningParams{}),
			},
		},
		{
			name: "genesis with valid params - pass",
			genesis: types.GenesisState{
				Params: types.NewParams(types.AlertParams{BondAmount: sdk.NewCoin("test", math.NewInt(0))}, nil, types.PruningParams{}),
			},
			valid: true,
		},
		{
			name: "genesis with an invalid alert - fail",
			genesis: types.NewGenesisState(types.NewParams(types.AlertParams{Enabled: true, BondAmount: sdk.NewCoin("test", math.NewInt(1000000)), MaxBlockAge: 1}, nil, types.PruningParams{}), []types.AlertWithStatus{
				types.NewAlertWithStatus(
					types.Alert{
						Height: 1,
						Signer: "",
					},
					types.NewAlertStatus(1, 2, time.Now(), 1),
				),
				types.NewAlertWithStatus(
					types.NewAlert(1, sdk.AccAddress("test"), oracletypes.NewCurrencyPair("BASE", "QUOTE")),
					types.NewAlertStatus(1, 2, time.Now(), 1),
				),
			}),
		},
		{
			name: "genesis with duplicate alerts - fail",
			genesis: types.NewGenesisState(types.NewParams(types.AlertParams{Enabled: true, BondAmount: sdk.NewCoin("test", math.NewInt(1000000)), MaxBlockAge: 1}, nil, types.PruningParams{}), []types.AlertWithStatus{
				types.NewAlertWithStatus(
					types.NewAlert(1, sdk.AccAddress("test"), oracletypes.NewCurrencyPair("BASE", "QUOTE")),
					types.NewAlertStatus(1, 2, time.Now(), 1),
				),
				types.NewAlertWithStatus(
					types.NewAlert(1, sdk.AccAddress("test1"), oracletypes.NewCurrencyPair("BASE", "QUOTE")),
					types.NewAlertStatus(1, 2, time.Now(), 0),
				),
			}),
		},
		{
			name: "genesis with valid non-duplicate alerts - pass",
			genesis: types.NewGenesisState(types.NewParams(types.AlertParams{Enabled: true, BondAmount: sdk.NewCoin("test", math.NewInt(1000000)), MaxBlockAge: 1}, nil, types.PruningParams{}), []types.AlertWithStatus{
				types.NewAlertWithStatus(
					types.NewAlert(1, sdk.AccAddress("test"), oracletypes.NewCurrencyPair("BASE", "QUOTE")),
					types.NewAlertStatus(1, 2, time.Now(), 1),
				),
				types.NewAlertWithStatus(
					types.NewAlert(0, sdk.AccAddress("test"), oracletypes.NewCurrencyPair("BASE2", "QUOTE2")),
					types.NewAlertStatus(1, 2, time.Now(), 0),
				),
			}),
			valid: true,
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
