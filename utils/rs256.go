package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type RS256Signer struct {
	privateKey *rsa.PrivateKey
}

func NewRS256Signature(privateKey *rsa.PrivateKey) *RS256Signer {
	return &RS256Signer{privateKey: privateKey}
}

func (s *RS256Signer) Sign(data []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, d)
	if err != nil {
		return nil, err
	}

	return signature, nil

}

func RS256Verify(pub []byte, data, signature []byte) error {
	// Decode PEM encoded public key
	block, _ := pem.Decode(pub)
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	publicKey := pubKey.(*rsa.PublicKey)

	h := sha256.New()
	h.Write([]byte(data))
	d := h.Sum(nil)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, d, signature)
	if err != nil {
		return errors.New("invalid signature")
	}
	return nil
}
