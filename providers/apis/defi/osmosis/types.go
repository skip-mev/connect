package osmosis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
)

const (
	Name              = "osmosis_api"
	QueryURLCharacter = "?"
	URLSeparator      = "/"
	URLSuffix         = "osmosis/poolmanager/v2/pools/%s/prices%sbase_asset_denom=%s&quote_asset_denom=%s"
)

// CreateURL creates the properly formatted osmosis query URL for spot price.
func CreateURL(baseURL string, poolID uint64, baseAsset, quoteAsset string) (string, error) {
	return strings.Join(
		[]string{
			baseURL,
			fmt.Sprintf(URLSuffix, strconv.FormatUint(poolID, 10), QueryURLCharacter, baseAsset, quoteAsset),
		},
		URLSeparator,
	), nil
}

// NoOsmosisMetadataForTickerError is returned when there is no metadata associated with a given ticker.
func NoOsmosisMetadataForTickerError(ticker string) error {
	return fmt.Errorf("no osmosis metadata for ticker: %s", ticker)
}

// metadataCache is a synchronous data-structure that holds the metadata for each ticker.
type metadataCache struct {
	// metaDataPerTicker is a map from ticker to metadata.
	metaDataPerTicker map[string]TickerMetadata

	// RWMutex is used to synchronize access to the metadata cache.
	mtx sync.RWMutex
}

func newMetadataCache() *metadataCache {
	return &metadataCache{
		metaDataPerTicker: make(map[string]TickerMetadata),
	}
}

func (mc *metadataCache) updateMetaDataCache(ticker types.ProviderTicker) (TickerMetadata, error) {
	mc.mtx.Lock()
	defer mc.mtx.Unlock()
	if metadata, ok := mc.metaDataPerTicker[ticker.String()]; ok {
		return metadata, nil
	}

	metadata, err := unmarshalMetadataJSON(ticker.GetJSON())
	if err != nil {
		return TickerMetadata{}, fmt.Errorf("error unmarshalling metadata for ticker %s: %w", ticker.String(), err)
	}
	if err := metadata.ValidateBasic(); err != nil {
		return TickerMetadata{}, fmt.Errorf("metadata for ticker %s is invalid: %w", ticker.String(), err)
	}
	mc.metaDataPerTicker[ticker.String()] = metadata

	return metadata, nil
}

// getMetadataPerTicker returns the metadata for the given ticker.
func (mc *metadataCache) getMetadataPerTicker(ticker types.ProviderTicker) (TickerMetadata, bool) {
	mc.mtx.RLock()
	defer mc.mtx.RUnlock()

	metaData, ok := mc.metaDataPerTicker[ticker.String()]
	return metaData, ok
}

// TickerMetadata represents the metadata associated with a ticker's corresponding
// raydium pool.
type TickerMetadata struct {
	// PoolID is the unique uint ID of the osmosis pool.
	PoolID uint64 `json:"pool_id"`

	// BaseTokenDenom is the identifier (on osmosis) of the quote token.
	BaseTokenDenom string `json:"base_token_denom"`

	// QuoteTokenDenom is the identifier (on osmosis) of the quote token.
	QuoteTokenDenom string `json:"quote_token_denom"`
}

// ValidateBasic checks that the pool and token information is formatted properly.
func (metadata TickerMetadata) ValidateBasic() error {
	if metadata.BaseTokenDenom == "" || metadata.QuoteTokenDenom == "" {
		return fmt.Errorf("base token denom or quote token denom cannot be empty")
	}

	return nil
}

// unmarshalMetadataJSON unmarshals the given metadata string into a TickerMetadata,
// this method assumes that the metadata string is valid json, otherwise an error is returned.
func unmarshalMetadataJSON(metadata string) (TickerMetadata, error) {
	// unmarshal the metadata string into a TickerMetadata
	var tickerMetadata TickerMetadata
	if err := json.Unmarshal([]byte(metadata), &tickerMetadata); err != nil {
		return TickerMetadata{}, err
	}

	return tickerMetadata, nil
}

// DefaultAPIConfig is the default configuration for the Osmosis API price fetcher.
var DefaultAPIConfig = config.APIConfig{
	Enabled:          true,
	Name:             Name,
	Timeout:          5 * time.Second,
	Interval:         2000 * time.Millisecond, // 2s block times on osmosis
	ReconnectTimeout: 5 * time.Second,
	MaxQueries:       10, // only run 10 queries concurrently to prevent rate limiting
	Atomic:           false,
	BatchSize:        1,
	Endpoints: []config.Endpoint{
		{
			URL: "https://osmosis-api.polkachu.com",
		},
	},
	MaxBlockHeightAge: 30 * time.Second,
}

type SpotPriceResponse struct {
	SpotPrice string `json:"spot_price"`
}

type WrappedSpotPriceResponse struct {
	SpotPriceResponse
	BlockHeight uint64 `json:"block_height"`
}
