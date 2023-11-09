package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestRSASignAndVerify(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("data to sign")

	signer := NewRS256Signature(privateKey)
	signature, err := signer.Sign(data)
	if err != nil {
		t.Fatal(err)
	}

	// print signature length and public key length
	t.Log("signature length:", len(signature))
	t.Log("public key length:", len(signer.privateKey.PublicKey.N.Bytes()))

	// public key bytes
	pub := signer.privateKey.PublicKey.N.Bytes()
	err = RS256Verify(pub, data, signature)
	if err != nil {
		t.Fatal(err)
	}

	// Test invalid signature
	err = RS256Verify(pub, data, []byte("invalid signature"))
	if err == nil {
		t.Fatal("expected error for invalid signature")
	}
}
