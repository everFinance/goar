package utils_test

import (
	"fmt"
	"testing"

	"github.com/daqiancode/goar/utils"
	"github.com/stretchr/testify/assert"
)

func TestReadBuffer(t *testing.T) {
	buf := utils.NewReadBuffer([]byte("this is a read buffer"))
	r := make([]byte, 50)
	for {
		n, err := buf.Read(r)
		assert.Nil(t, err)
		if n == 0 {
			break
		}
		fmt.Println(string(r))
	}
}
