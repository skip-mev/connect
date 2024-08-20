package types_test

import (
	"testing"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/x/alerts/types"
)

func TestParamsValidation(t *testing.T) {
	type testCase struct {
		name string
		// params to validate
		params types.Params
		// expected result
		valid bool
	}

	pk := secp256k1.GenPrivKey()
	pkany, err := codectypes.NewAnyWithValue(pk.PubKey())
	require.NoError(t, err)

	cases := []testCase{
		{
			name: "not-enabled, but bond amount is non-zero - fail",
			params: types.NewParams(
				types.AlertParams{BondAmount: sdk.NewCoin("test", math.NewInt(1000000))}, nil, types.PruningParams{}),
		},
		{
			name:   "not enabled and bond amount is zero - pass",
			params: types.NewParams(types.AlertParams{BondAmount: sdk.NewCoin("test", math.NewInt(0))}, nil, types.PruningParams{}),
			valid:  true,
		},
		{
			name:   "not enabled and max-block-age is non-zero -fail",
			params: types.NewParams(types.AlertParams{BondAmount: sdk.NewCoin("test", math.NewInt(0)), MaxBlockAge: 1}, nil, types.PruningParams{}),
		},
		{
			name:   "enabled, but max-block-age is zero - fail",
			params: types.NewParams(types.AlertParams{Enabled: true, BondAmount: sdk.NewCoin("test", math.NewInt(1))}, nil, types.PruningParams{}),
		},
		{
			name:   "enabled, but bond amount is zero - fail",
			params: types.NewParams(types.AlertParams{Enabled: true, BondAmount: sdk.NewCoin("test", math.NewInt(0)), MaxBlockAge: 1}, nil, types.PruningParams{}),
		},
		{
			name:   "enabled, but bond amount is negative - fail",
			params: types.NewParams(types.AlertParams{Enabled: true, BondAmount: sdk.Coin{Denom: "test", Amount: math.NewInt(-1)}, MaxBlockAge: 1}, nil, types.PruningParams{}),
		},
		{
			name:   "enabled, but bond amount is non-zero - pass",
			params: types.NewParams(types.AlertParams{Enabled: true, BondAmount: sdk.NewCoin("test", math.NewInt(1000000)), MaxBlockAge: 1}, nil, types.PruningParams{}),
			valid:  true,
		},
		{
			name: "valid alert params, but invalid conclusion verification params - fail",
			params: types.NewParams(types.AlertParams{Enabled: true, BondAmount: sdk.NewCoin("test", math.NewInt(1000000)), MaxBlockAge: 1}, &types.MultiSigConclusionVerificationParams{
				Signers: []*codectypes.Any{pkany, pkany},
			}, types.PruningParams{}),
		},
		{
			name: "valid alert params, but invalid pruning params (disabled + non-zero blocks-to-prune) - fail",
			params: types.NewParams(types.AlertParams{Enabled: true, BondAmount: sdk.NewCoin("test", math.NewInt(1000000)), MaxBlockAge: 1}, nil, types.PruningParams{
				Enabled:       false,
				BlocksToPrune: 10,
			}),
		},
		{
			name: "valid alert params, but invalid pruning params (enabled + zero blocks-to-prune) - fail",
			params: types.NewParams(types.AlertParams{Enabled: true, BondAmount: sdk.NewCoin("test", math.NewInt(1000000)), MaxBlockAge: 1}, nil, types.PruningParams{
				Enabled: true,
			}),
		},
		{
			name: "valid alert params, and valid pruning params (enabled + non-zero blocks-to-prune) - pass",
			params: types.NewParams(types.AlertParams{Enabled: true, BondAmount: sdk.NewCoin("test", math.NewInt(1000000)), MaxBlockAge: 1}, nil, types.PruningParams{
				Enabled:       true,
				BlocksToPrune: 10,
			}),
			valid: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.valid && err != nil {
				t.Errorf("expected params to be valid, got error: %v", err)
			}
			if !tc.valid && err == nil {
				t.Errorf("expected params to be invalid, got nil error")
			}
		})
	}
}

func TestDefaultParamsValidate(t *testing.T) {
	// test that the default params are valid
	params := types.DefaultParams("test", nil)
	require.NoError(t, params.Validate())
}
