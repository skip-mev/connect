package raydium_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	oracleconfig "github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/apis/defi/raydium"
	"github.com/skip-mev/slinky/providers/apis/defi/raydium/mocks"
)

// TestMultiJSONRPCClient tests the MultiJSONRPCClient.
func TestMultiJSONRPCClient(t *testing.T) {
	client1 := mocks.NewSolanaJSONRPCClient(t)
	client2 := mocks.NewSolanaJSONRPCClient(t)
	client3 := mocks.NewSolanaJSONRPCClient(t)
	client := raydium.NewMultiJSONRPCClient([]raydium.SolanaJSONRPCClient{client1, client2, client3}, zap.NewNop())

	t.Run("test MultiJSONRPCClient From endpoints", func(t *testing.T) {
		t.Run("invalid endpoint", func(t *testing.T) {
			endpoint := oracleconfig.Endpoint{}

			_, err := raydium.NewMultiJSONRPCClientFromEndpoints([]oracleconfig.Endpoint{endpoint}, zap.NewNop())
			require.Error(t, err)
			require.True(t, strings.Contains(err.Error(), "invalid endpoint"))
		})

		t.Run("endpoints with / wo authentication", func(t *testing.T) {
			endpoints := []oracleconfig.Endpoint{
				{
					URL: "http://localhost:8899",
				},
				{
					URL: "http://localhost:8899/",
					Authentication: oracleconfig.Authentication{
						APIKey:       "test",
						APIKeyHeader: "X-API-Key",
					},
				},
			}

			_, err := raydium.NewMultiJSONRPCClientFromEndpoints(zap.NewNop(), endpoints)
			require.NoError(t, err)
		})
	})

	// test adherence to the context
	t.Run("test failures in underlying client", func(t *testing.T) {
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
		}, fmt.Errorf("error"))

		resp, err := client.GetMultipleAccountsWithOpts(ctx, accounts, opts)
		require.NoError(t, err)

		require.Equal(t, uint64(2), resp.RPCContext.Context.Slot)
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
