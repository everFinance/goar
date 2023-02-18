package example

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/daqiancode/goar"
	"github.com/daqiancode/goar/types"
	"github.com/daqiancode/goar/utils"
	"github.com/everFinance/everpay-go/sdk"
	"github.com/everFinance/goether"
	"github.com/stretchr/testify/assert"
)

var (
	signer01 *goether.Signer
	signer02 *goar.Signer
)

func init() {
	var err error
	signer01, err = goether.NewSigner("1f534ac18009182c07d266fe4a7903c0bcc8a66190f0967b719b2b3974a69c2f")
	if err != nil {
		return
	}

	rsaKey := `{ "kty": "RSA",
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
	 "ENtqTq3iiDkeiyPWD7pNRfiwIJnY5Zf97yXakxe04usHXWKmZulllttqsDkfHOXkBxRxHxqqTgOLuRpNsLrpI5MAxs8uSl13A70LUzHldnE8ePgt0688UpoI5Iw9oV2RdF_LvSrsgpa-SeexXxbZqXWpDNeUxYt2S327cS8HmrnETKy9z9VoVFmCT6_NCnxOaOTwr67dPBnGnW7nT3499m_aqmikCNjcmkfYihED6S2jZBRHPaSDM7JPPyQSEyRkGjR4z9JzhLOvbJf8tDKSE00JXJClmbpX-5qRcNt0gcJy6ceYQs-c94I24yGpunMMSwGo2i1-sGNwH1wj5-gv1Q" }`

	signer02, err = goar.NewSigner([]byte(rsaKey))
	if err != nil {
		return
	}
}

func TestBundleToArweave(t *testing.T) {
	// sig item01 by ecc signer
	itemSigner01, err := goar.NewItemSigner(signer01)
	assert.NoError(t, err)
	item01, err := itemSigner01.CreateAndSignItem([]byte("aa bb cc"), "", "", []types.Tag{
		{Name: "Content-Type", Value: "application/txt"},
		{Name: "App-Version", Value: "2.0.0"},
	})
	assert.NoError(t, err)

	// sig item02 by rsa signer
	itemSigner02, err := goar.NewItemSigner(signer02)
	assert.NoError(t, err)
	item02, err := itemSigner02.CreateAndSignItem([]byte("dd ee ff"), "", "", []types.Tag{
		{Name: "Content-Type", Value: "application/txt"},
	})
	assert.NoError(t, err)

	// assemble bundle
	bundle, err := utils.NewBundle(item01, item02)
	assert.NoError(t, err)

	// send to arweave
	wal, err := goar.NewWalletFromPath("jwkKey.json", "https://arweave.net")
	assert.NoError(t, err)
	tx, err := wal.SendBundleTx(context.TODO(), 0, bundle.BundleBinary, []types.Tag{
		{Name: "App", Value: "goar"},
	})
	assert.NoError(t, err)
	t.Log(tx.ID)
}

func TestBundleItemToArseeding(t *testing.T) {
	// sig item01 by ecc signer
	itemSigner01, err := goar.NewItemSigner(signer01)
	assert.NoError(t, err)
	item01, err := itemSigner01.CreateAndSignItem([]byte("aa bb cc"), "", "", []types.Tag{
		{Name: "Content-Type", Value: "application/txt"},
		{Name: "App-Version", Value: "2.0.0"},
	})
	assert.NoError(t, err)

	// submit item to arseeding
	resp, err := utils.SubmitItemToArSeed(item01, "USDC", "https://seed.everpay.io")
	assert.NoError(t, err)
	t.Log("itemId: ", resp.ItemId)

	// payment fee
	paySdk, err := sdk.New(signer01, "https://api.everpay.io")
	assert.NoError(t, err)
	tokenSymbol := resp.Currency
	amount, _ := new(big.Int).SetString(resp.Fee, 10)
	to := resp.Bundler
	everTx, err := paySdk.Transfer(tokenSymbol, amount, to, "")
	assert.NoError(t, err)
	t.Log("paymentId: ", everTx.HexHash())
}

func TestBundleItemToBundlr(t *testing.T) {
	itemSigner01, err := goar.NewItemSigner(signer01)
	assert.NoError(t, err)
	item01, err := itemSigner01.CreateAndSignItem([]byte("aa bb cc"), "", "", []types.Tag{
		{Name: "Content-Type", Value: "application/txt"},
		{Name: "App-Version", Value: "2.0.0"},
	})
	assert.NoError(t, err)

	// submit item to bundlr
	resp, err := utils.SubmitItemToBundlr(item01, "https://node1.bundlr.network")
	assert.NoError(t, err)
	t.Log("itemId: ", resp.Id)
}

func TestVerifyBundleItem(t *testing.T) {
	cli := goar.NewClient("https://arweave.net")
	// id := "K0JskpURZ-zZ7m01txR7hArvsBDDi08S6-6YIVQoc_Y" // big size data
	// id := "mTm5-TtpsfJvUCPXflFe-P7HO6kOy4E2pGbt6-DUs40"

	// goar test tx
	// id := "ipVFFrAkLosTtk-M3J6wYq3MKpfE6zK75nMIC-oLVXw"
	// id := "2ZFhlTJlFbj8XVmBtnBHS-y6Clg68trcRgIKBNemTM8"
	// id := "WNGKdWsGqyhh7Y4vMcQL0GHFzNiyeqASIJn-Z1IjJE0"
	id := "lt24bnUGms5XLZeVamSPHePl4M2ClpLQyRxZI7weH1k"
	data, err := cli.DownloadChunkData(id)
	assert.NoError(t, err)
	bd, err := utils.DecodeBundle(data)
	assert.NoError(t, err)
	for _, item := range bd.Items {
		err = utils.VerifyBundleItem(item)
		assert.NoError(t, err)
	}
}

