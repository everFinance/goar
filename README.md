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
	wallet, err := wallet.NewFromPath("./test-keyfile.json")
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
	wallet, err := wallet.NewFromPath("./test-keyfile.json")
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