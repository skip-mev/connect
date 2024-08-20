package dydx_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/apis/dydx"
	"github.com/skip-mev/connect/v2/providers/apis/marketmap"
	apihandlers "github.com/skip-mev/connect/v2/providers/base/api/handlers"
	apihandlermocks "github.com/skip-mev/connect/v2/providers/base/api/handlers/mocks"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
	apimetricsmocks "github.com/skip-mev/connect/v2/providers/base/api/metrics/mocks"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	mmclient "github.com/skip-mev/connect/v2/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

func TestDefaultSwitchOverProvider(t *testing.T) {
	cases := []struct {
		name    string
		logger  *zap.Logger
		api     config.APIConfig
		rh      apihandlers.RequestHandler
		metrics metrics.APIMetrics
		err     bool
	}{
		{
			name:    "nil logger",
			logger:  nil,
			api:     config.APIConfig{},
			rh:      nil,
			metrics: nil,
			err:     true,
		},
		{
			name:    "wrong api name",
			logger:  zap.NewNop(),
			api:     dydx.DefaultAPIConfig,
			rh:      nil,
			metrics: nil,
			err:     true,
		},
		{
			name:   "missing endpoints",
			logger: zap.NewNop(),
			api: config.APIConfig{
				Name:             dydx.SwitchOverAPIHandlerName,
				Atomic:           true,
				Enabled:          true,
				Timeout:          20 * time.Second, // Set a high timeout to account for slow API responses in the case where many markets are queried.
				Interval:         10 * time.Second,
				ReconnectTimeout: 2000 * time.Millisecond,
				MaxQueries:       1,
				Endpoints: []config.Endpoint{
					{
						URL: "http://localhost:1317", // REST endpoint (HTTP/HTTPS prefix)
					},
				},
			},
			rh:      nil,
			metrics: nil,
			err:     true,
		},
		{
			name:    "nil request handler",
			logger:  zap.NewNop(),
			api:     dydx.DefaultSwitchOverAPIConfig,
			rh:      nil,
			metrics: metrics.NewNopAPIMetrics(),
			err:     true,
		},
		{
			name:    "nil metrics",
			logger:  zap.NewNop(),
			api:     dydx.DefaultSwitchOverAPIConfig,
			rh:      apihandlermocks.NewRequestHandler(t),
			metrics: nil,
			err:     true,
		},
		{
			name:    "valid",
			logger:  zap.NewNop(),
			api:     dydx.DefaultSwitchOverAPIConfig,
			rh:      apihandlermocks.NewRequestHandler(t),
			metrics: metrics.NewNopAPIMetrics(),
			err:     false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := dydx.NewDefaultSwitchOverMarketMapFetcher(tc.logger, tc.api, tc.rh, tc.metrics)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestNewSwitchOverAPIHandler(t *testing.T) {
	cases := []struct {
		name             string
		logger           *zap.Logger
		pricesFetcher    mmclient.MarketMapFetcher
		marketmapFetcher mmclient.MarketMapFetcher
		err              bool
	}{
		{
			name:             "nil logger",
			logger:           nil,
			pricesFetcher:    nil,
			marketmapFetcher: nil,
			err:              true,
		},
		{
			name:             "nil prices fetcher",
			logger:           zap.NewNop(),
			pricesFetcher:    nil,
			marketmapFetcher: apihandlermocks.NewAPIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse](t),
			err:              true,
		},
		{
			name:             "nil marketmap fetcher",
			logger:           zap.NewNop(),
			pricesFetcher:    apihandlermocks.NewAPIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse](t),
			marketmapFetcher: nil,
			err:              true,
		},
		{
			name:             "valid",
			logger:           zap.NewNop(),
			pricesFetcher:    apihandlermocks.NewAPIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse](t),
			marketmapFetcher: apihandlermocks.NewAPIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse](t),
			err:              false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := dydx.NewSwitchOverFetcher(tc.logger, tc.pricesFetcher, tc.marketmapFetcher, metrics.NewNopAPIMetrics())
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestSwitchOverProvider_Fetch(t *testing.T) {
	pf := apihandlermocks.NewAPIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse](t)
	mmf := apihandlermocks.NewAPIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse](t)
	metrics := apimetricsmocks.NewAPIMetrics(t)

	fetcher, err := dydx.NewSwitchOverFetcher(zap.NewNop(), pf, mmf, metrics)
	require.NoError(t, err)

	cases := []struct {
		name             string
		pricesFetcher    func(*apihandlermocks.APIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse])
		marketmapFetcher func(*apihandlermocks.APIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse])
		metrics          func(*apimetricsmocks.APIMetrics)
		resp             mmclient.MarketMapResponse
	}{
		{
			name: "market map returns no resolved markets",
			pricesFetcher: func(pf *apihandlermocks.APIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse]) {
				resp := mmclient.NewMarketMapResponse(
					mmclient.ResolvedMarketMap{
						dydx.DYDXChain: mmclient.NewMarketMapResult(
							&mmtypes.MarketMapResponse{},
							time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						),
					},
					make(mmclient.UnResolvedMarketMap),
				)
				pf.On("Fetch", mock.Anything, mock.Anything).Return(resp).Once()
			},
			marketmapFetcher: func(mmf *apihandlermocks.APIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse]) {
				resp := mmclient.NewMarketMapResponse(
					make(mmclient.ResolvedMarketMap),
					make(mmclient.UnResolvedMarketMap),
				)
				mmf.On("Fetch", mock.Anything, mock.Anything).Return(resp).Once()
			},
			metrics: func(m *apimetricsmocks.APIMetrics) {
				m.On("AddProviderResponse", dydx.Name, strings.ToLower(dydx.DYDXChain.String()), providertypes.OK).Once()
			},
			resp: mmclient.NewMarketMapResponse(
				mmclient.ResolvedMarketMap{
					dydx.DYDXChain: mmclient.NewMarketMapResult(
						&mmtypes.MarketMapResponse{},
						time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					),
				},
				make(mmclient.UnResolvedMarketMap),
			),
		},
		{
			name: "market map returns resolved markets",
			pricesFetcher: func(*apihandlermocks.APIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse]) {
			},
			marketmapFetcher: func(mmf *apihandlermocks.APIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse]) {
				resp := mmclient.NewMarketMapResponse(
					mmclient.ResolvedMarketMap{
						dydx.DYDXChain: mmclient.NewMarketMapResult(
							&mmtypes.MarketMapResponse{},
							time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						),
					},
					make(mmclient.UnResolvedMarketMap),
				)
				mmf.On("Fetch", mock.Anything, mock.Anything).Return(resp).Once()
			},
			metrics: func(m *apimetricsmocks.APIMetrics) {
				m.On("AddProviderResponse", marketmap.Name, strings.ToLower(dydx.DYDXChain.String()), providertypes.OK).Once()
			},
			resp: mmclient.NewMarketMapResponse(
				mmclient.ResolvedMarketMap{
					dydx.DYDXChain: mmclient.NewMarketMapResult(
						&mmtypes.MarketMapResponse{},
						time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					),
				},
				make(mmclient.UnResolvedMarketMap),
			),
		},
		{
			name: "market map returns error after switch over (should not make request to x/prices)",
			pricesFetcher: func(*apihandlermocks.APIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse]) {
			},
			marketmapFetcher: func(mmf *apihandlermocks.APIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse]) {
				resp := mmclient.NewMarketMapResponse(
					make(mmclient.ResolvedMarketMap),
					mmclient.UnResolvedMarketMap{
						dydx.DYDXChain: providertypes.UnresolvedResult{
							ErrorWithCode: providertypes.NewErrorWithCode(
								fmt.Errorf("error"),
								providertypes.ErrorAPIGeneral,
							),
						},
					},
				)
				mmf.On("Fetch", mock.Anything, mock.Anything).Return(resp).Once()
			},
			metrics: func(m *apimetricsmocks.APIMetrics) {
				m.On("AddProviderResponse", marketmap.Name, strings.ToLower(dydx.DYDXChain.String()), providertypes.ErrorAPIGeneral).Once()
			},
			resp: mmclient.NewMarketMapResponse(
				make(mmclient.ResolvedMarketMap),
				mmclient.UnResolvedMarketMap{
					dydx.DYDXChain: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(
							fmt.Errorf("error"),
							providertypes.ErrorAPIGeneral,
						),
					},
				},
			),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.pricesFetcher(pf)
			tc.marketmapFetcher(mmf)
			tc.metrics(metrics)

			resp := fetcher.Fetch(context.Background(), nil)
			require.Equal(t, tc.resp, resp)
		})
	}
}
