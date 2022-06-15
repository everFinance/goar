/*
	js-lib:  https://github.com/Bundler-Network/arbundles
	ANS-104 format: https://github.com/joshbenaron/arweave-standards/blob/ans104/ans/ANS-104.md
*/

package utils

import (
	"crypto/ed25519"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goether"
	"strconv"
)

func NewBundle(items ...types.BundleItem) (*types.Bundle, error) {
	headers := make([]byte, 0) // length is 64 * len(items)
	binaries := make([]byte, 0)

	for _, d := range items {
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
	bdBinary = append(bdBinary, LongTo32ByteArray(len(items))...)
	bdBinary = append(bdBinary, headers...)
	bdBinary = append(bdBinary, binaries...)
	return &types.Bundle{
		Items:        items,
		BundleBinary: bdBinary,
	}, nil
}

func DecodeBundle(bundleBinary []byte) (*types.Bundle, error) {
	// length must more than 32
	if len(bundleBinary) < 32 {
		return nil, errors.New("binary length must more than 32")
	}
	itemsNum := ByteArrayToLong(bundleBinary[:32])

	if len(bundleBinary) < 32+itemsNum*64 {
		return nil, errors.New("binary length incorrect")
	}

	bd := &types.Bundle{
		Items:        make([]types.BundleItem, 0),
		BundleBinary: bundleBinary,
	}
	bundleItemStart := 32 + itemsNum*64
	for i := 0; i < itemsNum; i++ {
		headerBegin := 32 + i*64
		end := headerBegin + 64
		headerByte := bundleBinary[headerBegin:end]
		itemBinaryLength := ByteArrayToLong(headerByte[:32])
		id := Base64Encode(headerByte[32:64])

		bundleItemBytes := bundleBinary[bundleItemStart : bundleItemStart+itemBinaryLength]
		bundleItem, err := DecodeBundleItem(bundleItemBytes)
		if err != nil {
			return nil, err
		}
		if bundleItem.Id != id {
			return nil, fmt.Errorf("bundleItem.Id != id, bundleItem.Id: %s, id: %s", bundleItem.Id, id)
		}
		bd.Items = append(bd.Items, *bundleItem)
		bundleItemStart += itemBinaryLength
	}
	return bd, nil
}

func DecodeBundleItem(itemBinary []byte) (*types.BundleItem, error) {
	sigType := ByteArrayToLong(itemBinary[:2])
	sigMeta, ok := types.SigConfigMap[sigType]
	if !ok {
		return nil, fmt.Errorf("not support sigType:%d", sigType)
	}
	sigLength := sigMeta.SigLength
	sigBy := itemBinary[2 : sigLength+2]
	signature := Base64Encode(sigBy)
	idhash := sha256.Sum256(sigBy)
	id := Base64Encode(idhash[:])

	ownerLength := sigMeta.PubLength
	owner := Base64Encode(itemBinary[sigLength+2 : sigLength+2+ownerLength])
	target := ""
	anchor := ""
	position := 2 + sigLength + ownerLength

	tagsStart := position + 2
	anchorPresentByte := position + 1

	targetPersent := itemBinary[position] == 1
	if targetPersent {
		tagsStart += 32
		anchorPresentByte += 32
		target = Base64Encode(itemBinary[position+1 : position+1+32])
	}
	anchorPersent := itemBinary[anchorPresentByte] == 1
	if anchorPersent {
		tagsStart += 32
		anchor = Base64Encode(itemBinary[anchorPresentByte+1 : anchorPresentByte+1+32])
	}

	numOfTags := ByteArrayToLong(itemBinary[tagsStart : tagsStart+8])

	var tagsBytesLength int
	tags := []types.Tag{}
	if numOfTags > 0 {
		tagsBytesLength = ByteArrayToLong(itemBinary[tagsStart+8 : tagsStart+16])
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
		SignatureType: sigType,
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

func NewBundleItem(owner string, signatureType int, target, anchor string, data []byte, tags []types.Tag) *types.BundleItem {
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

func BundleItemSignData(d types.BundleItem) ([]byte, error) {
	var err error
	tagsBy := make([]byte, 0)
	if len(d.ItemBinary) > 0 { // verify logic
		tagsBy, err = GetBundleItemTagsBytes(d.ItemBinary)
		if err != nil {
			return nil, err
		}
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
	dataList = append(dataList, Base64Encode([]byte(strconv.Itoa(d.SignatureType))))
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
	signMsg, err := BundleItemSignData(d)
	if err != nil {
		return fmt.Errorf("signMsg, err := d.GetSignatureData(); err : %v", err)
	}
	sign, err := Base64Decode(d.Signature)
	if err != nil {
		return fmt.Errorf("utils.Base64Decode(d.Signature) error: %v", err)
	}
	// Verify Id is correct
	idBytes := sha256.Sum256(sign)
	id := Base64Encode(idBytes[:])
	if id != d.Id {
		return fmt.Errorf("verify Id is not equal; id: %s, recId: %s", d.Id, id)
	}
	switch d.SignatureType {
	case types.ArweaveSignType:
		// Verify Signature is correct
		pubKey, err := OwnerToPubKey(d.Owner)
		if err != nil {
			return fmt.Errorf("utils.OwnerToPubKey(d.Owner), err: %v", err)
		}
		return Verify(signMsg, pubKey, sign)

	case types.ED25519SignType, types.SolanaSignType:
		pubkey, err := Base64Decode(d.Owner)
		if err != nil {
			return err
		}
		if !ed25519.Verify(pubkey, signMsg, sign) {
			return errors.New("verify ed25519 signature failed")
		}

	case types.EthereumSignType:
		signer, err := ItemSignerAddr(d)
		if err != nil {
			return err
		}

		addr, err := goether.Ecrecover(accounts.TextHash(signMsg), sign)
		if err != nil {
			return err
		}
		if signer != addr.String() {
			return errors.New("verify ecc sign failed")
		}
	default:
		return errors.New("not support the signType")
	}
	return nil
}

func GetBundleItemTagsBytes(itemBinary []byte) ([]byte, error) {
	sigType := ByteArrayToLong(itemBinary[:2])
	sigMeta, ok := types.SigConfigMap[sigType]
	if !ok {
		return nil, fmt.Errorf("not support sigType:%d", sigType)
	}
	sigLength := sigMeta.SigLength
	ownerLength := sigMeta.PubLength
	position := 2 + sigLength + ownerLength
	tagsStart := position + 2

	anchorPresentByte := position + 1
	if len(itemBinary) < anchorPresentByte {
		return nil, errors.New("itemBinary incorrect")
	}
	targetPersent := itemBinary[position] == 1
	if targetPersent {
		tagsStart += 32
		anchorPresentByte += 32
	}
	anchorPersent := itemBinary[anchorPresentByte] == 1
	if anchorPersent {
		tagsStart += 32
	}

	numOfTags := ByteArrayToLong(itemBinary[tagsStart : tagsStart+8])

	if numOfTags > 0 {
		tagsBytesLength := ByteArrayToLong(itemBinary[tagsStart+8 : tagsStart+16])
		tagsBytes := itemBinary[tagsStart+16 : tagsStart+16+tagsBytesLength]
		return tagsBytes, nil
	} else {
		return []byte{}, nil
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

	bytesArr = append(bytesArr, ShortTo2ByteArray(d.SignatureType)...)
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

func ItemSignerAddr(b types.BundleItem) (string, error) {
	switch b.SignatureType {
	case types.ArweaveSignType:
		return OwnerToAddress(b.Owner)

	case types.ED25519SignType, types.SolanaSignType:
		by, err := Base64Decode(b.Owner)
		if err != nil {
			return "", err
		}
		return base58.Encode(by), nil
	case types.EthereumSignType:
		pubkey, err := Base64Decode(b.Owner)
		if err != nil {
			return "", err
		}
		pk, err := crypto.UnmarshalPubkey(pubkey)
		if err != nil {
			err = fmt.Errorf("can not unmarshal pubkey: %v", err)
			return "", err
		}
		return crypto.PubkeyToAddress(*pk).String(), nil

	default:
		return "", errors.New("not support the signType")
	}
}
