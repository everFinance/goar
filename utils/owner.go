package utils

import "crypto/sha256"

func OwnerToAddress(owner string) (address string) {
	by, _ := Base64Decode(owner)
	addr := sha256.Sum256(by)
	address = Base64Encode(addr[:])
	return
}
