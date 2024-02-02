package petri

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	tmloadtest "github.com/informalsystems/tm-load-test/pkg/loadtest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/skip-mev/petri/cosmosutil/v2"
	"github.com/skip-mev/petri/loadtest/v2"
	petritypes "github.com/skip-mev/petri/types/v2"

	"github.com/skip-mev/slinky/tests/simapp"
)

const (
	envKeepAlive = "ORACLE_INTEGRATION_KEEPALIVE"
)

// SlinkyIntegrationSuite is the testify suite for building and running slinky testapp networks using petri
type SlinkyIntegrationSuite struct {
	suite.Suite

	logger *zap.Logger

	spec *petritypes.ChainConfig

	chain petritypes.ChainI
}

func NewSlinkyIntegrationSuite(spec *petritypes.ChainConfig) *SlinkyIntegrationSuite {
	return &SlinkyIntegrationSuite{
		spec: spec,
	}
}

func (s *SlinkyIntegrationSuite) SetupSuite() {
	// create the logger
	var err error
	s.logger, err = zap.NewDevelopment()
	s.Require().NoError(err)

	// create the chain
	s.chain, err = GetChain(context.Background(), s.logger)
	s.Require().NoError(err)

	//initialize the chain
	err = s.chain.Init(context.Background())
	s.Require().NoError(err)
}

func (s *SlinkyIntegrationSuite) TearDownSuite() {
	// get the oracle integration-test suite keep alive env
	if ok := os.Getenv(envKeepAlive); ok == "" {
		return
	}
	err := s.chain.Teardown(context.Background())
	s.Require().NoError(err)
}

func (s *SlinkyIntegrationSuite) TestSlinkyLoad() {
	err := s.chain.WaitForHeight(context.Background(), 1)
	s.Require().NoError(err)

	encCfg := cosmosutil.EncodingConfig{
		InterfaceRegistry: s.chain.GetInterfaceRegistry(),
		Codec:             codec.NewProtoCodec(s.chain.GetInterfaceRegistry()),
		TxConfig:          s.chain.GetTxConfig(),
	}

	clientFactory, err := loadtest.NewDefaultClientFactory(
		loadtest.ClientFactoryConfig{
			Chain:                 s.chain,
			Seeder:                cosmosutil.NewInteractingWallet(s.chain, s.chain.GetFaucetWallet(), encCfg),
			WalletConfig:          s.spec.WalletConfig,
			AmountToSend:          1000000000,
			SkipSequenceIncrement: false,
			EncodingConfig:        encCfg,
			MsgGenerator:          s.genMsg,
		},
		simapp.ModuleBasics,
	)
	s.Require().NoError(err)

	err = tmloadtest.RegisterClientFactory("slinky", clientFactory)
	s.Require().NoError(err)

	var endpoints []string
	for _, val := range s.chain.GetValidators() {
		endpoint, err := val.GetTMClient(context.Background())
		s.Require().NoError(err)

		url := strings.Replace(endpoint.Remote(), "http", "ws", -1)

		endpoints = append(endpoints, fmt.Sprintf("%s/websocket", url))
	}

	cfg := tmloadtest.Config{
		ClientFactory:        "slinky",
		Connections:          1,
		Endpoints:            endpoints,
		Time:                 60,
		SendPeriod:           1,
		Rate:                 350,
		Size:                 250,
		Count:                -1,
		BroadcastTxMethod:    "async",
		EndpointSelectMethod: "supplied",
	}
	err = tmloadtest.ExecuteStandalone(cfg)
	s.Require().NoError(err)
}

func (s *SlinkyIntegrationSuite) genMsg(senderAddress []byte) ([]sdk.Msg, petritypes.GasSettings, error) {
	return []sdk.Msg{
			&bank.MsgSend{
				FromAddress: string(senderAddress),
				ToAddress:   "cosmos1qy3523p8x9z0j6z3qg3y7t4v6gj6z9q8r9m9x5",
				Amount:      sdk.NewCoins(sdk.NewInt64Coin("stake", 1)),
			},
		}, petritypes.GasSettings{
			Gas:         200000,
			GasDenom:    "stake",
			PricePerGas: 0,
		}, nil
}
