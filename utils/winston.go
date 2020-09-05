package utils

import (
	"math/big"
)

// WinstonToAR 1 Winston = 0.000000000001 AR
func WinstonToAR(w *big.Int) *big.Float {
	return new(big.Float).Quo(
		new(big.Float).SetInt(w),
		new(big.Float).SetUint64(1000000000000),
	)
}

// ARToWinston 1 AR = 1000000000000
func ARToWinston(a *big.Float) *big.Int {
	wFloat := new(big.Float).Mul(a, new(big.Float).SetUint64(1000000000000))
	w, _ := new(big.Int).SetString(wFloat.Text('f', 0), 10)
	return w
}
