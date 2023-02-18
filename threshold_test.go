package goar

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/daqiancode/goar/utils"
	tcrsa "github.com/everFinance/ttcrsa"
	"github.com/stretchr/testify/assert"
)

// TestCreateKeyPair Secret key creation and threshold signature
func TestCreateTcKeyPair(t *testing.T) {
	exampleData := []byte("aaabbbcccddd112233")
	signHashed := sha256.Sum256(exampleData)
	salt := sha256.Sum256([]byte("everHash salt aaa"))

	/* -------------------------- Key pair that generates RSA threshold signature on the server side ----------------------------*/
	bitSize := 1024 // If the values are 2048 and 4096, then the generation functions below will perform minute-level times, and we need 4096 bits as the maximum safety level for production environments.
	l := 5
	k := 3
	// Keyshares are distributed to each signer, and KeyMeta stores publicKey, k, l and other public information, which should be sent to the signer together
	keyShares, keyMeta, err := CreateTcKeyPair(bitSize, k, l)
	if err != nil {
		panic(err)
	}

	ts, err := NewTcSign(keyMeta, exampleData, salt[:])
	if err != nil {
		panic(err)
	}

	/* -------------------------- Distribute KeyShare to the signatories ----------------------------*/
	signer01 := keyShares[0]
	signer02 := keyShares[1]
	signer03 := keyShares[2]
	signer04 := keyShares[3]
	signer05 := keyShares[4]

	/* -------------------------- Each signer signs the data received and submits it to the server ----------------------------*/
	// sign data

	signedData01, err := ts.ThresholdSign(signer01)
	if err != nil {
		panic(err)
	}

	signedData02, err := ts.ThresholdSign(signer02)
	if err != nil {
		panic(err)
	}

	signedData03, err := ts.ThresholdSign(signer03)
	if err != nil {
		panic(err)
	}

	signedData04, err := ts.ThresholdSign(signer04)
	if err != nil {
		panic(err)
	}

	signedData05, err := ts.ThresholdSign(signer05)
	if err != nil {
		panic(err)
	}

	/* -------------------------- After receiving the signature data submitted by the signers, the server verifies the signature and assembles the signature ----------------------------*/
	signedShares := tcrsa.SigShareList{
		signedData01,
		signedData02,
		signedData03,
		signedData04,
		signedData05,
	}

	ts, err = NewTcSign(keyMeta, exampleData, salt[:])
	assert.NoError(t, err)
	signature, err := ts.AssembleSigShares(signedShares)
	if err != nil {
		panic(err)
	}

	err = utils.Verify(exampleData, keyMeta.PublicKey, signature)
	if err != nil {
		panic(err)
	}

	/* --------------------------------------------------------------------------------------------------------------*/
	/* -------------------------- Next, we will test the threshold based on the above environment ----------------------------*/
	// As can be seen from the above, L =5, K =3, there are 5 signers, and the threshold is 3.
	// However, in the above process, we have submitted the signatures of all 5 signers, and the goal of threshold test has not been reached.
	// The threshold test is as follows

	// 1. Submit the signature data of signer1,2,3 and verify it
	signedShares123 := tcrsa.SigShareList{
		signedData01,
		signedData02,
		signedData03,
	}
	// assemble
	signature123, err := ts.AssembleSigShares(signedShares123)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature123, nil)
	if err != nil {
		panic(err)
	}

	// 2. Submit signer 3,2,1
	signedShares321 := tcrsa.SigShareList{
		signedData03,
		signedData02,
		signedData01,
	}
	// assemble
	signature321, err := ts.AssembleSigShares(signedShares321)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature321, nil)
	if err != nil {
		panic(err)
	}

	// 3. Submit signer 3,1,2
	signedShares312 := tcrsa.SigShareList{
		signedData03,
		signedData01,
		signedData02,
	}
	// assemble
	signature312, err := ts.AssembleSigShares(signedShares312)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature312, nil)
	if err != nil {
		panic(err)
	}

	// 4. Submit signer 1,3,5
	signedShares135 := tcrsa.SigShareList{
		signedData01,
		signedData03,
		signedData05,
	}
	// assemble
	signature135, err := ts.AssembleSigShares(signedShares135)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature135, nil)
	if err != nil {
		panic(err)
	}

	// 5. Submit signer 5, 1, 4
	signedShares514 := tcrsa.SigShareList{
		signedData05,
		signedData01,
		signedData04,
	}
	// assemble
	signature514, err := ts.AssembleSigShares(signedShares514)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature514, nil)
	if err != nil {
		panic(err)
	}

	// 6. Submit signer 2,3,4,5
	signedShares2345 := tcrsa.SigShareList{
		signedData02,
		signedData03,
		signedData04,
		signedData05,
	}
	// assemble
	signature2345, err := ts.AssembleSigShares(signedShares2345)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature2345, nil)
	if err != nil {
		panic(err)
	}

	// 7. Submit signer 5,4,2,3
	signedShares5423 := tcrsa.SigShareList{
		signedData05,
		signedData04,
		signedData02,
		signedData03,
	}
	// assemble
	signature5423, err := ts.AssembleSigShares(signedShares5423)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature5423, nil)
	if err != nil {
		panic(err)
	}

	// 8. Submit 5， 4， 3，2，1
	signedShares54321 := tcrsa.SigShareList{
		signedData05,
		signedData04,
		signedData03,
		signedData02,
		signedData01,
	}
	// assemble
	signature54321, err := ts.AssembleSigShares(signedShares54321)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature54321, nil)
	if err != nil {
		panic(err)
	}

	// 9. Submit 4，3
	signedShares43 := tcrsa.SigShareList{
		signedData04,
		signedData03,
	}
	// assemble
	_, err = ts.AssembleSigShares(signedShares43)
	assert.EqualError(t, err, "insufficient number of signature shares. provided: 2, needed: 3")

	/*
		Verify that the refactored signature data is committed.
		So let's start with the conclusion：
		threshold k=3, Pass the test [a,b,c,x,x,x]; [a,b,c,d,x,x];
		Fail test [a,a,b,c,d,e,x,x,x]; [a,b,b,c,x,x,x];
		why: No matter how many signatures you submit, the Join method only takes the first 3(k) pieces of the signedShares array to assemble the final signature.

	*/

	// 10. Submit 1，1，3，4; Is the signature of the assembly data signer01 signer01, signer03, so only two valid signatures, below the threshold
	signedShares1134 := tcrsa.SigShareList{
		signedData01,
		signedData01,
		signedData03,
		signedData04,
	}
	// assemble
	signature1134, err := ts.AssembleSigShares(signedShares1134)
	assert.EqualError(t, err, "crypto/rsa: verification error")
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature1134, nil)
	assert.EqualError(t, err, "crypto/rsa: verification error")

	// 11. submit 1，2，2，3; Same thing as above
	signedShares1223 := tcrsa.SigShareList{
		signedData01,
		signedData02,
		signedData02,
		signedData03,
	}
	// assemble
	signature1223, err := ts.AssembleSigShares(signedShares1223)
	assert.EqualError(t, err, "crypto/rsa: verification error")
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature1223, nil)
	assert.EqualError(t, err, "crypto/rsa: verification error")

	// 12. submit 3，2，5，5，5，4，2，1，3；Can get 3, 2 and 5 signature data to meet the threshold value
	signedShares325554213 := tcrsa.SigShareList{
		signedData03,
		signedData02,
		signedData05,
		signedData05,
		signedData05,
		signedData04,
		signedData02,
		signedData01,
		signedData03,
	}
	// assemble
	signature325554213, err := ts.AssembleSigShares(signedShares325554213)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature325554213, nil)
	if err != nil {
		panic(err)
	}

}

