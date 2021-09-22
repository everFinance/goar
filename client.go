package goar

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/everFinance/sandy_log/log"

	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
)

// arweave HTTP API: https://docs.arweave.org/developers/server/http-api

type Client struct {
	client *http.Client
	url    string
}

func NewClient(nodeUrl string, proxyUrl ...string) *Client {
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
		return nil, ErrBadGateway
	}

	info = &types.NetworkInfo{}
	err = json.Unmarshal(body, info)
	return
}

func (c *Client) Peers() ([]string, error) {
	body, _, err := c.httpGet("peers")
	if err != nil {
		return nil, ErrBadGateway
	}

	peers := make([]string, 0)
	err = json.Unmarshal(body, &peers)
	return peers, err
}

// GetTransactionByID status: Pending/Invalid hash/overspend
func (c *Client) GetTransactionByID(id string) (tx *types.Transaction, err error) {
	body, statusCode, err := c.httpGet(fmt.Sprintf("tx/%s", id))
	if err != nil {
		return nil, ErrBadGateway
	}

	switch statusCode {
	case 200:
		// json unmarshal
		tx = &types.Transaction{}
		err = json.Unmarshal(body, tx)
		return
	case 202:
		return nil, ErrPendingTx
	case 400:
		return nil, ErrInvalidId
	case 404:
		return nil, ErrNotFound
	default:
		return nil, ErrBadGateway
	}
}

// GetTransactionStatus
func (c *Client) GetTransactionStatus(id string) (*types.TxStatus, error) {
	body, code, err := c.httpGet(fmt.Sprintf("tx/%s/status", id))
	if err != nil {
		return nil, ErrBadGateway
	}

	switch code {
	case 200:
		// json unmarshal
		txStatus := &types.TxStatus{}
		err = json.Unmarshal(body, txStatus)
		return txStatus, err
	case 202:
		return nil, ErrPendingTx
	case 404:
		return nil, ErrNotFound
	default:
		return nil, ErrBadGateway
	}
}

func (c *Client) GetTransactionField(id string, field string) (string, error) {
	body, statusCode, err := c.httpGet(fmt.Sprintf("tx/%v/%v", id, field))
	if err != nil {
		return "", ErrBadGateway
	}

	switch statusCode {
	case 200:
		return string(body), nil
	case 202:
		return "", ErrPendingTx
	case 400:
		return "", ErrInvalidId
	case 404:
		return "", ErrNotFound
	default:
		return "", ErrBadGateway
	}
}

