package rsa_crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/everFinance/goar/client"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/niclabs/tcrsa"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"testing"
)

// 测试 秘钥创建和门限签名
func TestCreateKeyPair(t *testing.T) {
	exampleData := []byte("aaabbbcccddd112233")
	signHashed := sha256.Sum256(exampleData)

	/* -------------------------- 服务器端生成rsa 门限签名的 key pair ----------------------------*/
	bitSize := 1024 // 如果为2048 以及4096 则下面的生成函数会执行分钟级别的时间，生产环境我们需要4096 位作为最高安全级别。
	l := 5
	k := 3
	// keyShares 为分发给每个签名者，keyMeta 里面存储了 publicKey, k, l 等公开的信息，需要一起发送给签名者
	keyShares, keyMeta, err := CreateKeyPair(bitSize, k, l)
	if err != nil {
		panic(err)
	}

	// 构建 pss 加密算法的签名数据(ar tx 目前只支持 pss 的签名)
	signDataByPss, err := tcrsa.PreparePssDocumentHash(keyMeta.PublicKey.N.BitLen(), crypto.SHA256, signHashed[:], nil)
	if err != nil {
		panic(err)
	}

	/* -------------------------- 把keyShare 分发给各个签名者 ----------------------------*/
	signer01 := keyShares[0]
	signer02 := keyShares[1]
	signer03 := keyShares[2]
	signer04 := keyShares[3]
	signer05 := keyShares[4]

	/* -------------------------- 各个签名者对收到的数据进行签名并提交到服务器 ----------------------------*/
	// 分别对数据进行签名
	signedData01, err := signer01.Sign(signDataByPss, crypto.SHA256, keyMeta)
	if err != nil {
		panic(err)
	}

	signedData02, err := signer02.Sign(signDataByPss, crypto.SHA256, keyMeta)
	if err != nil {
		panic(err)
	}

	signedData03, err := signer03.Sign(signDataByPss, crypto.SHA256, keyMeta)
	if err != nil {
		panic(err)
	}

	signedData04, err := signer04.Sign(signDataByPss, crypto.SHA256, keyMeta)
	if err != nil {
		panic(err)
	}

	signedData05, err := signer05.Sign(signDataByPss, crypto.SHA256, keyMeta)
	if err != nil {
		panic(err)
	}

	/* -------------------------- 服务器收到签名者们提交的签名数据之后进行验证签名、组装签名 ----------------------------*/
	// 收集好签名者的签名数据到一个数组中
	signedShares := tcrsa.SigShareList{
		signedData01,
		signedData02,
		signedData03,
		signedData04,
		signedData05,
	}

	// 验证每个收集的签名者的签名。在实际过程中是服务器收到签名者提交的签名就要做一次验证验证通过之后再放入上面的数组
	for _, sd := range signedShares {
		err = sd.Verify(signDataByPss, keyMeta)
		if err != nil {
			panic(err)
		}
	}

	// 组装签名
	signature, err := signedShares.Join(signDataByPss, keyMeta)
	if err != nil {
		panic(err)
	}
	// 最后通过 rsa 原生的pss 验证签名方法来验证聚合之后的签名
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature, nil)
	if err != nil {
		panic(err)
	}

	/* -------------------------- 以上流程是一个门限签名的完整流程，下面我们基于上面的环境来进行门限测试 ----------------------------*/
	// 从上面能看到l=5 k=3, 签名者有5个，门限是3，但是上面流程我们把5个签名者的签名都提交上去了，没有达到门限测试的目的，下面测试门限

	// 1. 提交signer1,2,3 的签名数据并验证
	signedShares123 := tcrsa.SigShareList{
		signedData01,
		signedData02,
		signedData03,
	}
	// 组装
	signature123, err := signedShares123.Join(signDataByPss, keyMeta)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature123, nil)
	if err != nil {
		panic(err)
	}

	// 2. 提交 signer 3,2,1
	signedShares321 := tcrsa.SigShareList{
		signedData03,
		signedData02,
		signedData01,
	}
	// 组装
	signature321, err := signedShares321.Join(signDataByPss, keyMeta)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature321, nil)
	if err != nil {
		panic(err)
	}

	// 3. 提交 signer 3,1,2
	signedShares312 := tcrsa.SigShareList{
		signedData03,
		signedData01,
		signedData02,
	}
	// 组装
	signature312, err := signedShares312.Join(signDataByPss, keyMeta)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature312, nil)
	if err != nil {
		panic(err)
	}

	// 4. 提交 signer 1,3,5
	signedShares135 := tcrsa.SigShareList{
		signedData01,
		signedData03,
		signedData05,
	}
	// 组装
	signature135, err := signedShares135.Join(signDataByPss, keyMeta)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature135, nil)
	if err != nil {
		panic(err)
	}

	// 5. 提交 signer 5, 1, 4
	signedShares514 := tcrsa.SigShareList{
		signedData05,
		signedData01,
		signedData04,
	}
	// 组装
	signature514, err := signedShares514.Join(signDataByPss, keyMeta)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature514, nil)
	if err != nil {
		panic(err)
	}

	// 6. 提交 signer 2,3,4,5
	signedShares2345 := tcrsa.SigShareList{
		signedData02,
		signedData03,
		signedData04,
		signedData05,
	}
	// 组装
	signature2345, err := signedShares2345.Join(signDataByPss, keyMeta)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature2345, nil)
	if err != nil {
		panic(err)
	}

	// 7. 提交 signer 5,4,2,3
	signedShares5423 := tcrsa.SigShareList{
		signedData05,
		signedData04,
		signedData02,
		signedData03,
	}
	// 组装
	signature5423, err := signedShares5423.Join(signDataByPss, keyMeta)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature5423, nil)
	if err != nil {
		panic(err)
	}

	// 8. 提交 5， 4， 3，2，1
	signedShares54321 := tcrsa.SigShareList{
		signedData05,
		signedData04,
		signedData03,
		signedData02,
		signedData01,
	}
	// 组装
	signature54321, err := signedShares54321.Join(signDataByPss, keyMeta)
	if err != nil {
		panic(err)
	}
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature54321, nil)
	if err != nil {
		panic(err)
	}

	// 9. 提交 4，3
	signedShares43 := tcrsa.SigShareList{
		signedData04,
		signedData03,
	}
	// 组装
	_, err = signedShares43.Join(signDataByPss, keyMeta)
	assert.EqualError(t, err, "insufficient number of signature shares. provided: 2, needed: 3")

	/*
		验证提交重构签名数据的情况。
		先给出结论：
		门限值k=3, 能通过的情况[a,b,c,x,x,x]; [a,b,c,d,x,x];
		不能通过的情况[a,a,b,c,d,e,x,x,x]; [a,b,b,c,x,x,x];
		原因： 不管你提交多少签名上来，Join 方法只取 signedShares 数组中的前3(k)个数据来组装最终的签名。

	*/

	// 10. 提交 1，1，3，4; 组装的签名数据是signer01,signer01,signer03, 这样只有2个有效的签名，低于门限值
	signedShares1134 := tcrsa.SigShareList{
		signedData01,
		signedData01,
		signedData03,
		signedData04,
	}
	// 组装
	signature1134, err := signedShares1134.Join(signDataByPss, keyMeta)
	assert.NoError(t, err)
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature1134, nil)
	assert.EqualError(t, err, "crypto/rsa: verification error")

	// 11. 提交1，2，2，3; 和上面情况一致
	signedShares1223 := tcrsa.SigShareList{
		signedData01,
		signedData02,
		signedData02,
		signedData03,
	}
	// 组装
	signature1223, err := signedShares1223.Join(signDataByPss, keyMeta)
	assert.NoError(t, err)
	// verify
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature1223, nil)
	assert.EqualError(t, err, "crypto/rsa: verification error")

	// 12. 提交 3，2，5，5，5，4，2，1，3； 是可以的，能取到3，2，5 这三个签名数据满足门限值
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
	// 组装
	signature325554213, err := signedShares325554213.Join(signDataByPss, keyMeta)
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

