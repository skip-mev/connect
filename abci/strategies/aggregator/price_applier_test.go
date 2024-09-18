package aggregator_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/abci/strategies/aggregator"
	"github.com/skip-mev/connect/v2/abci/strategies/aggregator/mocks"
	"github.com/skip-mev/connect/v2/abci/strategies/codec"
	"github.com/skip-mev/connect/v2/abci/testutils"
	abcimocks "github.com/skip-mev/connect/v2/abci/types/mocks"

	"cosmossdk.io/log"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	vetypes "github.com/skip-mev/connect/v2/abci/ve/types"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
)

func TestPriceApplier(t *testing.T) {
	veCodec := codec.NewDefaultVoteExtensionCodec()
	extCommitcodec := codec.NewDefaultExtendedCommitCodec()

	va := mocks.NewVoteAggregator(t)

	ok := abcimocks.NewOracleKeeper(t)

	pa := aggregator.NewOraclePriceApplier(
		va,
		ok,
		veCodec,
		extCommitcodec,
		log.NewNopLogger(),
	)

	t.Run("if extracting oracle votes fails, fail", func(t *testing.T) {
		ctx := sdk.Context{}

		// first tx is garbage, fail
		prices, err := pa.ApplyPricesFromVoteExtensions(ctx, &abcitypes.RequestFinalizeBlock{
			Txs: [][]byte{[]byte("garbage")},
		})

		require.Error(t, err)
		require.Nil(t, prices)
	})

	t.Run("if vote aggregation fails, fail", func(t *testing.T) {
		prices := map[uint64][]byte{
			1: []byte("price1"),
		}
		ca := sdk.ConsAddress("val1")

		vote1, err := testutils.CreateExtendedVoteInfo(
			ca,
			prices,
			veCodec,
		)
		require.NoError(t, err)

		_, extCommitInfoBz, err := testutils.CreateExtendedCommitInfo(
			[]abcitypes.ExtendedVoteInfo{vote1},
			extCommitcodec,
		)
		require.NoError(t, err)

		ctx := sdk.Context{}

		// fail vote aggregation
		va.On("AggregateOracleVotes", ctx, []aggregator.Vote{
			{
				OracleVoteExtension: vetypes.OracleVoteExtension{
					Prices: prices,
				},
				ConsAddress: ca,
			},
		}).Return(nil, fmt.Errorf("fail")).Once()

		returnedPrices, err := pa.ApplyPricesFromVoteExtensions(ctx, &abcitypes.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz},
		})

		require.Error(t, err)
		require.Nil(t, returnedPrices)
	})

	t.Run("ignore negative prices", func(t *testing.T) {
		priceBz := big.NewInt(-100).Bytes()

		prices := map[uint64][]byte{
			1: priceBz,
		}

		ca := sdk.ConsAddress("val1")

		vote1, err := testutils.CreateExtendedVoteInfo(
			ca,
			prices,
			veCodec,
		)
		require.NoError(t, err)

		_, extCommitInfoBz, err := testutils.CreateExtendedCommitInfo(
			[]abcitypes.ExtendedVoteInfo{vote1},
			extCommitcodec,
		)
		require.NoError(t, err)

		ctx := sdk.Context{}

		// succeed vote aggregation
		cp := connecttypes.NewCurrencyPair("BTC", "USD")
		va.On("AggregateOracleVotes", ctx, []aggregator.Vote{
			{
				OracleVoteExtension: vetypes.OracleVoteExtension{
					Prices: prices,
				},
				ConsAddress: ca,
			},
		}).Return(map[connecttypes.CurrencyPair]*big.Int{
			cp: big.NewInt(-100),
		}, nil)

		ok.On("GetAllCurrencyPairs", ctx).Return(
			[]connecttypes.CurrencyPair{cp},
		)

		_, err = pa.ApplyPricesFromVoteExtensions(ctx, &abcitypes.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz},
		})

		require.NoError(t, err)
	})

	t.Run("update prices in state", func(t *testing.T) {
		priceBz := big.NewInt(100).Bytes()

		prices1 := map[uint64][]byte{
			1: priceBz,
		}

		prices2 := map[uint64][]byte{
			1: big.NewInt(200).Bytes(),
		}

		ca1 := sdk.ConsAddress("val1")
		ca2 := sdk.ConsAddress("val2")

		vote1, err := testutils.CreateExtendedVoteInfo(
			ca1,
			prices1,
			veCodec,
		)
		require.NoError(t, err)

		vote2, err := testutils.CreateExtendedVoteInfo(
			ca2,
			prices2,
			veCodec,
		)
		require.NoError(t, err)

		_, extCommitInfoBz, err := testutils.CreateExtendedCommitInfo(
			[]abcitypes.ExtendedVoteInfo{vote1, vote2},
			extCommitcodec,
		)
		require.NoError(t, err)

		ctx := sdk.Context{}.WithBlockHeader(cmtproto.Header{
			Time: time.Now(),
		}).WithBlockHeight(1)

		// succeed vote aggregation
		cp := connecttypes.NewCurrencyPair("BTC", "USD")

		va.On("AggregateOracleVotes", ctx, []aggregator.Vote{
			{
				OracleVoteExtension: vetypes.OracleVoteExtension{
					Prices: prices1,
				},
				ConsAddress: ca1,
			},
			{
				OracleVoteExtension: vetypes.OracleVoteExtension{
					Prices: prices2,
				},
				ConsAddress: ca2,
			},
		}).Return(map[connecttypes.CurrencyPair]*big.Int{
			cp: big.NewInt(150),
		}, nil)

		// return multiple prices
		ok.On("GetAllCurrencyPairs", ctx).Return(
			[]connecttypes.CurrencyPair{cp, connecttypes.NewCurrencyPair("ETH", "USD")}, // ignore last cp
		)

		ok.On("SetPriceForCurrencyPair", ctx, cp, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			qp := args.Get(2).(oracletypes.QuotePrice)

			require.Equal(t, qp.Price.BigInt(), big.NewInt(150))
			require.Equal(t, qp.BlockTimestamp, ctx.BlockHeader().Time)
			require.Equal(t, qp.BlockHeight, uint64(ctx.BlockHeight())) //nolint:gosec
		})

		prices, err := pa.ApplyPricesFromVoteExtensions(ctx, &abcitypes.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz},
		})
		require.NoError(t, err)
		require.Equal(t, map[connecttypes.CurrencyPair]*big.Int{
			cp: big.NewInt(150),
		}, prices)

		// get prices from validators
		expPrices := map[connecttypes.CurrencyPair]*big.Int{
			cp: big.NewInt(150),
		}
		va.On("GetPriceForValidator", ca1).Return(
			expPrices,
		)
		valPrices := pa.GetPricesForValidator(ca1)
		require.Equal(t, expPrices, valPrices)
	})
}
