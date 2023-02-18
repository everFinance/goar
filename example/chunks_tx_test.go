package example

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/daqiancode/goar"
	"github.com/daqiancode/goar/types"
	"github.com/daqiancode/goar/utils"
	"github.com/stretchr/testify/assert"
)

const (
	privateKey = "./testKey.json" // your private key file
	arNode     = "https://arweave.net"
)

var wallet *goar.Wallet

func init() {
	var err error
	wallet, err = goar.NewWalletFromPath(privateKey, arNode)
	if err != nil {
		panic(err)
	}
}

func assemblyDataTx(bigData []byte, wallet *goar.Wallet, tags []types.Tag) (*types.Transaction, error) {
	reward, err := wallet.Client.GetTransactionPrice(bigData, nil)
	if err != nil {
		return nil, err
	}
	tx := &types.Transaction{
		Format:   2,
		Target:   "",
		Quantity: "0",
		Tags:     utils.TagsEncode(tags),
		Data:     utils.Base64Encode(bigData),
		DataSize: fmt.Sprintf("%d", len(bigData)),
		Reward:   fmt.Sprintf("%d", reward),
	}
	anchor, err := wallet.Client.GetTransactionAnchor()
	if err != nil {
		return nil, err
	}
	tx.LastTx = anchor
	tx.Owner = wallet.Owner()

	signData, err := utils.GetSignatureData(tx)
	if err != nil {
		return nil, err
	}

	sign, err := wallet.Signer.SignMsg(signData)
	if err != nil {
		return nil, err
	}

	txHash := sha256.Sum256(sign)
	tx.ID = utils.Base64Encode(txHash[:])

	tx.Signature = utils.Base64Encode(sign)
	return tx, nil
}

// test upload post big size data by chunks
func Test_PostBigDataByChunks(t *testing.T) {
	filePath := "./testFile/2mbFile.pdf"
	bigData, err := ioutil.ReadFile(filePath)
	assert.NoError(t, err)

	tags := []types.Tag{{Name: "Content-Type", Value: "application/pdf"}, {Name: "goar", Value: "testdata"}}
	tx, err := assemblyDataTx(bigData, wallet, tags)
	assert.NoError(t, err)
	t.Log("txHash: ", tx.ID)

	// uploader Transaction
	uploader, err := goar.CreateUploader(wallet.Client, tx, nil)
	assert.NoError(t, err)
	assert.NoError(t, uploader.Once())
}

// test retry upload(断点重传) post big size data by tx id
func Test_RetryUploadDataByTxId(t *testing.T) {
	filePath := "./testFile/3mPhoto.jpg"
	bigData, err := ioutil.ReadFile(filePath)
	assert.NoError(t, err)

	tags := []types.Tag{{Name: "Content-Type", Value: "application/jpg"}, {Name: "goar", Value: "testdata"}}

	tx, err := assemblyDataTx(bigData, wallet, tags)
	assert.NoError(t, err)
	t.Log("txHash: ", tx.ID)

	// 1. post this tx without data
	tx.Data = ""
	body, status, err := wallet.Client.SubmitTransaction(tx)
	assert.NoError(t, err)
	t.Logf("post tx without data; body: %s, status: %d", string(body), status)

	// 2. watcher this tx from ar chain and must be sure the tx on chain
	getTxOnchain := func() bool {
		_, err := wallet.Client.GetTransactionByID(tx.ID)
		if err != nil {
			t.Log("watcher tx status: ", string(body))
			return false
		} else {
			return true
		}
	}

	// watcher tx
	for !getTxOnchain() {
		t.Log("sleep 10s ...")
		time.Sleep(10 * time.Second)
	}

	// get uploader by txId and post big data by chunks
	uploader, err := goar.CreateUploader(wallet.Client, tx.ID, bigData)
	assert.NoError(t, err)
	assert.NoError(t, uploader.Once())
}

// test continue upload(断点续传) big size data by last time uploader
func Test_ContinueUploadDataByLastUploader(t *testing.T) {
	filePath := "./testFile/1.8mPhoto.jpg"
	bigData, err := ioutil.ReadFile(filePath)
	assert.NoError(t, err)

	tags := []types.Tag{{Name: "Content-Type", Value: "application/jpg"}, {Name: "goar", Value: "1.8mbPhoto"}}
	tx, err := assemblyDataTx(bigData, wallet, tags)
	assert.NoError(t, err)
	t.Log("txHash: ", tx.ID)

	// 1. Upload a portion of data, for test when uploaded chunk == 2 and stop upload
	uploader, err := goar.CreateUploader(wallet.Client, tx, nil)
	assert.NoError(t, err)
	// only upload 2 chunks to ar chain
	for !uploader.IsComplete() && uploader.ChunkIndex <= 2 {
		err := uploader.UploadChunk()
		assert.NoError(t, err)
	}

	// then store uploader object to file
	jsonUploader, err := json.Marshal(uploader)
	assert.NoError(t, err)
	err = ioutil.WriteFile("./jsonUploaderFile.json", jsonUploader, 0777)
	assert.NoError(t, err)
	t.Log("sleep time ...")
	time.Sleep(5 * time.Second) // sleep 5s

	// 2. read uploader object from jsonUploader.json file and continue upload by last time uploader
	uploaderBuf, err := ioutil.ReadFile("./jsonUploaderFile.json")
	assert.NoError(t, err)
	lastUploader := &goar.TransactionUploader{}
	err = json.Unmarshal(uploaderBuf, lastUploader)
	assert.NoError(t, err)

	// new uploader object by last time uploader
	newUploader, err := goar.CreateUploader(wallet.Client, lastUploader.FormatSerializedUploader(), bigData)
	assert.NoError(t, err)
	assert.NoError(t, newUploader.Once())

	// end remove jsonUploaderFile.json file
	_ = os.Remove("./jsonUploaderFile.json")
}

func Test_aa(t *testing.T) {
	t.Log("address: ", wallet.Signer.Address)

	ownerBy := wallet.Signer.PubKey.N.Bytes()
	t.Log("length: ", len(ownerBy))
	owner := utils.Base64Encode(ownerBy)
	t.Log("owner:", owner)
}

func Test_dd(t *testing.T) {
	prv, err := rsa.GenerateKey(rand.Reader, 4096)
	assert.NoError(t, err)
	t.Log("size: ", len(prv.PublicKey.N.Bytes()))
}
