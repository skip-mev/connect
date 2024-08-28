package oracle

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/types"
)

// Start starts the (blocking) oracle. This will initialize the oracle
// with the relevant price and market mapper providers, and then start all of them.
func (o *OracleImpl) Start(ctx context.Context) error {
	o.logger.Info("starting oracle")
	o.running.Store(true)
	defer o.running.Store(false)
	if err := o.Init(ctx); err != nil {
		o.logger.Error("failed to initialize oracle", zap.Error(err))
		return err
	}

	// Set the main context for the oracle.
	ctx, _ = o.setMainCtx(ctx)

	// Start all price providers which have tickers.
	for name, state := range o.priceProviders {
		providerTickers, err := types.ProviderTickersFromMarketMap(name, o.marketMap)
		if err != nil {
			o.logger.Error("failed to create provider market map", zap.String("provider", name), zap.Error(err))
			return err
		}

		// Update the provider's state.
		_, err = o.UpdateProviderState(providerTickers, state)
		if err != nil {
			o.logger.Error("failed to update provider state", zap.String("provider", name), zap.Error(err))
			return err
		}
	}

	// Start the market map provider.
	if o.mmProvider != nil {
		o.logger.Info("starting marketmap provider")

		o.wg.Add(1)
		go func() {
			defer o.wg.Done()
			o.execProviderFn(ctx, o.mmProvider)
		}()

		o.wg.Add(1)
		go func() {
			defer o.wg.Done()
			o.listenForMarketMapUpdates(ctx)
		}()
	}

	// Start price fetch loop.
	ticker := time.NewTicker(o.cfg.UpdateInterval)
	defer ticker.Stop()
	o.metrics.SetConnectBuildInfo()

	for {
		select {
		case <-ctx.Done():
			o.Stop()
			o.logger.Info("oracle stopped via context")
			return ctx.Err()
		case <-ticker.C:
			o.fetchAllPrices()
		}
	}
}

// Stop stops the oracle. This is a synchronous operation that will
// wait for all providers to exit.
func (o *OracleImpl) Stop() {
	o.logger.Info("stopping oracle")
	if _, cancel := o.getMainCtx(); cancel != nil {
		o.logger.Info("cancelling context")
		cancel()
	}

	o.logger.Info("waiting for routines to stop")
	o.wg.Wait()
	o.logger.Info("oracle exited successfully")
}

func (o *OracleImpl) IsRunning() bool { return o.running.Load() }

// execProviderFn starts a provider and recovers from any panics that occur.
func (o *OracleImpl) execProviderFn(
	ctx context.Context,
	p generalProvider,
) {
	defer func() {
		if r := recover(); r != nil {
			o.logger.Error("recovered from panic", zap.Error(fmt.Errorf("%v", r)))
		}
	}()

	if ctx == nil {
		o.logger.Error("main context is nil; cannot start provider", zap.String("provider", p.Name()))
		return
	}

	err := p.Start(ctx)
	o.logger.Error("provider exited", zap.String("provider", p.Name()), zap.Error(err))
}

// getMainCtx returns the main context for the oracle.
func (o *OracleImpl) getMainCtx() (context.Context, context.CancelFunc) {
	o.mut.Lock()
	defer o.mut.Unlock()

	return o.mainCtx, o.mainCancel
}

// setMainCtx sets the main context for the oracle.
func (o *OracleImpl) setMainCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	o.mut.Lock()
	defer o.mut.Unlock()

	o.mainCtx, o.mainCancel = context.WithCancel(ctx)
	return o.mainCtx, o.mainCancel
}
