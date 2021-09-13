package bundles

import (
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
)

type ArweaveBundles interface {
	Sign(w *goar.Wallet) (DataItem, error)
	AddTag(name, value string)
	Verify() bool
	DecodeData() ([]byte, error)
	DecodeTag(tag types.Tag) (types.Tag, error)
	DecodeTagAt(index int) (types.Tag, error)
	UnpackTags() (map[string][]string, error)
	BundleData(datas ...DataItem) (BundleData, error)
	UnBundleData(txData []byte) ([]DataItem, error)
}
