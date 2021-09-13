package bundles

import (
	"encoding/hex"
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	privateKey = "../example/testKey.json" // your private key file
	arNode     = "https://arweave.net"
)

var w *goar.Wallet

func init() {
	var err error
	w, err = goar.NewWalletFromPath(privateKey, arNode)
	if err != nil {
		panic(err)
	}
}

// bundle data test
// func TestDataItemJson_BundleData(t *testing.T) {
// 	// 1. new dataItem
// 	owner := utils.Base64Encode(w.PubKey.N.Bytes())
// 	item01, err := newDataItemJson(owner, "0", "", "1", []byte("this is a data bundle tx test item03"), []types.Tag{{Name: "GOAR", Value: "test01-bundle"}})
// 	assert.NoError(t, err)
// 	signedItem01, err := item01.Sign(w)
// 	assert.NoError(t, err)
//
// 	target := "Goueytjwney8mRqbWBwuxbk485svPUWxFQojteZpTx8"
// 	item02, err := newDataItemJson(owner, "0", target, "2", []byte("this is a data bundle tx test04"), []types.Tag{{Name: "GOAR", Value: "test02-bundle"}})
// 	assert.NoError(t, err)
// 	signedItem02, err := item02.Sign(w)
// 	assert.NoError(t, err)
//
// 	// 2. verify and assemble dataItem to BundleData
// 	bundleData, err := BundleDataItems(signedItem01, signedItem02)
// 	if err != nil {
// 		panic(err)
// 		return
// 	}
//
// 	// 3. json serialization bundle data
// 	bd, err := json.Marshal(&bundleData)
// 	assert.NoError(t, err)
//
// 	// 4. send transaction include bundle data to ar chain
// 	id, err := w.SendData(bd, BundleTags)
// 	assert.NoError(t, err)
// 	t.Log(id)
// }

// // unBundle data test
// func TestDataItemJson_UnBundleData(t *testing.T) {
// 	id := "A41r5OgQ2qwx0kkYEbbBQZJosnqY54Uz82O8W2upi6g"
// 	c := goar.NewClient(arNode)
// 	// 1. get bundle txData type transaction txData
// 	txData, err := c.GetTransactionData(id, "json")
// 	assert.NoError(t, err)
// 	// 2. unBundle txData
// 	items, err := UnBundleDataItems(txData)
// 	assert.NoError(t, err)
//
// 	// decode tags for test
// 	for i, item := range items {
// 		tags := item.Tags
// 		items[i].Tags, _ = utils.TagsDecode(tags)
// 	}
//
// 	assert.Equal(t, 2, len(items))
// 	t.Log(items)
// }

func TestBundleData_SubmitBundleTx(t *testing.T) {
	target := "Fkj5J8CDLC9Jif4CzgtbiXJBnwXLSrp5AaIllleH_yY"
	tags := []types.Tag{
		{Name: "App-Name", Value: "myApp"},
		{Name: "App-Version", Value: "1.0.0"},
	}
	itemData01, err := CreateDataItem(w, []byte("test02"), w.PubKey.N.Bytes(), 1, target, "", tags)
	assert.NoError(t, err)
	// itemData02, err := CreateDataItem(w,[]byte("sandy test goar bundle tx data02"),w.PubKey.N.Bytes(),1,target,"",tags)
	// assert.NoError(t, err)

	// bd, err := BundleDataItems(*itemData01,*itemData02)
	// assert.NoError(t, err)
	// txId, err := bd.SubmitBundleTx(w,nil)
	// assert.NoError(t, err)
	// t.Log(txId)

	t.Log(len(itemData01.binary))
	sig, _ := utils.Base64Decode(itemData01.Signature)
	t.Log(hex.EncodeToString(sig))
	t.Log(hex.EncodeToString(itemData01.binary))

	t.Log(hex.EncodeToString(w.PubKey.N.Bytes()))

	t.Log(itemData01.binary)

	// BUNDLER := "http://bundler.arweave.net:10000"
	// resp, err := http.DefaultClient.Post(BUNDLER + "/tx", "application/octet-stream", bytes.NewReader(itemData01.binary))
	// if err != nil {
	// 	return
	// }
	// defer resp.Body.Close()
	//
	// statusCode := resp.StatusCode
	// t.Log(statusCode)
	// body, err := ioutil.ReadAll(resp.Body)
	// assert.NoError(t, err)
	// t.Log(string(body))
}

func TestBundleDataItems(t *testing.T) {
	tags := []types.Tag{
		{Name: "App-Name", Value: "myApp"},
		{Name: "App-Version", Value: "1.0.0"},
	}
	bb, err := serializeTags(tags)
	assert.NoError(t, err)
	t.Log(hex.EncodeToString(bb))
	t.Log(bb)

	aa := "04104170702d4e616d650a6d79417070164170702d56657273696f6e0a312e302e3000"
	aaBy, err := hex.DecodeString(aa)
	assert.NoError(t, err)
	t.Log(aaBy)

}
