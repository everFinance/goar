package goar_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/daqiancode/goar"
	"github.com/daqiancode/goar/types"
	"github.com/daqiancode/goar/utils"
	"github.com/stretchr/testify/assert"
)

var testWallet *goar.Wallet
var err error

func init() {
	clientUrl := "https://arweave.net"
	testWallet, err = goar.NewWallet([]byte(`{ "kty": "RSA",
	"n":
	 "nQ9iy1fRM2xrgggjHhN1xZUnOkm9B4KFsJzH70v7uLMVyDqfyIJEVXeJ4Jhk_8KpjzYQ1kYfnCMjeXnhTUfY3PbeqY4PsK5nTje0uoOe1XGogeGAyKr6mVtKPhBku-aq1gz7LLRHndO2tvLRbLwX1931vNk94bSfJPYgMfU7OXxFXbTdKU38W6u9ShoaJGgUQI1GObd_sid1UVniCmu7P-99XPkixqyacsrkHzBajGz1S7jGmpQR669KWE9Z0unvH0KSHxAKoDD7Q7QZO7_4ujTBaIFwy_SJUxzVV8G33xvs7edmRdiqMdVK5W0LED9gbS4dv_aee9IxUJQqulSqZphPgShIiGNl9TcL5iUi9gc9cXR7ISyavos6VGiem_A-S-5f-_OKxoeZzvgAQda8sD6jtBTTuM5eLvgAbosbaSi7zFYCN7zeFdB72OfvCh72ZWSpBMH3dkdxsKCDmXUXvPdDLEnnRS87-MP5RV9Z6foq_YSEN5MFTMDdo4CpFGYl6mWTP6wUP8oM3Mpz3-_HotwSZEjASvWtiff2tc1fDHulVMYIutd52Fis_FKj6K1fzpiDYVA1W3cV4P28Q1-uF3CZ8nJEa5FXchB9lFrXB4HvsJVG6LPSt-y2R9parGi1_kEc6vOYIesKspgZ0hLyIKtqpTQFiPgKRlyUc-WEn5E",
	"e": "AQAB",
	"d":
	 "LIqif_yFrcm_q37XRr5KFiC4oUUsQKb5dx7fbLPlzXmsYb6OdfTLoFloVrOhYQ85uw2gNMRqToOAmgDAroQDspaoivlo5bhwP7R4orSVJP84xKzJMx-aNke3hGZtywQdytqfmQv_i3jxRm0Si33EXUnrWQVbEVmCEJ9kfgaIJ0NhALQ8TGx7dxv7cLp6U3zY0X2_PrsVkday5MFS45Wt4vHuYaGeBS4KFygHDflOlKiJ4FGksU3wzyBFO0o1tST21ayxd_G6sbdyar72sQU-asBvYU3kSVMuZs20i1C67qEizk1jqcdKbRuKRApqqs7ub8g2U6yDQaZYqft7KqC8Oi6X-VrUFNMttzvrTmJ0nhzGs1g7SJnZp1NJ1yvqwHOoisUjd0kyqBsDBzT_q0UsYOTM9yi3Ve021R7ghVQPGuvbDmRz8rCm_XusWRp2pzE1R14vot6G67w_2eXGNeD8GYjLKEX35DJwQqhqmD4AVR03u1slY7QK_EI2TxD2aCSpsBToX2Rix0E18dfPCkPmb4WAXccgyMKyR07w_31eypp1qdcQ4cNB_3-2QvM6kbzH3LJkRyPZzamiTB4pKoaQRDqKas0_QctmIdjAWNDr79q4cDcWB0de6X8H6QTjJoL28Lf6pvXXatLoHMAaVHWZpq-yz04-6oiA4pY6DFIs8AE",
	"p":
	 "yL6EzJ5NP0T4vUyWDjkGRI2JAqD9wJhaNd3jt-2Wip3aOQBcmWyagM2tRfbQunwWRquVGd6oSW-Tfr5aBesY3C_EKkbiXt7CSJ9iWz15FBWlLe5lSLbVi5tIcY7iV_D3iFND4mX39gmILX7p_sPmt6z5ZpHKv9MvypjuU2jKScC0HER5fgB-xJAhOQ45xfxuqjpDY6vVWUms6ZSTF3NOrLIrcupNdxUpWNjxKI6A48cvfKdAwJi-vtJcEZFrCJCv98UTgZFSgiWyUVcW2Bt9BYJQRKffm7cuC4aT_xwYq4sMcALoEq0CXIQ9XK4rRPH2bq6E2aZ2ll57OUBWR_1KEQ",
	"q":
	 "yEqoTJqr4EyWoxa-qz5x4FEylw9ueTqLmnZGBbRZaVwJnEY-ydoOaty9ZOtGA6g8EfLVGhtYy_B51479gpDTtygV846313WPDU3drTsPvGNmuA0dD5w4Vva-759Zmc4BcoJA7UYnlr4QRaLQuNAmf9RCgA-yTtHppt5Lre7tGG5K4j7uHRPzp3NhKPq8WosbydN0kFz-a2Vn9kZjoqI3m1JqvX9wRGFzOZlUnQTWA7MaGv4SgoqMMu8PyKQLjhbIJyvnIGI7NZjCi1uw0H03skdgC8bBEL_uSp6Go8nr1apGJ6o1x5_hriofgD9DvpKMkEsx4uTxRvxFys8w3jh9gQ",
	"dp":
	 "Il7kc_hit4OCpz62roa6-P_Wxplz-Qbc4z4zoClQzjkKxRm3wRkkNwuAMGt6_4MBeWYlaEGERNaSxW-oED1Zi1GuX6K1XZL8ZtzLRV34HiU6m-umcdXEKFwVAkR5op8Cctf21ouo8fpd05RYUiOOnEJEjXhG46MwGpsmqydVA124OOLMfnNtQRCAb7ls0OZQuFqzcRxZsij4LyIeMTSv8seqwsk1LD92Td0PJWeIz_cpvUkRwCgm-Jsh4mwojFXhmyWmGlgcbWYw6tZjder28_uE7Mxlb87kVlrbeiGAY9ax8Xe97nyq29ZUf0re47YeAINnAbELuuFAbeQDId5PUQ",
	"dq":
	 "woDVvUZ64OAfbRNaZ_vFJHxVr6K5uppjFcYDq-h-57UMVClXMjhCxf3FIqrjnAuVAi0aSzcBXVMTT4S5pUC1iOkxoAsZdu_f0qCqRF7VojG5f8SkUxN3FuSZeSP7JESM3UGmgYUeTuIV9TnujXr92CctyST1GFv7FiRLxAYBUzdQGzPXkn9cn2GJmf0cSqVKgA2L5eGY5HxeoCes_DOh4oD_zTRjttQXzHidVbprhr43_Lx9By46hf_oCQVdf0eaaYfV9HnQW_UT_7c0FtNy8fskR2tk87ofU3Fs-MPO9PhdFonRniEiTTr0ylslk3zHahzLvjZsJG457ICWSUb8gQ",
	"qi":
	 "ENtqTq3iiDkeiyPWD7pNRfiwIJnY5Zf97yXakxe04usHXWKmZulllttqsDkfHOXkBxRxHxqqTgOLuRpNsLrpI5MAxs8uSl13A70LUzHldnE8ePgt0688UpoI5Iw9oV2RdF_LvSrsgpa-SeexXxbZqXWpDNeUxYt2S327cS8HmrnETKy9z9VoVFmCT6_NCnxOaOTwr67dPBnGnW7nT3499m_aqmikCNjcmkfYihED6S2jZBRHPaSDM7JPPyQSEyRkGjR4z9JzhLOvbJf8tDKSE00JXJClmbpX-5qRcNt0gcJy6ceYQs-c94I24yGpunMMSwGo2i1-sGNwH1wj5-gv1Q" }`),
		clientUrl)
	if err != nil {
		panic(err)
	}
}

