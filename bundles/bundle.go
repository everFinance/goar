/*
	js-lib:  https://github.com/ArweaveTeam/arweave-bundles
	ANS-102 format: https://github.com/ArweaveTeam/arweave-standards/blob/master/ans/ANS-102.md
*/

package bundles

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	slog "github.com/zyjblockchain/sandy_log/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	log = slog.NewLog("bundles", slog.LevelDebug, false)
)

const (
	BUNDLER         = "http://bundler.arweave.net:10000"
	MIN_BINARY_SIZE = 1044
)

type BundleData struct {
	Items        []DataItem `json:"items"`
	bundleBinary []byte
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

	itemBinary []byte
}

type BundlerResp struct {
	Id        string `json:"id"`
	Signature string `json:"signature"`
	N         string `json:"n"`
}

func newDataItem(owner, signatureType, target, anchor string, data []byte, tags []types.Tag) (DataItem, error) {
	dataItem := DataItem{
		SignatureType: signatureType,
		Signature:     "",
		Owner:         owner,
		Target:        target,
		Anchor:        anchor,
		Tags:          tags,
		Data:          utils.Base64Encode(data),
		Id:            "",
		itemBinary:    make([]byte, 0),
	}
	return dataItem, nil
}

func (d DataItem) getSignatureData() ([]byte, error) {
	tagsBy, err := serializeTags(d.Tags)
	if err != nil {
		return nil, err
	}
	// deep hash
	dataList := make([]interface{}, 0)
	dataList = append(dataList, utils.Base64Encode([]byte("dataitem")))
	dataList = append(dataList, utils.Base64Encode([]byte("1")))
	dataList = append(dataList, utils.Base64Encode([]byte(d.SignatureType)))
	dataList = append(dataList, d.Owner)
	dataList = append(dataList, d.Target)
	dataList = append(dataList, d.Anchor)
	dataList = append(dataList, utils.Base64Encode(tagsBy))
	dataList = append(dataList, d.Data)

	hash := utils.DeepHash(dataList)
	deepHash := hash[:]
	return deepHash, nil
}

func (d *DataItem) Sign(w *goar.Wallet) (DataItem, error) {
	// sign item
	signatureData, err := d.getSignatureData()
	if err != nil {
		return DataItem{}, err
	}
	fmt.Printf("signData: %s", hex.EncodeToString(signatureData))
	signatureBytes, err := utils.Sign(signatureData, w.PrvKey)
	if err != nil {
		return DataItem{}, errors.New(fmt.Sprintf("signature error: %v", err))
	}
	id := sha256.Sum256(signatureBytes)
	d.Id = utils.Base64Encode(id[:])
	d.Signature = utils.Base64Encode(signatureBytes)
	return *d, nil
}

func (d DataItem) Verify() error {
	// Get signature data and signature present in di.
	signatureData, err := d.getSignatureData()
	if err != nil {
		return fmt.Errorf("signatureData, err := d.getSignatureData(); err : %v", err)
	}
	signatureBytes, err := utils.Base64Decode(d.Signature)
	if err != nil {
		return fmt.Errorf("utils.Base64Decode(d.Signature) error: %v", err)
	}
	// Verify Id is correct
	idBytes := sha256.Sum256(signatureBytes)
	id := utils.Base64Encode(idBytes[:])
	if id != d.Id {
		return fmt.Errorf("verify Id is not equal; id: %s, recId: %s", d.Id, id)
	}
	// Verify Signature is correct
	pubKey, err := utils.OwnerToPubKey(d.Owner)
	if err != nil {
		return fmt.Errorf("utils.OwnerToPubKey(d.Owner), err: %v", err)
	}
	if err := utils.Verify(signatureData, pubKey, signatureBytes); err != nil {
		return fmt.Errorf("utils.Verify(signatureData,pubKey,signatureBytes); err: %v", err)
	}
	return nil
}

// ANS-104

// CreateDataItem This will create a single DataItem in bundleBinary format (Uint8Array)
func CreateDataItem(w *goar.Wallet, data []byte, owner []byte, signatureType int, target string, anchor string, tags []types.Tag) (di DataItem, err error) {
	if len(owner) != 512 {
		return di, errors.New("public key is not correct length")
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
		return di, err
	}
	tagsLength := 16 + len(tagsBytes)

	dataLenght := len(data)

	length := 2 + 512 + len(owner) + targetLength + anchorLength + tagsLength + dataLenght

	dataItemJs, err := newDataItem(utils.Base64Encode(owner), strconv.Itoa(signatureType), target, anchor, data, tags)
	if err != nil {
		return di, err
	}
	// sign
	dataItemJs, err = dataItemJs.Sign(w)
	if err != nil {
		return di, err
	}
	// Create array with set length
	bytesArr := make([]byte, 0, length)
	bytesArr = append(bytesArr, shortTo2ByteArray(signatureType)...)
	// Push bytes for `signature`
	sig, err := utils.Base64Decode(dataItemJs.Signature)
	if err != nil {
		return di, err
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
	dataItemJs.itemBinary = bytesArr
	return dataItemJs, nil
}

// SendToBundler send bundle dataItem to bundler gateway
func (d DataItem) SendToBundler() (*BundlerResp, error) {
	// post to bundler
	resp, err := http.DefaultClient.Post(BUNDLER+"/tx", "application/octet-stream", bytes.NewReader(d.itemBinary))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("send to bundler request failed; http code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	// json unmarshal
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll(resp.Body) error: %v", err)
	}
	br := &BundlerResp{}
	if err := json.Unmarshal(body, br); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(body,br) failed; err: %v", err)
	}
	return br, nil
}

