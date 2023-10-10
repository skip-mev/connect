package ve_test

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/abci/testutils"
	"github.com/skip-mev/slinky/abci/ve"
	abcitypes "github.com/skip-mev/slinky/abci/ve/types"
	"github.com/skip-mev/slinky/service"
	"github.com/skip-mev/slinky/service/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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

type VoteExtenstionTestSuite struct {
	suite.Suite
	ctx sdk.Context
}

func (s *VoteExtenstionTestSuite) SetupTest() {
	s.ctx = testutils.CreateBaseSDKContext(s.T())
}

func TestVoteExtenstionTestSuite(t *testing.T) {
	suite.Run(t, new(VoteExtenstionTestSuite))
}

func (s *VoteExtenstionTestSuite) TestExtendVoteExtension() {
	cases := []struct {
		name             string
		oracleService    func() service.OracleService
		expectedResponse *abcitypes.OracleVoteExtension
	}{
		{
			name: "oracle service returns no prices",
			oracleService: func() service.OracleService {
				mockServer := mocks.NewOracleService(s.T())

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
				mockServer := mocks.NewOracleService(s.T())

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
				mockServer := mocks.NewOracleService(s.T())

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
				mockServer := mocks.NewOracleService(s.T())

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
				mockServer := mocks.NewOracleService(s.T())

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
				mockServer := mocks.NewOracleService(s.T())

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
		s.Run(tc.name, func() {
			h := ve.NewVoteExtensionHandler(
				log.NewTestLogger(s.T()),
				tc.oracleService(),
				time.Second*1,
			)

			resp, err := h.ExtendVoteHandler()(s.ctx, &cometabci.RequestExtendVote{})
			s.Require().NoError(err)
			s.Require().NotNil(resp)

			ve := &abcitypes.OracleVoteExtension{}
			s.Require().NoError(ve.Unmarshal(resp.VoteExtension))
			s.Require().Equal(tc.expectedResponse.Prices, ve.Prices)
		})
	}
}

func (s *VoteExtenstionTestSuite) TestVerifyVoteExtension() {
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

				ve, err := testutils.CreateVoteExtensionBytes(
					prices,
					timestamp,
					1,
				)
				s.Require().NoError(err)

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

				ve, err := testutils.CreateVoteExtensionBytes(
					prices,
					timestamp,
					1,
				)
				s.Require().NoError(err)

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

				ve, err := testutils.CreateVoteExtensionBytes(
					prices,
					timestamp,
					1,
				)
				s.Require().NoError(err)

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

				ve, err := testutils.CreateVoteExtensionBytes(
					prices,
					timestamp,
					0,
				)
				s.Require().NoError(err)

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

				ve, err := testutils.CreateVoteExtensionBytes(
					prices,
					timestamp,
					1,
				)
				s.Require().NoError(err)

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
		s.Run(tc.name, func() {
			handler := ve.NewVoteExtensionHandler(
				log.NewTestLogger(s.T()),
				mocks.NewOracleService(s.T()),
				time.Second*1,
			).VerifyVoteExtensionHandler()

			resp, err := handler(s.ctx, tc.getReq())
			s.Require().Equal(tc.expectedResponse, resp)

			if tc.expectedError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