func TestPubKey(t *testing.T) {

	pubKey := testWallet.Signer.PubKey
	assert.Equal(t, "nQ9iy1fRM2xrgggjHhN1xZUnOkm9B4KFsJzH70v7uLMVyDqfyIJEVXeJ4Jhk_8KpjzYQ1kYfnCMjeXnhTUfY3PbeqY4PsK5nTje0uoOe1XGogeGAyKr6mVtKPhBku-aq1gz7LLRHndO2tvLRbLwX1931vNk94bSfJPYgMfU7OXxFXbTdKU38W6u9ShoaJGgUQI1GObd_sid1UVniCmu7P-99XPkixqyacsrkHzBajGz1S7jGmpQR669KWE9Z0unvH0KSHxAKoDD7Q7QZO7_4ujTBaIFwy_SJUxzVV8G33xvs7edmRdiqMdVK5W0LED9gbS4dv_aee9IxUJQqulSqZphPgShIiGNl9TcL5iUi9gc9cXR7ISyavos6VGiem_A-S-5f-_OKxoeZzvgAQda8sD6jtBTTuM5eLvgAbosbaSi7zFYCN7zeFdB72OfvCh72ZWSpBMH3dkdxsKCDmXUXvPdDLEnnRS87-MP5RV9Z6foq_YSEN5MFTMDdo4CpFGYl6mWTP6wUP8oM3Mpz3-_HotwSZEjASvWtiff2tc1fDHulVMYIutd52Fis_FKj6K1fzpiDYVA1W3cV4P28Q1-uF3CZ8nJEa5FXchB9lFrXB4HvsJVG6LPSt-y2R9parGi1_kEc6vOYIesKspgZ0hLyIKtqpTQFiPgKRlyUc-WEn5E", base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes()))
}

