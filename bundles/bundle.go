/*
	js-lib:  https://github.com/ArweaveTeam/arweave-bundles
	ANS-102 format: https://github.com/ArweaveTeam/arweave-standards/blob/master/ans/ANS-102.md
*/

package bundles

import (
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/hamba/avro"
	slog "github.com/zyjblockchain/sandy_log/log"
	"math/big"
	"strconv"
)

var (
	log = slog.NewLog("bundles", slog.LevelDebug, false)
)

type BundleData struct {
	Items  []DataItem `json:"items"`
	binary []byte
}

type DataItem struct {
	SignatureType string      `json:"signatureType"`
	Signature     string      `json:"signature"`
	Owner         string      `json:"owner"`  //  utils.Base64Encode(wallet.PubKey.N.Bytes())
	Target        string      `json:"target"` // optional
	Anchor        string      `json:"anchor"` // optional
	Tags          []types.Tag `json:"tags"`
	Data          string      `json:"data"`
	Id            string      `json:"id"`

	binary []byte
}

func newDataItemJson(owner, signatureType, target, anchor string, data []byte, tags []types.Tag) (DataItem, error) {
	encTags := utils.TagsEncode(tags)
	dataItem := DataItem{
		SignatureType: signatureType,
		Signature:     "",
		Owner:         owner,
		Target:        target,
		Anchor:        anchor,
		Tags:          encTags,
		Data:          utils.Base64Encode(data),
		Id:            "",
		binary:        make([]byte, 0),
	}

	// verify tags
	if !VerifyEncodedTags(encTags) {
		return DataItem{}, errors.New("verify encoded tags failed")
	} else {
		return dataItem, nil
	}
}

func (d DataItem) getSignatureData() []byte {
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
	dataList = append(dataList, utils.Base64Encode([]byte(d.SignatureType)))
	dataList = append(dataList, d.Owner)
	dataList = append(dataList, d.Target)
	dataList = append(dataList, d.Anchor)
	dataList = append(dataList, tags)
	dataList = append(dataList, d.Data)

	hash := utils.DeepHash(dataList)
	deepHash := hash[:]
	return deepHash
}

func (d *DataItem) Sign(w *goar.Wallet) (DataItem, error) {
	// sign item
	signatureData := d.getSignatureData()
	signatureBytes, err := utils.Sign(signatureData, w.PrvKey)
	if err != nil {
		return DataItem{}, errors.New(fmt.Sprintf("signature error: %v", err))
	}
	id := sha256.Sum256(signatureBytes)
	d.Id = utils.Base64Encode(id[:])
	d.Signature = utils.Base64Encode(signatureBytes)
	return *d, nil
}

func (d *DataItem) AddTag(name, value string) {
	newTag := types.Tag{
		Name:  name,
		Value: value,
	}
	oldTags := d.Tags

	d.Tags = append(oldTags, utils.TagsEncode([]types.Tag{newTag})...)
}

