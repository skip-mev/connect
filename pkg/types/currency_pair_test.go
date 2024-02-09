package types_test

import (
	"testing"

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
			types.CurrencyPairString("A", "B"),
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
