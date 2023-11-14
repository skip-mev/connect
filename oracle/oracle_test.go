package oracle_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/holiman/uint256"
	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/mocks"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type OracleTestSuite struct {
	suite.Suite
	random *rand.Rand

	// Oracle config
	oracle        *oracle.Oracle
	oracleTicker  time.Duration
	providers     []*mocks.Provider
	currencyPairs []oracletypes.CurrencyPair
	aggregationFn aggregator.AggregateFn
	ctx           context.Context
}

func TestOracleSuite(t *testing.T) {
	suite.Run(t, new(OracleTestSuite))
}

func (suite *OracleTestSuite) SetupTest() {
	suite.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	// Oracle set up
	suite.oracleTicker = 1 * time.Second
	suite.currencyPairs = []oracletypes.CurrencyPair{
		oracletypes.NewCurrencyPair("BITCOIN", "USD"),
		oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
		oracletypes.NewCurrencyPair("COSMOS", "USD"),
	}
	suite.aggregationFn = aggregator.ComputeMedian()

	suite.ctx = context.TODO()
}

func (suite *OracleTestSuite) TestProviders() {
	cases := []struct {
		name        string
		fetchPrices func() map[oracletypes.CurrencyPair]*uint256.Int
	}{
		{
			name: "no providers",
			fetchPrices: func() map[oracletypes.CurrencyPair]*uint256.Int {
				suite.providers = []*mocks.Provider{}

				return map[oracletypes.CurrencyPair]*uint256.Int{}
			},
		},
		{
			name: "one provider with no prices",
			fetchPrices: func() map[oracletypes.CurrencyPair]*uint256.Int {
				staticProvider := suite.createStaticProvider(
					"static",
					map[oracletypes.CurrencyPair]aggregator.QuotePrice{},
				)

				suite.providers = []*mocks.Provider{
					staticProvider,
				}

				return map[oracletypes.CurrencyPair]*uint256.Int{}
			},
		},
		{
			name: "one provider with static prices",
			fetchPrices: func() map[oracletypes.CurrencyPair]*uint256.Int {
				staticProvider := suite.createStaticProvider(
					"static",
					map[oracletypes.CurrencyPair]aggregator.QuotePrice{
						suite.currencyPairs[0]: {
							Price:     uint256.NewInt(1),
							Timestamp: time.Now(),
						},
						suite.currencyPairs[1]: {
							Price:     uint256.NewInt(2),
							Timestamp: time.Now(),
						},
						suite.currencyPairs[2]: {
							Price:     uint256.NewInt(3),
							Timestamp: time.Now(),
						},
					},
				)

				suite.providers = []*mocks.Provider{
					staticProvider,
				}

				return map[oracletypes.CurrencyPair]*uint256.Int{
					suite.currencyPairs[0]: uint256.NewInt(1),
					suite.currencyPairs[1]: uint256.NewInt(2),
					suite.currencyPairs[2]: uint256.NewInt(3),
				}
			},
		},
		{
			name: "two providers with static prices",
			fetchPrices: func() map[oracletypes.CurrencyPair]*uint256.Int {
				staticProvider1 := suite.createStaticProvider(
					"static1",
					map[oracletypes.CurrencyPair]aggregator.QuotePrice{
						suite.currencyPairs[0]: {
							Price:     uint256.NewInt(1),
							Timestamp: time.Now(),
						},
						suite.currencyPairs[1]: {
							Price:     uint256.NewInt(2),
							Timestamp: time.Now(),
						},
						suite.currencyPairs[2]: {
							Price:     uint256.NewInt(3),
							Timestamp: time.Now(),
						},
					},
				)

				staticProvider2 := suite.createStaticProvider(
					"static2",
					map[oracletypes.CurrencyPair]aggregator.QuotePrice{
						suite.currencyPairs[0]: {
							Price:     uint256.NewInt(4),
							Timestamp: time.Now(),
						},
						suite.currencyPairs[1]: {
							Price:     uint256.NewInt(5),
							Timestamp: time.Now(),
						},
						suite.currencyPairs[2]: {
							Price:     uint256.NewInt(6),
							Timestamp: time.Now(),
						},
					},
				)

				suite.providers = []*mocks.Provider{
					staticProvider1,
					staticProvider2,
				}

				return map[oracletypes.CurrencyPair]*uint256.Int{
					suite.currencyPairs[0]: uint256.NewInt(2),
					suite.currencyPairs[1]: uint256.NewInt(3),
					suite.currencyPairs[2]: uint256.NewInt(4),
				}
			},
		},
		{
			name: "one provider with randomized prices",
			fetchPrices: func() map[oracletypes.CurrencyPair]*uint256.Int {
				randomizedProvider := suite.createRandomizedProvider(
					"randomized",
					suite.currencyPairs,
				)

				suite.providers = []*mocks.Provider{
					randomizedProvider,
				}

				return suite.aggregateProviderData(suite.providers)
			},
		},
		{
			name: "two providers with randomized prices",
			fetchPrices: func() map[oracletypes.CurrencyPair]*uint256.Int {
				randomizedProvider1 := suite.createRandomizedProvider(
					"randomized1",
					suite.currencyPairs,
				)

				randomizedProvider2 := suite.createRandomizedProvider(
					"randomized2",
					suite.currencyPairs,
				)

				suite.providers = []*mocks.Provider{
					randomizedProvider1,
					randomizedProvider2,
				}

				return suite.aggregateProviderData(suite.providers)
			},
		},
		{
			name: "one normal static provider and one panic provider",
			fetchPrices: func() map[oracletypes.CurrencyPair]*uint256.Int {
				staticProvider := suite.createStaticProvider(
					"static",
					map[oracletypes.CurrencyPair]aggregator.QuotePrice{
						suite.currencyPairs[0]: {
							Price:     uint256.NewInt(1),
							Timestamp: time.Now(),
						},
						suite.currencyPairs[1]: {
							Price:     uint256.NewInt(2),
							Timestamp: time.Now(),
						},
						suite.currencyPairs[2]: {
							Price:     uint256.NewInt(3),
							Timestamp: time.Now(),
						},
					},
				)

				panicProvider := suite.createPanicProvider(
					"panic",
				)

				suite.providers = []*mocks.Provider{
					staticProvider,
					panicProvider,
				}

				return map[oracletypes.CurrencyPair]*uint256.Int{
					suite.currencyPairs[0]: uint256.NewInt(1),
					suite.currencyPairs[1]: uint256.NewInt(2),
					suite.currencyPairs[2]: uint256.NewInt(3),
				}
			},
		},
		{
			name: "one normal static provider and one timeout provider",
			fetchPrices: func() map[oracletypes.CurrencyPair]*uint256.Int {
				staticProvider := suite.createStaticProvider(
					"static",
					map[oracletypes.CurrencyPair]aggregator.QuotePrice{
						suite.currencyPairs[0]: {
							Price:     uint256.NewInt(1),
							Timestamp: time.Now(),
						},
						suite.currencyPairs[1]: {
							Price:     uint256.NewInt(2),
							Timestamp: time.Now(),
						},
						suite.currencyPairs[2]: {
							Price:     uint256.NewInt(3),
							Timestamp: time.Now(),
						},
					},
				)

				timeoutProvider := suite.createTimeoutProvider(
					"timeout",
				)

				suite.providers = []*mocks.Provider{
					staticProvider,
					timeoutProvider,
				}

				return map[oracletypes.CurrencyPair]*uint256.Int{
					suite.currencyPairs[0]: uint256.NewInt(1),
					suite.currencyPairs[1]: uint256.NewInt(2),
					suite.currencyPairs[2]: uint256.NewInt(3),
				}
			},
		},
		{
			name: "multiple random providers with one timeout and panic provider",
			fetchPrices: func() map[oracletypes.CurrencyPair]*uint256.Int {
				randomizedProvider1 := suite.createRandomizedProvider(
					"randomized1",
					suite.currencyPairs,
				)

				randomizedProvider2 := suite.createRandomizedProvider(
					"randomized2",
					suite.currencyPairs,
				)

				timeoutProvider := suite.createTimeoutProvider(
					"timeout",
				)

				panicProvider := suite.createPanicProvider(
					"panic",
				)

				suite.providers = []*mocks.Provider{
					randomizedProvider1,
					randomizedProvider2,
					timeoutProvider,
					panicProvider,
				}

				return suite.aggregateProviderData(
					[]*mocks.Provider{
						randomizedProvider1,
						randomizedProvider2,
					},
				)
			},
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			// Reset oracle
			suite.SetupTest()

			expectedPrices := tc.fetchPrices()

			oracleConfig := config.OracleConfig{
				InProcess:      true,
				RemoteAddress:  "",
				UpdateInterval: suite.oracleTicker,
			}

			tempProviders := make([]oracle.Provider, len(suite.providers))
			for i, provider := range suite.providers {
				tempProviders[i] = provider
			}
			factory := func(log.Logger, config.OracleConfig) ([]oracle.Provider, error) {
				return tempProviders, nil
			}

			// Initialize oracle
			var err error
			suite.oracle, err = oracle.New(
				log.NewTestLogger(suite.T()),
				oracleConfig,
				factory,
				suite.aggregationFn,
				nil,
			)
			suite.Require().NoError(err)

			// Start oracle
			go func() {
				suite.oracle.Start(suite.ctx)
			}()

			// Wait for oracle to update prices
			time.Sleep(suite.oracleTicker * 2)
			suite.oracle.Stop()
			time.Sleep(suite.oracleTicker * 2)

			// Check prices
			prices := suite.oracle.GetPrices()
			for pair, price := range expectedPrices {
				suite.Require().Contains(prices, pair)

				suite.Require().Equal(
					price,
					prices[pair],
				)
			}

			// Check oracle status
			suite.Require().Eventually(
				func() bool {
					return !suite.oracle.IsRunning()
				},
				5*suite.oracleTicker,
				suite.oracleTicker/3,
			)
		})
	}
}

