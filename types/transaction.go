package types

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/everFinance/goar/merkle"
	"github.com/everFinance/goar/utils"
	"github.com/everFinance/sandy_log/log"
)

type Transaction struct {
	Format    int    `json:"format"`
	ID        string `json:"id"`
	LastTx    string `json:"last_tx"`
	Owner     string `json:"owner"`
	Tags      []Tag  `json:"tags"`
	Target    string `json:"target"`
	Quantity  string `json:"quantity"`
	Data      string `json:"data"` // base64.encode
	DataSize  string `json:"data_size"`
	DataRoot  string `json:"data_root"`
	Reward    string `json:"reward"`
	Signature string `json:"signature"`

	// Computed when needed.
	Chunks *merkle.Chunks `json:"-"`
}

func (tx *Transaction) AddSignature(signature []byte) {
	txId := sha256.Sum256(signature)
	tx.ID = utils.Base64Encode(txId[:])
	tx.Signature = utils.Base64Encode(signature)
}

func (tx *Transaction) PrepareChunks(data []byte) {
	// Note: we *do not* use `this.Data`, the caller may be
	// operating on a Transaction with an zero length Data field.
	// This function computes the chunks for the Data passed in and
	// assigns the result to this Transaction. It should not read the
	// Data *from* this Transaction.
	fmt.Printf("Tx data size: %fMB \n", float64(len(data))/1024.0/1024.0)
	if tx.Chunks == nil && len(data) > 0 {
		chunks := merkle.GenerateChunks(data)
		tx.Chunks = &chunks
		tx.DataRoot = utils.Base64Encode(tx.Chunks.DataRoot)
	}

	if tx.Chunks == nil && len(data) == 0 {
		tx.Chunks = &merkle.Chunks{
			DataRoot: make([]byte, 0),
			Chunks:   make([]merkle.Chunk, 0),
			Proofs:   make([]*merkle.Proof, 0),
		}
		tx.DataRoot = ""
	}
	return
}

type GetChunk struct {
	DataRoot string `json:"data_root"`
	DataSize string `json:"data_size"`
	DataPath string `json:"data_path"`
	Offset   string `json:"offset"`
	Chunk    string `json:"chunk"`
}

// Returns a chunk in a format suitable for posting to /chunk.
// Similar to `PrepareChunks()` this does not operate `this.Data`,
// instead using the Data passed in.
func (tx *Transaction) GetChunk(idx int, data []byte) (*GetChunk, error) {
	if tx.Chunks == nil {
		return nil, errors.New("Chunks have not been prepared")
	}

	proof := tx.Chunks.Proofs[idx]
	chunk := tx.Chunks.Chunks[idx]

	return &GetChunk{
		DataRoot: tx.DataRoot,
		DataSize: tx.DataSize,
		DataPath: utils.Base64Encode(proof.Proof),
		Offset:   strconv.Itoa(proof.Offest),
		Chunk:    utils.Base64Encode(data[chunk.MinByteRange:chunk.MaxByteRange]),
	}, nil
}

func (gc *GetChunk) Marshal() ([]byte, error) {
	return json.Marshal(gc)
}

func (tx *Transaction) SignTransaction(pubKey *rsa.PublicKey, prvKey *rsa.PrivateKey) error {
	tx.Owner = utils.Base64Encode(pubKey.N.Bytes())

	signData, err := GetSignatureData(tx)
	if err != nil {
		return err
	}
	sig, err := utils.Sign(signData, prvKey)
	if err != nil {
		return err
	}

	tx.AddSignature(sig)
	return nil
}

func GetSignatureData(tx *Transaction) ([]byte, error) {
	switch tx.Format {
	case 1:
		tags := make([]byte, 0)
		dcTags, err := TagsDecode(tx.Tags)
		if err != nil {
			return nil, err
		}
		for _, tag := range dcTags {
			tags = append(tags, merkle.ConcatBuffer([]byte(tag.Name), []byte(tag.Value))...)
		}

		data, err := utils.Base64Decode(tx.Data)
		if err != nil {
			return nil, err
		}
		owner, err := utils.Base64Decode(tx.Owner)
		if err != nil {
			return nil, err
		}
		target, err := utils.Base64Decode(tx.Target)
		if err != nil {
			return nil, err
		}

		lastTx, err := utils.Base64Decode(tx.LastTx)
		if err != nil {
			return nil, err
		}
		return merkle.ConcatBuffer(
			owner,
			target,
			data,
			[]byte(tx.Quantity),
			[]byte(tx.Reward),
			lastTx,
			tags,
		), nil

	case 2:
		data, err := utils.Base64Decode(tx.Data)
		if err != nil {
			log.Errorf("data, err := utils.Base64Decode(tx.Data) error:%v", err)
			return nil, err
		}
		tx.PrepareChunks(data)
		tags := [][]string{}
		for _, tag := range tx.Tags {
			tags = append(tags, []string{
				tag.Name, tag.Value,
			})
		}

		dataList := []interface{}{}
		dataList = append(dataList, utils.Base64Encode([]byte(fmt.Sprintf("%d", tx.Format))))
		dataList = append(dataList, tx.Owner)
		dataList = append(dataList, tx.Target)
		dataList = append(dataList, utils.Base64Encode([]byte(tx.Quantity)))
		dataList = append(dataList, utils.Base64Encode([]byte(tx.Reward)))
		dataList = append(dataList, tx.LastTx)
		dataList = append(dataList, tags)
		dataList = append(dataList, utils.Base64Encode([]byte(tx.DataSize)))
		dataList = append(dataList, tx.DataRoot)

		hash := utils.DeepHash(dataList)
		deepHash := hash[:]
		return deepHash, nil

	default:
		return nil, errors.New(fmt.Sprintf("Unexpected Transaction format: %d", tx.Format))
	}
}


func (tx Transaction) VerifyTransaction() (err error) {
	sig, err := utils.Base64Decode(tx.Signature)
	if err != nil {
		return
	}

	// verify ID
	id := sha256.Sum256(sig)
	if utils.Base64Encode(id[:]) != tx.ID {
		err = fmt.Errorf("wrong id")
		return
	}

	signData, err := GetSignatureData(&tx)
	if err != nil {
		return
	}

	pubKey, err := utils.OwnerToPubKey(tx.Owner)
	if err != nil {
		return
	}

	return utils.Verify(signData, pubKey, sig)
}
