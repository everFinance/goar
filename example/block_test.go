package example

import (
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIndepHash(t *testing.T) {
	// height_2_0
	height := int64(422250)
	cli := goar.NewClient("https://arweave.net")
	b, err := cli.GetBlockByHeight(height)
	assert.NoError(t, err)
	indepHash := utils.GenerateIndepHash(*b)
	assert.Equal(t, "5VTARz7bwDO4GqviCSI9JXm8_JOtoQwF-QCZm0Gt2gVgwdzSY3brOtOD46bjMz09", indepHash)

	height = int64(422244)
	b, err = cli.GetBlockByHeight(height)
	assert.NoError(t, err)
	indepHash = utils.GenerateIndepHash(*b)
	assert.Equal(t, "ueAvRMVZRvH9aLeIL-LRkJ0TUWLydQcFKh6RiKZ2aLaSonoUZ16y1tQ74nKBCPYp", indepHash)

	// height_2_4
	height = int64(633720)
	b, err = cli.GetBlockByHeight(height)
	assert.NoError(t, err)
	indepHash = utils.GenerateIndepHash(*b)
	assert.Equal(t, "7GofRKP53XhLBgTUPoBRTWUq8ncHnOSpsMRjxw5fs5nu8x1png1gCWm7STt68nhq", indepHash)

	height = int64(633719)
	b, err = cli.GetBlockByHeight(height)
	assert.NoError(t, err)
	indepHash = utils.GenerateIndepHash(*b)
	assert.Equal(t, "OzC3T3l0bTcYpEESxliBKfPQ57JIiA7KvOTeKZVXCDjv9UMAlv2-ua04VtRXQ2AR", indepHash)

	// height_2_5
	height = int64(812970)
	b, err = cli.GetBlockByHeight(height)
	assert.NoError(t, err)
	indepHash = utils.GenerateIndepHash(*b)
	assert.Equal(t, "nIq5881hbLMH5vPsv0mwrP6Je-4-0fp0AOSf2UbsQ1jnoA3SfSOYZm4dd6X3g2lu", indepHash)

	height = int64(812969)
	b, err = cli.GetBlockByHeight(height)
	assert.NoError(t, err)
	indepHash = utils.GenerateIndepHash(*b)
	assert.Equal(t, "R58RTTzEKmSyaqhRik4fvl9AkN3g98QEntvZuoly02uwm8J4fZbcvgv9wEEgN5Ne", indepHash)
}
