package ethmulticlient_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/apis/defi/ethmulticlient"
	"github.com/skip-mev/connect/v2/providers/apis/defi/ethmulticlient/mocks"
)

func TestMultiClient(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	testcases := []struct {
		name            string
		client          ethmulticlient.EVMClient
		args            []rpc.BatchElem
		expectedResults []interface{}
		err             error
	}{
		{
			name: "no elems, no-ops",
			client: ethmulticlient.NewMultiRPCClient(
				logger,
				config.APIConfig{},
				[]ethmulticlient.EVMClient{},
			),
			args: []rpc.BatchElem{},
			err:  nil,
		},
		{
			name: "single client failure height request",
			client: ethmulticlient.NewMultiRPCClient(
				logger,
				config.APIConfig{
					Endpoints: []config.Endpoint{{URL: "http://localhost:8545"}},
				},
				[]ethmulticlient.EVMClient{
					createEVMClientWithResponse(
						t,
						nil,
						[]string{"", ""},
						[]error{nil, fmt.Errorf("height req failed")},
					),
				},
			),
			args: []rpc.BatchElem{{}},
			err:  fmt.Errorf("endpoint request failed"),
		},
		{
			name: "single client failure hex height decode",
			client: ethmulticlient.NewMultiRPCClient(
				logger,
				config.APIConfig{
					Endpoints: []config.Endpoint{{URL: "http://localhost:8545"}},
				},
				[]ethmulticlient.EVMClient{
					createEVMClientWithResponse(
						t,
						nil,
						[]string{"", "zzzzzz"},
						[]error{nil, nil},
					),
				},
			),
			args: []rpc.BatchElem{{}},
			err:  fmt.Errorf("could not decode hex eth height"),
		},
		{
			name: "single client success",
			client: ethmulticlient.NewMultiRPCClient(
				logger,
				config.APIConfig{
					Endpoints: []config.Endpoint{{URL: "http://localhost:8545"}},
				},
				[]ethmulticlient.EVMClient{
					createEVMClientWithResponse(
						t,
						nil,
						[]string{"some value", "0x12c781c"},
						[]error{nil, nil},
					),
				},
			),
			args:            []rpc.BatchElem{{}},
			expectedResults: []interface{}{"some value"},
			err:             nil,
		},
		{
			name: "two clients one failed height request",
			client: ethmulticlient.NewMultiRPCClient(
				logger,
				config.APIConfig{
					Endpoints: []config.Endpoint{{URL: "http://localhost:8545"}, {URL: "http://localhost:8546"}},
				},
				[]ethmulticlient.EVMClient{
					createEVMClientWithResponse(
						t,
						nil,
						[]string{"", ""},
						[]error{nil, fmt.Errorf("height req failed")},
					),
					createEVMClientWithResponse(
						t,
						nil,
						[]string{"some value", "0x12c781c"},
						[]error{nil, nil},
					),
				},
			),
			args:            []rpc.BatchElem{{}},
			expectedResults: []interface{}{"some value"},
			err:             nil,
		},
		{
			name: "two clients different heights",
			client: ethmulticlient.NewMultiRPCClient(
				logger,
				config.APIConfig{
					Endpoints: []config.Endpoint{{URL: "http://localhost:8545"}, {URL: "http://localhost:8546"}},
				},
				[]ethmulticlient.EVMClient{
					createEVMClientWithResponse(
						t,
						nil,
						[]string{"value1", "0x12c781b"},
						[]error{nil, nil},
					),
					createEVMClientWithResponse(
						t,
						nil,
						[]string{"value2", "0x12c781c"},
						[]error{nil, nil},
					),
				},
			),
			args:            []rpc.BatchElem{{}},
			expectedResults: []interface{}{"value2"},
			err:             nil,
		},
		{
			name: "worst case scenario where no errors were returned but both clients returned height 0",
			client: ethmulticlient.NewMultiRPCClient(
				logger,
				config.APIConfig{
					Endpoints: []config.Endpoint{{URL: "http://localhost:8545"}, {URL: "http://localhost:8546"}},
				},
				[]ethmulticlient.EVMClient{
					createEVMClientWithResponse(
						t,
						nil,
						[]string{"value1", "0x0"},
						[]error{nil, nil},
					),
					createEVMClientWithResponse(
						t,
						nil,
						[]string{"value2", "0x0"},
						[]error{nil, nil},
					),
				},
			),
			args: []rpc.BatchElem{{}},
			err:  fmt.Errorf("no errors were encountered, however no go routine was able to report a height"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.client.BatchCallContext(context.TODO(), tc.args)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				for i, result := range tc.expectedResults {
					require.Equal(t, result, *tc.args[i].Result.(*string))
				}
			}
		})
	}
}

func createEVMClientWithResponse(
	t *testing.T,
	failedRequestErr error,
	responses []string,
	errs []error,
) ethmulticlient.EVMClient {
	t.Helper()

	c := mocks.NewEVMClient(t)
	if failedRequestErr != nil {
		c.On("BatchCallContext", mock.Anything, mock.Anything).Return(failedRequestErr)
	} else {
		c.On("BatchCallContext", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			elems, ok := args.Get(1).([]rpc.BatchElem)
			require.True(t, ok)

			require.True(t, ok)
			require.Equal(t, len(elems), len(responses))
			require.Equal(t, len(elems), len(errs))

			for i, elem := range elems {
				elem.Result = &responses[i]
				elem.Error = errs[i]
				elems[i] = elem
			}
		})
	}

	return c
}
