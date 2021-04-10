package bundles

import (
	"encoding/json"
	client2 "github.com/everFinance/goar/client"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	wallet2 "github.com/everFinance/goar/wallet"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	privateKey = "../example/testKey.json" // your private key file
	arNode     = "https://arweave.net"
)

var w *wallet2.Wallet

func init() {
	var err error
	w, err = wallet2.NewFromPath(privateKey, arNode)
	if err != nil {
		panic(err)
	}
}

// bundle data test
func TestDataItemJson_BundleData(t *testing.T) {
	// 1. new dataItem
	owner := utils.Base64Encode(w.PubKey.N.Bytes())
	item01, err := CreateDataItemJson(owner, "", "1", []byte("this is a data bundle tx test item03"), []types.Tag{{Name: "GOAR", Value: "test01-bundle"}})
	assert.NoError(t, err)
	signedItem01, err := item01.Sign(w)
	assert.NoError(t, err)

	target := "Goueytjwney8mRqbWBwuxbk485svPUWxFQojteZpTx8"
	item02, err := CreateDataItemJson(owner, target, "2", []byte("this is a data bundle tx test04"), []types.Tag{{Name: "GOAR", Value: "test02-bundle"}})
	assert.NoError(t, err)
	signedItem02, err := item02.Sign(w)
	assert.NoError(t, err)

	// 2. verify and assemble dataItem to BundleData
	bundleData, err := (DataItemJson{}).BundleData(signedItem01, signedItem02)
	if err != nil {
		panic(err)
		return
	}

	// 3. json serialization bundle data
	bd, err := json.Marshal(&bundleData)
	assert.NoError(t, err)

	// 4. send transaction include bundle data to ar chain
	id, state, err := w.SendData(bd, []types.Tag{{Name: "Bundle-Format", Value: "json"}, {Name: "Bundle-Version", Value: "1.0.0"}, {Name: "Content-Type", Value: "application/json"}})
	assert.NoError(t, err)
	t.Log(state)
	t.Log(id)
}

// unBundle data test
func TestDataItemJson_UnBundleData(t *testing.T) {
	id := "toiw3oFbG3B8cxeRVoc8MXrp8Zv0LJYHp5dVyPkEVkk"
	c := client2.New(arNode)
	// 1. get bundle txData type transaction txData
	txData, err := c.GetTransactionData(id, "json")
	assert.NoError(t, err)
	// 2. unBundle txData
	items, err := (DataItemJson{}).UnBundleData(txData)
	assert.NoError(t, err)

	// decode tags for test
	for i, item := range items {
		tags := item.Tags
		items[i].Tags, _ = types.TagsDecode(tags)
	}

	assert.Equal(t, 2, len(items))
	assert.Equal(t, "1", items[0].Nonce)
	assert.Equal(t, "2", items[1].Nonce)
}
