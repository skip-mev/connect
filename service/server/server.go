package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	gateway "github.com/cosmos/gogogateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/skip-mev/slinky/pkg/sync"
	"github.com/skip-mev/slinky/service"
	servicetypes "github.com/skip-mev/slinky/service/types"
)

const DefaultServerShutdownTimeout = 3 * time.Second

// OracleServer is the base implementation of the service.OracleServer interface, this is meant to
// serve requests from a remote OracleClient
type OracleServer struct {
	service.UnimplementedOracleServer

	// expected implementation of the oracle
	o servicetypes.Oracle

	// underlying grpc-server -- serves all grpc requests
	grpcSrv *grpc.Server

	// grpc-gateway mux -- serves all http grpc proxy requests
	gatewayMux *runtime.ServeMux

	// underlying http server
	httpSrv *http.Server

	// closer to handle graceful closures from multiple go-routines
	*sync.Closer

	// logger to log incoming requests
	logger *zap.Logger
}

// NewOracleServer returns a new instance of the OracleServer, given an implementation of the Oracle interface.
func NewOracleServer(o servicetypes.Oracle, logger *zap.Logger) *OracleServer {
	logger = logger.With(zap.String("server", "oracle"))

	os := &OracleServer{
		o:      o,
		logger: logger,
	}
	os.Closer = sync.NewCloser().WithCallback(func() {
		// if the server has been started, close it
		if os.httpSrv != nil {
			ctx, cf := context.WithTimeout(context.Background(), DefaultServerShutdownTimeout)
			_ = os.httpSrv.Shutdown(ctx)
			cf()
		}
	})

	return os
}

// routeRequest determines if the incoming http request is a grpc or http request and routes to the proper handler
func (os *OracleServer) routeRequest(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.HasPrefix(
		r.Header.Get("Content-Type"), "application/grpc") {
		os.grpcSrv.ServeHTTP(w, r)
	} else {
		os.gatewayMux.ServeHTTP(w, r)
	}
}

// StartServer starts the oracle gRPC server on the given host and port. The server is killed on any errors from the listener, or if ctx is cancelled.
// This method returns an error via any failure from the listener. This is a blocking call, i.e until the server is closed or the server errors,
// this method will block.
func (os *OracleServer) StartServer(ctx context.Context, host, port string) error {
	serverEndpoint := fmt.Sprintf("%s:%s", host, port)
	os.httpSrv = &http.Server{
		Addr:              serverEndpoint,
		ReadHeaderTimeout: DefaultServerShutdownTimeout,
	}
	// create grpc server
	os.grpcSrv = grpc.NewServer()
	// register oracle server
	service.RegisterOracleServer(os.grpcSrv, os)

	// register the grpc-gateway
	// it handles the http request and dials the server endpoint with the grpc request
	os.gatewayMux = runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &gateway.JSONPb{
			EmitDefaults: true,
			Indent:       "",
			OrigName:     true,
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := service.RegisterOracleHandlerFromEndpoint(ctx, os.gatewayMux, serverEndpoint, opts)
	if err != nil {
		return err
	}

	router := http.NewServeMux()
	router.HandleFunc("/", os.routeRequest)

	os.httpSrv.Handler = h2c.NewHandler(router, &http2.Server{})

	eg, ctx := errgroup.WithContext(ctx)

	// listen for ctx cancellation
	eg.Go(func() error {
		// if the context is closed, close the server + oracle
		<-ctx.Done()
		os.logger.Info("context cancelled, closing oracle")

		_ = os.Close()
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

		err = os.httpSrv.ListenAndServe()
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
func (os *OracleServer) Prices(ctx context.Context, req *service.QueryPricesRequest) (*service.QueryPricesResponse, error) {
	// check that the request is non-nil
	if req == nil {
		return nil, servicetypes.ErrorNilRequest
	}

	os.logger.Info("received request for prices")

	// check that oracle is running
	if !os.o.IsRunning() {
		os.logger.Error("oracle not running")
		return nil, servicetypes.ErrorOracleNotRunning
	}

	resCh := make(chan *service.QueryPricesResponse)

	// run the request in a goroutine, to unblock server + ctx cancellation
	go func() {
		// get the prices
		prices := os.o.GetPrices()

		// get the latest timestamp of the latest update from the oracle
		timestamp := os.o.GetLastSyncTime()

		resCh <- &service.QueryPricesResponse{
			Prices:    servicetypes.ToReqPrices(prices),
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
