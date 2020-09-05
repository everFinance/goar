package utils

import (
	"crypto/sha512"
	"fmt"
	"reflect"

	"github.com/everFinance/goar/types"
)

func DataHash(tx *types.Transaction) (deepHash []byte, err error) {
	if err = PrepareChunks(tx); err != nil {
		return
	}

	tags := [][]string{}
	for _, tag := range tx.Tags {
		tags = append(tags, []string{
			tag.Name, tag.Value,
		})
	}

	dataList := []interface{}{}
	dataList = append(dataList, Base64Encode([]byte(fmt.Sprintf("%d", tx.Format))))
	dataList = append(dataList, tx.Owner)
	dataList = append(dataList, tx.Target)
	dataList = append(dataList, Base64Encode([]byte(tx.Quantity)))
	dataList = append(dataList, Base64Encode([]byte(tx.Reward)))
	dataList = append(dataList, tx.LastTx)
	dataList = append(dataList, tags)
	dataList = append(dataList, Base64Encode([]byte(tx.DataSize)))
	dataList = append(dataList, tx.DataRoot)

	hash := DeepHash(dataList)
	deepHash = hash[:]
	return
}

func PrepareChunks(tx *types.Transaction) (err error) {
	if tx.Data == "" {
		return
	}

	data, _ := Base64Decode(tx.Data)
	chunks := GenerateChunks(data)
	tx.DataRoot = Base64Encode(chunks.DataRoot)

	// TODO, use chunks in tx

	return
}

func DeepHash(data []interface{}) [48]byte {
	tag := append([]byte("list"), []byte(fmt.Sprintf("%d", len(data)))...)
	tagHash := sha512.Sum384(tag)

	return deepHashChunk(data, tagHash)
}

func deepHashStr(str string) [48]byte {
	by, _ := Base64Decode(str)
	tag := append([]byte("blob"), []byte(fmt.Sprintf("%d", len(by)))...)
	tagHash := sha512.Sum384(tag)
	blobHash := sha512.Sum384(by)
	tagged := append(tagHash[:], blobHash[:]...)

	return sha512.Sum384(tagged)
}

func deepHashChunk(data []interface{}, acc [48]byte) [48]byte {
	if len(data) < 1 {
		return acc
	}

	dHash := [48]byte{}
	if reflect.TypeOf(data[0]).String() == "string" {
		dHash = deepHashStr(data[0].(string))
	} else {
		value := reflect.ValueOf(data[0])
		dData := []interface{}{}

		for i := 0; i < value.Len(); i++ {
			dData = append(dData, value.Index(i).Interface())
		}

		dHash = DeepHash(dData)
	}

	hashPair := append(acc[:], dHash[:]...)
	newAcc := sha512.Sum384(hashPair)
	return deepHashChunk(data[1:], newAcc)
}
