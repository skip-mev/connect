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
	legacy.RegisterAminoMsg(cdc, &MsgCreateMarkets{}, "slinky/x/marketmap/MsgCreateMarkets")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateMarkets{}, "slinky/x/marketmap/MsgUpdateMarkets")
	legacy.RegisterAminoMsg(cdc, &MsgParams{}, "slinky/x/marketmap/MsgParams")
	legacy.RegisterAminoMsg(cdc, &MsgUpsertMarkets{}, "slinky/x/marketmap/MsgUpsertMarkets")
	legacy.RegisterAminoMsg(cdc, &MsgRemoveMarkets{}, "slinky/x/marketmap/MsgRemoveMarkets")
}

// RegisterInterfaces registers the x/marketmap messages + message service w/ the InterfaceRegistry (registry).
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	// register the implementations of Msg-type
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateMarkets{},
		&MsgUpdateMarkets{},
		&MsgParams{},
		&MsgUpsertMarkets{},
		&MsgRemoveMarkets{},
		&MsgRemoveMarketAuthorities{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
