package oracle

import (
	"context"
	"fmt"
	"math/big"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var CtxErrors = map[error]struct{}{
	context.Canceled:         {},
	context.DeadlineExceeded: {},
}

// StartProviders starts all providers.
func (o *OracleImpl) StartProviders(ctx context.Context) {
	providerGroup, ctx := errgroup.WithContext(ctx)
	providerGroup.SetLimit(len(o.providers))

	for _, p := range o.providers {
		providerGroup.Go(o.execProviderFn(ctx, p))
	}

	o.providerCh <- providerGroup.Wait()
	close(o.providerCh)
}

// execProvider executes a given provider. The provider continues
// to concurrently run until the context is canceled.
func (o *OracleImpl) execProviderFn(
	ctx context.Context,
	p providertypes.Provider[oracletypes.CurrencyPair, *big.Int],
) func() error {
	return func() error {
		defer func() {
			if r := recover(); r != nil {
				o.logger.Error("recovered from panic", zap.Error(fmt.Errorf("%v", r)))
			}
		}()

		o.logger.Info("starting provider routine", zap.String("name", p.Name()))
		err := p.Start(ctx)
		o.logger.Info("provider exiting", zap.String("name", p.Name()), zap.Error(err))

		// If the context is canceled, or the deadline is exceeded,
		// we want to exit the provider and trigger the error group
		// to exit for all providers.
		if _, ok := CtxErrors[err]; ok {
			return err
		}

		// Otherwise, we gracefully exit the go routine.
		return nil
	}
}