func TestAddress(t *testing.T) {
	addr := testWallet.Signer.Address
	assert.Equal(t, "eIgnDk4vSKPe0lYB6yhCHDV1dOw3JgYHGocfj7WGrjQ", addr)
}

// test sand ar without data
func TestWallet_SendAR(t *testing.T) {
	// arNode := "https://arweave.net"
	// w, err := NewWalletFromPath("./example/testKey.json", arNode) // your wallet private key
	// assert.NoError(t, err)
	//
	// target := "cSYOy8-p1QFenktkDBFyRM3cwZSTrQ_J4EsELLho_UE"
	// amount := big.NewFloat(0.001)
	// tags := []types.Tag{
	// 	{Name: "GOAR", Value: "sendAR"},
	// }
	// tx,  err := w.SendAR(amount, target, tags)
	// assert.NoError(t, err)
	// t.Logf("tx hash: %s \n", tx.ID)
}

// test send small size file
func TestWallet_SendDataSpeedUp01(t *testing.T) {
	// arNode := "https://arweave.net"
	// w, err := NewWalletFromPath("./example/testKey.json", arNode) // your wallet private key
	// assert.NoError(t, err)
	//
	// // data := []byte("aaa this is a goar test small size file data") // small file
	// data := make([]byte, 255*1024)
	// for i := 0; i < len(data); i++ {
	// 	data[i] = byte('b' + i)
	// }
	// tags := []types.Tag{
	// 	{Name: "GOAR", Value: "SMDT"},
	// }
	// tx, err := w.SendDataSpeedUp(data, tags, 50)
	// assert.NoError(t, err)
	// t.Logf("tx hash: %s", tx.ID)
}

// test send big size file
func TestWallet_SendDataSpeedUp02(t *testing.T) {
	// proxyUrl := "http://127.0.0.1:8001"
	// arNode := "https://arweave.net"
	// w, err := NewWalletFromPath("./wallet/account1.json", arNode, proxyUrl) // your wallet private key
	// assert.NoError(t, err)
	//
	// data, err := ioutil.ReadFile("/Users/sandyzhou/Downloads/abc.jpeg")
	// if err != nil {
	// 	panic(err)
	// }
	// tags := []types.Tag{
	// 	{Name: "Sender", Value: "Jie"},
	// 	{Name: "Data-Introduce", Value: "Happy anniversary, my google and dearest! I‘m so grateful to have you in my life. I love you to infinity and beyond! (⁎⁍̴̛ᴗ⁍̴̛⁎)"},
	// }
	// tx, err := w.SendDataSpeedUp(data, tags, 10)
	// assert.NoError(t, err)
	// t.Logf("tx hash: %s", tx.ID)
}

func Test_SendPstTransfer(t *testing.T) {
	// w, err := NewWalletFromPath("./wallet/account1.json","https://arweave.net")
	// assert.NoError(t, err)
	//
	// contractId := "usjm4PCxUd5mtaon7zc97-dt-3qf67yPyqgzLnLqk5A"
	// target := "Ii5wAMlLNz13n26nYY45mcZErwZLjICmYd46GZvn4ck"
	// qty := big.NewInt(1)
	// arTx, err := w.SendPst(contractId,target,qty,nil,50)
	// assert.NoError(t, err)
	// t.Log(arTx.ID)
}

