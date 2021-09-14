package bundles

import (
	"bytes"
	"encoding/hex"
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/hamba/avro"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
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
// 	item01, err := newDataItem(owner, "0", "", "1", []byte("this is a data bundle tx test item03"), []types.Tag{{Name: "GOAR", Value: "test01-bundle"}})
// 	assert.NoError(t, err)
// 	signedItem01, err := item01.Sign(w)
// 	assert.NoError(t, err)
//
// 	target := "Goueytjwney8mRqbWBwuxbk485svPUWxFQojteZpTx8"
// 	item02, err := newDataItem(owner, "0", target, "2", []byte("this is a data bundle tx test04"), []types.Tag{{Name: "GOAR", Value: "test02-bundle"}})
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

	// post to bundler
	BUNDLER := "http://bundler.arweave.net:10000"
	resp, err := http.DefaultClient.Post(BUNDLER+"/tx", "application/octet-stream", bytes.NewReader(itemData01.itemBinary))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	t.Log(statusCode)
	assert.Equal(t, 200, statusCode)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	t.Log(string(body))

	bd, err := BundleDataItems(itemData01)
	assert.NoError(t, err)

	txId, err := bd.SubmitBundleTx(w, nil, 50)
	assert.NoError(t, err)
	t.Log(txId)
}

func TestBundleDataItems(t *testing.T) {
	tags := []types.Tag{
		{Name: "App-Name", Value: "myApp"},
		{Name: "App-Version", Value: "1.0.0"},
	}
	bb, err := serializeTags(tags)
	assert.NoError(t, err)
	t.Log(hex.EncodeToString(bb))
	t.Log(utils.Base64Encode(bb))

	aa := "04104170702d4e616d650a6d79417070164170702d56657273696f6e0a312e302e3000"
	aaBy, err := hex.DecodeString(aa)
	assert.NoError(t, err)
	t.Log(utils.Base64Encode(aaBy))
	tagsParser, err := avro.Parse(`{"type": "array", "items": {"type": "record", "name": "Tag", "fields": [{"name": "name", "type": "string"}, {"name": "value", "type": "string"}]}}`)
	assert.NoError(t, err)
	tt := make([]types.Tag, 0)
	err = avro.Unmarshal(tagsParser, aaBy, &tt)
	assert.NoError(t, err)
	t.Log(tt)

	ttbb := make([]types.Tag, 0)
	err = avro.Unmarshal(tagsParser, bb, &ttbb)
	assert.NoError(t, err)
	t.Log(ttbb)
}

func TestCreateDataItem(t *testing.T) {
	cli := goar.NewClient("https://arweave.net")
	// id := "ipVFFrAkLosTtk-M3J6wYq3MKpfE6zK75nMIC-oLVXw"
	id := "mTm5-TtpsfJvUCPXflFe-P7HO6kOy4E2pGbt6-DUs40"
	body, err := cli.DownloadChunkData(id)
	assert.NoError(t, err)
	t.Log("body: ", len(body))
	bd, err := RecoverBundleData(body)
	assert.NoError(t, err)
	for _, item := range bd.Items {
		err = item.Verify()
		assert.NoError(t, err)
	}
}
