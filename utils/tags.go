package utils

import (
	"github.com/everFinance/goar/types"
	"github.com/hamba/avro"
	"github.com/linkedin/goavro/v2"
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
func SerializeTags1(tags []types.Tag) ([]byte, error) {
	if len(tags) == 0 {
		return make([]byte, 0), nil
	}

	tagsParser, err := avro.Parse(`{"type": "array", "items": {"type": "record", "name": "Tag", "fields": [{"name": "name", "type": "string"}, {"name": "value", "type": "string"}]}}`)
	if err != nil {
		return nil, err
	}

	return avro.Marshal(tagsParser, tags)
}

func DeserializeTags1(data []byte) ([]types.Tag, error) {
	tagsParser, err := avro.Parse(`{"type": "array", "items": {"type": "record", "name": "Tag", "fields": [{"name": "name", "type": "string"}, {"name": "value", "type": "string"}]}}`)
	if err != nil {
		return nil, err
	}
	tags := make([]types.Tag, 0)
	err = avro.Unmarshal(tagsParser, data, &tags)
	return tags, err
}

const avroTagSchema = `{
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

func SerializeTags(tags []types.Tag) ([]byte, error) {
	if len(tags) == 0 {
		return make([]byte, 0), nil
	}

	codec, err := goavro.NewCodec(avroTagSchema)
	if err != nil {
		return nil, err
	}
	avroTags := []map[string]interface{}{}
	for _, tag := range tags {
		m := map[string]interface{}{"name": []byte(tag.Name), "value": []byte(tag.Value)}
		avroTags = append(avroTags, m)
	}

	data, err := codec.BinaryFromNative(nil, avroTags)
	if err != nil {
		return nil, err
	}

	return data, err
}

func DeserializeTags(data []byte) ([]types.Tag, error) {
	codec, err := goavro.NewCodec(avroTagSchema)
	if err != nil {
		return nil, err
	}

	avroTags, _, err := codec.NativeFromBinary(data)
	if err != nil {
		return nil, err
	}

	tags := []types.Tag{}

	for _, v := range avroTags.([]interface{}) {
		tag := v.(map[string]interface{})
		tags = append(tags, types.Tag{Name: string(tag["name"].([]byte)), Value: string(tag["value"].([]byte))})
	}
	return tags, err
}
