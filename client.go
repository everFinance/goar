package goar

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/inconshreveable/log15"
	"github.com/panjf2000/ants/v2"
	"github.com/tidwall/gjson"
	"gopkg.in/h2non/gentleman.v2"

	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
)

var log = log15.New("module", "goar")

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
			log.Error("url parse", "error", err)
			panic(err)
		}
		tr := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		httpClient = &http.Client{Transport: tr}
	}

	return &Client{client: httpClient, url: nodeUrl}
}

func NewTempConn() *Client {
	transport := http.Transport{DisableKeepAlives: true}
	cli := &http.Client{Transport: &transport}
	return &Client{client: cli}
}

func (c *Client) SetTempConnUrl(url string) {
	c.url = url
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

func (c *Client) GetInfo() (info *types.NetworkInfo, err error) {
	body, code, err := c.httpGet("info")
	if code == 429 {
		return nil, ErrRequestLimit
	}
	if err != nil {
		return nil, ErrBadGateway
	}
	if code != 200 {
		return nil, fmt.Errorf("get info error: %s", string(body))
	}

	info = &types.NetworkInfo{}
	err = json.Unmarshal(body, info)
	return
}

func (c *Client) GetPeers() ([]string, error) {
	body, code, err := c.httpGet("peers")
	if code == 429 {
		return nil, ErrRequestLimit
	}
	if err != nil {
		return nil, ErrBadGateway
	}
	if code != 200 {
		return nil, fmt.Errorf("get peers error: %s", string(body))
	}

	peers := make([]string, 0)
	err = json.Unmarshal(body, &peers)
	if err != nil {
		return nil, err
	}

	// filter local
	fpeers := make([]string, 0)
	for _, p := range peers {
		if strings.Contains(p, "127.0.0.") {
			continue
		}
		fpeers = append(fpeers, p)
	}

	return fpeers, nil
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
	case 429:
		return nil, ErrRequestLimit
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
	case 429:
		return nil, ErrRequestLimit
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
	case 429:
		return "", ErrRequestLimit
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

func (c *Client) GetTransactionData(id string, extension ...string) ([]byte, error) {
	urlPath := fmt.Sprintf("tx/%v/%v", id, "data")
	if extension != nil {
		urlPath = urlPath + "." + extension[0]
	}
	data, statusCode, err := c.httpGet(urlPath)
	if err != nil {
		return nil, fmt.Errorf("httpGet error: %v", err)
	}

	// When data is bigger than 12MiB statusCode == 400 NOTE: Data bigger than that has to be downloaded chunk by chunk.
	switch statusCode {
	case 200:
		if len(data) == 0 {
			return c.DownloadChunkData(id)
		}
		return data, nil
	case 400:
		return c.DownloadChunkData(id)
	case 202:
		return nil, ErrPendingTx
	case 404:
		return nil, ErrNotFound
	case 429:
		return nil, ErrRequestLimit
	default:
		return nil, ErrBadGateway
	}
}

func (c *Client) GetTransactionDataStream(id string, extension ...string) (*os.File, error) {
	urlPath := fmt.Sprintf("tx/%v/%v", id, "data")
	if extension != nil {
		urlPath = urlPath + "." + extension[0]
	}
	data, statusCode, err := c.httpGet(urlPath)
	if err != nil {
		return nil, fmt.Errorf("httpGet error: %v", err)
	}

	// When data is bigger than 12MiB statusCode == 400 NOTE: Data bigger than that has to be downloaded chunk by chunk.
	switch statusCode {
	case 200:
		if len(data) == 0 {
			return c.DownloadChunkDataStream(id)
		}
		dataFile, err := os.CreateTemp(".", "arTxData-")
		if err != nil {
			return nil, err
		}
		_, err = dataFile.Write(data)
		return dataFile, err
	case 400:
		return c.DownloadChunkDataStream(id)
	case 202:
		return nil, ErrPendingTx
	case 404:
		return nil, ErrNotFound
	case 429:
		return nil, ErrRequestLimit
	default:
		return nil, ErrBadGateway
	}
}

// GetTransactionDataByGateway
func (c *Client) GetTransactionDataByGateway(id string) (body []byte, err error) {
	urlPath := fmt.Sprintf("/%v/%v", id, "data")
	body, statusCode, err := c.httpGet(urlPath)
	if err != nil {
		return nil, fmt.Errorf("httpGet error: %v", err)
	}
	switch statusCode {
	case 200:
		if len(body) == 0 {
			return c.DownloadChunkData(id)
		}
		return body, nil
	case 400:
		return c.DownloadChunkData(id)
	case 202:
		return nil, ErrPendingTx
	case 404:
		return nil, ErrNotFound
	case 410:
		return nil, ErrInvalidId
	case 429:
		return nil, ErrRequestLimit
	default:
		return nil, ErrBadGateway
	}
}

func (c *Client) GetTransactionDataStreamByGateway(id string) (*os.File, error) {
	urlPath := fmt.Sprintf("/%v/%v", id, "data")
	body, statusCode, err := c.httpGet(urlPath)
	if err != nil {
		return nil, fmt.Errorf("httpGet error: %v", err)
	}
	switch statusCode {
	case 200:
		if len(body) == 0 {
			return c.DownloadChunkDataStream(id)
		}
		dataFile, err := os.CreateTemp(".", "arTxData-")
		if err != nil {
			return nil, err
		}
		_, err = dataFile.Write(body)
		return dataFile, err
	case 400:
		return c.DownloadChunkDataStream(id)
	case 202:
		return nil, ErrPendingTx
	case 404:
		return nil, ErrNotFound
	case 410:
		return nil, ErrInvalidId
	case 429:
		return nil, ErrRequestLimit
	default:
		return nil, ErrBadGateway
	}
}

func (c *Client) GetTransactionPrice(dataSize int, target *string) (reward int64, err error) {
	url := fmt.Sprintf("price/%d", dataSize)
	if target != nil {
		url = fmt.Sprintf("%v/%v", url, *target)
	}

	body, code, err := c.httpGet(url)
	if code == 429 {
		return 0, ErrRequestLimit
	}
	if err != nil {
		return
	}
	if code != 200 {
		return 0, fmt.Errorf("get reward error: %s", string(body))
	}

	reward, err = strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		return
	}

	// reward can not be 0
	if reward <= 0 {
		err = errors.New("reward must more than 0")
	}
	return
}

func (c *Client) GetTransactionAnchor() (anchor string, err error) {
	body, code, err := c.httpGet("tx_anchor")
	if code == 429 {
		return "", ErrRequestLimit
	}
	if err != nil {
		return
	}
	if code != 200 {
		return "", fmt.Errorf("get tx anchor err: %s", string(body))
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
	if statusCode == 429 {
		return nil, ErrRequestLimit
	}
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
	winstonAmt, err := c.GetWalletWinstonBalance(address)
	if err != nil {
		return nil, err
	}

	arAmount = utils.WinstonToAR(winstonAmt)
	return
}

func (c *Client) GetWalletWinstonBalance(address string) (arAmount *big.Int, err error) {
	body, code, err := c.httpGet(fmt.Sprintf("wallet/%s/balance", address))
	if code == 429 {
		return nil, ErrRequestLimit
	}
	if err != nil {
		return
	}
	if code != 200 {
		return nil, fmt.Errorf("get balance error: %s", string(body))
	}

	winstomStr := string(body)
	winstom, ok := new(big.Int).SetString(winstomStr, 10)
	if !ok {
		err = fmt.Errorf("invalid balance: %v", winstomStr)
		return
	}
	arAmount = winstom
	return
}

func (c *Client) GetLastTransactionID(address string) (id string, err error) {
	body, code, err := c.httpGet(fmt.Sprintf("wallet/%s/last_tx", address))
	if code == 429 {
		return "", ErrRequestLimit
	}
	if err != nil {
		return
	}
	if code != 200 {
		return "", fmt.Errorf("get last id error: %s", string(body))
	}

	id = string(body)
	return
}

// Block
func (c *Client) GetBlockByID(id string) (block *types.Block, err error) {
	body, code, err := c.httpGet(fmt.Sprintf("block/hash/%s", id))
	if err != nil {
		return
	}
	if code == 429 {
		return nil, ErrRequestLimit
	}

	if code != 200 {
		return nil, fmt.Errorf("get block by id error: %s", string(body))
	}
	block, err = utils.DecodeBlock(string(body))
	return
}

func (c *Client) GetBlockByHeight(height int64) (block *types.Block, err error) {
	body, code, err := c.httpGet(fmt.Sprintf("block/height/%d", height))
	if err != nil {
		return
	}
	if code == 429 {
		return nil, ErrRequestLimit
	}

	if code != 200 {
		return nil, fmt.Errorf("get block by height error: %s", string(body))
	}
	block, err = utils.DecodeBlock(string(body))
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
	body, err = io.ReadAll(resp.Body)
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
	body, err = io.ReadAll(resp.Body)
	return
}

// about chunk

func (c *Client) getChunk(offset int64) (*types.TransactionChunk, error) {
	_path := "chunk/" + strconv.FormatInt(offset, 10)
	body, statusCode, err := c.httpGet(_path)
	if err != nil {
		return nil, fmt.Errorf("httpGet getChunk error: %v", err)
	}

	switch statusCode {
	case 200:
		txChunk := &types.TransactionChunk{}
		if err := json.Unmarshal(body, txChunk); err != nil {
			return nil, err
		}
		return txChunk, nil
	case 404:
		return nil, ErrNotFound
	case 429:
		return nil, ErrRequestLimit
	default:
		return nil, ErrBadGateway
	}
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
	if statusCode == 429 {
		return nil, ErrRequestLimit
	}
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

// it's caller's responsibility to reserve or delete the tmp file created by this method

func (c *Client) DownloadChunkDataStream(id string) (*os.File, error) {
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
	dataFile, err := os.CreateTemp(".", "chunkData-")
	if err != nil {
		return nil, err
	}
	downloadSize := 0
	n := 0
	for i := 0; int64(i)+startOffset < endOffset; {
		chunkData, err := c.getChunkData(int64(i) + startOffset)
		if err != nil {
			return nil, err
		}
		downloadSize += len(chunkData)
		n, err = dataFile.Write(chunkData)
		if err != nil || n < len(chunkData) {
			return nil, fmt.Errorf("write chunkData to dataFile failed")
		}
		fmt.Printf("download chunk data; offset: %d/%d; size: %d/%d \n", int64(i)+startOffset, endOffset, downloadSize, size)
		i += len(chunkData)
	}
	return dataFile, nil
}

func (c *Client) ConcurrentDownloadChunkData(id string, concurrentNum int) ([]byte, error) {
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

	offsetArr := make([]int64, 0, 5)
	for i := 0; int64(i)+startOffset < endOffset; {
		offsetArr = append(offsetArr, int64(i)+startOffset)
		i += types.MAX_CHUNK_SIZE
	}

	if len(offsetArr) <= 3 { // not need concurrent get chunks
		return c.DownloadChunkData(id)
	}

	log.Debug("need download chunks length", "length", len(offsetArr))
	// concurrent get chunks
	type OffsetSort struct {
		Idx    int
		Offset int64
	}

	chunkArr := make([][]byte, len(offsetArr)-2)
	var (
		lock sync.Mutex
		wg   sync.WaitGroup
	)
	if concurrentNum <= 0 {
		concurrentNum = types.DEFAULT_CHUNK_CONCURRENT_NUM
	}
	p, _ := ants.NewPoolWithFunc(concurrentNum, func(i interface{}) {
		defer wg.Done()
		oss := i.(OffsetSort)
		chunkData, err := c.getChunkData(oss.Offset)
		if err != nil {
			count := 0
			for count < 2 {
				time.Sleep(1 * time.Second)
				chunkData, err = c.getChunkData(oss.Offset)
				if err == nil {
					break
				}
				log.Error("retry getChunkData failed and try again...", "err", err, "idx", oss.Idx, "offset", oss.Offset, "retryCount", count, "arId", id)
				if err != ErrRequestLimit {
					count++
				}
			}
		}
		lock.Lock()
		chunkArr[oss.Idx] = chunkData
		lock.Unlock()
	})

	defer p.Release()

	for i, offset := range offsetArr[:len(offsetArr)-2] {
		wg.Add(1)
		if err := p.Invoke(OffsetSort{Idx: i, Offset: offset}); err != nil {
			log.Error("p.Invoke(i)", "err", err, "i", i)
			return nil, err
		}
	}
	wg.Wait()

	// add latest 2 chunks
	start := offsetArr[len(offsetArr)-3] + types.MAX_CHUNK_SIZE
	for i := 0; int64(i)+start < endOffset; {
		chunkData, err := c.getChunkData(int64(i) + start)
		if err != nil {
			count := 0
			for count < 2 {
				time.Sleep(1 * time.Second)
				chunkData, err = c.getChunkData(int64(i) + start)
				if err == nil {
					break
				}
				log.Error("latest two chunks retry getChunkData failed and try again...", "err", err, "offset", int64(i)+start, "retryCount", count, "arId", id)
				if err != ErrRequestLimit {
					count++
				}
			}
		}
		if err != nil {
			return nil, errors.New("concurrent get latest two chunks failed")
		}
		chunkArr = append(chunkArr, chunkData)
		i += len(chunkData)
	}

	// assemble data
	data := make([]byte, 0, size)
	for _, chunk := range chunkArr {
		if chunk == nil {
			return nil, errors.New("concurrent get chunk failed")
		}
		data = append(data, chunk...)
	}
	return data, nil
}

// it's caller's responsibility to reserve or delete the tmp file created by this method

func (c *Client) ConcurrentDownloadChunkDataStream(id string, concurrentNum int) (dataFile *os.File, err error) {
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

	offsetArr := make([]int64, 0, 5)
	for i := 0; int64(i)+startOffset < endOffset; {
		offsetArr = append(offsetArr, int64(i))
		i += types.MAX_CHUNK_SIZE
	}

	log.Debug("need download chunks length", "length", len(offsetArr))

	dataFile, err = os.CreateTemp(".", "concurrent-load-data-")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			dataFile.Close()
			os.Remove(dataFile.Name())
		}
	}()

	if len(offsetArr) <= 3 { // not need concurrent get chunks
		var data []byte
		data, err = c.DownloadChunkData(id)
		if err != nil {
			return
		}
		_, err = dataFile.Write(data)
		return dataFile, err
	}

	type Offset struct {
		fileOffset  int64
		chunkOffset int64
	}
	var (
		lock sync.Mutex
		wg   sync.WaitGroup
	)
	if concurrentNum <= 0 {
		concurrentNum = types.DEFAULT_CHUNK_CONCURRENT_NUM
	}
	p, _ := ants.NewPoolWithFunc(concurrentNum, func(i interface{}) {
		defer wg.Done()
		oss := i.(Offset)
		chunkData, err := c.getChunkData(oss.chunkOffset)
		if err != nil {
			count := 0
			for count < 5 {
				time.Sleep(1 * time.Second)
				chunkData, err = c.getChunkData(oss.chunkOffset)
				if err == nil {
					break
				}
				log.Warn("retry getChunkData failed and try again...", "err", err, "idx", oss.fileOffset/types.MAX_CHUNK_SIZE, "offset", oss.chunkOffset, "retryCount", count, "arId", id)
				if err != ErrRequestLimit {
					count++
				}
			}
		}
		if err != nil {
			log.Error("getChunkData failed", "err", err, "arId", id, "idx", oss.fileOffset/types.MAX_CHUNK_SIZE, "offset", oss.chunkOffset)
			return
		}
		var n int
		lock.Lock()
		n, err = dataFile.WriteAt(chunkData, oss.fileOffset)
		if err != nil || n < len(chunkData) {
			log.Error("write dataFile error")
		}
		lock.Unlock()
	})

	defer p.Release()

	for i, offset := range offsetArr[:len(offsetArr)-2] {
		wg.Add(1)
		if err = p.Invoke(Offset{fileOffset: offset, chunkOffset: offset + startOffset}); err != nil {
			log.Error("p.Invoke(i)", "err", err, "i", i)
			return
		}
	}
	wg.Wait()
	_, err = dataFile.Seek(0, 2)
	if err != nil {
		return
	}
	// add latest 2 chunks
	start := offsetArr[len(offsetArr)-3] + startOffset + types.MAX_CHUNK_SIZE
	for i := 0; int64(i)+start < endOffset; {
		var chunkData []byte
		chunkData, err = c.getChunkData(int64(i) + start)
		if err != nil {
			count := 0
			for count < 5 {
				time.Sleep(1 * time.Second)
				chunkData, err = c.getChunkData(int64(i) + start)
				if err == nil {
					break
				}
				log.Error("latest two chunks retry getChunkData failed and try again...", "err", err, "offset", int64(i)+start, "retryCount", count, "arId", id)
				if err != ErrRequestLimit {
					count++
				}
			}
		}
		if err != nil {
			err = errors.New(fmt.Sprintf("concurrent get latest two chunks failed,err:%v", err))
			return
		}
		n := 0
		n, err = dataFile.Write(chunkData)
		if err != nil || n < len(chunkData) {
			err = fmt.Errorf("write dataFile error writeSize:%d, expectSize:%d", n, len(chunkData))
			return
		}
		i += len(chunkData)
	}

	return dataFile, nil
}

func (c *Client) GetUnconfirmedTx(arId string) (*types.Transaction, error) {
	_path := fmt.Sprintf("unconfirmed_tx/%s", arId)
	body, statusCode, err := c.httpGet(_path)
	if statusCode != 200 {
		return nil, errors.New("not found unconfirmed tx")
	}
	if err != nil {
		return nil, err
	}
	tx := &types.Transaction{}
	if err := json.Unmarshal(body, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (c *Client) GetPendingTxIds() ([]string, error) {
	body, statusCode, err := c.httpGet("/tx/pending")
	if statusCode != 200 {
		return nil, errors.New("get pending txIds failed")
	}
	if err != nil {
		return nil, err
	}
	res := make([]string, 0)
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetBlockHashList(from, to int) ([]string, error) {
	if from > to {
		return nil, errors.New("from must <= to")
	}
	body, statusCode, err := c.httpGet("/hash_list/" + strconv.Itoa(from) + "/" + strconv.Itoa(to))
	if statusCode != 200 {
		return nil, errors.New("get block hash list failed")
	}
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) ExistTxData(arId string) (bool, error) {
	offsetResponse, err := c.getTransactionOffset(arId)
	if err != nil {
		return false, err
	}
	endOffset := offsetResponse.Offset

	records, err := c.DataSyncRecord(endOffset, 1)
	if err != nil {
		return false, err
	}
	if len(records) == 0 {
		return false, errors.New("c.DataSyncRecord(endOffset,1) is null")
	}
	record := records[0]

	// if tx data has end offset 145 and size 10 (you can see it in GET /tx/<id>/offset),
	// you can query GET /data_sync_record/145/1
	// - you will receive {"<end>": "<start>"} => the node has the tx data if start =< 145 - 10
	mmp := gjson.Parse(record).Map()
	start := ""
	for _, val := range mmp {
		start = val.String()
		break
	}
	startNum, err := strconv.Atoi(start)
	if err != nil {
		return false, err
	}
	endOffsetNum, err := strconv.Atoi(endOffset)
	if err != nil {
		return false, err
	}
	sizeNum, err := strconv.Atoi(offsetResponse.Size)
	if err != nil {
		return false, err
	}

	return startNum <= endOffsetNum-sizeNum, nil
}

// DataSyncRecord you can use GET /data_sync_record/<end_offset>/<number_of_intervals>
// to fetch the first intervals with end offset >= end_offset;
// set Content-Type: application/json to get the reply in JSON
func (c *Client) DataSyncRecord(endOffset string, intervalsNum int) ([]string, error) {
	req := gentleman.New().URL(c.url).Request()
	req.AddPath("/data_sync_record/" + endOffset + "/" + strconv.Itoa(intervalsNum))
	req.SetHeader("Content-Type", "application/json")
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 429 {
		return nil, ErrRequestLimit
	}
	if !resp.Ok {
		return nil, errors.New("resp ok is false")
	}
	defer resp.Close()
	ss := gjson.ParseBytes(resp.Bytes()).Array()
	result := make([]string, 0, len(ss))
	for _, s := range ss {
		result = append(result, s.String())
	}
	return result, nil
}

func (c *Client) SubmitToWarp(tx *types.Transaction) ([]byte, error) {
	by, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(c.url)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "/gateway/sequencer/register")
	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(by))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

/**

 bundleBinary Data Format
+------------------+-----------------------------+--------------------------------------------+
| Number of Items  |      Headers (64 bytes)     |          Items' Binary Data                |
|      32 bytes    | 32 bytes (length) + 32 bytes |                                            |
|                  |       (ID) for each item    |                                            |
+------------------+-----------------------------+--------------------------------------------+
*/

func (c *Client) GetBundleItems(bundleInId string, itemsIds []string) (items []*types.BundleItem, err error) {
	offset, err := c.getTransactionOffset(bundleInId)
	if err != nil {
		return nil, err
	}

	size, err := strconv.ParseInt(offset.Size, 10, 64)
	if err != nil {
		return nil, err
	}
	endOffset, err := strconv.ParseInt(offset.Offset, 10, 64)
	startOffset := endOffset - size + 1

	firstChunk, err := c.getChunkData(startOffset)
	if err != nil {
		return nil, err
	}

	itemsNum := utils.ByteArrayToLong(firstChunk[:32])

	// get Headers endOffset
	bundleItemStart := 32 + itemsNum*64
	var containHeadersChunks []byte
	if len(firstChunk) < bundleItemStart {
		// To calculate headers, you need to pull several chunks and fetch an integer upwards
		chunkNum := int(math.Ceil(float64(bundleItemStart) / float64(types.MAX_CHUNK_SIZE)))

		for i := 0; i < chunkNum; i++ {
			chunk, err := c.getChunkData(startOffset + int64(i*types.MAX_CHUNK_SIZE))
			if err != nil {
				return nil, err
			}
			containHeadersChunks = append(containHeadersChunks, chunk...)
		}
	} else {
		containHeadersChunks = firstChunk
	}

	for i := 0; i < itemsNum; i++ {
		headerBegin := 32 + i*64
		end := headerBegin + 64

		headerByte := containHeadersChunks[headerBegin:end]
		itemBinaryLength := utils.ByteArrayToLong(headerByte[:32])
		id := utils.Base64Encode(headerByte[32:64])

		// if item is in itemsIds
		if utils.ContainsInSlice(itemsIds, id) {
			startChunkNum := bundleItemStart / types.MAX_CHUNK_SIZE
			startChunkOffset := bundleItemStart % types.MAX_CHUNK_SIZE
			data := make([]byte, 0, itemBinaryLength)

			for offset := startOffset + int64(startChunkNum*types.MAX_CHUNK_SIZE); offset <= startOffset+int64(bundleItemStart+itemBinaryLength); {
				if offset >= endOffset {
					break
				}
				chunk, err := c.getChunkData(offset)

				if err != nil {
					return nil, err
				}
				data = append(data, chunk...)
				offset += int64(len(chunk))
			}

			itemData := data[startChunkOffset : startChunkOffset+itemBinaryLength]
			item, err := utils.DecodeBundleItem(itemData)
			if err != nil {
				return nil, err
			}

			items = append(items, item)
		}
		// next itemBy start offset
		bundleItemStart += itemBinaryLength
	}

	return items, nil

}
