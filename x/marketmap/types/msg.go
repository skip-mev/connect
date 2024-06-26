package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCreateMarkets{}
	_ sdk.Msg = &MsgUpdateMarkets{}
	_ sdk.Msg = &MsgParams{}
	_ sdk.Msg = &MsgRemoveMarketAuthorities{}
	_ sdk.Msg = &MsgUpsertMarkets{}
)

// ValidateBasic asserts that the authority address in the upsert-markets message is formatted correctly.
// If also verifies that all markets w/in the message are valid, if no markets are present it returns an error.
func (m *MsgUpsertMarkets) ValidateBasic() error {
	// validate signer address
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return err
	}

	if len(m.Markets) == 0 {
		return fmt.Errorf("no markets to upsert")
	}

	seenTickers := make(map[string]struct{})
	for _, market := range m.Markets {
		ticker := market.Ticker.CurrencyPair.String()

		if _, seen := seenTickers[ticker]; seen {
			return fmt.Errorf("duplicate ticker: %s", ticker)
		}

		if err := market.ValidateBasic(); err != nil {
			return err
		}

		seenTickers[ticker] = struct{}{}
	}

	return nil
}

// ValidateBasic determines whether the information in the message is formatted correctly, specifically
// whether the signer is a valid acc-address.
func (m *MsgCreateMarkets) ValidateBasic() error {
	// validate signer address
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return err
	}

	if len(m.CreateMarkets) == 0 {
		return fmt.Errorf("no markets to create")
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

	if len(m.UpdateMarkets) == 0 {
		return fmt.Errorf("no markets to update")
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

// ValidateBasic determines whether the information in the message is formatted correctly, specifically
// whether the signer is a valid acc-address.
func (m *MsgRemoveMarketAuthorities) ValidateBasic() error {
	// validate signer address
	if _, err := sdk.AccAddressFromBech32(m.Admin); err != nil {
		return err
	}

	if len(m.RemoveAddresses) == 0 {
		return fmt.Errorf("addresses to remove cannot be nil")
	}

	seenAuthorities := make(map[string]struct{}, len(m.RemoveAddresses))
	for _, authority := range m.RemoveAddresses {
		if _, seen := seenAuthorities[authority]; seen {
			return fmt.Errorf("duplicate address %s found", authority)
		}

		if _, err := sdk.AccAddressFromBech32(authority); err != nil {
			return fmt.Errorf("invalid market authority string: %w", err)
		}

		seenAuthorities[authority] = struct{}{}
	}

	return nil
}
