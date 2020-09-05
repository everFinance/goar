package utils

import (
	"encoding/base64"

	"github.com/everFinance/goar/types"
)

func Base64Encode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func Base64Decode(data string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(data)
}

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
