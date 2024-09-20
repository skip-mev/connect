package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
)

var (
	_ sdk.Msg = &MsgAddCurrencyPairs{}
	_ sdk.Msg = &MsgRemoveCurrencyPairs{}
)

// NewMsgAddCurrencyPairs returns a new message from a set of currency-pairs and an authority.
func NewMsgAddCurrencyPairs(authority string, cps []connecttypes.CurrencyPair) MsgAddCurrencyPairs {
	return MsgAddCurrencyPairs{
		Authority:     authority,
		CurrencyPairs: cps,
	}
}

// ValidateBasic determines whether the information in the message is formatted correctly, specifically
// whether the authority is a valid acc-address, and that each CurrencyPair in the message is formatted correctly.

func (m *MsgAddCurrencyPairs) ValidateBasic() error {
	// validate authority address
	_, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		return err
	}

	// validate currency pairs
	for _, cp := range m.CurrencyPairs {
		if err := cp.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}

// NewMsgRemoveCurrencyPairs returns a new message to remove a set of currency-pairs from the x/oracle module's state.
func NewMsgRemoveCurrencyPairs(authority string, currencyPairIDs []string) MsgRemoveCurrencyPairs {
	return MsgRemoveCurrencyPairs{
		Authority:       authority,
		CurrencyPairIds: currencyPairIDs,
	}
}

// ValidateBasic determines whether the information in the message is valid, specifically
// whether the authority is a valid acc-address, and that each CurrencyPairID in the message is formatted correctly.
func (m *MsgRemoveCurrencyPairs) ValidateBasic() error {
	// validate authority address
	_, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		return err
	}

	// check that each CurrencyPairID is correctly formatted
	for _, id := range m.CurrencyPairIds {
		if _, err := connecttypes.CurrencyPairFromString(id); err != nil {
			return err
		}
	}

	return nil
}
