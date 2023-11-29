package ve_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/skip-mev/slinky/abci/strategies/compression"
	mockstrategies "github.com/skip-mev/slinky/abci/strategies/currencypair/mocks"
	"github.com/skip-mev/slinky/abci/testutils"
	"github.com/skip-mev/slinky/abci/ve"
	abcitypes "github.com/skip-mev/slinky/abci/ve/types"
	"github.com/skip-mev/slinky/service"
	"github.com/skip-mev/slinky/service/mocks"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	btcUSD     = oracletypes.NewCurrencyPair("BTC", "USD")
	ethUSD     = oracletypes.NewCurrencyPair("ETH", "USD")
	oneHundred = big.NewInt(100)
	twoHundred = big.NewInt(200)

	nilPrices   = map[string]string{}
	singlePrice = map[string]string{
		btcUSD.ToString(): oneHundred.String(),
	}
	multiplePrices = map[string]string{
		btcUSD.ToString(): oneHundred.String(),
		ethUSD.ToString(): twoHundred.String(),
	}
)

type VoteExtenstionTestSuite struct {
	suite.Suite
	ctx sdk.Context
}

func (s *VoteExtenstionTestSuite) SetupTest() {
	s.ctx = testutils.CreateBaseSDKContext(s.T())
}

func TestVoteExtensionTestSuite(t *testing.T) {
	suite.Run(t, new(VoteExtenstionTestSuite))
}

func (s *VoteExtenstionTestSuite) TestExtendVoteExtension() {
	cases := []struct {
		name             string
		oracleService    func() service.OracleService
		currencyPairID   func() *mockstrategies.CurrencyPairStrategy
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
			currencyPairID: func() *mockstrategies.CurrencyPairStrategy {
				return mockstrategies.NewCurrencyPairStrategy(s.T())
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
			currencyPairID: func() *mockstrategies.CurrencyPairStrategy {
				cps := mockstrategies.NewCurrencyPairStrategy(s.T())

				cps.On("ID", mock.Anything, btcUSD).Return(uint64(0), nil)

				return cps
			},
			expectedResponse: &abcitypes.OracleVoteExtension{
				Prices: map[uint64][]byte{
					0: oneHundred.Bytes(),
				},
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
			currencyPairID: func() *mockstrategies.CurrencyPairStrategy {
				cps := mockstrategies.NewCurrencyPairStrategy(s.T())

				cps.On("ID", mock.Anything, btcUSD).Return(uint64(0), nil)
				cps.On("ID", mock.Anything, ethUSD).Return(uint64(1), nil)

				return cps
			},
			expectedResponse: &abcitypes.OracleVoteExtension{
				Prices: map[uint64][]byte{
					0: oneHundred.Bytes(),
					1: twoHundred.Bytes(),
				},
			},
		},
		{
			name: "oracle service panics",
			oracleService: func() service.OracleService {
				mockServer := mocks.NewOracleService(s.T())

				mockServer.On("Prices", mock.Anything, mock.Anything).Panic("panic")

				return mockServer
			},
			currencyPairID: func() *mockstrategies.CurrencyPairStrategy {
				return mockstrategies.NewCurrencyPairStrategy(s.T())
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
			currencyPairID: func() *mockstrategies.CurrencyPairStrategy {
				return mockstrategies.NewCurrencyPairStrategy(s.T())
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
			currencyPairID: func() *mockstrategies.CurrencyPairStrategy {
				return mockstrategies.NewCurrencyPairStrategy(s.T())
			},
			expectedResponse: &abcitypes.OracleVoteExtension{
				Prices: nil,
			},
		},
		{
			name: "currency pair id strategy returns an error",
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
			currencyPairID: func() *mockstrategies.CurrencyPairStrategy {
				cps := mockstrategies.NewCurrencyPairStrategy(s.T())

				cps.On("ID", mock.Anything, btcUSD).Return(uint64(0), fmt.Errorf("error"))
				cps.On("ID", mock.Anything, ethUSD).Return(uint64(1), nil)

				return cps
			},
			expectedResponse: &abcitypes.OracleVoteExtension{
				Prices: map[uint64][]byte{
					1: twoHundred.Bytes(),
				},
			},
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			codec := compression.NewCompressionVoteExtensionCodec(
				compression.NewDefaultVoteExtensionCodec(),
				compression.NewZLibCompressor(),
			)

			h := ve.NewVoteExtensionHandler(
				log.NewTestLogger(s.T()),
				tc.oracleService(),
				time.Second*1,
				tc.currencyPairID(),
				codec,
			)

			resp, err := h.ExtendVoteHandler()(s.ctx, &cometabci.RequestExtendVote{})
			s.Require().NoError(err)
			s.Require().NotNil(resp)

			ve, err := codec.Decode(resp.VoteExtension)
			s.Require().NoError(err)

			s.Require().Equal(tc.expectedResponse.Prices, ve.Prices)
		})
	}
}

func (s *VoteExtenstionTestSuite) TestVerifyVoteExtension() {
	codec := compression.NewCompressionVoteExtensionCodec(
		compression.NewDefaultVoteExtensionCodec(),
		compression.NewZLibCompressor(),
	)

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
				prices := map[uint64][]byte{
					0: oneHundred.Bytes(),
					1: twoHundred.Bytes(),
				}

				ve, err := testutils.CreateVoteExtensionBytes(
					prices,
					codec,
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
				prices := map[uint64][]byte{}

				ve, err := testutils.CreateVoteExtensionBytes(
					prices,
					codec,
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
			name: "vote extension with malformed prices",
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				prices := map[uint64][]byte{
					0: make([]byte, 34),
				}

				ve, err := testutils.CreateVoteExtensionBytes(
					prices,
					codec,
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

	cpID := mockstrategies.NewCurrencyPairStrategy(s.T())

	for _, tc := range cases {
		s.Run(tc.name, func() {
			handler := ve.NewVoteExtensionHandler(
				log.NewTestLogger(s.T()),
				mocks.NewOracleService(s.T()),
				time.Second*1,
				cpID,
				codec,
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
