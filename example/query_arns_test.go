package example

import (
	"github.com/everFinance/goar/arns"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_QueryArNS(t *testing.T) {
	arID, err := queryArNS()
	assert.NoError(t, err)
	t.Log(arID)
}

func queryArNS() (arID string, err error) {
	dreUrl := "https://dre-3.warp.cc"
	arNSAddress := "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U"
	timeout := 10 * time.Second

	domain := "cookbook_ao"

	a := arns.NewArNS(dreUrl, arNSAddress, timeout)
	arID, err = a.QueryLatestRecord(domain)
	return
}
