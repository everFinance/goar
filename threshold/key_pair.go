package threshold

import (
	"errors"
	"fmt"
	"github.com/everFinance/sandy_log/log"
	tcrsa "github.com/everFinance/ttcrsa"
	"time"
)

// CreateKeyPair
// bitSize: Is used to generate key shares with a security level equivalent to a RSA private of that size.
// l: creates l key shares for a k-threshold signing scheme.
// k: The generated key shares have a threshold parameter of k
func CreateKeyPair(bitSize, k, l int) (shares tcrsa.KeyShareList, meta *tcrsa.KeyMeta, err error) {
	if bitSize > 4096 || bitSize < 512 {
		return nil, nil, errors.New(fmt.Sprintf("bitSize:%d parameter must in [512,4096]", bitSize))
	}
	if k <= 0 || l <= 1 || l < k || k < (l/2+1) {
		return nil, nil, errors.New(fmt.Sprintf("k: %d l: %d parameter incorrect; k must > 0, l must > 1, l must >= k, k must >= (l/2+1)", k, l))
	}

	now := time.Now()
	keyShares, keyMeta, err := tcrsa.NewKey(bitSize, uint16(k), uint16(l), nil)
	if err != nil {
		log.Errorf("tcrsa newKey error; bitSize: %d, k: %d, l: %d; err: %v", bitSize, k, l, err)
		return nil, nil, err
	}
	log.Debugf("Create bit size = %d rsa threshold keyPair spend time: %s", bitSize, time.Since(now).String())
	return keyShares, keyMeta, nil
}
