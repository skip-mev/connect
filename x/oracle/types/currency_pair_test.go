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
				Base:  "aB",
				Quote: "BB",
			},
			false,
		},
		{
			"if the Quote is not upper-case - fail",
			types.CurrencyPair{
				Base:  "BB",
				Quote: "aB",
			},
			false,
		},
		{
			"if the base string is empty - fail",
			types.CurrencyPair{
				Base:  "",
				Quote: "BB",
			},
			false,
		},
		{
			"if the quote string is empty - fail",
			types.CurrencyPair{
				Base:  "AA",
				Quote: "",
			},
			false,
		},
		{
			"if both Quote + Base are formatted correctly - pass",
			types.CurrencyPair{
				Base:  "BB",
				Quote: "AA",
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

func TestToFromString(t *testing.T) {
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
			"if the string is correctly formatted, return the original CurrencyPair",
			types.CurrencyPair{Base: "A", Quote: "B"}.ToString(),
			types.CurrencyPair{Base: "A", Quote: "B"},
			true,
		},
		{
			"if the string is not formatted upper-case, return the original CurrencyPair",
			"a/B",
			types.CurrencyPair{Base: "A", Quote: "B"},
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

func TestDecimals(t *testing.T) {
	tcs := []struct {
		name string
		cp   types.CurrencyPair
		dec  int
	}{
		{
			"if the quote is ethereum, return 18",
			types.CurrencyPair{Base: "A", Quote: "ETHEREUM"},
			18,
		},
		{
			"if the quote is not ethereum or eth, return 8",
			types.CurrencyPair{Base: "A", Quote: "B"},
			8,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.cp.Decimals(), tc.dec)
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
			assert.Equal(t, tc.cps.ValidateBasic() == nil, tc.valid)
		})
	}
}
