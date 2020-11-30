package example

import (
	"fmt"
	"testing"
	"time"

	"github.com/everFinance/goar/client"
	"github.com/everFinance/goar/types"
	wallet2 "github.com/everFinance/goar/wallet"
	"github.com/stretchr/testify/assert"
)

func Test_Client(t *testing.T) {
	// create client
	arNode := "https://arweave.net"
	c := client.New(arNode)
	txId := "hKMMPNh_emBf8v_at1tFzNYACisyMQNcKzeeE1QE9p8"

	// 1. getInfo
	nodeInfo, err := c.GetInfo()
	assert.NoError(t, err)
	t.Logf("%v", nodeInfo)

	// 2. full transaction via Id
	tx, state, err := c.GetTransactionByID(txId)
	assert.NoError(t, err)
	t.Logf("state: %s", state)
	t.Log(tx)

	// 3. get transaction field by id
	f, err := c.GetTransactionField(txId, "signature")
	assert.NoError(t, err)
	t.Log(f)

	// 4. get transaction data
	data, err := c.GetTransactionData(txId)
	assert.NoError(t, err)
	t.Log(string(data))
	data, err = c.GetTransactionData(txId, "html")
	assert.NoError(t, err)
	t.Log(string(data))

	// 5. get tx send current time reward
	reward, err := c.GetTransactionPrice(data, nil)
	assert.NoError(t, err)
	t.Log(reward)
	to := "1seRanklLU_1VTGkEk7P0xAwMJfA7owA1JHW5KyZKlY"
	reward, err = c.GetTransactionPrice([]byte{}, &to)
	assert.NoError(t, err)
	t.Log(reward)

	// 6. get anchor
	anchor, err := c.GetTransactionAnchor()
	assert.NoError(t, err)
	t.Log(anchor)

}

func Test_client2(t *testing.T) {
	arNode := "https://arweave.net"
	wallet, err := wallet2.NewFromPath("./testKey.json", arNode)
	assert.NoError(t, err)

	tag := []types.Tag{
		types.Tag{
			Name:  "TokenSymbol",
			Value: "DXN",
		},
		types.Tag{
			Name:  "Version",
			Value: "1.1.0",
		},
		types.Tag{
			Name:  "CreatedBy",
			Value: "ZYJ123",
		},
	}
	// 连续发送5 笔交易来测试交易打包顺序
	for i := 0; i < 5; i++ {
		data := fmt.Sprintf("nonce: %d", i)
		id, status, err := wallet.SendData([]byte(data), tag)
		t.Log(id)
		t.Log(status)
		t.Log(err)
		time.Sleep(30 * time.Second)
	}
}

func TestGetTransactionsStatus(t *testing.T) {
	arNode := "https://arweave.net"
	wallet, err := wallet2.NewFromPath("./testKey.json", arNode)
	assert.NoError(t, err)

	status, code, err := wallet.Client.GetTransactionStatus("ggt-x5Q_niHifdNzMxZrhiibKf0KQ-cJun0UIBBa-yA")
	assert.Equal(t, "Success", status)
	assert.Equal(t, 200, code)
	assert.NoError(t, err)
}

func Test_Arq(t *testing.T) {
	arqStr := `{
			"op": "and",
			"expr1": {
				"op": "equals",
				"expr1": "TokenSymbol",
				"expr2": "DXN"
			},
			"expr2": {
				"op": "equals",
				"expr1": "CreatedBy",
				"expr2": "zhou yu ji"
			}
		}`
	// create client
	arNode := "https://arweave.net"
	c := client.New(arNode)
	ids, err := c.Arql(arqStr)
	t.Log(len(ids))
	assert.NoError(t, err)
	sstr := make([]string, 0)
	sstr = append(sstr, "a")
	for _, val := range ids {
		t.Log(val)
	}

	// 冒泡排序
	for i := 0; i < len(ids); i++ {
		for j := 1; j < len(ids)-i; j++ {
			if ids[j] > ids[j-1] {
				ids[j], ids[j-1] = ids[j-1], ids[j]
			}
		}
	}
	t.Log(ids)
}
