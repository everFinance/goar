package example

import (
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBundle_SendBundleTx(t *testing.T) {
	privateKey := "./testKey.json" // your private key file
	arNode := "https://arweave.net"
	w, err := goar.NewWalletFromPath(privateKey, arNode)
	if err != nil {
		panic(err)
	}

	target := "Fkj5J8CDLC9Jif4CzgtbiXJBnwXLSrp5AaIllleH_yY"
	tags := []types.Tag{
		{Name: "goar", Value: "bundleTx"},
		{Name: "App-Version", Value: "2.0.0"},
	}

	item01, err := w.CreateAndSignBundleItem([]byte("goar bundle tx 01"), 1, target, "", tags)
	assert.NoError(t, err)
	item02, err := w.CreateAndSignBundleItem([]byte("goar bundle tx 02"), 1, target, "", tags)
	assert.NoError(t, err)
	item03, err := w.CreateAndSignBundleItem([]byte("goar bundle tx 03"), 1, target, "", tags)
	assert.NoError(t, err)
	item04, err := w.CreateAndSignBundleItem([]byte("goar bundle tx 04"), 1, target, "", tags)
	assert.NoError(t, err)
	item05, err := w.CreateAndSignBundleItem([]byte("goar bundle tx 05"), 1, target, "", tags)
	assert.NoError(t, err)
	item06, err := w.CreateAndSignBundleItem([]byte("goar bundle tx 06"), 1, target, "", tags)
	assert.NoError(t, err)
	item07, err := w.CreateAndSignBundleItem([]byte("goar bundle tx 07"), 1, target, "", tags)
	assert.NoError(t, err)
	item08, err := w.CreateAndSignBundleItem([]byte("goar bundle tx 08"), 1, target, "", tags)
	assert.NoError(t, err)
	item09, err := w.CreateAndSignBundleItem([]byte("goar bundle tx 09"), 1, target, "", tags)
	assert.NoError(t, err)
	item10, err := w.CreateAndSignBundleItem([]byte("goar bundle tx 10"), 1, target, "", tags)
	assert.NoError(t, err)

	items := []types.BundleItem{item01, item02, item03, item04, item05, item06, item07, item08, item09, item10}

	// // send item to bundler gateway
	// for _, item := range items {
	// 	resp, err := w.Client.SendItemToBundler(item.ItemBinary,"")
	// 	assert.NoError(t, err)
	// 	t.Log(resp.Id)
	// }
	resp, err := w.Client.BatchSendItemToBundler(items, "")
	assert.NoError(t, err)
	t.Log(resp)

	// assemble items
	arTxtags := []types.Tag{
		{Name: "GOAR", Value: "bundleTx"},
		{Name: "ACTION", Value: "test tx"},
	}
	bd, err := utils.NewBundle(items...)
	assert.NoError(t, err)

	txId, err := w.SendBundleTx(bd.BundleBinary, arTxtags)
	assert.NoError(t, err)
	t.Log(txId)
}

func TestVerifyBundleItem(t *testing.T) {
	cli := goar.NewClient("https://arweave.net")
	// id := "K0JskpURZ-zZ7m01txR7hArvsBDDi08S6-6YIVQoc_Y" // big size data
	// id := "mTm5-TtpsfJvUCPXflFe-P7HO6kOy4E2pGbt6-DUs40"

	// goar test tx
	// id := "ipVFFrAkLosTtk-M3J6wYq3MKpfE6zK75nMIC-oLVXw"
	// id := "2ZFhlTJlFbj8XVmBtnBHS-y6Clg68trcRgIKBNemTM8"
	// id := "WNGKdWsGqyhh7Y4vMcQL0GHFzNiyeqASIJn-Z1IjJE0"
	id := "lt24bnUGms5XLZeVamSPHePl4M2ClpLQyRxZI7weH1k"
	bd, err := cli.GetBundle(id)
	assert.NoError(t, err)
	for _, item := range bd.Items {
		err = utils.VerifyBundleItem(item)
		assert.NoError(t, err)
	}
}
