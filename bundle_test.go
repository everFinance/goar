package goar

import (
	"github.com/everFinance/goar/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	privateKey = "./example/testKey.json" // your private key file
	arNode     = "https://arweave.net"
)

func TestBundleData_SubmitBundleTx(t *testing.T) {
	w, err := NewWalletFromPath(privateKey, arNode)
	if err != nil {
		panic(err)
	}

	target := "Fkj5J8CDLC9Jif4CzgtbiXJBnwXLSrp5AaIllleH_yY"
	tags := []types.Tag{
		{Name: "goar", Value: "bundleTx"},
		{Name: "App-Version", Value: "2.0.0"},
	}

	item01 , err := w.CreateBundleDataItem([]byte("goar bundle tx 01"), 1, target, "", tags)
	assert.NoError(t, err)
	item02 , err :=  w.CreateBundleDataItem([]byte("goar bundle tx 02"),  1, target, "", tags)
	assert.NoError(t, err)
	item03 , err :=  w.CreateBundleDataItem([]byte("goar bundle tx 03"),  1, target, "", tags)
	assert.NoError(t, err)
	item04 , err :=  w.CreateBundleDataItem([]byte("goar bundle tx 04"),  1, target, "", tags)
	assert.NoError(t, err)
	item05 , err :=  w.CreateBundleDataItem([]byte("goar bundle tx 05"),  1, target, "", tags)
	assert.NoError(t, err)
	item06 , err :=  w.CreateBundleDataItem([]byte("goar bundle tx 06"), 1, target, "", tags)
	assert.NoError(t, err)
	item07 , err :=  w.CreateBundleDataItem( []byte("goar bundle tx 07"), 1, target, "", tags)
	assert.NoError(t, err)
	item08 , err :=  w.CreateBundleDataItem([]byte("goar bundle tx 08"), 1, target, "", tags)
	assert.NoError(t, err)
	item09 , err :=  w.CreateBundleDataItem([]byte("goar bundle tx 09"),  1, target, "", tags)
	assert.NoError(t, err)
	item10 , err :=  w.CreateBundleDataItem([]byte("goar bundle tx 10"), 1, target, "", tags)
	assert.NoError(t, err)

	items := []DataItem{item01,item02,item03,item04,item05,item06,item07,item08,item09,item10}

	// send item to bundler gateway
	for _, item := range items {
		resp, err := w.Client.SendToBundler(item.itemBinary)
		assert.NoError(t, err)
		t.Log(resp.Id)
	}

	// assemble items
	arTxtags := []types.Tag{
		{Name: "GOAR", Value: "bundleTx"},
		{Name: "ACTION", Value: "test tx"},
	}
	bd, err := NewBundleData(items...)
	assert.NoError(t, err)

	txId, err := w.SubmitBundleTx(bd.bundleBinary, arTxtags, 50)
	assert.NoError(t, err)
	t.Log(txId)
}

func TestVerifyDataItem(t *testing.T) {
	cli := NewClient("https://arweave.net")
	// id := "K0JskpURZ-zZ7m01txR7hArvsBDDi08S6-6YIVQoc_Y" // big size data
	// id := "mTm5-TtpsfJvUCPXflFe-P7HO6kOy4E2pGbt6-DUs40"

	// goar test tx
	// id := "ipVFFrAkLosTtk-M3J6wYq3MKpfE6zK75nMIC-oLVXw"
	// id := "2ZFhlTJlFbj8XVmBtnBHS-y6Clg68trcRgIKBNemTM8"
	id := "lt24bnUGms5XLZeVamSPHePl4M2ClpLQyRxZI7weH1k"
	bd, err := cli.GetBundleData(id)
	assert.NoError(t, err)
	for _, item := range bd.Items {
		err = item.Verify()
		assert.NoError(t, err)
	}
}