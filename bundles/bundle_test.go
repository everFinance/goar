package bundles

import (
	"encoding/json"
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

func TestCreateDataItemJson(t *testing.T) {
	owner := utils.Base64Encode(w.PubKey.N.Bytes())
	item01, err := CreateDataItemJson(owner, "", "99", []byte("this is a data bundle tx test item01"), []types.Tag{{Name: "GOAR", Value: "test01-bundle"}})
	assert.NoError(t, err)
	target := "Goueytjwney8mRqbWBwuxbk485svPUWxFQojteZpTx8"
	item02, err := CreateDataItemJson(owner, target, "100", []byte("this is a data bundle tx test02"), []types.Tag{{Name: "GOAR", Value: "test02-bundle"}})
	assert.NoError(t, err)
	signedItem01, err := item01.Sign(w)
	assert.NoError(t, err)
	signedItem02, err := item02.Sign(w)
	assert.NoError(t, err)

	bundleData, err := (DataItemJson{}).BundleData(signedItem01, signedItem02)
	if err != nil {
		panic(err)
		return
	}

	bd, err := json.Marshal(&bundleData)
	assert.NoError(t, err)
	id, state, err := w.SendData(bd, []types.Tag{{Name: "bundle-tx", Value: "goar-test"}})
	assert.NoError(t, err)
	t.Log(state)
	t.Log(id)
}
