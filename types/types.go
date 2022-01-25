package types

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
	RewardPool               interface{}   `json:"reward_pool"`
	WeaveSize                interface{}   `json:"weave_size"`
	BlockSize                interface{}   `json:"block_size"`
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

type TransactionChunk struct {
	Chunk    string `json:"chunk"`
	DataPath string `json:"data_path"`
	TxPath   string `json:"tx_path"`
}

type TransactionOffset struct {
	Size   string `json:"size"`
	Offset string `json:"offset"`
}

type TxStatus struct {
	BlockHeight           int    `json:"block_height"`
	BlockIndepHash        string `json:"block_indep_hash"`
	NumberOfConfirmations int    `json:"number_of_confirmations"`
}

type BundlerResp struct {
	Id        string `json:"id"`
	Signature string `json:"signature"`
	N         string `json:"n"`
}
