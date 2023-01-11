package goar

import "errors"

var (
	ErrNotFound     = errors.New("Not Found")
	ErrPendingTx    = errors.New("Pending")
	ErrInvalidId    = errors.New("Invalid ArId")
	ErrBadGateway   = errors.New("Bad Gateway")
	ErrRequestLimit = errors.New("Arweave gateway request limit")
)
