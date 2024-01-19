package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
)

func TestProviderMarketConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      config.MarketConfig
		expectedErr bool
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.ValidateBasic()
			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
