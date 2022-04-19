package goar

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
)

type Wallet struct {
	Client *Client
	Signer *Signer
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
	signer, err := NewSigner(b)
	if err != nil {
		return nil, err
	}

	w = &Wallet{
		Client: NewClient(clientUrl, proxyUrl...),
		Signer: signer,
	}

	return
}

func (w *Wallet) Owner() string {
	return w.Signer.Owner()
}

func (w *Wallet) SendAR(amount *big.Float, target string, tags []types.Tag) (types.Transaction, error) {
	return w.SendWinstonSpeedUp(utils.ARToWinston(amount), target, tags, 0)
}

func (w *Wallet) SendARSpeedUp(amount *big.Float, target string, tags []types.Tag, speedFactor int64) (types.Transaction, error) {
	return w.SendWinstonSpeedUp(utils.ARToWinston(amount), target, tags, speedFactor)
}

func (w *Wallet) SendWinston(amount *big.Int, target string, tags []types.Tag) (types.Transaction, error) {
	return w.SendWinstonSpeedUp(amount, target, tags, 0)
}

func (w *Wallet) SendWinstonSpeedUp(amount *big.Int, target string, tags []types.Tag, speedFactor int64) (types.Transaction, error) {
	reward, err := w.Client.GetTransactionPrice(nil, &target)
	if err != nil {
		return types.Transaction{}, err
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

func (w *Wallet) SendData(data []byte, tags []types.Tag) (types.Transaction, error) {
	return w.SendDataSpeedUp(data, tags, 0)
}

// SendDataSpeedUp set speedFactor for speed up
// eg: speedFactor = 10, reward = 1.1 * reward
func (w *Wallet) SendDataSpeedUp(data []byte, tags []types.Tag, speedFactor int64) (types.Transaction, error) {
	reward, err := w.Client.GetTransactionPrice(data, nil)
	if err != nil {
		return types.Transaction{}, err
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
func (w *Wallet) SendTransaction(tx *types.Transaction) (types.Transaction, error) {
	anchor, err := w.Client.GetTransactionAnchor()
	if err != nil {
		return types.Transaction{}, err
	}
	tx.LastTx = anchor
	tx.Owner = w.Owner()
	if err = w.Signer.SignTx(tx); err != nil {
		return types.Transaction{}, err
	}

	uploader, err := CreateUploader(w.Client, tx, nil)
	if err != nil {
		return types.Transaction{}, err
	}
	err = uploader.Once()
	return *tx, err
}

func (w *Wallet) SendPst(contractId string, target string, qty *big.Int, customTags []types.Tag, speedFactor int64) (types.Transaction, error) {
	maxQty := big.NewInt(9007199254740991) // swc support max js integer
	if qty.Cmp(maxQty) > 0 {
		return types.Transaction{}, fmt.Errorf("qty:%s can not more than max integer:%s", qty.String(), maxQty.String())
	}

	// assemble tx tags
	swcTags, err := utils.PstTransferTags(contractId, target, qty.Int64())
	if err != nil {
		return types.Transaction{}, err
	}

	if len(customTags) > 0 {
		// customTags can not include pstTags
		mmap := map[string]struct{}{
			"App-Name":    {},
			"App-Version": {},
			"Contract":    {},
			"Input":       {},
		}
		for _, tag := range customTags {
			if _, ok := mmap[tag.Name]; ok {
				return types.Transaction{}, errors.New("custom tags can not include smartweave tags")
			}
		}
		swcTags = append(swcTags, customTags...)
	}

	// rand data
	data := strconv.Itoa(rand.Intn(9999))
	// send data tx
	return w.SendDataSpeedUp([]byte(data), swcTags, speedFactor)
}
