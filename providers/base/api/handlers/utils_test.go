package handlers_test

import (
	"testing"
	"time"

	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base/api/handlers"
	"github.com/test-go/testify/require"
	"go.uber.org/zap"
)

func TestNewExponentialBackOffTicker(t *testing.T) {
	tc := []struct {
		name       string
		logger     *zap.Logger
		base       time.Duration
		multiplier float64
		max        time.Duration
		jitter     time.Duration
		err        bool
	}{
		{
			name:       "valid",
			logger:     logger,
			base:       1 * time.Second,
			multiplier: 2.0,
			max:        10 * time.Second,
			jitter:     1 * time.Second,
			err:        false,
		},
		{
			name:       "invalid logger",
			logger:     nil,
			base:       1 * time.Second,
			multiplier: 2.0,
			max:        10 * time.Second,
			jitter:     1 * time.Second,
			err:        true,
		},
		{
			name:       "invalid base",
			logger:     logger,
			base:       0,
			multiplier: 2.0,
			max:        10 * time.Second,
			jitter:     1 * time.Second,
			err:        true,
		},
		{
			name:       "invalid multiplier",
			logger:     logger,
			base:       1 * time.Second,
			multiplier: 0,
			max:        10 * time.Second,
			jitter:     1 * time.Second,
			err:        true,
		},
		{
			name:       "invalid max",
			logger:     logger,
			base:       1 * time.Second,
			multiplier: 2.0,
			max:        0,
			jitter:     1 * time.Second,
			err:        true,
		},
		{
			name:       "invalid jitter",
			logger:     logger,
			base:       1 * time.Second,
			multiplier: 2.0,
			max:        10 * time.Second,
			jitter:     -1 * time.Second,
			err:        true,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handlers.NewExponentialBackOffTicker(
				tt.logger,
				tt.base,
				tt.multiplier,
				tt.max,
				tt.jitter,
			)
			if tt.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTick(t *testing.T) {
	baseInterval := 1 * time.Second
	multiplier := 2.0
	maxInterval := 10 * time.Second
	jitter := 1 * time.Second

	// Set timestamps on the logger so that we can compare logs.
	logger, _ := zap.NewDevelopment()
	logger = logger.WithOptions(zap.IncreaseLevel(zap.DebugLevel))

	t.Run("tick with no resets or backoff", func(t *testing.T) {
		expTicker, err := handlers.NewExponentialBackOffTicker(
			logger,
			baseInterval,
			multiplier,
			maxInterval,
			jitter,
		)
		require.NoError(t, err)

		for i := 0; i < 5; i++ {
			duration := expTicker.Interval()
			require.Equal(t, baseInterval, duration)

			require.WithinDuration(t, time.Now().Add(baseInterval), <-expTicker.Tick(), 10*time.Millisecond)
		}
	})

	t.Run("tick with one backoff, no randomness and no resets", func(t *testing.T) {
		expTicker, err := handlers.NewExponentialBackOffTicker(
			logger,
			baseInterval,
			multiplier,
			maxInterval,
			0,
		)
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			duration := expTicker.Interval()
			require.Equal(t, baseInterval, duration)

			require.WithinDuration(t, time.Now().Add(baseInterval), <-expTicker.Tick(), 10*time.Millisecond)
		}

		// After the first backoff, the interval should be doubled.
		duration := expTicker.BackOff(false)
		require.Equal(t, baseInterval*2, duration)
		require.WithinDuration(t, time.Now().Add(baseInterval*2), <-expTicker.Tick(), 10*time.Millisecond)
	})

	t.Run("tick with max backoffs, no randomness and no resets", func(t *testing.T) {
		expTicker, err := handlers.NewExponentialBackOffTicker(
			logger,
			baseInterval,
			multiplier,
			maxInterval,
			0,
		)
		require.NoError(t, err)

		for i := 0; i < 10; i++ {
			duration := expTicker.Interval()

			backOffMultipler := math.Min(5, i)
			expectedDuration := baseInterval * time.Duration(backOffMultipler)
			require.Equal(t, expectedDuration, duration)
			require.True(t, duration <= maxInterval)

			require.WithinDuration(t, time.Now().Add(expectedDuration), <-expTicker.Tick(), 10*time.Millisecond)
		}
	})
}

func TestBackOff(t *testing.T) {
	tc := []struct {
		name             string
		base             time.Duration
		multiplier       float64
		max              time.Duration
		jitter           time.Duration
		statuses         []bool
		expectedInterval time.Duration
	}{
		{
			name:             "no backoff",
			base:             1 * time.Second,
			multiplier:       2.0,
			max:              10 * time.Second,
			jitter:           0,
			statuses:         []bool{},
			expectedInterval: 1 * time.Second,
		},
		{
			name:             "backoff once (failure) with no jitter",
			base:             1 * time.Second,
			multiplier:       2.0,
			max:              10 * time.Second,
			jitter:           0,
			statuses:         []bool{false},
			expectedInterval: 2 * time.Second,
		},
		{
			name:             "backoff 2 times (failure) with no jitter",
			base:             1 * time.Second,
			multiplier:       2.0,
			max:              10 * time.Second,
			jitter:           0,
			statuses:         []bool{false, false},
			expectedInterval: 4 * time.Second,
		},
		{
			name:             "backoff 3 times (failure) with no jitter",
			base:             1 * time.Second,
			multiplier:       2.0,
			max:              10 * time.Second,
			jitter:           0,
			statuses:         []bool{false, false, false},
			expectedInterval: 8 * time.Second,
		},
		{
			name:             "backoff 4 times (failure) with no jitter, reaches the max interval",
			base:             1 * time.Second,
			multiplier:       2.0,
			max:              10 * time.Second,
			jitter:           0,
			statuses:         []bool{false, false, false, false},
			expectedInterval: 10 * time.Second,
		},
		{
			name:             "backoff 1 time (success) with no jitter, reaches the base interval",
			base:             1 * time.Second,
			multiplier:       2.0,
			max:              10 * time.Second,
			jitter:           0,
			statuses:         []bool{true},
			expectedInterval: 1 * time.Second,
		},
		{
			name:             "backoff 2 times (failure followed by success) with no jitter",
			base:             1 * time.Second,
			multiplier:       2.0,
			max:              10 * time.Second,
			jitter:           0,
			statuses:         []bool{false, true},
			expectedInterval: 1 * time.Second,
		},
		{
			name:             "backoff 1 time (failure) with jitter",
			base:             1 * time.Second,
			multiplier:       2.0,
			max:              10 * time.Second,
			jitter:           1 * time.Second,
			statuses:         []bool{false},
			expectedInterval: 2 * time.Second,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			expTicker, err := handlers.NewExponentialBackOffTicker(
				logger,
				tt.base,
				tt.multiplier,
				tt.max,
				tt.jitter,
			)
			require.NoError(t, err)

			for _, status := range tt.statuses {
				expTicker.BackOff(status)
			}

			if tt.jitter == 0 {
				require.Equal(t, tt.expectedInterval, expTicker.Interval())
				return
			}

			// Jitter is non-zero, so we can't predict the exact interval.
			interval := expTicker.Interval()
			lowerBound := tt.expectedInterval - tt.jitter
			upperBound := tt.expectedInterval + tt.jitter
			require.True(t, interval >= lowerBound && interval <= upperBound)
		})
	}
}