func (suite *OracleTestSuite) TestShutDownWithContextCancel() {
	suite.SetupTest()

	oracleConfig := config.OracleConfig{
		InProcess:      true,
		RemoteAddress:  "",
		UpdateInterval: suite.oracleTicker,
	}
	factory := func(log.Logger, config.OracleConfig) ([]oracle.Provider, error) {
		return []oracle.Provider{
			suite.createStaticProvider(
				"static",
				map[oracletypes.CurrencyPair]aggregator.QuotePrice{},
			),
		}, nil
	}

	// Initialize oracle
	var err error
	suite.oracle, err = oracle.New(
		log.NewTestLogger(suite.T()),
		oracleConfig,
		factory,
		suite.aggregationFn,
		nil,
	)
	suite.Require().NoError(err)

	ctx, cancel := context.WithCancel(suite.ctx)

	// Start oracle
	go func() {
		suite.oracle.Start(ctx)
	}()

	// Wait for oracle to update prices
	time.Sleep(suite.oracleTicker * 2)
	cancel()
	time.Sleep(suite.oracleTicker * 2)

	// Check prices
	prices := suite.oracle.GetPrices()
	suite.Require().Empty(prices)

	// Check oracle status
	suite.Require().Eventually(
		func() bool {
			return !suite.oracle.IsRunning()
		},
		5*suite.oracleTicker,
		suite.oracleTicker/3,
	)
}

