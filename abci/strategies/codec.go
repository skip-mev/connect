package strategies

import (
	"bytes"
	"compress/zlib"
	"io"

	vetypes "github.com/skip-mev/slinky/abci/ve/types"
)

// VoteExtensionCodec is the interface for encoding / decoding vote extensions.
type VoteExtensionCodec interface {
	// Encode encodes the vote extension into a byte array.
	Encode(ve vetypes.OracleVoteExtension) ([]byte, error)

	// Decode decodes the vote extension from a byte array.
	Decode([]byte) (vetypes.OracleVoteExtension, error)
}

// NewDefaultVoteExtensionCodec returns a new DefaultVoteExtensionCodec.
func NewDefaultVoteExtensionCodec() *DefaultVoteExtensionCodec {
	return &DefaultVoteExtensionCodec{}
}

// DefaultVoteExtensionCodec is the default implementation of VoteExtensionCodec. It uses the
// vanilla implementations of Unmarshal / Marshal under the hood
type DefaultVoteExtensionCodec struct{}

func (codec *DefaultVoteExtensionCodec) Encode(ve vetypes.OracleVoteExtension) ([]byte, error) {
	return ve.Marshal()
}

func (codec *DefaultVoteExtensionCodec) Decode(bz []byte) (vetypes.OracleVoteExtension, error) {
	if len(bz) == 0 {
		return vetypes.OracleVoteExtension{}, nil
	}

	var ve vetypes.OracleVoteExtension
	return ve, ve.Unmarshal(bz)
}

type Compressor interface {
	Compress([]byte) ([]byte, error)
	Decompress([]byte) ([]byte, error)
}

// ZLibCompressor is a Compressor that uses zlib to compress / decompress byte arrays, this object is not thread-safe.
type ZLibCompressor struct{}

// NewZLibCompressor returns a new zlibDecompressor.
func NewZLibCompressor() *ZLibCompressor {
	return &ZLibCompressor{}
}

// Compress compresses the given byte array using zlib. It returns an error if the compression fails.
// This function is not thread-safe, and uses zlib.BestCompression as the compression level.
func (c *ZLibCompressor) Compress(bz []byte) ([]byte, error) {
	var b bytes.Buffer

	// we use the best compression level as size reduction is prioritized
	w, err := zlib.NewWriterLevel(&b, zlib.BestCompression)
	if err != nil {
		return nil, err
	}
	defer w.Close()

	// write and flush the buffer
	if _, err := w.Write(bz); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Decompress decompresses the given byte array using zlib. It returns an error if the decompression fails.
func (c *ZLibCompressor) Decompress(bz []byte) ([]byte, error) {
	if len(bz) == 0 {
		return nil, nil
	}
	r, err := zlib.NewReader(bytes.NewReader(bz))
	if err != nil {
		return nil, err
	}
	r.Close()

	// read bytes and return
	return io.ReadAll(r)
}

// CompressionVoteExtensionCodec is a VoteExtensionCodec that uses compression to encode / decode, given a
// VoteExtensionCodec that can encode / decode uncompressed vote extensions.
type CompressionVoteExtensionCodec struct {
	codec      VoteExtensionCodec
	compressor Compressor
}

// NewCompressionVoteExtensionCodec returns a new CompressionVoteExtensionCodec given an underlying codec.
func NewCompressionVoteExtensionCodec(codec VoteExtensionCodec, compressor Compressor) *CompressionVoteExtensionCodec {
	return &CompressionVoteExtensionCodec{
		codec:      codec,
		compressor: compressor,
	}
}

// Encode returns the encoded vote extension using the codec's Encode method and then compresses the result.
// This implementation uses zstd compression.
func (codec *CompressionVoteExtensionCodec) Encode(ve vetypes.OracleVoteExtension) ([]byte, error) {
	bz, err := codec.codec.Encode(ve)
	if err != nil {
		return nil, err
	}

	return codec.compressor.Compress(bz)
}

// Decode decompresses the vote extension using zstd and then decodes the result using the codec's Decode method.
func (codec *CompressionVoteExtensionCodec) Decode(bz []byte) (vetypes.OracleVoteExtension, error) {
	// decompress first
	bz, err := codec.compressor.Decompress(bz)
	if err != nil {
		return vetypes.OracleVoteExtension{}, err
	}

	return codec.codec.Decode(bz)
}
