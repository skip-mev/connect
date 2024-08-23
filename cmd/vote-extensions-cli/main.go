package main

import (
	"fmt"
	"math/big"
	"os"

	cmtabci "github.com/cometbft/cometbft/abci/types"
	cmthttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/spf13/cobra"

	"github.com/skip-mev/connect/v2/abci/strategies/codec"
)

var (
	rootCmd = &cobra.Command{
		Use:   "vote-extensions-cli",
		Short: "Inspect the vote extensions of a given node, at a given height",
		Long: `Use as follows to inspect the vote extensions of a given node, at a given height:
		
		vote-extensions-cli --node <http<s>://<url>:26657> --height <height> --extended-commit-codec <selector> --vote-extension-codec <selector>
		Where:
			--node: The node to query
			--height: The height to query. If not provided, the latest height will be used
			--extended-commit-codec: The codec to use to decode the extended commit. Options are 1: standard encoding (default), 2: z-lib compressed encoding, 3: zstd compressed encoding
			--vote-extension-codec: The codec to use to decode the vote extension. Options are 1: standard encoding (default), 2: z-lib compressed encoding, 3: zstd compressed encoding
		`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// create a comet-http client
			client, err := cmthttp.New(node, "/websocket")
			if err != nil {
				return err
			}

			// get the block-data
			height := &height
			if *height == 0 {
				height = nil
			}

			block, err := client.Block(cmd.Context(), height)
			if err != nil {
				return err
			}

			// decode the extended commit
			var extCommit cmtabci.ExtendedCommitInfo
			extCommitCodec, veCodec := codecsFromFlags(extendedCommitCodec, voteExtensionCodec)
			extCommit, err = extCommitCodec.Decode(block.Block.Txs[0])
			if err != nil {
				return fmt.Errorf("failed to decode extended commit: %w", err)
			}

			// log the round of this extended commit
			cmd.Println("Height:", block.Block.Height, "Round:", extCommit.Round)

			// decode the vote extensions
			for _, vote := range extCommit.Votes {
				// decode the vote extension
				ve, err := veCodec.Decode(vote.VoteExtension)
				if err != nil {
					return err
				}

				// per price, unmarshal via gob decoding
				for priceID, priceBz := range ve.Prices {
					cmd.Println("Price ID:", priceID)

					price := new(big.Int)
					if err := price.GobDecode(priceBz); err != nil {
						return err
					}
					cmd.Println("Price:", price)
				}

				// log the block id
				cmd.Println("Block ID:", vote.BlockIdFlag)

				// log the validator
				cmd.Println("Validator:", vote.Validator)
			}

			return nil
		},
	}

	// Flags.
	node                string
	height              int64
	extendedCommitCodec string
	voteExtensionCodec  string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&node, "node", "", "The node to query")
	rootCmd.PersistentFlags().Int64Var(&height, "height", 0, "The height to query. If not provided, the latest height will be used")
	rootCmd.PersistentFlags().StringVar(&extendedCommitCodec, "extended-commit-codec", "1", "The codec to use to decode the extended commit. Options are 1: standard encoding (default), 2: z-lib compressed encoding, 3: zstd compressed encoding")
	rootCmd.PersistentFlags().StringVar(&voteExtensionCodec, "vote-extension-codec", "1", "The codec to use to decode the vote extension. Options are 1: standard encoding (default), 2: z-lib compressed encoding, 3: zstd compressed encoding")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func codecsFromFlags(extCommitCodecFlag, veCodecFlag string) (codec.ExtendedCommitCodec, codec.VoteExtensionCodec) {
	var extCommitCodec codec.ExtendedCommitCodec
	var veCodec codec.VoteExtensionCodec

	switch extCommitCodecFlag {
	case "1":
		extCommitCodec = codec.NewDefaultExtendedCommitCodec()
	case "2":
		extCommitCodec = codec.NewCompressionExtendedCommitCodec(
			codec.NewDefaultExtendedCommitCodec(),
			codec.NewZLibCompressor(),
		)
	case "3":
		extCommitCodec = codec.NewCompressionExtendedCommitCodec(
			codec.NewDefaultExtendedCommitCodec(),
			codec.NewZStdCompressor(),
		)
	}

	switch veCodecFlag {
	case "1":
		veCodec = codec.NewDefaultVoteExtensionCodec()
	case "2":
		veCodec = codec.NewCompressionVoteExtensionCodec(
			codec.NewDefaultVoteExtensionCodec(),
			codec.NewZLibCompressor(),
		)
	case "3":
		veCodec = codec.NewCompressionVoteExtensionCodec(
			codec.NewDefaultVoteExtensionCodec(),
			codec.NewZStdCompressor(),
		)
	}

	return extCommitCodec, veCodec
}
