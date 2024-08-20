package uniswapv3_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/constants"
	"github.com/skip-mev/connect/v2/providers/apis/defi/uniswapv3"
)

func TestPoolConfig(t *testing.T) {
	t.Run("empty config", func(t *testing.T) {
		cfg := uniswapv3.PoolConfig{}
		require.Error(t, cfg.ValidateBasic())
	})

	t.Run("invalid address", func(t *testing.T) {
		cfg := uniswapv3.PoolConfig{
			Address: "invalid",
		}
		require.Error(t, cfg.ValidateBasic())
	})

	t.Run("invalid base decimals", func(t *testing.T) {
		cfg := uniswapv3.PoolConfig{
			Address:      "0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8",
			BaseDecimals: -1,
		}
		require.Error(t, cfg.ValidateBasic())
	})

	t.Run("invalid quote decimals", func(t *testing.T) {
		cfg := uniswapv3.PoolConfig{
			Address:       "0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8",
			BaseDecimals:  18,
			QuoteDecimals: -1,
		}
		require.Error(t, cfg.ValidateBasic())
	})

	t.Run("valid config", func(t *testing.T) {
		cfg := uniswapv3.PoolConfig{
			Address:       "0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8",
			BaseDecimals:  18,
			QuoteDecimals: 18,
		}
		require.NoError(t, cfg.ValidateBasic())
	})
}

func TestIsValidProviderName(t *testing.T) {
	type testcase struct {
		testName     string
		providerName string
		valid        bool
	}
	testcases := []testcase{
		{
			testName:     "valid base, invalid chain",
			providerName: fmt.Sprintf("%s%s%s", uniswapv3.BaseName, uniswapv3.NameSeparator, "arbitrum"),
			valid:        false,
		},
		{
			testName:     "valid base, invalid separator",
			providerName: fmt.Sprintf("%s%s%s", uniswapv3.BaseName, "*", constants.ETHEREUM),
			valid:        false,
		},
		{
			testName:     "invalid base",
			providerName: fmt.Sprintf("%s%s%s", "uniswapv2", uniswapv3.NameSeparator, constants.ETHEREUM),
			valid:        false,
		},
		{
			testName:     "valid provider eth",
			providerName: fmt.Sprintf("%s%s%s", uniswapv3.BaseName, uniswapv3.NameSeparator, constants.ETHEREUM),
			valid:        true,
		},
		{
			testName:     "valid provider base",
			providerName: fmt.Sprintf("%s%s%s", uniswapv3.BaseName, uniswapv3.NameSeparator, constants.BASE),
			valid:        true,
		},
	}
	// Also test that all ProviderNames are Valid
	for _, providerName := range uniswapv3.ProviderNames {
		testcases = append(testcases, testcase{
			testName:     fmt.Sprintf("valid-provider-name-%s", providerName),
			providerName: providerName,
			valid:        true,
		})
	}
	for _, tc := range testcases {
		t.Run(tc.testName, func(t *testing.T) {
			require.Equal(t, tc.valid, uniswapv3.IsValidProviderName(tc.providerName))
		})
	}
}
