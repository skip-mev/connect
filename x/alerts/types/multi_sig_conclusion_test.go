package types_test

import (
	"math/big"
	"testing"

	cmtabci "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256r1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/alerts/types"
)

type MultiSigConclusionTestSuite struct {
	// test-suite
	suite.Suite

	pks  []cryptotypes.PubKey
	pvks []cryptotypes.PrivKey
}

func TestMultiSigConclusionTestSuite(t *testing.T) {
	suite.Run(t, new(MultiSigConclusionTestSuite))
}

func (s *MultiSigConclusionTestSuite) SetupTest() {
	// create the set of public / private keys
	s.pks = make([]cryptotypes.PubKey, 3)
	s.pvks = make([]cryptotypes.PrivKey, 3)

	// create a secp256k1, ed25519, and secp256r1 key
	s.pvks[0] = secp256k1.GenPrivKey()
	s.pks[0] = s.pvks[0].PubKey()

	s.pvks[1] = ed25519.GenPrivKey()
	s.pks[1] = s.pvks[1].PubKey()

	var err error
	s.pvks[2], err = secp256r1.GenPrivKey()
	s.Require().NoError(err)
	s.pks[2] = s.pvks[2].PubKey()
}

type invalidPubKey struct {
	cryptotypes.PubKey
}

func (i invalidPubKey) Type() string {
	return "invalid"
}

func (s *MultiSigConclusionTestSuite) TestParams() {
	s.Run("test NewMultiSigVerificationParams()", func() {
		s.Run("test duplicate pub-keys fails", func() {
			_, err := types.NewMultiSigVerificationParams([]cryptotypes.PubKey{s.pks[0], s.pks[0]})
			s.Require().Error(err)
		})

		s.Run("test 0 pubkeys fails", func() {
			_, err := types.NewMultiSigVerificationParams([]cryptotypes.PubKey{})
			s.Require().Error(err)
		})

		s.Run("test all 3 pubkeys passes validate basic", func() {
			params, err := types.NewMultiSigVerificationParams(s.pks)
			s.Require().NoError(err)

			// passes validate-basic
			s.Require().NoError(params.ValidateBasic())
		})
	})

	s.Run("test ValidateBasic", func() {
		s.Run("signers map with invalid pubkey", func() {
			params := types.MultiSigConclusionVerificationParams{
				Signers: []*codectypes.Any{
					nil,
				},
			}

			s.Require().Error(params.ValidateBasic())
		})

		s.Run("empty signers array", func() {
			params := types.MultiSigConclusionVerificationParams{
				Signers: []*codectypes.Any{},
			}

			s.Require().Error(params.ValidateBasic())
		})

		s.Run("invalid public key type in array", func() {
			pkAny, err := codectypes.NewAnyWithValue(&invalidPubKey{
				PubKey: s.pks[0],
			})
			s.Require().NoError(err)

			params := types.MultiSigConclusionVerificationParams{
				Signers: []*codectypes.Any{pkAny},
			}

			s.Require().Error(params.ValidateBasic())
		})

		a, err := codectypes.NewAnyWithValue(s.pks[0])
		s.Require().NoError(err)

		s.Run("duplicate public-keys", func() {
			params := types.MultiSigConclusionVerificationParams{
				Signers: []*codectypes.Any{a, a},
			}

			s.Require().Error(params.ValidateBasic())
		})

		s.Run("public keys are all valid", func() {
			params := types.MultiSigConclusionVerificationParams{
				Signers: []*codectypes.Any{a},
			}

			s.Require().NoError(params.ValidateBasic())
		})
	})
}

