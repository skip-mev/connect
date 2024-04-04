package raydium_test

import (
	"context"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/providers/apis/defi/raydium"
	"github.com/skip-mev/slinky/providers/apis/defi/raydium/mocks"
)

// TestMultiJSONRPCClient tests the MultiJSONRPCClient.
func TestMultiJSONRPCClient(t *testing.T) {
	client1 := mocks.NewSolanaJSONRPCClient(t)
	client2 := mocks.NewSolanaJSONRPCClient(t)
	client3 := mocks.NewSolanaJSONRPCClient(t)
	client := raydium.NewMultiJSONRPCClient([]raydium.SolanaJSONRPCClient{client1, client2, client3}, zap.NewNop())

	// test adherence to the context
	t.Run("test adherence to context", func(t *testing.T) {
		accounts := []solana.PublicKey{{}}
		opts := &rpc.GetMultipleAccountsOpts{}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// mocks
		client1.On("GetMultipleAccountsWithOpts", ctx, accounts, opts).Return(&rpc.GetMultipleAccountsResult{}, nil).After(10 * time.Second)
		client2.On("GetMultipleAccountsWithOpts", ctx, accounts, opts).Return(&rpc.GetMultipleAccountsResult{}, nil).After(10 * time.Second)
		client3.On("GetMultipleAccountsWithOpts", ctx, accounts, opts).Return(&rpc.GetMultipleAccountsResult{}, nil).After(10 * time.Second)

		cancel()
		_, err := client.GetMultipleAccountsWithOpts(ctx, accounts, opts)
		require.Error(t, err, "no responses to filter")
	})

	// test correct aggregation of responses
	t.Run("test correct aggregation of responses", func(t *testing.T) {
		accounts := []solana.PublicKey{{}}
		opts := &rpc.GetMultipleAccountsOpts{}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// mocks
		client1.On("GetMultipleAccountsWithOpts", ctx, accounts, opts).Return(&rpc.GetMultipleAccountsResult{
			RPCContext: rpc.RPCContext{
				Context: rpc.Context{
					Slot: 1,
				},
			},
		}, nil)
		client2.On("GetMultipleAccountsWithOpts", ctx, accounts, opts).Return(&rpc.GetMultipleAccountsResult{
			RPCContext: rpc.RPCContext{
				Context: rpc.Context{
					Slot: 2,
				},
			},
		}, nil)
		client3.On("GetMultipleAccountsWithOpts", ctx, accounts, opts).Return(&rpc.GetMultipleAccountsResult{
			RPCContext: rpc.RPCContext{
				Context: rpc.Context{
					Slot: 3,
				},
			},
		}, nil)

		resp, err := client.GetMultipleAccountsWithOpts(ctx, accounts, opts)
		require.NoError(t, err)

		require.Equal(t, uint64(3), resp.RPCContext.Context.Slot)
	})
}
