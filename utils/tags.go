package utils

import (
	"github.com/everVision/goar/schema"
	"github.com/linkedin/goavro/v2"
)

func TagsEncode(tags []schema.Tag) []schema.Tag {
	base64Tags := []schema.Tag{}

	for _, tag := range tags {
		base64Tags = append(base64Tags, schema.Tag{
			Name:  Base64Encode([]byte(tag.Name)),
			Value: Base64Encode([]byte(tag.Value)),
		})
	}

	return base64Tags
}

func TagsDecode(base64Tags []schema.Tag) ([]schema.Tag, error) {
	tags := []schema.Tag{}

	for _, bt := range base64Tags {
		bName, err := Base64Decode(bt.Name)
		if err != nil {
			return nil, err
		}

		bValue, err := Base64Decode(bt.Value)
		if err != nil {
			return nil, err
		}

		tags = append(tags, schema.Tag{
			Name:  string(bName),
			Value: string(bValue),
		})
	}

	return tags, nil
}

// using bundle tx, avro serialize
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

func SerializeTags(tags []schema.Tag) ([]byte, error) {
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

func DeserializeTags(data []byte) ([]schema.Tag, error) {
	codec, err := goavro.NewCodec(avroTagSchema)
	if err != nil {
		return nil, err
	}

	avroTags, _, err := codec.NativeFromBinary(data)
	if err != nil {
		return nil, err
	}

	tags := []schema.Tag{}

	for _, v := range avroTags.([]interface{}) {
		tag := v.(map[string]interface{})
		tags = append(tags, schema.Tag{Name: string(tag["name"].([]byte)), Value: string(tag["value"].([]byte))})
	}
	return tags, err
}
