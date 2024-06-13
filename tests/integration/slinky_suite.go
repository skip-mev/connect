package integration

import (
	"context"
	"encoding/hex"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/skip-mev/slinky/providers/apis/marketmap"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/stretchr/testify/suite"

	oracleconfig "github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	envKeepAlive          = "ORACLE_INTEGRATION_KEEPALIVE"
	genesisAmount         = 1000000000
	defaultDenom          = "stake"
	validatorKey          = "validator"
	yes                   = "yes"
	deposit               = 1000000
	userMnemonic          = "foster poverty abstract scorpion short shrimp tilt edge romance adapt only benefit moral another where host egg echo ability wisdom lizard lazy pool roast"
	userAccountAddressHex = "877E307618AB73E009A978AC32E0264791F6D40A"
)

func DefaultOracleSidecar(image ibc.DockerImage) ibc.SidecarConfig {
	return ibc.SidecarConfig{
		ProcessName: "oracle",
		Image:       image,
		HomeDir:     "/oracle",
		Ports:       []string{"8080", "8081"},
		StartCmd: []string{
			"slinky",
			"--oracle-config", "/oracle/oracle.json",
		},
		ValidatorProcess: true,
		PreStart:         true,
	}
}

func DefaultOracleConfig(url string) oracleconfig.OracleConfig {
	cfg := marketmap.DefaultAPIConfig
	cfg.Endpoints = []oracleconfig.Endpoint{
		{
			URL: url,
		},
	}

	// Create the oracle config
	oracleConfig := oracleconfig.OracleConfig{
		UpdateInterval: 500 * time.Millisecond,
		MaxPriceAge:    1 * time.Minute,
		Host:           "0.0.0.0",
		Port:           "8080",
		Providers: map[string]oracleconfig.ProviderConfig{
			marketmap.Name: {
				Name: marketmap.Name,
				API:  cfg,
				Type: "market_map_provider",
			},
		},
	}

	return oracleConfig
}

func DefaultMarketMap() mmtypes.MarketMap {
	return mmtypes.MarketMap{}
}

func GetOracleSideCar(node *cosmos.ChainNode) *cosmos.SidecarProcess {
	if len(node.Sidecars) == 0 {
		panic("no sidecars found")
	}
	return node.Sidecars[0]
}

type SlinkyIntegrationSuite struct {
	suite.Suite

	spec *interchaintest.ChainSpec

	// add more fields here as necessary
	chain *cosmos.CosmosChain

	// oracle side-car config
	oracleConfig ibc.SidecarConfig

	// user
	user cosmos.User

	// default token denom
	denom string

	// authority address
	authority sdk.AccAddress

	// block time
	blockTime time.Duration

	// interchain constructor
	icc InterchainConstructor

	// interchain
	ic Interchain

	// chain constructor
	cc ChainConstructor
}

// Option is a function that modifies the SlinkyIntegrationSuite
type Option func(*SlinkyIntegrationSuite)

// WithDenom sets the token denom
func WithDenom(denom string) Option {
	return func(s *SlinkyIntegrationSuite) {
		s.denom = denom
	}
}

// WithAuthority sets the authority address
func WithAuthority(addr sdk.AccAddress) Option {
	return func(s *SlinkyIntegrationSuite) {
		s.authority = addr
	}
}

// WithBlockTime sets the block time
func WithBlockTime(t time.Duration) Option {
	return func(s *SlinkyIntegrationSuite) {
		s.blockTime = t
	}
}

// WithInterchainConstructor sets the interchain constructor
func WithInterchainConstructor(ic InterchainConstructor) Option {
	return func(s *SlinkyIntegrationSuite) {
		s.icc = ic
	}
}

// WithChainConstructor sets the chain constructor
func WithChainConstructor(cc ChainConstructor) Option {
	return func(s *SlinkyIntegrationSuite) {
		s.cc = cc
	}
}

func NewSlinkyIntegrationSuite(spec *interchaintest.ChainSpec, oracleImage ibc.DockerImage, opts ...Option) *SlinkyIntegrationSuite {
	suite := &SlinkyIntegrationSuite{
		spec:         spec,
		oracleConfig: DefaultOracleSidecar(oracleImage),
		denom:        defaultDenom,
		authority:    authtypes.NewModuleAddress(govtypes.ModuleName),
		blockTime:    10 * time.Second,
		icc:          DefaultInterchainConstructor,
		cc:           DefaultChainConstructor,
	}

	for _, opt := range opts {
		opt(suite)
	}

	return suite
}

