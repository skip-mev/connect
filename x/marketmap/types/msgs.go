package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var _ sdk.Msg = &MsgCreateMarket{}

// ValidateBasic determines whether the information in the message is formatted correctly, specifically
// whether the signer is a valid acc-address.
func (m *MsgCreateMarket) ValidateBasic() error {
	// validate signer address
	_, err := sdk.AccAddressFromBech32(m.Signer)

	return err
}
