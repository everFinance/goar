package goar

import (
	"context"
	"errors"
	"os"

	"github.com/everVision/goar/schema"
)

func (w *Wallet) SendBundleTxSpeedUp(ctx context.Context, concurrentNum int, bundleBinary interface{}, tags []schema.Tag, txSpeed int64) (schema.Transaction, error) {
	bundleTags := []schema.Tag{
		{Name: "Bundle-Format", Value: "binary"},
		{Name: "Bundle-Version", Value: "2.0.0"},
	}
	// check tags cannot include bundleTags Name
	mmap := map[string]struct{}{
		"Bundle-Format":  {},
		"Bundle-Version": {},
	}
	for _, tag := range tags {
		if _, ok := mmap[tag.Name]; ok {
			return schema.Transaction{}, errors.New("tags can not set bundleTags")
		}
	}

	txTags := make([]schema.Tag, 0)
	txTags = append(bundleTags, tags...)
	return w.SendDataConcurrentSpeedUp(ctx, concurrentNum, bundleBinary, txTags, txSpeed)
}

func (w *Wallet) SendBundleTx(ctx context.Context, concurrentNum int, bundleBinary []byte, tags []schema.Tag) (schema.Transaction, error) {
	return w.SendBundleTxSpeedUp(ctx, concurrentNum, bundleBinary, tags, 0)
}

func (w *Wallet) SendBundleTxStream(ctx context.Context, concurrentNum int, bundleReader *os.File, tags []schema.Tag) (schema.Transaction, error) {
	return w.SendBundleTxSpeedUp(ctx, concurrentNum, bundleReader, tags, 0)
}
