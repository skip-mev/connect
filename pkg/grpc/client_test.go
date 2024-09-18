package grpc_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	reflectionpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	connectgrpc "github.com/skip-mev/connect/v2/pkg/grpc"
)

func TestClient(t *testing.T) {
	// spin up a mock grpc-server + test connecting to it via diff addresses
	srv := grpc.NewServer()

	// listen on a random open port
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	// register reflection service on the server
	reflection.Register(srv)

	// start the server
	go func() {
		srv.Serve(lis)
	}()

	_, port, err := net.SplitHostPort(lis.Addr().String())
	if err != nil {
		t.Fatalf("failed to parse address: %v", err)
	}

	t.Run("try dialing via non supported GRPC target URL (i.e tcp prefix)", func(t *testing.T) {
		// try dialing via non supported GRPC target URL (i.e tcp prefix)
		client, err := connectgrpc.NewClient(fmt.Sprintf("tcp://localhost:%s", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		// ping the server
		_, err = reflectionpb.NewServerReflectionClient(client).ServerReflectionInfo(context.Background())
		require.NoError(t, err)
	})

	t.Run("try dialing via supported GRPC target URL (i.e host:port)", func(t *testing.T) {
		// try dialing via supported GRPC target URL (i.e host:port)
		client, err := connectgrpc.NewClient(fmt.Sprintf("localhost:%s", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		// ping the server
		_, err = reflectionpb.NewServerReflectionClient(client).ServerReflectionInfo(context.Background())
		require.NoError(t, err)
	})
}
