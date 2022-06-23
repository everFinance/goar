package example

import (
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBundle_SendBundleTx(t *testing.T) {
	// signer, err := goar.NewSignerFromPath("./testKey.json")
	// assert.NoError(t, err)
	//
	// itemSdk, err := goar.NewItemSigner(signer)
	// assert.NoError(t, err)
	//
	// tags := []types.Tag{
	// 	{Name: "Content-Type", Value: "image/jpeg"},
	// 	{Name: "App-Version", Value: "2.0.0"},
	// }
	//
	// data := []byte("123456")
	// item01, err := itemSdk.CreateAndSignItem(data, "", "", tags)
	// assert.NoError(t, err)
	//
	// err = utils.VerifyBundleItem(item01)
	// assert.NoError(t, err)

	// // send item to arseed
	// arseedUrl := "https://seed-dev.everpay.io"
	// resp, err := utils.SubmitItemToArSeed(item01,"USDT",arseedUrl)
	// assert.NoError(t, err)
	// t.Log(*resp)

	// // send item to bundlr network
	// bundlrUrl := "https://node1.bundlr.network"
	// resp, err := utils.SubmitItemToBundlr(item01, bundlrUrl)
	// assert.NoError(t, err)
	// t.Log(*resp)
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
