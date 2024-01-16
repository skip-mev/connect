package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/skip-mev/slinky/service"
)

var (
	host = flag.String("host", "localhost", "host for the grpc-service to listen on")
	port = flag.String("port", "8080", "port for the grpc-service to listen on")
)

func main() {
	// Channel with width for termination signal.
	sigs := make(chan os.Signal, 1)

	// Gracefully trigger close on interrupt or terminate signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Parse flags.
	flag.Parse()

	// Set up a connection to the server.
	url := fmt.Sprintf("%s:%s", *host, *port)
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new client
	client := service.NewOracleClient(conn)

	// Continuous loop
	for {
		select {
		case <-sigs:
			log.Printf("Received interrupt or terminate signal, exiting...\n")
			return
		default:
			// Call Prices RPC
			log.Printf("Calling Prices RPC...\n")
			resp, err := client.Prices(context.Background(), &service.QueryPricesRequest{})
			if err != nil {
				log.Fatalf("could not get prices: %v", err) //nolint
			}

			prices := resp.GetPrices()

			var keys []string
			for k := range prices {
				keys = append(keys, k)
			}

			// Sort the prices by the currency pair
			sort.Slice(keys, func(i, j int) bool {
				return keys[i] < keys[j]
			})

			// Log the response
			for _, key := range keys {
				log.Printf("Currency Pair, Price: (%s, %s)", key, prices[key])
			}

			// Wait for a bit before making the next request
			log.Printf("Sleeping for 10 seconds...\n\n")
			time.Sleep(time.Second * 10) // Adjust the interval as needed
		}
	}
}
