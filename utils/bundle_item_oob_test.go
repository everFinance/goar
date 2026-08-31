package utils

import (
	"testing"

	"github.com/everFinance/goar/types"
	"github.com/stretchr/testify/assert"
)

// A bundle item whose binary ends right after the target/anchor presence flags
// leaves fewer than tagsStart+8 bytes for the tag-count field. DecodeBundleItem
// must return an error instead of panicking with slice-out-of-range.
// (GetBundleItemTagsBytes already guards this; DecodeBundleItem did not.)
func TestDecodeBundleItem_TruncatedBeforeTagsDoesNotPanic(t *testing.T) {
	sc, ok := types.SigConfigMap[1]
	if !ok {
		t.Skip("sig type 1 not configured")
	}
	// 2 (sigType) + sig + owner + target-present + anchor-present, then nothing.
	position := 2 + sc.SigLength + sc.PubLength
	total := position + 2
	b := make([]byte, total)
	b[0] = 1 // sigType 1, little-endian

	assert.NotPanics(t, func() {
		_, err := DecodeBundleItem(b)
		assert.Error(t, err)
	})
}

// When numOfTags == 0 the tags branch is skipped, so the tagsStart+16 bound was
// never checked before slicing `data := itemBinary[tagsStart+16+...:]`. An item
// that is exactly tagsStart+8 bytes (zero tag count, nothing after) must error,
// not panic.
func TestDecodeBundleItem_TruncatedDataWithZeroTagsDoesNotPanic(t *testing.T) {
	sc, ok := types.SigConfigMap[1]
	if !ok {
		t.Skip("sig type 1 not configured")
	}
	position := 2 + sc.SigLength + sc.PubLength
	tagsStart := position + 2 // no target, no anchor
	b := make([]byte, tagsStart+8)
	b[0] = 1 // sigType 1; tag-count bytes are zero => numOfTags == 0

	assert.NotPanics(t, func() {
		_, err := DecodeBundleItem(b)
		assert.Error(t, err)
	})
}
