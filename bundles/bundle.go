package bundles

import (
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/wallet"
)

type DataItemJson struct {
	Owner     string
	Target    string
	Nonce     string
	Tags      []types.Tag
	Data      string
	Signature string
	Id        string
}

type BundleData struct {
	Item []DataItemJson
}

type ArweaveBundles interface {
	CreateData(owner, target, nonce, data string, tags []types.Tag) (DataItemJson, error)
	Sign(item DataItemJson, w wallet.Wallet) (DataItemJson, error)
	AddTag(item *DataItemJson, name, value string)
	Verify(item DataItemJson) bool
	DecodeData(item DataItemJson) []byte
	DecodeTag(tag types.Tag) []byte
	DecodeTagAt(item DataItemJson, index int) []byte
	UnpackTags(item DataItemJson) []string
	BundleData(data []DataItemJson) BundleData
	UnBundleData(txData []byte) []DataItemJson
}
