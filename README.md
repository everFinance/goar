# goar

### Install

```
go get github.com/everFinance/goar
```

### Example

#### Send AR or Winston

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

	id, err := wallet.SendAR(
  //id, err := wallet.SendWinston( 
		big.NewFloat(1.0), // AR amount
		{{target}}, // target address
		[]types.Tag{},
	)

	fmt.Println(id, err) // {{id}}, nil
}

```

#### Send Data

```golang
tx, err := wallet.SendData(
  []byte("123"), // Data bytes
  []types.Tag{
    types.Tag{
      Name:  "testSendData",
      Value: "123",
    },
  },
)

fmt.Println(id, err) // {{id}}, nil
```

#### Send Data SpeedUp

Arweave occasionally experiences congestion, and a low Reward can cause a transaction to fail; use speedUp to accelerate the transaction.

```golang
speedUp := int64(50) // means reward = reward * 150%
tx, err := wallet.SendDataSpeedUp(
  []byte("123"), // Data bytes
  []types.Tag{
    types.Tag{
      Name:  "testSendDataSpeedUp",
      Value: "123",
    },
  },speedUp)

fmt.Println(id, err) // {{id}}, nil
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
- [x] GenerateIndepHash

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

tx.LastTx = anchor
tx.Owner = utils.Base64Encode(w.PubKey.N.Bytes())

if err = utils.SignTransaction(tx, w.PubKey, w.PrvKey); err != nil {
  return
}

id = tx.ID

uploader, err := goar.CreateUploader(w.Client, tx, nil)
if err != nil {
  return
}

err = uploader.Once()
if err != nil {
  return
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
assert.NoError(t, uploader.Once())
```

##### NOTE: About all chunk transfer full example can be viewed in path `./example/chunks_tx_test.go`

---
### About Arweave Bundles
1. `goar` implemented creating,editing,reading and verifying bundles tx
2. This is the [ANS-104](https://github.com/joshbenaron/arweave-standards/blob/ans104/ans/ANS-104.md) standard protocol and refers to the [arbundles](https://github.com/Bundler-Network/arbundles) js-lib implement
3. more example can be viewed in path `./example/bundle_test.go`

#### CreateBundle
```go
w, err := goar.NewWalletFromPath(privateKey, arNode)
if err != nil {
  panic(err)
}

// Create Item
data := []byte("upload update...")   
signatureType := 1 // currently only supply type 1
target := "" // option 
anchor := "" // option
tags := []types.Tags{}{} // bundle item tags
item01, err := w.CreateAndSignBundleItem(data, 1, target, anchor, tags)    
// Same as create item
item02
item03
....

items := []types.BundleItem{item01, item02, item03 ...}	

bundle, err := utils.NewBundle(items...)

```

#### Send Item to Bundler
Bundler network provides guaranteed data seeding and instant data accessibility
```go
resp, err := w.Client.BatchSendItemToBundler(items,"") // The second parameter is the bundler gateway url，"" means use default url
```

#### Send Bundle Tx
```go
txId, err := w.SendBundleTx(bd.BundleBinary, arTxtags)
```

#### Get Bundle and Verify
```go
id := "lt24bnUGms5XLZeVamSPHePl4M2ClpLQyRxZI7weH1k"
bundle, err := cli.GetBundle(id)

// verify
for _, item := range bundle.Items {
  err = utils.VerifyBundleItem(item)
  assert.NoError(t, err)
}
```

### notice
if you call `w.Client.BatchSendItemToBundler(items,"")` 
and return `panic: send to bundler request failed; http code: 402`        
means that you have to pay ar to the bundler service address    
must use item signature address to transfer funds   

##### how to get bundler service address?
```go
curl --location --request GET 'https://node1.bundlr.network/info'

response:
{
"uptime": 275690.552536824,
"address": "OXcT1sVRSA5eGwt2k6Yuz8-3e3g9WJi5uSE99CWqsBs",
"gateway": "arweave.net"
}
```
This "address" is the bundler service receive ar address.    
You need to transfer a certain amount of ar to this address    
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
