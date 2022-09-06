package types

type KeyFile struct {
	D   string `json:"d"` // privKey
	Dp  string `json:"dp"`
	Dq  string `json:"dq"`
	E   string `json:"e"` // pubKey exp
	Ext bool   `json:"ext"`
	Kty string `json:"kty"`
	N   string `json:"n"` // pubKey
	P   string `json:"p"`
	Q   string `json:"q"`
	Qi  string `json:"qi"`
}
