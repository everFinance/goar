package utils

import (
	"encoding/hex"
	"testing"

	"github.com/everFinance/goar/types"
	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	tagsBase64 := []types.Tag{
		types.Tag{
			Name:  "QXBwLU5hbWU",
			Value: "U21hcnRXZWF2ZUFjdGlvbg",
		},
		types.Tag{
			Name:  "SW5wdXQ",
			Value: "eyJmdW5jdGlvbiI6InRyYW5zZmVyIiwicXR5Ijo1MDAsInRhcmdldCI6Ilp5aGhBTHdxazhuMnVyV1Y0RTNqSEJjNzd3YWE1RnItcUhscl9jdGlIQk0ifQ",
		},
	}
	tags := []types.Tag{
		types.Tag{
			Name:  "App-Name",
			Value: "SmartWeaveAction",
		},
		types.Tag{
			Name:  "Input",
			Value: `{"function":"transfer","qty":500,"target":"ZyhhALwqk8n2urWV4E3jHBc77waa5Fr-qHlr_ctiHBM"}`,
		},
	}

	tagsRes, err := TagsDecode(tagsBase64)
	assert.NoError(t, err)
	assert.Equal(t, tags[0].Name, tagsRes[0].Name)
	assert.Equal(t, tags[0].Value, tagsRes[0].Value)
	assert.Equal(t, tags[1].Name, tagsRes[1].Name)
	assert.Equal(t, tags[1].Value, tagsRes[1].Value)

	tagsRes = TagsEncode(tags)
	assert.Equal(t, tagsBase64[0].Name, tagsRes[0].Name)
	assert.Equal(t, tagsBase64[0].Value, tagsRes[0].Value)
	assert.Equal(t, tagsBase64[1].Name, tagsRes[1].Name)
	assert.Equal(t, tagsBase64[1].Value, tagsRes[1].Value)
}

func TestSerializeTags(t *testing.T) {
	tags := []types.Tag{{Name: "abc", Value: "123"}, {Name: "bbc", Value: "223"}}

	by, err := SerializeTags(tags)
	assert.NoError(t, err)
	t.Log(hex.EncodeToString(by))
	//   040661626306313233066262630632323300
	// 03200661626306313233066262630632323300
}