func (d DataItem) Verify() bool {
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

func (d *DataItem) DecodeData() ([]byte, error) {
	return utils.Base64Decode(d.Data)
}

func (d DataItem) DecodeTag(tag types.Tag) (types.Tag, error) {
	tags, err := utils.TagsDecode([]types.Tag{tag})
	if err != nil || len(tags) == 0 {
		return types.Tag{}, errors.New(fmt.Sprintf("types.TagsDecode([]types.Tag{tag}) error: %v", err))
	} else {
		return tags[0], nil
	}
}

func (d DataItem) DecodeTagAt(index int) (types.Tag, error) {
	if len(d.Tags) < index-1 {
		return types.Tag{}, errors.New(fmt.Sprintf("Invalid index %d when tags array has %d tags", index, len(d.Tags)))
	}
	return d.DecodeTag(d.Tags[index])
}

func (d DataItem) UnpackTags() (map[string][]string, error) {
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

// func BundleDataItems(datas ...DataItemJson) (BundleData, error) {
// 	// verify
// 	for _, data := range datas {
// 		if !data.Verify() {
// 			return BundleData{}, errors.New("verify dataItemJson error")
// 		}
// 	}
// 	return BundleData{
// 		Items: datas,
// 	}, nil
// }

// func UnBundleDataItems(txData []byte) ([]DataItemJson, error) {
// 	bundleData := BundleData{}
// 	if err := json.Unmarshal(txData, &bundleData); err != nil {
// 		return nil, errors.New(fmt.Sprintf("json.Unmarshal(txData, &bundleData) error: %v", err))
// 	}
//
// 	itemsArray := bundleData.Items
//
// 	// verify
// 	for index, item := range itemsArray {
// 		if !item.Verify() {
// 			return nil, errors.New(fmt.Sprintf("verify item faild; item index: %d", index))
// 		}
// 	}
// 	return itemsArray, nil
// }

// ANS-104

// CreateDataItem This will create a single DataItem in binary format (Uint8Array)
func CreateDataItem(w *goar.Wallet, data []byte, owner []byte, signatureType int, target string, anchor string, tags []types.Tag) (di *DataItem, err error) {
	if len(owner) != 512 {
		return nil, errors.New("public key is not correct length")
	}
	targetBytes := []byte{}
	if target != "" {
		targetBytes, err = utils.Base64Decode(target)
		if err != nil {
			return
		}
	}

	targetLength := len(targetBytes) + 1

	anchorBytes := []byte{}
	if anchor != "" {
		anchorBytes, err = utils.Base64Decode(anchor)
		if err != nil {
			return
		}
	}
	anchorLength := len(anchorBytes) + 1

	tagsBytes, err := serializeTags(tags)
	if err != nil {
		return nil, err
	}
	tagsLength := 16 + len(tagsBytes)

	dataLenght := len(data)

	length := 2 + 512 + len(owner) + targetLength + anchorLength + tagsLength + dataLenght

	dataItemJs, err := newDataItemJson(utils.Base64Encode(owner), strconv.Itoa(signatureType), target, anchor, data, tags)
	if err != nil {
		return nil, err
	}
	// sign
	dataItemJs, err = dataItemJs.Sign(w)
	if err != nil {
		return nil, err
	}
	// Create array with set length
	bytesArr := make([]byte, 0, length)
	bytesArr = append(bytesArr, shortTo2ByteArray(signatureType)...)
	// Push bytes for `signature`
	sig, err := utils.Base64Decode(dataItemJs.Signature)
	if err != nil {
		return nil, err
	}
	bytesArr = append(bytesArr, sig...)
	// Push bytes for `owner`
	bytesArr = append(bytesArr, owner...)
	// Push `presence byte` and push `target` if present
	// 64 + OWNER_LENGTH
	if target != "" {
		bytesArr = append(bytesArr, byte(1))
		bytesArr = append(bytesArr, targetBytes...)
	} else {
		bytesArr = append(bytesArr, byte(0))
	}

	// Push `presence byte` and push `anchor` if present
	// 64 + OWNER_LENGTH
	if anchor != "" {
		bytesArr = append(bytesArr, byte(1))
		bytesArr = append(bytesArr, anchorBytes...)
	} else {
		bytesArr = append(bytesArr, byte(0))
	}

	// push tags
	bytesArr = append(bytesArr, longTo8ByteArray(len(tags))...)
	bytesArr = append(bytesArr, longTo8ByteArray(len(tagsBytes))...)

	if tags != nil {
		bytesArr = append(bytesArr, tagsBytes...)
	}

	// push data
	bytesArr = append(bytesArr, data...)
	dataItemJs.binary = bytesArr
	return &dataItemJs, nil
}

func BundleDataItems(dataItems ...DataItem) (*BundleData, error) {
	headers := make([]byte, 0) // length is 64 * len(dataItems)
	binaries := make([]byte, 0)

	for _, d := range dataItems {
		header := make([]byte, 0, 64)
		header = append(header, longTo32ByteArray(len(d.binary))...)
		id, err := utils.Base64Decode(d.Id)
		if err != nil {
			return nil, err
		}
		header = append(header, id...)

		headers = append(headers, header...)
		binaries = append(binaries, d.binary...)
	}

	bdBinary := make([]byte, 0)
	bdBinary = append(bdBinary, longTo32ByteArray(len(dataItems))...)
	bdBinary = append(bdBinary, headers...)
	bdBinary = append(bdBinary, binaries...)
	return &BundleData{
		Items:  dataItems,
		binary: bdBinary,
	}, nil
}

func (b *BundleData) SubmitBundleTx(w *goar.Wallet, tags []types.Tag) (txId string, err error) {
	bundleTags := []types.Tag{
		{Name: "Bundle-Format", Value: "binary"},
		{Name: "Bundle-Version", Value: "2.0.0"},
	}
	txTags := make([]types.Tag, 0)
	txTags = append(bundleTags, tags...)
	txId, err = w.SendDataSpeedUp(b.binary, txTags, 50)
	return
}

func longTo8ByteArray(long int) []byte {
	// we want to represent the input as a 8-bytes array
	byteArray := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < len(byteArray); i++ {
		byt := long & 0xff
		byteArray[i] = byte(byt)
		long = (long - byt) / 256
	}
	return byteArray
}

func shortTo2ByteArray(long int) []byte {
	byteArray := []byte{0, 0}
	for i := 0; i < len(byteArray); i++ {
		byt := long & 0xff
		byteArray[i] = byte(byt)
		long = (long - byt) / 256
	}
	return byteArray
}

func longTo32ByteArray(long int) []byte {
	byteArray := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < len(byteArray); i++ {
		byt := long & 0xff
		byteArray[i] = byte(byt)
		long = (long - byt) / 256
	}
	return byteArray
}

func serializeTags(tags []types.Tag) ([]byte, error) {
	if len(tags) == 0 {
		return make([]byte, 0), nil
	}

	tagParser, err := avro.Parse(`{
		  "type": "record",
		  "name": "Tag",
		  "fields": [
			{ "name": "name", "type": "string" },
			{ "name": "value", "type": "string" }
		  ]
		}`)
	if err != nil {
		return nil, err
	}

	tagsParser, err := avro.Parse(`{
		  "type": "array",
		  "items": ` + tagParser.String() + `
		}`)
	if err != nil {
		return nil, err
	}
	return avro.Marshal(tagsParser, tags)
}
