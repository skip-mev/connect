package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/testutil"
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
			"if Base formatted incorrectly as defi, Quote standard - fail",
			slinkytypes.CurrencyPair{
				Base:  "bB,testAddress,testChain",
				Quote: "AA",
			},
			false,
		},
		{
			"if Quote formatted correctly as defi, Base incorrectly  - fail",
			slinkytypes.CurrencyPair{
				Base:  "bB",
				Quote: "AA,testAddress,testChain",
			},
			false,
		},
		{
			"if both Quote + Base are formatted incorrectly as defi - fail",
			slinkytypes.CurrencyPair{
				Base:  "bB,testAddress,testChain",
				Quote: "AA,testAddress,testChain",
			},
			false,
		},
		{
			"if Base formatted correctly as defi, Quote incorrectly - fail",
			slinkytypes.CurrencyPair{
				Base:  "BB,testAddress,testChain",
				Quote: "aA",
			},
			false,
		},
		{
			"if Quote formatted incorrectly as defi, Quote standard - fail",
			slinkytypes.CurrencyPair{
				Base:  "BB",
				Quote: "aA,testAddress,testChain",
			},
			false,
		},
		{
			"if both Quote + Base are formatted incorrectly as defi - fail",
			slinkytypes.CurrencyPair{
				Base:  "BB,testAddress,testChain",
				Quote: "aA,testAddress,testChain",
			},
			false,
		},
		{
			"Base defi asset too many fields - fail",
			slinkytypes.CurrencyPair{
				Base:  "BB,testAddress,testChain,extra",
				Quote: "AA",
			},
			false,
		},
		{
			"Base defi asset too few fields - fail",
			slinkytypes.CurrencyPair{
				Base:  "BB,testAddress",
				Quote: "AA",
			},
			false,
		},
		{
			"Quote defi asset too many fields - fail",
			slinkytypes.CurrencyPair{
				Base:  "BB",
				Quote: "AA,testAddress,testChain,extra",
			},
			false,
		},
		{
			"Quote defi asset too few fields - fail",
			slinkytypes.CurrencyPair{
				Base:  "BB",
				Quote: "AA,testAddress",
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
			"Base is too long - fail",
			slinkytypes.CurrencyPair{
				Base:  testutil.RandomString(slinkytypes.MaxCPFieldLength + 1),
				Quote: "AA",
			},
			false,
		},
		{
			"Quote is too long - fail",
			slinkytypes.CurrencyPair{
				Base:  "BB",
				Quote: testutil.RandomString(slinkytypes.MaxCPFieldLength + 1),
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
		{
			"if Base formatted incorrectly as defi, Quote standard but rest lowercase - fail",
			slinkytypes.CurrencyPair{
				Base:  "BB,testAddress,testChain",
				Quote: "AA",
			},
			false,
		},
		{
			"if Quote formatted incorrectly as Base, Quote standard but rest lowercase - fail",
			slinkytypes.CurrencyPair{
				Base:  "BB",
				Quote: "AA,testAddress,testChain",
			},
			false,
		},
		{
			"if both Quote + Base are formatted correctly as defi but rest lowercase - fail",
			slinkytypes.CurrencyPair{
				Base:  "BB,testAddress,testChain",
				Quote: "AA,testAddress,testChain",
			},
			false,
		},
		{
			"if Base formatted incorrectly as defi, Quote standard - pass",
			slinkytypes.CurrencyPair{
				Base:  "BB,TESTADDRESS,TESTCHAIN",
				Quote: "AA",
			},
			true,
		},
		{
			"if Quote formatted incorrectly as Base, Quote standard - pass",
			slinkytypes.CurrencyPair{
				Base:  "BB",
				Quote: "AA,TESTADDRESS,TESTCHAIN",
			},
			true,
		},
		{
			"if both Quote + Base are formatted correctly as defi - pass",
			slinkytypes.CurrencyPair{
				Base:  "BB,TESTADDRESS,TESTCHAIN",
				Quote: "AA,TESTADDRESS,TESTCHAIN",
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
			"if string is incorrectly formatted (defi), return an empty CurrencyPair",
			"a,a,a,a,a,a/a",
			slinkytypes.CurrencyPair{},
			false,
		},
		{
			"if string is incorrectly formatted (empty), return an empty CurrencyPair",
			"",
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
		{
			"if the string is not formatted upper-case, return the original CurrencyPair",
			"A/b",
			slinkytypes.CurrencyPair{Base: "A", Quote: "B"},
			true,
		},
		{
			"if the string is not formatted upper-case (defi), return the original CurrencyPair",
			"a,testAddress,testChain/B",
			slinkytypes.CurrencyPair{Base: "A,TESTADDRESS,TESTCHAIN", Quote: "B"},
			true,
		},
		{
			"if the string is not formatted upper-case (defi), return the original CurrencyPair",
			"a/b,testAddress,testChain",
			slinkytypes.CurrencyPair{Base: "A", Quote: "B,TESTADDRESS,TESTCHAIN"},
			true,
		},
		{
			"if the string is not formatted upper-case (defi), return the original CurrencyPair",
			"A,testAddress,testChain/B,testAddress,testChain",
			slinkytypes.CurrencyPair{Base: "A,TESTADDRESS,TESTCHAIN", Quote: "B,TESTADDRESS,TESTCHAIN"},
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
		{
			"if the CurrencyPairs are equal, return true - defi",
			slinkytypes.CurrencyPair{Base: "A,testAddress,testChain", Quote: "B"},
			slinkytypes.CurrencyPair{Base: "A,testAddress,testChain", Quote: "B"},
			true,
		},
		{
			"if the CurrencyPairs are not equal, return false - defi",
			slinkytypes.CurrencyPair{Base: "A,testAddress,testChain", Quote: "B"},
			slinkytypes.CurrencyPair{Base: "B,testAddress,testChain", Quote: "A"},
			false,
		},
		{
			"if the CurrencyPairs are equal, return true - defi",
			slinkytypes.CurrencyPair{Base: "A,testAddress,testChain", Quote: "B,testAddress,testChain"},
			slinkytypes.CurrencyPair{Base: "A,testAddress,testChain", Quote: "B,testAddress,testChain"},
			true,
		},
		{
			"if the CurrencyPairs are not equal, return false - defi",
			slinkytypes.CurrencyPair{Base: "A,testAddress,testChain", Quote: "B,testAddress,testChain"},
			slinkytypes.CurrencyPair{Base: "B,testAddress,testChain", Quote: "A,testAddress,testChain"},
			false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.cp1.Equal(tc.cp2), tc.eq)
		})
	}
}
