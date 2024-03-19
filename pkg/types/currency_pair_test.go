package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
)

func TestValidateBasic(t *testing.T) {
	tcs := []struct {
		name       string
		cp         slinkytypes.CurrencyPair
		expectPass bool
	}{
		{
			"if the Base is not upper-case - fail",
			slinkytypes.CurrencyPair{
				Base:  "aB",
				Quote: "BB",
			},
			false,
		},
		{
			"if the Quote is not upper-case - fail",
			slinkytypes.CurrencyPair{
				Base:  "BB",
				Quote: "aB",
			},
			false,
		},
		{
			"if the base string is empty - fail",
			slinkytypes.CurrencyPair{
				Base:  "",
				Quote: "BB",
			},
			false,
		},
		{
			"if the quote string is empty - fail",
			slinkytypes.CurrencyPair{
				Base:  "AA",
				Quote: "",
			},
			false,
		},
		{
			"if both Quote + Base are formatted correctly - pass",
			slinkytypes.CurrencyPair{
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
				require.Nil(t, err)
			default:
				require.NotNil(t, err)
			}
		})
	}
}

func TestToFromString(t *testing.T) {
	tcs := []struct {
		name string
		// string formatted CurrencyPair
		cps        string
		cp         slinkytypes.CurrencyPair
		expectPass bool
	}{
		{
			"if string is incorrectly formatted, return an empty CurrencyPair",
			"aa",
			slinkytypes.CurrencyPair{},
			false,
		},
		{
			"if the string is correctly formatted, return the original CurrencyPair",
			slinkytypes.CurrencyPairString("A", "B"),
			slinkytypes.CurrencyPair{Base: "A", Quote: "B"},
			true,
		},
		{
			"if the string is not formatted upper-case, return the original CurrencyPair",
			"a/B",
			slinkytypes.CurrencyPair{Base: "A", Quote: "B"},
			true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			cp, err := slinkytypes.CurrencyPairFromString(tc.cps)
			if tc.expectPass {
				require.Nil(t, err)
				require.Equal(t, cp, tc.cp)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

func TestDecimals(t *testing.T) {
	tcs := []struct {
		name string
		cp   slinkytypes.CurrencyPair
		dec  int
	}{
		{
			"if the quote is ethereum, return 18",
			slinkytypes.CurrencyPair{Base: "A", Quote: "ETHEREUM"},
			18,
		},
		{
			"if the quote is not ethereum or eth, return 8",
			slinkytypes.CurrencyPair{Base: "A", Quote: "B"},
			8,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.cp.LegacyDecimals(), tc.dec)
		})
	}
}

func TestEqual(t *testing.T) {
	tcs := []struct {
		name string
		cp1  slinkytypes.CurrencyPair
		cp2  slinkytypes.CurrencyPair
		eq   bool
	}{
		{
			"if the CurrencyPairs are equal, return true",
			slinkytypes.CurrencyPair{Base: "A", Quote: "B"},
			slinkytypes.CurrencyPair{Base: "A", Quote: "B"},
			true,
		},
		{
			"if the CurrencyPairs are not equal, return false",
			slinkytypes.CurrencyPair{Base: "A", Quote: "B"},
			slinkytypes.CurrencyPair{Base: "B", Quote: "A"},
			false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.cp1.Equal(tc.cp2), tc.eq)
		})
	}
}
