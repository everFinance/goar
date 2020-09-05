package utils

import (
	"crypto/sha256"
	"math/big"
)

const (
	MAX_CHUNK_SIZE = 256 * 1024
	MIN_CHUNK_SIZE = 32 * 1024

	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	wordBytes = wordBits / 8
)

type Chunks struct {
	DataRoot []byte
	Chunks   []Chunk
	Proofs   []Proof
}

type Chunk struct {
	DataHash     []byte
	MinByteRange int
	MaxByteRange int
}

type LeafNode struct {
	ID           []byte
	DataHash     []byte
	Type         string
	MinByteRange int
	MaxByteRange int
}

type Proof struct {
	Offest int
	Proof  []byte
}

func GenerateChunks(data []byte) (c Chunks) {
	chunks := chunkData(data)
	leaves := generateLeaves(chunks)
	// TODO, calculate root & proofs

	c.DataRoot = leaves[0].ID
	return
}

func chunkData(data []byte) (chunks []Chunk) {
	cursor := 0
	// TODO
	// for len(data) >= MAX_CHUNK_SIZE {
	// }

	hash := sha256.Sum256(data)
	chunks = append(chunks, Chunk{
		DataHash:     hash[:],
		MinByteRange: cursor,
		MaxByteRange: cursor + len(data),
	})
	return
}

func generateLeaves(chunks []Chunk) (leafs []LeafNode) {
	for _, chunk := range chunks {
		hDataHash := sha256.Sum256(chunk.DataHash)
		hMaxByteRange := sha256.Sum256(PaddedBigBytes(big.NewInt(int64(chunk.MaxByteRange)), 32))

		leafs = append(leafs, LeafNode{
			ID: hashArray(
				[][]byte{hDataHash[:], hMaxByteRange[:]},
			),
			DataHash:     chunk.DataHash,
			Type:         "leaf",
			MinByteRange: chunk.MinByteRange,
			MaxByteRange: chunk.MaxByteRange,
		})
	}

	return
}

func buildLayer() {
	//TODO
}

func generateProofs() {
	// TODO
}

func hashArray(data [][]byte) []byte {
	dataFlat := []byte{}
	for _, v := range data {
		dataFlat = append(dataFlat, v...)
	}

	hash := sha256.Sum256(dataFlat)
	return hash[:]
}

func PaddedBigBytes(bigint *big.Int, n int) []byte {
	if bigint.BitLen()/8 >= n {
		return bigint.Bytes()
	}
	ret := make([]byte, n)
	ReadBits(bigint, ret)
	return ret
}

func ReadBits(bigint *big.Int, buf []byte) {
	i := len(buf)
	for _, d := range bigint.Bits() {
		for j := 0; j < wordBytes && i > 0; j++ {
			i--
			buf[i] = byte(d)
			d >>= 8
		}
	}
}
