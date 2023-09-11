package integration

import (
	"context"
	"fmt"
	"strconv"

	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/pelletier/go-toml"
	oracleconfig "github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	abcitypes "github.com/skip-mev/slinky/abci/types"
)

const (
	oracleConfigPath = "oracle.toml"
	appConfigPath    = "config/app.toml"
)

// construct the network from a spec
// ChainBuilderFromChainSpec creates an interchaintest chain builder factory given a ChainSpec
// and returns the associated chain
func ChainBuilderFromChainSpec(t *testing.T, spec *interchaintest.ChainSpec) *cosmos.CosmosChain {
	// require that NumFullNodes == NumValidators == 4
	require.Equal(t, *spec.NumValidators, 4)

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{spec})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	require.Len(t, chains, 1)
	chain := chains[0]

	cosmosChain, ok := chain.(*cosmos.CosmosChain)
	require.True(t, ok)

	return cosmosChain
}

// Modify the application config of a node so that the oracle is running out of process
func SetOracleOutOfProcess(node *cosmos.ChainNode) {
	// read the app config from the node
	bz, err := node.ReadFile(context.Background(), appConfigPath)
	if err != nil {
		panic(err)
	}

	// unmarshal the application config
	var appConfig map[string]interface{}
	err = toml.Unmarshal(bz, &appConfig)
	if err != nil {
		panic(err)
	}

	// get the oracle config
	oracleConfig, ok := appConfig["oracle"].(map[string]interface{})
	if !ok {
		panic("oracle config not found")
	}

	// set the oracle config to out of process
	oracleConfig["in_process"] = false
	oracleConfig["timeout"] = "500ms"

	if len(node.Sidecars) == 0 {
		panic("no sidecars found")
	}

	// get the oracle sidecar
	oracle := node.Sidecars[0]
	// set the oracle port
	oracleConfig["remote_address"] = fmt.Sprintf("%s:%s", oracle.HostName(), "8080")

	appConfig["oracle"] = oracleConfig

	// write back
	bz, err = toml.Marshal(appConfig)
	if err != nil {
		panic(err)
	}

	err = node.WriteFile(context.Background(), bz, appConfigPath)
	if err != nil {
		panic(err)
	}
}

// AddSidecarToNode adds the sidecar configured by the given config to the given node. These are configured
// so that the sidecar is started before the node is started.
func AddSidecarToNode(node *cosmos.ChainNode, conf ibc.SidecarConfig) {
	// create the sidecar process
	node.NewSidecarProcess(
		context.Background(),
		true,
		conf.ProcessName,
		node.DockerClient,
		node.NetworkID,
		conf.Image,
		conf.HomeDir,
		conf.Ports,
		conf.StartCmd,
	)
}

// spin up the network (with side-cars enabled)
// BuildPOBInterchain creates a new Interchain testing env with the configured POB CosmosChain
func BuildPOBInterchain(t *testing.T, ctx context.Context, chain ibc.Chain) *interchaintest.Interchain {
	ic := interchaintest.NewInterchain()
	ic.AddChain(chain)

	// create docker network
	client, networkID := interchaintest.DockerSetup(t)

	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// build the interchain
	err := ic.Build(ctx, nil, interchaintest.InterchainBuildOptions{
		SkipPathCreation: true,
		Client:           client,
		NetworkID:        networkID,
		TestName:         t.Name(),
	})
	require.NoError(t, err)

	return ic
}

// SetOracleConfig writes the given oracle config to the given node
func SetOracleConfig(node *cosmos.ChainNode, conf oracleconfig.Config) {
	if len(node.Sidecars) != 1 {
		panic("expected node to have oracle sidecar")
	}

	oracle := node.Sidecars[0]

	// marshal the oracle config
	bz, err := toml.Marshal(conf)
	if err != nil {
		panic(err)
	}

	// write the oracle config to the node
	err = oracle.WriteFile(context.Background(), bz, oracleConfigPath)
	if err != nil {
		panic(err)
	}
}

