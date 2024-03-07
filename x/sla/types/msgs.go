package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgAddSLAs{}
	_ sdk.Msg = &MsgRemoveSLAs{}
	_ sdk.Msg = &MsgParams{}
)

// NewMsgAddSLAs returns a new message from a set of SLAs and an authority address.
func NewMsgAddSLAs(authority string, slas []PriceFeedSLA) MsgAddSLAs {
	return MsgAddSLAs{
		Authority: authority,
		SLAs:      slas,
	}
}

// ValidateBasic determines whether the information in the message is formatted correctly, specifically
// whether the authority is a valid acc-address, and that each SLA in the message is formatted correctly.
func (m *MsgAddSLAs) ValidateBasic() error {
	// validate authority address
	_, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		return err
	}

	// validate SLAs
	seen := make(map[string]struct{})
	for _, sla := range m.SLAs {
		if _, ok := seen[sla.ID]; ok {
			return fmt.Errorf("duplicate price feed sla id %s", sla.ID)
		}

		if err := sla.ValidateBasic(); err != nil {
			return err
		}

		seen[sla.ID] = struct{}{}
	}

	return nil
}

// NewMsgRemoveSLAs returns a new message to remove a set of SLAs from the x/sla module's state.
func NewMsgRemoveSLAs(authority string, slaIDs []string) MsgRemoveSLAs {
	return MsgRemoveSLAs{
		Authority: authority,
		IDs:       slaIDs,
	}
}

// ValidateBasic determines whether the information in the message is formatted correctly, specifically
// whether the authority is a valid acc-address, and that each SLA ID in the message is not empty.
func (m *MsgRemoveSLAs) ValidateBasic() error {
	// validate authority address
	_, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		return err
	}

	// validate SLA IDs
	seen := make(map[string]struct{})
	for _, id := range m.IDs {
		if _, ok := seen[id]; ok {
			return fmt.Errorf("duplicate price feed sla id %s", id)
		}

		if len(id) == 0 {
			return fmt.Errorf("sla id cannot be empty")
		}

		seen[id] = struct{}{}
	}

	return nil
}

// NewMsgParams returns a new message to update the x/sla module's parameters.
func NewMsgParams(authority string, params Params) MsgParams {
	return MsgParams{
		Authority: authority,
		Params:    params,
	}
}

// ValidateBasic determines whether the information in the message is formatted correctly, specifically
// whether the authority is a valid acc-address.
func (m *MsgParams) ValidateBasic() error {
	// validate authority address
	_, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		return err
	}

	return nil
}
