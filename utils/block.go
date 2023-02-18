package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/daqiancode/goar/types"
)

const (
	height_2_0 = int64(422250)
	height_2_4 = int64(633720)
	height_2_5 = int64(812970)
)

func GenerateIndepHash(b types.Block) string {
	if b.Height < height_2_0 { // not support arweave v1.0
		return b.IndepHash
	}

	bds := generateBlockDataSegment(b)
	list := make([]interface{}, 0)
	list = append(list, Base64Encode(bds))
	list = append(list, b.Hash)
	list = append(list, b.Nonce)
	if b.Height >= height_2_4 {
		list = append(list, poaToList(b.Poa))
	}
	hash := DeepHash(list)
	return Base64Encode(hash[:])
}

func generateBlockDataSegment(b types.Block) []byte {
	bdsBase := generateBlockDataSegmentBase(b)

	list := make([]interface{}, 0)
	list = append(list, Base64Encode(bdsBase))
	list = append(list, Base64Encode([]byte(fmt.Sprintf("%d", b.Timestamp))))
	list = append(list, Base64Encode([]byte(fmt.Sprintf("%d", b.LastRetarget))))
	list = append(list, Base64Encode([]byte(fmt.Sprintf("%v", b.Diff))))
	list = append(list, Base64Encode([]byte(fmt.Sprintf("%v", b.CumulativeDiff))))
	list = append(list, Base64Encode([]byte(fmt.Sprintf("%v", b.RewardPool))))
	list = append(list, b.WalletList)
	list = append(list, b.HashListMerkle)
	hash := DeepHash(list)
	return hash[:]
}

func generateBlockDataSegmentBase(b types.Block) []byte {
	props := make([]interface{}, 0)
	props = append(props, Base64Encode([]byte(fmt.Sprintf("%d", b.Height))))
	props = append(props, b.PreviousBlock)
	props = append(props, b.TxRoot)
	props = append(props, b.Txs)

	props = append(props, Base64Encode([]byte(fmt.Sprintf("%v", b.BlockSize))))
	props = append(props, Base64Encode([]byte(fmt.Sprintf("%v", b.WeaveSize))))

	if b.RewardAddr == "unclaimed" {
		props = append(props, Base64Encode([]byte("unclaimed")))
	} else {
		props = append(props, b.RewardAddr)
	}
	props = append(props, b.Tags)

	endProps := make([]interface{}, 0)
	if b.Height >= height_2_4 {
		props2 := make([]interface{}, 0)
		if b.Height >= height_2_5 {
			RateDividend := b.UsdToArRate[0]
			RateDivisor := b.UsdToArRate[1]

			ScheduledRateDividend := b.ScheduledUsdToArRate[0]
			ScheduledRateDivisor := b.ScheduledUsdToArRate[1]

			props2 = append(props2, Base64Encode([]byte(fmt.Sprintf("%s", RateDividend))))
			props2 = append(props2, Base64Encode([]byte(fmt.Sprintf("%s", RateDivisor))))
			props2 = append(props2, Base64Encode([]byte(fmt.Sprintf("%s", ScheduledRateDividend))))
			props2 = append(props2, Base64Encode([]byte(fmt.Sprintf("%s", ScheduledRateDivisor))))

			props2 = append(props2, Base64Encode([]byte(fmt.Sprintf("%s", b.Packing25Threshold))))
			props2 = append(props2, Base64Encode([]byte(fmt.Sprintf("%s", b.StrictDataSplitThreshold))))
		}
		endProps = append(props2, props...)
	} else {
		endProps = append(props, poaToList(b.Poa))
	}

	hash := DeepHash(endProps)
	return hash[:]
}

func poaToList(poa types.POA) []string {
	return []string{
		Base64Encode([]byte(poa.Option)),
		poa.TxPath,
		poa.DataPath,
		poa.Chunk,
	}
}

func DecodeBlock(body string) (*types.Block, error) {
	b := &types.Block{}
	// json unmarshal exist number precision problem
	decoder := json.NewDecoder(strings.NewReader(body))
	decoder.UseNumber()
	err := decoder.Decode(b)
	if err != nil {
		return nil, err
	}
	formatBlockFields(b)
	return b, err
}

func formatBlockFields(b *types.Block) {
	if _, ok := b.RewardPool.(string); !ok {
		by, _ := json.Marshal(b.RewardPool)
		b.RewardPool = string(by)
	}

	if _, ok := b.WeaveSize.(string); !ok {
		by, _ := json.Marshal(b.WeaveSize)
		b.WeaveSize = string(by)
	}

	if _, ok := b.BlockSize.(string); !ok {
		by, _ := json.Marshal(b.BlockSize)
		b.BlockSize = string(by)
	}
}
