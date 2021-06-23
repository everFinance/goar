package threshold

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	tcrsa "github.com/everFinance/ttcrsa"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestNewTcSign(t *testing.T) {
	keyMeta := &tcrsa.KeyMeta{}
	keyMetaBy, err := ioutil.ReadFile("keyMeta.json")
	assert.NoError(t, err)
	err = json.Unmarshal(keyMetaBy, keyMeta)
	assert.NoError(t, err)

	signData := []byte("aaabbbbccc")
	signHashed := sha256.Sum256(signData)

	for i := 0; i < 5; i++ {
		signDataByPss, err := tcrsa.PreparePssDocumentHash(keyMeta.PublicKey.N.BitLen(), crypto.SHA256, signHashed[:], &rsa.PSSOptions{
			SaltLength: 0,
			Hash:       crypto.SHA256,
		})
		assert.NoError(t, err)
		t.Log(hex.EncodeToString(signDataByPss))
	}
}
