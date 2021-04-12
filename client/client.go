package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zyjblockchain/sandy_log/log"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
)

type Client struct {
	client *http.Client
	url    string
}

func New(nodeUrl string, proxyUrl ...string) *Client {
	httpClient := http.DefaultClient
	// if exist proxy url
	if len(proxyUrl) > 0 {
		pUrl := proxyUrl[0]
		proxyUrl, err := url.Parse(pUrl)
		if err != nil {
			log.Errorf("url parse error: %v", err)
			panic(err)
		}
		tr := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		httpClient = &http.Client{Transport: tr}
	}

	return &Client{client: httpClient, url: nodeUrl}
}

func (c *Client) GetInfo() (info *types.NetworkInfo, err error) {
	body, _, err := c.httpGet("info")
	if err != nil {
		return
	}

	info = &types.NetworkInfo{}
	err = json.Unmarshal(body, info)
	return
}

// Transaction
// status: Pending/Invalid hash/overspend
func (c *Client) GetTransactionByID(id string) (tx *types.Transaction, status string, code int, err error) {
	body, statusCode, err := c.httpGet(fmt.Sprintf("tx/%s", id))
	if err != nil {
		return
	}

	code = statusCode
	if statusCode != 200 {
		status = string(body)
		return
	}

	tx = &types.Transaction{}
	err = json.Unmarshal(body, tx)
	return
}

// GetTransactionStatus
func (c *Client) GetTransactionStatus(id string) (status string, code int, err error) {
	body, code, err := c.httpGet(fmt.Sprintf("tx/%s/status", id))
	if code == 200 {
		return types.SuccessTxStatus, code, nil
	} else if code == 202 {
		return types.PendingTxStatus, code, nil
	} else {
		return string(body), code, err
	}
}

func (c *Client) GetTransactionField(id string, field string) (f string, err error) {
	url := fmt.Sprintf("tx/%v/%v", id, field)

	body, statusCode, err := c.httpGet(url)
	if statusCode != 200 {
		err = fmt.Errorf("not found data")
	}

	f = string(body)
	return
}

func (c *Client) GetTransactionData(id string, extension ...string) (body []byte, err error) {
	url := fmt.Sprintf("tx/%v/%v", id, "data")
	if extension != nil {
		url = url + "." + extension[0]
	}
	body, statusCode, err := c.httpGet(url)

	if statusCode == 400 {
		body, err = c.DownloadChunkData(id)
	} else if statusCode != 200 {
		err = fmt.Errorf("not found data")
	}

	return
}

func (c *Client) GetTransactionPrice(data []byte, target *string) (reward int64, err error) {
	url := fmt.Sprintf("price/%d", len(data))
	if target != nil {
		url = fmt.Sprintf("%v/%v", url, *target)
	}

	body, _, err := c.httpGet(url)
	if err != nil {
		return
	}

	return strconv.ParseInt(string(body), 10, 64)
}

func (c *Client) GetTransactionAnchor() (anchor string, err error) {
	body, _, err := c.httpGet("tx_anchor")
	if err != nil {
		return
	}

	anchor = string(body)
	return
}

func (c *Client) SubmitTransaction(tx *types.Transaction) (status string, code int, err error) {
	by, err := json.Marshal(tx)
	if err != nil {
		return
	}

	body, statusCode, err := c.httpPost("tx", by)
	status = string(body)
	code = statusCode
	return
}

func (c *Client) SubmitChunks(gc *types.GetChunk) (status string, code int, err error) {
	byteGc, err := gc.Marshal()
	if err != nil {
		return
	}

	var body []byte
	body, code, err = c.httpPost("chunk", byteGc)
	status = string(body)
	return
}

// Arql is Deprecated, recommended to use GraphQL
func (c *Client) Arql(arql string) (ids []string, err error) {
	body, _, err := c.httpPost("arql", []byte(arql))
	err = json.Unmarshal(body, &ids)
	return
}

func (c *Client) GraphQL(query string) ([]byte, error) {
	// generate query
	graQuery := struct {
		Query string `json:"query"`
	}{query}
	byQuery, err := json.Marshal(graQuery)
	if err != nil {
		return nil, err
	}

	// query from http client
	data, statusCode, err := c.httpPost("graphql", byQuery)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf(string(data))
	}

	// unwrap data
	res := struct {
		Data interface{}
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return json.Marshal(res.Data)
}

