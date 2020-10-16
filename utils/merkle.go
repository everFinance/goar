package utils

import (
	"crypto/sha256"
	"math"
	"math/big"

	"github.com/shopspring/decimal"
)

const (
	MAX_CHUNK_SIZE = 256 * 1024
	MIN_CHUNK_SIZE = 32 * 1024

	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	wordBytes      = wordBits / 8
	BranchNodeType = "branch"
	LeafNodeType   = "leaf"
)

type Chunks struct {
	DataRoot []byte
	Chunks   []Chunk
	Proofs   []*Proof
}

type Chunk struct {
	DataHash     []byte
	MinByteRange int
	MaxByteRange int
}

// Node include leaf node and branch node
type Node struct {
	ID           []byte
	Type         string // "branch" or "leaf"
	DataHash     []byte // only leaf node
	MinByteRange int    // only leaf node
	MaxByteRange int
	ByteRange    int   // only branch node
	LeftChild    *Node // only branch node
	RightChild   *Node // only branch node
}

type Proof struct {
	Offest int
	Proof  []byte
}

func GenerateChunks(data []byte) Chunks {
	chunks := chunkData(data)
	leaves := generateLeaves(chunks)
	root := buildLayer(leaves, 0) // leaf node level == 0
	proofs := generateProofs(root)

	// Discard the last chunk & proof if it's zero length.
	lastChunk := chunks[len(chunks)-1]
	if lastChunk.MaxByteRange-lastChunk.MinByteRange == 0 {
		chunks = chunks[:len(chunks)-1]
		proofs = proofs[:len(proofs)-1]
	}

	return Chunks{
		DataRoot: root.ID,
		Chunks:   chunks,
		Proofs:   proofs,
	}
}

func chunkData(data []byte) (chunks []Chunk) {
	cursor := 0
	var rest = data
	// if data length > max size
	for len(rest) >= MAX_CHUNK_SIZE {
		chunkSize := MAX_CHUNK_SIZE

		// 查看下一轮的chunkSize 是否小于最小的size，如果是则在这轮中调整chunk size 的大小
		nextChunkSize := len(rest) - MAX_CHUNK_SIZE
		if nextChunkSize > 0 && nextChunkSize < MIN_CHUNK_SIZE {
			dec := decimal.NewFromFloat(math.Ceil(float64(len(rest) / 2)))
			chunkSize = int(dec.IntPart())
		}

		chunk := rest[:chunkSize]
		dataHash := sha256.Sum256(chunk)
		cursor += len(chunk)
		chunks = append(chunks, Chunk{
			DataHash:     dataHash[:],
			MinByteRange: cursor - len(chunk),
			MaxByteRange: cursor,
		})

		rest = rest[chunkSize:]
	}

	hash := sha256.Sum256(rest)
	chunks = append(chunks, Chunk{
		DataHash:     hash[:],
		MinByteRange: cursor,
		MaxByteRange: cursor + len(rest),
	})
	return
}

func generateLeaves(chunks []Chunk) (leafs []*Node) {
	for _, chunk := range chunks {
		hDataHash := sha256.Sum256(chunk.DataHash)
		hMaxByteRange := sha256.Sum256(PaddedBigBytes(big.NewInt(int64(chunk.MaxByteRange)), 32))

		leafs = append(leafs, &Node{
			ID: hashArray(
				[][]byte{hDataHash[:], hMaxByteRange[:]},
			),
			Type:         LeafNodeType,
			DataHash:     chunk.DataHash,
			MinByteRange: chunk.MinByteRange,
			MaxByteRange: chunk.MaxByteRange,
			ByteRange:    0,
			LeftChild:    nil,
			RightChild:   nil,
		})
	}
	return
}

// buildLayer
func buildLayer(nodes []*Node, level int) (root *Node) {
	if len(nodes) == 1 {
		root = nodes[0]
		return
	}

	nextLayer := make([]*Node, 0, len(nodes)/2)
	for i := 0; i < len(nodes); i += 2 {
		leftNode := nodes[i]
		var rightNode *Node
		if i+1 < len(nodes) {
			rightNode = nodes[i+1]
		}
		nextLayer = append(nextLayer, hashBranch(leftNode, rightNode))
	}

	return buildLayer(nextLayer, level+1)
}

// hashBranch get branch node by child node
func hashBranch(leftNode, rightNode *Node) (branchNode *Node) {
	// 如果只有一个node，则该node 为branch node
	if rightNode == nil {
		return leftNode
	}
	hLeafNodeId := sha256.Sum256(leftNode.ID)
	hRightNodeId := sha256.Sum256(rightNode.ID)
	hLeafNodeMaxRange := sha256.Sum256(PaddedBigBytes(big.NewInt(int64(leftNode.MaxByteRange)), 32))
	id := hashArray([][]byte{hLeafNodeId[:], hRightNodeId[:], hLeafNodeMaxRange[:]})

	branchNode = &Node{
		Type:         BranchNodeType,
		ID:           id,
		DataHash:     nil,
		MinByteRange: 0,
		MaxByteRange: rightNode.MaxByteRange,
		ByteRange:    leftNode.MaxByteRange,
		LeftChild:    leftNode,
		RightChild:   rightNode,
	}

	return
}

func generateProofs(rootNode *Node) []*Proof {
	return resolveBranchProofs(rootNode, []byte{}, 0)
}

// resolveBranchProofs 从root node 递归搜索叶子节点并为其生成证明
func resolveBranchProofs(node *Node, proof []byte, depth int) (proofs []*Proof) {

	if node.Type == LeafNodeType {
		p := &Proof{
			Offest: node.MaxByteRange - 1,
			Proof: ConcatBuffer(
				proof,
				node.DataHash,
				PaddedBigBytes(big.NewInt(int64(node.MaxByteRange)), 32),
			),
		}
		proofs = append(proofs, p)
		return
	}

	if node.Type == BranchNodeType {
		partialProof := ConcatBuffer(
			proof,
			node.LeftChild.ID,
			node.RightChild.ID,
			PaddedBigBytes(big.NewInt(int64(node.ByteRange)), 32),
		)
		leftProofs := resolveBranchProofs(node.LeftChild, partialProof, depth+1)
		rightProofs := resolveBranchProofs(node.RightChild, partialProof, depth+1)
		proofs = append(append(proofs, leftProofs...), rightProofs...)
		return
	}

	// node type error then return nil
	return nil
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

// ConcatBuffer
func ConcatBuffer(buffers ...[]byte) []byte {
	totalLength := 0
	for i := 0; i < len(buffers); i++ {
		totalLength += len(buffers[i])
	}

	temp := make([]byte, 0, totalLength)

	for _, val := range buffers {
		temp = append(temp, val...)
	}
	return temp
}
