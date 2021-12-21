package goar

import (
	"errors"
	"fmt"
	"github.com/everFinance/goar/types"
)

func (c *Client) BroadcastData(txId string, data []byte, numOfNodes int64, peers ...string) error {
	var err error
	if len(peers) == 0 {
		peers, err = c.GetPeers()
		if err != nil {
			return err
		}
	}

	count := int64(0)
	pNode := NewTempConn()
	for _, peer := range peers {
		pNode.SetShortConnUrl("http://" + peer)
		uploader, err := CreateUploader(pNode, txId, data)
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

func (c *Client) GetTxDataFromPeers(txId string, peers ...string) ([]byte, error) {
	var err error
	if len(peers) == 0 {
		peers, err = c.GetPeers()
		if err != nil {
			return nil, err
		}
	}

	pNode := NewTempConn()
	for _, peer := range peers {
		pNode.SetShortConnUrl("http://" + peer)
		data, err := pNode.DownloadChunkData(txId)
		if err != nil {
			log.Error("get tx data", "err", err, "peer", peer)
			continue
		}
		return data, nil
	}

	return nil, errors.New("get tx data from peers failed")
}

func (c *Client) GetBlockFromPeers(height int64, peers ...string) (*types.Block, error) {
	var err error
	if len(peers) == 0 {
		peers, err = c.GetPeers()
		if err != nil {
			return nil, err
		}
	}

	pNode := NewTempConn()
	for _, peer := range peers {
		pNode.SetShortConnUrl("http://" + peer)
		block, err := pNode.GetBlockByHeight(height)
		if err != nil {
			fmt.Printf("get block error:%v, peer: %s, height: %d\n", err, peer, height)
			continue
		}
		fmt.Printf("success get block; peer: %s\n", peer)
		return block, nil
	}

	return nil, errors.New("get block from peers failed")
}

func (c *Client) GetTxFromPeers(arId string, peers ...string) (*types.Transaction, error) {
	var err error
	if len(peers) == 0 {
		peers, err = c.GetPeers()
		if err != nil {
			return nil, err
		}
	}

	pNode := NewTempConn()
	for _, peer := range peers {
		pNode.SetShortConnUrl("http://" + peer)
		tx, err := pNode.GetTransactionByID(arId)
		if err != nil {
			fmt.Printf("get tx error:%v, peer: %s, arTx: %s\n", err, peer, arId)
			continue
		}
		fmt.Printf("success get tx; peer: %s, arTx: %s\n", peer, arId)
		return tx, nil
	}

	return nil, fmt.Errorf("get tx failed; arId: %s", arId)
}

func (c *Client) GetUnconfirmedTxFromPeers(arId string, peers ...string) (*types.Transaction, error) {
	var err error
	if len(peers) == 0 {
		peers, err = c.GetPeers()
		if err != nil {
			return nil, err
		}
	}

	pNode := NewTempConn()
	for _, peer := range peers {
		pNode.SetShortConnUrl("http://" + peer)
		tx, err := pNode.GetUnconfirmedTx(arId)
		if err != nil {
			fmt.Printf("get tx error:%v, peer: %s, arTx: %s\n", err, peer, arId)
			continue
		}
		fmt.Printf("success get unconfirmed tx; peer: %s, arTx: %s\n", peer, arId)
		return tx, nil
	}

	return nil, fmt.Errorf("get unconfirmed tx failed; arId: %s", arId)
}
