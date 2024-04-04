package uniswap_test

import (
	"testing"

	"github.com/skip-mev/slinky/providers/apis/uniswap"
	"github.com/stretchr/testify/require"
)

func TestPoolConfig(t *testing.T) {
	t.Run("empty config", func(t *testing.T) {
		cfg := uniswap.PoolConfig{}
		require.Error(t, cfg.ValidateBasic())
	})

	t.Run("invalid address", func(t *testing.T) {
		cfg := uniswap.PoolConfig{
			Address: "invalid",
		}
		require.Error(t, cfg.ValidateBasic())
	})

	t.Run("invalid base decimals", func(t *testing.T) {
		cfg := uniswap.PoolConfig{
			Address:      "0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8",
			BaseDecimals: -1,
		}
		require.Error(t, cfg.ValidateBasic())
	})

	t.Run("invalid quote decimals", func(t *testing.T) {
		cfg := uniswap.PoolConfig{
			Address:       "0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8",
			BaseDecimals:  18,
			QuoteDecimals: -1,
		}
		require.Error(t, cfg.ValidateBasic())
	})

	t.Run("valid config", func(t *testing.T) {
		cfg := uniswap.PoolConfig{
			Address:       "0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8",
			BaseDecimals:  18,
			QuoteDecimals: 18,
		}
		require.NoError(t, cfg.ValidateBasic())
	})
}
