package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/everFinance/goar/types"
)

func SignTransaction(tx *types.Transaction, pubKey *rsa.PublicKey, prvKey *rsa.PrivateKey) (err error) {
	tx.Owner = Base64Encode(pubKey.N.Bytes())

	// data is not null, generate chunk in tx.data_root
	signData, err := DataHash(tx)
	if err != nil {
		return
	}

	sig, err := Sign(signData, prvKey)
	if err != nil {
		return err
	}

	id := sha256.Sum256(sig)

	tx.ID = Base64Encode(id[:])
	tx.Signature = Base64Encode(sig)
	return
}

func VerifyTransaction(tx types.Transaction) (err error) {
	sig, err := Base64Decode(tx.Signature)
	if err != nil {
		return
	}

	// verify ID
	id := sha256.Sum256(sig)
	if Base64Encode(id[:]) != tx.ID {
		err = fmt.Errorf("wrong id")
	}

	signData, err := DataHash(&tx)
	if err != nil {
		return
	}

	owner, err := Base64Decode(tx.Owner)
	if err != nil {
		return
	}

	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(owner),
		E: 65537, //"AQAB"
	}

	return verify(signData, pubKey, sig)
}

func Sign(msg []byte, prvKey *rsa.PrivateKey) ([]byte, error) {
	hashed := sha256.Sum256(msg)

	return rsa.SignPSS(rand.Reader, prvKey, crypto.SHA256, hashed[:], &rsa.PSSOptions{
		SaltLength: 0,
		Hash:       crypto.SHA256,
	})
}

func verify(msg []byte, pubKey *rsa.PublicKey, sign []byte) error {
	hashed := sha256.Sum256(msg)

	return rsa.VerifyPSS(pubKey, crypto.SHA256, hashed[:], sign, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA256,
	})
}
