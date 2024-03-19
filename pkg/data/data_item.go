package data

import (
	"encoding/base64"
	"encoding/binary"
	"errors"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/pkg/crypto"
)

const (
	MAX_TAGS             = 128
	MAX_TAG_KEY_LENGTH   = 1024
	MAX_TAG_VALUE_LENGTH = 3072
)

func NewDataItem(rawData []byte, s *goar.Signer, target string, anchor string, tags []Tag) (*DataItem, error) {
	rawOwner := []byte(s.PubKey.N.Bytes())
	rawTarget, err := base64.RawURLEncoding.DecodeString(target)
	if err != nil {
		return nil, err
	}
	rawAnchor := []byte(anchor)

	rawTags, err := decodeTags(&tags)
	if err != nil {
		return nil, err
	}

	chunks := []interface{}{
		[]byte("dataitem"),
		[]byte("1"),
		[]byte("1"),
		rawOwner,
		rawTarget,
		rawAnchor,
		rawTags,
		rawData,
	}
	signatureData := crypto.DeepHash(chunks)

	rawSignature, err := s.SignMsg(signatureData[:])
	if err != nil {
		return nil, err
	}
	raw := make([]byte, 0)
	raw = binary.LittleEndian.AppendUint16(raw, uint16(1))
	raw = append(raw, rawSignature...)
	raw = append(raw, rawOwner...)

	if target == "" {
		raw = append(raw, 0)
	} else {
		raw = append(raw, 1)
	}
	raw = append(raw, rawTarget...)

	if anchor == "" {
		raw = append(raw, 0)
	} else {
		raw = append(raw, 1)
	}
	raw = append(raw, rawAnchor...)
	numberOfTags := make([]byte, 8)
	binary.LittleEndian.PutUint16(numberOfTags, uint16(len(tags)))
	raw = append(raw, numberOfTags...)

	tagsLength := make([]byte, 8)
	binary.LittleEndian.PutUint16(tagsLength, uint16(len(rawTags)))
	raw = append(raw, tagsLength...)
	raw = append(raw, rawTags...)
	raw = append(raw, rawData...)
	rawID, err := crypto.SHA256(rawSignature)
	if err != nil {
		return nil, err
	}
	data := base64.RawURLEncoding.EncodeToString(rawData)
	return &DataItem{
		SignatureType: 1,
		Signature:     base64.RawURLEncoding.EncodeToString(rawSignature),
		ID:            base64.RawURLEncoding.EncodeToString(rawID),
		Owner:         s.Owner(),
		Target:        target,
		Anchor:        anchor,
		Tags:          tags,
		Data:          data,
		Raw:           raw,
	}, nil
}

// Decode a DataItem from bytes
func DecodeDataItem(raw []byte) (*DataItem, error) {
	N := len(raw)
	if N < 2 {
		return nil, errors.New("binary too small")
	}

	signatureType, signatureLength, publicKeyLength, err := getSignatureMetadata(raw[:2])
	if err != nil {
		return nil, err
	}

	signatureStart := 2
	signatureEnd := signatureLength + signatureStart
	signature := base64.RawURLEncoding.EncodeToString(raw[signatureStart:signatureEnd])
	rawId, err := crypto.SHA256(raw[signatureStart:signatureEnd])
	if err != nil {
		return nil, err
	}
	id := base64.RawURLEncoding.EncodeToString(rawId)
	ownerStart := signatureEnd
	ownerEnd := ownerStart + publicKeyLength
	owner := base64.RawURLEncoding.EncodeToString(raw[ownerStart:ownerEnd])

	position := ownerEnd
	target, position := getTarget(&raw, position)
	anchor, position := getAnchor(&raw, position)
	tags, position, err := encodeTags(&raw, position)
	if err != nil {
		return nil, err
	}
	data := base64.RawURLEncoding.EncodeToString(raw[position:])

	return &DataItem{
		ID:            id,
		SignatureType: signatureType,
		Signature:     signature,
		Owner:         owner,
		Target:        target,
		Anchor:        anchor,
		Tags:          *tags,
		Data:          data,
		Raw:           raw,
	}, nil
}

func VerifyDataItem(dataItem *DataItem) (bool, error) {

	// Verify ID
	rawSignature, err := base64.RawURLEncoding.DecodeString(dataItem.Signature)
	if err != nil {
		return false, err
	}
	rawId, err := crypto.SHA256(rawSignature)
	if err != nil {
		return false, err
	}
	id := base64.RawURLEncoding.EncodeToString(rawId)
	if id != dataItem.ID {
		return false, errors.New("invalid data item - signature and id don't match")
	}

	// Verify Signature Owner
	rawOwner, err := base64.RawURLEncoding.DecodeString(dataItem.Owner)
	if err != nil {
		return false, err
	}
	rawTarget, err := base64.RawURLEncoding.DecodeString(dataItem.Target)
	if err != nil {
		return false, err
	}
	rawAnchor := []byte(dataItem.Anchor)
	rawTags, err := decodeTags(&dataItem.Tags)
	if err != nil {
		return false, err
	}
	rawData, err := base64.RawURLEncoding.DecodeString(dataItem.Data)
	if err != nil {
		return false, err
	}
	chunks := []interface{}{
		[]byte("dataitem"),
		[]byte("1"),
		[]byte("1"),
		rawOwner,
		rawTarget,
		rawAnchor,
		rawTags,
		rawData,
	}
	signatureData := crypto.DeepHash(chunks)

	valid, err := verify(signatureData[:], rawSignature, dataItem.Owner)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, errors.New("invalid data item - signature failed to verify")
	}

	// VERIFY TAGS
	if len(dataItem.Tags) > MAX_TAGS {
		return false, errors.New("invalid data item - tags cannot be more than 128")
	}

	for _, tag := range dataItem.Tags {
		if len([]byte(tag.Name)) == 0 || len([]byte(tag.Name)) > MAX_TAG_KEY_LENGTH {
			return false, errors.New("invalid data item - tag key too long")
		}
		if len([]byte(tag.Value)) == 0 || len([]byte(tag.Value)) > MAX_TAG_VALUE_LENGTH {
			return false, errors.New("invalid data item - tag value too long")
		}
	}

	if len([]byte(dataItem.Anchor)) > 32 {
		return false, errors.New("invalid data item - anchor should be 32 bytes")
	}
	return true, nil
}
