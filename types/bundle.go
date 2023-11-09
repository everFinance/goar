package types

import (
	"os"
)

const (
	ArweaveSignType  = 1
	ED25519SignType  = 2
	EthereumSignType = 3
	SolanaSignType   = 4
	EC256SignType    = 5
	RS256SignType    = 6
)

type SigMeta struct {
	SigLength int
	PubLength int
	SigName   string
}

var SigConfigMap = map[int]SigMeta{
	ArweaveSignType: {
		SigLength: 512,
		PubLength: 512,
		SigName:   "arweave",
	},
	ED25519SignType: {
		SigLength: 64,
		PubLength: 32,
		SigName:   "ed25519",
	},
	EthereumSignType: {
		SigLength: 65,
		PubLength: 65,
		SigName:   "ethereum",
	},
	SolanaSignType: {
		SigLength: 64,
		PubLength: 32,
		SigName:   "solana",
	},
	EC256SignType: {
		SigLength: 70,
		PubLength: 32,
		SigName:   "ec256",
	},
	RS256SignType: {
		SigLength: 344,
		PubLength: 256,
		SigName:   "rs256",
	},
}

type Bundle struct {
	Items            []BundleItem `json:"items"`
	BundleBinary     []byte
	BundleDataReader *os.File
}

type BundleItem struct {
	SignatureType int    `json:"signatureType"`
	Signature     string `json:"signature"`
	Owner         string `json:"owner"`  //  utils.Base64Encode(pubkey)
	Target        string `json:"target"` // optional, if exist must length 32, and is base64 str
	Anchor        string `json:"anchor"` // optional, if exist must length 32, and is base64 str
	Tags          []Tag  `json:"tags"`
	Data          string `json:"data"`
	Id            string `json:"id"`
	TagsBy        string `json:"tagsBy"` // utils.Base64Encode(TagsBytes) for retry assemble item

	ItemBinary []byte   `json:"-"`
	DataReader *os.File `json:"-"`
}
