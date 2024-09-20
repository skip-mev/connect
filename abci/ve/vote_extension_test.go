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

	aggregatormocks "github.com/skip-mev/connect/v2/abci/strategies/aggregator/mocks"
	"github.com/skip-mev/connect/v2/abci/strategies/codec"
	codecmocks "github.com/skip-mev/connect/v2/abci/strategies/codec/mocks"
	mockstrategies "github.com/skip-mev/connect/v2/abci/strategies/currencypair/mocks"
	"github.com/skip-mev/connect/v2/abci/testutils"
	connectabci "github.com/skip-mev/connect/v2/abci/types"
	"github.com/skip-mev/connect/v2/abci/ve"
	abcitypes "github.com/skip-mev/connect/v2/abci/ve/types"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	client "github.com/skip-mev/connect/v2/service/clients/oracle"
	"github.com/skip-mev/connect/v2/service/clients/oracle/mocks"
	servicemetrics "github.com/skip-mev/connect/v2/service/metrics"
	metricsmocks "github.com/skip-mev/connect/v2/service/metrics/mocks"
	servicetypes "github.com/skip-mev/connect/v2/service/servers/oracle/types"
)

var (
	btcUSD     = connecttypes.NewCurrencyPair("BTC", "USD")
	ethUSD     = connecttypes.NewCurrencyPair("ETH", "USD")
	oneHundred = big.NewInt(100)
	twoHundred = big.NewInt(200)

	nilPrices   = map[string]string{}
	singlePrice = map[string]string{
		btcUSD.String(): oneHundred.String(),
	}
	multiplePrices = map[string]string{
		btcUSD.String(): oneHundred.String(),
		ethUSD.String(): twoHundred.String(),
	}
)

type VoteExtensionTestSuite struct {
	suite.Suite
	ctx sdk.Context
}

func (s *VoteExtensionTestSuite) SetupTest() {
	s.ctx = testutils.CreateBaseSDKContext(s.T())
}

func TestVoteExtensionTestSuite(t *testing.T) {
	suite.Run(t, new(VoteExtensionTestSuite))
}

