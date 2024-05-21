package handlers

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	// DefaultMultipler is the default multiplier used to increase the interval between ticks.
	DefaultMultipler = 1.5
)

// ExponentialBackOffTicker is a ticker that ticks exponentially increasing intervals. The
// algorithm is composed of:
//
//   - A base interval that is the minimum interval between ticks.
//   - A multiplier that is used to increase the interval between ticks.
//   - A maximum interval that is the maximum interval between ticks.
//   - A jitter that is used to add randomness to the interval between ticks.
//
// The interval between ticks is calculated as:
//   - interval = min(base * multiplier^attempts, max)
//   - interval = interval + random(-jitter, jitter)
//
// The ticker will tick at the calculated interval. Note that this implementation is thread-safe.
type ExponentialBackOffTicker struct {
	mut    sync.Mutex
	logger *zap.Logger

	// interval is the current interval between ticks.
	interval time.Duration
	// base is the base interval between ticks.
	base time.Duration
	// multiplier is the multiplier used to increase the interval between ticks.
	multiplier float64
	// max is the maximum interval between ticks.
	max time.Duration
	// jitter is the jitter used to add randomness to the interval between ticks.
	jitter time.Duration
}

// NewExponentialBackOffTicker creates a new ExponentialBackOffTicker.
func NewExponentialBackOffTicker(
	logger *zap.Logger,
	base time.Duration,
	multipler float64,
	max time.Duration,
	jitter time.Duration,
) (*ExponentialBackOffTicker, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	if base <= 0 {
		return nil, fmt.Errorf("base interval must be positive")
	}

	if multipler <= 1 {
		return nil, fmt.Errorf("multiplier must be greater than 1")
	}

	if max <= 0 {
		return nil, fmt.Errorf("max interval must be positive")
	}

	if jitter < 0 {
		return nil, fmt.Errorf("jitter must be non-negative")
	}

	return &ExponentialBackOffTicker{
		logger:     logger.With(zap.String("component", "exponential_backoff_ticker")),
		interval:   base,
		base:       base,
		multiplier: multipler,
		max:        max,
		jitter:     jitter,
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

// BackOff increases/decreases the interval between ticks. If success is true, the interval is
// decreased. Otherwise, the interval is increased. The increase/decrease is exponential.
func (t *ExponentialBackOffTicker) BackOff(success bool) time.Duration {
	t.mut.Lock()
	defer t.mut.Unlock()

	var updatedInterval time.Duration
	// If the operation was successful, decrease the interval. Otherwise, increase the interval.
	if !success {
		updatedInterval = t.interval * time.Duration(t.multiplier)
	} else {
		updatedInterval = t.base / time.Duration(t.multiplier)
	}

	// Add jitter to the interval.
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
