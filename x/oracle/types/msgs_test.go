package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/oracle/types"
)

func TestValidateBasicMsgAddCurrencyPairs(t *testing.T) {
	tcs := []struct {
		name       string
		msg        types.MsgAddCurrencyPairs
		expectPass bool
	}{
		{
			"if the authority is not an acc-address - fail",
			types.MsgAddCurrencyPairs{
				Authority: "abc",
			},
			false,
		},
		{
			"if any of the currency pairs are invalid - fail",
			types.MsgAddCurrencyPairs{
				Authority: sdk.AccAddress("abc").String(),
				CurrencyPairs: []slinkytypes.CurrencyPair{
					{Base: "A"},
				},
			},
			false,
		},
		{
			"if all currency pairs are valid + authority is valid - pass",
			types.MsgAddCurrencyPairs{
				Authority: sdk.AccAddress("abc").String(),
				CurrencyPairs: []slinkytypes.CurrencyPair{
					{Base: "A", Quote: "B"},
					{Base: "C", Quote: "D"},
				},
			},
			true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if !tc.expectPass {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestValidateBasicMsgRemoveCurrencyPairs(t *testing.T) {
	tcs := []struct {
		name       string
		msg        types.MsgRemoveCurrencyPairs
		expectPass bool
	}{
		{
			"if the authority is not an acc-address - fail",
			types.MsgRemoveCurrencyPairs{
				Authority: "abc",
			},
			false,
		},
		{
			"if any of the currency pairs are invalid - fail",
			types.MsgRemoveCurrencyPairs{
				Authority: sdk.AccAddress("abc").String(),
				CurrencyPairIds: []string{
					"AA",
				},
			},
			false,
		},
		{
			"if all currency pairs are valid + authority is valid - pass",
			types.MsgRemoveCurrencyPairs{
				Authority: sdk.AccAddress("abc").String(),
				CurrencyPairIds: []string{
					slinkytypes.CurrencyPairString("A", "B"),
					slinkytypes.CurrencyPairString("C", "D"),
				},
			},
			true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if !tc.expectPass {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}
