package bundles

import (
	"github.com/everFinance/goar/types"
	"github.com/hamba/avro"
)

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

	tagsParser, err := avro.Parse(`{"type": "array", "items": {"type": "record", "name": "Tag", "fields": [{"name": "name", "type": "string"}, {"name": "value", "type": "string"}]}}`)
	if err != nil {
		return nil, err
	}

	return avro.Marshal(tagsParser, tags)
}

func deserializeTags(data []byte) ([]types.Tag, error) {
	tagsParser, err := avro.Parse(`{"type": "array", "items": {"type": "record", "name": "Tag", "fields": [{"name": "name", "type": "string"}, {"name": "value", "type": "string"}]}}`)
	if err != nil {
		return nil, err
	}
	tags := make([]types.Tag, 0)
	err = avro.Unmarshal(tagsParser, data, &tags)
	return tags, err
}

func byteArrayToLong(b []byte) int {
	value := 0
	for i := len(b) - 1; i >= 0; i-- {
		value = value*256 + int(b[i])
	}
	return value
}
