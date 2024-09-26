package example

import (
	"context"
	"testing"

	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"

	"github.com/everFinance/goar"
	"github.com/stretchr/testify/assert"
)

func Test_Client(t *testing.T) {
	// create client
	arNode := "https://arweave.net"
	c := goar.NewClient(arNode)
	txId := "hKMMPNh_emBf8v_at1tFzNYACisyMQNcKzeeE1QE9p8"

	// 1. getInfo
	nodeInfo, err := c.GetInfo()
	assert.NoError(t, err)
	t.Logf("%v", nodeInfo)

	// 2. full transaction via Id
	tx, err := c.GetTransactionByID(txId)
	assert.NoError(t, err)
	t.Log(tx)

	// 3. get transaction field by id
	f, err := c.GetTransactionField(txId, "signature")
	assert.NoError(t, err)
	t.Log(f)

	// 4. get transaction data
	data, err := c.GetTransactionData(txId)
	assert.NoError(t, err)
	t.Log(string(data))
	data, err = c.GetTransactionData(txId, "html")
	assert.NoError(t, err)
	t.Log(string(data))

	// 5. get tx send current time reward
	reward, err := c.GetTransactionPrice(len(data), nil)
	assert.NoError(t, err)
	t.Log(reward)
	to := "1seRanklLU_1VTGkEk7P0xAwMJfA7owA1JHW5KyZKlY"
	reward, err = c.GetTransactionPrice(0, &to)
	assert.NoError(t, err)
	t.Log(reward)

	// 6. get anchor
	anchor, err := c.GetTransactionAnchor()
	assert.NoError(t, err)
	t.Log(anchor)

}

func Test_Arq1(t *testing.T) {
	arNode := "https://arweave.net"
	c := goar.NewClient(arNode)
	dd, err := c.GetTransactionPrice(712, nil)
	assert.NoError(t, err)
	t.Log(dd)
}

func Test_Arq(t *testing.T) {
	ctx := context.Background()
	arqStr := `{
			"op": "and",
			"expr1": {
				"op": "equals",
				"expr1": "TokenSymbol",
				"expr2": "DXN"
			},
			"expr2": {
				"op": "equals",
				"expr1": "CreatedBy",
				"expr2": "zhou yu ji"
			}
		}`
	// create client
	arNode := "https://arweave.net"
	c := goar.NewClient(arNode)
	ids, err := c.Arql(ctx, arqStr)
	t.Log(len(ids))
	assert.NoError(t, err)
	sstr := make([]string, 0)
	sstr = append(sstr, "a")
	for _, val := range ids {
		t.Log(val)
	}

	for i := 0; i < len(ids); i++ {
		for j := 1; j < len(ids)-i; j++ {
			if ids[j] > ids[j-1] {
				ids[j], ids[j-1] = ids[j-1], ids[j]
			}
		}
	}
	t.Log(ids)
}

func Test_SendFormatTx(t *testing.T) {
	// arNode := "https://arweave.net"
	// wallet, err := goar.NewWalletFromPath("./testKey.json", arNode)
	// assert.NoError(t, err)
	//
	// owner := utils.Base64Encode(wallet.PubKey.N.Bytes())
	//
	// target := "cSYOy8-p1QFenktkDBFyRM3cwZSTrQ_J4EsELLho_UE"
	// reward, err := wallet.Client.GetTransactionPrice(nil, &target)
	// assert.NoError(t, err)
	//
	// anchor, err := wallet.Client.GetTransactionAnchor()
	// assert.NoError(t, err)
	//
	// amount := big.NewInt(140000) // transfer amount
	// tags := []types.Tag{{Name: "Content-Type", Value: "application/json"}, {Name: "tcrsa", Value: "sandyTest"}}
	// tx := &types.Transaction{
	// 	Format:    1,
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
	//
	// sig, err := utils.sign(signData, wallet.PrvKey)
	// assert.NoError(t, err)
	// tx.AddSignature(sig)
	//
	// status, code, err := wallet.Client.SubmitTransaction(tx)
	// assert.NoError(t, err)
	// t.Log(status, code)
	// t.Log("from: ",wallet.Address)
	// t.Log("txHash: ", tx.ID)
}

func Test_SendMsg(t *testing.T) {
	signer02, err := goar.NewSignerFromPath("./testKey.json")
	assert.NoError(t, err)
	t.Log(signer02.Address)
	itemSigner02, err := goar.NewItemSigner(signer02)
	assert.NoError(t, err)

	defaultTags := []types.Tag{
		{Name: "Data-Protocol", Value: "ao"},
		{Name: "Variant", Value: "ao.TN.1"},
		{Name: "Type", Value: "Message"},
		{Name: "SDK", Value: "argo"},
	}
	tags := append(defaultTags, []types.Tag{
		{Name: "Action", Value: "Transfer"},
		{Name: "Recipient", Value: "AVm7zHYGzw9PmUXv4CFHgiK3QtbRTpBMS2VOxJAK-F4"},
		{Name: "Quantity", Value: "22"},
	}...)
	target := "7En3PZJ0BBROTuSazQ9yZmeL1ThwsVqk616gY3DMFZU"
	data := []byte("1234")
	item02, err := itemSigner02.CreateAndSignItem(data, target, "", tags)
	assert.NoError(t, err)
	muUrl := "https://mu.ao-testnet.xyz"
	t.Log("item", "id", item02.Id)
	resp, err := utils.SubmitItemToMU(item02, muUrl)
	assert.NoError(t, err)
	t.Log(string(resp))
}
