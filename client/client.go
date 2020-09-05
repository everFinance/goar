package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/everFinance/goar/types"
)

type Client struct {
	client *http.Client
	url    string
}

func New(url string) *Client {
	return &Client{client: http.DefaultClient, url: url}
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
func (c *Client) GetTransactionByID(id string) (tx *types.Transaction, status string, err error) {
	body, statusCode, err := c.httpGet(fmt.Sprintf("tx/%s", id))
	if err != nil {
		return
	}

	if statusCode != 200 {
		status = string(body)
		return
	}

	tx = &types.Transaction{}
	err = json.Unmarshal(body, tx)
	return
}

func (c *Client) GetTransactionField(id string, field string) (f string, err error) {
	// TODO
	return
}

func (c *Client) GetTransactionData(id string, extension string) (body []byte, err error) {
	url := fmt.Sprintf("tx/%v/%v", id, "data")
	if extension != "" {
		url = url + "." + extension
	}

	body, statusCode, err := c.httpGet(url)
	if statusCode != 200 {
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

func (c *Client) SubmitTransaction(tx *types.Transaction) (status string, err error) {
	by, err := json.Marshal(tx)
	if err != nil {
		return
	}

	body, _, err := c.httpPost("tx", by)
	status = string(body)
	return
}

func (c *Client) Arql(arql string) (ids []string, err error) {
	body, _, err := c.httpPost("arql", []byte(arql))
	err = json.Unmarshal(body, &ids)
	return
}

// Wallet
func (c *Client) GetWalletBalance(address string) (amount string, err error) {
	// TODO
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
