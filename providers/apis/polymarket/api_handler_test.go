package polymarket

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var candidateWinsElectionToken = types.DefaultProviderTicker{
	OffChainTicker: "95128817762909535143571435260705470642391662537976312011260538371392879420759",
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
		{
			name: "Invalid host",
			modifyConfig: func(cfg config.APIConfig) config.APIConfig {
				cfg.Endpoints[0].URL = "https://foobar.com/price"
				return cfg
			},
			expectError: true,
			errorMsg:    "invalid polymarket URL: expected",
		},
		{
			name: "Invalid endpoint path",
			modifyConfig: func(cfg config.APIConfig) config.APIConfig {
				cfg.Endpoints[0].URL = "https://" + host + "/foo"
				return cfg
			},
			expectError: true,
			errorMsg:    `invalid polymarket endpoint url path /foo`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultAPIConfig
			cfg.Endpoints = append([]config.Endpoint{}, DefaultAPIConfig.Endpoints...)
			modifiedConfig := tt.modifyConfig(cfg)
			_, err := NewAPIHandler(modifiedConfig)
			if tt.expectError {
				fmt.Println(err.Error())
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
			expectedURL: fmt.Sprintf(URL, candidateWinsElectionToken),
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
	id := candidateWinsElectionToken
	handler, err := NewAPIHandler(DefaultAPIConfig)
	require.NoError(t, err)
	testCases := []struct {
		name             string
		path             string
		noError          bool
		ids              []types.ProviderTicker
		responseBody     string
		expectedResponse types.PriceResponse
	}{
		{
			name:         "happy case from midpoint",
			path:         "/midpoint",
			ids:          []types.ProviderTicker{candidateWinsElectionToken},
			noError:      true,
			responseBody: `{ "mid": "0.45" }`,
			expectedResponse: types.NewPriceResponse(
				types.ResolvedPrices{
					id: types.NewPriceResult(big.NewFloat(0.45), time.Now().UTC()),
				},
				nil,
			),
		},
		{
			name:         "happy case from price",
			path:         "/price",
			ids:          []types.ProviderTicker{candidateWinsElectionToken},
			noError:      true,
			responseBody: `{ "price": "0.45" }`,
			expectedResponse: types.NewPriceResponse(
				types.ResolvedPrices{
					id: types.NewPriceResult(big.NewFloat(0.45), time.Now().UTC()),
				},
				nil,
			),
		},
		{
			name:         "bad path",
			path:         "/foobar",
			ids:          []types.ProviderTicker{candidateWinsElectionToken},
			responseBody: `{"mid": "234.3"}"`,
			expectedResponse: types.NewPriceResponseWithErr(
				[]types.ProviderTicker{candidateWinsElectionToken},
				providertypes.NewErrorWithCode(fmt.Errorf("unknown request path %q", "/foobar"), providertypes.ErrorFailedToDecode),
			),
		},
		{
			name:         "1.00 should resolve to 0.999...",
			path:         "/midpoint",
			ids:          []types.ProviderTicker{candidateWinsElectionToken},
			noError:      true,
			responseBody: `{ "mid": "1.00" }`,
			expectedResponse: types.NewPriceResponse(
				types.ResolvedPrices{
					id: types.NewPriceResult(big.NewFloat(priceAdjustmentMax), time.Now().UTC()),
				},
				nil,
			),
		},
		{
			name:         "0.00 should resolve to 0.00001",
			path:         "/midpoint",
			ids:          []types.ProviderTicker{candidateWinsElectionToken},
			noError:      true,
			responseBody: `{ "mid": "0.00" }`,
			expectedResponse: types.NewPriceResponse(
				types.ResolvedPrices{
					id: types.NewPriceResult(big.NewFloat(priceAdjustmentMin), time.Now().UTC()),
				},
				nil,
			),
		},
		{
			name:         "too many IDs",
			path:         "/midpoint",
			ids:          []types.ProviderTicker{candidateWinsElectionToken, candidateWinsElectionToken},
			responseBody: ``,
			expectedResponse: types.NewPriceResponseWithErr(
				[]types.ProviderTicker{candidateWinsElectionToken, candidateWinsElectionToken},
				providertypes.NewErrorWithCode(
					fmt.Errorf("expected 1 ticker, got 2"),
					providertypes.ErrorInvalidResponse,
				),
			),
		},
		{
			name:         "invalid JSON",
			path:         "/midpoint",
			ids:          []types.ProviderTicker{candidateWinsElectionToken},
			responseBody: `{"mid": "0fa3adk"}"`,
			expectedResponse: types.NewPriceResponseWithErr(
				[]types.ProviderTicker{candidateWinsElectionToken},
				providertypes.NewErrorWithCode(fmt.Errorf("failed to convert %q to float", "0fa3adk"), providertypes.ErrorFailedToDecode),
			),
		},
		{
			name:         "bad price - max",
			path:         "/midpoint",
			ids:          []types.ProviderTicker{candidateWinsElectionToken},
			responseBody: `{"mid": "1.0001"}"`,
			expectedResponse: types.NewPriceResponseWithErr(
				[]types.ProviderTicker{candidateWinsElectionToken},
				providertypes.NewErrorWithCode(fmt.Errorf("price exceeded 1.00"), providertypes.ErrorInvalidResponse),
			),
		},
		{
			name:         "bad price - negative",
			path:         "/midpoint",
			ids:          []types.ProviderTicker{candidateWinsElectionToken},
			responseBody: `{"mid": "-0.12"}"`,
			expectedResponse: types.NewPriceResponseWithErr(
				[]types.ProviderTicker{candidateWinsElectionToken},
				providertypes.NewErrorWithCode(fmt.Errorf("price must be greater than 0.00"), providertypes.ErrorInvalidResponse),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpInput := &http.Response{
				Body:    io.NopCloser(bytes.NewBufferString(tc.responseBody)),
				Request: &http.Request{URL: &url.URL{Path: tc.path}},
			}
			res := handler.ParseResponse(tc.ids, httpInput)

			// timestamps are off, repair here.
			if tc.noError {
				val := tc.expectedResponse.Resolved[tc.ids[0]]
				val.Timestamp = res.Resolved[tc.ids[0]].Timestamp
				tc.expectedResponse.Resolved[tc.ids[0]] = val
			}
			require.Equal(t, tc.expectedResponse.String(), res.String())
		})
	}
}
