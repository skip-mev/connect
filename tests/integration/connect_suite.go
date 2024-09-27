package integration

import (
	"context"
	"encoding/hex"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	connectabci "github.com/skip-mev/connect/v2/abci/ve/types"
	oracleconfig "github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/apis/marketmap"
	"github.com/skip-mev/connect/v2/providers/static"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
)

const (
	envKeepAlive          = "ORACLE_INTEGRATION_KEEPALIVE"
	genesisAmount         = 1000000000
	defaultDenom          = "stake"
	validatorKey          = "validator"
	yes                   = "yes"
	userMnemonic          = "foster poverty abstract scorpion short shrimp tilt edge romance adapt only benefit moral another where host egg echo ability wisdom lizard lazy pool roast"
	userAccountAddressHex = "877E307618AB73E009A978AC32E0264791F6D40A"
	gasPrice              = 100
)

func DefaultOracleSidecar(image ibc.DockerImage) ibc.SidecarConfig {
	return ibc.SidecarConfig{
		ProcessName: "oracle",
		Image:       image,
		HomeDir:     "/oracle",
		Ports:       []string{"8080", "8081"},
		StartCmd: []string{
			"connect",
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

type ConnectIntegrationSuite struct {
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

// Option is a function that modifies the ConnectIntegrationSuite
type Option func(*ConnectIntegrationSuite)

// WithDenom sets the token denom
func WithDenom(denom string) Option {
	return func(s *ConnectIntegrationSuite) {
		s.denom = denom
	}
}

// WithAuthority sets the authority address
func WithAuthority(addr sdk.AccAddress) Option {
	return func(s *ConnectIntegrationSuite) {
		s.authority = addr
	}
}

// WithBlockTime sets the block time
func WithBlockTime(t time.Duration) Option {
	return func(s *ConnectIntegrationSuite) {
		s.blockTime = t
	}
}

// WithInterchainConstructor sets the interchain constructor
func WithInterchainConstructor(ic InterchainConstructor) Option {
	return func(s *ConnectIntegrationSuite) {
		s.icc = ic
	}
}

// WithChainConstructor sets the chain constructor
func WithChainConstructor(cc ChainConstructor) Option {
	return func(s *ConnectIntegrationSuite) {
		s.cc = cc
	}
}

// CreateTx creates a new transaction to be signed by the given user, including a provided set of messages
func CreateTx(t *testing.T, chain *cosmos.CosmosChain, user cosmos.User, GasPrice int64, msgs ...sdk.Msg) []byte {
	bc := cosmos.NewBroadcaster(t, chain)

	ctx := context.Background()
	// create tx factory + Client Context
	txf, err := bc.GetFactory(ctx, user)
	require.NoError(t, err)

	cc, err := bc.GetClientContext(ctx, user)
	require.NoError(t, err)

	txf = txf.WithSimulateAndExecute(true)

	txf, err = txf.Prepare(cc)
	require.NoError(t, err)

	// get gas for tx
	txf.WithGas(25000000)

	// update sequence number
	txf = txf.WithSequence(txf.Sequence())
	txf = txf.WithGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(chain.Config().Denom, math.NewInt(GasPrice))).String())

	// sign the tx
	txBuilder, err := txf.BuildUnsignedTx(msgs...)
	require.NoError(t, err)

	require.NoError(t, tx.Sign(cc.CmdContext, txf, cc.GetFromName(), txBuilder, true))

	// encode and return
	bz, err := cc.TxConfig.TxEncoder()(txBuilder.GetTx())
	require.NoError(t, err)
	return bz
}

func NewConnectIntegrationSuite(spec *interchaintest.ChainSpec, oracleImage ibc.DockerImage, opts ...Option) *ConnectIntegrationSuite {
	suite := &ConnectIntegrationSuite{
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

func (s *ConnectIntegrationSuite) SetupSuite() {
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

func (s *ConnectIntegrationSuite) TearDownSuite() {
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

func (s *ConnectIntegrationSuite) Teardown() {
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

func (s *ConnectIntegrationSuite) SetupTest() {
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

type ConnectOracleIntegrationSuite struct {
	*ConnectIntegrationSuite
}

func NewConnectOracleIntegrationSuite(suite *ConnectIntegrationSuite) *ConnectOracleIntegrationSuite {
	return &ConnectOracleIntegrationSuite{
		ConnectIntegrationSuite: suite,
	}
}

func (s *ConnectOracleIntegrationSuite) TestOracleModule() {
	// query the oracle module grpc service for any CurrencyPairs
	s.Run("QueryCurrencyPairs - no currency-pairs reported", func() {
		resp, err := QueryCurrencyPairs(s.chain)
		s.Require().NoError(err)
		s.Require().True(len(resp.CurrencyPairs) == 0)
	})

	// pass a governance proposal to approve a new currency-pair, and check Prices are reported
	s.Run("Add a currency-pair and check Prices", func() {
		s.Require().NoError(s.AddCurrencyPairs(s.chain, s.user, 1.1, []mmtypes.Ticker{
			enabledTicker(connecttypes.CurrencyPair{
				Base:  "BTC",
				Quote: "USD",
			}),
		}...))

		// check that the currency-pair is added to state
		resp, err := QueryCurrencyPairs(s.chain)
		s.Require().NoError(err)
		s.Require().True(len(resp.CurrencyPairs) == 1)
		s.Require().Equal(resp.CurrencyPairs[0].Base, "BTC")
		s.Require().Equal(resp.CurrencyPairs[0].Quote, "USD")
	})

	s.Run("Add multiple Currency Pairs", func() {
		cp1 := connecttypes.NewCurrencyPair("ETH", "USD")
		cp2 := connecttypes.NewCurrencyPair("USDT", "USD")
		s.Require().NoError(s.AddCurrencyPairs(s.chain, s.user, 1.1, []mmtypes.Ticker{
			enabledTicker(cp1),
			enabledTicker(cp2),
		}...))

		resp, err := QueryCurrencyPairs(s.chain)
		s.Require().NoError(err)
		s.Require().True(len(resp.CurrencyPairs) == 3)

		s.Run("fail to remove an enabled market", func() {
			s.Require().Error(s.RemoveMarket(s.chain, []connecttypes.CurrencyPair{cp1}))

			// check not removed
			market, err := QueryMarket(s.chain, cp1)
			s.Require().NoError(err)
			s.Require().NotNil(market)
		})
	})

	s.Run("remove a disabled market", func() {
		disabledCP := connecttypes.NewCurrencyPair("DIS", "ABLE")
		s.Require().NoError(s.AddCurrencyPairs(s.chain, s.user, 1.1, []mmtypes.Ticker{
			disabledTicker(disabledCP),
		}...))

		market, err := QueryMarket(s.chain, disabledCP)
		s.Require().NoError(err)
		s.Require().NotNil(market)

		s.Require().NoError(s.RemoveMarket(s.chain, []connecttypes.CurrencyPair{disabledCP}))

		// check removed
		_, err = QueryMarket(s.chain, disabledCP)
		s.Require().Error(err)
	})

	s.Run("remove a non existent market", func() {
		nonexistentCP := connecttypes.NewCurrencyPair("NON", "EXIST")

		// check removed doesnt exist
		_, err := QueryMarket(s.chain, nonexistentCP)
		s.Require().Error(err)

		// tx will not error
		s.Require().NoError(s.RemoveMarket(s.chain, []connecttypes.CurrencyPair{nonexistentCP}))
	})
}

func translateGRPCAddr(chain *cosmos.CosmosChain) string {
	return chain.GetGRPCAddress()
}

func (s *ConnectOracleIntegrationSuite) TestNodeFailures() {
	ethusdcCP := connecttypes.NewCurrencyPair("ETH", "USDC")
	tickerETHUSDC := mmtypes.Ticker{
		CurrencyPair:     ethusdcCP,
		Decimals:         8,
		MinProviderCount: 1,
		Enabled:          true,
		Metadata_JSON:    "",
	}

	s.Require().NoError(s.AddCurrencyPairs(s.chain, s.user, 1.1, []mmtypes.Ticker{
		tickerETHUSDC,
	}...))

	cc, closeFn, err := GetChainGRPC(s.chain)
	s.Require().NoError(err)
	defer closeFn()

	id, err := getIDForCurrencyPair(context.Background(), oracletypes.NewQueryClient(cc), ethusdcCP)
	s.Require().NoError(err)

	zero := big.NewInt(0)
	zeroBz, err := zero.GobEncode()
	s.Require().NoError(err)

	// configure failing providers for various sets of nodes
	s.Run("all nodes report Prices", func() {
		// update all oracle configs
		for _, node := range s.chain.Nodes() {
			oracleConfig := DefaultOracleConfig(translateGRPCAddr(s.chain))
			oracleConfig.Providers[static.Name] = oracleconfig.ProviderConfig{
				Name: static.Name,
				API: oracleconfig.APIConfig{
					Enabled:          true,
					Timeout:          250 * time.Millisecond,
					Interval:         250 * time.Millisecond,
					ReconnectTimeout: 250 * time.Millisecond,
					MaxQueries:       1,
					Endpoints: []oracleconfig.Endpoint{
						{
							URL: "http://un-used-url.com",
						},
					},
					Atomic: true,
					Name:   static.Name,
				},
				Type: types.ConfigType,
			}

			oracle := GetOracleSideCar(node)
			SetOracleConfigsOnOracle(oracle, oracleConfig)
			s.Require().NoError(RestartOracle(node))
		}

		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []connectabci.OracleVoteExtension{
			{
				Prices: map[uint64][]byte{
					id: zeroBz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id: zeroBz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id: zeroBz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id: zeroBz,
				},
			},
		})
		s.Require().NoError(err)
		// query for the given currency pair
		resp, _, err := QueryCurrencyPair(s.chain, ethusdcCP, height)
		s.Require().NoError(err)
		s.Require().Equal(resp.Price.Int64(), int64(110000000))
	})

	s.Run("single oracle down, price updates", func() {
		// stop single node's oracle process and check that all Prices are reported
		node := s.chain.Nodes()[0]
		StopOracle(node)

		// expect the following vote-extensions
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []connectabci.OracleVoteExtension{
			{
				Prices: map[uint64][]byte{},
			},
			{
				Prices: map[uint64][]byte{
					id: zeroBz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id: zeroBz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id: zeroBz,
				},
			},
		})
		s.Require().NoError(err)

		_, oldNonce, err := QueryCurrencyPair(s.chain, ethusdcCP, height-1)
		s.Require().NoError(err)

		_, newNonce, err := QueryCurrencyPair(s.chain, ethusdcCP, height)
		s.Require().NoError(err)

		// expect update for height
		s.Require().Equal(newNonce, oldNonce+1)

		// start the oracle again
		StartOracle(node)
	})

	s.Run("single node down, price updates", func() {
		// stop single node's oracle process and check that all prices are reported
		node := s.chain.Nodes()[0]
		StopOracle(node)

		// expect the following vote-extensions
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []connectabci.OracleVoteExtension{
			{
				Prices: map[uint64][]byte{},
			},
			{
				Prices: map[uint64][]byte{
					id: zeroBz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id: zeroBz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id: zeroBz,
				},
			},
		})
		s.Require().NoError(err)

		_, oldNonce, err := QueryCurrencyPair(s.chain, ethusdcCP, height-1)
		s.Require().NoError(err)

		_, newNonce, err := QueryCurrencyPair(s.chain, ethusdcCP, height)
		s.Require().NoError(err)

		// expect update for height
		s.Require().Equal(newNonce, oldNonce+1)

		StartOracle(node)
	})

	s.Run("only 1 node reports a price (oracles are down)", func() {
		// shut down all oracles except for one
		for _, node := range s.chain.Nodes()[1:] {
			StopOracle(node)
		}

		// expect the given oracle reports
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []connectabci.OracleVoteExtension{
			{
				Prices: map[uint64][]byte{},
			},
			{
				Prices: map[uint64][]byte{},
			},
			{
				Prices: map[uint64][]byte{},
			},
			{
				Prices: map[uint64][]byte{
					id: zeroBz,
				},
			},
		})
		s.Require().NoError(err)

		_, oldNonce, err := QueryCurrencyPair(s.chain, ethusdcCP, height-1)
		s.Require().NoError(err)

		_, newNonce, err := QueryCurrencyPair(s.chain, ethusdcCP, height)
		s.Require().NoError(err)

		// expect no update for the height
		s.Require().Equal(newNonce, oldNonce)

		// start all oracles again
		for _, node := range s.chain.Nodes()[1:] {
			StartOracle(node)
		}
	})
}

func enabledTicker(pair connecttypes.CurrencyPair) mmtypes.Ticker {
	return mmtypes.Ticker{
		CurrencyPair:     pair,
		Decimals:         8,
		MinProviderCount: 1,
		Enabled:          true,
		Metadata_JSON:    "",
	}
}

func disabledTicker(pair connecttypes.CurrencyPair) mmtypes.Ticker {
	return mmtypes.Ticker{
		CurrencyPair:     pair,
		Decimals:         8,
		MinProviderCount: 1,
		Enabled:          false,
		Metadata_JSON:    "",
	}
}

func (s *ConnectOracleIntegrationSuite) TestMultiplePriceFeeds() {
	ethusdcCP := connecttypes.NewCurrencyPair("ETH", "USDC")
	ethusdtCP := connecttypes.NewCurrencyPair("ETH", "USDT")
	ethusdCP := connecttypes.NewCurrencyPair("ETH", "USD")

	// add multiple tickers
	tickers := []mmtypes.Ticker{
		enabledTicker(ethusdcCP),
		enabledTicker(ethusdtCP),
		enabledTicker(ethusdCP),
	}

	s.Require().NoError(s.AddCurrencyPairs(s.chain, s.user, 1.1, tickers...))

	cc, closeFn, err := GetChainGRPC(s.chain)
	s.Require().NoError(err)
	defer closeFn()

	// get the currency pair ids
	ctx := context.Background()
	id1, err := getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), ethusdcCP)
	s.Require().NoError(err)

	id2, err := getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), ethusdtCP)
	s.Require().NoError(err)

	id3, err := getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), ethusdCP)
	s.Require().NoError(err)

	zero := big.NewInt(0)
	zeroBz, err := zero.GobEncode()
	s.Require().NoError(err)

	// start all oracles
	for _, node := range s.chain.Nodes() {
		oracleConfig := DefaultOracleConfig(translateGRPCAddr(s.chain))
		oracleConfig.Providers[static.Name] = oracleconfig.ProviderConfig{
			Name: static.Name,
			API: oracleconfig.APIConfig{
				Enabled:          true,
				Timeout:          250 * time.Millisecond,
				Interval:         250 * time.Millisecond,
				ReconnectTimeout: 250 * time.Millisecond,
				MaxQueries:       1,
				Endpoints: []oracleconfig.Endpoint{
					{
						URL: "http://un-used-url.com",
					},
				},
				Atomic: true,
				Name:   static.Name,
			},
			Type: types.ConfigType,
		}

		oracle := GetOracleSideCar(node)
		SetOracleConfigsOnOracle(oracle, oracleConfig)
		s.Require().NoError(RestartOracle(node))
	}

	s.Run("all oracles running for multiple price feeds", func() {
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []connectabci.OracleVoteExtension{
			{
				Prices: map[uint64][]byte{
					id1: zeroBz,
					id2: zeroBz,
					id3: zeroBz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id1: zeroBz,
					id2: zeroBz,
					id3: zeroBz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id1: zeroBz,
					id2: zeroBz,
					id3: zeroBz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id1: zeroBz,
					id2: zeroBz,
					id3: zeroBz,
				},
			},
		})
		s.Require().NoError(err)

		// query for the given currency pair
		for _, ticker := range tickers {
			resp, _, err := QueryCurrencyPair(s.chain, ticker.CurrencyPair, height)
			s.Require().NoError(err)
			s.Require().Equal(int64(110000000), resp.Price.Int64())
		}
	})

	s.Run("all oracles running for multiple price feeds, except for one", func() {
		// stop first node's oracle, and update Prices to report
		node := s.chain.Nodes()[0]
		StopOracle(node)

		oracleConfig := DefaultOracleConfig(translateGRPCAddr(s.chain))

		// set only a provider (no marketmap)
		oracleConfig.Providers[static.Name] = oracleconfig.ProviderConfig{
			Name: static.Name,
			API: oracleconfig.APIConfig{
				Enabled:          true,
				Timeout:          250 * time.Millisecond,
				Interval:         250 * time.Millisecond,
				ReconnectTimeout: 250 * time.Millisecond,
				MaxQueries:       1,
				Endpoints: []oracleconfig.Endpoint{
					{
						URL: "http://un-used-url.com",
					},
				},
				Atomic: true,
				Name:   static.Name,
			},
			Type: types.ConfigType,
		}

		oracle := GetOracleSideCar(node)
		SetOracleConfigsOnOracle(oracle, oracleConfig)
		s.Require().NoError(RestartOracle(node))
		s.Require().NoError(RestartOracle(node))

		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []connectabci.OracleVoteExtension{
			{
				Prices: map[uint64][]byte{},
			},
			{
				Prices: map[uint64][]byte{
					id1: zeroBz,
					id2: zeroBz,
					id3: zeroBz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id1: zeroBz,
					id2: zeroBz,
					id3: zeroBz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id1: zeroBz,
					id2: zeroBz,
					id3: zeroBz,
				},
			},
		})
		s.Require().NoError(err)

		// query for the given currency pair
		for _, ticker := range tickers {
			resp, _, err := QueryCurrencyPair(s.chain, ticker.CurrencyPair, height)
			s.Require().NoError(err)
			s.Require().Equal(int64(110000000), resp.Price.Int64())
		}
	})
}

func getIDForCurrencyPair(ctx context.Context, client oracletypes.QueryClient, cp connecttypes.CurrencyPair) (uint64, error) {
	// query for the given currency pair
	resp, err := client.GetPrice(ctx, &oracletypes.GetPriceRequest{
		CurrencyPair: cp.String(),
	})
	if err != nil {
		return 0, err
	}

	return resp.Id, nil
}
