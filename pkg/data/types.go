package data

const (
	Arweave  = 1
	ED25519  = 2
	Ethereum = 3
	Solana   = 4
)

type SignatureMeta struct {
	SignatureLength int
	PublicKeyLength int
	Name            string
}

var SignatureConfig = map[int]SignatureMeta{
	Arweave: {
		SignatureLength: 512,
		PublicKeyLength: 512,
		Name:            "arweave",
	},
	ED25519: {
		SignatureLength: 64,
		PublicKeyLength: 32,
		Name:            "ed25519",
	},
	Ethereum: {
		SignatureLength: 65,
		PublicKeyLength: 65,
		Name:            "ethereum",
	},
	Solana: {
		SignatureLength: 64,
		PublicKeyLength: 32,
		Name:            "solana",
	},
}

type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DataItem struct {
	ID            string `json:"id"`
	Signature     string `json:"signature"`
	SignatureType int    `json:"signature_type"`
	Owner         string `json:"owner"`  //  utils.Base64Encode(pubkey)
	Target        string `json:"target"` // optional, if exist must length 32, and is base64 str
	Anchor        string `json:"anchor"` // optional, if exist must length 32, and is base64 str
	Tags          []Tag  `json:"tags"`
	Data          string `json:"data"`
	Raw           []byte
}

type BundleHeader struct {
	id   int
	size int
	raw  []byte
}

type Bundle struct {
	Headers []BundleHeader `json:"bundle_header"`
	Items   []DataItem     `json:"items"`
	RawData string         `json:"raw_data"`
}
