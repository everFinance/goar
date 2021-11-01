package goar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
)

func (c *Client) GetBundle(arId string) (*types.Bundle, error) {
	data, err := c.DownloadChunkData(arId)
	if err != nil {
		return nil, err
	}
	return utils.DecodeBundle(data)
}

// SendItemToBundler send bundle bundleItem to bundler gateway
func (c *Client) SendItemToBundler(itemBinary []byte, gateway string) (*types.BundlerResp, error) {
	if gateway == "" {
		gateway = types.BUNDLER_HOST
	}
	// post to bundler
	resp, err := http.DefaultClient.Post(gateway+"/tx", "application/octet-stream", bytes.NewReader(itemBinary))
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

func (c *Client) BatchSendItemToBundler(bundleItems []types.BundleItem, gateway string) ([]*types.BundlerResp, error) {
	respList := make([]*types.BundlerResp, 0, len(bundleItems))
	for _, item := range bundleItems {
		itemBinary := item.ItemBinary
		if len(itemBinary) == 0 {
			if err := utils.GenerateItemBinary(&item); err != nil {
				return nil, err
			}
			itemBinary = item.ItemBinary
		}
		resp, err := c.SendItemToBundler(itemBinary, gateway)
		if err != nil {
			return nil, err
		}
		respList = append(respList, resp)
	}
	return respList, nil
}
