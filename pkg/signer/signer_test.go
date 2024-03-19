package goar

import (
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSigner(t *testing.T) {
	signer, err := NewSignerFromPath("../../example/testKey.json")
	assert.NoError(t, err)
	tags := []types.Tag{
		{Name: "GOAR", Value: "sendAR"},
	}
	tx := &types.Transaction{
		Format:    2,
		LastTx:    "",
		Owner:     signer.Owner(),
		Tags:      utils.TagsEncode(tags),
		Target:    "cSYOy8-p1QFenktkDBFyRM3cwZSTrQ_J4EsELLho_UE",
		Quantity:  "1000000000",
		Data:      "",
		DataSize:  "0",
		DataRoot:  "",
		Reward:    "",
		Signature: "",
	}
	err = signer.SignTx(tx)
	assert.NoError(t, err)

	err = utils.VerifyTransaction(*tx)
	assert.NoError(t, err)
}
