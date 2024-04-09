package types_test

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

func TestValidateBasic(t *testing.T) {
	type testCase struct {
		name      string
		params    types.Params
		expectErr bool
	}

	testCases := []testCase{
		{
			name:      "valid default params",
			params:    types.DefaultParams(),
			expectErr: false,
		},
		{
			name: "valid multiple authorities",
			params: types.Params{
				MarketAuthorities: []string{authtypes.NewModuleAddress(authtypes.ModuleName).String(), types.DefaultMarketAuthority},
				Admin:             types.DefaultAdmin,
			},
			expectErr: false,
		},
		{
			name: "invalid admin",
			params: types.Params{
				MarketAuthorities: []string{authtypes.NewModuleAddress(authtypes.ModuleName).String(), types.DefaultMarketAuthority},
				Admin:             "invalid",
			},
			expectErr: true,
		},
		{
			name: "invalid duplicate authority",
			params: types.Params{
				MarketAuthorities: []string{types.DefaultMarketAuthority, types.DefaultMarketAuthority},
				Admin:             types.DefaultAdmin,
			},
			expectErr: true,
		},
		{
			name: "invalid authority string",
			params: types.Params{
				MarketAuthorities: []string{"incorrect"},
				Admin:             types.DefaultAdmin,
			},
			expectErr: true,
		},
		{
			name: "invalid nil authority",
			params: types.Params{
				MarketAuthorities: nil,
				Admin:             types.DefaultAdmin,
			},
			expectErr: true,
		},
		{
			name:      "invalid empty params",
			params:    types.Params{},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.ValidateBasic()
			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
