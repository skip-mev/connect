package metrics

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/skip-mev/connect/v2/oracle/config"
	connectgrpc "github.com/skip-mev/connect/v2/pkg/grpc"

	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
)

//go:generate mockery --name NodeClient --filename mock_node_client.go
type NodeClient interface {
	DeriveNodeIdentifier() (string, error)
}

type NodeClientImpl struct {
	conn *grpc.ClientConn
}

func NewNodeClient(endpoint config.Endpoint) (NodeClient, error) {
	conn, err := connectgrpc.NewClient(
		endpoint.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithNoProxy(),
	)
	if err != nil {
		return nil, err
	}

	return &NodeClientImpl{
		conn,
	}, nil
}

func (nc *NodeClientImpl) DeriveNodeIdentifier() (string, error) {
	svcclient := cmtservice.NewServiceClient(nc.conn)

	info, err := svcclient.GetNodeInfo(context.Background(), &cmtservice.GetNodeInfoRequest{})
	if err != nil {
		return "", err
	}

	moniker := strings.ReplaceAll(info.DefaultNodeInfo.Moniker, " ", "-")
	network := info.DefaultNodeInfo.Network

	identifier := fmt.Sprintf("%s.%s", network, moniker)

	return identifier, nil
}
