package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	cometcfg "github.com/cometbft/cometbft/config"
	cometjson "github.com/cometbft/cometbft/libs/json"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/server"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/skip-mev/slinky/tests/simapp"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

var (
	numValidators   = 4
	minGasPrice     = sdk.NewDecCoinFromDec(simapp.BondDenom, sdk.MustNewDecFromStr("0.02")).String()
	initBalanceStr  = sdk.NewInt64Coin(simapp.BondDenom, 1000000000000000000).String()
	stakeAmount     = sdk.NewInt(100000000000)
	stakeAmountCoin = sdk.NewCoin(simapp.BondDenom, stakeAmount)
)

type (
	TestAccount struct {
		PrivateKey *secp256k1.PrivKey
		Address    sdk.AccAddress
	}

	IntegrationTestSuite struct {
		suite.Suite

		tmpDirs      []string
		chain        *chain
		dkrPool      *dockertest.Pool
		dkrNet       *dockertest.Network
		valResources []*dockertest.Resource
	}
)

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up e2e integration test suite...")

	var err error
	s.chain, err = newChain()
	s.Require().NoError(err)

	s.T().Logf("starting e2e infrastructure; chain-id: %s; datadir: %s", s.chain.id, s.chain.dataDir)

	s.dkrPool, err = dockertest.NewPool("")
	s.Require().NoError(err)

	s.dkrNet, err = s.dkrPool.CreateNetwork(fmt.Sprintf("%s-testnet", s.chain.id))
	s.Require().NoError(err)

	// The bootstrapping phase is as follows:
	//
	// 1. Initialize TestApp validator nodes.
	// 2. Create and initialize TestApp validator genesis files, i.e. setting
	// 		delegate keys for validators.
	// 3. Start TestApp network.
	s.initNodes()
	s.initGenesis()
	s.initValidatorConfigs()
	s.runValidators()
}

func (s *IntegrationTestSuite) TearDownSuite() {
	if str := os.Getenv("E2E_SKIP_CLEANUP"); len(str) > 0 {
		skipCleanup, err := strconv.ParseBool(str)
		s.Require().NoError(err)

		if skipCleanup {
			return
		}
	}

	s.T().Log("tearing down e2e integration test suite...")

	for _, vc := range s.valResources {
		s.Require().NoError(s.dkrPool.Purge(vc))
	}

	s.Require().NoError(s.dkrPool.RemoveNetwork(s.dkrNet))

	os.RemoveAll(s.chain.dataDir)
	for _, td := range s.tmpDirs {
		os.RemoveAll(td)
	}
}

func (s *IntegrationTestSuite) initNodes() {
	s.Require().NoError(s.chain.createAndInitValidators(numValidators))

	// initialize a genesis file for the first validator
	val0ConfigDir := s.chain.validators[0].configDir()

	for _, val := range s.chain.validators {
		valAddr, err := val.keyInfo.GetAddress()
		s.Require().NoError(err)

		s.Require().NoError(initGenesisFile(val0ConfigDir, "", initBalanceStr, valAddr))
	}

	// copy the genesis file to the remaining validators
	for _, val := range s.chain.validators[1:] {
		_, err := copyFile(
			filepath.Join(val0ConfigDir, "config", "genesis.json"),
			filepath.Join(val.configDir(), "config", "genesis.json"),
		)
		s.Require().NoError(err)
	}
}

func (s *IntegrationTestSuite) initGenesis() {
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config

	config.SetRoot(s.chain.validators[0].configDir())
	config.Moniker = s.chain.validators[0].moniker

	genFilePath := config.GenesisFile()
	appGenState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFilePath)
	s.T().Log("starting e2e infrastructure; validator_0 config:", genFilePath)
	s.Require().NoError(err)

	// x/gov
	var govGenState govtypesv1.GenesisState
	s.Require().NoError(cdc.UnmarshalJSON(appGenState[govtypes.ModuleName], &govGenState))

	votingPeriod := 5 * time.Second
	govGenState.Params.VotingPeriod = &votingPeriod
	govGenState.Params.MinDeposit = sdk.NewCoins(sdk.NewCoin(simapp.BondDenom, sdk.NewInt(100)))

	bz, err := cdc.MarshalJSON(&govGenState)
	s.Require().NoError(err)
	appGenState[govtypes.ModuleName] = bz

	var genUtilGenState genutiltypes.GenesisState
	s.Require().NoError(cdc.UnmarshalJSON(appGenState[genutiltypes.ModuleName], &genUtilGenState))

	// x/genutil genesis txs
	genTxs := make([]json.RawMessage, len(s.chain.validators))
	for i, val := range s.chain.validators {
		createValMsg, err := val.buildCreateValidatorMsg(stakeAmountCoin)
		s.Require().NoError(err)

		signedTx, err := val.signMsg(createValMsg)
		s.Require().NoError(err)

		txRaw, err := cdc.MarshalJSON(signedTx)
		s.Require().NoError(err)

		genTxs[i] = txRaw
	}

	genUtilGenState.GenTxs = genTxs

	bz, err = cdc.MarshalJSON(&genUtilGenState)
	s.Require().NoError(err)
	appGenState[genutiltypes.ModuleName] = bz

	bz, err = json.MarshalIndent(appGenState, "", "  ")
	s.Require().NoError(err)

	genDoc.AppState = bz

	bz, err = cometjson.MarshalIndent(genDoc, "", "  ")
	s.Require().NoError(err)

	// write the updated genesis file to each validator
	for _, val := range s.chain.validators {
		writeFile(filepath.Join(val.configDir(), "config", "genesis.json"), bz)
	}
}

