package goar

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/everFinance/sandy_log/log"
	tcrsa "github.com/everFinance/ttcrsa"
)

// CreateTcKeyPair
// bitSize: Is used to generate key shares with a security level equivalent to a RSA private of that size.
// l: creates l key shares for a k-threshold signing scheme.
// k: The generated key shares have a threshold parameter of k
func CreateTcKeyPair(bitSize, k, l int) (shares tcrsa.KeyShareList, meta *tcrsa.KeyMeta, err error) {
	if bitSize > 4096 || bitSize < 512 {
		return nil, nil, errors.New(fmt.Sprintf("bitSize:%d parameter must in [512,4096]", bitSize))
	}
	if k <= 0 || l <= 1 || l < k || k < (l/2+1) {
		return nil, nil, errors.New(fmt.Sprintf("k: %d l: %d parameter incorrect; k must > 0, l must > 1, l must >= k, k must >= (l/2+1)", k, l))
	}

	now := time.Now()
	keyShares, keyMeta, err := tcrsa.NewKey(bitSize, uint16(k), uint16(l), nil)
	if err != nil {
		log.Errorf("tcrsa newKey error; bitSize: %d, k: %d, l: %d; err: %v", bitSize, k, l, err)
		return nil, nil, err
	}
	log.Debugf("Create bit size = %d rsa threshold keyPair spend time: %s", bitSize, time.Since(now).String())
	return keyShares, keyMeta, nil
}

type TcSign struct {
	keyMeta  *tcrsa.KeyMeta
	signData []byte
	pssData  []byte
}

func NewTcSign(meta *tcrsa.KeyMeta, signData []byte) (*TcSign, error) {
	signHashed := sha256.Sum256(signData)

	signDataByPss, err := tcrsa.PreparePssDocumentHash(meta.PublicKey.N.BitLen(), crypto.SHA256, signHashed[:], &rsa.PSSOptions{
		SaltLength: 0,
		Hash:       crypto.SHA256,
	})
	if err != nil {
		return nil, err
	}

	return &TcSign{
		keyMeta:  meta,
		signData: signData,
		pssData:  signDataByPss,
	}, nil
}

// for signer
// ThresholdSignTx single share sign tx
func (ts *TcSign) ThresholdSign(signShare *tcrsa.KeyShare) (*tcrsa.SigShare, error) {

	signedData, err := signShare.Sign(ts.pssData, crypto.SHA256, ts.keyMeta)
	if err != nil {
		return nil, err
	}

	// verify
	if err := signedData.Verify(ts.pssData, ts.keyMeta); err != nil {
		return nil, err
	}
	return signedData, nil
}

// for server hub
// AssembleSigShares
func (ts *TcSign) AssembleSigShares(signedShares tcrsa.SigShareList) ([]byte, error) {
	// verify each signer share
	for _, sd := range signedShares {
		if err := sd.Verify(ts.pssData, ts.keyMeta); err != nil {
			log.Errorf("verify signer %d sign failed; err: %v", sd.Id, err)
			return nil, err
		}
	}
	signature, err := signedShares.Join(ts.pssData, ts.keyMeta)
	if err != nil {
		log.Errorf("signedShares.Join(signDataByPss, meta) error: %v", err)
		return nil, err
	}

	// verify
	signHashed := sha256.Sum256(ts.signData)
	if err := rsa.VerifyPSS(ts.keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature, nil); err != nil {
		log.Errorf("verify signature error; %v", err)
		return nil, err
	}
	return signature, nil
}

// VerifySigShare verify share sig
func (ts *TcSign) VerifySigShare(sigShareData string) error {
	// unmarshal share sig data
	ss := &tcrsa.SigShare{}
	if err := json.Unmarshal([]byte(sigShareData), ss); err != nil {
		return err
	}
	return ss.Verify(ts.pssData, ts.keyMeta)
}
