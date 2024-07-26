package osmosis

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/pkg/http"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
)

var (
	_ Client = &ClientImpl{}
	_ Client = &MultiClientImpl{}
)

// Client is the expected interface for an osmosis client.
//
//go:generate mockery --name Client --output ./mocks/ --case underscore
type Client interface {
	SpotPrice(ctx context.Context,
		poolID uint64,
		baseAsset,
		quoteAsset string,
	) (SpotPriceResponse, error)
}

// ClientImpl is an implementation of a client to Osmosis using a
// poolmanager Query Client.
type ClientImpl struct {
	api         config.APIConfig
	apiMetrics  metrics.APIMetrics
	redactedURL string
	endpoint    config.Endpoint
	httpClient  *http.Client
}

func NewClient(
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
	endpoint config.Endpoint,
) (Client, error) {
	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	if api.Name != Name {
		return nil, fmt.Errorf("invalid config: name (%s) expected (%s)", api.Name, Name)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("invalid config: disabled (%v)", api.Enabled)
	}

	if apiMetrics == nil {
		return nil, fmt.Errorf("invalid config: apiMetrics is nil")
	}

	redactedURL := metrics.RedactedEndpointURL(0)

	return &ClientImpl{
		api:         api,
		apiMetrics:  apiMetrics,
		redactedURL: redactedURL,
		endpoint:    endpoint,
		httpClient:  http.NewClient(),
	}, nil
}

// SpotPrice uses the underlying x/poolmanager client to access spot prices.
func (c *ClientImpl) SpotPrice(ctx context.Context, poolID uint64, baseAsset, quoteAsset string) (SpotPriceResponse, error) {
	start := time.Now()
	defer func() {
		c.apiMetrics.ObserveProviderResponseLatency(c.api.Name, c.redactedURL, time.Since(start))
	}()

	url, err := CreateURL(c.endpoint.URL, poolID, baseAsset, quoteAsset)
	if err != nil {
		return SpotPriceResponse{}, err
	}

	resp, err := c.httpClient.GetWithContext(ctx, url)
	if err != nil {
		return SpotPriceResponse{}, err
	}

	c.apiMetrics.AddHTTPStatusCode(c.api.Name, resp)

	var spotPriceResponse SpotPriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&spotPriceResponse); err != nil {
		return SpotPriceResponse{}, err
	}

	c.apiMetrics.AddHTTPStatusCode(c.api.Name, resp)
	return spotPriceResponse, nil
}

// MultiClientImpl is an Osmosis client that wraps a set of multiple Clients.
type MultiClientImpl struct {
	logger     *zap.Logger
	api        config.APIConfig
	apiMetrics metrics.APIMetrics

	clients []Client
}

// NewMultiClient creates a new Client.
func NewMultiClient(
	logger *zap.Logger,
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
	clients []Client,
) (Client, error) {
	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	if api.Name != Name {
		return nil, fmt.Errorf("invalid config: name (%s) expected (%s)", api.Name, Name)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("invalid config: disabled (%v)", api.Enabled)
	}

	if apiMetrics == nil {
		return nil, fmt.Errorf("invalid config: apiMetrics is nil")
	}

	return &MultiClientImpl{
		logger:     logger,
		api:        api,
		apiMetrics: apiMetrics,
		clients:    clients,
	}, nil
}

// NewMultiClientFromEndpoints creates a new Client from a list of endpoints.
func NewMultiClientFromEndpoints(
	logger *zap.Logger,
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
) (Client, error) {
	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	if api.Name != Name {
		return nil, fmt.Errorf("invalid config: name (%s) expected (%s)", api.Name, Name)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("invalid config: disabled (%v)", api.Enabled)
	}

	if apiMetrics == nil {
		return nil, fmt.Errorf("invalid config: apiMetrics is nil")
	}

	clients := make([]Client, 0, len(api.Endpoints))
	for _, endpoint := range api.Endpoints {
		c, err := NewClient(api, apiMetrics, endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to create client: %w", err)
		}

		clients = append(clients, c)
	}

	return &MultiClientImpl{
		logger:     logger,
		api:        api,
		apiMetrics: apiMetrics,
		clients:    clients,
	}, nil
}

// SpotPrice delegates the request to all underlying clients and applies a filter to the
// set of responses.
func (mc *MultiClientImpl) SpotPrice(ctx context.Context, poolID uint64, baseAsset, quoteAsset string) (SpotPriceResponse, error) {
	resps := make([]SpotPriceResponse, len(mc.clients))

	var wg sync.WaitGroup
	wg.Add(len(mc.clients))

	for i := range mc.clients {
		url := mc.api.Endpoints[i].URL

		index := i
		go func(index int, client Client) {
			defer wg.Done()
			resp, err := client.SpotPrice(ctx, poolID, baseAsset, quoteAsset)
			if err != nil {
				mc.logger.Error("failed to spot price in sub client", zap.String("url", url), zap.Error(err))
				return
			}

			mc.logger.Debug("successfully fetched accounts", zap.String("url", url))

			resps[index] = resp
		}(index, mc.clients[i])
	}

	wg.Wait()

	return filterSpotPriceResponses(resps)
}

// filterSpotPriceResponses currently just chooses a random response as there is no way to differentiate.
func filterSpotPriceResponses(responses []SpotPriceResponse) (SpotPriceResponse, error) {
	if len(responses) == 0 {
		return SpotPriceResponse{}, fmt.Errorf("no responses found")
	}

	perm := rand.Perm(len(responses))
	for _, i := range perm {
		resp := responses[perm[i]]
		if resp.SpotPrice != "" {
			return resp, nil
		}
	}

	return SpotPriceResponse{}, fmt.Errorf("no responses found")
}