// GetKeyPairByLocal
func GetKeyPairFormLocalFile() (shares tcrsa.KeyShareList, meta *tcrsa.KeyMeta, err error) {
	dd, err := ioutil.ReadFile("keyMeta.json")
	if err != nil {
		return nil, nil, err
	}
	ee, err := ioutil.ReadFile("keyShares.json")
	if err != nil {
		return nil, nil, err
	}

	keyMeta := &tcrsa.KeyMeta{}
	err = json.Unmarshal(dd, keyMeta)
	if err != nil {
		return nil, nil, err
	}

	keyShares := tcrsa.KeyShareList{}
	err = json.Unmarshal(ee, &keyShares)
	if err != nil {
		return nil, nil, err
	}
	return keyShares, keyMeta, nil
}

// TestCreateKeyPair3 get address
// func TestCreateKeyPair3(t *testing.T) {
// 	keyMeta := &tcrsa.KeyMeta{}
// 	keyMetaBy, err := ioutil.ReadFile("keyMeta.json") // replace your key
// 	assert.NoError(t, err)
// 	err = json.Unmarshal(keyMetaBy, keyMeta)
// 	assert.NoError(t, err)
// 	addr := sha256.Sum256(keyMeta.PublicKey.N.Bytes())
// 	t.Log("address: ", utils.Base64Encode(addr[:])) // KKzL8og7VFLNwxbwW6cpUY_WkE5jFjWL26cTvKfWYms
// }

