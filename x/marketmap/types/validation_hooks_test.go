package types_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

func TestDefaultDeleteMarketValidationHooks_ValidateMarket(t *testing.T) {
	tests := []struct {
		name    string
		market  types.Market
		wantErr bool
	}{
		{
			name: "valid - disabled market",
			market: types.Market{
				Ticker: types.Ticker{
					CurrencyPair: slinkytypes.CurrencyPair{
						Base:  "BTC",
						Quote: "USD",
					},
					Decimals:         3,
					MinProviderCount: 3,
					Enabled:          false,
					Metadata_JSON:    "",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid - enabled market",
			market: types.Market{
				Ticker: types.Ticker{
					CurrencyPair: slinkytypes.CurrencyPair{
						Base:  "BTC",
						Quote: "USD",
					},
					Decimals:         3,
					MinProviderCount: 3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hooks := types.DefaultDeleteMarketValidationHooks()

			err := hooks.ValidateMarket(context.Background(), tt.market)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
