package types_test

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/x/marketmap/types"
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
				MarketAuthorities: []string{authtypes.NewModuleAddress(authtypes.ModuleName).String(), authtypes.NewModuleAddress(govtypes.ModuleName).String()},
				Admin:             authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			},
			expectErr: false,
		},
		{
			name: "invalid admin",
			params: types.Params{
				MarketAuthorities: []string{authtypes.NewModuleAddress(authtypes.ModuleName).String(), authtypes.NewModuleAddress(govtypes.ModuleName).String()},
				Admin:             "invalid",
			},
			expectErr: true,
		},
		{
			name: "invalid duplicate authority",
			params: types.Params{
				MarketAuthorities: []string{authtypes.NewModuleAddress(govtypes.ModuleName).String(), authtypes.NewModuleAddress(govtypes.ModuleName).String()},
				Admin:             authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			},
			expectErr: true,
		},
		{
			name: "invalid authority string",
			params: types.Params{
				MarketAuthorities: []string{"incorrect"},
				Admin:             authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			},
			expectErr: true,
		},
		{
			name: "invalid nil authority",
			params: types.Params{
				MarketAuthorities: nil,
				Admin:             authtypes.NewModuleAddress(govtypes.ModuleName).String(),
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
