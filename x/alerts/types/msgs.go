package types

import (
	"fmt"

	"cosmossdk.io/x/tx/signing"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/proto"

	alertv1 "github.com/skip-mev/connect/v2/api/slinky/alerts/v1"
)

var (
	_ sdk.Msg = &MsgAlert{}
	_ sdk.Msg = &MsgConclusion{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// ProvideMsgAlertGetSigners provides the CustomGetSigners method for the MsgAlert.
func ProvideMsgAlertGetSigners() signing.CustomGetSigner {
	return signing.CustomGetSigner{
		MsgType: proto.MessageName(&alertv1.MsgAlert{}),
		Fn: func(m proto.Message) ([][]byte, error) {
			msg := m.(*alertv1.MsgAlert)
			signer, err := sdk.AccAddressFromBech32(msg.Alert.Signer)
			if err != nil {
				return nil, err
			}
			return [][]byte{signer}, nil
		},
	}
}

// NewMsgAlert creates a new MsgAlert, given the alert details.
func NewMsgAlert(a Alert) *MsgAlert {
	return &MsgAlert{
		Alert: a,
	}
}

// ValidateBasic performs basic validation of the MsgAlert fields, i.e. on the underlying Alert.
func (msg *MsgAlert) ValidateBasic() error {
	return msg.Alert.ValidateBasic()
}

// NewMsgConclusion creates a new MsgConclusion, given the conclusion, and the
// conclusion signer.
func NewMsgConclusion(c Conclusion, signer sdk.AccAddress) *MsgConclusion {
	msg := &MsgConclusion{
		Signer: signer.String(),
	}

	var err error
	if c != nil {
		// marshal to any
		msg.Conclusion, err = codectypes.NewAnyWithValue(c)
		if err != nil {
			return nil
		}
	}

	return msg
}

// ValidateBasic performs basic validation of the MsgConclusion fields, i.e. that the signer address
// is valid, and the alert is non-nil.
func (msg *MsgConclusion) ValidateBasic() error {
	// check signer validity
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return err
	}

	// check conclusion validity
	if msg.Conclusion == nil {
		return fmt.Errorf("conclusion cannot be nil")
	}

	// unmarshal the any into a conclusion
	var c Conclusion
	if err := pc.UnpackAny(msg.Conclusion, &c); err != nil {
		return err
	}

	return c.ValidateBasic()
}

// NewMsgUpdateParams creates a new MsgUpdateParams, given the new params, and the authority address
// that is allowed to update the params.
func NewMsgUpdateParams(params Params, authority sdk.AccAddress) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority.String(),
		Params:    params,
	}
}

// ValidateBasic checks that the params in the msg are valid, and that the authority address is valid
// if either check fails the method will error.
func (m *MsgUpdateParams) ValidateBasic() error {
	// check authority address
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return err
	}

	// check params
	return m.Params.Validate()
}
