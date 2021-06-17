package threshold

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"github.com/everFinance/sandy_log/log"
	"github.com/niclabs/tcrsa"
)

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
