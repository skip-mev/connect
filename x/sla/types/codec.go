package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/sla interfaces (messages) on the
// provided LegacyAmino codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgAddSLAs{}, "sla/MsgAddSLAs")
	legacy.RegisterAminoMsg(cdc, &MsgRemoveSLAs{}, "sla/MsgRemoveSLAs")
	legacy.RegisterAminoMsg(cdc, &MsgParams{}, "sla/MsgParams")
}

// RegisterInterfaces registers the x/sla interfaces (messages + msg server) on the
// provided InterfaceRegistry.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddSLAs{},
		&MsgRemoveSLAs{},
		&MsgParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
