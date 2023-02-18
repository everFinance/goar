package goar

import (
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io/ioutil"

	"github.com/daqiancode/goar/types"
	"github.com/daqiancode/goar/utils"
	"github.com/everFinance/gojwk"
)

type Signer struct {
	Address string
	PubKey  *rsa.PublicKey
	PrvKey  *rsa.PrivateKey
}

func NewSignerFromPath(path string) (*Signer, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return NewSigner(b)
}

func NewSigner(b []byte) (*Signer, error) {
	key, err := gojwk.Unmarshal(b)
	if err != nil {
		return nil, err
	}
	pubKey, err := key.DecodePublicKey()
	if err != nil {
		return nil, err
	}
	pub, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		err = fmt.Errorf("pubKey type error")
		return nil, err
	}

	prvKey, err := key.DecodePrivateKey()
	if err != nil {
		return nil, err
	}
	prv, ok := prvKey.(*rsa.PrivateKey)
	if !ok {
		err = fmt.Errorf("prvKey type error")
		return nil, err
	}
	addr := sha256.Sum256(pub.N.Bytes())
	return &Signer{
		Address: utils.Base64Encode(addr[:]),
		PubKey:  pub,
		PrvKey:  prv,
	}, nil
}

func (s *Signer) SignTx(tx *types.Transaction) error {
	return utils.SignTransaction(tx, s.PrvKey)
}

func (s *Signer) Owner() string {
	return utils.Base64Encode(s.PubKey.N.Bytes())
}

func (s *Signer) SignMsg(msg []byte) ([]byte, error) {
	return utils.Sign(msg, s.PrvKey)
}
