package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/skip-mev/slinky/x/oracle/types"
)

func TestGetSignersMsgAddCurrencyPairs(t *testing.T) {
	// create a msgAddCurrencyPairs
	auth := sdk.AccAddress([]byte("abc")).String()
	msg := types.NewMsgAddCurrencyPairs(auth, nil)
	// get signers
	signer := msg.GetSigners()
	assert.Equal(t, signer[0].String(), auth)
}

func TestGetSignersMsgRemoveCurrencyPairs(t *testing.T) {
	// create a msgAddCurrencyPairs
	auth := sdk.AccAddress([]byte("abc")).String()
	msg := types.NewMsgRemoveCurrencyPairs(auth, nil)
	// get signers
	signer := msg.GetSigners()
	assert.Equal(t, signer[0].String(), auth)
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
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
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
				Authority: sdk.AccAddress([]byte("abc")).String(),
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
					types.CurrencyPair{Base: "A", Quote: "B"}.ToString(),
					types.CurrencyPair{Base: "C", Quote: "D"}.ToString(),
				},
			},
			true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if !tc.expectPass {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
