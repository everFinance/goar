package goar

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/daqiancode/goar/types"
	"github.com/daqiancode/goar/utils"
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
	reward, err := w.Client.GetTransactionPrice(0, target)
	if err != nil {
		return types.Transaction{}, err
	}
	tx := NewSendTokenTransaction(target, *amount, reward*(100+speedFactor)/100, tags...)
	return w.SendTransaction(tx)
}

func (w *Wallet) SendFile(file io.ReadSeeker, fileSize int64, tags ...types.Tag) (types.Transaction, error) {
	return w.SendFileSpeedUp(file, fileSize, 0, tags...)
}

// SendDataSpeedUp set speedFactor for speed up
// eg: speedFactor = 10, reward = 1.1 * reward
func (w *Wallet) SendFileSpeedUp(file io.ReadSeeker, fileSize int64, speedFactor int64, tags ...types.Tag) (types.Transaction, error) {
	reward, err := w.Client.GetTransactionPrice(fileSize)
	if err != nil {
		return types.Transaction{}, err
	}

	tx := NewSendFileTransaction(file, fileSize, reward*(100+speedFactor)/100, tags...)
	return w.SendTransaction(tx)
}

func (w *Wallet) SendData(data []byte, tags ...types.Tag) (types.Transaction, error) {
	return w.SendDataSpeedUp(data, 0, tags...)
}

// SendDataSpeedUp set speedFactor for speed up
// eg: speedFactor = 10, reward = 1.1 * reward
func (w *Wallet) SendDataSpeedUp(data []byte, speedFactor int64, tags ...types.Tag) (types.Transaction, error) {
	reward, err := w.Client.GetTransactionPrice(int64(len(data)))
	if err != nil {
		return types.Transaction{}, err
	}
	tx := NewSendFileTransaction(utils.NewReadBuffer(data), int64(len(data)), reward*(100+speedFactor)/100, tags...)
	return w.SendTransaction(tx)
}

func (w *Wallet) SendDataConcurrentSpeedUp(ctx context.Context, concurrentNum int, file io.ReadSeeker, fileSize int64, tags []types.Tag, speedFactor int64) (types.Transaction, error) {
	reward, err := w.Client.GetTransactionPrice(fileSize)
	if err != nil {
		return types.Transaction{}, err
	}

	tx := NewSendFileTransaction(file, fileSize, reward*(100+speedFactor)/100, tags...)
	return w.SendTransactionConcurrent(ctx, concurrentNum, tx)
}

// SendTransaction: if send success, should return pending
func (w *Wallet) SendTransaction(tx *types.Transaction) (types.Transaction, error) {
	uploader, err := w.getUploader(tx)
	if err != nil {
		return types.Transaction{}, err
	}
	err = uploader.Once()
	return *tx, err
}

func (w *Wallet) SendTransactionConcurrent(ctx context.Context, concurrentNum int, tx *types.Transaction) (types.Transaction, error) {
	uploader, err := w.getUploader(tx)
	if err != nil {
		return types.Transaction{}, err
	}
	err = uploader.ConcurrentOnce(ctx, concurrentNum)
	return *tx, err
}

func (w *Wallet) getUploader(tx *types.Transaction) (*TransactionUploader, error) {
	anchor, err := w.Client.GetTransactionAnchor()
	if err != nil {
		return nil, err
	}
	tx.LastTx = anchor
	tx.Owner = w.Owner()
	if err = w.Signer.SignTx(tx); err != nil {
		return nil, err
	}
	return CreateUploader(w.Client, tx, tx.File, tx.FileSize)
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
	return w.SendDataSpeedUp([]byte(data), speedFactor, swcTags...)
}

func (w *Wallet) SignTransaction(tx *types.Transaction) error {
	anchor, err := w.Client.GetTransactionAnchor()
	if err != nil {
		return err
	}
	tx.LastTx = anchor
	tx.Owner = w.Owner()
	signData, err := utils.GetSignatureData(tx)
	if err != nil {
		return err
	}
	sign, err := w.Signer.SignMsg(signData)
	if err != nil {
		return err
	}
	txHash := sha256.Sum256(sign)
	tx.ID = utils.Base64Encode(txHash[:])
	tx.Signature = utils.Base64Encode(sign)
	return nil
}

// NewFileTransaction https://docs.arweave.org/developers/server/http-api
func NewSendFileTransaction(file io.ReadSeeker, fileSize, reward int64, tags ...types.Tag) *types.Transaction {
	return &types.Transaction{
		Format:   2,
		Target:   "",
		Quantity: "0",
		Tags:     utils.TagsEncode(tags),
		DataSize: fmt.Sprintf("%d", fileSize),
		Reward:   fmt.Sprintf("%d", reward),
		File:     file,
		FileSize: fileSize,
	}
}

// NewSendTokenTransaction https://docs.arweave.org/developers/server/http-api
// Send a amount of tokens to target
func NewSendTokenTransaction(target string, amount big.Int, reward int64, tags ...types.Tag) *types.Transaction {
	return &types.Transaction{
		Format:   2,
		Target:   target,
		Quantity: amount.String(),
		Tags:     utils.TagsEncode(tags),
		Data:     "",
		DataSize: "0",
		Reward:   fmt.Sprintf("%d", reward),
	}
}
