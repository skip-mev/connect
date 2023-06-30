package abci_test

import (
	"fmt"
	"time"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/skip-mev/slinky/abci"
	abcitypes "github.com/skip-mev/slinky/abci/types"
	"github.com/skip-mev/slinky/service"
	mocks "github.com/skip-mev/slinky/service/mocks"
	"github.com/stretchr/testify/mock"
)

var (
	nilPrices   = map[string]string{}
	singlePrice = map[string]string{
		"BTC/USD": "100",
	}
	multiplePrices = map[string]string{
		"BTC/USD": "100",
		"ETH/USD": "200",
	}
)

func (suite *ABCITestSuite) TestExtendVoteExtension() {
	cases := []struct {
		name             string
		oracleService    func() service.OracleService
		expectedResponse *abcitypes.OracleVoteExtension
	}{
		{
			name: "oracle service returns no prices",
			oracleService: func() service.OracleService {
				mockServer := mocks.NewOracleService(suite.T())

				mockServer.On("Prices", mock.Anything, mock.Anything).Return(
					&service.QueryPricesResponse{
						Prices: nilPrices,
					},
					nil,
				)

				return mockServer
			},
			expectedResponse: &abcitypes.OracleVoteExtension{
				Prices: nil,
			},
		},
		{
			name: "oracle service returns a single price",
			oracleService: func() service.OracleService {
				mockServer := mocks.NewOracleService(suite.T())

				mockServer.On("Prices", mock.Anything, mock.Anything).Return(
					&service.QueryPricesResponse{
						Prices: singlePrice,
					},
					nil,
				)

				return mockServer
			},
			expectedResponse: &abcitypes.OracleVoteExtension{
				Prices: singlePrice,
			},
		},
		{
			name: "oracle service returns multiple prices",
			oracleService: func() service.OracleService {
				mockServer := mocks.NewOracleService(suite.T())

				mockServer.On("Prices", mock.Anything, mock.Anything).Return(
					&service.QueryPricesResponse{
						Prices: multiplePrices,
					},
					nil,
				)

				return mockServer
			},
			expectedResponse: &abcitypes.OracleVoteExtension{
				Prices: multiplePrices,
			},
		},
		{
			name: "oracle service panics",
			oracleService: func() service.OracleService {
				mockServer := mocks.NewOracleService(suite.T())

				mockServer.On("Prices", mock.Anything, mock.Anything).Panic("panic")

				return mockServer
			},
			expectedResponse: &abcitypes.OracleVoteExtension{
				Prices: nil,
			},
		},
		{
			name: "oracle service returns an nil response",
			oracleService: func() service.OracleService {
				mockServer := mocks.NewOracleService(suite.T())

				mockServer.On("Prices", mock.Anything, mock.Anything).Return(
					nil,
					nil,
				)

				return mockServer
			},
			expectedResponse: &abcitypes.OracleVoteExtension{
				Prices: nil,
			},
		},
		{
			name: "oracle service returns an error",
			oracleService: func() service.OracleService {
				mockServer := mocks.NewOracleService(suite.T())

				mockServer.On("Prices", mock.Anything, mock.Anything).Return(
					nil,
					fmt.Errorf("error"),
				)

				return mockServer
			},
			expectedResponse: &abcitypes.OracleVoteExtension{
				Prices: nil,
			},
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			h := abci.NewVoteExtensionHandler(
				log.NewTestLogger(suite.T()),
				tc.oracleService(),
				time.Second*1,
			)

			resp, err := h.ExtendVoteHandler()(suite.ctx, &cometabci.RequestExtendVote{})
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)

			ve := &abcitypes.OracleVoteExtension{}
			suite.Require().NoError(ve.Unmarshal(resp.VoteExtension))
			suite.Require().Equal(tc.expectedResponse.Prices, ve.Prices)
		})
	}
}

func (suite *ABCITestSuite) TestVerifyVoteExtension() {
	cases := []struct {
		name             string
		getReq           func() *cometabci.RequestVerifyVoteExtension
		expectedResponse *cometabci.ResponseVerifyVoteExtension
		expectedError    bool
	}{
		{
			name: "empty vote extension",
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				return &cometabci.RequestVerifyVoteExtension{}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		{
			name: "malformed bytes",
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: []byte("malformed"),
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
		{
			name: "valid vote extension",
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				prices := map[string]string{
					"BTC/USD": "0x1",
					"ETH/USD": "0x2",
				}
				timestamp := time.Now()

				ve := suite.createVoteExtensionBytes(
					prices,
					timestamp,
					1,
				)

				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: ve,
					Height:        1,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		{
			name: "vote extension with no prices",
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				prices := map[string]string{}
				timestamp := time.Now()

				ve := suite.createVoteExtensionBytes(
					prices,
					timestamp,
					1,
				)

				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: ve,
					Height:        1,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		{
			name: "vote extension with invalid prices",
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				prices := map[string]string{
					"Bitcoin": "0x1",
				}
				timestamp := time.Now()

				ve := suite.createVoteExtensionBytes(
					prices,
					timestamp,
					1,
				)

				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: ve,
					Height:        1,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
		{
			name: "vote extension with invalid heights",
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				prices := map[string]string{
					"BTC/USD": "0x1",
					"ETH/USD": "0x2",
				}
				timestamp := time.Now()

				ve := suite.createVoteExtensionBytes(
					prices,
					timestamp,
					0,
				)

				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: ve,
					Height:        1,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
		{
			name: "vote extension with malformed prices",
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				prices := map[string]string{
					"BTC/USD": "malformed",
				}
				timestamp := time.Now()

				ve := suite.createVoteExtensionBytes(
					prices,
					timestamp,
					1,
				)

				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: ve,
					Height:        1,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			handler := abci.NewVoteExtensionHandler(
				log.NewTestLogger(suite.T()),
				mocks.NewOracleService(suite.T()),
				time.Second*1,
			).VerifyVoteExtensionHandler()

			resp, err := handler(suite.ctx, tc.getReq())
			suite.Require().Equal(tc.expectedResponse, resp)

			if tc.expectedError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
