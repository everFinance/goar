package data

import (
	"encoding/base64"
	"errors"
)

func DecodeBundle(data []byte) (*Bundle, error) {
	// length must more than 32
	if len(data) < 32 {
		return nil, errors.New("binary length must more than 32")
	}
	headers, N := decodeBundleHeader(&data)
	bundle := &Bundle{
		Items:   make([]DataItem, N),
		RawData: base64.RawURLEncoding.EncodeToString(data),
	}
	bundleStart := 32 + 64*N
	for i := 0; i < N; i++ {
		header := (*headers)[i]
		bundleEnd := bundleStart + header.size
		dataItem, err := DecodeDataItem(data[bundleStart:bundleEnd])
		if err != nil {
			return nil, err
		}
		bundle.Items[i] = *dataItem
		bundleStart = bundleEnd
	}
	return bundle, nil
}

func NewBundle(dataItems *[]DataItem) (*Bundle, error) {
	bundle := &Bundle{}

	headers, err := generateBundleHeader(dataItems)
	if err != nil {
		return nil, err
	}

	bundle.Headers = *headers
	bundle.Items = *dataItems
	N := len(*dataItems)

	var sizeBytes []byte
	var headersBytes []byte
	var dataItemsBytes []byte

	for i := 0; i < N; i++ {
		headersBytes = append(headersBytes, (*headers)[i].raw...)
		dataItemsBytes = append(dataItemsBytes, (*headers)[i].raw...)
	}

	bundle.RawData = base64.RawURLEncoding.EncodeToString(append(sizeBytes, append(headersBytes, dataItemsBytes...)...))
	return bundle, nil
}

func ValidateBundle(data []byte) (bool, error) {
	// length must more than 32
	if len(data) < 32 {
		return false, errors.New("binary length must more than 32")
	}
	headers, N := decodeBundleHeader(&data)
	dataItemSize := 0
	for i := 0; i < N; i++ {
		dataItemSize += (*headers)[i].size
	}
	return len(data) == dataItemSize+32+64*N, nil
}
