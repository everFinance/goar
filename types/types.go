package types

import (
	"github.com/everFinance/goar/utils"
)

const (
	SuccessTxStatus = "Success"
	PendingTxStatus = "Pending"
)

type NetworkInfo struct {
	Network          string `json:"network"`
	Version          int64  `json:"version"`
	Release          int64  `json:"release"`
	Height           int64  `json:"height"`
	Current          string `json:"current"`
	Blocks           int64  `json:"blocks"`
	Peers            int64  `json:"peers"`
	QueueLength      int64  `json:"queue_length"`
	NodeStateLatency int64  `json:"node_state_latency"`
}

type Block struct {
	Nonce         string        `json:"nonce"`
	PreviousBlock string        `json:"previous_block"`
	Timestamp     int64         `json:"timestamp"`
	LastRetarget  int64         `json:"last_retarget"`
	Diff          string        `json:"diff"`
	Height        int64         `json:"height"`
	Hash          string        `json:"hash"`
	IndepHash     string        `json:"indep_hash"`
	Txs           []interface{} `json:"txs"`
	WalletList    string        `json:"wallet_list"`
	RewardAddr    string        `json:"reward_addr"`
	Tags          []interface{} `json:"tags"`
	RewardPool    int           `json:"reward_pool"`
	WeaveSize     int           `json:"weave_size"`
	BlockSize     int           `json:"block_size"`
}

type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func TagsEncode(tags []Tag) []Tag {
	base64Tags := []Tag{}

	for _, tag := range tags {
		base64Tags = append(base64Tags, Tag{
			Name:  utils.Base64Encode([]byte(tag.Name)),
			Value: utils.Base64Encode([]byte(tag.Value)),
		})
	}

	return base64Tags
}

func TagsDecode(base64Tags []Tag) ([]Tag, error) {
	tags := []Tag{}

	for _, bt := range base64Tags {
		bName, err := utils.Base64Decode(bt.Name)
		if err != nil {
			return nil, err
		}

		bValue, err := utils.Base64Decode(bt.Value)
		if err != nil {
			return nil, err
		}

		tags = append(tags, Tag{
			Name:  string(bName),
			Value: string(bValue),
		})
	}

	return tags, nil
}
