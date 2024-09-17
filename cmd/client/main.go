package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	connectgrpc "github.com/skip-mev/connect/v2/pkg/grpc"
	"github.com/skip-mev/connect/v2/service/servers/oracle/types"
)

var (
	rootCmd = &cobra.Command{
		Use:   "client",
		Short: "Continuously calls the Prices RPC on the sidecar.",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			// Channel with width for termination signal.
			sigs := make(chan os.Signal, 1)

			// Gracefully trigger close on interrupt or terminate signals.
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

			// Set up a connection to the server.
			url := fmt.Sprintf("%s:%s", host, port)
			conn, err := connectgrpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()

			// Create a new client
			client := types.NewOracleClient(conn)

			// Continuous loop
			for {
				select {
				case <-sigs:
					log.Printf("Received interrupt or terminate signal, exiting...\n")
					return
				default:
					// Call Prices RPC
					log.Printf("Calling Prices RPC...\n")
					resp, err := client.Prices(context.Background(), &types.QueryPricesRequest{})
					if err != nil {
						log.Fatalf("could not get prices: %v", err)
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
		},
	}
	// host stores the host portion of the sidecar RPC to query.
	host string
	// port stores the port of the sidecar RPC to query.
	port string
)

func init() {
	rootCmd.Flags().StringVarP(
		&host,
		"host",
		"",
		"localhost",
		"host of the grpc-service, which the client will connect to",
	)
	rootCmd.Flags().StringVarP(
		&port,
		"port",
		"",
		"8080",
		"port the grpc-service, which the client will connect to",
	)
}

func main() {
	rootCmd.Execute()
}
