package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/oracle/types"
)

func TestGetSignersMsgAddCurrencyPairs(t *testing.T) {
	// create a msgAddCurrencyPairs
	auth := sdk.AccAddress("abc").String()
	msg := types.NewMsgAddCurrencyPairs(auth, nil)
	// get signers
	signer := msg.GetSigners()
	require.Equal(t, signer[0].String(), auth)
}

func TestGetSignersMsgRemoveCurrencyPairs(t *testing.T) {
	// create a msgAddCurrencyPairs
	auth := sdk.AccAddress("abc").String()
	msg := types.NewMsgRemoveCurrencyPairs(auth, nil)
	// get signers
	signer := msg.GetSigners()
	require.Equal(t, signer[0].String(), auth)
}

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
				Authority: sdk.AccAddress([]byte("abc")).String(),
				CurrencyPairs: []types.CurrencyPair{
					{Base: "A"},
				},
			},
			false,
		},
		{
			"if all currency pairs are valid + authority is valid - pass",
			types.MsgAddCurrencyPairs{
				Authority: sdk.AccAddress([]byte("abc")).String(),
				CurrencyPairs: []types.CurrencyPair{
					{Base: "A", Quote: "B", Decimals: types.DefaultDecimals},
					{Base: "C", Quote: "D", Decimals: types.DefaultDecimals},
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
				Authority: sdk.AccAddress([]byte("abc")).String(),
				CurrencyPairIds: []string{
					types.NewCurrencyPair("A", "B", types.DefaultDecimals).Ticker(),
					types.NewCurrencyPair("C", "D", types.DefaultDecimals).Ticker(),
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
