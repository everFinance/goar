package goar

import (
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
)

const (
	BUNDLER         = "http://bundler.arweave.net:10000"
	MIN_BINARY_SIZE = 1044
)

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

func newDataItem(owner, signatureType, target, anchor string, data []byte, tags []types.Tag) (*DataItem, error) {
	dataItem := &DataItem{
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
	var err error
	tagsBy := make([]byte,0)
	if len(d.itemBinary) > 0 { // verify logic
		tagsBy = getTagsBytes(d.itemBinary)
	} else {
		tagsBy, err = utils.SerializeTags(d.Tags)
		if err != nil {
			return nil, err
		}
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

func (d *DataItem) Sign(prvKey  *rsa.PrivateKey) error {
	// sign item
	signatureData, err := d.getSignatureData()
	if err != nil {
		return  err
	}
	signatureBytes, err := utils.Sign(signatureData, prvKey)
	if err != nil {
		return errors.New(fmt.Sprintf("signature error: %v", err))
	}
	id := sha256.Sum256(signatureBytes)
	d.Id = utils.Base64Encode(id[:])
	d.Signature = utils.Base64Encode(signatureBytes)
	return nil
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

func getTagsBytes(itemBinary []byte) []byte {
	tagsStart := 2 + 512 + 512 + 2
	anchorPresentByte := 1027
	if len(itemBinary) < anchorPresentByte {
		return []byte{}
	}
		targetPersent := itemBinary[1026] == 1
	if targetPersent {
		tagsStart += 32
		anchorPresentByte += 32 // 1059
	}
	anchorPersent := itemBinary[anchorPresentByte] == 1
	if anchorPersent {
		tagsStart += 32
	}

	numOfTags := utils.ByteArrayToLong(itemBinary[tagsStart : tagsStart+8])
	tagsBytesLength := utils.ByteArrayToLong(itemBinary[tagsStart+8 : tagsStart+16])

	if numOfTags > 0 {
		tagsBytes := itemBinary[tagsStart+16 : tagsStart+16+tagsBytesLength]
		return tagsBytes
	} else {
		return []byte{}
	}
}

func recoverDataItem(itemBinary []byte) (*DataItem, error) {
	if len(itemBinary) < MIN_BINARY_SIZE {
		return nil, errors.New("itemBinary length incorrect")
	}
	sigType := utils.ByteArrayToLong(itemBinary[:2])
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

	numOfTags := utils.ByteArrayToLong(itemBinary[tagsStart : tagsStart+8])
	tagsBytesLength := utils.ByteArrayToLong(itemBinary[tagsStart+8 : tagsStart+16])

	tags := []types.Tag{}
	if numOfTags > 0 {
		tagsBytes := itemBinary[tagsStart+16 : tagsStart+16+tagsBytesLength]
		// parser tags
		tgs, err := utils.DeserializeTags(tagsBytes)
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



