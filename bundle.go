/*
	js-lib:  https://github.com/Bundler-Network/arbundles
	ANS-104 format: https://github.com/joshbenaron/arweave-standards/blob/ans104/ans/ANS-104.md
*/

package goar

import (
	"errors"
	"fmt"
	"github.com/everFinance/goar/utils"
)

type BundleData struct {
	Items        []DataItem `json:"items"`
	bundleBinary []byte
}

func NewBundleData(dataItems ...DataItem) (*BundleData, error) {
	headers := make([]byte, 0) // length is 64 * len(dataItems)
	binaries := make([]byte, 0)

	for _, d := range dataItems {
		header := make([]byte, 0, 64)
		header = append(header, utils.LongTo32ByteArray(len(d.itemBinary))...)
		id, err := utils.Base64Decode(d.Id)
		if err != nil {
			return nil, err
		}
		header = append(header, id...)

		headers = append(headers, header...)
		binaries = append(binaries, d.itemBinary...)
	}

	bdBinary := make([]byte, 0)
	bdBinary = append(bdBinary, utils.LongTo32ByteArray(len(dataItems))...)
	bdBinary = append(bdBinary, headers...)
	bdBinary = append(bdBinary, binaries...)
	return &BundleData{
		Items:        dataItems,
		bundleBinary: bdBinary,
	}, nil
}

func RecoverBundleData(bundleBinary []byte) (*BundleData, error) {
	// length must more than 32
	if len(bundleBinary) < 32 {
		return nil, errors.New("binary length must more than 32")
	}
	dataItemsNum := utils.ByteArrayToLong(bundleBinary[:32])

	if len(bundleBinary) < 32+dataItemsNum*64 {
		return nil, errors.New("binary length incorrect")
	}

	bd := &BundleData{
		Items:        make([]DataItem, 0),
		bundleBinary: bundleBinary,
	}
	dataItemStart := 32 + dataItemsNum*64
	for i := 0; i < dataItemsNum; i++ {
		headerBegin := 32 + i*64
		end := headerBegin + 64
		headerByte := bundleBinary[headerBegin:end]
		itemBinaryLength := utils.ByteArrayToLong(headerByte[:32])
		id := utils.Base64Encode(headerByte[32:64])

		dataItemBytes := bundleBinary[dataItemStart : dataItemStart+itemBinaryLength]
		dataItem, err := recoverDataItem(dataItemBytes)
		if err != nil {
			return nil, err
		}
		if dataItem.Id != id {
			return nil, fmt.Errorf("dataItem.Id != id, dataItem.Id: %s, id: %s", dataItem.Id, id)
		}
		bd.Items = append(bd.Items, *dataItem)
		dataItemStart += itemBinaryLength
	}
	return bd, nil
}
