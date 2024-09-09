package validation

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	oraclemetrics "github.com/skip-mev/connect/v2/oracle/metrics"
)

const (
	DefaultNumChecks                    = 1000
	DefaultRequiredPriceLivenessPercent = 99.0
	DefaultValidationPeriod             = 10 * time.Minute
	DefaultBurnInPeriod                 = 60 * time.Second
)

type Validator struct {
	logger  *zap.Logger
	metrics oraclemetrics.Metrics
	cfg     Config
}

// NewDefaultValidator returns a new validator service with the default configuration.
func NewDefaultValidator(logger *zap.Logger, metrics oraclemetrics.Metrics) Validator {
	return NewValidator(logger, metrics, DefaultConfig())
}

// NewValidator returns a new validator service with given configuration.
func NewValidator(logger *zap.Logger, metrics oraclemetrics.Metrics, cfg Config) Validator {
	return Validator{
		logger:  logger.Named("validation"),
		metrics: metrics,
		cfg:     cfg,
	}
}

// Config includes information for configuring a validation service instance.
type Config struct {
	// BurnInPeriod is the amount of time to let the sidecar run before checking validation.
	// This ensures that any markets or prices that take a few iterations to land will be populated,
	// eliminating false positives.
	BurnInPeriod time.Duration
	// ValidationPeriod is the amount of time the validation service will check that prices are landing correctly.
	ValidationPeriod time.Duration
	// NumChecks is the number of times the validation service will check validity over the ValidationPeriod.
	NumChecks int
	// RequiredPriceLivenessPercent is the percentage of liveness each price must demonstrate to be considered "valid".
	RequiredPriceLivenessPercent float64
}

// DefaultConfig returns a default validation config.
func DefaultConfig() Config {
	return Config{
		BurnInPeriod:                 DefaultBurnInPeriod,
		NumChecks:                    DefaultNumChecks,
		RequiredPriceLivenessPercent: DefaultRequiredPriceLivenessPercent,
		ValidationPeriod:             DefaultValidationPeriod,
	}
}

// Validate checks the validity of fields in the Config.
func (c *Config) Validate() error {
	if c.NumChecks <= 0 {
		return fmt.Errorf("num checks must be greater than zero")
	}

	if c.RequiredPriceLivenessPercent <= 0 {
		return fmt.Errorf("required price liveness percent must be greater than zero")
	}

	return nil
}

type LivenessResults map[string]float64

// Run runs the validation service, checking for missing prices and accumulating liveness data.
func (v *Validator) Run(ctx context.Context) (LivenessResults, error) {
	if err := v.cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}

	v.logger.Info("running in validation mode", zap.Duration("validation period (s)",
		v.cfg.ValidationPeriod))

	v.logger.Info("waiting for burn in period to end", zap.Duration("period (s)", v.cfg.BurnInPeriod))
	time.Sleep(v.cfg.BurnInPeriod)

	missingPricesMap := make(map[string]int)
	tickTime := v.cfg.ValidationPeriod / time.Duration(v.cfg.NumChecks)
	ticker := time.NewTicker(tickTime)

	v.logger.Info("beginning validation")

	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				v.logger.Info("checking for missing prices")

				missingPrices := v.metrics.GetMissingPrices()
				if len(missingPrices) > 0 {
					v.logger.Warn("currently missing prices", zap.Any("prices", missingPrices))
					for _, price := range missingPrices {
						missingPricesMap[price]++
					}
				}
			case <-done:
				return
			case <-ctx.Done():
				close(done)
				return
			}
		}
	}()

	// check for early exits
	numSleeps := 0
	for numSleeps < v.cfg.NumChecks {
		select {
		case <-done:
			v.logger.Info("context canceled - exiting early")
			return nil, nil
		default:
			time.Sleep(tickTime)
			numSleeps++
		}
	}

	ticker.Stop()
	close(done)

	resultsMap := make(LivenessResults)
	invalidTickers := make([]string, 0)
	for pairID, ticksMissing := range missingPricesMap {
		livenessRate := float64(v.cfg.NumChecks-ticksMissing) / float64(v.cfg.NumChecks)
		resultsMap[pairID] = livenessRate

		if livenessRate < v.cfg.RequiredPriceLivenessPercent {
			v.logger.Error("liveness for price is below required rate", zap.String("pairID", pairID),
				zap.Float64("liveness_rate", livenessRate), zap.Float64("required", v.cfg.RequiredPriceLivenessPercent))
			invalidTickers = append(invalidTickers, pairID)
		}
	}

	if len(invalidTickers) > 0 {
		return resultsMap, fmt.Errorf("invalid pairs below liveness threshold: %v", invalidTickers)
	}

	v.logger.Info("finished validation", zap.Any("liveness stats", resultsMap))

	return resultsMap, nil
}
