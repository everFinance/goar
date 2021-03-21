package transfer

import (
	"github.com/everFinance/goar/client"
	"github.com/everFinance/goar/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Tx(t *testing.T) {
	c := client.New("https://arweave.net")
	id := "HepOYsiov0ijeJIIE6o2sajI_vekdD1wSaFvP0-Ujw8"

	// body, statusCode, err := c.HttpGet(fmt.Sprintf("tx/%s", id))
	// assert.NoError(t, err)
	// t.Log("body: ", string(body))
	// t.Log("code: ", statusCode)
	// tx := &TransactionChunks{}
	// err = json.Unmarshal(body, tx)
	// assert.NoError(t, err)
	// t.Log(string(tx.Data))

	// url := fmt.Sprintf("tx/%v/%v", id, "data")
	// body, statusCode, err := c.HttpGet(url)
	// assert.NoError(t, err)
	// t.Log("code: ", statusCode)
	// t.Log(string(body))

	tx, status, err := c.GetTransactionByID(id)
	assert.NoError(t, err)
	t.Log(status)
	t.Log(tx.Data)
	t.Log(tx)
}

func TestVerifyTransaction(t *testing.T) {
	tx := TransactionChunks{
		Format: 2,
		ID:     "QLRuAX12XtlNwDZ1UdPTfGV5grt-V1BJksTjyvbGeM4",
		LastTx: "JrWwVEXja1AI7mVB2v24Ye_Dd2WIUTVaHdmpFrse37dOqnqPy-6UtsPuXeHl4SCG",
		Owner:  "oA4bDVEjBSVf324cQPfyJcE-90rwPo1xOvrvc7g2-Lminag5lT0JZAVfcSjg1vgeoivSu2I0yO-jznomZ1m4H2CUJGY8Hc3wEsO3SUqxPIEaOKuFvezeAZpRuh__SlzYGfvhwLfoGf7KJ6UlvcuNm49xyIfsRGTc52u-fTDevqvBtz5YtYsyk6LcMoCMDoPzE4ldTBZ8V3ucaFXz-kbwkN1aU2Ph8MYfKIySOtjeVsCxU5BBRn39cHm4dVbDqEZu8QT-l26d8QjITl91tpTatzGivvFioR5_dJEtlo1xkNsOYIsPXT6WrwgPkRYflqhJ8GW9wqp7QCoekYIxr2nQXxOm0lLsPoA_gx4ZCDGsfUtumWdk2UV-cKsvQo-iuxPPVWX-7bWTP9aMKDTyyywW9Ho6r99ggeSEpNYYGqYBmLb-v2S0L-8EdvjAQ2M6yCw0qpNN6XhFvNVcPZXATrhhDIijp-KXDHx2zy2IDubjShCmZOSRrJy5OV6EfJOpPEgBXO-CdKtVk48uh8htb5SMEZ48hDyeU-4Htjaryz-N_M1n9Rv9ffqVQf-pIrP_cpw0DzxcrvFsFRYnAkFLnXlY9FX5mha2FW5veeL0ZMCN7ETEVovCrO1Q8_sV9v2rs6S8NuMjA4RW93nKJjqvAJFkt7rMpv4a_tuftZvuX8OzxOc",
		Tags: []types.Tag{
			types.Tag{
				Name:  "QXBwLU5hbWU",
				Value: "U21hcnRXZWF2ZUFjdGlvbg",
			},
			types.Tag{
				Name:  "QXBwLVZlcnNpb24",
				Value: "MC4zLjA",
			},
			types.Tag{
				Name:  "Q29udHJhY3Q",
				Value: "dHJ1ZQ",
			},
			types.Tag{
				Name:  "SW5wdXQ",
				Value: "eyJmdW5jdGlvbiI6InRyYW5zZmVyIiwicXR5Ijo1MDAsInRhcmdldCI6Ilp5aGhBTHdxazhuMnVyV1Y0RTNqSEJjNzd3YWE1RnItcUhscl9jdGlIQk0ifQ",
			},
		},
		Target:    "",
		Quantity:  "0",
		DataSize:  "4",
		DataRoot:  "z3rQGxyiqdQuOh2dxDst176oOKmW3S9MwQNTEh4DK1U",
		Reward:    "1195477",
		Signature: "CjGBVPulxEzsKBi83L7dhAQtgWf7vjT5JDYazJyi4-p37nVA0ghQcGbeGjKy9HO4t-dgLqDExw_PDbtr_9SRbRmkBNEPgnlFVT4U82MqxgJrVm6adMlJtvC-Vw-O3nfT3ObtaBOCxUFflOcTrPAW7V4p0MmXwU3u_xw4hPVYW9Da1c2SnwFzWDU5mG0y8pego9ZNWM9bfYylQz25fOgfDeWJgHZ5g540EfH2wC55obx_qCezBVCFd-hiiznP5UXMplR6exQM_fBomfMFd7TAfhYkBV-eRqykmj68xGQOS4ynwKFWajrM4BiP-6fc68bQn8PLYjtcvBAdhH9J8zPZaArY7ozRwnmveLe-lfQG7pQUDKpwXIOUcr6N3wBotN1Tm37k6Lp-hGi24zQhndZmf6S6mrcodanvXKBUYgMqs6TrEHSNFzX69WmoxTdW13COv0txY_wePB_RYRlnCuwOEiNj396_pZoTHdxe2Qvl86ZP_rlCVvpmVPmBxon0i6kdvxC02w5rWnKBh2YK-wyed47SyCgX6EVEThTxKcQeXvs6yIwxYOhH044_oSOzOouOAZqigtGy6BDYJu2Y4jQq9N55SRCP1VB6F1AQsHPAMmnPnewwXEJl-vG3MvESZzjEhtG5KtN4uwzGMIPLtY3dba1EUWfcxgwfUekUIREcKv4",
	}

	assert.NoError(t, VerifyTransaction(tx))
}

