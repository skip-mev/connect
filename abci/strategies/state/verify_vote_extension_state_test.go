package state_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/log"
	pruningtypes "cosmossdk.io/store/pruning/types"
	storetypes "cosmossdk.io/store/types"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/abci/strategies/state"
	statemocks "github.com/skip-mev/slinky/abci/strategies/state/mocks"
	"github.com/skip-mev/slinky/tests/simapp"
)

// test SetSlinkyAppStatePruningParams.
func TestSetSlinkyAppStatePruningParams(t *testing.T) {
	t.Run("test that KeepRecent is updated if it is too low", func(t *testing.T) {
		app := simapp.NewSimApp(
			log.NewNopLogger(),
			dbm.NewMemDB(),
			nil,
			true,
			simtestutil.EmptyAppOptions{},
			baseapp.SetPruning(pruningtypes.PruningOptions{
				KeepRecent: 1,
				Interval:   10,
				Strategy:   pruningtypes.PruningNothing,
			}),
			state.SetSlinkyAppStatePruningParams(),
		)

		require.Equal(t, app.CommitMultiStore().GetPruning().KeepRecent, uint64(state.VoteExtensionVerificationHeightOffset))
	})

	t.Run("test that KeepRecent is not updated if it is already high enough", func(t *testing.T) {
		app := simapp.NewSimApp(
			log.NewNopLogger(),
			dbm.NewMemDB(),
			nil,
			true,
			simtestutil.EmptyAppOptions{},
			baseapp.SetPruning(pruningtypes.PruningOptions{
				KeepRecent: 10,
				Interval:   10,
				Strategy:   pruningtypes.PruningNothing,
			}),
			state.SetSlinkyAppStatePruningParams(),
		)

		require.Equal(t, app.CommitMultiStore().GetPruning().KeepRecent, uint64(10))
	})
}

// test verify VoteExtensionState
// test with a mock + test with a non-mocked application.
func TestAppStates(t *testing.T) {
	// mock
	app := statemocks.NewApplication(t)
	appState := state.NewBaseAppState(app)

	t.Run("if app's retention height is < state.VoteExtensionVerificationHeightOffset - fail", func(t *testing.T) {
		ctx := sdk.Context{}.WithBlockHeight(10)
		app.On("GetBlockRetentionHeight", int64(10)).Return(int64(9)).Once()

		ctx, err := appState.VerifyVoteExtensionState(ctx)
		require.Error(t, err)
		require.Equal(t, ctx, ctx)
	})

	t.Run("if cache multi-store with version fails - fail", func(t *testing.T) {
		ctx := sdk.Context{}.WithBlockHeight(10)
		app.On("GetBlockRetentionHeight", int64(10)).Return(int64(8)).Once()
		app.On("CommitMultiStore").Return(mockCommitMultiStore{expError: true}).Once()

		ctx, err := appState.VerifyVoteExtensionState(ctx)
		require.Error(t, err)
		require.Equal(t, ctx, ctx)
	})

	t.Run("happy path", func(t *testing.T) {
		ctx := sdk.Context{}.WithBlockHeight(10)
		app.On("GetBlockRetentionHeight", int64(10)).Return(int64(8)).Once()
		app.On("CommitMultiStore").Return(mockCommitMultiStore{expError: false}).Once()

		_, err := appState.VerifyVoteExtensionState(ctx)
		require.NoError(t, err)
	})
}

// commitMultiStore mock.
type mockCommitMultiStore struct {
	storetypes.CommitMultiStore
	expError bool
}

func (m mockCommitMultiStore) CacheMultiStoreWithVersion(int64) (storetypes.CacheMultiStore, error) {
	var err error
	if m.expError {
		err = fmt.Errorf("error")
	}
	return nil, err
}
