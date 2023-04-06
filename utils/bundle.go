/*
	js-lib:  https://github.com/Bundler-Network/arbundles
	ANS-104 format: https://github.com/joshbenaron/arweave-standards/blob/ans104/ans/ANS-104.md
*/

package utils

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/everFinance/arseeding/schema"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goether"
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
		if len(id) != 32 {
			return nil, errors.New("item id length must 32")
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

// it's caller's responsibility to delete tmp file after handle bundleData

func NewBundleStream(items ...types.BundleItem) (*types.Bundle, error) {
	headers := make([]byte, 0) // length is 64 * len(items)
	headers = append(headers, LongTo32ByteArray(len(items))...)
	dataReader, err := os.CreateTemp(".", "bundleData-")
	if err != nil {
		return nil, err
	}
	for _, d := range items {
		header := make([]byte, 0, 64)
		if d.DataReader == nil {
			return nil, errors.New("NewBundleStream method dataReader can't be null")
		}
		itemInfo, err := d.DataReader.Stat()
		if err != nil {
			return nil, err
		}
		metaBy, err := generateItemMetaBinary(&d)
		if err != nil {
			return nil, err
		}
		itemBinaryLen := len(metaBy) + int(itemInfo.Size())
		header = append(header, LongTo32ByteArray(itemBinaryLen)...)
		id, err := Base64Decode(d.Id)
		if err != nil {
			return nil, err
		}
		if len(id) != 32 {
			return nil, errors.New("item id length must 32")
		}
		header = append(header, id...)
		headers = append(headers, header...)
	}
	_, err = io.Copy(dataReader, bytes.NewBuffer(headers))
	if err != nil {
		return nil, err
	}
	for _, d := range items {
		_, err = d.DataReader.Seek(0, 0)
		if err != nil {
			return nil, err
		}
		binaryReader, err := GenerateItemBinaryStream(&d)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(dataReader, binaryReader)
		if err != nil {
			return nil, err
		}
		_, err = d.DataReader.Seek(0, 0)
		if err != nil {
			return nil, err
		}
	}
	_, err = dataReader.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	b := &types.Bundle{
		Items:            items,
		BundleBinary:     make([]byte, 0),
		BundleDataReader: dataReader,
	}
	return b, nil
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
		if len(bundleBinary) < end {
			return nil, errors.New("binary length incorrect")
		}
		headerByte := bundleBinary[headerBegin:end]
		itemBinaryLength := ByteArrayToLong(headerByte[:32])
		id := Base64Encode(headerByte[32:64])
		if len(bundleBinary) < bundleItemStart+itemBinaryLength || itemBinaryLength < 0 {
			return nil, errors.New("binary length incorrect")
		}
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

// it's caller's responsibility to delete all tmp file after handle all bundle item

func DecodeBundleStream(bundleData *os.File) (*types.Bundle, error) {
	// length must more than 32
	itemsNumBy := make([]byte, 32, 32)
	n, err := bundleData.Read(itemsNumBy)
	if n < 32 || err != nil {
		return nil, errors.New("binary length must more than 32")
	}
	itemsNum := ByteArrayToLong(itemsNumBy)
	bd := &types.Bundle{
		Items: make([]types.BundleItem, 0),
	}
	bundleItemStart := 32 + itemsNum*64
	for i := 0; i < itemsNum; i++ {
		headerBegin := 32 + i*64
		headerByte := make([]byte, 64, 64)
		n, err = bundleData.ReadAt(headerByte, int64(headerBegin))
		if n < 64 || err != nil {
			return nil, errors.New("binary length incorrect")
		}
		itemBinaryLength := ByteArrayToLong(headerByte[:32])
		id := Base64Encode(headerByte[32:64])
		itemReader, err := os.CreateTemp(".", "bundleItem-")
		if err != nil {
			return nil, errors.New("CreateTempItemFile error")
		}
		_, err = bundleData.Seek(int64(bundleItemStart), 0)
		if err != nil {
			return nil, errors.New("seek bundleData failed")
		}
		n, err1 := io.CopyN(itemReader, bundleData, int64(itemBinaryLength))
		if int(n) < itemBinaryLength || err1 != nil {
			return nil, errors.New("binary length incorrect")
		}
		_, err = itemReader.Seek(0, 0)
		if err != nil {
			return nil, errors.New("seek itemData failed")
		}
		bundleItem, err2 := DecodeBundleItemStream(itemReader)
		itemReader.Close()
		os.Remove(itemReader.Name())
		if err2 != nil {
			return nil, err2
		}

		if bundleItem.Id != id {
			return nil, fmt.Errorf("bundleItem.Id != id, bundleItem.Id: %s, id: %s", bundleItem.Id, id)
		}
		bd.Items = append(bd.Items, *bundleItem)
		bundleItemStart += itemBinaryLength
	}
	_, err = bundleData.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	bd.BundleDataReader = bundleData
	return bd, nil
}

func DecodeBundleItem(itemBinary []byte) (*types.BundleItem, error) {
	if len(itemBinary) < 2 {
		return nil, errors.New("itemBinary incorrect")
	}
	sigType := ByteArrayToLong(itemBinary[:2])
	sigMeta, ok := types.SigConfigMap[sigType]
	if !ok {
		return nil, fmt.Errorf("not support sigType:%d", sigType)
	}
	sigLength := sigMeta.SigLength
	if len(itemBinary) < sigLength+2 {
		return nil, errors.New("itemBinary incorrect")
	}
	sigBy := itemBinary[2 : sigLength+2]
	signature := Base64Encode(sigBy)
	idhash := sha256.Sum256(sigBy)
	id := Base64Encode(idhash[:])

	ownerLength := sigMeta.PubLength
	if len(itemBinary) < sigLength+2+ownerLength {
		return nil, errors.New("itemBinary incorrect")
	}
	owner := Base64Encode(itemBinary[sigLength+2 : sigLength+2+ownerLength])
	target := ""
	anchor := ""
	position := 2 + sigLength + ownerLength

	tagsStart := position + 2
	anchorPresentByte := position + 1
	if len(itemBinary) < position {
		return nil, errors.New("itemBinary incorrect")
	}
	targetPersent := itemBinary[position] == 1
	if targetPersent {
		tagsStart += 32
		anchorPresentByte += 32
		if len(itemBinary) < position+1+32 {
			return nil, errors.New("itemBinary incorrect")
		}
		target = Base64Encode(itemBinary[position+1 : position+1+32])
	}
	if len(itemBinary) < anchorPresentByte {
		return nil, errors.New("itemBinary incorrect")
	}
	anchorPersent := itemBinary[anchorPresentByte] == 1
	if anchorPersent {
		tagsStart += 32
		if len(itemBinary) < anchorPresentByte+1+32 {
			return nil, errors.New("itemBinary incorrect")
		}
		anchor = Base64Encode(itemBinary[anchorPresentByte+1 : anchorPresentByte+1+32])
	}

	numOfTags := ByteArrayToLong(itemBinary[tagsStart : tagsStart+8])

	var tagsBytesLength int
	tags := []types.Tag{}
	tagsBytes := make([]byte, 0)
	if numOfTags > 0 {
		if len(itemBinary) < tagsStart+16 {
			return nil, errors.New("itemBinary incorrect")
		}
		tagsBytesLength = ByteArrayToLong(itemBinary[tagsStart+8 : tagsStart+16])
		if len(itemBinary) < tagsStart+16+tagsBytesLength || tagsStart+16+tagsBytesLength < 0 {
			return nil, errors.New("itemBinary incorrect")
		}
		tagsBytes = itemBinary[tagsStart+16 : tagsStart+16+tagsBytesLength]
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
		TagsBy:        Base64Encode(tagsBytes),
		ItemBinary:    itemBinary,
	}, nil
}

func DecodeBundleItemStream(itemBinary io.Reader) (*types.BundleItem, error) {
	sigTypeBy := make([]byte, 2, 2)
	n, err := itemBinary.Read(sigTypeBy)
	if err != nil || n < 2 {
		return nil, errors.New("itemBinary incorrect")
	}
	sigType := ByteArrayToLong(sigTypeBy)
	sigMeta, ok := types.SigConfigMap[sigType]
	if !ok {
		return nil, fmt.Errorf("not support sigType:%d", sigType)
	}
	sigLength := sigMeta.SigLength
	sigBy := make([]byte, sigLength, sigLength)
	n, err = itemBinary.Read(sigBy)
	if err != nil || n < sigLength {
		return nil, errors.New("itemBinary incorrect")
	}
	signature := Base64Encode(sigBy)
	idhash := sha256.Sum256(sigBy)
	id := Base64Encode(idhash[:])

	ownerLength := sigMeta.PubLength
	ownerBy := make([]byte, ownerLength, ownerLength)
	n, err = itemBinary.Read(ownerBy)
	if err != nil || n < ownerLength {
		return nil, errors.New("itemBinary incorrect")
	}
	owner := Base64Encode(ownerBy)
	target := ""
	anchor := ""

	targetPresentByte := make([]byte, 1, 1)
	n, err = itemBinary.Read(targetPresentByte)
	if err != nil || n < 1 {
		return nil, errors.New("itemBinary incorrect")
	}
	if targetPresentByte[0] == 1 {
		targetBy := make([]byte, 32, 32)
		n, err = itemBinary.Read(targetBy)
		if err != nil || n < 32 {
			return nil, errors.New("itemBinary incorrect")
		}
		target = Base64Encode(targetBy)
	}

	anchorPresentByte := make([]byte, 1, 1)
	n, err = itemBinary.Read(anchorPresentByte)
	if err != nil || n < 1 {
		return nil, errors.New("itemBinary incorrect")
	}
	if anchorPresentByte[0] == 1 {
		anchorBy := make([]byte, 32, 32)
		n, err = itemBinary.Read(anchorBy)
		if err != nil || n < 32 {
			return nil, errors.New("itemBinary incorrect")
		}
		anchor = Base64Encode(anchorBy)
	}

	numOfTagsBy := make([]byte, 8, 8)
	n, err = itemBinary.Read(numOfTagsBy)
	if err != nil || n < 8 {
		return nil, errors.New("itemBinary incorrect")
	}
	numOfTags := ByteArrayToLong(numOfTagsBy)

	tagsBytesLengthBy := make([]byte, 8, 8)
	n, err = itemBinary.Read(tagsBytesLengthBy)
	if err != nil || n < 8 {
		return nil, errors.New("itemBinary incorrect")
	}
	tagsBytesLength := ByteArrayToLong(tagsBytesLengthBy)

	tags := []types.Tag{}
	tagsBytes := make([]byte, 0)
	if numOfTags > 0 {
		tagsBytes = make([]byte, tagsBytesLength, tagsBytesLength)
		n, err = itemBinary.Read(tagsBytes)
		if err != nil || n < tagsBytesLength {
			return nil, errors.New("itemBinary incorrect")
		}
		// parser tags
		tgs, err := DeserializeTags(tagsBytes)
		if err != nil {
			return nil, err
		}
		tags = tgs
	}
	dataReader, err := os.CreateTemp(".", "itemData-")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(dataReader, itemBinary)
	if err != nil {
		os.Remove(dataReader.Name())
		return nil, err
	}
	_, err = dataReader.Seek(0, 0)
	if err != nil {
		os.Remove(dataReader.Name())
		return nil, err
	}
	return &types.BundleItem{
		SignatureType: sigType,
		Signature:     signature,
		Owner:         owner,
		Target:        target,
		Anchor:        anchor,
		Tags:          tags,
		Data:          "",
		Id:            id,
		TagsBy:        Base64Encode(tagsBytes),
		ItemBinary:    make([]byte, 0),
		DataReader:    dataReader,
	}, nil
}

func NewBundleItemStream(owner string, signatureType int, target, anchor string, data io.Reader, tags []types.Tag) (*types.BundleItem, error) {
	return newBundleItem(owner, signatureType, target, anchor, data, tags)
}

func NewBundleItem(owner string, signatureType int, target, anchor string, data []byte, tags []types.Tag) (*types.BundleItem, error) {
	return newBundleItem(owner, signatureType, target, anchor, data, tags)
}

func newBundleItem(owner string, signatureType int, target, anchor string, data interface{}, tags []types.Tag) (*types.BundleItem, error) {
	if target != "" {
		targetBy, err := Base64Decode(target)
		if err != nil {
			return nil, err
		}
		if len(targetBy) != 32 {
			return nil, errors.New("taget length must be 32")
		}
	}
	if anchor != "" {
		anchorBy, err := Base64Decode(anchor)
		if err != nil {
			return nil, err
		}
		if len(anchorBy) != 32 {
			return nil, errors.New("anchor length must be 32")
		}
	}
	tagsBytes, err := SerializeTags(tags)
	if err != nil {
		return nil, err
	}
	item := &types.BundleItem{
		SignatureType: signatureType,
		Signature:     "",
		Owner:         owner,
		Target:        target,
		Anchor:        anchor,
		Tags:          tags,
		Id:            "",
		TagsBy:        Base64Encode(tagsBytes),
		ItemBinary:    make([]byte, 0),
	}
	if _, ok := data.(*os.File); ok {
		item.DataReader = data.(*os.File)
	} else if _, ok = data.([]byte); ok {
		item.Data = Base64Encode(data.([]byte))
	}
	return item, nil
}

func BundleItemSignData(d types.BundleItem) ([]byte, error) {
	// deep hash
	dataList := make([]interface{}, 0)
	dataList = append(dataList, Base64Encode([]byte("dataitem")))
	dataList = append(dataList, Base64Encode([]byte("1")))
	dataList = append(dataList, Base64Encode([]byte(strconv.Itoa(d.SignatureType))))
	dataList = append(dataList, d.Owner)
	dataList = append(dataList, d.Target)
	dataList = append(dataList, d.Anchor)
	dataList = append(dataList, d.TagsBy)
	if d.DataReader != nil {
		dataList = append(dataList, d.DataReader)
	} else {
		dataList = append(dataList, d.Data)
	}

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
	if len(itemBinary) < 2 {
		return nil, errors.New("itemBinary incorrect")
	}

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
	if len(itemBinary) < position {
		return nil, errors.New("itemBinary incorrect")
	}
	targetPersent := itemBinary[position] == 1
	if targetPersent {
		tagsStart += 32
		anchorPresentByte += 32
	}
	if len(itemBinary) < anchorPresentByte {
		return nil, errors.New("itemBinary incorrect")
	}
	anchorPersent := itemBinary[anchorPresentByte] == 1
	if anchorPersent {
		tagsStart += 32
	}

	if len(itemBinary) < tagsStart+8 {
		return nil, errors.New("itemBinary incorrect")
	}
	numOfTags := ByteArrayToLong(itemBinary[tagsStart : tagsStart+8])

	if numOfTags > 0 {
		if len(itemBinary) < tagsStart+16 {
			return nil, errors.New("itemBinary incorrect")
		}
		tagsBytesLength := ByteArrayToLong(itemBinary[tagsStart+8 : tagsStart+16])
		if len(itemBinary) < tagsStart+16+tagsBytesLength || tagsStart+16+tagsBytesLength < 0 {
			return nil, errors.New("itemBinary incorrect")
		}
		tagsBytes := itemBinary[tagsStart+16 : tagsStart+16+tagsBytesLength]
		return tagsBytes, nil
	} else {
		return []byte{}, nil
	}
}

func generateItemMetaBinary(d *types.BundleItem) ([]byte, error) {
	if len(d.Signature) == 0 {
		return nil, errors.New("must be sign")
	}

	var err error
	targetBytes := []byte{}
	if d.Target != "" {
		targetBytes, err = Base64Decode(d.Target)
		if err != nil {
			return nil, err
		}
		if len(targetBytes) != 32 {
			return nil, errors.New("targetBytes length must 32")
		}
	}
	anchorBytes := []byte{}
	if d.Anchor != "" {
		anchorBytes, err = Base64Decode(d.Anchor)
		if err != nil {
			return nil, err
		}
		if len(anchorBytes) != 32 {
			return nil, errors.New("anchorBytes length must 32")
		}
	}
	tagsBytes := make([]byte, 0)
	if len(d.Tags) > 0 {
		tagsBytes, err = Base64Decode(d.TagsBy)
		if err != nil {
			return nil, err
		}
	}

	sigMeta, ok := types.SigConfigMap[d.SignatureType]
	if !ok {
		return nil, fmt.Errorf("not support sigType:%d", d.SignatureType)
	}

	sigLength := sigMeta.SigLength
	ownerLength := sigMeta.PubLength

	// Create array with set length
	bytesArr := make([]byte, 0, 2+sigLength+ownerLength)

	bytesArr = append(bytesArr, ShortTo2ByteArray(d.SignatureType)...)
	// Push bytes for `signature`
	sig, err := Base64Decode(d.Signature)
	if err != nil {
		return nil, err
	}

	if len(sig) != sigLength {
		return nil, errors.New("signature length incorrect")
	}

	bytesArr = append(bytesArr, sig...)
	// Push bytes for `ownerByte`
	ownerByte, err := Base64Decode(d.Owner)
	if err != nil {
		return nil, err
	}
	if len(ownerByte) != ownerLength {
		return nil, errors.New("signature length incorrect")
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
	return bytesArr, nil
}

func GenerateItemBinary(d *types.BundleItem) (by []byte, err error) {
	metaBinary, err := generateItemMetaBinary(d)
	if err != nil {
		return nil, err
	}

	by = append(by, metaBinary...)
	// push data
	data := make([]byte, 0)
	if len(d.Data) > 0 {
		data, err = Base64Decode(d.Data)
		if err != nil {
			return nil, err
		}
		by = append(by, data...)
	}
	return
}

func GenerateItemBinaryStream(d *types.BundleItem) (binaryReader io.Reader, err error) {
	metaBinary, err := generateItemMetaBinary(d)
	if err != nil {
		return nil, err
	}

	metaBuf := bytes.NewBuffer(metaBinary)
	if d.DataReader == nil {
		return metaBuf, nil
	} else {
		_, err = d.DataReader.Seek(0, 0)
		if err != nil {
			return nil, err
		}
		// note: DataReader must seek(0,0) after call DataReader.read(), otherwise BinaryReader will change
		return io.MultiReader(metaBuf, d.DataReader), nil
	}
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

func SubmitItemToBundlr(item types.BundleItem, bundlrUrl string) (*types.BundlrResp, error) {
	itemBinary := item.ItemBinary
	if len(itemBinary) == 0 {
		var err error
		itemBinary, err = GenerateItemBinary(&item)
		if err != nil {
			return nil, err
		}
	}
	// post to bundler
	resp, err := http.DefaultClient.Post(bundlrUrl+"/tx", "application/octet-stream", bytes.NewReader(itemBinary))
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
	br := &types.BundlrResp{}
	if err := json.Unmarshal(body, br); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(body,br) failed; err: %v", err)
	}
	return br, nil
}

func SubmitItemToArSeed(item types.BundleItem, currency, arseedUrl string) (*schema.RespOrder, error) {
	itemBinary := item.ItemBinary
	if len(itemBinary) == 0 {
		var err error
		itemBinary, err = GenerateItemBinary(&item)
		if err != nil {
			return nil, err
		}
		itemBinary = item.ItemBinary
	}
	resp, err := http.DefaultClient.Post(arseedUrl+"/bundle/tx/"+currency, "application/octet-stream", bytes.NewReader(itemBinary))
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
	br := &schema.RespOrder{}
	if err := json.Unmarshal(body, br); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(body,br) failed; err: %v", err)
	}
	return br, nil
}
