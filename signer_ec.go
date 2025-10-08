package goar

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
)

type EcSigner struct {
	private *ecdsa.PrivateKey
}

func NewEcSigner(prvHex string) (*EcSigner, error) {
	k, err := crypto.HexToECDSA(prvHex)
	if err != nil {
		return nil, err
	}
	return &EcSigner{private: k}, nil
}

func (s *EcSigner) SignTx(tx *types.Transaction) error {
	payload, err := utils.GetSignatureData(tx)
	if err != nil {
		return err
	}
	hash := sha256.Sum256(payload)
	signature, err := crypto.Sign(hash[:], s.private)
	if err != nil {
		return err
	}
	txId := sha256.Sum256(signature)
	tx.ID = utils.Base64Encode(txId[:])
	tx.Signature = utils.Base64Encode(signature)
	return nil
}

func (s *EcSigner) Public() []byte {
	return crypto.CompressPubkey(&s.private.PublicKey)
}
func (s *EcSigner) Owner() string {
	return utils.Base64Encode(s.Public())
}

func (s *EcSigner) Address() string {
	addr := sha256.Sum256(s.Public())
	return utils.Base64Encode(addr[:])
}

func VerifyEcdsaTxSig(tx *types.Transaction) error {
	if tx.Signature == "" {
		return fmt.Errorf("no signature to verify")
	}

	signatureData, err := utils.GetSignatureData(tx)
	if err != nil {
		return err
	}

	hash := sha256.Sum256(signatureData)

	signatureBytes, err := utils.Base64Decode(tx.Signature)
	if err != nil {
		return err
	}

	publicKey, err := crypto.SigToPub(hash[:], signatureBytes)
	if err != nil {
		return err
	}

	ok := crypto.VerifySignature(
		crypto.FromECDSAPub(publicKey),
		hash[:],
		signatureBytes[:64],
	)
	if !ok {
		return errors.New("signature incorrect")
	}
	return nil
}

func GetEcdsaTxOwner(tx *types.Transaction) (string, error) {
	publicKey, err := recoverPubkey(tx)
	if err != nil {
		return "", err
	}
	return utils.Base64Encode(crypto.CompressPubkey(publicKey)), nil
}

func recoverPubkey(tx *types.Transaction) (*ecdsa.PublicKey, error) {
	if tx.Signature == "" {
		return nil, errors.New("not signed")
	}

	payload, err := utils.GetSignatureData(tx)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(payload)

	signature, err := utils.Base64Decode(tx.Signature)
	if err != nil {
		return nil, err
	}

	publicKey, err := crypto.SigToPub(hash[:], signature)
	return publicKey, err
}

func OwnerToAddress(owner string) (string, error) {
	public, err := utils.Base64Decode(owner)
	if err != nil {
		return "", err
	}
	addr := sha256.Sum256(public)
	return utils.Base64Encode(addr[:]), nil
}
