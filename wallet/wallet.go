package wallet

import (
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"github.com/everFinance/goar/uploader"
	"io/ioutil"
	"math/big"

	"github.com/everFinance/goar/client"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/everFinance/gojwk"
)

type Wallet struct {
	Client  *client.Client
	PubKey  *rsa.PublicKey
	PrvKey  *rsa.PrivateKey
	Address string
}

// proxyUrl: option
func NewFromPath(path string, clientUrl string, proxyUrl ...string) (*Wallet, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return New(b, clientUrl, proxyUrl...)
}

func New(b []byte, clientUrl string, proxyUrl ...string) (w *Wallet, err error) {
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
		Client:  client.New(clientUrl, proxyUrl...),
		PubKey:  pub,
		PrvKey:  prv,
		Address: utils.Base64Encode(addr[:]),
	}

	return
}

func (w *Wallet) SendAR(amount *big.Float, target string, tags []types.Tag) (id, status string, err error) {
	return w.SendWinston(utils.ARToWinston(amount), target, tags)
}

func (w *Wallet) SendWinston(amount *big.Int, target string, tags []types.Tag) (id, status string, err error) {
	reward, err := w.Client.GetTransactionPrice(nil, &target)
	if err != nil {
		return
	}

	tx := &types.Transaction{
		Format:   2,
		Target:   target,
		Quantity: amount.String(),
		Tags:     types.TagsEncode(tags),
		Data:     []byte{},
		DataSize: "0",
		Reward:   fmt.Sprintf("%d", reward),
	}

	return w.SendTransaction(tx)
}

// SendDataSpeedUp set speedFactor for speed up
// eg: speedFactor = 10, reward = 1.1 * reward
func (w *Wallet) SendDataSpeedUp(data []byte, tags []types.Tag, speedFactor int64) (id, status string, err error) {
	reward, err := w.Client.GetTransactionPrice(data, nil)
	if err != nil {
		return
	}

	tx := &types.Transaction{
		Format:   2,
		Target:   "",
		Quantity: "0",
		Tags:     types.TagsEncode(tags),
		Data:     data,
		DataSize: fmt.Sprintf("%d", len(data)),
		Reward:   fmt.Sprintf("%d", reward*(100+speedFactor)/100),
	}

	return w.SendTransaction(tx)
}

func (w *Wallet) SendData(data []byte, tags []types.Tag) (id, status string, err error) {
	reward, err := w.Client.GetTransactionPrice(data, nil)
	if err != nil {
		return
	}

	tx := &types.Transaction{
		Format:   2,
		Target:   "",
		Quantity: "0",
		Tags:     types.TagsEncode(tags),
		Data:     data,
		DataSize: fmt.Sprintf("%d", len(data)),
		Reward:   fmt.Sprintf("%d", reward),
	}

	return w.SendTransaction(tx)
}

// SendTransaction: if send success, should return pending
func (w *Wallet) SendTransaction(tx *types.Transaction) (id, status string, err error) {
	anchor, err := w.Client.GetTransactionAnchor()
	if err != nil {
		return
	}
	tx.LastTx = anchor

	if err = tx.SignTransaction(w.PubKey, w.PrvKey); err != nil {
		return
	}

	id = tx.ID

	uploader, err := uploader.CreateUploader(w.Client, tx, nil)
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
