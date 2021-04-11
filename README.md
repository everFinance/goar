# goar

### Install

```
go get github.com/everFinance/goar
```

### Example

Send winston

```golang
package main

import (
	"fmt"
	"math/big"

	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/wallet"
)

func main() {
	wallet, err := wallet.NewFromPath("./test-keyfile.json", "https://arweave.net")
	if err != nil {
		panic(err)
	}

	id, stat, err := wallet.SendWinston(
		big.NewInt(1), // Winston amount
		{{target}}, // target address
		[]types.Tag{
			types.Tag{
				Name:  "testSendWinston",
				Value: "1",
			},
		},
	)

	fmt.Println(id, stat, err) // {{id}}, Pending, nil
}

```

Send Data

```golang
package main

import (
	"fmt"

	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/wallet"
)

func main() {
	wallet, err := wallet.NewFromPath("./test-keyfile.json", "https://arweave.net")
	if err != nil {
		panic(err)
	}

	id, stat, err := wallet.SendData(
		[]byte("123"), // Data bytes
		[]types.Tag{
			types.Tag{
				Name:  "testSendData",
				Value: "123",
			},
		},
	)

	fmt.Println(id, stat, err) // {{id}}, Pending, nil
}
```

### Golang Package

#### client

- [x] GetInfo
- [x] GetTransactionByID
- [x] GetTransactionStatus
- [x] GetTransactionField
- [x] GetTransactionData
- [x] GetTransactionPrice
- [x] GetTransactionAnchor
- [x] SubmitTransaction
- [x] Arql(Deprecated)
- [x] GraphQL
- [x] GetWalletBalance
- [x] GetLastTransactionID
- [x] GetBlockByID
- [x] GetBlockByHeight

Initialize the instance:

```golang
arClient := New("https://arweave.net")
```

#### wallet

- [x] SendAR
- [x] SendWinston
- [x] SendData
- [x] SendDataSpeedUp
- [x] SendTransaction

Initialize the instance, use a keyfile.json:

```golang
arWallet := NewFromPath("./keyfile.json")
```

### Development

#### Test

```
make test
```
---
### About chunks
1. First, we use Chunk transactions for all types of transactions in this library, so we only support transactions where format equals 2.
2. Second, the library already encapsulates a common interface for sending transactions : e.g `SendAR; SendData`. The user only needs to call this interface to send the transaction and do not need to worry about the usage of chunks.
3. The third，If the user needs to control the transaction such as breakpoint retransmission and breakpoint continuation operations. Here is how to do it.

#### chunked uploading advanced options
##### upload all transaction data
The method of sumbitting a data transaction is to use chunk uploading. This method will allow larger transaction zises,resuming a transaction upload if it's interrupted and give progress updates while uploading.
Simple example:
```
    arNode := "https://arweave.net"
	w, err := NewFromPath("../example/testKey.json", arNode) // your wallet private key
    anchor, err := w.Client.GetTransactionAnchor()
	if err != nil {
		return
	}
	data, err := ioutil.ReadFile("./2.3MBPhoto.jpg")
	if err != nil {
	    return
	}
	tx.LastTx = anchor
    reward, err := w.Client.GetTransactionPrice(data, nil)
	if err != nil {
		return
	}

	tx := &types.Transaction{
		Format:   2,
		Target:   "",
		Quantity: "0",
		Tags:     types.TagsEncode(tags),
		Data:     data,
		DataSize: fmt.Sprintf("%d", len(data)),
		Reward:   fmt.Sprintf("%d", reward*(100+speedFactor)/100),
	}
	if err = tx.SignTransaction(w.PubKey, w.PrvKey); err != nil {
		return
	}

	id = tx.ID

	uploader, err := uploader.CreateUploader(w.Client, tx, nil)
	if err != nil {
		return
	}
	for !uploader.IsComplete() {
		err = uploader.UploadChunk()
		if err != nil {
			return
		}
	}
```
##### Breakpoint continuingly
You can resume an upload from a saved uploader object, that you have persisted in storage some using json.marshal(uploader) at any stage of the upload.To resume, parse it back into an object and pass it to getUploader() along with the transactions data:
```

    uploaderBuf, err := ioutil.ReadFile("./jsonUploaderFile.json")
	lastUploader := &txType.TransactionUploader{}
	err = json.Unmarshal(uploaderBuf, lastUploader)
	assert.NoError(t, err)

	// new uploader object by last time uploader
	newUploader, err := txType.CreateUploader(wallet.Client, lastUploader.FormatSerializedUploader(), bigData)
	assert.NoError(t, err)
	for !newUploader.IsComplete() {
		err := newUploader.UploadChunk()
		assert.NoError(t, err)
	}
```
When resuming the upload, you must provide the same data as the original upload. When you serialize the uploader object with json.marshal() to save it somewhere, it will not include the data.
##### Breakpoint retransmission
You can also resume an upload from just the transaction ID and data, once it has been mined into a block. This can be useful if you didn't save the uploader somewhere but the upload got interrupted. This will re-upload all of the data from the beginning, since we don't know which parts have been uploaded:
```

    bigData, err := ioutil.ReadFile(filePath)
    txId := "myTxId"

    // get uploader by txId and post big data by chunks
	uploader, err := txType.CreateUploader(wallet.Client, txId, bigData)
	assert.NoError(t, err)
	for !uploader.IsComplete() {
		err := uploader.UploadChunk()
		assert.NoError(t, err)
	}
```

##### NOTE: About all chunk transfer full example can be viewed in path `./example/chunks_tx_test.go`
---