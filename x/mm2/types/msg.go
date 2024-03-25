package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCreateMarkets{}
	_ sdk.Msg = &MsgUpdateMarkets{}
	_ sdk.Msg = &MsgParams{}
)

// ValidateBasic determines whether the information in the message is formatted correctly, specifically
// whether the signer is a valid acc-address.
func (m *MsgCreateMarkets) ValidateBasic() error {
	// validate signer address
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return err
	}

	for _, market := range m.CreateMarkets {
		if err := market.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}

// ValidateBasic determines whether the information in the message is formatted correctly, specifically
// whether the signer is a valid acc-address.
func (m *MsgUpdateMarkets) ValidateBasic() error {
	// validate signer address
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return err
	}

	for _, market := range m.UpdateMarkets {
		if err := market.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}

// ValidateBasic determines whether the information in the message is formatted correctly, specifically
// whether the signer is a valid acc-address.
func (m *MsgParams) ValidateBasic() error {
	// validate signer address
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return err
	}

	return m.Params.ValidateBasic()
}