// RestartOracle restarts the oracle sidecar for a given node
func RestartOracle(node *cosmos.ChainNode) error {
	if len(node.Sidecars) != 1 {
		panic("expected node to have oracle sidecar")
	}

	oracle := node.Sidecars[0]

	if err := oracle.StopContainer(context.Background()); err != nil {
		return err
	}

	return oracle.StartContainer(context.Background())
}

// StopOracle stops the oracle sidecar for a given node
func StopOracle(node *cosmos.ChainNode) error {
	if len(node.Sidecars) != 1 { 
		panic("expected node to have oracle sidecar")
	}

	oracle := node.Sidecars[0]

	return oracle.StopContainer(context.Background())
}

// StartOracle starts the oracle sidecar for a given node
func StartOracle(node *cosmos.ChainNode) error {
	if len(node.Sidecars) != 1 {
		panic("expected node to have oracle sidecar")
	}

	oracle := node.Sidecars[0]

	return oracle.StartContainer(context.Background())
}

// QueryCurrencyPair queries the chain for the given CurrencyPair, this method returns the grpc response from the module
func QueryCurrencyPairs(chain *cosmos.CosmosChain) (*oracletypes.GetAllCurrencyPairsResponse, error) {
	// get grpc address
	grpcAddr := chain.GetHostGRPCAddress()

	// create the client
	cc, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer cc.Close()

	// create the oracle client
	client := oracletypes.NewQueryClient(cc)

	// query the currency pairs
	return client.GetAllCurrencyPairs(context.Background(), &oracletypes.GetAllCurrencyPairsRequest{})
}

// QueryCurrencyPair queries the price for the given currency-pair given a desired height to query from
func QueryCurrencyPair(chain *cosmos.CosmosChain, cp oracletypes.CurrencyPair, height uint64) (*oracletypes.QuotePrice, int64, error) {
	grpcAddr := chain.GetHostGRPCAddress()

	// create the client
	cc, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, 0, err
	}
	defer cc.Close()

	// create the oracle client
	client := oracletypes.NewQueryClient(cc)

	ctx := context.Background()

	md := metadata.New(map[string]string{
		grpctypes.GRPCBlockHeightHeader: strconv.FormatInt(int64(height), 10),
	})

	ctx = metadata.NewOutgoingContext(ctx, md)

	// query the currency pairs
	res, err := client.GetPrice(ctx, &oracletypes.GetPriceRequest{
		CurrencyPairSelector: &oracletypes.GetPriceRequest_CurrencyPair{
			CurrencyPair: &cp,
		},
	},)
	if err != nil {
		return nil, 0, err
	}
	
	return res.Price, int64(res.Nonce), nil
}

// Submit proposal creates and submits a proposal to the chain
func SubmitProposal(chain *cosmos.CosmosChain, deposit sdk.Coin, submitter string, msgs ...sdk.Msg) (string, error) {
	// build the proposal
	rand := rand.Str(10)
	prop, err := chain.BuildProposal(msgs, rand, rand, rand, deposit.String())

	if err != nil {
		return "", err
	}

	// submit the proposal
	tx, err := chain.SubmitProposal(context.Background(), submitter, prop)
	return tx.ProposalID, err
}

// PassProposal given a proposal id, vote for the proposal and wait for it to pass
func PassProposal(chain *cosmos.CosmosChain, propId string, timeout time.Duration) error {
	if err := WaitForProposalStatus(chain, propId, timeout, govtypesv1.StatusVotingPeriod); err != nil {
		return fmt.Errorf("proposal did not enter voting period: %v", err)
	}

	// have all nodes vote on the proposal
	wg := errgroup.Group{}
	for _, node := range chain.Nodes() {
		n := node // pin
		wg.Go(func() error {
			return n.VoteOnProposal(context.Background(), validatorKey, propId, yes)
		})
	}
	if err := wg.Wait(); err != nil {
		return err
	}
	// wait for the proposal to pass
	if err := WaitForProposalStatus(chain, propId, timeout, govtypesv1.StatusPassed); err != nil {
		return fmt.Errorf("proposal did not pass: %v", err)
	}
	return nil
}

