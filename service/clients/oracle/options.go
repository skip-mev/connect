package oracle

// Option enables consumers to configure the behavior of an OracleClient on initialization.
type Option func(OracleClient)

// WithBlockingDial configures the OracleClient to block on dialing the remote oracle server.
func WithBlockingDial() Option {
	return func(c OracleClient) {
		client, ok := c.(*GRPCClient)
		if !ok {
			return
		}

		client.blockingDial = true
	}
}
