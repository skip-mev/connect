package types

import (
	"fmt"

	cmthash "github.com/cometbft/cometbft/crypto/tmhash"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256r1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ Conclusion                   = &MultiSigConclusion{}
	_ ConclusionVerificationParams = &MultiSigConclusionVerificationParams{}
)

const (
	Secp256k1Type = "secp256k1"
	Secp256r1Type = "secp256r1"
	Ed25519Type   = "ed25519"
)

// NewMultiSigConclusion creates the parameters necessary for verification of a MultiSigConclusion. This method will
// derive the addresses corresponding to the provided public-keys, and will store the address -> public key mapping,
// this method will fail if no public-keys are given, or if duplicates are provided.
func NewMultiSigVerificationParams(pks []cryptotypes.PubKey) (ConclusionVerificationParams, error) {
	signers := make(map[string]struct{})

	// expect there to be at least 1 signer
	if len(pks) == 0 {
		return nil, fmt.Errorf("no signers provided")
	}

	pbks := make([]*codectypes.Any, 0)

	for _, pk := range pks {
		// get signer address from pubkey
		signer := sdk.AccAddress(pk.Address()).String()

		// check for duplicate signers
		if _, ok := signers[signer]; ok {
			return nil, fmt.Errorf("duplicate signer: %s", signer)
		}

		// pack the pubkey into an any
		pbk, err := codectypes.NewAnyWithValue(pk)
		if err != nil {
			return nil, err
		}

		pbks = append(pbks, pbk)

		// store signer in map
		signers[signer] = struct{}{}
	}

	return &MultiSigConclusionVerificationParams{
		Signers: pbks,
	}, nil
}

// ValidateBasic checks that all signers correspond to their pubkeys, and that the signer addresses are valid bech32 addresses.
// It also checks that there is at least 1 signer. This method also validates that the signers are not duplicated.
//
// NOTICE: the public-keys given must be of type secp2561k1, secp256r1, or ed25519.
func (params MultiSigConclusionVerificationParams) ValidateBasic() error {
	// check that there is at least 1 signer
	if len(params.Signers) == 0 {
		return fmt.Errorf("no signers provided")
	}

	signers := make(map[string]struct{})

	// check validity of pubkeys
	for _, pk := range params.Signers {
		var pkv cryptotypes.PubKey
		if err := pc.UnpackAny(pk, &pkv); err != nil {
			return err
		}

		// check that the pubkey is non-nil
		if pkv == nil {
			return fmt.Errorf("nil pubkey")
		}

		// check that the pubkey is of a valid type
		if err := validatePkType(pkv); err != nil {
			return err
		}

		// get signer address from pubkey
		signer := sdk.AccAddress(pkv.Address()).String()
		if _, ok := signers[signer]; ok {
			return fmt.Errorf("duplicate signer: %s", signer)
		}

		// store signer in map
		signers[signer] = struct{}{}
	}

	return nil
}

func validatePkType(pkv cryptotypes.PubKey) error {
	switch {
	case pkv.Type() == Secp256k1Type:
		if _, ok := pkv.(*secp256k1.PubKey); !ok {
			return fmt.Errorf("invalid secp256k1 pubkey")
		}
	case pkv.Type() == Secp256r1Type:
		if _, ok := pkv.(*secp256r1.PubKey); !ok {
			return fmt.Errorf("invalid secp256r1 pubkey")
		}
	case pkv.Type() == Ed25519Type:
		if _, ok := pkv.(*ed25519.PubKey); !ok {
			return fmt.Errorf("invalid ed25519 pubkey")
		}
	default:
		return fmt.Errorf("invalid pubkey type: %s", pkv.Type())
	}

	return nil
}

// ValidateBasic validates the Conclusion. Specifically, it validates that the signers are valid bech32 addresses, and that the
// sub-fields are non-nil.
func (c MultiSigConclusion) ValidateBasic() error {
	// check that the alert is valid
	if err := c.Alert.ValidateBasic(); err != nil {
		return err
	}

	// check that the price-bound is valid
	if err := c.PriceBound.ValidateBasic(); err != nil {
		return err
	}

	if len(c.Signatures) == 0 {
		return fmt.Errorf("no signatures provided")
	}

	// check that each signer is a vlaid bech32 address
	for _, signature := range c.Signatures {
		if _, err := sdk.AccAddressFromBech32(signature.Signer); err != nil {
			return err
		}
	}

	return nil
}

// Verify verifies the conclusion. Specifically, it verifies that the signatures are valid, and that the oracle-data is valid.
func (c MultiSigConclusion) Verify(params ConclusionVerificationParams) error {
	// check that the params are valid
	if err := params.ValidateBasic(); err != nil {
		return err
	}

	// check that the conclusion is valid
	if err := c.ValidateBasic(); err != nil {
		return err
	}

	// assert the type of the params
	multiSigParams, ok := params.(*MultiSigConclusionVerificationParams)
	if !ok {
		return fmt.Errorf("invalid params type: %T", params)
	}

	sigBytes, err := c.SignBytes()
	if err != nil {
		return err
	}

	signatures := signaturesToSignersMap(c.Signatures)

	// verify the signatures of the Conclusion data
	for _, pkb := range multiSigParams.Signers {
		// get the signer address associated w/ pkb
		var pkv cryptotypes.PubKey
		if err := pc.UnpackAny(pkb, &pkv); err != nil {
			return err
		}

		signer := sdk.AccAddress(pkv.Address()).String()
		sig, ok := signatures[signer]
		if !ok {
			return fmt.Errorf("no signature provided for signer: %s", signer)
		}

		// verify the signature
		if !pkv.VerifySignature(sigBytes, sig) {
			return fmt.Errorf("signature verification failed for signer: %s", signer)
		}

	}

	return nil
}

func signaturesToSignersMap(signatures []Signature) map[string][]byte {
	signaturesMap := make(map[string][]byte)

	for _, sig := range signatures {
		signaturesMap[sig.Signer] = sig.Signature
	}

	return signaturesMap
}

// SignBytes returns the bytes that should be signed by the signers. I.e the 20-byte truncated hash of the marshalled oracle-data,
// the alert UID, the price-bound, and the status.
func (c MultiSigConclusion) SignBytes() ([]byte, error) {
	bz := make([]byte, 0)

	// append oracle-data
	oracleData, err := c.ExtendedCommitInfo.Marshal()
	if err != nil {
		return nil, err
	}

	bz = append(bz, oracleData...)

	// append alert
	bz = append(bz, c.Alert.UID()...)

	// append price-bound
	priceBoundBz, err := c.PriceBound.Marshal()
	if err != nil {
		return nil, err
	}
	bz = append(bz, priceBoundBz...)

	// append status
	switch c.Status {
	case true:
		bz = append(bz, 1)
	case false:
		bz = append(bz, 0)
	}
	return cmthash.SumTruncated(bz), nil
}
