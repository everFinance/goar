package goar

import (
	"errors"
	"fmt"
)

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
