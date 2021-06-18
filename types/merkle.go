package types

type Chunks struct {
	DataRoot []byte   `json:"data_root"`
	Chunks   []Chunk  `json:"chunks"`
	Proofs   []*Proof `json:"proofs"`
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
