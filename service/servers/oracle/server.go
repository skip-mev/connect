package oracle

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/pkg/sync"
	"github.com/skip-mev/slinky/service/servers/oracle/types"
)

const (
	Transport = "tcp"
)

// OracleServer is the base implementation of the service.OracleServer interface, this is meant to
// serve requests from a remote OracleClient
type OracleServer struct { //nolint
	types.UnimplementedOracleServer

	// expected implementation of the oracle
	o oracle.Oracle

	// underlying grpc-server
	srv *grpc.Server

	// closer to handle graceful closures from multiple go-routines
	*sync.Closer

	// logger to log incoming requests
	logger *zap.Logger
}

// NewOracleServer returns a new instance of the OracleServer, given an implementation of the Oracle interface.
func NewOracleServer(o oracle.Oracle, logger *zap.Logger) *OracleServer {
	logger = logger.With(zap.String("server", "oracle"))

	os := &OracleServer{
		o:      o,
		logger: logger,
	}
	os.Closer = sync.NewCloser().WithCallback(func() {
		// if the server has been started, close it
		if os.srv != nil {
			os.srv.GracefulStop()
		}
	})

	return os
}

// StartServer starts the oracle gRPC server on the given host and port. The server is killed on any errors from the listener, or if ctx is cancelled.
// This method returns an error via any failure from the listener. This is a blocking call, i.e until the server is closed or the server errors,
// this method will block.
func (os *OracleServer) StartServer(ctx context.Context, host, port string) error {
	// create grpc server
	os.srv = grpc.NewServer()

	// register oracle server
	types.RegisterOracleServer(os.srv, os)

	// create listener
	listener, err := net.Listen(Transport, fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return fmt.Errorf("[grpc server]: error creating listener: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	// listen for ctx cancellation
	eg.Go(func() error {
		// if the context is closed, close the server + oracle
		<-ctx.Done()
		os.logger.Info("context cancelled, closing oracle")

		os.Close()
		return nil
	})

	// start the oracle, return error if it fails
	eg.Go(func() error {
		// start the oracle
		os.logger.Info("starting oracle")
		return os.o.Start(ctx)
	})

	// start the server
	eg.Go(func() error {
		// serve, and return any errors
		os.logger.Info(
			"starting grpc server",
			zap.String("host", host),
			zap.String("port", port),
		)

		err := os.srv.Serve(listener)
		if err != nil {
			return fmt.Errorf("[grpc server]: error serving: %v", err)
		}

		return nil
	})

	// wait for everything to finish
	return eg.Wait()
}

// Prices calls the underlying oracle's implementation of GetPrices. It defers to the ctx in the request, and errors if the context is cancelled
// for any reason, or if the oracle errors
func (os *OracleServer) Prices(ctx context.Context, req *types.QueryPricesRequest) (*types.QueryPricesResponse, error) {
	// check that the request is non-nil
	if req == nil {
		return nil, ErrNilRequest
	}

	os.logger.Info("received request for prices")

	// check that oracle is running
	if !os.o.IsRunning() {
		os.logger.Error("oracle not running")
		return nil, ErrOracleNotRunning
	}

	resCh := make(chan *types.QueryPricesResponse)

	// run the request in a goroutine, to unblock server + ctx cancellation
	go func() {
		// get the prices
		prices := os.o.GetPrices()

		// get the latest timestamp of the latest update from the oracle
		timestamp := os.o.GetLastSyncTime()

		resCh <- &types.QueryPricesResponse{
			Prices:    ToReqPrices(prices),
			Timestamp: timestamp,
		}
	}()

	// defer to context closure
	select {
	case <-ctx.Done():
		os.logger.Error("context cancelled")
		return nil, context.Canceled
	case resp := <-resCh:
		return resp, nil
	}
}

// Close closes the underlying oracle server, and blocks until all open requests have been satisfied
func (os *OracleServer) Close() error {
	// close + close server if necessary
	os.Closer.Close()
	return nil
}

// Done returns a channel that is closed when the oracle server is closed
func (os *OracleServer) Done() <-chan struct{} {
	return os.Closer.Done()
}
