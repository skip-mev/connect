package prometheus_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/service/servers/prometheus"
)

// Test that Starting the server fails if the address is incorrect.
func TestStart(t *testing.T) {
	t.Run("Start fails with incorrect address", func(t *testing.T) {
		address := ":8081"

		ps, err := prometheus.NewPrometheusServer(address, nil)
		require.Nil(t, ps)
		require.Error(t, err, "invalid prometheus server address: :8080")
	})

	t.Run("Start succeeds with correct address", func(t *testing.T) {
		address := "0.0.0.0:8081"

		ps, err := prometheus.NewPrometheusServer(address, zap.NewNop())
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
