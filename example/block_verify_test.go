package example

import (
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIndepHash(t *testing.T) {
	height := int64(859655)
	cli := goar.NewClient("https://arweave.net")
	b, err := cli.GetBlockByHeight(height)
	assert.NoError(t, err)
	indepHash := utils.IndepHash(*b)
	assert.Equal(t, "vbiDXtb8uDyQolW_gNiSyTwkUB2AVtTfjC0xgbXjP4OmwVeTYeK8C-S2svVZsH1-", utils.Base64Encode(indepHash))
}
