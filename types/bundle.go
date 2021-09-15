package types

type Bundle struct {
	Items        []BundleItem `json:"items"`
	BundleBinary []byte
}

type BundleItem struct {
	SignatureType string `json:"signatureType"`
	Signature     string `json:"signature"`
	Owner         string `json:"owner"`  //  utils.Base64Encode(wallet.PubKey.N.Bytes())
	Target        string `json:"target"` // optional
	Anchor        string `json:"anchor"` // optional
	Tags          []Tag  `json:"tags"`
	Data          string `json:"data"`
	Id            string `json:"id"`

	ItemBinary []byte
}
