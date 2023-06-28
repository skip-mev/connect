package abci

import (
	"bytes"
	"fmt"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtcrypto "github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protoio "github.com/cosmos/gogoproto/io"
	"github.com/cosmos/gogoproto/proto"

	cryptoenc "github.com/cometbft/cometbft/crypto/encoding"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
)

// ValidateVoteExtensions defines a helper function for verifying vote extension
// signatures that may be passed or manually injected into a block proposal from
// a proposer in ProcessProposal. It returns an error if any signature is invalid
// or if unexpected vote extensions and/or signatures are found or less than 2/3
// power is received.
func ValidateVoteExtensions( // TODO(nikhil): use sdk's ValidateVoteExtensions method once it's merged
	ctx sdk.Context,
	valStore baseapp.ValidatorStore,
	currentHeight int64,
	chainID string,
	extCommit abci.ExtendedCommitInfo,
) error {
	cp := ctx.ConsensusParams()
	extsEnabled := cp.Abci != nil && cp.Abci.VoteExtensionsEnableHeight < currentHeight

	// skip first block + any block before VEs were enabled
	if currentHeight <= cp.Abci.VoteExtensionsEnableHeight {
		return nil
	}

	marshalDelimitedFn := func(msg proto.Message) ([]byte, error) {
		var buf bytes.Buffer
		if err := protoio.NewDelimitedWriter(&buf).WriteMsg(msg); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}

	sumVP := math.ZeroInt()
	for _, vote := range extCommit.Votes {
		if !extsEnabled {
			if len(vote.VoteExtension) > 0 {
				return fmt.Errorf("vote extensions disabled; received non-empty vote extension at height %d", currentHeight)
			}
			if len(vote.ExtensionSignature) > 0 {
				return fmt.Errorf("vote extensions disabled; received non-empty vote extension signature at height %d", currentHeight)
			}

			continue
		}

		if len(vote.ExtensionSignature) == 0 {
			return fmt.Errorf("vote extensions enabled; received empty vote extension signature at height %d", currentHeight)
		}

		valConsAddr := cmtcrypto.Address(vote.Validator.Address)

		validator, err := valStore.GetValidatorByConsAddr(ctx, valConsAddr)
		if err != nil {
			return fmt.Errorf("failed to get validator %X: %w", valConsAddr, err)
		}
		if validator == nil {
			return fmt.Errorf("validator %X not found", valConsAddr)
		}

		cmtPubKeyProto, err := validator.CmtConsPublicKey()
		if err != nil {
			return fmt.Errorf("failed to get validator %X public key: %w", valConsAddr, err)
		}

		cmtPubKey, err := cryptoenc.PubKeyFromProto(cmtPubKeyProto)
		if err != nil {
			return fmt.Errorf("failed to convert validator %X public key: %w", valConsAddr, err)
		}

		cve := cmtproto.CanonicalVoteExtension{
			Extension: vote.VoteExtension,
			Height:    currentHeight - 1, // the vote extension was signed in the previous height
			Round:     int64(extCommit.Round),
			ChainId:   chainID,
		}
		ctx.Logger().Info("Validating vote extension", "vote_extension", cve)

		extSignBytes, err := marshalDelimitedFn(&cve)
		if err != nil {
			return fmt.Errorf("failed to encode CanonicalVoteExtension: %w", err)
		}

		if !cmtPubKey.VerifySignature(extSignBytes, vote.ExtensionSignature) {
			return fmt.Errorf("failed to verify validator %X vote extension signature", valConsAddr)
		}

		sumVP = sumVP.Add(validator.BondedTokens())
	}

	// Ensure we have at least 2/3 voting power that submitted valid vote
	// extensions.
	totalVP := valStore.TotalBondedTokens(ctx)
	percentSubmitted := math.LegacyNewDecFromInt(sumVP).Quo(math.LegacyNewDecFromInt(totalVP))
	if percentSubmitted.LT(baseapp.VoteExtensionThreshold) {
		return fmt.Errorf("insufficient cumulative voting power received to verify vote extensions; got: %s, expected: >=%s", percentSubmitted, baseapp.VoteExtensionThreshold)
	}

	return nil
}

// NoOpValidateVoteExtensions is a no-op validation method (purely used for testing)
func NoOpValidateVoteExtensions(
	_ sdk.Context,
	_ int64,
	_ abci.ExtendedCommitInfo,
) error {
	return nil
}
