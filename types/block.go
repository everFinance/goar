package types

type Block struct {
	Nonce                    string        `json:"nonce"`
	PreviousBlock            string        `json:"previous_block"`
	Timestamp                int64         `json:"timestamp"`
	LastRetarget             int64         `json:"last_retarget"`
	Diff                     interface{}   `json:"diff"`
	Height                   int64         `json:"height"`
	Hash                     string        `json:"hash"`
	IndepHash                string        `json:"indep_hash"`
	Txs                      []string      `json:"txs"`
	TxRoot                   string        `json:"tx_root"`
	TxTree                   interface{}   `json:"tx_tree"`
	HashList                 interface{}   `json:"hash_list"`
	HashListMerkle           string        `json:"hash_list_merkle"`
	WalletList               string        `json:"wallet_list"`
	RewardAddr               string        `json:"reward_addr"`
	Tags                     []interface{} `json:"tags"`
	RewardPool               interface{}   `json:"reward_pool"` // always string
	WeaveSize                interface{}   `json:"weave_size"`  // always string
	BlockSize                interface{}   `json:"block_size"`  // always string
	CumulativeDiff           interface{}   `json:"cumulative_diff"`
	SizeTaggedTxs            interface{}   `json:"size_tagged_txs"`
	Poa                      POA           `json:"poa"`
	UsdToArRate              []string      `json:"usd_to_ar_rate"`
	ScheduledUsdToArRate     []string      `json:"scheduled_usd_to_ar_rate"`
	Packing25Threshold       string        `json:"packing_2_5_threshold"`
	StrictDataSplitThreshold string        `json:"strict_data_split_threshold"`
}

type POA struct {
	Option   string `json:"option"`
	TxPath   string `json:"tx_path"`
	DataPath string `json:"data_path"`
	Chunk    string `json:"chunk"`
}
