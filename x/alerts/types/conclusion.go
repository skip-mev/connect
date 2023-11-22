package types

import (
	cmtabci "github.com/cometbft/cometbft/abci/types"
	"github.com/gogo/protobuf/proto"
)

// Conclusion defines the basic meta-data necessary for the alerts module to resolve an alert. At a minimum, this
// message must contain the OracleData referenced in the conclusion, the Alert that the Conclusion corresponds to,
// metadata necessary for verification (i.e signatures), and the status of the conclusion itself.
type Conclusion interface {
	proto.Message

	// performs stateless validation on the Claim
	ValidateBasic() error

	// Marshal the Claim to bytes (this is what will be stored in state)
	Marshal() ([]byte, error)

	// Unmarshal unmarshals the claim bytes into the receiver
	Unmarshal([]byte) error

	// Verify verifies the conclusion, given a ConclusionVerificationParams
	Verify(ConclusionVerificationParams) error

	// GetExtendedCommitInfo returns the ExtendedCommitInfo referenced in the conclusion.
	GetExtendedCommitInfo() cmtabci.ExtendedCommitInfo

	// Alert returns the Alert that the Conclusion corresponds to.
	GetAlert() Alert

	// Status returns the status of the conclusion.
	GetStatus() bool

	// PriceBounds returns the price-bounds of the conclusion.
	GetPriceBound() PriceBound

	// GetCurrencyPairID returns the ID of the CurrencyPair for which the alert is filed
	GetCurrencyPairID() uint64
}

// ConclusionVerificationParams defines the parameters necessary to verify a conclusion.
type ConclusionVerificationParams interface {
	proto.Message

	// performs stateless validation on the Claim
	ValidateBasic() error

	// Marshal the Claim to bytes (this is what will be stored in state)
	Marshal() ([]byte, error)

	// Unmarshal unmarshals the claim bytes into the receiver
	Unmarshal([]byte) error
}
