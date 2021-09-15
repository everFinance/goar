/*
	js-lib:  https://github.com/Bundler-Network/arbundles
	ANS-104 format: https://github.com/joshbenaron/arweave-standards/blob/ans104/ans/ANS-104.md
*/

package utils

import (
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/everFinance/goar/types"
	"strconv"
)

func NewBundle(dataItems ...types.BundleItem) (*types.Bundle, error) {
	headers := make([]byte, 0) // length is 64 * len(dataItems)
	binaries := make([]byte, 0)

	for _, d := range dataItems {
		header := make([]byte, 0, 64)
		header = append(header, LongTo32ByteArray(len(d.ItemBinary))...)
		id, err := Base64Decode(d.Id)
		if err != nil {
			return nil, err
		}
		header = append(header, id...)

		headers = append(headers, header...)
		binaries = append(binaries, d.ItemBinary...)
	}

	bdBinary := make([]byte, 0)
	bdBinary = append(bdBinary, LongTo32ByteArray(len(dataItems))...)
	bdBinary = append(bdBinary, headers...)
	bdBinary = append(bdBinary, binaries...)
	return &types.Bundle{
		Items:        dataItems,
		BundleBinary: bdBinary,
	}, nil
}

func DecodeBundle(bundleBinary []byte) (*types.Bundle, error) {
	// length must more than 32
	if len(bundleBinary) < 32 {
		return nil, errors.New("binary length must more than 32")
	}
	dataItemsNum := ByteArrayToLong(bundleBinary[:32])

	if len(bundleBinary) < 32+dataItemsNum*64 {
		return nil, errors.New("binary length incorrect")
	}

	bd := &types.Bundle{
		Items:        make([]types.BundleItem, 0),
		BundleBinary: bundleBinary,
	}
	dataItemStart := 32 + dataItemsNum*64
	for i := 0; i < dataItemsNum; i++ {
		headerBegin := 32 + i*64
		end := headerBegin + 64
		headerByte := bundleBinary[headerBegin:end]
		itemBinaryLength := ByteArrayToLong(headerByte[:32])
		id := Base64Encode(headerByte[32:64])

		dataItemBytes := bundleBinary[dataItemStart : dataItemStart+itemBinaryLength]
		dataItem, err := DecodeBundleItem(dataItemBytes)
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

func DecodeBundleItem(itemBinary []byte) (*types.BundleItem, error) {
	if len(itemBinary) < types.MIN_BUNDLE_BINARY_SIZE {
		return nil, errors.New("ItemBinary length incorrect")
	}
	sigType := ByteArrayToLong(itemBinary[:2])
	signature := Base64Encode(itemBinary[2:514])
	idhash := sha256.Sum256(itemBinary[2:514])
	id := Base64Encode(idhash[:])
	owner := Base64Encode(itemBinary[514:1026])
	target := ""
	anchor := ""
	tagsStart := 2 + 512 + 512 + 2
	anchorPresentByte := 1027
	targetPersent := itemBinary[1026] == 1
	if targetPersent {
		tagsStart += 32
		anchorPresentByte += 32 // 1059
		target = Base64Encode(itemBinary[1027 : 1027+32])
	}
	anchorPersent := itemBinary[anchorPresentByte] == 1
	if anchorPersent {
		tagsStart += 32
		anchor = Base64Encode(itemBinary[anchorPresentByte+1 : anchorPresentByte+1+32])
	}

	numOfTags := ByteArrayToLong(itemBinary[tagsStart : tagsStart+8])
	tagsBytesLength := ByteArrayToLong(itemBinary[tagsStart+8 : tagsStart+16])

	tags := []types.Tag{}
	if numOfTags > 0 {
		tagsBytes := itemBinary[tagsStart+16 : tagsStart+16+tagsBytesLength]
		// parser tags
		tgs, err := DeserializeTags(tagsBytes)
		if err != nil {
			return nil, err
		}
		tags = tgs
	}

	data := itemBinary[tagsStart+16+tagsBytesLength:]

	return &types.BundleItem{
		SignatureType: fmt.Sprintf("%d", sigType),
		Signature:     signature,
		Owner:         owner,
		Target:        target,
		Anchor:        anchor,
		Tags:          tags,
		Data:          Base64Encode(data),
		Id:            id,
		ItemBinary:    itemBinary,
	}, nil
}

func NewBundleItem(owner, signatureType, target, anchor string, data []byte, tags []types.Tag) *types.BundleItem {
	return &types.BundleItem{
		SignatureType: signatureType,
		Signature:     "",
		Owner:         owner,
		Target:        target,
		Anchor:        anchor,
		Tags:          tags,
		Data:          Base64Encode(data),
		Id:            "",
		ItemBinary:    make([]byte, 0),
	}
}

func SignBundleItem(d *types.BundleItem, prvKey *rsa.PrivateKey) error {
	// sign item
	signatureData, err := BundleItemSignData(*d)
	if err != nil {
		return err
	}
	signatureBytes, err := Sign(signatureData, prvKey)
	if err != nil {
		return errors.New(fmt.Sprintf("signature error: %v", err))
	}
	id := sha256.Sum256(signatureBytes)
	d.Id = Base64Encode(id[:])
	d.Signature = Base64Encode(signatureBytes)
	return nil
}

func BundleItemSignData(d types.BundleItem) ([]byte, error) {
	var err error
	tagsBy := make([]byte, 0)
	if len(d.ItemBinary) > 0 { // verify logic
		tagsBy = GetBundleItemTagsBytes(d.ItemBinary)
	} else {
		tagsBy, err = SerializeTags(d.Tags)
		if err != nil {
			return nil, err
		}
	}

	// deep hash
	dataList := make([]interface{}, 0)
	dataList = append(dataList, Base64Encode([]byte("dataitem")))
	dataList = append(dataList, Base64Encode([]byte("1")))
	dataList = append(dataList, Base64Encode([]byte(d.SignatureType)))
	dataList = append(dataList, d.Owner)
	dataList = append(dataList, d.Target)
	dataList = append(dataList, d.Anchor)
	dataList = append(dataList, Base64Encode(tagsBy))
	dataList = append(dataList, d.Data)

	hash := DeepHash(dataList)
	deepHash := hash[:]
	return deepHash, nil
}

func VerifyBundleItem(d types.BundleItem) error {
	// Get signature data and signature present in di.
	signatureData, err := BundleItemSignData(d)
	if err != nil {
		return fmt.Errorf("signatureData, err := d.GetSignatureData(); err : %v", err)
	}
	signatureBytes, err := Base64Decode(d.Signature)
	if err != nil {
		return fmt.Errorf("utils.Base64Decode(d.Signature) error: %v", err)
	}
	// Verify Id is correct
	idBytes := sha256.Sum256(signatureBytes)
	id := Base64Encode(idBytes[:])
	if id != d.Id {
		return fmt.Errorf("verify Id is not equal; id: %s, recId: %s", d.Id, id)
	}
	// Verify Signature is correct
	pubKey, err := OwnerToPubKey(d.Owner)
	if err != nil {
		return fmt.Errorf("utils.OwnerToPubKey(d.Owner), err: %v", err)
	}
	if err := Verify(signatureData, pubKey, signatureBytes); err != nil {
		return fmt.Errorf("utils.Verify(signatureData,pubKey,signatureBytes); err: %v", err)
	}
	return nil
}

func GetBundleItemTagsBytes(itemBinary []byte) []byte {
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

	numOfTags := ByteArrayToLong(itemBinary[tagsStart : tagsStart+8])
	tagsBytesLength := ByteArrayToLong(itemBinary[tagsStart+8 : tagsStart+16])

	if numOfTags > 0 {
		tagsBytes := itemBinary[tagsStart+16 : tagsStart+16+tagsBytesLength]
		return tagsBytes
	} else {
		return []byte{}
	}
}

func GenerateItemBinary(d *types.BundleItem) (err error) {
	if len(d.Signature) == 0 {
		return errors.New("must be sign")
	}

	targetBytes := []byte{}
	if d.Target != "" {
		targetBytes, err = Base64Decode(d.Target)
		if err != nil {
			return
		}
	}
	anchorBytes := []byte{}
	if d.Anchor != "" {
		anchorBytes, err = Base64Decode(d.Anchor)
		if err != nil {
			return
		}
	}
	tagsBytes, err := SerializeTags(d.Tags)
	if err != nil {
		return err
	}

	// Create array with set length
	bytesArr := make([]byte, 0, 1044)
	signType, err := strconv.Atoi(d.SignatureType)
	if err != nil {
		return err
	}
	bytesArr = append(bytesArr, ShortTo2ByteArray(signType)...)
	// Push bytes for `signature`
	sig, err := Base64Decode(d.Signature)
	if err != nil {
		return err
	}
	bytesArr = append(bytesArr, sig...)
	// Push bytes for `ownerByte`
	ownerByte, err := Base64Decode(d.Owner)
	if err != nil {
		return err
	}
	bytesArr = append(bytesArr, ownerByte...)
	// Push `presence byte` and push `target` if present
	// 64 + OWNER_LENGTH
	if d.Target != "" {
		bytesArr = append(bytesArr, byte(1))
		bytesArr = append(bytesArr, targetBytes...)
	} else {
		bytesArr = append(bytesArr, byte(0))
	}

	// Push `presence byte` and push `anchor` if present
	// 64 + OWNER_LENGTH
	if d.Anchor != "" {
		bytesArr = append(bytesArr, byte(1))
		bytesArr = append(bytesArr, anchorBytes...)
	} else {
		bytesArr = append(bytesArr, byte(0))
	}

	// push tags
	bytesArr = append(bytesArr, LongTo8ByteArray(len(d.Tags))...)
	bytesArr = append(bytesArr, LongTo8ByteArray(len(tagsBytes))...)

	if len(d.Tags) > 0 {
		bytesArr = append(bytesArr, tagsBytes...)
	}

	// push data
	data, err := Base64Decode(d.Data)
	if err != nil {
		return err
	}
	bytesArr = append(bytesArr, data...)
	d.ItemBinary = bytesArr
	return nil
}
