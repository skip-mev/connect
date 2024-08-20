package strategies

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"

	incentivetypes "github.com/skip-mev/connect/v2/x/incentives/types"
)

// RegisterLegacyAminoCodec registers the necessary x/incentives interfaces (messages) on the
// cdc. These types are used for amino serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// register the ValidatorAlertIncentive
	legacy.RegisterAminoMsg(cdc, &ValidatorAlertIncentive{}, "slinky/x/alerts/ValidatorAlertIncentive")
}

// RegisterInterfaces registers the x/incentives messages + message service w/ the InterfaceRegistry (registry).
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*incentivetypes.Incentive)(nil),
		&ValidatorAlertIncentive{},
	)
}
