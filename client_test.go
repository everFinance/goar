package goar

import (
	"testing"

	"github.com/daqiancode/goar/utils"
	"github.com/stretchr/testify/assert"
)

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

func TestClient_Arql(t *testing.T) {
	// client := NewClient("https://arweave.dev")
	// id := "PvLGaQzn9MOwucO91uuMGRnq8pj1qlwbURPqhmW0UiM"
	//
	// status, err := client.GetTransactionStatus(id)
	// assert.NoError(t, err)
	// t.Log(status)
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
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, txStatus)
	tx, err = cli.GetTransactionByID(txId)
	assert.Equal(t, ErrNotFound, err)
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

func TestClient_GetTransactionTags(t *testing.T) {
	arNode := "https://arweave.net"
	cli := NewClient(arNode)
	id := "gdXUJuj9EZm99TmeES7zRHCJtnJoP3XgYo_7KJNV8Vw"
	tags, err := cli.GetTransactionTags(id)
	assert.NoError(t, err)
	assert.Equal(t, "App", tags[0].Name)
	assert.Equal(t, "Version", tags[1].Name)
	assert.Equal(t, "Owner", tags[2].Name)
}

func TestClient_GetBlockByHeight(t *testing.T) {
	arNode := "https://arweave.net"
	cli := NewClient(arNode)
	block, err := cli.GetBlockByHeight(793791)
	assert.NoError(t, err)
	assert.Equal(t, "ci2uJhYmdldgkHbScDClCwAA0eqn7dCduAEpLfRorSA", block.Nonce)
}

func TestClient_GetTransactionDataByGateway(t *testing.T) {
	arNode := "https://arweave.net"
	cli := NewClient(arNode)
	id := "3S44SVxPWAqtadjehWR3bW1gP4B6Qsii4bnx9yz0_0s"
	data, err := cli.GetTransactionDataByGateway(id)
	assert.NoError(t, err)
	t.Log(len(data))
}

func TestClient_GetPeers(t *testing.T) {
	arNode := "https://arweave.net"
	cli := NewClient(arNode)
	peers, err := cli.GetPeers()
	assert.NoError(t, err)
	t.Log(len(peers))
}

func Test_GetTxDataFromPeers(t *testing.T) {
	cli := NewClient("https://arweave.net")
	txId := "J5FY1Ovd6JJ49WFHfCf-1wDM1TbaPSdKnGIB_8ePErE"
	data, err := cli.GetTxDataFromPeers(txId)

	assert.NoError(t, err)

	assert.NoError(t, err)
	t.Log(len(data))

	// verify data root
	chunks := utils.GenerateChunksBytes(data)
	dataRoot := utils.Base64Encode(chunks.DataRoot)
	tx, err := cli.GetTransactionByID(txId)
	assert.NoError(t, err)
	assert.Equal(t, tx.DataRoot, dataRoot)
}

// func TestClient_BroadcastData(t *testing.T) {
// 	cli := NewClient("https://arweave.net")
// 	txId := "J5FY1Ovd6JJ49WFHfCf-1wDM1TbaPSdKnGIB_8ePErE"
// 	data, err := cli.GetTransactionData(txId, "json")
// 	assert.NoError(t, err)

// 	err = cli.BroadcastData(txId, data, 20)
// 	assert.NoError(t, err)
// }

func TestClient_GetBlockFromPeers(t *testing.T) {
	cli := NewClient("https://arweave.net")
	block, err := cli.GetBlockFromPeers(793755)
	assert.NoError(t, err)
	t.Log(block.Txs)
}

func TestClient_GetTxFromPeers(t *testing.T) {
	cli := NewClient("https://arweave.net")
	arId := "5MiJDf2gFh4w3RXs1iXRrM9V8UwtnxX6xFATgxUqUN4"
	tx, err := cli.GetTxFromPeers(arId)
	assert.NoError(t, err)
	t.Log(tx)
}

func TestClient_GetUnconfirmedTx(t *testing.T) {
	cli := NewClient("https://arweave.net")
	arId := "5MiJDf2gFh4w3RXs1iXRrM9V8UwtnxX6xFATgxUqUN4"
	tx, err := cli.GetUnconfirmedTx(arId)
	assert.NoError(t, err)
	t.Log(tx)
}

func TestClient_GetUnconfirmedTxFromPeers(t *testing.T) {
	cli := NewClient("https://arweave.net")
	arId := "5MiJDf2gFh4w3RXs1iXRrM9V8UwtnxX6xFATgxUqUN4"
	tx, err := cli.GetUnconfirmedTxFromPeers(arId)
	assert.NoError(t, err)
	t.Log(tx)
}

func TestNewClient(t *testing.T) {
	cli := NewClient("https://arweave.net")
	res, err := cli.GetPendingTxIds()
	assert.NoError(t, err)
	t.Log("pending tx number:", len(res))
}

func TestNewTempConn(t *testing.T) {
	c := NewClient("https://arweave.net")
	peers, err := c.GetPeers()
	assert.NoError(t, err)
	pNode := NewTempConn()
	for _, peer := range peers {
		pNode.SetTempConnUrl("http://" + peer)
		offset, err := pNode.getTransactionOffset("pEYudvF0HjIU-2vKdhNZ9Dgr_bueucXaeRrbPhI90ew")
		if err != nil {
			t.Log("err", err, "perr", peer)
			continue
		}
		t.Logf("offset: %s, peer: %s", offset.Offset, peer)
	}
}

func TestClient_GetBlockHashList(t *testing.T) {
	c := NewClient("https://arweave.net")
	from := 1095730
	to := 1095750
	list, err := c.GetBlockHashList(from, to)
	assert.NoError(t, err)
	t.Log(list)
}

func TestClient_GetBlockHashList2(t *testing.T) {
	// c := NewClient("https://arweave.net")
	// peers, err := c.GetPeers()
	// assert.NoError(t, err)
	// pNode := NewTempConn()
	// for _, peer := range peers {
	// 	pNode.SetTempConnUrl("http://" + peer)
	// 	from := 1095740
	// 	to := 1095750
	// 	list, err := c.GetBlockHashList(from, to)
	// 	if err != nil {
	// 		t.Log("err", err, "perr", peer)
	// 		continue
	// 	}
	// 	t.Log(peer)
	// 	t.Log(list)
	// }
}

// func TestClient_ConcurrentDownloadChunkData(t *testing.T) {
// 	c := NewClient("https://arweave.net")
// 	arId := "trMxnk1aVVb_Nafg18tstoLS6SvUOpNcoSQ2qFazWio"
// 	data, err := c.ConcurrentDownloadChunkData(arId, 0)
// 	// data , err := c.DownloadChunkData(arId)
// 	assert.NoError(t, err)
// 	ioutil.WriteFile("nannan.gif", data, 0666)
// 	chunks := utils.GenerateChunks(data)
// 	dataRoot := utils.Base64Encode(chunks.DataRoot)
// 	t.Log(dataRoot)
// 	t.Log(len(data))
// }

func TestClient_ExistTxData(t *testing.T) {
	c := NewClient("https://arweave.net")
	arId := "trMxnk1aVVb_Nafg18tstoLS6SvUOpNcoSQ2qFazWio"
	exist, err := c.ExistTxData(arId)
	assert.NoError(t, err)
	t.Log(exist)
}
