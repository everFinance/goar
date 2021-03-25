package transfer

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/everFinance/goar/client"
	"github.com/everFinance/goar/merkle"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/zyjblockchain/sandy_log/log"
	"math/big"
	"strconv"
)

type TransactionChunks struct {
	Format    int         `json:"format"`
	ID        string      `json:"id"`
	LastTx    string      `json:"last_tx"`
	Owner     string      `json:"owner"`
	Tags      []types.Tag `json:"tags"`
	Target    string      `json:"target"`
	Quantity  string      `json:"quantity"`
	Data      []byte      `json:"data"`
	DataSize  string      `json:"data_size"`
	DataRoot  string      `json:"data_root"`
	Reward    string      `json:"reward"`
	Signature string      `json:"signature"`

	// Computed when needed.
	Chunks *merkle.Chunks `json:"-"`
}

func (tx *TransactionChunks) FormatTransaction() *types.Transaction {
	return &types.Transaction{
		Format:    tx.Format,
		ID:        tx.ID,
		LastTx:    tx.LastTx,
		Owner:     tx.Owner,
		Tags:      tx.Tags,
		Target:    tx.Target,
		Quantity:  tx.Quantity,
		Data:      tx.Data,
		DataSize:  tx.DataSize,
		DataRoot:  tx.DataRoot,
		Reward:    tx.Reward,
		Signature: tx.Signature,
	}
}

func (tx *TransactionChunks) PrepareChunks(data []byte) {
	// Note: we *do not* use `this.Data`, the caller may be
	// operating on a Transaction with an zero length Data field.
	// This function computes the chunks for the Data passed in and
	// assigns the result to this Transaction. It should not read the
	// Data *from* this Transaction.

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
// Similar to `prepareChunks()` this does not operate `this.Data`,
// instead using the Data passed in.
func (tx *TransactionChunks) GetChunk(idx int, data []byte) (*GetChunk, error) {
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

// GetUploader
// @param upload: TransactionChunks | SerializedUploader | string,
// @param Data the Data of the Transaction. Required when resuming an upload.
func GetUploader(api *client.Client, upload interface{}, data []byte) (*TransactionUploader, error) {
	var (
		uploader *TransactionUploader
		err      error
	)

	if tt, ok := upload.(*TransactionChunks); ok {
		uploader, err = NewTransactionUploader(tt, api)
		if err != nil {
			return nil, err
		}
		return uploader, nil
	}

	if id, ok := upload.(string); ok {
		// upload 返回为 SerializedUploader 类型
		upload, err = (&TransactionUploader{Client: api}).FromTransactionId(id)
		if err != nil {
			log.Errorf("(&TransactionUploader{Client: api}).FromTransactionId(id) error: %v", err)
			return nil, err
		}
	} else {
		// 最后 upload 为 SerializedUploader type
		newUpload, ok := upload.(*SerializedUploader)
		if !ok {
			panic("upload params error")
		}
		upload = newUpload
	}

	uploader, err = (&TransactionUploader{Client: api}).FromSerialized(upload.(*SerializedUploader), data)
	return uploader, err
}

func (tx *TransactionChunks) SignTransaction(pubKey *rsa.PublicKey, prvKey *rsa.PrivateKey) error {
	tx.Owner = utils.Base64Encode(pubKey.N.Bytes())

	signData, err := GetSignatureData(tx)
	if err != nil {
		return err
	}
	sig, err := utils.Sign(signData, prvKey)
	if err != nil {
		return err
	}

	id := sha256.Sum256(sig)
	tx.ID = utils.Base64Encode(id[:])
	tx.Signature = utils.Base64Encode(sig)
	return nil
}

func GetSignatureData(tx *TransactionChunks) ([]byte, error) {
	switch tx.Format {
	case 1:
		// todo
		return nil, errors.New("current do not support format is 1 tx")
	case 2:
		tx.PrepareChunks(tx.Data)
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

func VerifyTransaction(tx TransactionChunks) (err error) {
	sig, err := utils.Base64Decode(tx.Signature)
	if err != nil {
		return
	}

	// verify ID
	id := sha256.Sum256(sig)
	if utils.Base64Encode(id[:]) != tx.ID {
		err = fmt.Errorf("wrong id")
	}

	signData, err := GetSignatureData(&tx)
	if err != nil {
		return
	}

	owner, err := utils.Base64Decode(tx.Owner)
	if err != nil {
		return
	}

	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(owner),
		E: 65537, //"AQAB"
	}

	return utils.Verify(signData, pubKey, sig)
}
