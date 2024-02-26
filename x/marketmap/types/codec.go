package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/marketmap interfaces (messages) on the
// cdc. These types are used for amino serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// register the msg-types
	legacy.RegisterAminoMsg(cdc, &MsgUpdateMarketMap{}, "slinky/x/marketmap/MsgUpdateMarketMap")
	legacy.RegisterAminoMsg(cdc, &MsgParams{}, "slinky/x/marketmap/MsgParams")
}

// RegisterInterfaces registers the x/marketmap messages + message service w/ the InterfaceRegistry (registry).
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	// register the implementations of Msg-type
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateMarketMap{},
		&MsgParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
