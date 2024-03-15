package orchestrator

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var CtxErrors = map[error]struct{}{
	context.Canceled:         {},
	context.DeadlineExceeded: {},
}

// Start starts the provider orchestrator. This will initialize the provider orchestrator
// with the relevant price providers and market mapper, and then start all of them.
func (o *ProviderOrchestrator) Start(ctx context.Context) error {
	o.logger.Info("starting provider orchestrator")
	if err := o.Init(); err != nil {
		o.logger.Error("failed to initialize provider orchestrator", zap.Error(err))
		return err
	}

	// Create a new error group for the provider orchestrator.
	o.errGroup, ctx = errgroup.WithContext(ctx)

	// Set tthe main context for the provider orchestrator.
	ctx, _ = o.setMainCtx(ctx)

	// Start all of the price providers.
	for _, state := range o.providers {
		o.errGroup.Go(o.execProviderFn(ctx, state.Provider))
	}

	// Start the market map provider.
	if mapper := o.mapper.Mapper; mapper != nil {
		o.errGroup.Go(o.execProviderFn(ctx, mapper))
		o.errGroup.Go(o.listenForMarketMapUpdates(ctx))
	}

	return nil
}

// Stop stops the provider orchestrator.
func (o *ProviderOrchestrator) Stop() error {
	o.logger.Info("stopping provider orchestrator")
	if _, cancel := o.getMainCtx(); cancel != nil {
		cancel()

		if o.errGroup != nil {
			return nil
		}

		// Wait for all of the price providers to exit.
		if err := o.errGroup.Wait(); err != nil {
			o.logger.Error("provider orchestrator exited with error", zap.Error(err))
			return err
		}

		o.logger.Info("provider orchestrator exited successfully")
	}

	return nil
}

// execProviderFn returns a function that starts the provider. This function is used
// to start the provider in a go routine.
func (o *ProviderOrchestrator) execProviderFn(
	ctx context.Context,
	p GeneralProvider,
) func() error {
	return func() error {
		defer func() {
			if r := recover(); r != nil {
				o.logger.Error("recovered from panic", zap.Error(fmt.Errorf("%v", r)))
			}
		}()

		// If the context is canceled, or the deadline is exceeded,
		// we want to exit the provider and trigger the error group
		// to exit for all providers.
		err := p.Start(ctx)
		if _, ok := CtxErrors[err]; ok {
			return err
		}

		// Otherwise, we gracefully exit the go routine.
		return nil
	}
}
