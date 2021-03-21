package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignAndVerify(t *testing.T) {
	rightKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	wrongKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	msg := []byte("123")

	sig, err := Sign(msg[:], rightKey)
	assert.NoError(t, err)
	assert.NoError(t, Verify(msg, &rightKey.PublicKey, sig))
	assert.Error(t, Verify(msg, &wrongKey.PublicKey, sig))
}
