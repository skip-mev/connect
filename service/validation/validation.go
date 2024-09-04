package validation

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	oraclemetrics "github.com/skip-mev/connect/v2/oracle/metrics"
)

const (
	DefaultNumChecks             = 1000
	DefaultRequiredPriceLiveness = 99.0
	DefaultValidationPeriod      = 10 * time.Minute
	DefaultBurnInPeriod          = 60 * time.Second
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
	BurnInPeriod          time.Duration
	ValidationPeriod      time.Duration
	NumChecks             int
	RequiredPriceLiveness float64
}

// DefaultConfig returns a default validation config.
func DefaultConfig() Config {
	return Config{
		BurnInPeriod:          DefaultBurnInPeriod,
		NumChecks:             DefaultNumChecks,
		RequiredPriceLiveness: DefaultRequiredPriceLiveness,
		ValidationPeriod:      DefaultValidationPeriod,
	}
}

type LivenessResults map[string]float64

// Run runs the validation service, checking for missing prices and accumulating liveness data.
func (v *Validator) Run() (LivenessResults, error) {
	v.logger.Info("running in validation mode", zap.Duration("validation period (s)",
		v.cfg.ValidationPeriod))
	v.logger.Info("waiting for burn in period to end", zap.Duration("period (s)", v.cfg.BurnInPeriod))
	time.Sleep(v.cfg.BurnInPeriod)
	missingPricesMap := make(map[string]int)
	ticker := time.NewTicker(v.cfg.ValidationPeriod / time.Duration(v.cfg.NumChecks))
	v.logger.Info("beginning validation")

	go func() {
		for range ticker.C {
			v.logger.Info("checking for missing prices")

			missingPrices := v.metrics.GetMissingPrices()
			if len(missingPrices) > 0 {
				v.logger.Warn("currently missing prices", zap.Any("prices", missingPrices))
				for _, price := range missingPrices {
					missingPricesMap[price]++
				}
			}
		}
	}()

	time.Sleep(v.cfg.ValidationPeriod)
	ticker.Stop()

	resultsMap := make(LivenessResults)
	invalidTickers := make([]string, 0)
	for pairID, ticksMissing := range missingPricesMap {
		livenessRate := float64(v.cfg.NumChecks-ticksMissing) / float64(v.cfg.NumChecks)
		resultsMap[pairID] = livenessRate

		if livenessRate < v.cfg.RequiredPriceLiveness {
			v.logger.Error("liveness for price is below required rate", zap.String("pairID", pairID),
				zap.Float64("liveness_rate", livenessRate), zap.Float64("required", v.cfg.RequiredPriceLiveness))
			invalidTickers = append(invalidTickers, pairID)
		}
	}

	if len(invalidTickers) > 0 {
		return resultsMap, fmt.Errorf("invalid pairs below liveness threshold: %v", invalidTickers)
	}

	v.logger.Info("finished validation", zap.Any("liveness stats", resultsMap))

	return resultsMap, nil
}
