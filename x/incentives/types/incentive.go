package types

import (
	"github.com/cosmos/gogoproto/proto"
)

type (
	// Incentive is the interface contract that must be fulfilled to allow for
	// slashing/rewarding arbitrary events. To add a new incentive, implement
	// this interface along with the corresponding strategy callback function.
	// Each incentive type must be registered with its corresponding strategy
	// function in the keeper.
	Incentive interface {
		proto.Message

		// ValidateBasic does a stateless check on the incentive.
		ValidateBasic() error

		// Type returns the incentive type.
		Type() string

		// Marshall the incentive into bytes.
		Marshal() ([]byte, error)

		// Unmarshal the incentive from bytes.
		Unmarshal([]byte) error

		// Copy the incentive.
		Copy() Incentive
	}
)
