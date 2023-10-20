package types_test

import (
	"testing"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/x/alerts/types"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	cases := []testCase{
		{
			"not-enabled, but bond amount is non-zero - fail",
			types.NewParams(
				types.AlertParams{false, sdk.NewCoin("test", math.NewInt(1000000)), 0}, nil, types.PruningParams{}),
			false,
		},
		{
			"not enabled and bond amount is zero - pass",
			types.NewParams(types.AlertParams{false, sdk.NewCoin("test", math.NewInt(0)), 0}, nil, types.PruningParams{}),
			true,
		},
		{
			"not enabled and max-block-age is non-zero -fail",
			types.NewParams(types.AlertParams{false, sdk.NewCoin("test", math.NewInt(0)), 1}, nil, types.PruningParams{}),
			false,
		},
		{
			"enabled, but max-block-age is zero - fail",
			types.NewParams(types.AlertParams{true, sdk.NewCoin("test", math.NewInt(1)), 0}, nil, types.PruningParams{}),
			false,
		},
		{
			"enabled, but bond amount is zero - fail",
			types.NewParams(types.AlertParams{true, sdk.NewCoin("test", math.NewInt(0)), 1}, nil, types.PruningParams{}),
			false,
		},
		{
			"enabled, but bond amount is negative - fail",
			types.NewParams(types.AlertParams{true, sdk.Coin{Denom: "test", Amount: math.NewInt(-1)}, 1}, nil, types.PruningParams{}),
			false,
		},
		{
			"enabled, but bond amount is non-zero - pass",
			types.NewParams(types.AlertParams{true, sdk.NewCoin("test", math.NewInt(1000000)), 1}, nil, types.PruningParams{}),
			true,
		},
		{
			"valid alert params, but invalid conclusion verification params - fail",
			types.NewParams(types.AlertParams{true, sdk.NewCoin("test", math.NewInt(1000000)), 1}, &types.MultiSigConclusionVerificationParams{
				Signers: []*codectypes.Any{pkany, pkany},
			}, types.PruningParams{}),
			false,
		},
		{
			"valid alert params, but invalid pruning params (disabled + non-zero blocks-to-prune) - fail",
			types.NewParams(types.AlertParams{true, sdk.NewCoin("test", math.NewInt(1000000)), 1}, nil, types.PruningParams{
				Enabled:       false,
				BlocksToPrune: 10,
			}),
			false,
		},
		{
			"valid alert params, but invalid pruning params (enabled + zero blocks-to-prune) - fail",
			types.NewParams(types.AlertParams{true, sdk.NewCoin("test", math.NewInt(1000000)), 1}, nil, types.PruningParams{
				Enabled: true,
			}),
			false,
		},
		{
			"valid alert params, and valid pruning params (enabled + non-zero blocks-to-prune) - pass",
			types.NewParams(types.AlertParams{true, sdk.NewCoin("test", math.NewInt(1000000)), 1}, nil, types.PruningParams{
				Enabled:       true,
				BlocksToPrune: 10,
			}),
			true,
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
	assert.NoError(t, params.Validate())
}
