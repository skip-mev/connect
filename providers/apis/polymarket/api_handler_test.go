package polymarket

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
)

var candidateWinsElectionToken = types.DefaultProviderTicker{
	OffChainTicker: "0xc6485bb7ea46d7bb89beb9c91e7572ecfc72a6273789496f78bc5e989e4d1638/95128817762909535143571435260705470642391662537976312011260538371392879420759",
}

func TestNewAPIHandler(t *testing.T) {
	tests := []struct {
		name         string
		modifyConfig func(config.APIConfig) config.APIConfig
		expectError  bool
		errorMsg     string
	}{
		{
			name: "Valid configuration",
			modifyConfig: func(cfg config.APIConfig) config.APIConfig {
				return cfg // No modifications
			},
			expectError: false,
		},
		{
			name: "Invalid name",
			modifyConfig: func(cfg config.APIConfig) config.APIConfig {
				cfg.Name = "InvalidName"
				return cfg
			},
			expectError: true,
			errorMsg:    "expected api config name polymarket_api, got InvalidName",
		},
		{
			name: "Too many endpoints",
			modifyConfig: func(cfg config.APIConfig) config.APIConfig {
				cfg.Endpoints = append(cfg.Endpoints, cfg.Endpoints...)
				return cfg
			},
			expectError: true,
			errorMsg:    "invalid polymarket endpoint config: expected 1 endpoint got 2",
		},
		{
			name: "Disabled API",
			modifyConfig: func(cfg config.APIConfig) config.APIConfig {
				cfg.Enabled = false
				return cfg
			},
			expectError: true,
			errorMsg:    "api config for polymarket_api is not enabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultAPIConfig
			cfg.Endpoints = append([]config.Endpoint{}, DefaultAPIConfig.Endpoints...)
			modifiedConfig := tt.modifyConfig(cfg)
			_, err := NewAPIHandler(modifiedConfig)
			if tt.expectError {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCreateURL(t *testing.T) {
	testCases := []struct {
		name        string
		pts         []types.ProviderTicker
		expectedURL string
		expErr      string
	}{
		{
			name:   "empty",
			pts:    []types.ProviderTicker{},
			expErr: "expected 1 ticker, got 0",
		},
		{
			name: "too many",
			pts: []types.ProviderTicker{
				candidateWinsElectionToken,
				candidateWinsElectionToken,
			},
			expErr: "expected 1 ticker, got 2",
		},
		{
			name: "happy case",
			pts: []types.ProviderTicker{
				candidateWinsElectionToken,
			},
			expectedURL: fmt.Sprintf(URL, "0xc6485bb7ea46d7bb89beb9c91e7572ecfc72a6273789496f78bc5e989e4d1638"),
		},
	}
	h, err := NewAPIHandler(DefaultAPIConfig)
	require.NoError(t, err)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, err := h.CreateURL(tc.pts)
			if tc.expErr != "" {
				require.ErrorContains(t, err, tc.expErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, url, tc.expectedURL)
			}
		})
	}
}

func TestParseResponse(t *testing.T) {
	handler, err := NewAPIHandler(DefaultAPIConfig)
	require.NoError(t, err)
	testCases := map[string]struct {
		data          string
		ticker        []types.ProviderTicker
		expectedErr   string
		expectedPrice *big.Float
	}{
		"happy path": {
			data: `{"tokens": [{
          "token_id": "95128817762909535143571435260705470642391662537976312011260538371392879420759",
          "outcome": "Yes",
          "price": 1}]}]}`,
			ticker:        []types.ProviderTicker{candidateWinsElectionToken},
			expectedPrice: big.NewFloat(1.00),
		},
		"zero resolution": {
			data: `{"tokens": [{
          "token_id": "95128817762909535143571435260705470642391662537976312011260538371392879420759",
          "outcome": "Yes",
          "price": 0}]}]}`,
			ticker:        []types.ProviderTicker{candidateWinsElectionToken},
			expectedPrice: big.NewFloat(priceAdjustmentMin),
		},
		"other values work": {
			data: `{"tokens": [{
          "token_id": "95128817762909535143571435260705470642391662537976312011260538371392879420759",
          "outcome": "Yes",
          "price": 0.325}]}]}`,
			ticker:        []types.ProviderTicker{candidateWinsElectionToken},
			expectedPrice: big.NewFloat(0.325),
		},
		"token not in response": {
			data: `{"tokens": [{
          "token_id": "35128817762909535143571435260705470642391662537976312011260538371392879420759",
          "outcome": "Yes",
          "price": 0.325}]}]}`,
			ticker:      []types.ProviderTicker{candidateWinsElectionToken},
			expectedErr: "token ID 95128817762909535143571435260705470642391662537976312011260538371392879420759 not found in response",
		},
		"bad response data": {
			data: `{"tokens": [{
          "token_id":z,
          "outcome": "Yes",
          "price": 0.325}]}]}`,
			ticker:      []types.ProviderTicker{candidateWinsElectionToken},
			expectedErr: "failed to decode market response",
		},
		"too many tickers": {
			data:        `{"tokens": []}`,
			ticker:      []types.ProviderTicker{candidateWinsElectionToken, candidateWinsElectionToken},
			expectedErr: "expected 1 ticker, got 2",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			httpInput := &http.Response{
				Body: io.NopCloser(bytes.NewBufferString(tc.data)),
			}
			res := handler.ParseResponse(tc.ticker, httpInput)
			if tc.expectedErr != "" {
				require.Contains(t, res.UnResolved[tc.ticker[0]].Error(), tc.expectedErr)
			} else {
				gotPrice := res.Resolved[tc.ticker[0]].Value
				require.Equal(t, gotPrice.Cmp(tc.expectedPrice), 0, "expected %v, got %v", tc.expectedPrice, gotPrice)
				require.Equal(t, len(res.Resolved), len(tc.ticker))
			}
		})
	}
}
