package crypto

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"reflect"
)


func SHA256(data []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	r := h.Sum(nil)
	return r, nil
}



func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func unpackArray(s any) []any {
	v := reflect.ValueOf(s)
	r := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		r[i] = v.Index(i).Interface()
	}
	return r
}

func DeepHash(data any) [48]byte {
	if typeof(data) == "[]uint8" {
		tag := append([]byte("blob"), []byte(fmt.Sprintf("%d", len(data.([]byte))))...)
		tagHashed := sha512.Sum384(tag)
		dataHashed := sha512.Sum384(data.([]byte))
		r := append(tagHashed[:], dataHashed[:]...)
		rHashed := sha512.Sum384(r)
		return rHashed
	} else {
		_data := unpackArray(data)
		tag := append([]byte("list"), []byte(fmt.Sprintf("%d", len(_data)))...)
		tagHashed := sha512.Sum384(tag)
		return deepHashChunk(_data, tagHashed)
	}
}
func deepHashChunk(data []any, acc [48]byte) [48]byte {
	if len(data) < 1 {
		return acc
	}
	dHash := DeepHash(data[0])
	hashPair := append(acc[:], dHash[:]...)
	newAcc := sha512.Sum384(hashPair)
	return deepHashChunk(data[1:], newAcc)
}
