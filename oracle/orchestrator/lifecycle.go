package orchestrator

import (
	"context"
	"fmt"

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

// Start starts the provider orchestrator. This will initialize the provider orchestrator
// with the relevant price and market mapper providers, and then start all of them.
func (o *ProviderOrchestrator) Start(ctx context.Context) error {
	o.logger.Info("starting provider orchestrator")
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

	return nil
}

// Stop stops the provider orchestrator. This is a synchronous operation that will
// wait for all providers to exit.
func (o *ProviderOrchestrator) Stop() {
	o.logger.Info("stopping provider orchestrator")
	if _, cancel := o.getMainCtx(); cancel != nil {
		cancel()
	}

	o.wg.Wait()
	o.logger.Info("provider orchestrator exited successfully")
}

// execProviderFn starts a provider and recovers from any panics that occur.
func (o *ProviderOrchestrator) execProviderFn(
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
func (o *ProviderOrchestrator) setMainCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	o.mut.Lock()
	defer o.mut.Unlock()

	o.mainCtx, o.mainCancel = context.WithCancel(ctx)
	return o.mainCtx, o.mainCancel
}

// getMainCtx returns the main context for the provider orchestrator.
func (o *ProviderOrchestrator) getMainCtx() (context.Context, context.CancelFunc) {
	o.mut.Lock()
	defer o.mut.Unlock()

	return o.mainCtx, o.mainCancel
}
