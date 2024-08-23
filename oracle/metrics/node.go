package metrics

import (
	"context"
	"fmt"
	"strings"

	"github.com/skip-mev/connect/v2/oracle/config"
	slinkygrpc "github.com/skip-mev/connect/v2/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
)

type NodeClient struct {
	conn *grpc.ClientConn
}

func NewNodeClient(endpoint config.Endpoint) (*NodeClient, error) {
	conn, err := slinkygrpc.NewClient(
		endpoint.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithNoProxy(),
	)
	if err != nil {
		return nil, err
	}

	return &NodeClient{
		conn,
	}, nil
}

func (nc *NodeClient) DeriveNodeIdentifier() (string, error) {
	svcclient := cmtservice.NewServiceClient(nc.conn)

	info, err := svcclient.GetNodeInfo(context.Background(), &cmtservice.GetNodeInfoRequest{})

	if err != nil {
		return "", err
	}

	moniker := strings.ReplaceAll(info.DefaultNodeInfo.Moniker, " ", "-")
	network := info.DefaultNodeInfo.Network

	identifier := fmt.Sprintf("%s_%s", network, moniker)

	return identifier, nil
}
