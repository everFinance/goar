package goar

import (
	"fmt"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/tidwall/gjson"
	"math/rand"
	"strconv"
)

func (w *Wallet) WarpTransfer(contractId, target string, qty int64) (id string, err error) {
	tags, err := utils.PstTransferTags(contractId, target, qty, true)
	if err != nil {
		return
	}
	data := strconv.Itoa(rand.Intn(9999))

	tx := &types.Transaction{
		Format:    2,
		ID:        "",
		LastTx:    "p7vc1iSP6bvH_fCeUFa9LqoV5qiyW-jdEKouAT0XMoSwrNraB9mgpi29Q10waEpO", // default
		Owner:     w.Signer.Owner(),
		Tags:      utils.TagsEncode(tags),
		Target:    "",
		Quantity:  "0",
		Data:      utils.Base64Encode([]byte(data)),
		DataSize:  fmt.Sprintf("%d", len(data)),
		DataRoot:  "",
		Reward:    "0",
		Signature: "",
	}
	if err = w.Signer.SignTx(tx); err != nil {
		return
	}
	// send to wrap gateway
	result, err := w.Client.SubmitToWarp(tx) // {"id":"BQQyqbsULPNpyKgwVSSn8z0-3Km_y1GPMiLU1-eR_lc"}
	if err != nil {
		return
	}
	return gjson.ParseBytes(result).Get("id").String(), nil
}
