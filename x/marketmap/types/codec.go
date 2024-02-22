package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
)

// RegisterLegacyAminoCodec registers the necessary x/marketmap interfaces (messages) on the
// cdc. These types are used for amino serialization.
func RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {
}

// RegisterInterfaces registers the x/marketmap messages + message service w/ the InterfaceRegistry (registry).
func RegisterInterfaces(_ types.InterfaceRegistry) {
}
