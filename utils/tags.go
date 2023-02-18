package utils

import (
	"github.com/daqiancode/goar/types"
	"github.com/hamba/avro"
)

func TagsEncode(tags []types.Tag) []types.Tag {
	base64Tags := []types.Tag{}

	for _, tag := range tags {
		base64Tags = append(base64Tags, types.Tag{
			Name:  Base64Encode([]byte(tag.Name)),
			Value: Base64Encode([]byte(tag.Value)),
		})
	}

	return base64Tags
}

func TagsDecode(base64Tags []types.Tag) ([]types.Tag, error) {
	tags := []types.Tag{}

	for _, bt := range base64Tags {
		bName, err := Base64Decode(bt.Name)
		if err != nil {
			return nil, err
		}

		bValue, err := Base64Decode(bt.Value)
		if err != nil {
			return nil, err
		}

		tags = append(tags, types.Tag{
			Name:  string(bName),
			Value: string(bValue),
		})
	}

	return tags, nil
}

// using bundle tx, avro serialize
func SerializeTags(tags []types.Tag) ([]byte, error) {
	if len(tags) == 0 {
		return make([]byte, 0), nil
	}

	tagsParser, err := avro.Parse(`{"type": "array", "items": {"type": "record", "name": "Tag", "fields": [{"name": "name", "type": "string"}, {"name": "value", "type": "string"}]}}`)
	if err != nil {
		return nil, err
	}

	return avro.Marshal(tagsParser, tags)
}

func DeserializeTags(data []byte) ([]types.Tag, error) {
	tagsParser, err := avro.Parse(`{"type": "array", "items": {"type": "record", "name": "Tag", "fields": [{"name": "name", "type": "string"}, {"name": "value", "type": "string"}]}}`)
	if err != nil {
		return nil, err
	}
	tags := make([]types.Tag, 0)
	err = avro.Unmarshal(tagsParser, data, &tags)
	return tags, err
}
