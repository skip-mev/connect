package e2e

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	cmtclient "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// createClientContext creates a client.Context for use in integration tests.
// Note, it assumes all queries and broadcasts go to the first node.
func (s *IntegrationTestSuite) createClientContext() client.Context {
	node := s.valResources[0]

	rpcURI := node.GetHostPort("26657/tcp")
	gRPCURI := node.GetHostPort("9090/tcp")

	rpcClient, err := client.NewClientFromNode(rpcURI)
	s.Require().NoError(err)

	grpcClient, err := grpc.Dial(gRPCURI, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}...)
	s.Require().NoError(err)

	return client.Context{}.
		WithNodeURI(rpcURI).
		WithClient(rpcClient).
		WithGRPCClient(grpcClient).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithCodec(encodingConfig.Codec).
		WithChainID(s.chain.id).
		WithBroadcastMode(flags.BroadcastSync)
}

// waitForABlock will wait until the current block height has increased by a single block.
func (s *IntegrationTestSuite) waitForABlock() {
	height := s.queryCurrentHeight()
	s.Require().Eventually(
		func() bool {
			return s.queryCurrentHeight() >= height+1
		},
		10*time.Second,
		50*time.Millisecond,
	)
}

// queryBuilderParams returns the params of the builder module.
func (s *IntegrationTestSuite) queryAllCurrencyPairs() []oracletypes.CurrencyPair {
	queryClient := oracletypes.NewQueryClient(s.createClientContext())

	resp, err := queryClient.GetAllCurrencyPairs(context.Background(), &oracletypes.GetAllCurrencyPairsRequest{})
	s.Require().NoError(err)

	return resp.CurrencyPairs
}

// queryPriceForCurrencyPair returns the price for a given currency pair.
func (s *IntegrationTestSuite) queryPriceForCurrencyPair(base, quote string) (*oracletypes.GetPriceResponse, error) {
	queryClient := oracletypes.NewQueryClient(s.createClientContext())

	cp := oracletypes.NewCurrencyPair(base, quote)

	req := &oracletypes.GetPriceRequest{
		CurrencyPairSelector: &oracletypes.GetPriceRequest_CurrencyPairId{
			CurrencyPairId: cp.ToString(),
		},
	}

	return queryClient.GetPrice(context.Background(), req)
}

// queryCurrentHeight returns the current block height.
func (s *IntegrationTestSuite) queryCurrentHeight() uint64 {
	queryClient := cmtclient.NewServiceClient(s.createClientContext())

	req := &cmtclient.GetLatestBlockRequest{}
	resp, err := queryClient.GetLatestBlock(context.Background(), req)
	s.Require().NoError(err)

	return uint64(resp.SdkBlock.Header.Height)
}

// queryValidators returns the validators of the network.
func (s *IntegrationTestSuite) queryValidators() []stakingtypes.Validator {
	queryClient := stakingtypes.NewQueryClient(s.createClientContext())

	req := &stakingtypes.QueryValidatorsRequest{}
	resp, err := queryClient.Validators(context.Background(), req)
	s.Require().NoError(err)

	return resp.Validators
}
