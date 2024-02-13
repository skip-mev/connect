package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateMarket{}

// GetSigners gets the address that must sign this message.
func (m *MsgCreateMarket) GetSigners() []sdk.AccAddress {
	// convert from string to acc address
	addr, _ := sdk.AccAddressFromBech32(m.Signer)
	return []sdk.AccAddress{addr}
}

// ValidateBasic determines whether the information in the message is formatted correctly, specifically
// whether the signer is a valid acc-address.
func (m *MsgCreateMarket) ValidateBasic() error {
	// validate signer address
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return err
	}

	if err := m.Ticker.ValidateBasic(); err != nil {
		return err
	}

	if len(m.Paths) == 0 {
		return fmt.Errorf("at least one path is required for a ticker to be calculated")
	}

	for _, path := range m.Paths {
		if err := path.ValidateBasic(); err != nil {
			return err
		}
	}

	if uint64(len(m.ProvidersToOffChainTickers)) < m.Ticker.MinProviderCount {
		return fmt.Errorf("this ticker must have at least %d providers; got %d",
			m.Ticker.MinProviderCount,
			len(m.ProvidersToOffChainTickers),
		)
	}

	seenProviders := make(map[string]struct{})
	for providerName, offChainTicker := range m.ProvidersToOffChainTickers {
		// check for duplicate providers
		if _, seen := seenProviders[providerName]; seen {
			return fmt.Errorf("duplicate provider found: %s", providerName)
		}
		seenProviders[providerName] = struct{}{}

		if offChainTicker == "" {
			return fmt.Errorf("got empty off chain ticker for provider %s", providerName)
		}
	}

	return nil
}
