package goar

import (
	"crypto/sha256"
	"errors"
	"io"

	"github.com/daqiancode/goar/types"
	"github.com/daqiancode/goar/utils"
	"github.com/everFinance/goether"
)

type ItemSigner struct {
	signType   int
	signer     interface{}
	owner      string // only rsa has owner
	signerAddr string
}

func NewItemSigner(signer interface{}) (*ItemSigner, error) {
	signType, signerAddr, owner, err := reflectSigner(signer)
	if err != nil {
		return nil, err
	}
	return &ItemSigner{
		signType:   signType,
		signer:     signer,
		owner:      owner,
		signerAddr: signerAddr,
	}, nil
}

func (i *ItemSigner) CreateAndSignItem(data []byte, target string, anchor string, tags []types.Tag) (types.BundleItem, error) {
	bundleItem, err := utils.NewBundleItem(i.owner, i.signType, target, anchor, data, tags)
	if err != nil {
		return types.BundleItem{}, err
	}
	// sign
	if err := SignBundleItem(i.signType, i.signer, bundleItem); err != nil {
		return types.BundleItem{}, err
	}
	if err := utils.GenerateItemBinary(bundleItem); err != nil {
		return types.BundleItem{}, err
	}
	return *bundleItem, nil
}

func (i *ItemSigner) CreateAndSignItemStream(data io.Reader, target string, anchor string, tags []types.Tag) (types.BundleItem, error) {
	bundleItem, err := utils.NewBundleItemStream(i.owner, i.signType, target, anchor, data, tags)
	if err != nil {
		return types.BundleItem{}, err
	}
	// sign
	if err := SignBundleItem(i.signType, i.signer, bundleItem); err != nil {
		return types.BundleItem{}, err
	}
	if err := utils.GenerateItemBinary(bundleItem); err != nil {
		return types.BundleItem{}, err
	}
	return *bundleItem, nil
}

func reflectSigner(signer interface{}) (signType int, signerAddr, owner string, err error) {
	if s, ok := signer.(*Signer); ok {
		signType = types.ArweaveSignType
		signerAddr = s.Address
		owner = s.Owner()
		return
	}
	if s, ok := signer.(*goether.Signer); ok {
		signType = types.EthereumSignType
		signerAddr = s.Address.String()
		owner = utils.Base64Encode(s.GetPublicKey())
		return
	}
	err = errors.New("not support this signer")
	return
}

func SignBundleItem(signatureType int, signer interface{}, item *types.BundleItem) error {
	signMsg, err := utils.BundleItemSignData(*item)
	if err != nil {
		return err
	}
	var sigData []byte
	switch signatureType {
	case types.ArweaveSignType:
		arSigner, ok := signer.(*Signer)
		if !ok {
			return errors.New("signer must be goar signer")
		}
		sigData, err = utils.Sign(signMsg, arSigner.PrvKey)
		if err != nil {
			return err
		}

	case types.EthereumSignType:
		ethSigner, ok := signer.(*goether.Signer)
		if !ok {
			return errors.New("signer not goether signer")
		}
		sigData, err = ethSigner.SignMsg(signMsg)
		if err != nil {
			return err
		}
	default:
		// todo come soon supprot ed25519
		return errors.New("not supprot this signType")
	}
	id := sha256.Sum256(sigData)
	item.Id = utils.Base64Encode(id[:])
	item.Signature = utils.Base64Encode(sigData)
	return nil
}
