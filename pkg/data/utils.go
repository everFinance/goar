package data

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/linkedin/goavro/v2"
)

const avroTagSchema = `
{
	"type": "array",
	"items": {
		"type": "record",
		"name": "Tag",
		"fields": [
			{ "name": "name", "type": "bytes" },
			{ "name": "value", "type": "bytes" }
		]
	}
}`

func getSignatureMetadata(data []byte) (SignatureType int, SignatureLength int, PublicKeyLength int, err error) {
	SignatureType = int(binary.LittleEndian.Uint16(data))
	signatureMeta, ok := SignatureConfig[SignatureType]
	if !ok {
		return -1, -1, -1, fmt.Errorf("unsupported signature type:%d", SignatureType)
	}
	SignatureLength = signatureMeta.SignatureLength
	PublicKeyLength = signatureMeta.PublicKeyLength
	err = nil
	return
}

func getTarget(data *[]byte, position int) (string, int) {
	target := ""
	if (*data)[position] == 1 {
		target = base64.RawURLEncoding.EncodeToString((*data)[position+1 : position+1+32])
		position += 32
	}
	return target, position + 1
}

func getAnchor(data *[]byte, position int) (string, int) {
	anchor := ""
	if (*data)[position] == 1 {
		anchor = string((*data)[position+1 : position+1+32])
		position += 32
	}
	return anchor, position + 1
}

func encodeTags(data *[]byte, startAt int) (*[]Tag, int, error) {
	tags := &[]Tag{}
	tagsEnd := startAt + 8 + 8
	numberOfTags := int(binary.LittleEndian.Uint16((*data)[startAt : startAt+8]))
	numberOfTagBytesStart := startAt + 8
	numberOfTagBytesEnd := numberOfTagBytesStart + 8
	numberOfTagBytes := int(binary.LittleEndian.Uint16((*data)[numberOfTagBytesStart:numberOfTagBytesEnd]))
	if numberOfTags > 127 {
		return tags, tagsEnd, errors.New("invalid data item - max tags 127")
	}
	if numberOfTags > 0 && numberOfTagBytes > 0 {
		bytesDataStart := numberOfTagBytesEnd
		bytesDataEnd := numberOfTagBytesEnd + numberOfTagBytes
		bytesData := (*data)[bytesDataStart:bytesDataEnd]

		tags, err := decodeAvro(bytesData)
		if err != nil {
			return nil, tagsEnd, err
		}
		tagsEnd = bytesDataEnd
		return tags, tagsEnd, nil
	}
	return tags, tagsEnd, nil
}

func decodeAvro(data []byte) (*[]Tag, error) {
	codec, err := goavro.NewCodec(avroTagSchema)
	if err != nil {
		return nil, err
	}

	avroTags, _, err := codec.NativeFromBinary(data)
	if err != nil {
		return nil, err
	}

	tags := &[]Tag{}

	for _, v := range avroTags.([]interface{}) {
		tag := v.(map[string]any)
		*tags = append(*tags, Tag{Name: string(tag["name"].([]byte)), Value: string(tag["value"].([]byte))})
	}
	return tags, err
}

func decodeTags(tags *[]Tag) ([]byte, error) {
	if len(*tags) > 0 {
		data, err := encodeAvro(tags)
		if err != nil {
			return nil, err
		}

		return data, nil
	}
	return nil, nil
}

func encodeAvro(tags *[]Tag) ([]byte, error) {
	codec, err := goavro.NewCodec(avroTagSchema)
	if err != nil {
		return nil, err
	}

	avroTags := []map[string]any{}

	for _, tag := range *tags {
		m := map[string]any{"name": []byte(tag.Name), "value": []byte(tag.Value)}
		avroTags = append(avroTags, m)
	}
	data, err := codec.BinaryFromNative(nil, avroTags)
	if err != nil {
		return nil, err
	}

	return data, err
}

func generateBundleHeader(d *[]DataItem) (*[]BundleHeader, error) {
	headers := []BundleHeader{}

	for _, dataItem := range *d {
		idBytes, err := base64.RawURLEncoding.DecodeString(dataItem.ID)
		if err != nil {
			return nil, err
		}

		id := int(binary.LittleEndian.Uint16(idBytes))
		size := len(dataItem.Raw)
		raw := make([]byte, 64)
		binary.LittleEndian.PutUint16(raw, uint16(size))
		binary.LittleEndian.AppendUint16(raw, uint16(id))
		headers = append(headers, BundleHeader{id: id, size: size, raw: raw})
	}
	return &headers, nil
}

func decodeBundleHeader(data *[]byte) (*[]BundleHeader, int) {
	N := int(binary.LittleEndian.Uint32((*data)[:32]))
	headers := []BundleHeader{}
	for i := 32; i < 32+64*N; i += 64 {
		size := int(binary.LittleEndian.Uint16((*data)[i : i+32]))
		id := int(binary.LittleEndian.Uint16((*data)[i+32 : i+64]))
		headers = append(headers, BundleHeader{id: id, size: size})
	}
	return &headers, N
}
