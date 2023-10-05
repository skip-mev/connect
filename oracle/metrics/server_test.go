package metrics_test

import (
	"net/http"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/skip-mev/slinky/oracle/metrics"
	"github.com/stretchr/testify/require"
)

// Test that Starting the server fails if the address is incorrect.
func TestStart(t *testing.T) {
	t.Run("Start fails with incorrect address", func(t *testing.T) {
		address := ":8080"

		ps, err := metrics.NewPrometheusServer(address, nil)
		require.Nil(t, ps)
		require.Error(t, err, "invalid prometheus server address: :8080")
	})

	t.Run("Start succeeds with correct address", func(t *testing.T) {
		address := "0.0.0.0:8080"

		ps, err := metrics.NewPrometheusServer(address, log.NewTestLogger(t))
		require.NotNil(t, ps)
		require.NoError(t, err)

		// start the server
		go ps.Start()

		time.Sleep(1 * time.Second)

		// ping the server
		require.True(t, pingServer("http://"+address))

		// close the server
		ps.Close()

		// expect the server to be closed within 3 seconds
		select {
		case <-ps.Done():
		case <-time.After(3 * time.Second):
		}
	})
}

func pingServer(address string) bool {
	timeout := 5 * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(address)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
