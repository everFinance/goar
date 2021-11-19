package utils

import (
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"

	"github.com/everFinance/goar/types"
)

func PrepareChunks(tx *types.Transaction, data []byte) {
	// Note: we *do not* use `this.Data`, the caller may be
	// operating on a Transaction with an zero length Data field.
	// This function computes the chunks for the Data passed in and
	// assigns the result to this Transaction. It should not read the
	// Data *from* this Transaction.
	if tx.Chunks == nil && len(data) > 0 {
		chunks := GenerateChunks(data)
		tx.Chunks = &chunks
		tx.DataRoot = Base64Encode(tx.Chunks.DataRoot)
	}

	if tx.Chunks == nil && len(data) == 0 {
		tx.Chunks = &types.Chunks{
			DataRoot: make([]byte, 0),
			Chunks:   make([]types.Chunk, 0),
			Proofs:   make([]*types.Proof, 0),
		}
	}
	return
}

// Returns a chunk in a format suitable for posting to /chunk.
// Similar to `PrepareChunks()` this does not operate `this.Data`,
// instead using the Data passed in.
func GetChunk(tx types.Transaction, idx int, data []byte) (*types.GetChunk, error) {
	if tx.Chunks == nil {
		return nil, errors.New("Chunks have not been prepared")
	}

	proof := tx.Chunks.Proofs[idx]
	chunk := tx.Chunks.Chunks[idx]

	return &types.GetChunk{
		DataRoot: tx.DataRoot,
		DataSize: tx.DataSize,
		DataPath: Base64Encode(proof.Proof),
		Offset:   strconv.Itoa(proof.Offest),
		Chunk:    Base64Encode(data[chunk.MinByteRange:chunk.MaxByteRange]),
	}, nil
}

func SignTransaction(tx *types.Transaction, prvKey *rsa.PrivateKey) error {
	signData, err := GetSignatureData(tx)
	if err != nil {
		return err
	}
	sig, err := Sign(signData, prvKey)
	if err != nil {
		return err
	}

	txId := sha256.Sum256(sig)
	tx.ID = Base64Encode(txId[:])
	tx.Signature = Base64Encode(sig)
	return nil
}

func GetSignatureData(tx *types.Transaction) ([]byte, error) {
	switch tx.Format {
	case 1:
		tags := make([]byte, 0)
		dcTags, err := TagsDecode(tx.Tags)
		if err != nil {
			return nil, err
		}
		for _, tag := range dcTags {
			tags = append(tags, ConcatBuffer([]byte(tag.Name), []byte(tag.Value))...)
		}

		data, err := Base64Decode(tx.Data)
		if err != nil {
			return nil, err
		}
		owner, err := Base64Decode(tx.Owner)
		if err != nil {
			return nil, err
		}
		target, err := Base64Decode(tx.Target)
		if err != nil {
			return nil, err
		}

		lastTx, err := Base64Decode(tx.LastTx)
		if err != nil {
			return nil, err
		}
		return ConcatBuffer(
			owner,
			target,
			data,
			[]byte(tx.Quantity),
			[]byte(tx.Reward),
			lastTx,
			tags,
		), nil

	case 2:
		data, err := Base64Decode(tx.Data)
		if err != nil {
			return nil, err
		}
		PrepareChunks(tx, data)
		tags := [][]string{}
		for _, tag := range tx.Tags {
			tags = append(tags, []string{
				tag.Name, tag.Value,
			})
		}

		dataList := []interface{}{}
		dataList = append(dataList, Base64Encode([]byte(fmt.Sprintf("%d", tx.Format))))
		dataList = append(dataList, tx.Owner)
		dataList = append(dataList, tx.Target)
		dataList = append(dataList, Base64Encode([]byte(tx.Quantity)))
		dataList = append(dataList, Base64Encode([]byte(tx.Reward)))
		dataList = append(dataList, tx.LastTx)
		dataList = append(dataList, tags)
		dataList = append(dataList, Base64Encode([]byte(tx.DataSize)))
		dataList = append(dataList, tx.DataRoot)

		hash := DeepHash(dataList)
		deepHash := hash[:]
		return deepHash, nil

	default:
		return nil, errors.New(fmt.Sprintf("Unexpected Transaction format: %d", tx.Format))
	}
}

func VerifyTransaction(tx types.Transaction) (err error) {
	sig, err := Base64Decode(tx.Signature)
	if err != nil {
		return
	}

	// verify ID
	id := sha256.Sum256(sig)
	if Base64Encode(id[:]) != tx.ID {
		err = fmt.Errorf("wrong id")
		return
	}

	signData, err := GetSignatureData(&tx)
	if err != nil {
		return
	}

	pubKey, err := OwnerToPubKey(tx.Owner)
	if err != nil {
		return
	}

	return Verify(signData, pubKey, sig)
}
