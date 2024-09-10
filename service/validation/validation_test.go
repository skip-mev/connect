package validation_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/service/validation"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  validation.Config
		wantErr bool
	}{
		{
			name:    "empty invalid",
			wantErr: true,
		},
		{
			name:    "valid default",
			config:  validation.DefaultConfig(),
			wantErr: false,
		},
		{
			name: "invalid percent",
			config: validation.Config{
				BurnInPeriod:                 validation.DefaultValidationPeriod,
				ValidationPeriod:             validation.DefaultValidationPeriod,
				NumChecks:                    validation.DefaultNumChecks,
				RequiredPriceLivenessPercent: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid percent",
			config: validation.Config{
				BurnInPeriod:                 validation.DefaultValidationPeriod,
				ValidationPeriod:             validation.DefaultValidationPeriod,
				NumChecks:                    0,
				RequiredPriceLivenessPercent: validation.DefaultRequiredPriceLivenessPercent,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
