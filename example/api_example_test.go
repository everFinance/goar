package example

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/everFinance/goar/client"
	"github.com/stretchr/testify/assert"
	"testing"
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

	// 7. Arql

}

func Test_rsa(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	assert.NoError(t, err)

}
