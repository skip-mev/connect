package types

import (
	"errors"
	"time"

	servicemetrics "github.com/skip-mev/connect/v2/service/metrics"
)

// RecordLatencyAndStatus is used by the ABCI handlers to record their e2e latency, and the status of the request
// to their corresponding metrics objects.
func RecordLatencyAndStatus(
	metrics servicemetrics.Metrics, latency time.Duration, err error, method servicemetrics.ABCIMethod,
) {
	// observe latency
	metrics.ObserveABCIMethodLatency(method, latency)

	// increment the number of extend vote requests
	var label servicemetrics.Labeller
	if err != nil {
		_ = errors.As(err, &label)
	} else {
		label = servicemetrics.Success{}
	}
	metrics.AddABCIRequest(method, label)
}