func (suite *OracleTestSuite) TestShutDownWithStop() {
	suite.SetupTest()

	oracleConfig := config.OracleConfig{
		InProcess:      true,
		RemoteAddress:  "",
		UpdateInterval: suite.oracleTicker,
	}
	factory := func(log.Logger, config.OracleConfig) ([]oracle.Provider, error) {
		return []oracle.Provider{
			suite.createStaticProvider(
				"static",
				map[oracletypes.CurrencyPair]aggregator.QuotePrice{},
			),
		}, nil
	}

	// Initialize oracle
	var err error
	suite.oracle, err = oracle.New(
		log.NewTestLogger(suite.T()),
		oracleConfig,
		factory,
		suite.aggregationFn,
		nil,
	)
	suite.Require().NoError(err)

	// Start oracle
	go func() {
		suite.oracle.Start(suite.ctx)
	}()

	// Wait for oracle to update prices
	time.Sleep(suite.oracleTicker * 2)
	suite.oracle.Stop()
	time.Sleep(suite.oracleTicker * 2)

	// Check prices
	prices := suite.oracle.GetPrices()
	suite.Require().Empty(prices)

	// Check oracle status
	suite.Require().Eventually(
		func() bool {
			return !suite.oracle.IsRunning()
		},
		5*suite.oracleTicker,
		suite.oracleTicker/3,
	)
}