// Wallet
func (c *Client) GetWalletBalance(address string) (arAmount *big.Float, err error) {
	body, _, err := c.httpGet(fmt.Sprintf("wallet/%s/balance", address))
	if err != nil {
		return
	}

	winstomStr := string(body)
	winstom, ok := new(big.Int).SetString(winstomStr, 10)
	if !ok {
		err = fmt.Errorf("invalid balance: %v", winstomStr)
		return
	}

	arAmount = utils.WinstonToAR(winstom)
	return
}

func (c *Client) GetLastTransactionID(address string) (id string, err error) {
	body, _, err := c.httpGet(fmt.Sprintf("wallet/%s/last_tx", address))
	if err != nil {
		return
	}

	id = string(body)
	return
}

// Block
func (c *Client) GetBlockByID(id string) (block *types.Block, err error) {
	body, _, err := c.httpGet(fmt.Sprintf("block/hash/%s", id))
	if err != nil {
		return
	}

	block = &types.Block{}
	err = json.Unmarshal(body, block)
	return
}

func (c *Client) GetBlockByHeight(height int64) (block *types.Block, err error) {
	body, _, err := c.httpGet(fmt.Sprintf("block/height/%d", height))
	if err != nil {
		return
	}

	block = &types.Block{}
	err = json.Unmarshal(body, block)
	return
}

func (c *Client) httpGet(_path string) (body []byte, statusCode int, err error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return
	}

	u.Path = path.Join(u.Path, _path)

	resp, err := c.client.Get(u.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func (c *Client) httpPost(_path string, payload []byte) (body []byte, statusCode int, err error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return
	}

	u.Path = path.Join(u.Path, _path)

	resp, err := c.client.Post(u.String(), "application/json", bytes.NewReader(payload))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)
	return
}

// about chunk

func (c *Client) getChunk(offset int64) (*types.TransactionChunkResponse, error) {
	_path := "chunk/" + strconv.FormatInt(offset, 10)
	body, statusCode, err := c.httpGet(_path)
	if statusCode != 200 {
		return nil, errors.New("not found chunk data")
	}
	if err != nil {
		return nil, err
	}
	txChunk := &types.TransactionChunkResponse{}
	if err := json.Unmarshal(body, txChunk); err != nil {
		return nil, err
	}
	return txChunk, nil
}

func (c *Client) getChunkData(offset int64) ([]byte, error) {
	chunk, err := c.getChunk(offset)
	if err != nil {
		return nil, err
	}
	return utils.Base64Decode(chunk.Chunk)
}

func (c *Client) getTransactionOffset(id string) (*types.TransactionOffsetResponse, error) {
	_path := fmt.Sprintf("tx/%s/offset", id)
	body, statusCode, err := c.httpGet(_path)
	if statusCode != 200 {
		return nil, errors.New("not found tx offset")
	}
	if err != nil {
		return nil, err
	}
	txOffset := &types.TransactionOffsetResponse{}
	if err := json.Unmarshal(body, txOffset); err != nil {
		return nil, err
	}
	return txOffset, nil
}

func (c *Client) DownloadChunkData(id string) ([]byte, error) {
	offsetResponse, err := c.getTransactionOffset(id)
	if err != nil {
		return nil, err
	}
	size, err := strconv.ParseInt(offsetResponse.Size, 10, 64)
	if err != nil {
		return nil, err
	}
	endOffset, err := strconv.ParseInt(offsetResponse.Offset, 10, 64)
	if err != nil {
		return nil, err
	}
	startOffset := endOffset - size + 1
	data := make([]byte, 0, size)
	for i := 0; int64(i)+startOffset < endOffset; {
		chunkData, err := c.getChunkData(int64(i) + startOffset)
		if err != nil {
			return nil, err
		}
		data = append(data, chunkData...)
		fmt.Printf("download chunk data; offset: %d/%d; size: %d/%d \n", int64(i)+startOffset, endOffset, len(data), size)
		i += len(chunkData)
	}
	return data, nil
}