func TestGetSignatureData(t *testing.T) {
	tx := TransactionChunks{
		Format: 2,
		ID:     "QLRuAX12XtlNwDZ1UdPTfGV5grt-V1BJksTjyvbGeM4",
		LastTx: "JrWwVEXja1AI7mVB2v24Ye_Dd2WIUTVaHdmpFrse37dOqnqPy-6UtsPuXeHl4SCG",
		Owner:  "oA4bDVEjBSVf324cQPfyJcE-90rwPo1xOvrvc7g2-Lminag5lT0JZAVfcSjg1vgeoivSu2I0yO-jznomZ1m4H2CUJGY8Hc3wEsO3SUqxPIEaOKuFvezeAZpRuh__SlzYGfvhwLfoGf7KJ6UlvcuNm49xyIfsRGTc52u-fTDevqvBtz5YtYsyk6LcMoCMDoPzE4ldTBZ8V3ucaFXz-kbwkN1aU2Ph8MYfKIySOtjeVsCxU5BBRn39cHm4dVbDqEZu8QT-l26d8QjITl91tpTatzGivvFioR5_dJEtlo1xkNsOYIsPXT6WrwgPkRYflqhJ8GW9wqp7QCoekYIxr2nQXxOm0lLsPoA_gx4ZCDGsfUtumWdk2UV-cKsvQo-iuxPPVWX-7bWTP9aMKDTyyywW9Ho6r99ggeSEpNYYGqYBmLb-v2S0L-8EdvjAQ2M6yCw0qpNN6XhFvNVcPZXATrhhDIijp-KXDHx2zy2IDubjShCmZOSRrJy5OV6EfJOpPEgBXO-CdKtVk48uh8htb5SMEZ48hDyeU-4Htjaryz-N_M1n9Rv9ffqVQf-pIrP_cpw0DzxcrvFsFRYnAkFLnXlY9FX5mha2FW5veeL0ZMCN7ETEVovCrO1Q8_sV9v2rs6S8NuMjA4RW93nKJjqvAJFkt7rMpv4a_tuftZvuX8OzxOc",
		Tags: []types.Tag{
			types.Tag{
				Name:  "QXBwLU5hbWU",
				Value: "U21hcnRXZWF2ZUFjdGlvbg",
			},
			types.Tag{
				Name:  "QXBwLVZlcnNpb24",
				Value: "MC4zLjA",
			},
			types.Tag{
				Name:  "Q29udHJhY3Q",
				Value: "dHJ1ZQ",
			},
			types.Tag{
				Name:  "SW5wdXQ",
				Value: "eyJmdW5jdGlvbiI6InRyYW5zZmVyIiwicXR5Ijo1MDAsInRhcmdldCI6Ilp5aGhBTHdxazhuMnVyV1Y0RTNqSEJjNzd3YWE1RnItcUhscl9jdGlIQk0ifQ",
			},
		},
		Target:    "",
		Quantity:  "0",
		DataSize:  "4",
		DataRoot:  "z3rQGxyiqdQuOh2dxDst176oOKmW3S9MwQNTEh4DK1U",
		Reward:    "1195477",
		Signature: "CjGBVPulxEzsKBi83L7dhAQtgWf7vjT5JDYazJyi4-p37nVA0ghQcGbeGjKy9HO4t-dgLqDExw_PDbtr_9SRbRmkBNEPgnlFVT4U82MqxgJrVm6adMlJtvC-Vw-O3nfT3ObtaBOCxUFflOcTrPAW7V4p0MmXwU3u_xw4hPVYW9Da1c2SnwFzWDU5mG0y8pego9ZNWM9bfYylQz25fOgfDeWJgHZ5g540EfH2wC55obx_qCezBVCFd-hiiznP5UXMplR6exQM_fBomfMFd7TAfhYkBV-eRqykmj68xGQOS4ynwKFWajrM4BiP-6fc68bQn8PLYjtcvBAdhH9J8zPZaArY7ozRwnmveLe-lfQG7pQUDKpwXIOUcr6N3wBotN1Tm37k6Lp-hGi24zQhndZmf6S6mrcodanvXKBUYgMqs6TrEHSNFzX69WmoxTdW13COv0txY_wePB_RYRlnCuwOEiNj396_pZoTHdxe2Qvl86ZP_rlCVvpmVPmBxon0i6kdvxC02w5rWnKBh2YK-wyed47SyCgX6EVEThTxKcQeXvs6yIwxYOhH044_oSOzOouOAZqigtGy6BDYJu2Y4jQq9N55SRCP1VB6F1AQsHPAMmnPnewwXEJl-vG3MvESZzjEhtG5KtN4uwzGMIPLtY3dba1EUWfcxgwfUekUIREcKv4",
	}

	signData, err := GetSignatureData(&tx)
	assert.NoError(t, err)
	assert.Equal(t, []uint8([]byte{0xe2, 0x9f, 0xa3, 0x1b, 0xf2, 0x35, 0xe8, 0xc0, 0x9d, 0x97, 0x12, 0xe3, 0x40, 0x5, 0x4a, 0x26, 0xba, 0x8f, 0x67, 0x47, 0x3b, 0x68, 0x72, 0x25, 0xc0, 0x6d, 0x92, 0xd8, 0xb0, 0x59, 0x8a, 0x88, 0x89, 0x10, 0x41, 0xf9, 0x60, 0xee, 0x25, 0x23, 0x82, 0xbd, 0x14, 0xcc, 0xc, 0xeb, 0xa9, 0x23}), signData)
}
