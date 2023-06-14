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

// Spin up a local oracle and make requests to it.
//
// Meant only for testing purposes.
func main() {
	// Oracle config
	logger := NewLogger()
	providerTimeout := 2 * time.Second
	oracleTicker := 5 * time.Second

	// Currency pairs each of the providers will
	// be fetching prices for.
	currencyPairs := []types.CurrencyPair{
		{Base: "BITCOIN", Quote: "USD"},
		{Base: "ETHEREUM", Quote: "USD"},
		{Base: "COSMOS", Quote: "USD"},
		{Base: "POLKADOT", Quote: "USD"},
		{Base: "POLYGON", Quote: "USD"},
	}

	// Define the providers
	providers := []types.Provider{
		mock.NewTimeoutMockProvider(oracleTicker),
		mock.NewFailingMockProvider(),
		coinbase.NewProvider(logger, currencyPairs),
		coingecko.NewProvider(logger, currencyPairs),
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
	go func() {
		if err := client.Start(context.Background()); err != nil {
			panic(err)
		}
	}()

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