func (s *MultiSigConclusionTestSuite) TestConclusion() {
	low := big.NewInt(1)
	high := big.NewInt(2)
	invalidPriceBound := types.PriceBound{
		High: low.String(),
		Low:  high.String(),
	}
	priceBound := types.PriceBound{
		Low:  low.String(),
		High: high.String(),
	}
	s.Run("test ValidateBasic()", func() {
		alert := types.NewAlert(1, sdk.AccAddress("abc"), slinkytypes.NewCurrencyPair("A", "B"))

		cases := []struct {
			name       string
			conclusion types.MultiSigConclusion
			valid      bool
		}{
			{
				"invalid alert - fail",
				types.MultiSigConclusion{
					ExtendedCommitInfo: cmtabci.ExtendedCommitInfo{},
					Alert:              types.Alert{},
					PriceBound:         types.PriceBound{},
				},
				false,
			},
			{
				"invalid price-bound - fail",
				types.MultiSigConclusion{
					ExtendedCommitInfo: cmtabci.ExtendedCommitInfo{},
					Alert:              alert,
					PriceBound:         invalidPriceBound,
				},
				false,
			},
			{
				"invalid signer address",
				types.MultiSigConclusion{
					ExtendedCommitInfo: cmtabci.ExtendedCommitInfo{},
					Alert:              alert,
					PriceBound:         types.PriceBound{},
					Signatures: []types.Signature{
						{
							"invalid",
							nil,
						},
					},
				},
				false,
			},
			{
				"no signers",
				types.MultiSigConclusion{
					ExtendedCommitInfo: cmtabci.ExtendedCommitInfo{},
					Alert:              alert,
					PriceBound:         priceBound,
					Signatures:         nil,
				},
				false,
			},
			{
				"valid conclusion",
				types.MultiSigConclusion{
					ExtendedCommitInfo: cmtabci.ExtendedCommitInfo{},
					Alert:              alert,
					PriceBound:         priceBound,
					Signatures: []types.Signature{
						{
							sdk.AccAddress(s.pks[0].Address()).String(),
							nil,
						},
					},
				},
				true,
			},
		}

		for _, tc := range cases {
			s.Run(tc.name, func() {
				err := tc.conclusion.ValidateBasic()
				if tc.valid {
					s.Require().NoError(err)
				} else {
					s.Require().Error(err)
				}
			})
		}
	})

	s.Run("test Verify()", func() {
		params, err := types.NewMultiSigVerificationParams(s.pks)
		s.Require().NoError(err)
		alert := types.NewAlert(1, sdk.AccAddress("abc"), slinkytypes.NewCurrencyPair("A", "B"))

		s.Run("invalid params - fail", func() {
			msc := types.MultiSigConclusion{}

			s.Require().Error(msc.Verify(&types.MultiSigConclusionVerificationParams{
				Signers: []*codectypes.Any{nil},
			}))
		})

		s.Run("invalid conclusion - fail", func() {
			msc := types.MultiSigConclusion{
				ExtendedCommitInfo: cmtabci.ExtendedCommitInfo{},
				Alert:              types.Alert{},
				PriceBound:         priceBound,
			}

			s.Require().Error(msc.Verify(params))
		})

		s.Run("invalid signature - fail", func() {
			msc := types.MultiSigConclusion{
				ExtendedCommitInfo: cmtabci.ExtendedCommitInfo{},
				Alert:              alert,
				PriceBound:         priceBound,
				Signatures: []types.Signature{
					{
						s.pks[0].Address().String(),
						nil,
					},
				},
			}

			s.Require().Error(msc.Verify(params))
		})

		s.Run("valid conclusion - success", func() {
			msc := types.MultiSigConclusion{
				ExtendedCommitInfo: cmtabci.ExtendedCommitInfo{},
				Alert:              alert,
				PriceBound:         priceBound,
			}

			signBytes, err := msc.SignBytes()
			s.Require().NoError(err)

			sig0, err := s.pvks[0].Sign(signBytes)
			s.Require().NoError(err)

			sig1, err := s.pvks[1].Sign(signBytes)
			s.Require().NoError(err)

			sig2, err := s.pvks[2].Sign(signBytes)
			s.Require().NoError(err)

			msc.Signatures = []types.Signature{
				{
					sdk.AccAddress(s.pks[0].Address()).String(),
					sig0,
				},
				{
					sdk.AccAddress(s.pks[1].Address()).String(),
					sig1,
				},
				{
					sdk.AccAddress(s.pks[2].Address()).String(),
					sig2,
				},
			}

			s.Require().NoError(msc.Verify(params))
		})
	})
}
