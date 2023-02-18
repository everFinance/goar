package utils

import (
	"bytes"
	"crypto/sha256"
	"io"
	"math"
	"math/big"

	"github.com/daqiancode/goar/types"
	"github.com/shopspring/decimal"
)

func GenerateChunksBytes(data []byte) types.Chunks {
	chunks := chunkDataBytes(data)
	leaves := generateLeaves(chunks)
	root := buildLayer(leaves, 0) // leaf node level == 0
	proofs := generateProofs(root)

	// Discard the last chunk & proof if it's zero length.
	lastChunk := chunks[len(chunks)-1]
	if lastChunk.MaxByteRange-lastChunk.MinByteRange == 0 {
		chunks = chunks[:len(chunks)-1]
		proofs = proofs[:len(proofs)-1]
	}

	return types.Chunks{
		DataRoot: root.ID,
		Chunks:   chunks,
		Proofs:   proofs,
	}
}
func chunkDataBytes(data []byte) (chunks []types.Chunk) {
	cursor := 0
	var rest = data
	// if data length > max size
	for len(rest) >= types.MAX_CHUNK_SIZE {
		chunkSize := types.MAX_CHUNK_SIZE

		// 查看下一轮的chunkSize 是否小于最小的size，如果是则在这轮中调整chunk size 的大小
		nextChunkSize := len(rest) - types.MAX_CHUNK_SIZE
		if nextChunkSize > 0 && nextChunkSize < types.MIN_CHUNK_SIZE {
			dec := decimal.NewFromFloat(math.Ceil(float64(len(rest) / 2)))
			chunkSize = int(dec.IntPart())
		}

		chunk := rest[:chunkSize]
		dataHash := sha256.Sum256(chunk)
		cursor += len(chunk)
		chunks = append(chunks, types.Chunk{
			DataHash:     dataHash[:],
			MinByteRange: cursor - len(chunk),
			MaxByteRange: cursor,
		})

		rest = rest[chunkSize:]
	}

	hash := sha256.Sum256(rest)
	chunks = append(chunks, types.Chunk{
		DataHash:     hash[:],
		MinByteRange: cursor,
		MaxByteRange: cursor + len(rest),
	})
	return
}

/**
 * Generates the data_root, chunks & proofs
 * needed for a transaction.
 *
 * This also checks if the last chunk is a zero-length
 * chunk and discards that chunk and proof if so.
 * (we do not need to upload this zero length chunk)
 *
 * @param data
 */
func GenerateChunks(data io.ReadSeeker, fileSize int64) (types.Chunks, error) {
	_, err := data.Seek(0, io.SeekStart)
	if err != nil {
		return types.Chunks{}, err
	}
	chunks, err := chunkData(data, fileSize)
	if err != nil {
		return types.Chunks{}, err
	}
	leaves := generateLeaves(chunks)
	root := buildLayer(leaves, 0) // leaf node level == 0
	proofs := generateProofs(root)

	// Discard the last chunk & proof if it's zero length.
	lastChunk := chunks[len(chunks)-1]
	if lastChunk.MaxByteRange-lastChunk.MinByteRange == 0 {
		chunks = chunks[:len(chunks)-1]
		proofs = proofs[:len(proofs)-1]
	}

	return types.Chunks{
		DataRoot: root.ID,
		Chunks:   chunks,
		Proofs:   proofs,
	}, nil
}

func chunkData(data io.ReadSeeker, fileSize int64) ([]types.Chunk, error) {
	n := fileSize / types.MAX_CHUNK_SIZE
	w := fileSize % types.MAX_CHUNK_SIZE
	if w != 0 {
		n++
	}
	buffer := make([]byte, types.MAX_CHUNK_SIZE)
	chunks := make([]types.Chunk, n)
	var i int64 = 0
	var p int = 0
	for ; i < n; i++ {
		if i >= n-2 && n > 1 && w != 0 {
			buffer = make([]byte, (types.MAX_CHUNK_SIZE+w)/2+1)
		}
		rc, err := data.Read(buffer)
		if err != nil {
			return nil, err
		}
		p += rc
		dataHash := sha256.Sum256(buffer[:rc])
		chunks[i] = types.Chunk{
			DataHash:     dataHash[:],
			MinByteRange: p - rc,
			MaxByteRange: p,
		}
	}

	return chunks, nil

}

