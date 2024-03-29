package utils

import (
	"crypto/sha512"
	"fmt"
	"io"
	"reflect"
)

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

func deepHashStream(data io.Reader) [48]byte {
	hash := sha512.New384()
	n, err := io.Copy(hash, data)
	if err != nil {
		panic(err)
	}
	tag := append([]byte("blob"), []byte(fmt.Sprintf("%d", n))...)
	tagHash := sha512.Sum384(tag)
	blobHash := hash.Sum(nil)
	tagged := append(tagHash[:], blobHash...)

	return sha512.Sum384(tagged)
}

func deepHashChunk(data []interface{}, acc [48]byte) [48]byte {
	if len(data) < 1 {
		return acc
	}

	dHash := [48]byte{}
	if _, ok := data[0].(io.Reader); ok {
		dHash = deepHashStream(data[0].(io.Reader))
	} else if reflect.TypeOf(data[0]).String() == "string" {
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
