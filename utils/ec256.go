package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"math/big"
)

// ES256签名
func SignES256(privateKey *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	// Step 1: 对消息内容进行SHA-256哈希
	hash := sha256.Sum256(message)

	// Step 2: 使用ECDSA算法进行签名
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, err
	}

	// 将签名结果序列化为DER编码格式
	signature, err := asn1.Marshal(struct {
		R, S *big.Int
	}{r, s})

	return signature, err
}

// 验证ES256签名
func VerifyES256(publicKey *ecdsa.PublicKey, message, signature []byte) (bool, error) {
	// 解析签名
	sigStruct := struct {
		R, S *big.Int
	}{}
	if _, err := asn1.Unmarshal(signature, &sigStruct); err != nil {
		return false, err
	}

	// 对消息内容进行SHA-256哈希
	hash := sha256.Sum256(message)

	// 验证签名
	valid := ecdsa.Verify(publicKey, hash[:], sigStruct.R, sigStruct.S)
	return valid, nil
}
