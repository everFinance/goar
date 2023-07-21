package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	rightKey, err := GenerateRsaKey(4096)
	assert.NoError(t, err)
	wrongKey, err := GenerateRsaKey(4096)
	assert.NoError(t, err)
	msg := []byte("123")

	sig, err := Sign(msg[:], rightKey)
	assert.NoError(t, err)
	assert.NoError(t, Verify(msg, &rightKey.PublicKey, sig))
	assert.Error(t, Verify(msg, &wrongKey.PublicKey, sig))
}
