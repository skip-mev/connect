package polymarket

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/types"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var candidateWinsElectionToken = types.DefaultProviderTicker{
	OffChainTicker: "95128817762909535143571435260705470642391662537976312011260538371392879420759",
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
		noError          bool
		ids              []types.ProviderTicker
		responseBody     string
		expectedResponse types.PriceResponse
	}{
		{
			name:         "happy case",
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
			name:         "1.00 should resolve to 0.999...",
			ids:          []types.ProviderTicker{candidateWinsElectionToken},
			noError:      true,
			responseBody: `{ "price": "1.00" }`,
			expectedResponse: types.NewPriceResponse(
				types.ResolvedPrices{
					id: types.NewPriceResult(big.NewFloat(priceAdjustment), time.Now().UTC()),
				},
				nil,
			),
		},
		{
			name:         "too many IDs",
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
			ids:          []types.ProviderTicker{candidateWinsElectionToken},
			responseBody: `{"price": "0fa3adk"}"`,
			expectedResponse: types.NewPriceResponseWithErr(
				[]types.ProviderTicker{candidateWinsElectionToken},
				providertypes.NewErrorWithCode(fmt.Errorf("failed to convert %q to float", "0fa3adk"), providertypes.ErrorFailedToDecode),
			),
		},
		{
			name:         "bad price - max",
			ids:          []types.ProviderTicker{candidateWinsElectionToken},
			responseBody: `{"price": "1.0001"}"`,
			expectedResponse: types.NewPriceResponseWithErr(
				[]types.ProviderTicker{candidateWinsElectionToken},
				providertypes.NewErrorWithCode(fmt.Errorf("price exceeded 1.00"), providertypes.ErrorInvalidResponse),
			),
		},
		{
			name:         "bad price - negative",
			ids:          []types.ProviderTicker{candidateWinsElectionToken},
			responseBody: `{"price": "-0.12"}"`,
			expectedResponse: types.NewPriceResponseWithErr(
				[]types.ProviderTicker{candidateWinsElectionToken},
				providertypes.NewErrorWithCode(fmt.Errorf("price must be greater than 0.00"), providertypes.ErrorInvalidResponse),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpInput := &http.Response{
				Body: io.NopCloser(bytes.NewBufferString(tc.responseBody)),
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
