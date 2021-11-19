package goar

import (
	"errors"
	"fmt"
	"github.com/everFinance/goar/types"
)

func (c *Client) BroadcastData(txId string, data []byte, numOfNodes int64) error {
	peers, err := c.GetPeers()
	if err != nil {
		return err
	}

	count := int64(0)
	for _, peer := range peers {
		arNode := NewClient("http://" + peer)
		uploader, err := CreateUploader(arNode, txId, data)
		if err != nil {
			continue
		}

		if err = uploader.Once(); err != nil {
			continue
		}

		count++
		if count >= numOfNodes {
			return nil
		}
	}

	return fmt.Errorf("upload tx data to peers failed, txId: %s", txId)
}

func (c *Client) GetTxDataFromPeers(txId string) ([]byte, error) {
	peers, err := c.GetPeers()
	if err != nil {
		return nil, err
	}

	for _, peer := range peers {
		pNode := NewClient("http://" + peer)
		data, err := pNode.DownloadChunkData(txId)
		if err != nil {
			log.Error("get tx data", "err", err, "peer", peer)
			continue
		}
		return data, nil
	}

	return nil, errors.New("get tx data from peers failed")
}

func (c *Client) GetBlockFromPeers(height int64) (*types.Block, error) {
	peers, err := c.GetPeers()
	if err != nil {
		return nil, err
	}

	for _, peer := range peers {
		pNode := NewClient("http://" + peer)
		block, err := pNode.GetBlockByHeight(height)
		if err != nil {
			fmt.Printf("get block error:%v, peer: %s\n", err, peer)
			continue
		}
		fmt.Printf("success get block; peer: %s\n", peer)
		return block, nil
	}

	return nil, errors.New("get block from peers failed")
}

func (c *Client) GetTxFromPeers(arId string) (*types.Transaction, error) {
	peers, err := c.GetPeers()
	if err != nil {
		return nil, err
	}

	for _, peer := range peers {
		pNode := NewClient("http://" + peer)
		tx, err := pNode.GetTransactionByID(arId)
		if err != nil {
			fmt.Printf("get tx error:%v, peer: %s\n", err, peer)
			continue
		}
		fmt.Printf("success get tx; peer: %s\n", peer)
		return tx, nil
	}

	return nil, errors.New("get tx from peers failed")
}
