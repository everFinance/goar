package example

import (
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func Test_SendData(t *testing.T) {
	arNode := "https://arweave.net"
	w, err := goar.NewWalletFromPath("./wallet/account1.json", arNode) // your wallet private key
	assert.NoError(t, err)

	data, err := ioutil.ReadFile("/Users/local/Downloads/abc.jpeg") // local file path
	if err != nil {
		panic(err)
	}
	tags := []types.Tag{
		{Name: "xxxx", Value: "sssss"},
		{Name: "yyyyyy", Value: "kkkkkk"},
	}
	tx, err := w.SendDataSpeedUp(data, tags, 10)
	assert.NoError(t, err)
	t.Logf("tx hash: %s", tx.ID)
}

func Test_LoadData(t *testing.T) {
	arCli := goar.NewClient("https://arweave.net")

	arId := "r90Z_PuhD-louq6uzLTI-xWMfB5TzIti30o7QvW-6A4"
	data, err := arCli.GetTransactionData(arId)
	assert.NoError(t, err)
	t.Log(len(data))
}

func TestDownloadDataStream(t *testing.T) {
	arCli := goar.NewClient("https://arweave.net")

	arId := "BF85hzl9HobCkLKrKET1MRd2pr_XRqB2dAWQEZYDTRE" // 300KB
	// arId := "cqCdSEKu-A272DuwFpKPBdyEsxXHT92gxoorS3Y-sbM" // image size:12MB
	dataFile, err := arCli.DownloadChunkDataStream(arId)
	assert.NoError(t, err)
	dataFile.Close()
}

func TestConcurrentDownloadStream(t *testing.T) {
	arCli := goar.NewClient("https://arweave.net")

	arId := "cqCdSEKu-A272DuwFpKPBdyEsxXHT92gxoorS3Y-sbM"
	dataFile, data, err := arCli.ConcurrentDownloadChunkDataStream(arId, 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(data))
	dataFile.Close()
}
