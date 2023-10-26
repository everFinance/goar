package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
)

func signWithRS256(data string, privateKey string) (string, error) {
	// 解码 PEM 编码的私钥
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", errors.New("failed to decode PEM block containing private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 计算散列值
	hasher := sha256.New()
	hasher.Write([]byte(data))
	hashed := hasher.Sum(nil)

	// 创建签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", signature), nil

}

func verifyWithRS256(data string, signatureHex string, publicKey string) (bool, error) {
	// 解码 PEM 编码的公钥
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		return false, errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("not an RSA public key")
	}

	// 将十六进制签名转换回字节
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, err
	}

	// 计算散列值
	hasher := sha256.New()
	hasher.Write([]byte(data))
	hashed := hasher.Sum(nil)

	// 验证签名
	err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hashed, signature)
	if err != nil {
		return false, err
	}

	return true, nil
}
