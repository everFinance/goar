# goar

### Install

```
go get github.com/daqiancode/goar
```

### Example:
```golang

import (
	"fmt"
	"math/big"
	"github.com/daqiancode/goar/types"
	"github.com/daqiancode/goar"
)
func TestUploaderFile(t *testing.T) {
	w, err := goar.NewWalletFromPath("arweave-key.json", "https://arweave.net")
	assert.NoError(t, err)
	t.Log(w.Signer.Address)

	file, err := os.Open("test_resources/test.mp4")
	if err != nil {
		panic(err)
	}
	stat, err := file.Stat()
	assert.Nil(t, err)
	fileSize := stat.Size()
	reward, err := w.Client.GetTransactionPrice(fileSize)
	assert.Nil(t, err)
	tx := goar.NewSendFileTransaction(file, fileSize, reward, types.Tag{Name: "Content-Type", Value: "video/mp4"})
	err = w.SignTransaction(tx)
	assert.Nil(t, err)

	uploader, err := goar.CreateUploader(w.Client, tx, file, fileSize)
	totalSent := 0
	lastTime := time.Now().Unix()
	lastTotal := 0
	callbackCount := 0
	uploader.ProgressCallback = func(bytesSent int) {
		callbackCount += 1
		totalSent += bytesSent
		fmt.Println(bytesSent, totalSent, fileSize)
		fmt.Println("progress: ", totalSent/int(fileSize))
		if callbackCount%10 == 0 {
			now := time.Now().Unix()
			duration := now - lastTime
			if duration > 0 {
				speed := (totalSent - lastTotal) / 1024 / int(duration)
				fmt.Print("speed: ", speed, "KB/s")
				lastTotal = totalSent
				lastTime = now
			}

		}

	}
	assert.Nil(t, err)
	uploader.ConcurrentOnce(context.Background(), 10)
	// err = uploader.Once()
	assert.Nil(t, err)
	fmt.Println("https://arweave.net/" + tx.ID)
}

func TestSendBytes(t *testing.T) {
	w, err := goar.NewWalletFromPath("arweave-key.json", "https://arweave.net")
	assert.NoError(t, err)
	t.Log(w.Signer.Address)

	file := utils.NewReadBuffer([]byte("hello world"))
	assert.Nil(t, err)
	reward, err := w.Client.GetTransactionPrice(int64(file.Len()))
	assert.Nil(t, err)
	tags := []types.Tag{{Name: "content-type", Value: "text/plain"}}

	tx := goar.NewSendFileTransaction(file, int64(file.Len()), reward, tags...)
	err = w.SignTransaction(tx)
	assert.Nil(t, err)
	uploader, err := goar.CreateUploader(w.Client, tx, file, int64(file.Len()))
	assert.Nil(t, err)
	err = uploader.Once()
	assert.Nil(t, err)
	fmt.Println("https://arweave.net/" + tx.ID)
}


func main() {
	wallet, err := goar.NewWalletFromPath("./arweave-key.json", "https://arweave.net")
	if err != nil {
		panic(err)
	}

	tx, err := wallet.SendAR(
  //id, err := wallet.SendWinston( 
		big.NewFloat(1.0), // AR amount
		{{target}}, // target address
		[]types.Tag{},
	)

	fmt.Println(tx.ID, err)
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
- [x] BatchSendItemToBundler
- [x] GetBundle
- [x] GetTxDataFromPeers
- [x] BroadcastData
- [x] GetUnconfirmedTx
- [x] GetPendingTxIds
- [x] GetBlockHashList
- [x] ConcurrentDownloadChunkData

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
- [x] CreateAndSignBundleItem
- [x] SendBundleTxSpeedUp
- [x] SendBundleTx
- [x] SendPst

Initialize the instance, use a keyfile.json:

```golang
arWallet := goar.NewWalletFromPath("./keyfile.json")

// if your network is not good, you can config http proxy
proxyUrl := "http://127.0.0.1:8001"
arWallet := NewWalletFromPath("./keyfile.json", "https://arweave.net", proxyUrl)
```

#### Signer

- [x] SignTx
- [x] SignMsg
- [x] Owner

```golang
signer := goar.NewSignerFromPath("./keyfile.json")
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
- [x] NewBundle
- [x] NewBundleItem
- [x] SubmitItemToBundlr
- [x] SubmitItemToArSeed

#### RSA Threshold Cryptography

- [x] CreateTcKeyPair
- [x] ThresholdSign
- [x] AssembleSigShares
- [x] VerifySigShare

