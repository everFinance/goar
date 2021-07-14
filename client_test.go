package goar

import (
	"testing"

	"github.com/everFinance/goar/utils"
	"github.com/stretchr/testify/assert"
)

// import (
// 	"fmt"
// 	"testing"

// func TestGetTransactionByID(t *testing.T) {
// 	client := NewClient("https://arweave.net")
// 	fmt.Println(client.GetTransactionByID("FgcKlptyDXSgEonYfy5cNBimq7GJ4h8h6L6pxuuYOBc"))
// }

// func TestGetTransactionPrice(t *testing.T) {
// 	client := NewClient("https://arweave.net")
// 	target := ""
// 	reward, err := client.GetTransactionPrice([]byte("123"), &target)
// 	assert.NoError(t, err)
// 	fmt.Println(reward)
// }

// func TestGetLastTransactionID(t *testing.T) {
// 	client := NewClient("https://arweave.net")
// 	lastTx, err := client.GetLastTransactionID("dQzTM9hXV5MD1fRniOKI3MvPF_-8b2XDLmpfcMN9hi8")
// 	assert.NoError(t, err)
// 	fmt.Println(lastTx)
// }

// func TestGetTransactionAnchor(t *testing.T) {
// 	client := NewClient("https://arweave.net")
// 	fmt.Println(client.GetTransactionAnchor())
// }

// func TestSubmitTransaction(t *testing.T) {
// 	client := NewClient("https://arweave.net")
// 	fmt.Println(
// 		client.SubmitTransaction(&types.Transaction{
// 			ID: "n1iKT3trKn6Uvd1d8XyOqKBy8r-8SSBtGA62m3puK5k",
// 		}),
// 	)
// }

// func TestArql(t *testing.T) {
// 	client := NewClient("https://arweave.net")
// 	fmt.Println(
// 		client.Arql(`
// 		{
// 			"op": "and",
// 			"expr1": {
// 				"op": "equals",
// 				"expr1": "TokenSymbol",
// 				"expr2": "ROL"
// 			},
// 			"expr2": {
// 				"op": "equals",
// 				"expr1": "CreatedBy",
// 				"expr2": "dQzTM9hXV5MD1fRniOKI3MvPF_-8b2XDLmpfcMN9hi8"
// 			}
// 		}
// 		`),
// 	)
// }

// func TestGraphQL(t *testing.T) {
// 	client := NewClient("https://arweave.net")
// 	data, err := client.GraphQL(`
// 	{
// 		transactions(
// 			tags: [
// 					{
// 							name: "TokenSymbol",
// 							values: "ROL"
// 					},
// 			]
// 			sort: HEIGHT_ASC
// 		) {
// 			edges {
// 				node {
// 					id
// 					tags {
// 						name
// 						value
// 					}
// 				}
// 			}
// 		}
// 	}`)
// 	assert.NoError(t, err)
// 	t.Log(string(data))
// }

// func TestGetWalletBalance(t *testing.T) {
// 	client := NewClient("https://arweave.net")
// 	fmt.Println(
// 		client.GetWalletBalance("dQzTM9hXV5MD1fRniOKI3MvPF_-8b2XDLmpfcMN9hi8"),
// 	)
// }

func TestClient_DownloadChunkData(t *testing.T) {
	// client := NewClient("https://arweave.net")
	// id := "ybEmme6TE3JKwnSYciPCjnAINwi_CWthomsxBes-kYk"
	// data, err := client.GetTransactionData(id, "jpg")
	// assert.NoError(t, err)
	//
	// t.Log(len(data))
	// err = ioutil.WriteFile("photo.jpg", data, 0777)
	// assert.NoError(t, err)
}

func TestClient_GetTransactionData(t *testing.T) {
	// proxy := "http://127.0.0.1:8001"
	client := NewClient("https://arweave.net")
	id := "lSHWbAfjJsK0so08BTTmHO_n809fGW2DYOMySsXHNuI"
	data, err := client.GetTransactionData(id, "json")
	if err != nil {
		t.Log(err.Error())
	}

	t.Log(string(data))
}

func TestNew(t *testing.T) {
	data := []byte("this is a goar test small size file data")
	a := utils.Base64Encode(data)
	t.Log(a)
}

func TestClient_VerifyTx(t *testing.T) {
	// txId := "XOzxw5kaYJrt9Vljj23pA5_6b63kY2ydQ0lPfnhksMA"
	txId := "_fVj-WyEtXV3URXlNkSnHVGupl7_DM1UWZ64WMdhPkU"
	client := NewClient("https://arweave.net")
	tx, err := client.GetTransactionByID(txId)
	assert.NoError(t, err)
	t.Log(tx.Format)
	t.Log(utils.TagsDecode(tx.Tags))
	err = utils.VerifyTransaction(*tx)
	assert.NoError(t, err)
}

func TestGetTransaction(t *testing.T) {
	arNode := "https://arweave.net"
	cli := NewClient(arNode)

	// on chain tx
	txId := "ggt-x5Q_niHifdNzMxZrhiibKf0KQ-cJun0UIBBa-yA"
	txStatus, err := cli.GetTransactionStatus(txId)
	assert.NoError(t, err)
	assert.Equal(t, 575660, txStatus.BlockHeight)
	tx, err := cli.GetTransactionByID(txId)
	assert.NoError(t, err)
	assert.Equal(t, "0pu7-Otb-AH6SSSX_rfUmpTkwh3Nmhpztd_IT8nYXDwBE6P3B-eJSBuaTBeLypx4", tx.LastTx)

	// not exist tx
	txId = "KPlEyCrcs2rDHBFn2f0UUn2NZQKfawGb_EnBfip8ayA"
	txStatus, err = cli.GetTransactionStatus(txId)
	assert.Equal(t, `{"status":404,"error":"Not Found"}`, err.Error())
	assert.Nil(t, txStatus)
	tx, err = cli.GetTransactionByID(txId)
	assert.Equal(t, `{"status":404,"error":"Not Found"}`, err.Error())
	assert.Nil(t, tx)

	// // pending tx
	// txId = "muANv_lsyZKC5C8fTxQaC2dCCyGDao8z35ECuGdIBP8" // need send a new tx create pending status
	// txStatus, err = cli.GetTransactionStatus(txId)
	// assert.Equal(t, "Pending",err.Error())
	// assert.Nil(t, txStatus)
	// tx, err = cli.GetTransactionByID(txId)
	// assert.Equal(t, "Pending",err.Error())
	// assert.Nil(t, txStatus)
}
