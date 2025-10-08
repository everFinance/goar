package goar

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyEcdsaTxSig(t *testing.T) {
	// id := "hmtw7VXo-yfn_Gj5g5ZpcWU1P5ZMziU4EMhap1oPjyE"
	tx := &types.Transaction{
		Format: 2,
		ID:     "",
		LastTx: "Pu9WXD63YpX4TQe3iwQKF8-bkbtPx7WpqHtUvBcuqO9Bvq44SyX6dTOMfOCERX1t",
		Owner:  "",
		Tags: []types.Tag{
			{Name: "VGVzdA", Value: "ZWNkc2EtdHg"},
		},
		Target:     "Qa8AAZv-sEhQRIm7xZr3CVLtlzIH8NezaY0GZhURcAc",
		Quantity:   "1",
		Data:       "",
		DataReader: nil,
		DataSize:   "0",
		DataRoot:   "",
		Reward:     "3429726",
		Signature:  "ElI9De-2Z5AT3aoq5TYCz-qHJGa_tmYqqYwC7Rf89_05O5wVr2U-SuH9WPF-iRr1HQGjR-7_XVw0a2LyW0rf4gE",
	}

	sig, _ := utils.Base64Decode(tx.Signature)
	sigHash := sha256.Sum256(sig)
	txId := utils.Base64Encode(sigHash[:])
	t.Log("txId: ", txId)
	owner, err := GetEcdsaTxOwner(tx)
	assert.NoError(t, err)
	t.Log("owner: ", owner)
	address, err := OwnerToAddress(owner)
	assert.NoError(t, err)
	t.Log("address: ", address)

	/*
			{
		        "address": "RymI02hes920xGugzRJ3L54eGg-jVU-_R2uCI057_nU",
		        "key": "AlFhxNH-6NmRDVEukOvvgsrEOWgxLi5x_h-r9Cg_dP27"
		      }
	*/
	err = VerifyEcdsaTxSig(tx)
	assert.NoError(t, err)
}

func genEcc() (private string, public string, address string) {
	priv, _ := crypto.GenerateKey()

	address = crypto.PubkeyToAddress(priv.PublicKey).String()
	private = hex.EncodeToString(crypto.FromECDSA(priv))
	public = hex.EncodeToString(crypto.FromECDSAPub(&priv.PublicKey))
	return
}

func TestGenEcc(t *testing.T) {
	private, public, address := genEcc()
	t.Log("private: ", private)
	t.Log("public: ", public)
	t.Log("address: ", address)
	/*
		private:  ad5071a13d81ddc568b98f839c27209dfa158c8c2a65649c4375abcdcbfb0f91
		public:  046b8783cad5c8bcd9a6f16cba813b4a248014d07463379c2e388c22a71a9421e17a6c6b20c2626f36c39739833079ceb0e40b5620df2841d32bdebb9e7806770f
		address:  0x1d8c6356304B20a18221CcBD22150D5f294eD02f
	*/
}

func TestNewEcSigner(t *testing.T) {
	private := "ad5071a13d81ddc568b98f839c27209dfa158c8c2a65649c4375abcdcbfb0f91"
	ecSigner, err := NewEcSigner(private)
	assert.NoError(t, err)
	t.Log("address: ", ecSigner.Address())

	tx := &types.Transaction{
		Format: 2,
		LastTx: "5Yx84D8oaCAlPqMe7vmDGuEwXtIOy_rMOHlpgm7yovYoW__Mj8S2gC8UuWmLs-61",
		Owner:  "",
		Tags: []types.Tag{
			{Name: "VGVzdA", Value: "ZWNkc2EtdHg"},
		},
		Target:    "cSYOy8-p1QFenktkDBFyRM3cwZSTrQ_J4EsELLho_UE",
		Quantity:  "100000000",
		Data:      "",
		DataSize:  "0",
		DataRoot:  "",
		Reward:    "7050359",
		Signature: "",
	}
	err = ecSigner.SignTx(tx)
	assert.NoError(t, err)

	t.Log("txId: ", tx.ID)
	t.Log("txSig: ", tx.Signature)

	owner, err := GetEcdsaTxOwner(tx)
	assert.NoError(t, err)
	t.Log("owner: ", owner)
	address, err := OwnerToAddress(owner)
	assert.NoError(t, err)
	t.Log("address: ", address)

	err = VerifyEcdsaTxSig(tx)
	assert.NoError(t, err)

	cli := NewClient("https://arweave.net")
	status, code, err := cli.SubmitTransaction(tx)
	t.Log("status: ", status)
	t.Log("code: ", code)
	assert.NoError(t, err)
}
