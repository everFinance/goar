package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

func Sign(msg []byte, prvKey *rsa.PrivateKey) ([]byte, error) {
	hashed := sha256.Sum256(msg)

	return rsa.SignPSS(rand.Reader, prvKey, crypto.SHA256, hashed[:], &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA256,
	})
}

func Verify(msg []byte, pubKey *rsa.PublicKey, sign []byte) error {
	hashed := sha256.Sum256(msg)

	return rsa.VerifyPSS(pubKey, crypto.SHA256, hashed[:], sign, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA256,
	})
}
