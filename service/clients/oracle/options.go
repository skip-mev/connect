package oracle

// Option enables consumers to configure the behavior of an OracleClient on initialization.
type Option func(OracleClient)

// WithBlockingDial configures the OracleClient to block on dialing the remote oracle server.
//
// NOTICE: This option is not recommended to be used in practice. See the [GRPC docs](https://github.com/grpc/grpc-go/blob/master/Documentation/anti-patterns.md)
func WithBlockingDial() Option {
	return func(c OracleClient) {
		client, ok := c.(*GRPCClient)
		if !ok {
			return
		}

		client.blockingDial = true
	}
}
