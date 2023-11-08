package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var pc *codec.ProtoCodec

func init() {
	ir := codectypes.NewInterfaceRegistry()
	pc = codec.NewProtoCodec(ir)

	// register crypto types
	cryptocodec.RegisterInterfaces(ir)
}

// RegisterLegacyAminoCodec registers the necessary x/authz interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// register the conclusion interfaces / MultiSigConclusion implementation
	cdc.RegisterInterface((*Conclusion)(nil), nil)
	cdc.RegisterConcrete(&MultiSigConclusion{}, "slinky/x/alerts/Conclusion", nil)

	// register the conclusion verification params interfaces / MultiSigConclusionVerificationParams implementation
	cdc.RegisterInterface((*ConclusionVerificationParams)(nil), nil)
	cdc.RegisterConcrete(&MultiSigConclusionVerificationParams{}, "slinky/x/alerts/ConclusionVerificationParams", nil)

	// register the msg-types
	legacy.RegisterAminoMsg(cdc, &MsgAlert{}, "slinky/x/alerts/MsgAlert")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "slinky/x/alerts/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgConclusion{}, "slinky/x/alerts/MsgConclusion")
}

// RegisterInterfaces registers the x/alerts messages + message service w/ the InterfaceRegistry (registry).
//
// This method will also update the internal package's codec so that any subsequent attempts in the package to unmarshal
// anys will reference the registry passed here.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	// register the ConclusionVerificationParams interface + implementations
	registry.RegisterInterface(
		"slinky.alerts.v1.ConclusionVerificationParams",
		(*ConclusionVerificationParams)(nil),
		&MultiSigConclusionVerificationParams{},
	)

	// register the Conclusion interface + implementations
	registry.RegisterInterface(
		"slinky.alerts.v1.Conclusion",
		(*Conclusion)(nil),
		&MultiSigConclusion{},
	)

	// register the alert Msg-type
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgAlert{},
		&MsgUpdateParams{},
		&MsgConclusion{},
	)

	// update the package's codec
	pc = codec.NewProtoCodec(registry)
}
