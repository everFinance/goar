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
	"github.com/everFinance/goar"
)

func main() {
	wallet, err := goar.NewWalletFromPath("./test-keyfile.json", "https://arweave.net")
	if err != nil {
		panic(err)
	}

	id, err := wallet.SendWinston(
		big.NewInt(1), // Winston amount
		{{target}}, // target address
		[]types.Tag{
			types.Tag{
				Name:  "testSendWinston",
				Value: "1",
			},
		},
	)

	fmt.Println(id, err) // {{id}}, nil
}

```

Send Data

```golang
package main

import (
	"fmt"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar"
)

func main() {
	wallet, err := goar.NewWalletFromPath("./test-keyfile.json", "https://arweave.net")
	if err != nil {
		panic(err)
	}

	id, err := wallet.SendData(
		[]byte("123"), // Data bytes
		[]types.Tag{
			types.Tag{
				Name:  "testSendData",
				Value: "123",
			},
		},
	)

	fmt.Println(id, err) // {{id}}, nil
}
```

Send Data SpeedUp

```golang
package main

import (
	"fmt"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar"
)

func main() {
	wallet, err := goar.NewWalletFromPath("./test-keyfile.json", "https://arweave.net")
	if err != nil {
		panic(err)
	}

	speedUp := int64(50) // means reward = reward * 150%
	id, err := wallet.SendDataSpeedUp(
		[]byte("123"), // Data bytes
		[]types.Tag{
			types.Tag{
				Name:  "testSendData",
				Value: "123",
			},
		},speedUp)

	fmt.Println(id, err) // {{id}}, nil
}
```
### Components

#### Client

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
arClient := goar.NewClient("https://arweave.net")

// if your network is not good, you can config http proxy
proxyUrl := "http://127.0.0.1:8001"
arClient := goar.NewClient("https://arweave.net", proxyUrl)
```

#### Wallet

- [x] SendAR
- [x] SendARSpeedUp
- [x] SendWinston
- [x] SendWinstonSpeedUp
- [x] SendData
- [x] SendDataSpeedUp
- [x] SendTransaction

Initialize the instance, use a keyfile.json:

```golang
arWallet := goar.NewWalletFromPath("./keyfile.json")

// if your network is not good, you can config http proxy
proxyUrl := "http://127.0.0.1:8001"
arWallet := NewWalletFromPath("./keyfile.json", "https://arweave.net", proxyUrl)
```

#### Utils

Package for Arweave develop toolkit.

- [x] Base64Encode
- [x] Base64Decode
- [x] Sign
- [x] Verify
- [x] DeepHash
- [x] GenerateChunks
- [x] ValidatePath
- [x] OwnerToAddress
- [x] OwnerToPubKey
- [x] TagsEncode
- [x] TagsDecode
- [x] PrepareChunks
- [x] GetChunk
- [x] SignTransaction
- [x] GetSignatureData
- [x] VerifyTransaction

#### RSA Threshold Cryptography

- [x] CreateTcKeyPair
- [x] ThresholdSign
- [x] AssembleSigShares
- [x] VerifySigShare

[Threshold Signature Usage Guidelines](https://github.com/everFinance/goar/wiki/GOAR--RSA-Threshold-Signature-Usage-Guidelines)    

Create RSA Threshold Cryptography:

```golang
bitSize := 512 // If the values are 2048 and 4096, then the generation functions below will perform minute-level times, and we need 4096 bits as the maximum safety level for production environments.
l := 5
k := 3
keyShares, keyMeta, err := goar.CreateTcKeyPair(bitSize, k, l)
```

New sign instance:

```golang
exampleData := []byte("aaabbbcccddd112233") // need sign data
ts, err := goar.NewTcSign(keyMeta, exampleData)

// signer threshold sign
signer01 := keyShares[0]
signedData01, err := ts.ThresholdSign(signer01)

// assemble sign
signedShares := tcrsa.SigShareList{
signedData01,
...
}
signature, err := ts.AssembleSigShares(signedShares)

// verify share sign 
err := ts.VerifySigShare(signer01)
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

```golang
    arNode := "https://arweave.net"
	w, err := goar.NewWalletFromPath("../example/testKey.json", arNode) // your wallet private key
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
		Tags:     utils.TagsEncode(tags),
		Data:     utils.Base64Encode(data),
		DataSize: fmt.Sprintf("%d", len(data)),
		Reward:   fmt.Sprintf("%d", reward*(100+speedFactor)/100),
	}
	if err = utils.SignTransaction(tx, w.PubKey, w.PrvKey); err != nil {
		return
	}

	id = tx.ID

	uploader, err := goar.CreateUploader(w.Client, tx, nil)
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

```golang
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

```golang
    bigData, err := ioutil.ReadFile(filePath)
    txId := "myTxId"

    // get uploader by txId and post big data by chunks
	uploader, err := goar.CreateUploader(wallet.Client, txId, bigData)
	assert.NoError(t, err)
	for !uploader.IsComplete() {
		err := uploader.UploadChunk()
		assert.NoError(t, err)
	}
```

##### NOTE: About all chunk transfer full example can be viewed in path `./example/chunks_tx_test.go`
---
