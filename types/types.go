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

type BundlrResp struct {
	Id                  string   `json:"id"`
	Signature           string   `json:"signature"`
	N                   string   `json:"n"`
	Public              string   `json:"public"`
	Block               int64    `json:"block"`
	ValidatorSignatures []string `json:"validatorSignatures"` 
}