// AddCurrencyPairs creates + submits the proposal to add the given currency-pairs to state, votes for the prop w/ all nodes,
// and waits for the proposal to pass.
func AddCurrencyPairs(chain *cosmos.CosmosChain, authority, denom string, deposit int64, timeout time.Duration, user cosmos.User, cps ...oracletypes.CurrencyPair) error {
	propId, err := SubmitProposal(chain, sdk.NewCoin(denom, math.NewInt(deposit)), user.KeyName(), []sdk.Msg{&oracletypes.MsgAddCurrencyPairs{
		Authority: authority,
		CurrencyPairs: cps,
	}}...)

	if err != nil {
		return err
	}
	
	return PassProposal(chain, propId, timeout)
}

// RemoveCurrencyPairs creates + submits the proposal to remove the given currency-pairs from state, votes for the prop w/ all nodes,
// and waits for the proposal to pass.
func RemoveCurrencyPairs(chain *cosmos.CosmosChain, authority, denom string, deposit int64, timeout time.Duration, user cosmos.User, cpIDs ...string) error {
	propId, err := SubmitProposal(chain, sdk.NewCoin(denom, math.NewInt(deposit)), user.KeyName(), []sdk.Msg{&oracletypes.MsgRemoveCurrencyPairs{
		Authority: authority,
		CurrencyPairIds: cpIDs,
	}}...)

	if err != nil {
		return err
	}
	
	return PassProposal(chain, propId, timeout)
}

// QueryProposal queries the chain for a given proposal
func QueryProposal(chain *cosmos.CosmosChain, propID string) (*govtypes.QueryProposalResponse, error) {
	// get grpc address
	grpcAddr := chain.GetHostGRPCAddress()

	// create the client
	cc, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer cc.Close()

	// create the oracle client
	client := govtypes.NewQueryClient(cc)

	propId, err := strconv.ParseUint(propID, 10, 64)
	if err != nil {
		return nil, err
	}
	// query the currency pairs
	return client.Proposal(context.Background(), &govtypes.QueryProposalRequest{
		ProposalId: propId,
	})
}

// WaitForVotingPeriod, waits for the deposit period for the proposal to end
func WaitForProposalStatus(chain *cosmos.CosmosChain, propID string, timeout time.Duration, status govtypes.ProposalStatus) error {
	return testutil.WaitForCondition(timeout, 1 * time.Second, func() (bool, error) {
		prop, err := QueryProposal(chain, propID)
		if err != nil {
			return false, err
		}
		fmt.Println("\nproposal", prop) // golint:ignore
		return prop.Proposal.Status == status, nil
	})
}

// WaitForHeight waits for the giuve height to be reached
func WaitForHeight(chain *cosmos.CosmosChain, height uint64, timeout time.Duration) error {
	return testutil.WaitForCondition(timeout, 1 * time.Second, func() (bool, error) {
		h, err := chain.Height(context.Background())
		if err != nil {
			return false, err
		}

		return h >= height, nil
	})
}

// WaitForOracleUpdate waits for the first oracle update. This method returns the height that the oracle update occurred
// it returns an error if there is no oracle update by the timeout
func WaitForOracleUpdate(chain *cosmos.CosmosChain, timeout time.Duration, cp oracletypes.CurrencyPair) (uint64, error) {
	client := chain.Nodes()[0].Client 
	var height int64

	if err := testutil.WaitForCondition(timeout, 1 * time.Second, func() (bool, error) {
		blockHeight, err := chain.Height(context.Background())
		if err != nil {
			return false, err
		}
		height = int64(blockHeight)
		
		block, err := client.Block(context.Background(), &height)
		if err != nil {
			return false, err
		}

		// check if the first tx is an oracle update
		if len(block.Block.Txs) == 0 {
			return false, err
		}

		var ve abcitypes.OracleVoteExtension
		if err := ve.Unmarshal(block.Block.Txs[0]); err != nil {
			return false, err
		}

		// check if the currency-pair has an update included for it
		_, ok := ve.Prices[cp.ToString()]

		return ok, nil
	}); err != nil {
		return 0, err
	}

	return uint64(height), nil
}
