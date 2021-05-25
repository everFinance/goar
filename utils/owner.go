package utils

import (
	"crypto/rsa"
	"crypto/sha256"
	"math/big"
)

func OwnerToAddress(owner string) (string, error) {
	by, err := Base64Decode(owner)
	if err != nil {
		return "", err
	}
	addr := sha256.Sum256(by)
	return Base64Encode(addr[:]), nil
}

func OwnerToPubKey(owner string) (*rsa.PublicKey, error) {
	by, err := Base64Decode(owner)
	if err != nil {
		return nil, err
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(by),
		E: 65537, //"AQAB"
	}, nil
}
