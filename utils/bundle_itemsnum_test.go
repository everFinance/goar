package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// A bundle's item count is read from 32 untrusted bytes into a plain int, which
// can overflow to a negative value. DecodeBundle must reject that instead of
// silently returning an empty bundle (the `i < itemsNum` loop would never run
// and `32+itemsNum*64` would be negative, so the length check passed).
func TestDecodeBundle_RejectsNegativeItemsNum(t *testing.T) {
	// Low 8 bytes all 0xff => ByteArrayToLong == -1.
	itemsNum := make([]byte, 32)
	for i := 0; i < 8; i++ {
		itemsNum[i] = 0xff
	}
	assert.Equal(t, -1, ByteArrayToLong(itemsNum))

	input := make([]byte, 40)
	copy(input[:32], itemsNum)

	_, err := DecodeBundle(input)
	assert.Error(t, err, "a negative items number must be rejected, not treated as an empty bundle")
}
