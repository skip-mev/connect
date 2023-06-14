package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/coinbase"
	"github.com/skip-mev/slinky/providers/coingecko"
	"github.com/skip-mev/slinky/providers/mock"
	"github.com/skip-mev/slinky/service"
	"github.com/skip-mev/slinky/service/client"
)

func main() {
	// Oracle config
	logger := NewLogger()
	providerTimeout := 5 * time.Second
	oracleTicker := 5 * time.Second

	currencyPairs := []types.CurrencyPair{
		{Base: "BITCOIN", Quote: "USD"},
		{Base: "ETHEREUM", Quote: "USD"},
		{Base: "COSMOS", Quote: "USD"},
	}

	// Coinbase provider
	coinbaseProvider := coinbase.NewProvider(
		logger,
		currencyPairs,
	)

	// CoinGecko provider
	coingeckoProvider := coingecko.NewProvider(
		logger,
		currencyPairs,
	)

	// Mock providers for testing
	// mockProvider := mock.NewMockProvider()
	faillingMockProvider := mock.NewFailingMockProvider()
	timeoutMockProvider := mock.NewTimeoutMockProvider(oracleTicker)

	// Define the providers
	providers := []types.Provider{
		// mockProvider,
		timeoutMockProvider,
		faillingMockProvider,
		coinbaseProvider,
		coingeckoProvider,
	}

	// Initializing the oracle
	oracle := oracle.New(
		logger,
		providerTimeout,
		oracleTicker,
		providers,
		types.ComputeMedian(),
	)

	// Client set up and start
	client := client.NewLocalClient(oracle)

	// Start the oracle
	go func() {
		if err := client.Start(context.Background()); err != nil {
			panic(err)
		}
	}()

	// Wait for the oracle to start
	for !oracle.IsRunning() {
		time.Sleep(100 * time.Millisecond)
	}

	// Start up a local client that makes requests to the oracle
	for {
		time.Sleep(10 * time.Second)

		// Get the latest prices
		prices, err := client.Prices(context.Background(), &service.QueryPricesRequest{})
		if err != nil {
			logger.Error("failed to get prices", err)
			continue
		}

		// Print the prices
		responseString := "\n\tPrices: \n"
		for pair, price := range prices.Prices {
			responseString += fmt.Sprintf("\t\t%s: %s\n", pair, price)
		}

		fmt.Println(responseString)
	}
}

type Logger struct{}

var _ log.Logger = Logger{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l Logger) Info(key string, msgs ...interface{}) {
	var msg string
	for _, m := range msgs {
		msg += fmt.Sprintf("%s ", m)
	}

	msg = fmt.Sprintf("Info: %s %s", key, msg)
	fmt.Println(msg)
}

func (l Logger) Debug(key string, msgs ...interface{}) {
	var msg string
	for _, m := range msgs {
		msg += fmt.Sprintf("%s ", m)
	}

	msg = fmt.Sprintf("Info: %s %s", key, msg)
	fmt.Println(msg)
}

func (l Logger) Error(key string, msgs ...interface{}) {
	var msg string
	for _, m := range msgs {
		msg += fmt.Sprintf("%s ", m)
	}

	msg = fmt.Sprintf("Info: %s %s", key, msg)
	fmt.Println(msg)
}

func (l Logger) With(msgs ...interface{}) log.Logger {
	return Logger{}
}
