package oracle

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/types"
)

// generalProvider is an interface for a provider that implements the base provider.
type generalProvider interface {
	// Start starts the provider.
	Start(ctx context.Context) error
	// Name is the provider's name.
	Name() string
}

// Start starts the (blocking) provider orchestrator. This will initialize the provider orchestrator
// with the relevant price and market mapper providers, and then start all of them.
func (o *OracleImpl) Start(ctx context.Context) error {
	o.logger.Info("starting provider orchestrator")
	o.running.Store(true)
	defer o.running.Store(false)

	if err := o.Init(ctx); err != nil {
		o.logger.Error("failed to initialize provider orchestrator", zap.Error(err))
		return err
	}

	// Set the main context for the provider orchestrator.
	ctx, _ = o.setMainCtx(ctx)

	// Start all price providers which have tickers.
	for name, state := range o.providers {
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
	o.metrics.SetSlinkyBuildInfo()

	for {
		select {
		case <-ctx.Done():
			o.Stop()
			o.logger.Info("orchestrator stopped via context")
			return ctx.Err()

		case <-o.closer.Done():
			o.logger.Info("orchestrator stopped via closer")
			return nil

		case <-ticker.C:
			o.fetchAllPrices()
		}
	}
}

// Stop stops the provider orchestrator. This is a synchronous operation that will
// wait for all providers to exit.
func (o *OracleImpl) Stop() {
	o.logger.Info("stopping provider orchestrator")
	if _, cancel := o.getMainCtx(); cancel != nil {
		cancel()
	}

	o.wg.Wait()
	o.logger.Info("provider orchestrator exited successfully")
	o.closer.Close()
	<-o.closer.Done()
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

// setMainCtx sets the main context for the provider orchestrator.
func (o *OracleImpl) setMainCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	o.mut.Lock()
	defer o.mut.Unlock()

	o.mainCtx, o.mainCancel = context.WithCancel(ctx)
	return o.mainCtx, o.mainCancel
}

// getMainCtx returns the main context for the provider orchestrator.
func (o *OracleImpl) getMainCtx() (context.Context, context.CancelFunc) {
	o.mut.Lock()
	defer o.mut.Unlock()

	return o.mainCtx, o.mainCancel
}