func (c *Client) GetTransactionTags(id string) ([]types.Tag, error) {
	jsTags, err := c.GetTransactionField(id, "tags")
	if err != nil {
		return nil, err
	}

	tags := make([]types.Tag, 0)
	if err := json.Unmarshal([]byte(jsTags), &tags); err != nil {
		return nil, err
	}
	tags, err = utils.TagsDecode(tags)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (c *Client) GetTransactionData(id string, extension ...string) (body []byte, err error) {
	urlPath := fmt.Sprintf("tx/%v/%v", id, "data")
	if extension != nil {
		urlPath = urlPath + "." + extension[0]
	}
	body, statusCode, err := c.httpGet(urlPath)

	// When data is bigger than 12MiB statusCode == 400 NOTE: Data bigger than that has to be downloaded chunk by chunk.
	if statusCode == 400 || len(body) == 0 {
		body, err = c.DownloadChunkData(id)
		return
	} else if statusCode == 200 {
		return body, nil
	} else if statusCode == 202 {
		return nil, ErrPendingTx
	} else if statusCode == 404 {
		return nil, ErrNotFound
	} else {
		return nil, ErrBadGateway
	}
}

// GetTransactionDataByGateway
func (c *Client) GetTransactionDataByGateway(id string) (body []byte, err error) {
	urlPath := fmt.Sprintf("/%v/%v", id, "data")
	body, statusCode, err := c.httpGet(urlPath)
	switch statusCode {
	case 200:
		return body, nil
	case 400:
		return c.DownloadChunkData(id)
	case 202:
		return nil, ErrPendingTx
	case 404:
		return nil, ErrNotFound
	case 410:
		return nil, ErrInvalidId
	default:
		return nil, ErrBadGateway
	}
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

func (c *Client) getChunk(offset int64) (*types.TransactionChunk, error) {
	_path := "chunk/" + strconv.FormatInt(offset, 10)
	body, statusCode, err := c.httpGet(_path)
	if statusCode != 200 {
		return nil, errors.New("not found chunk data")
	}
	if err != nil {
		return nil, err
	}
	txChunk := &types.TransactionChunk{}
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

func (c *Client) getTransactionOffset(id string) (*types.TransactionOffset, error) {
	_path := fmt.Sprintf("tx/%s/offset", id)
	body, statusCode, err := c.httpGet(_path)
	if statusCode != 200 {
		return nil, errors.New("not found tx offset")
	}
	if err != nil {
		return nil, err
	}
	txOffset := &types.TransactionOffset{}
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

func (c *Client) GetTxDataFromPeers(txId string) ([]byte, error) {
	peers, err := c.Peers()
	if err != nil {
		return nil, err
	}
	for _, peer := range peers {
		if strings.Contains(peer, "127.0") {
			continue
		}
		arNode := NewClient("http://" + peer)
		data, err := arNode.GetTransactionData(txId)
		if err != nil {
			fmt.Printf("get tx data error:%v, peer: %s\n", err, peer)
			continue
		}
		return data, nil
	}
	return nil, errors.New("get tx data from peers failed")
}

func (c *Client) UploadTxDataToPeers(txId string, data []byte) error {
	peers, err := c.Peers()
	if err != nil {
		return err
	}

	count := 0
	for _, peer := range peers {
		if strings.Contains(peer, "127.0") {
			continue
		}
		fmt.Printf("upload peer: %s, count: %d\n", peer, count)
		arNode := NewClient("http://" + peer)
		uploader, err := CreateUploader(arNode, txId, data)
		if err != nil {
			continue
		}
	Loop:
		for !uploader.IsComplete() {
			if err := uploader.UploadChunk(); err != nil {
				break Loop
			}
			if uploader.LastResponseStatus != 200 {
				break Loop
			}
		}
		if uploader.IsComplete() { // upload success
			count++
		}
		if count > 20 {
			return nil
		}
	}
	return fmt.Errorf("upload tx data to peers failed, txId: %s", txId)
}

// push to bundler gateway

// SendItemToBundler send bundle bundleItem to bundler gateway
func (c *Client) SendItemToBundler(itemBinary []byte) (*types.BundlerResp, error) {
	// post to bundler
	resp, err := http.DefaultClient.Post(types.BUNDLER_HOST+"/tx", "application/octet-stream", bytes.NewReader(itemBinary))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("send to bundler request failed; http code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	// json unmarshal
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll(resp.Body) error: %v", err)
	}
	br := &types.BundlerResp{}
	if err := json.Unmarshal(body, br); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(body,br) failed; err: %v", err)
	}
	return br, nil
}

func (c *Client) BatchSendItemToBundler(bundleItems []types.BundleItem) ([]*types.BundlerResp, error) {
	respList := make([]*types.BundlerResp, 0, len(bundleItems))
	for _, item := range bundleItems {
		itemBinary := item.ItemBinary
		if len(itemBinary) == 0 {
			if err := utils.GenerateItemBinary(&item); err != nil {
				return nil, err
			}
			itemBinary = item.ItemBinary
		}
		resp, err := c.SendItemToBundler(itemBinary)
		if err != nil {
			return nil, err
		}
		respList = append(respList, resp)
	}
	return respList, nil
}

func (c *Client) GetBundle(arId string) (*types.Bundle, error) {
	data, err := c.DownloadChunkData(arId)
	if err != nil {
		return nil, err
	}
	return utils.DecodeBundle(data)
}
