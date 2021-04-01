package types

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
