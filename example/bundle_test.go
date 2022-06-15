package example

import (
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBundle_SendBundleTx(t *testing.T) {
	signer, err := goar.NewSignerFromPath("./testKey.json")
	assert.NoError(t, err)
	arseedUrl := "http://127.0.0.1:8080"

	itemSdk, err := goar.NewItemSdk(signer, arseedUrl)
	assert.NoError(t, err)

	target := "Fkj5J8CDLC9Jif4CzgtbiXJBnwXLSrp5AaIllleH_yY"
	tags := []types.Tag{
		{Name: "goar", Value: "bundleTx"},
		{Name: "App-Version", Value: "2.0.0"},
	}

	item01, err := itemSdk.CreateAndSignItem([]byte("goar bundle tx 01"), target, "", tags)
	assert.NoError(t, err)
	item02, err := itemSdk.CreateAndSignItem([]byte("goar bundle tx 02"), target, "", tags)
	assert.NoError(t, err)
	item03, err := itemSdk.CreateAndSignItem([]byte("goar bundle tx 03"), target, "", tags)
	assert.NoError(t, err)
	item04, err := itemSdk.CreateAndSignItem([]byte("goar bundle tx 04"), target, "", tags)
	assert.NoError(t, err)
	item05, err := itemSdk.CreateAndSignItem([]byte("goar bundle tx 05"), target, "", tags)
	assert.NoError(t, err)
	item06, err := itemSdk.CreateAndSignItem([]byte("goar bundle tx 06"), target, "", tags)
	assert.NoError(t, err)
	item07, err := itemSdk.CreateAndSignItem([]byte("goar bundle tx 07"), target, "", tags)
	assert.NoError(t, err)
	item08, err := itemSdk.CreateAndSignItem([]byte("goar bundle tx 08"), target, "", tags)
	assert.NoError(t, err)
	item09, err := itemSdk.CreateAndSignItem([]byte("goar bundle tx 09"), target, "", tags)
	assert.NoError(t, err)
	item10, err := itemSdk.CreateAndSignItem([]byte("goar bundle tx 10"), target, "", tags)
	assert.NoError(t, err)

	items := []types.BundleItem{item01, item02, item03, item04, item05, item06, item07, item08, item09, item10}

	// send item to arseed gateway
	for _, item := range items {
		resp, err := itemSdk.SubmitItem(item, "USDC")
		assert.NoError(t, err)
		t.Log(resp)
	}
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
	data, err := cli.DownloadChunkData(id)
	assert.NoError(t, err)
	bd, err := utils.DecodeBundle(data)
	assert.NoError(t, err)
	for _, item := range bd.Items {
		err = utils.VerifyBundleItem(item)
		assert.NoError(t, err)
	}
}
