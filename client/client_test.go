package client

import (
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

// import (
// 	"fmt"
// 	"testing"

// func TestGetTransactionByID(t *testing.T) {
// 	client := New("https://arweave.net")
// 	fmt.Println(client.GetTransactionByID("FgcKlptyDXSgEonYfy5cNBimq7GJ4h8h6L6pxuuYOBc"))
// }

// func TestGetTransactionPrice(t *testing.T) {
// 	client := New("https://arweave.net")
// 	target := ""
// 	reward, err := client.GetTransactionPrice([]byte("123"), &target)
// 	assert.NoError(t, err)
// 	fmt.Println(reward)
// }

// func TestGetLastTransactionID(t *testing.T) {
// 	client := New("https://arweave.net")
// 	lastTx, err := client.GetLastTransactionID("dQzTM9hXV5MD1fRniOKI3MvPF_-8b2XDLmpfcMN9hi8")
// 	assert.NoError(t, err)
// 	fmt.Println(lastTx)
// }

// func TestGetTransactionAnchor(t *testing.T) {
// 	client := New("https://arweave.net")
// 	fmt.Println(client.GetTransactionAnchor())
// }

// func TestSubmitTransaction(t *testing.T) {
// 	client := New("https://arweave.net")
// 	fmt.Println(
// 		client.SubmitTransaction(&types.Transaction{
// 			ID: "n1iKT3trKn6Uvd1d8XyOqKBy8r-8SSBtGA62m3puK5k",
// 		}),
// 	)
// }

// func TestArql(t *testing.T) {
// 	client := New("https://arweave.net")
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
// 	client := New("https://arweave.net")
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
// 	client := New("https://arweave.net")
// 	fmt.Println(
// 		client.GetWalletBalance("dQzTM9hXV5MD1fRniOKI3MvPF_-8b2XDLmpfcMN9hi8"),
// 	)
// }

func TestClient_DownloadChunkData(t *testing.T) {
	// client := New("https://arweave.net")
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
	client := New("https://arweave.net")
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
	client := New("https://arweave.net")
	tx, status, code, err := client.GetTransactionByID(txId)
	assert.NoError(t, err)
	t.Log(status, code)
	t.Log(tx.Format)
	t.Log(types.TagsDecode(tx.Tags))
	err = tx.VerifyTransaction()
	assert.NoError(t, err)
}
