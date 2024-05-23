package http

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	// DefaultMultiplicativeIncrease is the default multiplicative increase factor for the exponential backoff.
	DefaultMultiplicativeIncrease float64 = 1.5
	// DefaultMultiplicativeDecrease is the default multiplicative decrease factor for the exponential backoff.
	DefaultMultiplicativeDecrease float64 = 0.9
)

// ExponentialBackOffTicker is a ticker that ticks exponentially increasing intervals. Every operation
// can be marked as successful or unsuccessful. If the operation is successful, the interval is decreased.
// Otherwise, the interval is increased. The increase/decrease is exponential. The interval is bounded by
// a minimum and maximum value. Jitter can be added to the interval to add randomness. This implementation
// is thread-safe.
type ExponentialBackOffTicker struct {
	mut    sync.Mutex
	logger *zap.Logger

	// interval is the current interval between ticks.
	interval time.Duration
	// base is the base interval between ticks.
	base time.Duration
	// multiplicativeIncrease is the multiplicative increase factor for the exponential backoff.
	multiplicativeIncrease float64
	// multiplicativeDecrease is the multiplicative decrease factor for the exponential backoff.
	multiplicativeDecrease float64
	// max is the maximum interval between ticks.
	max time.Duration
	// jitter is the jitter used to add randomness to the interval between ticks.
	jitter time.Duration
}

// NewExponentialBackOffTicker creates a new ExponentialBackOffTicker.
func NewExponentialBackOffTicker(
	logger *zap.Logger,
	base time.Duration,
	max time.Duration,
	jitter time.Duration,
	multiplicativeIncrease float64,
	multiplicativeDecrease float64,
) (*ExponentialBackOffTicker, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	if base <= 0 {
		return nil, fmt.Errorf("base interval must be positive; got %v", base)
	}

	if multiplicativeIncrease <= 1 {
		return nil, fmt.Errorf("multiplicative increase must be greater than 1; got %v", multiplicativeIncrease)
	}

	if multiplicativeDecrease >= 1 || multiplicativeDecrease <= 0 {
		return nil, fmt.Errorf("multiplicative decrease must be less than 1; got %v", multiplicativeDecrease)
	}

	if max <= 0 || max < base {
		return nil, fmt.Errorf("max interval must be positive and greater than the base interval; got max: %v, base: %v", max, base)
	}

	if jitter < 0 {
		return nil, fmt.Errorf("jitter must be non-negative; got %v", jitter)
	}

	return &ExponentialBackOffTicker{
		logger:                 logger.With(zap.String("component", "exponential_backoff_ticker")),
		interval:               max,
		base:                   base,
		max:                    max,
		jitter:                 jitter,
		multiplicativeIncrease: multiplicativeIncrease,
		multiplicativeDecrease: multiplicativeDecrease,
	}, nil
}

// Reset resets the ticker to the base interval.
func (t *ExponentialBackOffTicker) Reset() time.Duration {
	t.mut.Lock()
	defer t.mut.Unlock()

	t.logger.Info("resetting interval", zap.Duration("interval", t.base))
	t.interval = t.base
	return t.base
}

// Throttle increases/decreases the interval between ticks. If success is true, the interval is
// decreased. Otherwise, the interval is increased. The increase/decrease is exponential.
func (t *ExponentialBackOffTicker) Throttle(rateLimitSeen bool) time.Duration {
	t.mut.Lock()
	defer t.mut.Unlock()

	// If the operation was successful, decrease the interval. Otherwise, increase the interval.
	var multiple float64
	if rateLimitSeen {
		multiple = t.multiplicativeIncrease
	} else {
		multiple = t.multiplicativeDecrease
	}

	// Add jitter to the interval.
	updatedInterval := time.Duration(float64(t.interval) * multiple)
	if t.jitter > 0 {
		updatedInterval += time.Duration(rand.Int63n(int64(t.jitter)))
	}

	// Ensure the interval is within the bounds.
	if updatedInterval > t.max {
		updatedInterval = t.max
	} else if updatedInterval < t.base {
		updatedInterval = t.base
	}

	t.logger.Info(
		"updating exponential backoff interval",
		zap.Duration("current_interval", t.interval),
		zap.Duration("updated_interval", updatedInterval),
		zap.Bool("seen_rate_limit", rateLimitSeen),
	)
	t.interval = updatedInterval

	return t.interval
}

// Tick returns a channel that will receive a tick at the calculated interval.
func (t *ExponentialBackOffTicker) Tick() <-chan time.Time {
	t.mut.Lock()
	defer t.mut.Unlock()

	return time.After(t.interval)
}

// Interval returns the current interval between ticks.
func (t *ExponentialBackOffTicker) Interval() time.Duration {
	t.mut.Lock()
	defer t.mut.Unlock()

	return t.interval
}
