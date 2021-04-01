package types

import "github.com/everFinance/goar/utils"

type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func TagsEncode(tags []Tag) []Tag {
	base64Tags := []Tag{}

	for _, tag := range tags {
		base64Tags = append(base64Tags, Tag{
			Name:  utils.Base64Encode([]byte(tag.Name)),
			Value: utils.Base64Encode([]byte(tag.Value)),
		})
	}

	return base64Tags
}

func TagsDecode(base64Tags []Tag) ([]Tag, error) {
	tags := []Tag{}

	for _, bt := range base64Tags {
		bName, err := utils.Base64Decode(bt.Name)
		if err != nil {
			return nil, err
		}

		bValue, err := utils.Base64Decode(bt.Value)
		if err != nil {
			return nil, err
		}

		tags = append(tags, Tag{
			Name:  string(bName),
			Value: string(bValue),
		})
	}

	return tags, nil
}
