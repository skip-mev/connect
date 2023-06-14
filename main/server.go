package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/oracle/utils"
	"github.com/skip-mev/slinky/providers/coinbase"
	"github.com/skip-mev/slinky/providers/mock"
	"github.com/skip-mev/slinky/service/client"
)

func main() {
	// Oracle config
	logger := NewLogger()
	providerTimeout := 5 * time.Second
	oracleTicker := 5 * time.Second

	// Coinbase provider
	coinbaseProvider := coinbase.NewProvider(logger)
	coinbaseProvider.SetPairs(
		types.CurrencyPair{Base: "BTC", Quote: "USD"},
	)

	// Mock providers for testing
	mockProvider := mock.NewMockProvider()
	faillingMockProvider := mock.NewFailingMockProvider()
	timeoutMockProvider := mock.NewTimeoutMockProvider()

	// Define the providers
	providers := []types.Provider{
		timeoutMockProvider,
		mockProvider,
		coinbaseProvider,
		faillingMockProvider,
	}

	// Initializing the oracle
	oracle := oracle.New(
		logger,
		providerTimeout,
		oracleTicker,
		providers,
		utils.ComputeMedian(),
	)

	// Client set up and start
	client := client.NewLocalClient(oracle)
	if err := client.Start(context.Background()); err != nil {
		panic(err)
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
