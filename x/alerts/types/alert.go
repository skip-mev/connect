package types

import (
	"fmt"
	"time"

	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
)

const AlertUIDLen = 20

// ConclusionStatus wraps the status of an alert.
type ConclusionStatus uint64

const (
	// Unconcluded is the default status of an alert (no conclusion has been submitted for it yet).
	Unconcluded ConclusionStatus = iota
	// Concluded is the status of an alert that has been concluded.
	Concluded
)

// String implements fmt.Stringer.
func (cs ConclusionStatus) String() string {
	switch cs {
	case Unconcluded:
		return "Unconcluded"
	case Concluded:
		return "Concluded"
	default:
		return "unknown"
	}
}

// NewAlert returns a new Alert.
func NewAlert(height uint64, signer sdk.AccAddress, cp slinkytypes.CurrencyPair) Alert {
	return Alert{
		Height:       height,
		Signer:       signer.String(),
		CurrencyPair: cp,
	}
}

// ValidateBasic performs stateless validation on the Claim, i.e. checks that the signer and the currency-pair are valid.
func (a *Alert) ValidateBasic() error {
	// validate the currency-pair
	if err := a.CurrencyPair.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid alert: %w", err)
	}

	// validate the signer
	if _, err := sdk.AccAddressFromBech32(a.Signer); err != nil {
		return fmt.Errorf("invalid alert: %w", err)
	}

	return nil
}

// UID returns the Unique Identifier of this Claim, this is how the BondAddress is defined and
// this will be used to derive the key under which this Claim will be stored. This method appends
// the Height, signer, currency-pair strings into a byte-array, and returns the first 20-bytes of
// the hash of that array.
func (a *Alert) UID() []byte {
	heightBz := []byte(fmt.Sprintf("%d", a.Height))
	signerBz := []byte(a.Signer)
	currencyPairBz := []byte(a.CurrencyPair.String())
	return tmhash.SumTruncated(append(append(heightBz, signerBz...), currencyPairBz...))
}

// NewAlertStatus returns a new AlertStatus.
func NewAlertStatus(submissionHeight uint64, purgeHeight uint64, blockTimestamp time.Time, status ConclusionStatus) AlertStatus {
	return AlertStatus{
		SubmissionHeight:    submissionHeight,
		PurgeHeight:         purgeHeight,
		ConclusionStatus:    uint64(status),
		SubmissionTimestamp: uint64(blockTimestamp.UTC().Unix()),
	}
}

// ValidateBasic performs a basic validation of the ConclusionStatus, i.e. that the submissionHeight
// is before the purgeHeight, and that the status is either Unconcluded or Concluded.
func (a *AlertStatus) ValidateBasic() error {
	if a.SubmissionHeight >= a.PurgeHeight {
		return fmt.Errorf("invalid alert status: submission height must be before purge height")
	}

	if ConclusionStatus(a.ConclusionStatus) != Unconcluded && ConclusionStatus(a.ConclusionStatus) != Concluded {
		return fmt.Errorf("invalid alert status: status must be either Unconcluded or Concluded")
	}

	return nil
}

// NewAlertWithStatus returns a new AlertWithStatus.
func NewAlertWithStatus(alert Alert, status AlertStatus) AlertWithStatus {
	return AlertWithStatus{
		Alert:  alert,
		Status: status,
	}
}

// ValidateBasic validates that both the alert status and the alert are valid.
func (aws *AlertWithStatus) ValidateBasic() error {
	if err := aws.Alert.ValidateBasic(); err != nil {
		return err
	}

	return aws.Status.ValidateBasic()
}
