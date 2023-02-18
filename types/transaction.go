package types

import (
	"encoding/json"
	"io"
)

type Transaction struct {
	Format    int    `json:"format"`
	ID        string `json:"id"`
	LastTx    string `json:"last_tx"`
	Owner     string `json:"owner"` // utils.Base64Encode(wallet.PubKey.N.Bytes())
	Tags      []Tag  `json:"tags"`
	Target    string `json:"target"`
	Quantity  string `json:"quantity"`
	Data      string `json:"data"` // base64.encode
	DataSize  string `json:"data_size"`
	DataRoot  string `json:"data_root"`
	Reward    string `json:"reward"`
	Signature string `json:"signature"`

	// Computed when needed.
	Chunks   *Chunks       `json:"-"`
	File     io.ReadSeeker `json:"-"`
	FileSize int64         `json:"-"`
}

type GetChunk struct {
	DataRoot  string `json:"data_root"`
	DataSize  string `json:"data_size"`
	DataPath  string `json:"data_path"`
	Offset    string `json:"offset"`
	Chunk     string `json:"chunk"`
	ChunkSize int    `json:"-"`
}

func (gc *GetChunk) Marshal() ([]byte, error) {
	return json.Marshal(gc)
}