func BundleDataItems(dataItems ...DataItem) (*BundleData, error) {
	headers := make([]byte, 0) // length is 64 * len(dataItems)
	binaries := make([]byte, 0)

	for _, d := range dataItems {
		header := make([]byte, 0, 64)
		header = append(header, longTo32ByteArray(len(d.itemBinary))...)
		id, err := utils.Base64Decode(d.Id)
		if err != nil {
			return nil, err
		}
		header = append(header, id...)

		headers = append(headers, header...)
		binaries = append(binaries, d.itemBinary...)
	}

	bdBinary := make([]byte, 0)
	bdBinary = append(bdBinary, longTo32ByteArray(len(dataItems))...)
	bdBinary = append(bdBinary, headers...)
	bdBinary = append(bdBinary, binaries...)
	return &BundleData{
		Items:        dataItems,
		bundleBinary: bdBinary,
	}, nil
}

func (b *BundleData) SubmitBundleTx(w *goar.Wallet, tags []types.Tag, txSpeed int64) (txId string, err error) {
	bundleTags := []types.Tag{
		{Name: "Bundle-Format", Value: "binary"},
		{Name: "Bundle-Version", Value: "2.0.0"},
	}
	txTags := make([]types.Tag, 0)
	txTags = append(bundleTags, tags...)
	txId, err = w.SendDataSpeedUp(b.bundleBinary, txTags, txSpeed)
	return
}

func RecoverBundleData(bundleBinary []byte) (*BundleData, error) {
	// length must more than 32
	if len(bundleBinary) < 32 {
		return nil, errors.New("binary length must more than 32")
	}
	dataItemsNum := byteArrayToLong(bundleBinary[:32])

	if len(bundleBinary) < 32+dataItemsNum*64 {
		return nil, errors.New("binary length incorrect")
	}

	bd := &BundleData{
		Items:        make([]DataItem, 0),
		bundleBinary: bundleBinary,
	}
	dataItemStart := 32 + dataItemsNum*64
	for i := 0; i < dataItemsNum; i++ {
		headerBegin := 32 + i*64
		end := headerBegin + 64
		headerByte := bundleBinary[headerBegin:end]
		itemBinaryLength := byteArrayToLong(headerByte[:32])
		id := utils.Base64Encode(headerByte[32:64])

		dataItemBytes := bundleBinary[dataItemStart : dataItemStart+itemBinaryLength]
		dataItem, err := recoverDataItem(dataItemBytes)
		if err != nil {
			return nil, err
		}
		if dataItem.Id != id {
			return nil, fmt.Errorf("dataItem.Id != id, dataItem.Id: %s, id: %s", dataItem.Id, id)
		}
		bd.Items = append(bd.Items, *dataItem)
		dataItemStart += itemBinaryLength
	}
	return bd, nil
}

func recoverDataItem(itemBinary []byte) (*DataItem, error) {
	if len(itemBinary) < MIN_BINARY_SIZE {
		return nil, errors.New("itemBinary length incorrect")
	}
	sigType := byteArrayToLong(itemBinary[:2])
	signature := utils.Base64Encode(itemBinary[2:514])
	idhash := sha256.Sum256(itemBinary[2:514])
	id := utils.Base64Encode(idhash[:])
	owner := utils.Base64Encode(itemBinary[514:1026])
	target := ""
	anchor := ""
	tagsStart := 2 + 512 + 512 + 2
	anchorPresentByte := 1027
	targetPersent := itemBinary[1026] == 1
	if targetPersent {
		tagsStart += 32
		anchorPresentByte += 32 // 1059
		target = utils.Base64Encode(itemBinary[1027 : 1027+32])
	}
	anchorPersent := itemBinary[anchorPresentByte] == 1
	if anchorPersent {
		tagsStart += 32
		anchor = utils.Base64Encode(itemBinary[anchorPresentByte+1 : anchorPresentByte+1+32])
	}

	numOfTags := byteArrayToLong(itemBinary[tagsStart : tagsStart+8])
	tagsBytesLength := byteArrayToLong(itemBinary[tagsStart+8 : tagsStart+16])

	tags := []types.Tag{}
	if numOfTags > 0 {
		tagsBytes := itemBinary[tagsStart+16 : tagsStart+16+tagsBytesLength]
		// parser tags
		tgs, err := deserializeTags(tagsBytes)
		if err != nil {
			return nil, err
		}
		tags = tgs
	}

	data := itemBinary[tagsStart+16+tagsBytesLength:]

	return &DataItem{
		SignatureType: fmt.Sprintf("%d", sigType),
		Signature:     signature,
		Owner:         owner,
		Target:        target,
		Anchor:        anchor,
		Tags:          tags,
		Data:          utils.Base64Encode(data),
		Id:            id,
		itemBinary:    itemBinary,
	}, nil
}