[Threshold Signature Usage Guidelines](https://github.com/daqiancode/goar/wiki/GOAR--RSA-Threshold-Signature-Usage-Guidelines)    

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
assert.NoError(t, uploader.Once())
```

##### NOTE: About all chunk transfer full example can be viewed in path `./example/chunks_tx_test.go`

---
### About Arweave Bundles
1. `goar` implemented creating,editing,reading and verifying bundles tx
2. This is the [ANS-104](https://github.com/joshbenaron/arweave-standards/blob/ans104/ans/ANS-104.md) standard protocol and refers to the [arbundles](https://github.com/Bundler-Network/arbundles) js-lib implement

#### Create Bundle Item
```go
signer, err := goar.NewSignerFromPath("./testKey.json") // rsa signer
// or 
signer, err := goether.NewSigner("0x.....") // ecdsa signer

// Create Item
data := []byte("aa bb cc dd")
target := "" // option 
anchor := "" // option
tags := []types.Tags{}{} // option bundle item tags
item01, err := itemSigner.CreateAndSignItem(data, target, anchor, tags)    
// Same as create item
item02
item03
....

```
#### assemble bundle and send to arweave network 
You can send items directly to the arweave network
```go

items := []types.BundleItem{item01, item02, item03 ...}
bundle, err := utils.NewBundle(items...)

w, err := goar.NewWalletFromPath("./key.json", arNode)

arTxTags := []types.Tags{}{} // option
tx, err := w.SendBundleTx(bd.BundleBinary, arTxtags)

```

#### Send Item to [Arseeding](https://github.com/everFinance/arseeding) gateway
Arseeding provides guaranteed data seeding and instant data accessibility
```go
arseedUrl := "https://seed.everpay.io"
currency := "USDC" // used for payment fee currency
resp, err := utils.SubmitItemToArSeed(item01,currency,arseedUrl)
```

#### Send Item to Bundler gateway
Bundler provides guaranteed data seeding and instant data accessibility
```go
bundlrUrl := "https://node1.bundlr.network"
resp, err := utils.SubmitItemToBundlr(item01, bundlrUrl)
```

#### Verify Bundle Items
```go

// verify
for _, item := range bundle.Items {
  err = utils.VerifyBundleItem(item)
  assert.NoError(t, err)
}
```
check [bundle example](./example/bundle_test.go) 

#### About Arseeding
if you can `utils.SubmitItemToArseed(item,currency,arseedUrl)` 
and you will get the following return response   
```go
{
    "ItemId": "5rEb7c6OjMQIYjl6P7AJIb4bB9CLMBSxhZ9N7BVbRCk",
    "bundler": "Fkj5J8CDLC9Jif4CzgtbiXJBnwXLSrp5AaIllleH_yY",
    "currency": "USDT",
    "decimals": 6,
    "fee": "701",
    "paymentExpiredTime": 1656044994,
    "expectedBlock": 960751
}
```
After you transfer 0.000701 USDT to bundler using everpay, arseeding will upload the item to arweave.   
For more usage, jump to [docs](https://github.com/everFinance/arseeding/blob/main/README.md)

#### About Bundlr
if you call `utils.SubmitItemToBundlr(item,bundlrUrl)` 
and return `panic: send to bundler request failed; http code: 402`        
means that you have to pay ar to the bundler service address    
must use item signature address to transfer funds   

##### how to get bundler service address?
```go
curl --location --request GET 'https://node1.bundlr.network/info'

response:
{
    "version": "0.2.0",
    "addresses": {
        "arweave": "OXcT1sVRSA5eGwt2k6Yuz8-3e3g9WJi5uSE99CWqsBs",
        "ethereum": "0xb4DE0833771eae55040b698aF5eB06d59E142C82",
        "matic": "0xb4DE0833771eae55040b698aF5eB06d59E142C82",
        "bnb": "0xb4DE0833771eae55040b698aF5eB06d59E142C82",
        "avalanche": "0xb4DE0833771eae55040b698aF5eB06d59E142C82",
        "solana": "DHyDV2ZjN3rB6qNGXS48dP5onfbZd3fAEz6C5HJwSqRD",
        "arbitrum": "0xb4DE0833771eae55040b698aF5eB06d59E142C82",
        "boba-eth": "0xb4DE0833771eae55040b698aF5eB06d59E142C82",
        "boba": "0xb4DE0833771eae55040b698aF5eB06d59E142C82",
        "chainlink": "0xb4DE0833771eae55040b698aF5eB06d59E142C82",
        "kyve": "0xb4DE0833771eae55040b698aF5eB06d59E142C82",
        "fantom": "0xb4DE0833771eae55040b698aF5eB06d59E142C82",
        "near": "bundlr1.near",
        "algorand": "DL7ZTTQMTFNXRF3367OTSNAZ3L2X676OJ4GGB3DXMUJ37CCKJ5RJMEO6RI"
    },
    "gateway": "arweave.net"
}
```
This "addresses" are the bundler service receive address.    
You need to transfer a certain amount of token to this address    
and wait for 25 blocks to confirm the transaction before you can use the bundler service.    

You can also use the following api to query the balance in the bundler service.   
```
curl --location --request GET 'https://node1.bundlr.network/account/balance?address=Ii5wAMlLNz13n26nYY45mcZErwZLjICmYd46GZvn4ck'

response:
{
    "balance": 1000000000
}
```

---

The original & this project need refactory
