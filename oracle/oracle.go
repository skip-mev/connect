package oracle

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/oracle/provider"
	"github.com/skip-mev/slinky/oracle/types"
	ssync "github.com/skip-mev/slinky/pkg/sync"
	"golang.org/x/sync/errgroup"
)

// Oracle implements the core component responsible for fetching exchange rates
// for a given set of tickers and determining exchange rates.
type Oracle struct {
	logger          log.Logger
	closer          *ssync.Closer
	providerTimeout time.Duration
	oracleTicker    time.Duration
	lastPriceSync   time.Time
	providerPairs   map[string][]types.CurrencyPair
	endpoints       map[string]provider.Endpoint
	priceProviders  map[string]provider.Provider
	requiredPairs   []types.CurrencyPair

	mtx    sync.RWMutex
	prices map[string]sdk.Dec
}

func New(
	logger log.Logger,
	pTimeout, oTicker time.Duration,
	providerPairs map[string][]types.CurrencyPair,
	endpoints map[string]provider.Endpoint,
	requiredPairs []types.CurrencyPair,
) *Oracle {
	return &Oracle{
		logger:          logger.With("module", "oracle"),
		closer:          ssync.NewCloser(),
		providerTimeout: pTimeout,
		oracleTicker:    oTicker,
		providerPairs:   providerPairs,
		endpoints:       endpoints,
		requiredPairs:   requiredPairs,
		priceProviders:  make(map[string]provider.Provider),
	}
}

// Start starts the (blocking) oracle process. It will return when the context
// is cancelled or the oracle is stopped.
func (o *Oracle) Start(ctx context.Context) error {
	ticker := time.NewTicker(o.oracleTicker)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			o.closer.Close()
			return ctx.Err()

		case <-o.closer.Done():
			return nil

		case <-ticker.C:
			o.logger.Debug("starting oracle tick")

			if err := o.tick(ctx); err != nil {
				o.logger.Error("oracle tick failed", "err", err)
			}

			o.lastPriceSync = time.Now().UTC()
		}
	}
}

// Stop stops the oracle process and waits for it to gracefully exit.
func (o *Oracle) Stop() {
	o.closer.Close()
	<-o.closer.Done()
}

// tick executes a single oracle tick. It fetches prices from all configured
// providers and computes the aggregated price for each pair (ticker) per provider.
func (o *Oracle) tick(ctx context.Context) error {
	o.logger.Debug("executing oracle tick")

	g := new(errgroup.Group)
	priceAgg := types.NewPriceAggregator()

	// How an application determines which providers to use and for which pairs
	// can be done in a variety of ways. For demo purposes, we presume they are
	// locally configured. However, providers can be governed by governance.
	for providerName, currencyPairs := range o.providerPairs {
		priceProvider, err := o.getOrSetProvider(ctx, providerName)
		if err != nil {
			return err
		}

		// Launch a goroutine to fetch ticker prices and candles from the provider
		// for the given set of tickers.
		g.Go(func() error {
			doneCh := make(chan bool, 1)
			errCh := make(chan error, 1)

			var (
				prices  map[string]types.TickerPrice
				candles map[string][]types.Candle
				err     error
			)

			go func() {
				prices, err = priceProvider.GetTickerPrices(currencyPairs...)
				if err != nil {
					o.logger.Error("failed to fetch ticker prices from provider", "provider", providerName, "err", err)
					errCh <- err
				}

				candles, err = priceProvider.GetCandlePrices(currencyPairs...)
				if err != nil {
					o.logger.Error("failed to fetch candle prices from provider", "provider", providerName, "err", err)
					errCh <- err
				}

				doneCh <- true
			}()

			select {
			case <-doneCh:
				break

			case err := <-errCh:
				return err

			case <-time.After(o.providerTimeout):
				return fmt.Errorf("provider %s timed out", providerName)
			}

			// aggregate and collect prices based on the base currency per provider
			for _, pair := range currencyPairs {
				success := priceAgg.SetTickerPricesAndCandles(providerName, prices, candles, pair)
				if !success {
					return fmt.Errorf("failed to find any exchange rates in provider responses")
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	// compute aggregated prices
	computedPrices, err := o.ComputeOraclePrices(priceAgg)
	if err != nil {
		return err
	}

	// ensure we have fetched aggregated prices for all required pairs
	for _, cp := range o.requiredPairs {
		if _, ok := computedPrices[cp.Base]; !ok {
			return fmt.Errorf("failed to find price for %s", cp.Base)
		}
	}

	o.SetPrices(computedPrices)

	return nil
}

func (o *Oracle) SetPrices(prices map[string]sdk.Dec) {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	o.prices = prices
}

func (o *Oracle) GetPrices() map[string]sdk.Dec {
	o.mtx.RLock()
	defer o.mtx.RUnlock()

	p := make(map[string]sdk.Dec, len(o.prices))
	for k, v := range o.prices {
		p[k] = v
	}

	return p
}

// ComputeOraclePrices takes aggregated price points and candles from all
// providers returns the aggregated price per ticker. TVWAP on candles is preferred.
// If we cannot compute TVWAP, we fallback to VWAP on price points..
func (o *Oracle) ComputeOraclePrices(providerAgg *types.PriceAggregator) (prices map[string]sdk.Dec, err error) {
	// attempt to use candles for TVWAP calculations
	tvwapPrices, err := ComputeCandleTVWAP(providerAgg.GetProviderCandles())
	if err != nil {
		return nil, err
	}

	if len(tvwapPrices) > 0 {
		return tvwapPrices, nil
	}

	vwapPrices := ComputeVWAP(providerAgg.GetProviderPrices())

	return vwapPrices, nil
}

func (o *Oracle) getOrSetProvider(ctx context.Context, providerName string) (provider.Provider, error) {
	var (
		priceProvider provider.Provider
		ok            bool
	)

	priceProvider, ok = o.priceProviders[providerName]
	if !ok {
		// TODO: Create providers...
		// o.priceProviders[providerName] = priceProvider
	}

	return priceProvider, nil
}
