package orchestrator

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// generalProvider is a interface for a provider that implements the base provider.
type generalProvider interface {
	// Start starts the provider.
	Start(ctx context.Context) error
	// Stop stops the provider.
	Name() string
}

// Start starts the provider orchestrator. This will initialize the provider orchestrator
// with the relevant price and market mapper providers, and then start all of them.
func (o *ProviderOrchestrator) Start(ctx context.Context) error {
	o.logger.Info("starting provider orchestrator")
	if err := o.Init(); err != nil {
		o.logger.Error("failed to initialize provider orchestrator", zap.Error(err))
		return err
	}

	// Set tthe main context for the provider orchestrator.
	ctx, _ = o.setMainCtx(ctx)

	// Start all of the price providers.
	for _, state := range o.providers {
		o.wg.Add(1)
		go func() {
			defer o.wg.Done()
			o.execProviderFn(ctx, state.Provider)
		}()
	}

	// Start the market map provider.
	if o.mmProvider != nil {
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
// wait for all of the providers to exit.
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