// 获取门限秘钥地址
func TestCreateKeyPair3(t *testing.T) {
	keyMeta := &tcrsa.KeyMeta{}
	keyMetaBy, err := ioutil.ReadFile("keyMeta.json")
	assert.NoError(t, err)
	err = json.Unmarshal(keyMetaBy, keyMeta)
	assert.NoError(t, err)
	addr := sha256.Sum256(keyMeta.PublicKey.N.Bytes())
	t.Log("address: ", utils.Base64Encode(addr[:])) // KKzL8og7VFLNwxbwW6cpUY_WkE5jFjWL26cTvKfWYms
}

// 测试 通过门限签名发送 ar 交易
func TestCreateKeyPair2(t *testing.T) {
	cli := client.New("https://arweave.net")

	target := "Ii5wAMlLNz13n26nYY45mcZErwZLjICmYd46GZvn4ck"
	reward, err := cli.GetTransactionPrice(nil, &target)
	assert.NoError(t, err)
	// anchor, err := cli.GetTransactionAnchor()
	anchor, err := cli.GetLastTransactionID("KKzL8og7VFLNwxbwW6cpUY_WkE5jFjWL26cTvKfWYms")
	assert.NoError(t, err)
	t.Log("lastTx: ", anchor)
	// 读取本地生产的 门限签名秘钥对的公钥部分，注：测试使用的是 4096 bit 的秘钥，需要先单独生成并放到本地。
	keyMeta := &tcrsa.KeyMeta{}
	keyMetaBy, err := ioutil.ReadFile("keyMeta.json")
	assert.NoError(t, err)
	err = json.Unmarshal(keyMetaBy, keyMeta)
	assert.NoError(t, err)

	owner := utils.Base64Encode(keyMeta.PublicKey.N.Bytes())

	amount := big.NewInt(130000) // 转账余额
	tags := []types.Tag{{Name: "Content-Type", Value: "application/json"}, {Name: "tcrsa", Value: "sandyTest"}}
	tx := &types.Transaction{
		Format:    2,
		ID:        "",
		LastTx:    anchor,
		Owner:     owner,
		Tags:      types.TagsEncode(tags),
		Target:    target,
		Quantity:  amount.String(),
		Data:      []byte{},
		DataSize:  "0",
		DataRoot:  "",
		Reward:    fmt.Sprintf("%d", reward),
		Signature: "",
		Chunks:    nil,
	}
	signData, err := types.GetSignatureData(tx)
	assert.NoError(t, err)
	t.Log("signData: ", signData)

	// 门限签名
	keyShares := tcrsa.KeyShareList{}
	keySharesBy, err := ioutil.ReadFile("keyShares.json")
	assert.NoError(t, err)
	err = json.Unmarshal(keySharesBy, &keyShares)
	assert.NoError(t, err)

	signHashed := sha256.Sum256(signData)
	signDataByPss, err := tcrsa.PreparePssDocumentHash(keyMeta.PublicKey.N.BitLen(), crypto.SHA256, signHashed[:], &rsa.PSSOptions{
		SaltLength: 0,
		Hash:       crypto.SHA256,
	})
	assert.NoError(t, err)

	/* -------------------------- 把 keyShare 分发给各个签名者 ----------------------------*/
	signer01 := keyShares[0]
	signer02 := keyShares[1]
	signer03 := keyShares[2]
	signer04 := keyShares[3]
	signer05 := keyShares[4]

	/* -------------------------- 各个签名者对收到的数据进行签名并提交到服务器 ----------------------------*/
	// 分别对数据进行签名
	signedData01, err := signer01.Sign(signDataByPss, crypto.SHA256, keyMeta)
	if err != nil {
		panic(err)
	}
	t.Log(signedData01.Id)

	signedData02, err := signer02.Sign(signDataByPss, crypto.SHA256, keyMeta)
	if err != nil {
		panic(err)
	}
	t.Log(signedData02.Id)

	signedData03, err := signer03.Sign(signDataByPss, crypto.SHA256, keyMeta)
	if err != nil {
		panic(err)
	}
	t.Log(signedData03.Id)

	signedData04, err := signer04.Sign(signDataByPss, crypto.SHA256, keyMeta)
	if err != nil {
		panic(err)
	}
	t.Log(signedData04.Id)

	signedData05, err := signer05.Sign(signDataByPss, crypto.SHA256, keyMeta)
	if err != nil {
		panic(err)
	}
	t.Log(signedData05.Id)

	/* -------------------------- 服务器收到签名者们提交的签名数据之后进行验证签名、组装签名 ----------------------------*/
	// 收集好签名者的签名数据到一个数组中
	signedShares := tcrsa.SigShareList{
		// signedData01,
		signedData02,
		signedData03,
		signedData04,
		// signedData05,
	}

	// 验证每个收集的签名者的签名。在实际过程中是服务器收到签名者提交的签名就要做一次验证验证通过之后再放入上面的数组
	for _, sd := range signedShares {
		err = sd.Verify(signDataByPss, keyMeta)
		if err != nil {
			panic(err)
		}
	}

	// 组装签名
	signature, err := signedShares.Join(signDataByPss, keyMeta)
	if err != nil {
		panic(err)
	}
	// 最后通过 rsa 原生的pss 验证签名方法来验证聚合之后的签名
	err = rsa.VerifyPSS(keyMeta.PublicKey, crypto.SHA256, signHashed[:], signature, nil)
	if err != nil {
		panic(err)
	}
	// assemble tx and send to ar chain
	txId := sha256.Sum256(signature)
	tx.ID = utils.Base64Encode(txId[:])
	t.Log("txHash: ", tx.ID)
	tx.Signature = utils.Base64Encode(signature)

	status, code, err := cli.SubmitTransaction(tx)
	assert.NoError(t, err)
	t.Log("status: ", status)
	t.Log("code: ", code)
}