// go test  -run ^TestUploaderFile$ github.com/daqiancode/goar -v
func TestUploaderFile(t *testing.T) {
	w, err := goar.NewWalletFromPath("arweave-key.json", "https://arweave.net")
	assert.NoError(t, err)
	t.Log(w.Signer.Address)

	file, err := os.Open("test_resources/test.mp4")
	if err != nil {
		panic(err)
	}
	stat, err := file.Stat()
	assert.Nil(t, err)
	fileSize := stat.Size()
	reward, err := w.Client.GetTransactionPrice(fileSize)
	assert.Nil(t, err)
	tx := goar.NewSendFileTransaction(file, fileSize, reward, types.Tag{Name: "Content-Type", Value: "video/mp4"})
	err = w.SignTransaction(tx)
	assert.Nil(t, err)

	uploader, err := goar.CreateUploader(w.Client, tx, file, fileSize)
	totalSent := 0
	lastTime := time.Now().Unix()
	lastTotal := 0
	callbackCount := 0
	uploader.ProgressCallback = func(bytesSent int) {
		callbackCount += 1
		totalSent += bytesSent
		fmt.Println(bytesSent, totalSent, fileSize)
		fmt.Println("progress: ", totalSent/int(fileSize))
		if callbackCount%10 == 0 {
			now := time.Now().Unix()
			duration := now - lastTime
			if duration > 0 {
				speed := (totalSent - lastTotal) / 1024 / int(duration)
				fmt.Print("speed: ", speed, "KB/s")
				lastTotal = totalSent
				lastTime = now
			}

		}

	}
	assert.Nil(t, err)
	uploader.ConcurrentOnce(context.Background(), 10)
	// err = uploader.Once()
	assert.Nil(t, err)
	fmt.Println("https://arweave.net/" + tx.ID)
}

func TestSendBytes(t *testing.T) {
	w, err := goar.NewWalletFromPath("arweave-key.json", "https://arweave.net")
	assert.NoError(t, err)
	t.Log(w.Signer.Address)

	file := utils.NewReadBuffer([]byte("hello world"))
	assert.Nil(t, err)
	reward, err := w.Client.GetTransactionPrice(int64(file.Len()))
	assert.Nil(t, err)
	tags := []types.Tag{{Name: "content-type", Value: "text/plain"}}

	tx := goar.NewSendFileTransaction(file, int64(file.Len()), reward, tags...)
	err = w.SignTransaction(tx)
	assert.Nil(t, err)
	uploader, err := goar.CreateUploader(w.Client, tx, file, int64(file.Len()))
	assert.Nil(t, err)
	err = uploader.Once()
	assert.Nil(t, err)
	fmt.Println("https://arweave.net/" + tx.ID)
}

