package rsa_crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/niclabs/tcrsa"
	"github.com/zyjblockchain/sandy_log/log"
)

func AssembleSigShares(signedShares tcrsa.SigShareList, meta *tcrsa.KeyMeta, signData []byte) ([]byte, error) {
	signHashed := sha256.Sum256(signData)
	signDataByPss, err := tcrsa.PreparePssDocumentHash(meta.PublicKey.N.BitLen(), crypto.SHA256, signHashed[:], &rsa.PSSOptions{
		SaltLength: 0,
		Hash:       crypto.SHA256,
	})

	// verify each signer share
	for _, sd := range signedShares {
		if err := sd.Verify(signDataByPss, meta); err != nil {
			log.Errorf("verify signer %d sign failed; err: %v", sd.Id, err)
			return nil, err
		}
	}
	signature, err := signedShares.Join(signDataByPss, meta)
	if err != nil {
		log.Errorf("signedShares.Join(signDataByPss, meta) error: %v", err)
		return nil, err
	}

	// verify
	if err := rsa.VerifyPSS(meta.PublicKey, crypto.SHA256, signHashed[:], signature, nil); err != nil {
		log.Errorf("verify signature error; %v", err)
		return nil, err
	}
	return signature, nil
}

func AddTxSignature(tx *types.Transaction, signature []byte) *types.Transaction {
	txId := sha256.Sum256(signature)
	tx.ID = utils.Base64Encode(txId[:])

	tx.Signature = utils.Base64Encode(signature)
	return tx
}
