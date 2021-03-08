package utils

import (
	"encoding/json"
	"errors"
	"github.com/everFinance/goar/types"
	"strconv"
)

type Transaction struct {
	Format    int         `json:"format"`
	ID        string      `json:"id"`
	LastTx    string      `json:"last_tx"`
	Owner     string      `json:"owner"`
	Tags      []types.Tag `json:"tags"`
	Target    string      `json:"target"`
	Quantity  string      `json:"quantity"`
	Data      []byte      `json:"data"`
	DataSize  int         `json:"data_size"`
	DataRoot  string      `json:"data_root"`
	Reward    string      `json:"reward"`
	Signature string      `json:"signature"`

	// Computed when needed.
	Chunks *Chunks
}

func (tx *Transaction) PrepareChunks(data []byte) {
	// Note: we *do not* use `this.data`, the caller may be
	// operating on a transaction with an zero length data field.
	// This function computes the chunks for the data passed in and
	// assigns the result to this transaction. It should not read the
	// data *from* this transaction.

	if tx.Chunks == nil && len(data) > 0 {
		*tx.Chunks = GenerateChunks(data)
		tx.DataRoot = Base64Encode(tx.Chunks.DataRoot)
	}

	if tx.Chunks == nil && len(data) == 0 {
		tx.Chunks = &Chunks{
			DataRoot: make([]byte, 0),
			Chunks:   make([]Chunk, 0),
			Proofs:   make([]*Proof, 0),
		}
		tx.DataRoot = ""
	}
	return
}

type GetChunk struct {
	DataRoot string
	DataSize string
	DataPath string
	Offset   string
	Chunk    string
}

// Returns a chunk in a format suitable for posting to /chunk.
// Similar to `prepareChunks()` this does not operate `this.data`,
// instead using the data passed in.
func (tx *Transaction) GetChunk(idx int, data []byte) (*GetChunk, error) {
	if tx.Chunks == nil {
		return nil, errors.New("Chunks have not been prepared")
	}

	if len(tx.Chunks.Proofs) >= idx || len(tx.Chunks.Chunks) >= idx {
		return nil, errors.New("len(tx.Chunks.Proofs) >= idx || len(tx.Chunks.Chunks) >= idx")
	}

	proof := tx.Chunks.Proofs[idx]
	chunk := tx.Chunks.Chunks[idx]

	return &GetChunk{
		DataRoot: tx.DataRoot,
		DataSize: strconv.Itoa(tx.DataSize),
		DataPath: Base64Encode(proof.Proof),
		Offset:   strconv.Itoa(proof.Offest),
		Chunk:    Base64Encode(data[chunk.MinByteRange:chunk.MaxByteRange]),
	}, nil
}

func (gc *GetChunk) Marshal() ([]byte, error) {
	return json.Marshal(gc)
}
