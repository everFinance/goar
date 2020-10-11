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
		big.NewInt(1),
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
		[]byte("123"),
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

### Development

#### Test

```
make test
```