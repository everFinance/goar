package utils

import (
	"encoding/base64"
	"io"
)

func Base64Encode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func Base64Decode(data string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(data)
}

func Base64EncodeReader(data io.Reader, fileSize int64) (string, error) {
	buffer := make([]byte, fileSize)
	_, err := data.Read(buffer)
	if err != nil {
		return "", err
	}
	return Base64Encode(buffer), nil
}
