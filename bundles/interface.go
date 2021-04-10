package bundles

import (
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/wallet"
)

type ArweaveBundles interface {
	Sign(w *wallet.Wallet) (DataItemJson, error)
	AddTag(name, value string)
	Verify() bool
	DecodeData() ([]byte, error)
	DecodeTag(tag types.Tag) (types.Tag, error)
	DecodeTagAt(index int) (types.Tag, error)
	UnpackTags() (map[string][]string, error)
	BundleData(datas ...DataItemJson) (BundleData, error)
	UnBundleData(txData []byte) ([]DataItemJson, error)
}
