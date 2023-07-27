package goar

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
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
	fmt.Println("local goar is run")
	return w.SendWinstonSpeedUp(utils.ARToWinston(amount), target, tags, 0)
}

func (w *Wallet) SendARSpeedUp(amount *big.Float, target string, tags []types.Tag, speedFactor int64) (types.Transaction, error) {
	return w.SendWinstonSpeedUp(utils.ARToWinston(amount), target, tags, speedFactor)
}

func (w *Wallet) SendWinston(amount *big.Int, target string, tags []types.Tag) (types.Transaction, error) {
	return w.SendWinstonSpeedUp(amount, target, tags, 0)
}

func (w *Wallet) SendWinstonSpeedUp(amount *big.Int, target string, tags []types.Tag, speedFactor int64) (types.Transaction, error) {
	reward, err := w.Client.GetTransactionPrice(0, &target)
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

func (w *Wallet) SendDataStream(data *os.File, tags []types.Tag) (types.Transaction, error) {
	return w.SendDataStreamSpeedUp(data, tags, 0)
}

// SendDataSpeedUp set speedFactor for speed up
// eg: speedFactor = 10, reward = 1.1 * reward
func (w *Wallet) SendDataSpeedUp(data []byte, tags []types.Tag, speedFactor int64) (types.Transaction, error) {
	reward, err := w.Client.GetTransactionPrice(len(data), nil)
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

func (w *Wallet) SendDataStreamSpeedUp(data *os.File, tags []types.Tag, speedFactor int64) (types.Transaction, error) {
	fileInfo, err := data.Stat()
	if err != nil {
		return types.Transaction{}, err
	}
	reward, err := w.Client.GetTransactionPrice(int(fileInfo.Size()), nil)
	if err != nil {
		return types.Transaction{}, err
	}

	tx := &types.Transaction{
		Format:     2,
		Target:     "",
		Quantity:   "0",
		Tags:       utils.TagsEncode(tags),
		Data:       "",
		DataReader: data,
		DataSize:   fmt.Sprintf("%d", fileInfo.Size()),
		Reward:     fmt.Sprintf("%d", reward*(100+speedFactor)/100),
	}

	return w.SendTransaction(tx)
}

func (w *Wallet) SendDataConcurrentSpeedUp(ctx context.Context, concurrentNum int, data interface{}, tags []types.Tag, speedFactor int64) (types.Transaction, error) {
	var reward int64
	var dataLen int
	isByteArr := true
	if _, isByteArr = data.([]byte); isByteArr {
		dataLen = len(data.([]byte))
	} else {
		fileInfo, err := data.(*os.File).Stat()
		if err != nil {
			return types.Transaction{}, err
		}
		dataLen = int(fileInfo.Size())
	}
	reward, err := w.Client.GetTransactionPrice(dataLen, nil)
	if err != nil {
		return types.Transaction{}, err
	}

	tx := &types.Transaction{
		Format:   2,
		Target:   "",
		Quantity: "0",
		Tags:     utils.TagsEncode(tags),
		DataSize: fmt.Sprintf("%d", dataLen),
		Reward:   fmt.Sprintf("%d", reward*(100+speedFactor)/100),
	}

	if isByteArr {
		tx.Data = utils.Base64Encode(data.([]byte))
	} else {
		tx.DataReader = data.(*os.File)
	}

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
	return CreateUploader(w.Client, tx, nil)
}

func (w *Wallet) SendPst(contractId string, target string, qty *big.Int, customTags []types.Tag, speedFactor int64) (types.Transaction, error) {
	maxQty := big.NewInt(9007199254740991) // swc support max js integer
	if qty.Cmp(maxQty) > 0 {
		return types.Transaction{}, fmt.Errorf("qty:%s can not more than max integer:%s", qty.String(), maxQty.String())
	}

	// assemble tx tags
	swcTags, err := utils.PstTransferTags(contractId, target, qty.Int64(), false)
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