func TestTransactionUploader_ConcurrentUploadChunks(t *testing.T) {
	w, err := goar.NewWalletFromPath("", "https://arweave.net")
	assert.NoError(t, err)
	t.Log(w.Signer.Address)
	signer01 := w.Signer
	// sig item01 by ecc signer
	itemSigner01, err := goar.NewItemSigner(signer01)
	assert.NoError(t, err)
	d1, err := ioutil.ReadFile("test.go")
	if err != nil {
		panic(err)
	}
	item01, err := itemSigner01.CreateAndSignItem(d1, "", "", []types.Tag{
		{Name: "Content-Type", Value: "text/html"},
		{Name: "Owner", Value: "Vv"},
	})
	assert.NoError(t, err)
	t.Log("item01", "id", item01.Id)

	d2, err := ioutil.ReadFile("/Users/sandyzhou/Downloads/2.jpeg")
	if err != nil {
		panic(err)
	}
	item02, err := itemSigner01.CreateAndSignItem(d2, "", "", []types.Tag{
		{Name: "Content-Type", Value: "image/jpeg"},
		{Name: "Owner", Value: "Vv"},
	})
	assert.NoError(t, err)
	t.Log("item02", "id", item02.Id)

	d3, err := ioutil.ReadFile("/Users/sandyzhou/Downloads/3.jpeg")
	if err != nil {
		panic(err)
	}
	item03, err := itemSigner01.CreateAndSignItem(d3, "", "", []types.Tag{
		{Name: "Content-Type", Value: "image/jpeg"},
		{Name: "Owner", Value: "Vv"},
	})
	assert.NoError(t, err)
	t.Log("item03", "id", item03.Id)

	d4, err := ioutil.ReadFile("/Users/sandyzhou/Downloads/4.jpeg")
	if err != nil {
		panic(err)
	}
	item04, err := itemSigner01.CreateAndSignItem(d4, "", "", []types.Tag{
		{Name: "Content-Type", Value: "image/jpeg"},
		{Name: "Owner", Value: "Vv"},
	})
	assert.NoError(t, err)

	d5, err := ioutil.ReadFile("/Users/sandyzhou/Downloads/5.jpeg")
	if err != nil {
		panic(err)
	}
	item05, err := itemSigner01.CreateAndSignItem(d5, "", "", []types.Tag{
		{Name: "Content-Type", Value: "image/jpeg"},
		{Name: "Owner", Value: "Vv"},
	})
	assert.NoError(t, err)
	t.Log("item05", "id", item05.Id)

	d6, err := ioutil.ReadFile("/Users/sandyzhou/Downloads/6.jpeg")
	if err != nil {
		panic(err)
	}
	item06, err := itemSigner01.CreateAndSignItem(d6, "", "", []types.Tag{
		{Name: "Content-Type", Value: "image/jpeg"},
		{Name: "Owner", Value: "Vv"},
	})
	assert.NoError(t, err)
	t.Log("item06", "id", item06.Id)

	d7, err := ioutil.ReadFile("/Users/sandyzhou/Downloads/7.jpeg")
	if err != nil {
		panic(err)
	}
	item07, err := itemSigner01.CreateAndSignItem(d7, "", "", []types.Tag{
		{Name: "Content-Type", Value: "image/jpeg"},
		{Name: "Owner", Value: "Vv"},
	})
	assert.NoError(t, err)
	t.Log("item07", "id", item07.Id)

	d8, err := ioutil.ReadFile("/Users/sandyzhou/Downloads/8.jpeg")
	if err != nil {
		panic(err)
	}
	item08, err := itemSigner01.CreateAndSignItem(d8, "", "", []types.Tag{
		{Name: "Content-Type", Value: "image/jpeg"},
		{Name: "Owner", Value: "Vv"},
	})
	assert.NoError(t, err)
	t.Log("item08", "id", item08.Id)

	d9, err := ioutil.ReadFile("/Users/sandyzhou/Downloads/9.jpeg")
	if err != nil {
		panic(err)
	}
	item09, err := itemSigner01.CreateAndSignItem(d9, "", "", []types.Tag{
		{Name: "Content-Type", Value: "image/jpeg"},
		{Name: "Owner", Value: "Vv"},
	})
	assert.NoError(t, err)
	t.Log("item09", "id", item09.Id)

	d10, err := ioutil.ReadFile("/Users/sandyzhou/Downloads/10.jpeg")
	if err != nil {
		panic(err)
	}
	item10, err := itemSigner01.CreateAndSignItem(d10, "", "", []types.Tag{
		{Name: "Content-Type", Value: "image/jpeg"},
		{Name: "Owner", Value: "Vv"},
	})
	assert.NoError(t, err)
	t.Log("item10", "id", item10.Id)

	d11, err := ioutil.ReadFile("/Users/sandyzhou/Downloads/11.jpeg")
	if err != nil {
		panic(err)
	}
	item11, err := itemSigner01.CreateAndSignItem(d11, "", "", []types.Tag{
		{Name: "Content-Type", Value: "image/jpeg"},
		{Name: "Owner", Value: "Vv"},
	})
	assert.NoError(t, err)
	t.Log("item11", "id", item11.Id)

	// assemble bundle
	bundle, err := utils.NewBundle(item01, item02, item03, item04, item05, item06, item07, item08, item09, item10, item11)
	assert.NoError(t, err)

	t.Log(len(bundle.BundleBinary))
	// send to arweave
	// ctx ,cancel := context.WithTimeout(context.Background(),100*time.Millisecond)
	// defer cancel()
	// tx, err := w.SendBundleTx(ctx, 0,bundle.BundleBinary, []types.Tag{
	// 	{Name: "APP", Value: "Goar"},
	// 	{Name: "Protocol-Name", Value: "BAR"},
	// 	{Name: "Action", Value: "Burn"},
	// 	{Name: "App-Name", Value: "SmartWeaveAction"},
	// 	{Name: "App-Version", Value: "0.3.0"},
	// 	{Name: "Input", Value: `{"function":"mint"}`},
	// 	{Name: "Contract", Value: "VFr3Bk-uM-motpNNkkFg4lNW1BMmSfzqsVO551Ho4hA"},
	// })
	// assert.NoError(t, err)
	// t.Log(tx.ID)
}
