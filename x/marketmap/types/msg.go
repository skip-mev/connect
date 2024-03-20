package types

import (
	"fmt"

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
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return err
	}

	for _, market := range m.CreateMarkets {
		if err := market.Ticker.ValidateBasic(); err != nil {
			return err
		}

		for _, path := range market.Paths.Paths {
			if err := path.ValidateBasic(); err != nil {
				return err
			}
		}

		if uint64(len(market.Providers.Providers)) < market.Ticker.MinProviderCount {
			return fmt.Errorf("this ticker must have at least %d providers; got %d",
				market.Ticker.MinProviderCount,
				len(market.Providers.Providers),
			)
		}

		seenProviders := make(map[string]struct{})
		for _, provider := range market.Providers.Providers {
			if err := provider.ValidateBasic(); err != nil {
				return err
			}

			// check for duplicate providers
			if _, seen := seenProviders[provider.Name]; seen {
				return fmt.Errorf("duplicate provider found: %s", provider.Name)
			}
			seenProviders[provider.Name] = struct{}{}
		}
	}

	return nil
}

// ValidateBasic determines whether the information in the message is formatted correctly, specifically
// whether the signer is a valid acc-address.
func (m *MsgUpdateMarkets) ValidateBasic() error {
	// validate signer address
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return err
	}

	for _, market := range m.UpdateMarkets {
		if err := market.Ticker.ValidateBasic(); err != nil {
			return err
		}

		for _, path := range market.Paths.Paths {
			if err := path.ValidateBasic(); err != nil {
				return err
			}
		}

		if uint64(len(market.Providers.Providers)) < market.Ticker.MinProviderCount {
			return fmt.Errorf("this ticker must have at least %d providers; got %d",
				market.Ticker.MinProviderCount,
				len(market.Providers.Providers),
			)
		}

		seenProviders := make(map[string]struct{})
		for _, provider := range market.Providers.Providers {
			if err := provider.ValidateBasic(); err != nil {
				return err
			}

			// check for duplicate providers
			if _, seen := seenProviders[provider.Name]; seen {
				return fmt.Errorf("duplicate provider found: %s", provider.Name)
			}
			seenProviders[provider.Name] = struct{}{}
		}
	}

	return nil
}
