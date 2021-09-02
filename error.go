package goar

import "errors"

var (
	ErrNotFound   = errors.New("Not Found")
	ErrPendingTx  = errors.New("Pending")
	ErrInvalidId  = errors.New("Invalid ar tx id")
	ErrBadGateway = errors.New("Bad gateway")
)