func generateLeaves(chunks []types.Chunk) (leafs []*types.Node) {
	for _, chunk := range chunks {
		// hDataHash := sha256.Sum256(chunk.DataHash)
		// hMaxByteRange := sha256.Sum256(PaddedBigBytes(big.NewInt(int64(chunk.MaxByteRange)), 32))

		leafs = append(leafs, &types.Node{
			// ID: hashArray(
			// 	[][]byte{hDataHash[:], hMaxByteRange[:]},
			// ),
			ID: Hash([][]byte{
				Hash([][]byte{chunk.DataHash}),
				Hash([][]byte{intToBuffer(chunk.MaxByteRange)}),
			}),
			Type:         types.LeafNodeType,
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
func buildLayer(nodes []*types.Node, level int) (root *types.Node) {
	if len(nodes) == 1 {
		root = nodes[0]
		return
	}

	nextLayer := make([]*types.Node, 0, len(nodes)/2)
	for i := 0; i < len(nodes); i += 2 {
		leftNode := nodes[i]
		var rightNode *types.Node
		if i+1 < len(nodes) {
			rightNode = nodes[i+1]
		}
		nextLayer = append(nextLayer, hashBranch(leftNode, rightNode))
	}

	return buildLayer(nextLayer, level+1)
}

// hashBranch get branch node by child node
func hashBranch(leftNode, rightNode *types.Node) (branchNode *types.Node) {
	// 如果只有一个node，则该node 为branch node
	if rightNode == nil {
		return leftNode
	}
	hLeafNodeId := sha256.Sum256(leftNode.ID)
	hRightNodeId := sha256.Sum256(rightNode.ID)
	// hLeafNodeMaxRange := sha256.Sum256(PaddedBigBytes(big.NewInt(int64(leftNode.MaxByteRange)), 32))
	hLeafNodeMaxRange := sha256.Sum256(intToBuffer(leftNode.MaxByteRange))
	// id := hashArray([][]byte{hLeafNodeId[:], hRightNodeId[:], hLeafNodeMaxRange[:]})
	id := Hash([][]byte{hLeafNodeId[:], hRightNodeId[:], hLeafNodeMaxRange[:]})
	branchNode = &types.Node{
		Type:         types.BranchNodeType,
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

func generateProofs(rootNode *types.Node) []*types.Proof {
	return resolveBranchProofs(rootNode, []byte{}, 0)
}

// resolveBranchProofs 从root node 递归搜索叶子节点并为其生成证明
func resolveBranchProofs(node *types.Node, proof []byte, depth int) (proofs []*types.Proof) {

	if node.Type == types.LeafNodeType {
		p := &types.Proof{
			Offest: node.MaxByteRange - 1,
			Proof: ConcatBuffer(
				proof,
				node.DataHash,
				// PaddedBigBytes(big.NewInt(int64(node.MaxByteRange)), 32),
				intToBuffer(node.MaxByteRange),
			),
		}
		proofs = append(proofs, p)
		return
	}

	if node.Type == types.BranchNodeType {
		partialProof := ConcatBuffer(
			proof,
			node.LeftChild.ID,
			node.RightChild.ID,
			// PaddedBigBytes(big.NewInt(int64(node.ByteRange)), 32),
			intToBuffer(node.ByteRange),
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
		for j := 0; j < types.WordBytes && i > 0; j++ {
			i--
			buf[i] = byte(d)
			d >>= 8
		}
	}
}

// ConcatBuffer 更好的 slice append
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

type ValidateResult struct {
	Offset     int
	LeftBound  int
	RightBound int
	ChunkSize  int
}

// 验证 merkle path
func ValidatePath(id []byte, dest, leftBound, rightBound int, path []byte) (*ValidateResult, bool) {
	if rightBound <= 0 {
		return nil, false
	}

	if dest >= rightBound {
		return ValidatePath(id, 0, rightBound-1, rightBound, path)
	}
	if dest < 0 {
		return ValidatePath(id, 0, 0, rightBound, path)
	}
	if len(path) == types.HASH_SIZE+types.NOTE_SIZE {
		pathData := path[0:types.HASH_SIZE]
		endOffsetBuffer := path[len(pathData) : len(pathData)+types.NOTE_SIZE]

		pathDataHash := Hash([][]byte{
			Hash([][]byte{pathData}),
			Hash([][]byte{endOffsetBuffer}),
		})
		result := arrayCompare(id, pathDataHash)
		if result {
			return &ValidateResult{
				Offset:     rightBound - 1,
				LeftBound:  leftBound,
				RightBound: rightBound,
				ChunkSize:  rightBound - leftBound,
			}, true
		}
		return nil, false
	}

	left := path[0:types.HASH_SIZE]
	right := path[len(left) : len(left)+types.HASH_SIZE]
	offsetBuffer := path[len(left)+len(right) : len(left)+len(right)+types.NOTE_SIZE]
	offset := bufferToInt(offsetBuffer)

	remainder := path[len(left)+len(right)+len(offsetBuffer):]

	pathHash := Hash([][]byte{
		Hash([][]byte{left}),
		Hash([][]byte{right}),
		Hash([][]byte{offsetBuffer}),
	})

	if arrayCompare(id, pathHash) {
		if dest < offset {
			return ValidatePath(left, dest, leftBound, int(math.Min(float64(rightBound), float64(offset))), remainder)
		}
		return ValidatePath(right, dest, int(math.Max(float64(leftBound), float64(offset))), rightBound, remainder)
	}
	return nil, false
}

func bufferToInt(buf []byte) int {
	value := 0
	for i := 0; i < len(buf); i++ {
		value *= 256
		value += int(buf[i])
	}
	return value
}

func intToBuffer(note int) []byte {
	buffer := make([]byte, types.NOTE_SIZE)

	for i := len(buffer) - 1; i >= 0; i-- {
		byt := note % 256
		buffer[i] = byte(byt)
		note = (note - byt) / 256 // todo 在js 中 /  是包含小数的 eg: 1/2 = 0.5
	}
	return buffer
}

func arrayCompare(a, b []byte) bool {
	return bytes.Equal(a, b)
}

func Hash(data [][]byte) []byte {
	byte32 := sha256.Sum256(ConcatBuffer(data...))
	return byte32[:]
}
