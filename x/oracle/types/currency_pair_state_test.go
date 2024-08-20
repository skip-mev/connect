package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/x/oracle/types"
)

func TestCurrencyPairState(t *testing.T) {
	tcs := []struct {
		name  string
		cps   types.CurrencyPairState
		valid bool
	}{
		{
			"non-zero nonce, and nil price - invalid",
			types.CurrencyPairState{
				Nonce: 1,
				Price: nil,
			},
			false,
		},
		{
			"zero nonce, and non-nil price - invalid",
			types.CurrencyPairState{
				Nonce: 0,
				Price: &types.QuotePrice{
					Price: math.NewInt(1),
				},
			},
			false,
		},
		{
			"zero nonce, and nil price - valid",
			types.CurrencyPairState{
				Nonce: 0,
				Price: nil,
			},
			true,
		},
		{
			"non-zero nonce, and non-nil price - valid",
			types.CurrencyPairState{
				Nonce: 1,
				Price: &types.QuotePrice{
					Price: math.NewInt(1),
				},
			},
			true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.cps.ValidateBasic() == nil, tc.valid)
		})
	}
}
