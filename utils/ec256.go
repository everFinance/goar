package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"
)

type EC256Signer struct {
	privateKey *ecdsa.PrivateKey
}

func NewEC256Signature(privateKey *ecdsa.PrivateKey) *EC256Signer {
	return &EC256Signer{privateKey: privateKey}
}

func (e *EC256Signer) Sign(message []byte) ([]byte, error) {
	hash := sha256.Sum256(message)

	r, s, err := ecdsa.Sign(rand.Reader, e.privateKey, hash[:])
	if err != nil {
		return nil, err
	}

	signature, err := asn1.Marshal(struct {
		R, S *big.Int
	}{r, s})

	return signature, err
}

func EC256Verify(pub []byte, data, signature []byte) (bool, error) {

	// Decode PEM block into DER format
	block, _ := pem.Decode(pub)
	if block == nil {
		fmt.Println("Failed to decode PEM block")
		return false, nil
	}

	// Parse DER public key
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("Failed to parse public key:", err)
		return false, err
	}

	// Type assert to ECDSA public key
	pubKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("Not a ECDSA public key")
		return false, nil
	}

	sigStruct := struct {
		R, S *big.Int
	}{}
	if _, err := asn1.Unmarshal(signature, &sigStruct); err != nil {
		return false, err
	}

	hash := sha256.Sum256(data)

	valid := ecdsa.Verify(pubKey, hash[:], sigStruct.R, sigStruct.S)
	return valid, nil
}

//
