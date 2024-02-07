package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/assert"

	"github.com/skip-mev/slinky/x/oracle/types"
)

func TestValidateBasic(t *testing.T) {
	tcs := []struct {
		name       string
		cp         types.CurrencyPair
		expectPass bool
	}{
		{
			"if the Base is not upper-case - fail",
			types.CurrencyPair{
				Base:     "aB",
				Quote:    "BB",
				Decimals: types.DefaultDecimals,
			},
			false,
		},
		{
			"if the Quote is not upper-case - fail",
			types.CurrencyPair{
				Base:     "BB",
				Quote:    "aB",
				Decimals: types.DefaultDecimals,
			},
			false,
		},
		{
			"if the base string is empty - fail",
			types.CurrencyPair{
				Base:     "",
				Quote:    "BB",
				Decimals: types.DefaultDecimals,
			},
			false,
		},
		{
			"if the quote string is empty - fail",
			types.CurrencyPair{
				Base:     "AA",
				Quote:    "",
				Decimals: types.DefaultDecimals,
			},
			false,
		},
		{
			"invalid decimals - fail",
			types.CurrencyPair{
				Base:     "BB",
				Quote:    "AA",
				Decimals: 0,
			},
			false,
		},
		{
			"if both Quote + Base are formatted correctly - pass",
			types.CurrencyPair{
				Base:     "BB",
				Quote:    "AA",
				Decimals: types.DefaultDecimals,
			},
			true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cp.ValidateBasic()
			switch tc.expectPass {
			case true:
				assert.Nil(t, err)
			default:
				assert.NotNil(t, err)
			}
		})
	}
}

func TestToFromFullString(t *testing.T) {
	tcs := []struct {
		name string
		// string formatted CurrencyPair
		cps        string
		cp         types.CurrencyPair
		expectPass bool
	}{
		{
			"if string is incorrectly formatted, return an empty CurrencyPair",
			"aa",
			types.CurrencyPair{},
			false,
		},
		{
			"if string is incorrectly formatted, return an empty CurrencyPair",
			"a/a/a",
			types.CurrencyPair{},
			false,
		},
		{
			"if the string is correctly formatted, return the original CurrencyPair",
			types.CurrencyPairFullString("A", "B", types.DefaultDecimals),
			types.CurrencyPair{Base: "A", Quote: "B", Decimals: types.DefaultDecimals},
			true,
		},
		{
			"if the string is not formatted upper-case, return the original CurrencyPair",
			"a/B/8",
			types.CurrencyPair{Base: "A", Quote: "B", Decimals: types.DefaultDecimals},
			true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			cp, err := types.CurrencyPairFromString(tc.cps)
			if tc.expectPass {
				assert.Nil(t, err)
				assert.Equal(t, cp, tc.cp)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestCurrencyPairState(t *testing.T) {
	tcs := []struct {
		name  string
		cps   types.CurrencyPairState
		valid bool
	}{
		{
			"non-zero nonce, and nil price - invalid",
			types.CurrencyPairState{
				Nonce:    1,
				Price:    nil,
				Decimals: types.DefaultDecimals,
			},
			false,
		},
		{
			"0 decimals - invalid",
			types.CurrencyPairState{
				Nonce: 1,
				Price: &types.QuotePrice{
					Price: math.NewInt(1),
				},
				Decimals: 0,
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
				Decimals: types.DefaultDecimals,
			},
			false,
		},
		{
			"zero nonce, and nil price - valid",
			types.CurrencyPairState{
				Nonce:    0,
				Price:    nil,
				Decimals: types.DefaultDecimals,
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
				Decimals: types.DefaultDecimals,
			},
			true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.cps.ValidateBasic() == nil, tc.valid)
		})
	}
}
