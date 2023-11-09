package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignAndVerifyES256(t *testing.T) {

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	message := []byte("test message")

	singer := NewEC256Signature(privateKey)

	sig, err := singer.Sign(message)
	assert.NoError(t, err)

	// print signature length and public key length
	t.Log("signature length:", len(sig))
	t.Log("public key length:", len(singer.privateKey.PublicKey.X.Bytes()))

	// get public key bytes
	pubKey := elliptic.Marshal(privateKey.PublicKey.Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)
	valid, err := EC256Verify(pubKey, message, sig)
	assert.NoError(t, err)
	assert.True(t, valid)

	wrongKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// get public key bytes
	wrongPubKey := elliptic.Marshal(wrongKey.PublicKey.Curve, wrongKey.PublicKey.X, wrongKey.PublicKey.Y)

	valid, err = EC256Verify(wrongPubKey, message, sig)
	assert.False(t, valid)

	modifiedMessage := []byte("modified message")
	valid, err = EC256Verify(pubKey, modifiedMessage, sig)
	assert.False(t, valid)
}
