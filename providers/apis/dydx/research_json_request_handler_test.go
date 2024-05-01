package dydx_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/skip-mev/slinky/providers/apis/dydx"
	dydxtypes "github.com/skip-mev/slinky/providers/apis/dydx/types"
	apihandlermocks "github.com/skip-mev/slinky/providers/base/api/handlers/mocks"
	"github.com/stretchr/testify/require"
)

func TestResearchJSONRequestHandler_Do(t *testing.T) {
	mockReqHandler := apihandlermocks.NewRequestHandler(t)
	rh := dydx.NewResearchJSONRequestHandler(mockReqHandler)

	testURL := "https://raw.githubusercontent.com/dydxprotocol/v4-web/main/public/configs/otherMarketData.json"
	t.Run("check errors in http request", func(t *testing.T) {
		// mock
		ctx := context.Background()
		mockReqHandler.On("Do", ctx, testURL).Return(nil, fmt.Errorf("error making request")).Once()

		// test
		resp, err := rh.Do(ctx, testURL)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("check that non-OK status is propagated w/ no body change", func(t *testing.T) {
		// mock
		ctx := context.Background()
		mockReqHandler.On("Do", ctx, testURL).Return(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}, nil).Once()

		// test
		resp, err := rh.Do(ctx, testURL)
		require.Error(t, err)
		require.NotNil(t, resp)
		require.Equal(t, 500, resp.StatusCode)
	})

	t.Run("check that errors unmarshalling body return no response", func(t *testing.T) {
		// mock
		ctx := context.Background()
		mockReqHandler.On("Do", ctx, testURL).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}, nil).Once()

		// test
		resp, err := rh.Do(ctx, testURL)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("check that errors converting research json to QueryAllMarketsParamsResponse return no response", func(t *testing.T) {
		// mock
		ctx := context.Background()
		mockReqHandler.On("Do", ctx, testURL).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(bytes.NewBufferString(`{
				"1INCH": {
				}
			}`)),
		}, nil)

		// test
		resp, err := rh.Do(ctx, testURL)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("check that successful response is returned with body marshalled to QueryAllMarketsParamsResponse", func(t *testing.T) {
		// mock
		ctx := context.Background()
		researchJSON := dydxtypes.ResearchJSON{
			"1INCH": {
				"params": {
					ID:                0,
					Pair:              "1INCH-USD",
					Exponent:          -10.0,
					MinPriceChangePpm: 4000,
					MinExchanges:      3,
					ExchangeConfigJSON: []dydxtypes.ExchangeMarketConfigJson{
						{
							ExchangeName:   "Binance",
							Ticker:         "1INCHUSDT",
							AdjustByMarket: "USDT-USD",
						},
						{
							ExchangeName: "CoinbasePro",
							Ticker:       "1INCH-USD",
						},
						{
							ExchangeName:   "Gate",
							Ticker:         "1INCH_USDT",
							AdjustByMarket: "USDT-USD",
						},
						{
							ExchangeName:   "Kucoin",
							Ticker:         "1INCH-USDT",
							AdjustByMarket: "USDT-USD",
						},
						{
							ExchangeName:   "Mexc",
							Ticker:         "1INCH_USDT",
							AdjustByMarket: "USDT-USD",
						},
						{
							ExchangeName:   "Okx",
							Ticker:         "1INCH-USDT",
							AdjustByMarket: "USDT-USD",
						},
					},
				},
			},
		}
		bz, err := json.Marshal(researchJSON)
		require.NoError(t, err)

		mockReqHandler.On("Do", ctx, testURL).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBuffer(bz)),
		}, nil)

		// test
		resp, err := rh.Do(ctx, testURL)
		require.NoError(t, err)
		require.NotNil(t, resp)
		defer resp.Body.Close()
		// read the response body as a QueryAllMarketsParamsResponse
		var qamp dydxtypes.QueryAllMarketParamsResponse
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&qamp))

		// check the response
		require.Equal(t, dydxtypes.QueryAllMarketParamsResponse{
			MarketParams: []dydxtypes.MarketParam{
				{
					Id:                 0,
					Pair:               "1INCH-USD",
					Exponent:           int32(-10),
					MinExchanges:       3,
					MinPriceChangePpm:  4000,
					ExchangeConfigJson: "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"1INCHUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"1INCH-USD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"1INCH_USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"1INCH-USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Mexc\",\"ticker\":\"1INCH_USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"1INCH-USDT\",\"adjustByMarket\":\"USDT-USD\"}]}",
				},
			},
		}, qamp)
	})
}
