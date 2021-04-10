/*
	js-lib:  https://github.com/ArweaveTeam/arweave-bundles
	ANS-102 format: https://github.com/ArweaveTeam/arweave-standards/blob/master/ans/ANS-102.md
*/

package bundles

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/everFinance/goar/wallet"
	slog "github.com/zyjblockchain/sandy_log/log"
	"math/big"
)

var log = slog.NewLog("bundles", slog.LevelDebug, false)

type BundleData struct {
	Items []DataItemJson `json:"items"`
}

type DataItemJson struct {
	Owner     string      `json:"owner"` //  utils.Base64Encode(wallet.PubKey.N.Bytes())
	Target    string      `json:"target"`
	Nonce     string      `json:"nonce"`
	Tags      []types.Tag `json:"tags"`
	Data      string      `json:"data"`
	Signature string      `json:"signature"`
	Id        string      `json:"id"`
}

func CreateDataItemJson(owner, target, nonce string, data []byte, tags []types.Tag) (DataItemJson, error) {
	encTags := types.TagsEncode(tags)
	dataItem := DataItemJson{
		Owner:     owner,
		Target:    target,
		Nonce:     nonce,
		Tags:      encTags,
		Data:      utils.Base64Encode(data),
		Signature: "",
		Id:        "",
	}

	// verify tags
	if !VerifyEncodedTags(encTags) {
		return DataItemJson{}, errors.New("verify encoded tags failed")
	} else {
		return dataItem, nil
	}
}

func (d DataItemJson) getSignatureData() []byte {
	tags := [][]string{}
	for _, tag := range d.Tags {
		tags = append(tags, []string{
			tag.Name, tag.Value,
		})
	}

	// deep hash
	dataList := make([]interface{}, 0)
	dataList = append(dataList, utils.Base64Encode([]byte("dataitem")))
	dataList = append(dataList, utils.Base64Encode([]byte("1")))
	dataList = append(dataList, d.Owner)
	dataList = append(dataList, d.Target)
	dataList = append(dataList, d.Nonce)
	dataList = append(dataList, tags)
	dataList = append(dataList, d.Data)

	hash := utils.DeepHash(dataList)
	deepHash := hash[:]
	return deepHash
}

func (d *DataItemJson) Sign(w *wallet.Wallet) (DataItemJson, error) {
	// sign item
	signatureData := d.getSignatureData()
	signatureBytes, err := utils.Sign(signatureData, w.PrvKey)
	if err != nil {
		return DataItemJson{}, errors.New(fmt.Sprintf("signature error: %v", err))
	}
	id := sha256.Sum256(signatureBytes)
	d.Id = utils.Base64Encode(id[:])
	d.Signature = utils.Base64Encode(signatureBytes)
	return *d, nil
}

func (d *DataItemJson) AddTag(name, value string) {
	newTag := types.Tag{
		Name:  name,
		Value: value,
	}
	oldTags := d.Tags

	d.Tags = append(oldTags, types.TagsEncode([]types.Tag{newTag})...)
}

func (d DataItemJson) Verify() bool {
	// Get signature data and signature present in di.
	signatureData := d.getSignatureData()
	signatureBytes, err := utils.Base64Decode(d.Signature)
	if err != nil {
		log.Error("utils.Base64Decode(d.Signature) error", "error", err)
		return false
	}
	// Verify Id is correct
	idBytes := sha256.Sum256(signatureBytes)
	if utils.Base64Encode(idBytes[:]) != d.Id {
		log.Error("verify Id is not equal")
		return false
	}

	// Verify Signature is correct
	owner, err := utils.Base64Decode(d.Owner)
	if err != nil {
		log.Error(" utils.Base64Decode(d.Owner) error", "error", err)
		return false
	}
	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(owner),
		E: 65537, //"AQAB"
	}
	if err := utils.Verify(signatureData, pubKey, signatureBytes); err != nil {
		log.Error("utils.Verify(signatureData,pubKey,signatureBytes)", "error", err)
		return false
	}

	// Verify tags array is valid.
	if !VerifyEncodedTags(d.Tags) {
		log.Error("VerifyEncodedTags(d.Tags) failed")
		return false
	}
	return true
}

func (d *DataItemJson) DecodeData() ([]byte, error) {
	return utils.Base64Decode(d.Data)
}

func (d DataItemJson) DecodeTag(tag types.Tag) (types.Tag, error) {
	tags, err := types.TagsDecode([]types.Tag{tag})
	if err != nil || len(tags) == 0 {
		return types.Tag{}, errors.New(fmt.Sprintf("types.TagsDecode([]types.Tag{tag}) error: %v", err))
	} else {
		return tags[0], nil
	}
}

func (d DataItemJson) DecodeTagAt(index int) (types.Tag, error) {
	if len(d.Tags) < index-1 {
		return types.Tag{}, errors.New(fmt.Sprintf("Invalid index %d when tags array has %d tags", index, len(d.Tags)))
	}
	return d.DecodeTag(d.Tags[index])
}

func (d DataItemJson) UnpackTags() (map[string][]string, error) {
	tagsMap := make(map[string][]string)
	for _, tag := range d.Tags {
		tt, err := d.DecodeTag(tag)
		if err != nil {
			return nil, err
		}

		name := tt.Name
		val := tt.Value
		if _, ok := tagsMap[name]; !ok {
			tagsMap[name] = make([]string, 0)
		}
		tagsMap[name] = append(tagsMap[name], val)
	}
	return tagsMap, nil
}

func (d DataItemJson) BundleData(datas ...DataItemJson) (BundleData, error) {
	// verify
	for _, data := range datas {
		if !data.Verify() {
			return BundleData{}, errors.New("verify dataItemJson error")
		}
	}
	return BundleData{
		Items: datas,
	}, nil
}

func (d DataItemJson) UnBundleData(txData []byte) ([]DataItemJson, error) {
	bundleData := BundleData{}
	if err := json.Unmarshal(txData, &bundleData); err != nil {
		return nil, errors.New(fmt.Sprintf("json.Unmarshal(txData, &bundleData) error: %v", err))
	}

	itemsArray := bundleData.Items

	// verify
	for index, item := range itemsArray {
		if !item.Verify() {
			return nil, errors.New(fmt.Sprintf("verify item faild; item index: %d", index))
		}
	}
	return itemsArray, nil
}
