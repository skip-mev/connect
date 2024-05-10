package integration

import (
	"context"
	"math/big"
	"os"
	"time"

	"github.com/skip-mev/slinky/providers/apis/marketmap"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/stretchr/testify/suite"

	slinkyabci "github.com/skip-mev/slinky/abci/ve/types"
	oracleconfig "github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/static"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	envKeepAlive  = "ORACLE_INTEGRATION_KEEPALIVE"
	genesisAmount = 1000000000
	defaultDenom  = "stake"
	validatorKey  = "validator"
	yes           = "yes"
	deposit       = 1000000
)

func DefaultOracleSidecar(image ibc.DockerImage) ibc.SidecarConfig {
	return ibc.SidecarConfig{
		ProcessName: "oracle",
		Image:       image,
		HomeDir:     "/oracle",
		Ports:       []string{"8080", "8081"},
		StartCmd: []string{
			"slinky",
			"--oracle-config-path", "/oracle/oracle.json",
		},
		ValidatorProcess: true,
		PreStart:         true,
	}
}

func DefaultOracleConfig(url string) oracleconfig.OracleConfig {
	cfg := marketmap.DefaultAPIConfig
	cfg.URL = url

	// Create the oracle config
	oracleConfig := oracleconfig.OracleConfig{
		UpdateInterval: 500 * time.Millisecond,
		MaxPriceAge:    1 * time.Minute,
		Host:           "0.0.0.0",
		Port:           "8080",
		Providers: []oracleconfig.ProviderConfig{
			{
				Name: "marketmap_api",
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
}

func NewSlinkyIntegrationSuite(spec *interchaintest.ChainSpec, oracleImage ibc.DockerImage) *SlinkyIntegrationSuite {
	return &SlinkyIntegrationSuite{
		spec:         spec,
		oracleConfig: DefaultOracleSidecar(oracleImage),
		denom:        defaultDenom,
		authority:    authtypes.NewModuleAddress(govtypes.ModuleName),
		blockTime:    10 * time.Second,
	}
}

func (s *SlinkyIntegrationSuite) WithDenom(denom string) *SlinkyIntegrationSuite {
	s.denom = denom
	return s
}

func (s *SlinkyIntegrationSuite) WithAuthority(addr sdk.AccAddress) *SlinkyIntegrationSuite {
	s.authority = addr
	return s
}

func (s *SlinkyIntegrationSuite) WithBlockTime(t time.Duration) *SlinkyIntegrationSuite {
	s.blockTime = t
	return s
}

func (s *SlinkyIntegrationSuite) SetupSuite() {
	// create the chain
	s.chain = ChainBuilderFromChainSpec(s.T(), s.spec)

	s.chain.WithPrestartNodes(func(c *cosmos.CosmosChain) {
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
			node.WithPrestartNode(func(n *cosmos.ChainNode) {
				SetOracleConfigsOnApp(n)
			})
		}
	})

	// start the chain
	BuildPOBInterchain(s.T(), context.Background(), s.chain)
	users := interchaintest.GetAndFundTestUsers(s.T(), context.Background(), s.T().Name(), math.NewInt(genesisAmount), s.chain)
	s.user = users[0]

	resp, err := UpdateMarketMapParams(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, mmtypes.Params{
		MarketAuthorities: []string{s.user.FormattedAddress()},
		Admin:             s.user.FormattedAddress(),
	})
	s.Require().NoError(err, resp)
}

func (s *SlinkyIntegrationSuite) TearDownSuite() {
	// get the oracle integration-test suite keep alive env
	if ok := os.Getenv(envKeepAlive); ok == "" {
		return
	}

	// keep the chain running
	s.T().Log("Keeping the chain running")
	select {}
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

type SlinkyOracleIntegrationSuite struct {
	*SlinkyIntegrationSuite
}

func NewSlinkyOracleIntegrationSuite(suite *SlinkyIntegrationSuite) *SlinkyOracleIntegrationSuite {
	return &SlinkyOracleIntegrationSuite{
		SlinkyIntegrationSuite: suite,
	}
}

func (s *SlinkyOracleIntegrationSuite) TestOracleModule() {
	// query the oracle module grpc service for any CurrencyPairs
	s.Run("QueryCurrencyPairs - no currency-pairs reported", func() {
		resp, err := QueryCurrencyPairs(s.chain)
		s.Require().NoError(err)
		s.Require().True(len(resp.CurrencyPairs) == 0)
	})

	// pass a governance proposal to approve a new currency-pair, and check Prices are reported
	s.Run("Add a currency-pair and check Prices", func() {
		s.Require().NoError(s.AddCurrencyPairs(s.chain, s.user, 1.1, []slinkytypes.CurrencyPair{
			{
				Base:  "BTC",
				Quote: "USD",
			},
		}...))

		// check that the currency-pair is added to state
		resp, err := QueryCurrencyPairs(s.chain)
		s.Require().NoError(err)
		s.Require().True(len(resp.CurrencyPairs) == 1)
		s.Require().Equal(resp.CurrencyPairs[0].Base, "BTC")
		s.Require().Equal(resp.CurrencyPairs[0].Quote, "USD")
	})

	s.Run("Add multiple Currency Pairs", func() {
		cp1 := slinkytypes.NewCurrencyPair("ETH", "USD")
		cp2 := slinkytypes.NewCurrencyPair("USDT", "USD")
		s.Require().NoError(s.AddCurrencyPairs(s.chain, s.user, 1.1, []slinkytypes.CurrencyPair{
			cp1, cp2,
		}...))

		resp, err := QueryCurrencyPairs(s.chain)
		s.Require().NoError(err)
		s.Require().True(len(resp.CurrencyPairs) == 3)
	})
}

func translateGRPCAddr(chain *cosmos.CosmosChain) string {
	return chain.GetGRPCAddress()
}

func (s *SlinkyOracleIntegrationSuite) TestNodeFailures() {
	ethusdc := constants.ETHEREUM_USDC

	s.Require().NoError(s.AddCurrencyPairs(s.chain, s.user, 1.1, []slinkytypes.CurrencyPair{
		ethusdc,
	}...))

	cc, closeFn, err := GetChainGRPC(s.chain)
	s.Require().NoError(err)
	defer closeFn()

	id, err := getIDForCurrencyPair(context.Background(), oracletypes.NewQueryClient(cc), ethusdc)
	s.Require().NoError(err)

	zero := big.NewInt(0)
	zeroBz, err := zero.GobEncode()
	s.Require().NoError(err)

	// configure failing providers for various sets of nodes
	s.Run("all nodes report Prices", func() {
		// update all oracle configs
		for _, node := range s.chain.Nodes() {
			oracleConfig := DefaultOracleConfig(translateGRPCAddr(s.chain))
			oracleConfig.Providers = append(oracleConfig.Providers, oracleconfig.ProviderConfig{
				Name: static.Name,
				API: oracleconfig.APIConfig{
					Enabled:          true,
					Timeout:          250 * time.Millisecond,
					Interval:         250 * time.Millisecond,
					ReconnectTimeout: 250 * time.Millisecond,
					MaxQueries:       1,
					URL:              "http://un-used-url.com",
					Atomic:           true,
					Name:             static.Name,
				},
				Type: types.ConfigType,
			})

			oracle := GetOracleSideCar(node)
			SetOracleConfigsOnOracle(oracle, oracleConfig)
			s.Require().NoError(RestartOracle(node))
		}

		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
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
		resp, _, err := QueryCurrencyPair(s.chain, ethusdc, height)
		s.Require().NoError(err)
		s.Require().Equal(resp.Price.Int64(), int64(110000000))
	})

	s.Run("single oracle down, price updates", func() {
		// stop single node's oracle process and check that all Prices are reported
		node := s.chain.Nodes()[0]
		StopOracle(node)

		// expect the following vote-extensions
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
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

		_, oldNonce, err := QueryCurrencyPair(s.chain, ethusdc, height-1)
		s.Require().NoError(err)

		_, newNonce, err := QueryCurrencyPair(s.chain, ethusdc, height)
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
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
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

		_, oldNonce, err := QueryCurrencyPair(s.chain, ethusdc, height-1)
		s.Require().NoError(err)

		_, newNonce, err := QueryCurrencyPair(s.chain, ethusdc, height)
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
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
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

		_, oldNonce, err := QueryCurrencyPair(s.chain, ethusdc, height-1)
		s.Require().NoError(err)

		_, newNonce, err := QueryCurrencyPair(s.chain, ethusdc, height)
		s.Require().NoError(err)

		// expect no update for the height
		s.Require().Equal(newNonce, oldNonce)

		// start all oracles again
		for _, node := range s.chain.Nodes()[1:] {
			StartOracle(node)
		}
	})
}

func (s *SlinkyOracleIntegrationSuite) TestMultiplePriceFeeds() {
	ethusdc := constants.ETHEREUM_USDC
	ethusdt := constants.ETHEREUM_USDT
	ethusd := constants.ETHEREUM_USD

	// add multiple currency pairs
	cps := []slinkytypes.CurrencyPair{
		ethusdc,
		ethusdt,
		ethusd,
	}

	s.Require().NoError(s.AddCurrencyPairs(s.chain, s.user, 1.1, cps...))

	cc, closeFn, err := GetChainGRPC(s.chain)
	s.Require().NoError(err)
	defer closeFn()

	// get the currency pair ids
	ctx := context.Background()
	id1, err := getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), ethusdc)
	s.Require().NoError(err)

	id2, err := getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), ethusdt)
	s.Require().NoError(err)

	id3, err := getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), ethusd)
	s.Require().NoError(err)

	zero := big.NewInt(0)
	zeroBz, err := zero.GobEncode()
	s.Require().NoError(err)

	// start all oracles
	for _, node := range s.chain.Nodes() {
		oracleConfig := DefaultOracleConfig(translateGRPCAddr(s.chain))
		oracleConfig.Providers = append(oracleConfig.Providers, oracleconfig.ProviderConfig{
			Name: static.Name,
			API: oracleconfig.APIConfig{
				Enabled:          true,
				Timeout:          250 * time.Millisecond,
				Interval:         250 * time.Millisecond,
				ReconnectTimeout: 250 * time.Millisecond,
				MaxQueries:       1,
				URL:              "http://un-used-url.com",
				Atomic:           true,
				Name:             static.Name,
			},
			Type: types.ConfigType,
		})

		oracle := GetOracleSideCar(node)
		SetOracleConfigsOnOracle(oracle, oracleConfig)
		s.Require().NoError(RestartOracle(node))
	}

	s.Run("all oracles running for multiple price feeds", func() {
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
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
		for _, cp := range cps {
			resp, _, err := QueryCurrencyPair(s.chain, cp, height)
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
		oracleConfig.Providers = []oracleconfig.ProviderConfig{{
			Name: static.Name,
			API: oracleconfig.APIConfig{
				Enabled:          true,
				Timeout:          250 * time.Millisecond,
				Interval:         250 * time.Millisecond,
				ReconnectTimeout: 250 * time.Millisecond,
				MaxQueries:       1,
				URL:              "http://un-used-url.com",
				Atomic:           true,
				Name:             static.Name,
			},
			Type: types.ConfigType,
		}}

		oracle := GetOracleSideCar(node)
		SetOracleConfigsOnOracle(oracle, oracleConfig)
		s.Require().NoError(RestartOracle(node))
		s.Require().NoError(RestartOracle(node))

		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
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
		for _, cp := range cps {
			resp, _, err := QueryCurrencyPair(s.chain, cp, height)
			s.Require().NoError(err)
			s.Require().Equal(int64(110000000), resp.Price.Int64())
		}
	})
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
