package types_test

import (
	"testing"

	"github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/assert"
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
