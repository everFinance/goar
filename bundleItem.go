package goar

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/everFinance/goether"
	"gopkg.in/h2non/gentleman.v2"
	"net/http"
)

type ItemSdk struct {
	signType   int
	signer     interface{}
	owner      string // only rsa has owner
	signerAddr string
	cli        *gentleman.Client
}

func NewItemSdk(signer interface{}, serverUrl string) (*ItemSdk, error) {
	signType, signerAddr, owner, err := reflectSigner(signer)
	if err != nil {
		return nil, err
	}
	return &ItemSdk{
		signType:   signType,
		signer:     signer,
		owner:      owner,
		signerAddr: signerAddr,
		cli:        gentleman.New().URL(serverUrl),
	}, nil
}

func (i *ItemSdk) CreateAndSignItem(data []byte, target string, anchor string, tags []types.Tag) (types.BundleItem, error) {
	bundleItem := utils.NewBundleItem(i.owner, i.signType, target, anchor, data, tags)
	// sign
	if err := SignBundleItem(i.signType, i.signer, bundleItem); err != nil {
		return types.BundleItem{}, err
	}
	if err := utils.GenerateItemBinary(bundleItem); err != nil {
		return types.BundleItem{}, err
	}
	return *bundleItem, nil
}

func (i *ItemSdk) SubmitItem(item types.BundleItem, currency string) (*types.BundlerResp, error) {
	req := i.cli.Request()
	req.Method("POST")
	req.Path(fmt.Sprintf("/bundle/tx/%s", currency))
	req.SetHeader("content-type", "application/octet-stream")

	itemBinary := item.ItemBinary
	if len(itemBinary) == 0 {
		if err := utils.GenerateItemBinary(&item); err != nil {
			return nil, err
		}
		itemBinary = item.ItemBinary
	}
	req.Body(bytes.NewReader(itemBinary))

	resp, err := req.Send()
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("send to bundler request failed; http code: %d, errMsg:%s", resp.StatusCode, resp.String())
	}
	br := &types.BundlerResp{}
	err = resp.JSON(br)
	return br, err
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