func TestDecodeBundleItemStream(t *testing.T) {
	cli := goar.NewClient("https://arweave.net")

	id := "lt24bnUGms5XLZeVamSPHePl4M2ClpLQyRxZI7weH1k"
	data, err := cli.DownloadChunkData(id)
	assert.NoError(t, err)
	bd, err := utils.DecodeBundle(data)
	assert.NoError(t, err)
	item := bd.Items[0]
	err = os.WriteFile("test.item", item.ItemBinary, 0644)
	assert.NoError(t, err)
	itemReader, err := os.Open("test.item")
	defer itemReader.Close()
	assert.NoError(t, err)
	item2, err := utils.DecodeBundleItemStream(itemReader)
	assert.NoError(t, err)
	assert.Equal(t, item2.Id, item.Id)
	bundleData, err := io.ReadAll(item2.DataReader)
	assert.NoError(t, err)
	assert.Equal(t, item.Data, utils.Base64Encode(bundleData))
}

func TestIOBuffer(t *testing.T) {
	itemReader, err := os.Open("test.item")
	assert.NoError(t, err)
	assert.Equal(t, itemReader.Name(), "test.item")
	itemBinary, err := ioutil.ReadFile("test.item")
	assert.NoError(t, err)
	defer itemReader.Close()
	item, err := utils.DecodeBundleItemStream(itemReader)
	assert.NoError(t, err)
	_, err = item.DataReader.Seek(0, 0)
	assert.NoError(t, err)
	itemBinary2, err := io.ReadAll(item.DataReader)
	assert.NoError(t, err)
	assert.Equal(t, item.DataReader.Name(), "test.item")
	assert.Equal(t, itemBinary2, itemBinary)
}

func TestCreateBundleItemStream(t *testing.T) {
	data0, err := ioutil.ReadFile("../go.mod")
	assert.NoError(t, err)
	data1, err := os.Open("../go.mod")
	assert.NoError(t, err)
	itemSigner, err := goar.NewItemSigner(signer01)
	assert.NoError(t, err)
	item0, err := itemSigner.CreateAndSignItem(data0, "", "", []types.Tag{
		{Name: "Content-Type", Value: "application/txt"},
		{Name: "App-Version", Value: "2.0.0"},
	})
	assert.NoError(t, err)
	item1, err := itemSigner.CreateAndSignItemStream(data1, "", "", []types.Tag{
		{Name: "Content-Type", Value: "application/txt"},
		{Name: "App-Version", Value: "2.0.0"},
	})
	assert.Equal(t, item0.Signature, item1.Signature)
}

func TestPointerAndValueCopy(t *testing.T) {
	data0, err := ioutil.ReadFile("../go.mod")
	assert.NoError(t, err)
	data, err := os.Open("../go.mod")
	assert.NoError(t, err)
	itemSigner, err := goar.NewItemSigner(signer01)
	assert.NoError(t, err)
	item, err := itemSigner.CreateAndSignItemStream(data, "", "", []types.Tag{
		{Name: "Content-Type", Value: "application/txt"},
		{Name: "App-Version", Value: "2.0.0"},
	})
	assert.NoError(t, err)
	err = changeData(item)
	assert.NoError(t, err)
	data2, err := io.ReadAll(item.DataReader)
	assert.NoError(t, err)
	assert.NotEqual(t, data0, data2)
	assert.Equal(t, data0[1:], data2)
	n, err := item.DataReader.Seek(0, 2)
	assert.Equal(t, int(n), len(data0))
}

func TestReader(t *testing.T) {
	data, err := os.Open("../go.mod")
	defer data.Close()
	assert.NoError(t, err)
	data2, err := ioutil.ReadFile("../go.mod")
	itemSigner, err := goar.NewItemSigner(signer01)
	assert.NoError(t, err)
	item, err := itemSigner.CreateAndSignItemStream(data, "", "", []types.Tag{
		{Name: "Content-Type", Value: "application/txt"},
		{Name: "App-Version", Value: "2.0.0"},
	})
	assert.NoError(t, err)
	item2, err := itemSigner.CreateAndSignItem(data2, "", "", []types.Tag{
		{Name: "Content-Type", Value: "application/txt"},
		{Name: "App-Version", Value: "2.0.0"},
	})
	assert.NoError(t, err)

	binary, err := io.ReadAll(item.BinaryReader)
	assert.NoError(t, err)

	assert.Equal(t, binary, item2.ItemBinary)
}

func TestVerifyItemStream(t *testing.T) {
	itemReader, err := os.Open("test.item")
	defer itemReader.Close()
	assert.NoError(t, err)
	item, err := utils.DecodeBundleItemStream(itemReader)
	assert.NoError(t, err)
	err = utils.VerifyBundleItem(*item)
	assert.NoError(t, err)
}

func changeData(item types.BundleItem) error {
	b := make([]byte, 1)
	n, err := item.DataReader.Read(b)
	if n < 1 || err != nil {
		return errors.New("err")
	}
	return nil
}
