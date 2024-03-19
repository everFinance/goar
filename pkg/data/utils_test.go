package data

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBundleHeader(t *testing.T) {
	data, err := os.ReadFile("../test/stubs/bundleHeader")
	if err != nil {
		log.Fatal(err)
	}
	headers, N := decodeBundleHeader(&data)
	assert.Equal(t, N, 1)
	assert.Equal(t, (*headers)[0].size, 1115)
	assert.Equal(t, (*headers)[0].id, 39234)
}

func TestGenerateBundleHeader(t *testing.T) {
	data, err := os.ReadFile("../test/stubs/1115BDataItem")
	assert.NoError(t, err)

	dataItem, err := DecodeDataItem(data)
	assert.NoError(t, err)
	headers, err := generateBundleHeader(&[]DataItem{*dataItem})

	assert.NoError(t, err)
	assert.Equal(t, (*headers)[0].size, 1115)
	assert.Equal(t, (*headers)[0].id, 39234)
}

func TestEncodeAvro(t *testing.T) {
	data := []byte{6, 24, 67, 111, 110, 116, 101, 110, 116, 45, 84, 121, 112, 101, 20, 116, 101, 120, 116, 47, 112, 108, 97, 105, 110, 16, 65, 112, 112, 45, 78, 97, 109, 101, 22, 65, 114, 68, 114, 105, 118, 101, 45, 67, 76, 73, 22, 65, 112, 112, 45, 86, 101, 114, 115, 105, 111, 110, 12, 49, 46, 50, 49, 46, 48, 0}
	tags := []Tag{
		{Name: "Content-Type", Value: "text/plain"},
		{Name: "App-Name", Value: "ArDrive-CLI"},
		{Name: "App-Version", Value: "1.21.0"},
	}

	rawTags, err := encodeAvro(&tags)
	assert.NoError(t, err)
	assert.ElementsMatch(t, data, rawTags)
}

func TestDecodeAvro(t *testing.T) {
	data := []byte{6, 24, 67, 111, 110, 116, 101, 110, 116, 45, 84, 121, 112, 101, 20, 116, 101, 120, 116, 47, 112, 108, 97, 105, 110, 16, 65, 112, 112, 45, 78, 97, 109, 101, 22, 65, 114, 68, 114, 105, 118, 101, 45, 67, 76, 73, 22, 65, 112, 112, 45, 86, 101, 114, 115, 105, 111, 110, 12, 49, 46, 50, 49, 46, 48, 0}
	tags, err := decodeAvro(data)
	assert.NoError(t, err)
	assert.ElementsMatch(t, *tags, []Tag{
		{Name: "Content-Type", Value: "text/plain"},
		{Name: "App-Name", Value: "ArDrive-CLI"},
		{Name: "App-Version", Value: "1.21.0"},
	})
}