func (suite *OracleTestSuite) TestShutDownProviderWithTimeout() {
	suite.SetupTest()

	tempProviders := []oracle.Provider{
		suite.createTimeoutProviderWithTimeout(
			"timeout",
			suite.oracleTicker*40,
			map[oracletypes.CurrencyPair]aggregator.QuotePrice{
				suite.currencyPairs[0]: {
					Price:     uint256.NewInt(1),
					Timestamp: time.Now(),
				},
			},
		),
	}

	oracleConfig := config.OracleConfig{
		InProcess:      true,
		RemoteAddress:  "",
		UpdateInterval: suite.oracleTicker,
	}
	factory := func(log.Logger, config.OracleConfig) ([]oracle.Provider, error) {
		return tempProviders, nil
	}

	// Initialize oracle
	var err error
	suite.oracle, err = oracle.New(
		log.NewTestLogger(suite.T()),
		oracleConfig,
		factory,
		suite.aggregationFn,
		nil,
	)
	suite.Require().NoError(err)

	// Start oracle
	go func() {
		suite.oracle.Start(suite.ctx)
	}()

	// Wait for oracle to update prices
	time.Sleep(suite.oracleTicker * 2)
	suite.oracle.Stop()
	time.Sleep(suite.oracleTicker * 3)

	// Check prices
	prices := suite.oracle.GetPrices()
	suite.Require().Empty(prices)

	// Check oracle status
	suite.Require().Eventually(
		func() bool {
			return !suite.oracle.IsRunning()
		},
		5*suite.oracleTicker,
		suite.oracleTicker/3,
	)
}

func (suite *OracleTestSuite) createTimeoutProviderWithTimeout(
	name string,
	timeout time.Duration,
	prices map[oracletypes.CurrencyPair]aggregator.QuotePrice,
) *mocks.Provider {
	provider := mocks.NewProvider(suite.T())
	provider.On("Name").Return(name)

	// GetPrices returns a timeout error
	provider.On("GetPrices", mock.Anything).Return(
		prices,
		nil,
	).After(timeout)

	return provider
}

func (suite *OracleTestSuite) createTimeoutProvider(
	name string,
) *mocks.Provider {
	provider := mocks.NewProvider(suite.T())
	provider.On("Name").Return(name)

	// GetPrices returns a timeout error
	provider.On("GetPrices", mock.Anything).Return(
		nil,
		fmt.Errorf("timeout error"),
	).After(suite.oracleTicker)

	return provider
}

func (suite *OracleTestSuite) createPanicProvider(
	name string,
) *mocks.Provider {
	provider := mocks.NewProvider(suite.T())
	provider.On("Name").Return(name)

	// GetPrices returns a timeout error
	provider.On("GetPrices", mock.Anything).Panic("not implemented")

	return provider
}

func (suite *OracleTestSuite) createStaticProvider(
	name string,
	prices map[oracletypes.CurrencyPair]aggregator.QuotePrice,
) *mocks.Provider {
	provider := mocks.NewProvider(suite.T())
	provider.On("Name").Return(name)

	// GetPrices returns static prices
	provider.On("GetPrices", mock.Anything).Return(
		prices,
		nil,
	)

	return provider
}

func (suite *OracleTestSuite) createRandomizedProvider(
	name string,
	currencyPairs []oracletypes.CurrencyPair,
) *mocks.Provider {
	provider := mocks.NewProvider(suite.T())
	provider.On("Name").Return(name)

	// GetPrices returns randomized prices
	provider.On("GetPrices", mock.Anything).Return(
		suite.getRandomizedPrices(currencyPairs),
		nil,
	)

	return provider
}

func (suite *OracleTestSuite) getRandomizedPrices(
	currencyPairs []oracletypes.CurrencyPair,
) map[oracletypes.CurrencyPair]aggregator.QuotePrice {
	prices := make(map[oracletypes.CurrencyPair]aggregator.QuotePrice)

	for _, pair := range currencyPairs {
		price := suite.random.Uint64()
		prices[pair] = aggregator.QuotePrice{
			Price:     uint256.NewInt(price),
			Timestamp: time.Now(),
		}
	}

	return prices
}

func (suite *OracleTestSuite) aggregateProviderData(
	providers []*mocks.Provider,
) map[oracletypes.CurrencyPair]*uint256.Int {
	// Aggregate prices from all providers
	priceAggregator := aggregator.NewPriceAggregator(suite.aggregationFn)

	for _, provider := range providers {
		prices, err := provider.GetPrices(context.Background())
		suite.Require().NoError(err)

		priceAggregator.SetProviderPrices(provider.Name(), prices)
	}

	return priceAggregator.GetPrices()
}
