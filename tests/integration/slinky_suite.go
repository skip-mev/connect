package integration

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

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
			"oracle",
			"--oracle-config-path", "/oracle/oracle.toml",
			"--host", "0.0.0.0",
			"--port", "8080",
		},
		ValidatorProcess: true,
		PreStart:         true,
	}
}

func DefaultOracleConfig(node *cosmos.ChainNode) oracleconfig.OracleConfig {
	oracle := GetOracleSideCar(node)

	// Create the oracle config
	oracleConfig := oracleconfig.OracleConfig{
		UpdateInterval: 500 * time.Millisecond,
		Metrics: oracleconfig.MetricsConfig{
			Enabled:                 true,
			PrometheusServerAddress: fmt.Sprintf("%s:%s", oracle.HostName(), "8081"),
		},
	}

	return oracleConfig
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
			oracleCfg := DefaultOracleConfig(node)
			SetOracleConfigsOnOracle(GetOracleSideCar(node), oracleCfg)

			// set the out-of-process oracle config for all nodes
			node.WithPrestartNode(func(n *cosmos.ChainNode) {
				SetOracleConfigsOnApp(n, oracleCfg)
			})
		}
	})

	// start the chain
	BuildPOBInterchain(s.T(), context.Background(), s.chain)
	users := interchaintest.GetAndFundTestUsers(s.T(), context.Background(), s.T().Name(), math.NewInt(genesisAmount), s.chain)
	s.user = users[0]
}

func (s *SlinkyIntegrationSuite) TearDownSuite() {
	// get the oracle integration-test suite keep alive env
	if ok := os.Getenv(envKeepAlive); ok == "" {
		return
	}

	// keep the chain running
	s.T().Log("Keeping the chain running")
	for {
	}
}

