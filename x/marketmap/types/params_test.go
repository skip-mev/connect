package types_test

import (
	"testing"

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
			name:      "invalid authority",
			params:    types.NewParams("incorrect"),
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
