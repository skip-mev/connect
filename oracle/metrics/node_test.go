package metrics_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	p2p "github.com/cometbft/cometbft/proto/tendermint/p2p"
	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/skip-mev/connect/v2/oracle/config"
	oraclemetrics "github.com/skip-mev/connect/v2/oracle/metrics"
)

// mocks a remote sdk grpc service.
type mockServiceServer struct {
	cmtservice.ServiceServer
}

func (m *mockServiceServer) GetNodeInfo(_ context.Context, _ *cmtservice.GetNodeInfoRequest) (*cmtservice.GetNodeInfoResponse, error) {
	return &cmtservice.GetNodeInfoResponse{
		DefaultNodeInfo: &p2p.DefaultNodeInfo{
			Network: "neutron-1",
			Moniker: "some🫵😹node moniker",
		},
	}, nil
}

func TestNodeClientImpl_DeriveNodeIdentifier(t *testing.T) {
	// mock the remote node
	srv := grpc.NewServer()
	mockSvcServer := &mockServiceServer{}
	cmtservice.RegisterServiceServer(srv, mockSvcServer)
	reflection.Register(srv)

	// let the os assign a port
	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	go func() {
		srv.Serve(lis)
	}()
	defer srv.Stop()

	// get the port from earlier
	_, port, err := net.SplitHostPort(lis.Addr().String())
	require.NoError(t, err)

	// test conn
	endpoint := config.Endpoint{URL: fmt.Sprintf("localhost:%s", port)}
	nodeClient, err := oraclemetrics.NewNodeClient(endpoint)
	require.NoError(t, err)

	// test DeriveNodeIdentifier
	identifier, err := nodeClient.DeriveNodeIdentifier()
	require.NoError(t, err)
	require.Equal(t, "neutron-1.some🫵😹node-moniker", identifier)
}