func (s *SlinkyIntegrationSuite) SetupTest() {
	// query for all currency-pairs
	resp, err := QueryCurrencyPairs(s.chain)
	s.Require().NoError(err)

	s.T().Log("Removing all currency-pairs", resp.CurrencyPairs)

	// reset the oracle services
	// start all oracles
	for _, node := range s.chain.Nodes() {
		oCfg := DefaultOracleConfig(node)

		SetOracleConfigsOnOracle(GetOracleSideCar(node), oCfg)
		RestartOracle(node)
	}

	if len(resp.CurrencyPairs) == 0 {
		return
	}

	ids := make([]string, len(resp.CurrencyPairs))
	for i, cp := range resp.CurrencyPairs {
		ids[i] = cp.ToString()
	}

	// remove all currency-pairs
	s.Require().NoError(RemoveCurrencyPairs(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, ids...))
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
		s.Require().NoError(AddCurrencyPairs(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, []oracletypes.CurrencyPair{
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

	// remove the currency-pair from state and check the Prices for that currency-pair are no longer reported
	s.Run("Remove a currency-pair and check Prices", func() {
		s.Require().NoError(RemoveCurrencyPairs(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, []string{oracletypes.NewCurrencyPair("BTC", "USD").ToString()}...))

		// check that the currency-pair is added to state
		resp, err := QueryCurrencyPairs(s.chain)
		s.Require().NoError(err)
		s.Require().True(len(resp.CurrencyPairs) == 0)
	})

	s.Run("Add multiple Currency Pairs and remove 1", func() {
		cp1 := oracletypes.NewCurrencyPair("ETH", "USD")
		cp2 := oracletypes.NewCurrencyPair("BTC", "USD")
		s.Require().NoError(AddCurrencyPairs(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, []oracletypes.CurrencyPair{
			cp1, cp2,
		}...))

		resp, err := QueryCurrencyPairs(s.chain)
		s.Require().NoError(err)
		s.Require().True(len(resp.CurrencyPairs) == 2)

		// remove btc from state
		s.Require().NoError(RemoveCurrencyPairs(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, []string{cp2.ToString()}...))

		// check that the currency-pair is removed from state
		resp, err = QueryCurrencyPairs(s.chain)
		s.Require().NoError(err)
		s.Require().True(len(resp.CurrencyPairs) == 1)
		s.Require().Equal(resp.CurrencyPairs[0].Base, "ETH")
		s.Require().Equal(resp.CurrencyPairs[0].Quote, "USD")
	})
}

func (s *SlinkyOracleIntegrationSuite) TestNodeFailures() {
	cp := oracletypes.CurrencyPair{
		Base:  "ETHEREUM",
		Quote: "USDC",
	}

	s.Require().NoError(AddCurrencyPairs(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, []oracletypes.CurrencyPair{
		cp,
	}...))

	cc, close, err := GetChainGRPC(s.chain)
	s.Require().NoError(err)

	defer close()

	id, err := getIDForCurrencyPair(context.Background(), oracletypes.NewQueryClient(cc), cp)
	s.Require().NoError(err)

	zero := big.NewInt(0)
	zeroBz, err := zero.GobEncode()
	s.Require().NoError(err)

	// configure failing providers for various sets of nodes
	s.Run("all nodes report Prices", func() {
		// update all oracle configs
		for _, node := range s.chain.Nodes() {
			oracle := GetOracleSideCar(node)

			oracleConfig := DefaultOracleConfig(node)
			oracleConfig.Providers = append(oracleConfig.Providers, oracleconfig.ProviderConfig{
				Name: "static-mock-provider",
				API: oracleconfig.APIConfig{
					Enabled:    true,
					Timeout:    250 * time.Millisecond,
					Interval:   250 * time.Millisecond,
					MaxQueries: 1,
					URL:        "http://un-used-url.com",
					Atomic:     true,
					Name:       "static-mock-provider",
				},
				Market: oracleconfig.MarketConfig{
					Name: "static-mock-provider",
					CurrencyPairToMarketConfigs: map[string]oracleconfig.CurrencyPairMarketConfig{
						cp.ToString(): {
							Ticker:       "1140",
							CurrencyPair: cp,
						},
					},
				},
			})
			oracleConfig.CurrencyPairs = append(oracleConfig.CurrencyPairs, cp)

			SetOracleConfigsOnOracle(oracle, oracleConfig)
			RestartOracle(node)
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
		resp, _, err := QueryCurrencyPair(s.chain, cp, height)
		s.Require().NoError(err)
		s.Require().Equal(resp.Price.Int64(), int64(1140))
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

		_, oldNonce, err := QueryCurrencyPair(s.chain, cp, height-1)
		s.Require().NoError(err)

		_, newNonce, err := QueryCurrencyPair(s.chain, cp, height)
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

		_, oldNonce, err := QueryCurrencyPair(s.chain, cp, height-1)
		s.Require().NoError(err)

		_, newNonce, err := QueryCurrencyPair(s.chain, cp, height)
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

		_, oldNonce, err := QueryCurrencyPair(s.chain, cp, height-1)
		s.Require().NoError(err)

		_, newNonce, err := QueryCurrencyPair(s.chain, cp, height)
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
	cp1 := oracletypes.NewCurrencyPair("ETHEREUM", "USDC")
	cp2 := oracletypes.NewCurrencyPair("ETHEREUM", "USDT")
	cp3 := oracletypes.NewCurrencyPair("ETHEREUM", "USD")

	// add multiple currency pairs
	cps := []oracletypes.CurrencyPair{
		cp1,
		cp2,
		cp3,
	}

	s.Require().NoError(AddCurrencyPairs(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, cps...))

	cc, close, err := GetChainGRPC(s.chain)
	s.Require().NoError(err)

	defer close()

	// get the currency pair ids
	ctx := context.Background()
	id1, err := getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), cp1)
	s.Require().NoError(err)

	id2, err := getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), cp2)
	s.Require().NoError(err)

	id3, err := getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), cp3)
	s.Require().NoError(err)

	zero := big.NewInt(0)
	zeroBz, err := zero.GobEncode()
	s.Require().NoError(err)

	// start all oracles
	for _, node := range s.chain.Nodes() {
		oracle := GetOracleSideCar(node)

		oracleConfig := DefaultOracleConfig(node)

		oracleConfig.Providers = append(oracleConfig.Providers, oracleconfig.ProviderConfig{
			Name: "static-mock-provider",
			API: oracleconfig.APIConfig{
				Enabled:    true,
				Timeout:    250 * time.Millisecond,
				Interval:   250 * time.Millisecond,
				MaxQueries: 1,
				URL:        "http://un-used-url.com",
				Atomic:     true,
				Name:       "static-mock-provider",
			},
			Market: oracleconfig.MarketConfig{
				Name: "static-mock-provider",
				CurrencyPairToMarketConfigs: map[string]oracleconfig.CurrencyPairMarketConfig{
					cp1.ToString(): {
						Ticker:       "1140",
						CurrencyPair: cp1,
					},
					cp2.ToString(): {
						Ticker:       "1141",
						CurrencyPair: cp2,
					},
					cp3.ToString(): {
						Ticker:       "1142",
						CurrencyPair: cp3,
					},
				},
			},
		})

		oracleConfig.CurrencyPairs = append(oracleConfig.CurrencyPairs, cp1)
		oracleConfig.CurrencyPairs = append(oracleConfig.CurrencyPairs, cp2)
		oracleConfig.CurrencyPairs = append(oracleConfig.CurrencyPairs, cp3)

		SetOracleConfigsOnOracle(oracle, oracleConfig)
		RestartOracle(node)
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
		for i, cp := range cps {
			resp, _, err := QueryCurrencyPair(s.chain, cp, height)
			s.Require().NoError(err)
			s.Require().Equal(int64(1140+i), resp.Price.Int64())
		}
	})

	s.Run("all oracles running for multiple price feeds, except for one", func() {
		// stop first node's oracle, and update Prices to report
		node := s.chain.Nodes()[0]
		StopOracle(node)

		oracle := GetOracleSideCar(node)

		oracleConfig := DefaultOracleConfig(node)

		oracleConfig.Providers = append(oracleConfig.Providers, oracleconfig.ProviderConfig{
			Name: "static-mock-provider",
			API: oracleconfig.APIConfig{
				Enabled:    true,
				Timeout:    250 * time.Millisecond,
				Interval:   250 * time.Millisecond,
				MaxQueries: 1,
				URL:        "http://un-used-url.com",
				Atomic:     true,
				Name:       "static-mock-provider",
			},
			Market: oracleconfig.MarketConfig{
				Name: "static-mock-provider",
				CurrencyPairToMarketConfigs: map[string]oracleconfig.CurrencyPairMarketConfig{
					cp1.ToString(): {
						Ticker:       "1140",
						CurrencyPair: cp1,
					},
					cp2.ToString(): {
						Ticker:       "1141",
						CurrencyPair: cp2,
					},
				},
			},
		})

		oracleConfig.CurrencyPairs = append(oracleConfig.CurrencyPairs, cp1)
		oracleConfig.CurrencyPairs = append(oracleConfig.CurrencyPairs, cp2)
		oracleConfig.CurrencyPairs = append(oracleConfig.CurrencyPairs, cp3)

		SetOracleConfigsOnOracle(oracle, oracleConfig)
		RestartOracle(node)

		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
			{
				Prices: map[uint64][]byte{
					id1: zeroBz,
					id2: zeroBz,
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
		for i, cp := range cps {
			resp, _, err := QueryCurrencyPair(s.chain, cp, height)
			s.Require().NoError(err)
			s.Require().Equal(int64(1140+i), resp.Price.Int64())
		}
	})
}

func getIDForCurrencyPair(ctx context.Context, client oracletypes.QueryClient, cp oracletypes.CurrencyPair) (uint64, error) {
	// query for the given currency pair
	resp, err := client.GetPrice(ctx, &oracletypes.GetPriceRequest{
		CurrencyPairSelector: &oracletypes.GetPriceRequest_CurrencyPair{
			CurrencyPair: &cp,
		},
	})
	if err != nil {
		return 0, err
	}

	return resp.Id, nil
}