func (s *VoteExtensionTestSuite) TestExtendVoteExtension() {
	cases := []struct {
		name                 string
		oracleService        func() client.OracleClient
		currencyPairStrategy func() *mockstrategies.CurrencyPairStrategy
		expectedResponse     *abcitypes.OracleVoteExtension
		extendVoteRequest    func() *cometabci.RequestExtendVote
		expectedError        bool
	}{
		{
			name: "nil request returns an error",
			oracleService: func() client.OracleClient {
				return mocks.NewOracleClient(s.T())
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				return mockstrategies.NewCurrencyPairStrategy(s.T())
			},
			extendVoteRequest: func() *cometabci.RequestExtendVote { return nil },
		},
		{
			name: "oracle service returns no prices",
			oracleService: func() client.OracleClient {
				mockServer := mocks.NewOracleClient(s.T())

				mockServer.On("Prices", mock.Anything, mock.Anything).Return(
					&servicetypes.QueryPricesResponse{
						Prices: nilPrices,
					},
					nil,
				)

				return mockServer
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				return mockstrategies.NewCurrencyPairStrategy(s.T())
			},
			expectedResponse: &abcitypes.OracleVoteExtension{
				Prices: nil,
			},
		},
		{
			name: "oracle service returns a single price",
			oracleService: func() client.OracleClient {
				mockServer := mocks.NewOracleClient(s.T())

				mockServer.On("Prices", mock.Anything, mock.Anything).Return(
					&servicetypes.QueryPricesResponse{
						Prices: singlePrice,
					},
					nil,
				)

				return mockServer
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				cps := mockstrategies.NewCurrencyPairStrategy(s.T())

				cps.On("ID", mock.Anything, btcUSD).Return(uint64(0), nil)
				cps.On("GetEncodedPrice", mock.Anything, btcUSD, oneHundred).Return(oneHundred.Bytes(), nil)

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
			oracleService: func() client.OracleClient {
				mockServer := mocks.NewOracleClient(s.T())

				mockServer.On("Prices", mock.Anything, mock.Anything).Return(
					&servicetypes.QueryPricesResponse{
						Prices: multiplePrices,
					},
					nil,
				)

				return mockServer
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				cps := mockstrategies.NewCurrencyPairStrategy(s.T())

				cps.On("ID", mock.Anything, btcUSD).Return(uint64(0), nil)
				cps.On("GetEncodedPrice", mock.Anything, btcUSD, oneHundred).Return(oneHundred.Bytes(), nil)

				cps.On("ID", mock.Anything, ethUSD).Return(uint64(1), nil)
				cps.On("GetEncodedPrice", mock.Anything, ethUSD, twoHundred).Return(twoHundred.Bytes(), nil)

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
			oracleService: func() client.OracleClient {
				mockServer := mocks.NewOracleClient(s.T())

				mockServer.On("Prices", mock.Anything, mock.Anything).Panic("panic")

				return mockServer
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				return mockstrategies.NewCurrencyPairStrategy(s.T())
			},
			expectedResponse: &abcitypes.OracleVoteExtension{
				Prices: nil,
			},
			expectedError: true,
		},
		{
			name: "oracle service returns an nil response",
			oracleService: func() client.OracleClient {
				mockServer := mocks.NewOracleClient(s.T())

				mockServer.On("Prices", mock.Anything, mock.Anything).Return(
					nil,
					nil,
				)

				return mockServer
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				return mockstrategies.NewCurrencyPairStrategy(s.T())
			},
			expectedResponse: &abcitypes.OracleVoteExtension{
				Prices: nil,
			},
		},
		{
			name: "oracle service returns an error",
			oracleService: func() client.OracleClient {
				mockServer := mocks.NewOracleClient(s.T())

				mockServer.On("Prices", mock.Anything, mock.Anything).Return(
					nil,
					fmt.Errorf("error"),
				)

				return mockServer
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				return mockstrategies.NewCurrencyPairStrategy(s.T())
			},
			expectedResponse: &abcitypes.OracleVoteExtension{
				Prices: nil,
			},
		},
		{
			name: "currency pair id strategy returns an error",
			oracleService: func() client.OracleClient {
				mockServer := mocks.NewOracleClient(s.T())

				mockServer.On("Prices", mock.Anything, mock.Anything).Return(
					&servicetypes.QueryPricesResponse{
						Prices: multiplePrices,
					},
					nil,
				)

				return mockServer
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				cps := mockstrategies.NewCurrencyPairStrategy(s.T())

				cps.On("ID", mock.Anything, btcUSD).Return(uint64(0), fmt.Errorf("error"))
				cps.On("ID", mock.Anything, ethUSD).Return(uint64(1), nil)
				cps.On("GetEncodedPrice", mock.Anything, ethUSD, twoHundred).Return(twoHundred.Bytes(), nil)

				return cps
			},
			expectedResponse: &abcitypes.OracleVoteExtension{
				Prices: map[uint64][]byte{
					1: twoHundred.Bytes(),
				},
			},
		},
		{
			name: "currency pair price strategy returns an error",
			oracleService: func() client.OracleClient {
				mockServer := mocks.NewOracleClient(s.T())

				mockServer.On("Prices", mock.Anything, mock.Anything).Return(
					&servicetypes.QueryPricesResponse{
						Prices: multiplePrices,
					},
					nil,
				)

				return mockServer
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				cps := mockstrategies.NewCurrencyPairStrategy(s.T())

				cps.On("ID", mock.Anything, btcUSD).Return(uint64(0), nil)
				cps.On("GetEncodedPrice", mock.Anything, btcUSD, oneHundred).Return(nil, fmt.Errorf("error"))

				cps.On("ID", mock.Anything, ethUSD).Return(uint64(1), nil)
				cps.On("GetEncodedPrice", mock.Anything, ethUSD, twoHundred).Return(twoHundred.Bytes(), nil)

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
			cdc := codec.NewCompressionVoteExtensionCodec(
				codec.NewDefaultVoteExtensionCodec(),
				codec.NewZLibCompressor(),
			)

			mockPriceApplier := aggregatormocks.NewPriceApplier(s.T())

			h := ve.NewVoteExtensionHandler(
				log.NewTestLogger(s.T()),
				tc.oracleService(),
				time.Second*1,
				tc.currencyPairStrategy(),
				cdc,
				mockPriceApplier,
				servicemetrics.NewNopMetrics(),
			)

			req := &cometabci.RequestExtendVote{}
			if tc.extendVoteRequest != nil {
				req = tc.extendVoteRequest()
			}
			if req != nil {
				finalizeBlockReq := &cometabci.RequestFinalizeBlock{
					Txs:    req.Txs,
					Height: req.Height,
				}
				mockPriceApplier.On("ApplyPricesFromVoteExtensions", s.ctx, finalizeBlockReq).Return(nil, nil)
			}

			resp, err := h.ExtendVoteHandler()(s.ctx, req)
			if !tc.expectedError {
				if resp == nil || len(resp.VoteExtension) == 0 {
					return
				}
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				ext, err := cdc.Decode(resp.VoteExtension)
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedResponse.Prices, ext.Prices)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *VoteExtensionTestSuite) TestVerifyVoteExtension() {
	cdc := codec.NewCompressionVoteExtensionCodec(
		codec.NewDefaultVoteExtensionCodec(),
		codec.NewZLibCompressor(),
	)

	cases := []struct {
		name                 string
		getReq               func() *cometabci.RequestVerifyVoteExtension
		currencyPairStrategy func() *mockstrategies.CurrencyPairStrategy
		expectedResponse     *cometabci.ResponseVerifyVoteExtension
		expectedError        bool
	}{
		{
			name: "nil request returns error",
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				return nil
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				return mockstrategies.NewCurrencyPairStrategy(s.T())
			},
			expectedResponse: nil,
			expectedError:    true,
		},
		{
			name: "empty vote extension",
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				return &cometabci.RequestVerifyVoteExtension{}
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				cpStrategy := mockstrategies.NewCurrencyPairStrategy(s.T())
				return cpStrategy
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		{
			name: "empty vote extension - 1 cp in prev state",
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				return &cometabci.RequestVerifyVoteExtension{}
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				cpStrategy := mockstrategies.NewCurrencyPairStrategy(s.T())
				return cpStrategy
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
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				cpStrategy := mockstrategies.NewCurrencyPairStrategy(s.T())
				return cpStrategy
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
		{
			name: "valid vote extension - 2 cp in prev state",
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				prices := map[uint64][]byte{
					0: oneHundred.Bytes(),
					1: twoHundred.Bytes(),
				}

				ext, err := testutils.CreateVoteExtensionBytes(
					prices,
					cdc,
				)
				s.Require().NoError(err)

				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: ext,
					Height:        1,
				}
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				cpStrategy := mockstrategies.NewCurrencyPairStrategy(s.T())
				cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(2), nil).Once()
				return cpStrategy
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		{
			name: "invalid vote extension - 1 cp in prev state - should fail",
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				prices := map[uint64][]byte{
					0: oneHundred.Bytes(),
					1: twoHundred.Bytes(),
				}

				ext, err := testutils.CreateVoteExtensionBytes(
					prices,
					cdc,
				)
				s.Require().NoError(err)

				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: ext,
					Height:        1,
				}
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				cpStrategy := mockstrategies.NewCurrencyPairStrategy(s.T())
				cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Once()
				return cpStrategy
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
		{
			name: "vote extension with no prices",
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				prices := map[uint64][]byte{}

				ext, err := testutils.CreateVoteExtensionBytes(
					prices,
					cdc,
				)
				s.Require().NoError(err)

				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: ext,
					Height:        1,
				}
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				cpStrategy := mockstrategies.NewCurrencyPairStrategy(s.T())
				cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(0), nil).Once()
				return cpStrategy
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

				ext, err := testutils.CreateVoteExtensionBytes(
					prices,
					cdc,
				)
				s.Require().NoError(err)

				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: ext,
					Height:        1,
				}
			},
			currencyPairStrategy: func() *mockstrategies.CurrencyPairStrategy {
				cpStrategy := mockstrategies.NewCurrencyPairStrategy(s.T())
				cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Once()
				return cpStrategy
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			mockPriceApplier := aggregatormocks.NewPriceApplier(s.T())
			handler := ve.NewVoteExtensionHandler(
				log.NewTestLogger(s.T()),
				mocks.NewOracleClient(s.T()),
				time.Second*1,
				tc.currencyPairStrategy(),
				cdc,
				mockPriceApplier,
				servicemetrics.NewNopMetrics(),
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

func (s *VoteExtensionTestSuite) TestExtendVoteLatency() {
	m := metricsmocks.NewMetrics(s.T())
	os := mocks.NewOracleClient(s.T())

	pamock := aggregatormocks.NewPriceApplier(s.T())
	handler := ve.NewVoteExtensionHandler(
		log.NewTestLogger(s.T()),
		os,
		time.Second*1,
		mockstrategies.NewCurrencyPairStrategy(s.T()),
		codec.NewDefaultVoteExtensionCodec(),
		pamock,
		m,
	)

	pamock.On("ApplyPricesFromVoteExtensions", s.ctx, mock.Anything, mock.Anything).Return(nil, nil)
	// mock
	os.On("Prices", mock.Anything, mock.Anything).Return(
		&servicetypes.QueryPricesResponse{
			Prices:    map[string]string{},
			Timestamp: time.Now(),
		},
		nil,
	).Run(func(_ mock.Arguments) {
		// sleep to simulate latency
		time.Sleep(100 * time.Millisecond)
	})

	m.On("ObserveABCIMethodLatency", servicemetrics.ExtendVote, mock.Anything).Run(func(args mock.Arguments) {
		latency := args.Get(1).(time.Duration)
		s.Require().True(latency > 100*time.Millisecond)
	})
	m.On("AddABCIRequest", servicemetrics.ExtendVote, servicemetrics.Success{})
	_, err := handler.ExtendVoteHandler()(s.ctx, &cometabci.RequestExtendVote{
		Height: 1,
		Txs:    [][]byte{},
	})
	s.Require().NoError(err)
}

func (s *VoteExtensionTestSuite) TestExtendVoteStatus() {
	s.Run("test nil request", func() {
		mockMetrics := metricsmocks.NewMetrics(s.T())
		pamock := aggregatormocks.NewPriceApplier(s.T())
		handler := ve.NewVoteExtensionHandler(
			log.NewTestLogger(s.T()),
			nil,
			time.Second*1,
			mockstrategies.NewCurrencyPairStrategy(s.T()),
			nil,
			pamock,
			mockMetrics,
		)

		expErr := connectabci.NilRequestError{
			Handler: servicemetrics.ExtendVote,
		}
		mockMetrics.On("ObserveABCIMethodLatency", servicemetrics.ExtendVote, mock.Anything)
		mockMetrics.On("AddABCIRequest", servicemetrics.ExtendVote, expErr)

		_, err := handler.ExtendVoteHandler()(s.ctx, nil)
		s.Require().NoError(err)
	})

	s.Run("test panic", func() {
		mockMetrics := metricsmocks.NewMetrics(s.T())
		pamock := aggregatormocks.NewPriceApplier(s.T())
		handler := ve.NewVoteExtensionHandler(
			log.NewTestLogger(s.T()),
			nil,
			time.Second*1,
			mockstrategies.NewCurrencyPairStrategy(s.T()),
			nil,
			pamock,
			mockMetrics,
		)

		pamock.On("ApplyPricesFromVoteExtensions", s.ctx, mock.Anything, mock.Anything).Return(nil, nil).Run(func(_ mock.Arguments) {
			panic("panic")
		})

		expErr := ve.ErrPanic{
			Err: fmt.Errorf("panic"),
		}
		mockMetrics.On("ObserveABCIMethodLatency", servicemetrics.ExtendVote, mock.Anything)
		mockMetrics.On("AddABCIRequest", servicemetrics.ExtendVote, expErr)

		_, err := handler.ExtendVoteHandler()(s.ctx, &cometabci.RequestExtendVote{})
		s.Require().Error(err, expErr)
	})

	s.Run("test pre-blocker failure", func() {
		mockMetrics := metricsmocks.NewMetrics(s.T())
		preBlockError := fmt.Errorf("pre-blocker failure")
		pamock := aggregatormocks.NewPriceApplier(s.T())

		handler := ve.NewVoteExtensionHandler(
			log.NewTestLogger(s.T()),
			nil,
			time.Second*1,
			mockstrategies.NewCurrencyPairStrategy(s.T()),
			nil,
			pamock,
			mockMetrics,
		)

		pamock.On("ApplyPricesFromVoteExtensions", s.ctx, mock.Anything, mock.Anything).Return(nil, preBlockError)

		expErr := ve.PreBlockError{
			Err: preBlockError,
		}
		mockMetrics.On("ObserveABCIMethodLatency", servicemetrics.ExtendVote, mock.Anything)
		mockMetrics.On("AddABCIRequest", servicemetrics.ExtendVote, expErr)

		_, err := handler.ExtendVoteHandler()(s.ctx, &cometabci.RequestExtendVote{})
		s.Require().NoError(err)
	})

	s.Run("test client failure", func() {
		mockMetrics := metricsmocks.NewMetrics(s.T())
		clientError := fmt.Errorf("client failure")
		mockClient := mocks.NewOracleClient(s.T())
		pamock := aggregatormocks.NewPriceApplier(s.T())
		handler := ve.NewVoteExtensionHandler(
			log.NewTestLogger(s.T()),
			mockClient,
			time.Second*1,
			mockstrategies.NewCurrencyPairStrategy(s.T()),
			nil,
			pamock,
			mockMetrics,
		)

		pamock.On("ApplyPricesFromVoteExtensions", s.ctx, mock.Anything, mock.Anything).Return(nil, nil)

		expErr := ve.OracleClientError{
			Err: clientError,
		}
		mockMetrics.On("ObserveABCIMethodLatency", servicemetrics.ExtendVote, mock.Anything)
		mockMetrics.On("AddABCIRequest", servicemetrics.ExtendVote, expErr)
		mockClient.On("Prices", mock.Anything, &servicetypes.QueryPricesRequest{}).Return(nil, clientError)

		_, err := handler.ExtendVoteHandler()(s.ctx, &cometabci.RequestExtendVote{})
		s.Require().NoError(err)
	})

	s.Run("test price transformation failures", func() {
		mockMetrics := metricsmocks.NewMetrics(s.T())
		transformationError := fmt.Errorf("incorrectly formatted CurrencyPair: \"BTCETH\"")
		mockClient := mocks.NewOracleClient(s.T())
		pamock := aggregatormocks.NewPriceApplier(s.T())
		handler := ve.NewVoteExtensionHandler(
			log.NewTestLogger(s.T()),
			mockClient,
			time.Second*1,
			mockstrategies.NewCurrencyPairStrategy(s.T()),
			nil,
			pamock,
			mockMetrics,
		)

		pamock.On("ApplyPricesFromVoteExtensions", s.ctx, mock.Anything, mock.Anything).Return(nil, nil)

		expErr := ve.TransformPricesError{
			Err: transformationError,
		}
		mockMetrics.On("ObserveABCIMethodLatency", servicemetrics.ExtendVote, mock.Anything)
		mockMetrics.On("AddABCIRequest", servicemetrics.ExtendVote, expErr)
		mockClient.On("Prices", mock.Anything, &servicetypes.QueryPricesRequest{}).Return(&servicetypes.QueryPricesResponse{
			Prices: map[string]string{
				"BTCETH": "1000",
			},
		}, nil)

		_, err := handler.ExtendVoteHandler()(s.ctx, &cometabci.RequestExtendVote{})
		s.Require().NoError(err)
	})

	s.Run("test codec failures", func() {
		mockMetrics := metricsmocks.NewMetrics(s.T())
		codecError := fmt.Errorf("codec error")
		mockClient := mocks.NewOracleClient(s.T())
		cdc := codecmocks.NewVoteExtensionCodec(s.T())
		pamock := aggregatormocks.NewPriceApplier(s.T())
		handler := ve.NewVoteExtensionHandler(
			log.NewTestLogger(s.T()),
			mockClient,
			time.Second*1,
			mockstrategies.NewCurrencyPairStrategy(s.T()),
			cdc,
			pamock,
			mockMetrics,
		)

		pamock.On("ApplyPricesFromVoteExtensions", s.ctx, mock.Anything, mock.Anything).Return(nil, nil)

		expErr := connectabci.CodecError{
			Err: codecError,
		}
		mockMetrics.On("ObserveABCIMethodLatency", servicemetrics.ExtendVote, mock.Anything)
		mockMetrics.On("AddABCIRequest", servicemetrics.ExtendVote, expErr)
		mockClient.On("Prices", mock.Anything, &servicetypes.QueryPricesRequest{}).Return(&servicetypes.QueryPricesResponse{
			Prices: map[string]string{},
		}, nil)
		cdc.On("Encode", abcitypes.OracleVoteExtension{
			Prices: map[uint64][]byte{},
		}).Return(nil, codecError)

		_, err := handler.ExtendVoteHandler()(s.ctx, &cometabci.RequestExtendVote{})
		s.Require().NoError(err)
	})

	s.Run("test success", func() {
		mockMetrics := metricsmocks.NewMetrics(s.T())
		mockClient := mocks.NewOracleClient(s.T())
		cdc := codecmocks.NewVoteExtensionCodec(s.T())
		pamock := aggregatormocks.NewPriceApplier(s.T())

		handler := ve.NewVoteExtensionHandler(
			log.NewTestLogger(s.T()),
			mockClient,
			time.Second*1,
			mockstrategies.NewCurrencyPairStrategy(s.T()),
			cdc,
			pamock,
			mockMetrics,
		)

		pamock.On("ApplyPricesFromVoteExtensions", s.ctx, mock.Anything, mock.Anything).Return(nil, nil)
		mockMetrics.On("ObserveABCIMethodLatency", servicemetrics.ExtendVote, mock.Anything)
		mockMetrics.On("AddABCIRequest", servicemetrics.ExtendVote, servicemetrics.Success{})
		mockClient.On("Prices", mock.Anything, &servicetypes.QueryPricesRequest{}).Return(&servicetypes.QueryPricesResponse{
			Prices: map[string]string{},
		}, nil)
		cdc.On("Encode", abcitypes.OracleVoteExtension{
			Prices: map[uint64][]byte{},
		}).Return(nil, nil)

		_, err := handler.ExtendVoteHandler()(s.ctx, &cometabci.RequestExtendVote{})
		s.Require().NoError(err)
	})
}

func (s *VoteExtensionTestSuite) TestVerifyVoteExtensionStatus() {
	s.Run("nil request", func() {
		mockMetrics := metricsmocks.NewMetrics(s.T())
		handler := ve.NewVoteExtensionHandler(
			log.NewTestLogger(s.T()),
			nil,
			time.Second*1,
			nil,
			nil,
			aggregatormocks.NewPriceApplier(s.T()),
			mockMetrics,
		)
		expErr := connectabci.NilRequestError{
			Handler: servicemetrics.VerifyVoteExtension,
		}
		mockMetrics.On("ObserveABCIMethodLatency", servicemetrics.VerifyVoteExtension, mock.Anything)
		mockMetrics.On("AddABCIRequest", servicemetrics.VerifyVoteExtension, expErr)

		_, err := handler.VerifyVoteExtensionHandler()(s.ctx, nil)
		s.Require().Error(err, expErr)
	})

	s.Run("codec error", func() {
		mockMetrics := metricsmocks.NewMetrics(s.T())
		codecError := fmt.Errorf("codec error")
		cdc := codecmocks.NewVoteExtensionCodec(s.T())
		handler := ve.NewVoteExtensionHandler(
			log.NewTestLogger(s.T()),
			nil,
			time.Second*1,
			nil,
			cdc,
			aggregatormocks.NewPriceApplier(s.T()),
			mockMetrics,
		)
		expErr := connectabci.CodecError{
			Err: codecError,
		}
		mockMetrics.On("ObserveABCIMethodLatency", servicemetrics.VerifyVoteExtension, mock.Anything)
		mockMetrics.On("AddABCIRequest", servicemetrics.VerifyVoteExtension, expErr)
		cdc.On("Decode", mock.Anything).Return(abcitypes.OracleVoteExtension{}, codecError)

		_, err := handler.VerifyVoteExtensionHandler()(s.ctx, &cometabci.RequestVerifyVoteExtension{
			VoteExtension: []byte{1, 2, 3},
		})
		s.Require().Error(err, expErr)
	})

	s.Run("invalid vote-extension", func() {
		mockMetrics := metricsmocks.NewMetrics(s.T())
		cdc := codecmocks.NewVoteExtensionCodec(s.T())
		cpStrategy := mockstrategies.NewCurrencyPairStrategy(s.T())
		cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Once()

		length := 34
		transformErr := fmt.Errorf("price bytes are too long: %d", length)

		handler := ve.NewVoteExtensionHandler(
			log.NewTestLogger(s.T()),
			nil,
			time.Second*1,
			cpStrategy,
			cdc,
			aggregatormocks.NewPriceApplier(s.T()),
			mockMetrics,
		)
		expErr := ve.ValidateVoteExtensionError{
			Err: transformErr,
		}
		mockMetrics.On("ObserveABCIMethodLatency", servicemetrics.VerifyVoteExtension, mock.Anything)
		mockMetrics.On("AddABCIRequest", servicemetrics.VerifyVoteExtension, expErr)
		cdc.On("Decode", mock.Anything).Return(abcitypes.OracleVoteExtension{
			Prices: map[uint64][]byte{
				1: make([]byte, length),
			},
		}, nil)

		_, err := handler.VerifyVoteExtensionHandler()(s.ctx, &cometabci.RequestVerifyVoteExtension{
			VoteExtension: []byte{1, 2, 3},
		})
		s.Require().Error(err, expErr)
	})

	s.Run("success", func() {
		mockMetrics := metricsmocks.NewMetrics(s.T())
		cdc := codecmocks.NewVoteExtensionCodec(s.T())
		cpStrategy := mockstrategies.NewCurrencyPairStrategy(s.T())
		cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Once()

		handler := ve.NewVoteExtensionHandler(
			log.NewTestLogger(s.T()),
			nil,
			time.Second*1,
			cpStrategy,
			cdc,
			aggregatormocks.NewPriceApplier(s.T()),
			mockMetrics,
		)

		mockMetrics.On("ObserveABCIMethodLatency", servicemetrics.VerifyVoteExtension, mock.Anything)
		mockMetrics.On("AddABCIRequest", servicemetrics.VerifyVoteExtension, servicemetrics.Success{})
		mockMetrics.On("ObserveMessageSize", servicemetrics.VoteExtension, mock.Anything)

		cdc.On("Decode", mock.Anything).Return(abcitypes.OracleVoteExtension{}, nil)

		_, err := handler.VerifyVoteExtensionHandler()(s.ctx, &cometabci.RequestVerifyVoteExtension{
			VoteExtension: []byte{1, 2, 3},
		})
		s.Require().NoError(err)
	})
}

func (s *VoteExtensionTestSuite) TestVoteExtensionSize() {
	mockMetrics := metricsmocks.NewMetrics(s.T())

	mockClient := mocks.NewOracleClient(s.T())
	cdc := codecmocks.NewVoteExtensionCodec(s.T())
	cpStrategy := mockstrategies.NewCurrencyPairStrategy(s.T())
	cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Once()
	pamock := aggregatormocks.NewPriceApplier(s.T())

	handler := ve.NewVoteExtensionHandler(
		log.NewTestLogger(s.T()),
		mockClient,
		time.Second*1,
		cpStrategy,
		cdc,
		pamock,
		mockMetrics,
	)

	voteExtension := make([]byte, 100)

	// mock metrics calls
	mockMetrics.On("ObserveABCIMethodLatency", servicemetrics.VerifyVoteExtension, mock.Anything)
	mockMetrics.On("AddABCIRequest", servicemetrics.VerifyVoteExtension, servicemetrics.Success{})
	mockMetrics.On("ObserveMessageSize", servicemetrics.VoteExtension, 100)

	// mock codec calls
	cdc.On("Decode", mock.Anything).Return(abcitypes.OracleVoteExtension{
		Prices: map[uint64][]byte{},
	}, nil)

	_, err := handler.VerifyVoteExtensionHandler()(s.ctx, &cometabci.RequestVerifyVoteExtension{
		VoteExtension: voteExtension,
	})
	s.Require().NoError(err)
}
