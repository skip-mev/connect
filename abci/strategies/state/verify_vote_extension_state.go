package state

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetSlinkyAppStatePruningParams sets the minimum block retention height for the application, and the
// underlying multi-store. This is used to ensure that the application can retrieve the state necessary
// to verify vote-extensions.
func SetSlinkyAppStatePruningParams() func(*baseapp.BaseApp) {
	return func(ba *baseapp.BaseApp) {
		// check what the app's CommitMultiStore's KeepRecent is
		if pruning := ba.CommitMultiStore().GetPruning(); pruning.KeepRecent < VoteExtensionVerificationHeightOffset {
			pruning.KeepRecent = VoteExtensionVerificationHeightOffset
			ba.CommitMultiStore().SetPruning(pruning)
		}
	}
}

const VoteExtensionVerificationHeightOffset = 2

// Application is the expected interface of an SDK based application.
//
//go:generate mockery --name Application --filename application.go
type Application interface {
	GetBlockRetentionHeight(commitHeight int64) int64
	CommitMultiStore() storetypes.CommitMultiStore
}

// AppStates is an interface used to retrieve application states for the purposes of
// Slinky ABCI methods
//
//go:generate mockery --name AppState --filename app_state.go
type AppState interface {
	// VerifyVoteExtensionState is used to get the state against which the vote-extensions
	// of height h - 1 are verified (h - 2 state). This is used to ensure parity between Process + PrepareProposal
	// and VerifyVoteExtension.
	VerifyVoteExtensionState(ctx sdk.Context) (sdk.Context, error)
}

// NewBaseAppState returns a new instance of the baseAppState. That returns a cached state for verifying vote-extensions.
func NewBaseAppState(app Application) AppState {
	return baseAppState{
		app: app,
	}
}

// NewNoopAppState returns a new instance of the noopAppState. That returns the same context for verifying vote-extensions.
func NewNoopAppState() AppState {
	return noopAppState{}
}

type baseAppState struct {
	app Application
}

func (b baseAppState) VerifyVoteExtensionState(ctx sdk.Context) (sdk.Context, error) {
	// check that the app's retention height is sufficient to retrieve the state for
	// verifying vote-extensions.
	retentionHeight := b.app.GetBlockRetentionHeight(ctx.BlockHeight())
	if ctx.BlockHeight()-retentionHeight < VoteExtensionVerificationHeightOffset {
		return ctx, fmt.Errorf("insufficient retention height for verifying vote-extensions: required: %d, available: %d", ctx.BlockHeight()-VoteExtensionVerificationHeightOffset, retentionHeight)
	}

	// retrieve the state at height h - 2
	multiStore, err := b.app.CommitMultiStore().CacheMultiStoreWithVersion(ctx.BlockHeight() - VoteExtensionVerificationHeightOffset)
	if err != nil {
		return ctx, fmt.Errorf("failed to retrieve state at height %d: %w", ctx.BlockHeight()-VoteExtensionVerificationHeightOffset, err)
	}

	return ctx.WithMultiStore(multiStore), nil
}

type noopAppState struct{}

func (n noopAppState) VerifyVoteExtensionState(ctx sdk.Context) (sdk.Context, error) {
	return ctx, nil
}
