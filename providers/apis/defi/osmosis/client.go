package osmosis

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/apis/defi/osmosis/queryproto"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
)

var (
	_ GRPCCLient = &GRPCCLientImpl{}
	_ GRPCCLient = &GRPCMultiClientImpl{}
)

// GRPCCLient is the expected interface for an osmosis grpc client.
//
//go:generate mockery --name GRPCCLient --output ./mocks/ --case underscore
type GRPCCLient interface {
	SpotPrice(grpcCtx context.Context,
		req *queryproto.SpotPriceRequest,
	) (*queryproto.SpotPriceResponse, error)
}

// GRPCCLientImpl is an implementation of a GPRC client to Osmosis using a
// poolmanager Query Client.
type GRPCCLientImpl struct {
	api         config.APIConfig
	apiMetrics  metrics.APIMetrics
	redactedURL string

	pmClient queryproto.QueryClient
}

func NewGRPCCLient(
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
	endpoint config.Endpoint,
) (GRPCCLient, error) {
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

	// TODO set up creds and API keys etc
	cc, err := grpc.NewClient(endpoint.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}

	pmClient := queryproto.NewQueryClient(cc)
	redactedURL := metrics.RedactedEndpointURL(0)

	return &GRPCCLientImpl{
		api:         api,
		apiMetrics:  apiMetrics,
		redactedURL: redactedURL,
		pmClient:    pmClient,
	}, nil
}

// SpotPrice uses the underlying x/poolmanager client to access spot prices.
func (c *GRPCCLientImpl) SpotPrice(grpcCtx context.Context, req *queryproto.SpotPriceRequest) (*queryproto.SpotPriceResponse, error) {
	start := time.Now()
	defer func() {
		c.apiMetrics.ObserveProviderResponseLatency(c.api.Name, c.redactedURL, time.Since(start))
	}()

	resp, err := c.pmClient.SpotPrice(grpcCtx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to spot price: %w", err)
	}

	c.apiMetrics.AddRPCStatusCode(c.api.Name, c.redactedURL, metrics.RPCCodeOK)
	return resp, nil
}

// GRPCMultiClientImpl is an Osmosis GRPC client that wraps a set of multiple Clients.
type GRPCMultiClientImpl struct {
	logger     *zap.Logger
	api        config.APIConfig
	apiMetrics metrics.APIMetrics

	clients []GRPCCLient
}

func NewGRPCMultiClient(
	logger *zap.Logger,
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
) (GRPCCLient, error) {
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

	var clients []GRPCCLient
	for _, endpoint := range api.Endpoints {
		c, err := NewGRPCCLient(api, apiMetrics, endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to create grpc client: %w", err)
		}

		clients = append(clients, c)
	}

	return &GRPCMultiClientImpl{
		logger:     logger,
		api:        api,
		apiMetrics: apiMetrics,
		clients:    clients,
	}, nil
}

// SpotPrice delegates the request to all underlying clients and applies a filter to the
// set of responses.
func (mc *GRPCMultiClientImpl) SpotPrice(grpcCtx context.Context, req *queryproto.SpotPriceRequest) (*queryproto.SpotPriceResponse, error) {
	resps := make([]*queryproto.SpotPriceResponse, len(mc.clients))

	var wg sync.WaitGroup
	wg.Add(len(mc.clients))

	// to do mega parallel
	for i, client := range mc.clients {
		url := mc.api.Endpoints[i].URL
		index := i
		go func(index int, client GRPCCLient) {
			// Observe the latency of the request.
			start := time.Now()
			defer func() {
				wg.Done()
				mc.apiMetrics.ObserveProviderResponseLatency(mc.api.Name, metrics.RedactedEndpointURL(index), time.Since(start))
			}()

			resp, err := client.SpotPrice(grpcCtx, req)
			if err != nil {
				mc.apiMetrics.AddRPCStatusCode(mc.api.Name, metrics.RedactedEndpointURL(index), metrics.RPCCodeError)
				mc.logger.Error("failed to fetch accounts", zap.String("url", url), zap.Error(err))
				return
			}

			mc.apiMetrics.AddRPCStatusCode(mc.api.Name, metrics.RedactedEndpointURL(index), metrics.RPCCodeOK)
			mc.logger.Debug("successfully fetched accounts", zap.String("url", url))

			resps[index] = resp
		}(index, client)

	}

	wg.Wait()

	return filterSpotPriceResponses(resps)
}

// filterSpotPriceResponses currently just chooses a random response as there is no way to differentiate.
// TODO differentiate.
func filterSpotPriceResponses(responses []*queryproto.SpotPriceResponse) (*queryproto.SpotPriceResponse, error) {
	var bestResp *queryproto.SpotPriceResponse

	if len(responses) == 0 {
		return nil, fmt.Errorf("no responses found")
	}

	idx := rand.Intn(len(responses))
	bestResp = responses[idx]

	return bestResp, nil
}
