package bundles

import (
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
)

const (
	MAX_TAG_KEY_LENGTH_BYTES   = 1024 * 1
	MAX_TAG_VALUE_LENGTH_BYTES = 1024 * 3
	MAX_TAG_COUNT              = 128
)

func VerifyEncodedTags(tags []types.Tag) bool {
	if len(tags) > MAX_TAG_COUNT {
		return false
	}
	// Search for something invalid
	for _, tag := range tags {
		if !verifyEncodeTagSize(tag) {
			return false
		}
	}
	return true
}

func verifyEncodeTagSize(tag types.Tag) bool {
	name, err := utils.Base64Decode(tag.Name)
	if err != nil || len(name) > MAX_TAG_KEY_LENGTH_BYTES {
		return false
	}

	value, err := utils.Base64Decode(tag.Value)
	if err != nil || len(value) > MAX_TAG_VALUE_LENGTH_BYTES {
		return false
	}
	return true
}
