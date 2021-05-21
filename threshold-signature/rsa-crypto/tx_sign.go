package rsa_crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"github.com/everFinance/goar/types"
	"github.com/niclabs/tcrsa"
)

// ThresholdSignArTx single share sign tx
func ThresholdSignArTx(meta *tcrsa.KeyMeta, signShare tcrsa.KeyShare, tx *types.Transaction) (*tcrsa.SigShare, error) {
	signData, err := types.GetSignatureData(tx)
	if err != nil {
		return nil, err
	}

	// sign
	signHashed := sha256.Sum256(signData)

	signDataByPss, err := tcrsa.PreparePssDocumentHash(meta.PublicKey.N.BitLen(), crypto.SHA256, signHashed[:], &rsa.PSSOptions{
		SaltLength: 0,
		Hash:       crypto.SHA256,
	})
	if err != nil {
		return nil, err
	}

	signedData, err := signShare.Sign(signDataByPss, crypto.SHA256, meta)
	if err != nil {
		return nil, err
	}

	// verify
	if err := signedData.Verify(signDataByPss, meta); err != nil {
		return nil, err
	}

	return signedData, nil
}