// TestCreateKeyPair2 send ar tx by threshold signature keypair
func TestCreateKeyPair2(t *testing.T) {
	// cli := client.New("https://arweave.net")
	//
	// target := "Ii5wAMlLNz13n26nYY45mcZErwZLjICmYd46GZvn4ck"
	// reward, err := cli.GetTransactionPrice(nil, &target)
	// assert.NoError(t, err)
	// // anchor, err := cli.GetTransactionAnchor() // for test
	// anchor, err := cli.GetLastTransactionID("KKzL8og7VFLNwxbwW6cpUY_WkE5jFjWL26cTvKfWYms")
	// assert.NoError(t, err)
	// t.Log("lastTx: ", anchor)
	// // read created threshold keypair for local file; need to be generated ahead of time;
	// keyMeta := &tcrsa.KeyMeta{}
	// keyMetaBy, err := ioutil.ReadFile("keyMeta.json")
	// assert.NoError(t, err)
	// err = json.Unmarshal(keyMetaBy, keyMeta)
	// assert.NoError(t, err)
	//
	// owner := utils.Base64Encode(keyMeta.PublicKey.N.Bytes())
	//
	// amount := big.NewInt(140000) // transfer amount
	// tags := []types.Tag{{Name: "Content-Type", Value: "application/json"}, {Name: "tcrsa", Value: "sandyTest"}}
	// tx := &types.Transaction{
	// 	Format:    2,
	// 	ID:        "",
	// 	LastTx:    anchor,
	// 	Owner:     owner,
	// 	Tags:      types.TagsEncode(tags),
	// 	Target:    target,
	// 	Quantity:  amount.String(),
	// 	Data:      "",
	// 	DataSize:  "0",
	// 	DataRoot:  "",
	// 	Reward:    fmt.Sprintf("%d", reward),
	// 	Signature: "",
	// 	Chunks:    nil,
	// }
	// signData, err := types.GetSignatureData(tx)
	// assert.NoError(t, err)
	// t.Log("signData: ", signData)
	//
	// // signature
	// keyShares := tcrsa.KeyShareList{}
	// keySharesBy, err := ioutil.ReadFile("keyShares.json")
	// assert.NoError(t, err)
	// err = json.Unmarshal(keySharesBy, &keyShares)
	// assert.NoError(t, err)
	//
	// ts, err := NewTcSign(keyMeta, signData)
	// assert.NoError(t, err)
	//
	// /* --------------------------distribute keyShares to the signers ----------------------------*/
	// signer01 := keyShares[0]
	// signer02 := keyShares[1]
	// signer03 := keyShares[2]
	// signer04 := keyShares[3]
	// signer05 := keyShares[4]
	//
	// /* -------------------------- signers to sign data ----------------------------*/
	// signedData01, err := ts.ThresholdSign(signer01)
	// if err != nil {
	// 	panic(err)
	// }
	// t.Log(signedData01.Id)
	// bb, _ := json.Marshal(signedData01)
	// t.Log(hex.EncodeToString(bb))
	//
	// signedData02, err := ts.ThresholdSign(signer02)
	// if err != nil {
	// 	panic(err)
	// }
	// t.Log(signedData02.Id)
	// bb, _ = json.Marshal(signedData02)
	// t.Log(hex.EncodeToString(bb))
	//
	// signedData03, err := ts.ThresholdSign(signer03)
	// if err != nil {
	// 	panic(err)
	// }
	// t.Log(signedData03.Id)
	// bb, _ = json.Marshal(signedData03)
	// t.Log(hex.EncodeToString(bb))
	//
	// signedData04, err := ts.ThresholdSign(signer04)
	// if err != nil {
	// 	panic(err)
	// }
	// t.Log(signedData04.Id)
	// bb, _ = json.Marshal(signedData04)
	// t.Log(hex.EncodeToString(bb))
	//
	// signedData05, err := ts.ThresholdSign(signer05)
	// if err != nil {
	// 	panic(err)
	// }
	// t.Log(signedData05.Id)
	// bb, _ = json.Marshal(signedData05)
	// t.Log(hex.EncodeToString(bb))
	//
	// /* -------------------------- After receiving the signature data submitted by the signers, the server verifies the signature and assembles the signature ----------------------------*/
	// // Collect the signer's signature data into an array
	// signedShares := tcrsa.SigShareList{
	// 	// signedData01,
	// 	signedData02,
	// 	signedData03,
	// 	signedData04,
	// 	// signedData05,
	// }
	//
	// // Verify the signature of each collected signer. And what happens in practice is that the server receives the signature submitted by the signer and then it verifies it and then it puts it in the array above
	// for _, sd := range signedShares {
	// 	err = sd.Verify(ts.pssData, keyMeta)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	//
	// // assemble signatures
	// signature, err := ts.AssembleSigShares(signedShares)
	// if err != nil {
	// 	panic(err)
	// }
	// // Finally, RSA native PSS verification signature method is used to verify the aggregated signature
	// signHashed := sha256.Sum256(signData)
	// err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature, nil)
	// if err != nil {
	// 	panic(err)
	// }
	// // assemble tx and send to ar chain
	// tx.AddSignature(signature)
	// t.Log("txHash: ", tx.ID)
	//
	// status, code, err := cli.SubmitTransaction(tx)
	// assert.NoError(t, err)
	// t.Log("status: ", status)
	// t.Log("code: ", code)
}
