package goar

import (
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"

	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/everFinance/gojwk"
)

type Wallet struct {
	Client  *Client
	PubKey  *rsa.PublicKey
	PrvKey  *rsa.PrivateKey
	Address string
}

// proxyUrl: option
func NewWalletFromPath(path string, clientUrl string, proxyUrl ...string) (*Wallet, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewWallet(b, clientUrl, proxyUrl...)
}

func NewWallet(b []byte, clientUrl string, proxyUrl ...string) (w *Wallet, err error) {
	key, err := gojwk.Unmarshal(b)
	if err != nil {
		return
	}

	pubKey, err := key.DecodePublicKey()
	if err != nil {
		return
	}
	pub, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		err = fmt.Errorf("pubKey type error")
		return
	}
	prvKey, err := key.DecodePrivateKey()
	if err != nil {
		return
	}
	prv, ok := prvKey.(*rsa.PrivateKey)
	if !ok {
		err = fmt.Errorf("prvKey type error")
		return
	}

	addr := sha256.Sum256(pub.N.Bytes())
	w = &Wallet{
		Client:  NewClient(clientUrl, proxyUrl...),
		PubKey:  pub,
		PrvKey:  prv,
		Address: utils.Base64Encode(addr[:]),
	}

	return
}

func (w *Wallet) SendAR(amount *big.Float, target string, tags []types.Tag) (id string, err error) {
	return w.SendWinstonSpeedUp(utils.ARToWinston(amount), target, tags, 0)
}

func (w *Wallet) SendARSpeedUp(amount *big.Float, target string, tags []types.Tag, speedFactor int64) (id string, err error) {
	return w.SendWinstonSpeedUp(utils.ARToWinston(amount), target, tags, speedFactor)
}

func (w *Wallet) SendWinston(amount *big.Int, target string, tags []types.Tag) (id string, err error) {
	return w.SendWinstonSpeedUp(amount, target, tags, 0)
}

func (w *Wallet) SendWinstonSpeedUp(amount *big.Int, target string, tags []types.Tag, speedFactor int64) (id string, err error) {
	reward, err := w.Client.GetTransactionPrice(nil, &target)
	if err != nil {
		return
	}

	tx := &types.Transaction{
		Format:   2,
		Target:   target,
		Quantity: amount.String(),
		Tags:     utils.TagsEncode(tags),
		Data:     "",
		DataSize: "0",
		Reward:   fmt.Sprintf("%d", reward*(100+speedFactor)/100),
	}

	return w.SendTransaction(tx)
}

func (w *Wallet) SendData(data []byte, tags []types.Tag) (id string, err error) {
	return w.SendDataSpeedUp(data, tags, 0)
}

// SendDataSpeedUp set speedFactor for speed up
// eg: speedFactor = 10, reward = 1.1 * reward
func (w *Wallet) SendDataSpeedUp(data []byte, tags []types.Tag, speedFactor int64) (id string, err error) {
	reward, err := w.Client.GetTransactionPrice(data, nil)
	if err != nil {
		return
	}

	tx := &types.Transaction{
		Format:   2,
		Target:   "",
		Quantity: "0",
		Tags:     utils.TagsEncode(tags),
		Data:     utils.Base64Encode(data),
		DataSize: fmt.Sprintf("%d", len(data)),
		Reward:   fmt.Sprintf("%d", reward*(100+speedFactor)/100),
	}

	return w.SendTransaction(tx)
}

// SendTransaction: if send success, should return pending
func (w *Wallet) SendTransaction(tx *types.Transaction) (id string, err error) {
	anchor, err := w.Client.GetTransactionAnchor()
	if err != nil {
		return
	}
	tx.LastTx = anchor

	if err = utils.SignTransaction(tx, w.PubKey, w.PrvKey); err != nil {
		return
	}

	id = tx.ID

	uploader, err := CreateUploader(w.Client, tx, nil)
	if err != nil {
		return
	}
	for !uploader.IsComplete() {
		err = uploader.UploadChunk()
		if err != nil {
			return
		}
	}
	return
}

// about bundle tx

func (w *Wallet) CreateBundleDataItem(data []byte, signatureType int, target string, anchor string, tags []types.Tag) (di DataItem, err error) {
	targetBytes := []byte{}
	if target != "" {
		targetBytes, err = utils.Base64Decode(target)
		if err != nil {
			return
		}
	}

	anchorBytes := []byte{}
	if anchor != "" {
		anchorBytes, err = utils.Base64Decode(anchor)
		if err != nil {
			return
		}
	}

	tagsBytes, err := utils.SerializeTags(tags)
	if err != nil {
		return di, err
	}

	ownerByte := w.PubKey.N.Bytes()
	dataItem, err := newDataItem(utils.Base64Encode(ownerByte), strconv.Itoa(signatureType), target, anchor, data, tags)
	if err != nil {
		return di, err
	}
	// sign
	err = dataItem.Sign(w.PrvKey)
	if err != nil {
		return di, err
	}
	// Create array with set length
	bytesArr := make([]byte, 0, 1044)
	bytesArr = append(bytesArr, utils.ShortTo2ByteArray(signatureType)...)
	// Push bytes for `signature`
	sig, err := utils.Base64Decode(dataItem.Signature)
	if err != nil {
		return di, err
	}
	bytesArr = append(bytesArr, sig...)
	// Push bytes for `ownerByte`
	bytesArr = append(bytesArr, ownerByte...)
	// Push `presence byte` and push `target` if present
	// 64 + OWNER_LENGTH
	if target != "" {
		bytesArr = append(bytesArr, byte(1))
		bytesArr = append(bytesArr, targetBytes...)
	} else {
		bytesArr = append(bytesArr, byte(0))
	}

	// Push `presence byte` and push `anchor` if present
	// 64 + OWNER_LENGTH
	if anchor != "" {
		bytesArr = append(bytesArr, byte(1))
		bytesArr = append(bytesArr, anchorBytes...)
	} else {
		bytesArr = append(bytesArr, byte(0))
	}

	// push tags
	bytesArr = append(bytesArr, utils.LongTo8ByteArray(len(tags))...)
	bytesArr = append(bytesArr, utils.LongTo8ByteArray(len(tagsBytes))...)

	if tags != nil {
		bytesArr = append(bytesArr, tagsBytes...)
	}

	// push data
	bytesArr = append(bytesArr, data...)
	dataItem.itemBinary = bytesArr
	return *dataItem, nil
}

func (w *Wallet) SubmitBundleTx(bundleBinary []byte, tags []types.Tag, txSpeed int64) (txId string, err error){
	bundleTags := []types.Tag{
		{Name: "Bundle-Format", Value: "binary"},
		{Name: "Bundle-Version", Value: "2.0.0"},
	}
	txTags := make([]types.Tag, 0)
	txTags = append(bundleTags, tags...)
	txId, err = w.SendDataSpeedUp(bundleBinary, txTags, txSpeed)
	return
}