func (s *IntegrationTestSuite) initValidatorConfigs() {
	for i, val := range s.chain.validators {
		tmCfgPath := filepath.Join(val.configDir(), "config", "config.toml")

		vpr := viper.New()
		vpr.SetConfigFile(tmCfgPath)
		s.Require().NoError(vpr.ReadInConfig())

		valConfig := cometcfg.DefaultConfig()
		s.Require().NoError(vpr.Unmarshal(valConfig))

		valConfig.P2P.ListenAddress = "tcp://0.0.0.0:26656"
		valConfig.P2P.AddrBookStrict = false
		valConfig.P2P.ExternalAddress = fmt.Sprintf("%s:%d", val.instanceName(), 26656)
		valConfig.RPC.ListenAddress = "tcp://0.0.0.0:26657"
		valConfig.StateSync.Enable = false
		valConfig.LogLevel = "info"

		var peers []string

		for j := 0; j < len(s.chain.validators); j++ {
			if i == j {
				continue
			}

			peer := s.chain.validators[j]
			peerID := fmt.Sprintf("%s@%s%d:26656", peer.nodeKey.ID(), peer.moniker, j)
			peers = append(peers, peerID)
		}

		valConfig.P2P.PersistentPeers = strings.Join(peers, ",")
		cometcfg.WriteConfigFile(tmCfgPath, valConfig)

		// set application configuration
		appCfgPath := filepath.Join(val.configDir(), "config", "app.toml")
		appConfig := srvconfig.DefaultConfig()
		appConfig.API.Enable = true
		appConfig.MinGasPrices = minGasPrice
		appConfig.API.Address = "tcp://0.0.0.0:1317"
		appConfig.GRPC.Address = "0.0.0.0:9090"

		srvconfig.WriteConfigFile(appCfgPath, appConfig)
	}
}

func (s *IntegrationTestSuite) runValidators() {
	s.T().Log("starting POB TestApp validator containers...")

	s.valResources = make([]*dockertest.Resource, len(s.chain.validators))
	for i, val := range s.chain.validators {
		runOpts := &dockertest.RunOptions{
			Name:      val.instanceName(),
			NetworkID: s.dkrNet.Network.ID,
			Mounts: []string{
				fmt.Sprintf("%s/:/root/.testapp", val.configDir()),
			},
			Repository: "docker.io/skip-mev/pob-e2e",
		}

		// expose the first validator for debugging and communication
		if val.index == 0 {
			runOpts.PortBindings = map[docker.Port][]docker.PortBinding{
				"1317/tcp":  {{HostIP: "", HostPort: "1317"}},
				"6060/tcp":  {{HostIP: "", HostPort: "6060"}},
				"6061/tcp":  {{HostIP: "", HostPort: "6061"}},
				"6062/tcp":  {{HostIP: "", HostPort: "6062"}},
				"6063/tcp":  {{HostIP: "", HostPort: "6063"}},
				"6064/tcp":  {{HostIP: "", HostPort: "6064"}},
				"6065/tcp":  {{HostIP: "", HostPort: "6065"}},
				"9090/tcp":  {{HostIP: "", HostPort: "9090"}},
				"26656/tcp": {{HostIP: "", HostPort: "26656"}},
				"26657/tcp": {{HostIP: "", HostPort: "26657"}},
			}
		}

		resource, err := s.dkrPool.RunWithOptions(runOpts, noRestart)
		s.Require().NoError(err)

		s.valResources[i] = resource
		s.T().Logf("started POB TestApp validator container: %s", resource.Container.ID)
	}

	rpcClient, err := rpchttp.New("tcp://localhost:26657", "/websocket")
	s.Require().NoError(err)

	s.Require().Eventually(
		func() bool {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			status, err := rpcClient.Status(ctx)
			if err != nil {
				return false
			}

			// let the node produce a few blocks
			if status.SyncInfo.CatchingUp || status.SyncInfo.LatestBlockHeight < 3 {
				return false
			}

			return true
		},
		2*time.Minute,
		time.Second,
		"POB TestApp node failed to produce blocks",
	)
}

func noRestart(config *docker.HostConfig) {
	// in this case we don't want the nodes to restart on failure
	config.RestartPolicy = docker.RestartPolicy{
		Name: "no",
	}
}
