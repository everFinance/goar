package data

import "github.com/everFinance/goar/utils"

func verify(data []byte, signature []byte, owner string) (bool, error) {
	publicKey, err := utils.OwnerToPubKey(owner)
	if err != nil {
		return false, err
	}
	err = utils.Verify(data, publicKey, signature)
	if err != nil {
		return false, err
	}
	return true, err
}