func (s *SlinkyIntegrationSuite) SetupSuite() {
	// update market-map params to add the user as the market-authority
	accountAddressBz, err := hex.DecodeString(userAccountAddressHex)
	if err != nil {
		panic(err)
	}
	accountAddress, err := bech32.ConvertAndEncode(s.spec.ChainConfig.Bech32Prefix, accountAddressBz)
	if err != nil {
		panic(err)
	}
	existingGenesisModifier := s.spec.ChainConfig.ModifyGenesis
	s.spec.ChainConfig.ModifyGenesis = func(cc ibc.ChainConfig, genesisBz []byte) ([]byte, error) {
		genesisBz, err := cosmos.ModifyGenesis([]cosmos.GenesisKV{
			cosmos.NewGenesisKV(
				"app_state.marketmap.params.admin",
				accountAddress,
			),
			cosmos.NewGenesisKV(
				"app_state.marketmap.params.market_authorities.0",
				accountAddress,
			),
			cosmos.NewGenesisKV(
				"app_state.staking.params.unbonding_time",
				"10s",
			),
		})(cc, genesisBz)
		if err != nil {
			return nil, err
		}

		return existingGenesisModifier(cc, genesisBz)
	}

	chains := s.cc(s.T(), s.spec)

	if len(chains) < 1 {
		panic("no chains created")
	}

	chains[0].WithPreStartNodes(func(c *cosmos.CosmosChain) {
		// for each node in the chain, set the sidecars
		for i := range c.Nodes() {
			// pin
			node := c.Nodes()[i]
			// add sidecars to node
			AddSidecarToNode(node, s.oracleConfig)

			// set config for the oracle
			oracleCfg := DefaultOracleConfig("localhost:9090")
			SetOracleConfigsOnOracle(GetOracleSideCar(node), oracleCfg)

			// set the out-of-process oracle config for all nodes
			node.WithPreStartNode(func(n *cosmos.ChainNode) {
				SetOracleConfigsOnApp(n)
			})
		}
	})

	// start the chain
	s.ic = s.icc(context.Background(), s.T(), chains)
	s.chain = chains[0]
	s.user, err = interchaintest.GetAndFundTestUserWithMnemonic(context.Background(), s.T().Name(), userMnemonic, math.NewInt(genesisAmount), s.chain)
	s.Require().NoError(err)
}

func (s *SlinkyIntegrationSuite) TearDownSuite() {
	defer s.Teardown()
	// get the oracle integration-test suite keep alive env
	if ok := os.Getenv(envKeepAlive); ok == "" {
		return
	}

	// await on a signal to keep the chain running
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s.T().Log("Keeping the chain running")
	<-sig
}

func (s *SlinkyIntegrationSuite) Teardown() {
	// stop all nodes + sidecars in the chain
	ctx := context.Background()
	if s.chain == nil {
		return
	}

	s.chain.StopAllNodes(ctx)
	s.chain.StopAllSidecars(ctx)

	// if there is a provider, stop that as well
	if s.chain.Provider != nil {
		s.chain.Provider.StopAllNodes(ctx)
		s.chain.Provider.StopAllSidecars(ctx)
	}
}

func (s *SlinkyIntegrationSuite) SetupTest() {
	s.TearDownSuite()
	s.SetupSuite()

	// reset the oracle services
	// start all oracles
	for _, node := range s.chain.Nodes() {
		oCfg := DefaultOracleConfig(translateGRPCAddr(s.chain))

		SetOracleConfigsOnOracle(GetOracleSideCar(node), oCfg)
		s.Require().NoError(RestartOracle(node))
	}
}

func translateGRPCAddr(chain *cosmos.CosmosChain) string {
	return chain.GetGRPCAddress()
}

type SlinkyOracleIntegrationSuite struct {
	*SlinkyIntegrationSuite
}

func NewSlinkyOracleIntegrationSuite(suite *SlinkyIntegrationSuite) *SlinkyOracleIntegrationSuite {
	return &SlinkyOracleIntegrationSuite{
		SlinkyIntegrationSuite: suite,
	}
}

func (s *SlinkyOracleIntegrationSuite) TestValidatorExit() {
	// Unbond the first validator
	ctx := context.Background()
	vals, err := s.chain.StakingQueryValidators(ctx, stakingtypes.Bonded.String())
	s.Require().NoError(err)
	val := vals[0].OperatorAddress

	heightBeforeUnBond, err := s.chain.Height(ctx)
	s.T().Logf("height before unbond: %d", heightBeforeUnBond)
	s.Require().NoError(err)

	wasErr := true
	s.Require().NoError(err)
	for _, node := range s.chain.Validators {
		height, err := s.chain.Height(ctx)
		s.Require().NoError(err)

		s.T().Logf("attempting to unbond at height: %d", height)
		err = node.StakingUnbond(ctx, validatorKey, val, vals[0].BondedTokens().String()+s.denom)
		s.T().Logf("unbond error: %v", err)
		if err == nil {
			wasErr = false
		}
	}
	s.Require().False(wasErr)
	height, err := s.chain.Height(ctx)
	s.Require().NoError(err)
	s.T().Logf("unbond successful after height: %d", height)

	// Determine the height at which the unbond happened.
}

func getIDForCurrencyPair(ctx context.Context, client oracletypes.QueryClient, cp slinkytypes.CurrencyPair) (uint64, error) {
	// query for the given currency pair
	resp, err := client.GetPrice(ctx, &oracletypes.GetPriceRequest{
		CurrencyPair: cp,
	})
	if err != nil {
		return 0, err
	}

	return resp.Id, nil
}
