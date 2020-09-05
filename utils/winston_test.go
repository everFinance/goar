package utils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWinstonToAR(t *testing.T) {
	w := new(big.Int).SetInt64(1000000000000)
	a := WinstonToAR(w)
	assert.Equal(t, "1", a.String())

	w = new(big.Int).SetInt64(1)
	a = WinstonToAR(w)
	assert.Equal(t, "1e-12", a.String())
}

func TestARToWinston(t *testing.T) {
	a := new(big.Float).SetFloat64(1.000001)
	w := ARToWinston(a)
	assert.Equal(t, "1000001000000", w.String())

	a = new(big.Float).SetFloat64(1e-12)
	w = ARToWinston(a)
	assert.Equal(t, "1", w.String())
}
