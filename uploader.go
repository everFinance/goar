package goar

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"

	"github.com/daqiancode/goar/types"
	"github.com/daqiancode/goar/utils"
	"github.com/shopspring/decimal"
)

type SerializedUploader struct {
	chunkIndex         int
	txPosted           bool
	transaction        *types.Transaction
	lastRequestTimeEnd int64
	lastResponseStatus int
	lastResponseError  string
}

type TransactionUploader struct {
	Client             *Client `json:"-"`
	ChunkIndex         int
	TxPosted           bool
	Transaction        *types.Transaction
	Data               io.ReadSeeker
	FileSize           int64
	LastRequestTimeEnd int64
	TotalErrors        int // Not serialized.
	LastResponseStatus int
	LastResponseError  string
	progressLock       sync.Mutex
	ProgressCallback   func(bytesSent int)
}

func newUploader(tt *types.Transaction, client *Client, file io.ReadSeeker, fileSize int64) (*TransactionUploader, error) {
	if tt.ID == "" {
		return nil, errors.New("Transaction is not signed.")
	}
	if tt.Chunks == nil {
		log.Warn("Transaction chunks not perpared")
	}
	// Make a copy of Transaction, zeroing the Data so we can serialize.
	tu := &TransactionUploader{
		Client:   client,
		Data:     file,
		FileSize: fileSize,
	}
	// da, err := utils.Base64Decode(tt.Data)
	// if err != nil {
	// 	log.Error("utils.Base64Decode(tt.Data)", "err", err)
	// 	return nil, err

	// }
	tu.Transaction = &types.Transaction{
		Format:   tt.Format,
		ID:       tt.ID,
		LastTx:   tt.LastTx,
		Owner:    tt.Owner,
		Tags:     tt.Tags,
		Target:   tt.Target,
		Quantity: tt.Quantity,
		// Data:      "",
		DataSize:  tt.DataSize,
		DataRoot:  tt.DataRoot,
		Reward:    tt.Reward,
		Signature: tt.Signature,
		Chunks:    tt.Chunks,
	}
	return tu, nil
}

// CreateUploader
// @param upload: Transaction | SerializedUploader | string,
// @param Data the Data of the Transaction. Required when resuming an upload.
func CreateUploader(api *Client, upload interface{}, data io.ReadSeeker, fileSize int64) (*TransactionUploader, error) {
	var (
		uploader *TransactionUploader
		err      error
	)

	if tt, ok := upload.(*types.Transaction); ok {
		uploader, err = newUploader(tt, api, data, fileSize)
		if err != nil {
			return nil, err
		}
		return uploader, nil
	}

	if id, ok := upload.(string); ok {
		// upload 返回为 SerializedUploader 类型
		upload, err = (&TransactionUploader{Client: api}).FromTransactionId(id)
		if err != nil {
			log.Error("(&TransactionUploader{Client: api}).FromTransactionId(id)", "err", err)
			return nil, err
		}
	} else {
		// 最后 upload 为 SerializedUploader type
		newUpload, ok := upload.(*SerializedUploader)
		if !ok {
			panic("upload params error")
		}
		upload = newUpload
	}

	uploader, err = (&TransactionUploader{Client: api}).FromSerialized(upload.(*SerializedUploader), data, fileSize)
	return uploader, err
}

func (tt *TransactionUploader) Once() (err error) {
	for !tt.IsComplete() {
		if err = tt.UploadChunk(); err != nil {
			return
		}

		if tt.LastResponseStatus != 200 {
			return errors.New(tt.LastResponseError)
		}
	}

	return
}

func (tt *TransactionUploader) IsComplete() bool {
	tChunks := tt.Transaction.Chunks
	if tChunks == nil {
		return false
	} else {
		return tt.TxPosted && (tt.ChunkIndex == len(tChunks.Chunks)) || tt.TxPosted && len(tChunks.Chunks) == 0
	}
}

func (tt *TransactionUploader) TotalChunks() int {
	if tt.Transaction.Chunks == nil {
		return 0
	} else {
		return len(tt.Transaction.Chunks.Chunks)
	}
}

func (tt *TransactionUploader) UploadedChunks() int {
	return tt.ChunkIndex
}

func (tt *TransactionUploader) PctComplete() float64 {
	val := decimal.NewFromInt(int64(tt.UploadedChunks())).Div(decimal.NewFromInt(int64(tt.TotalChunks())))
	fval, _ := val.Float64()
	return math.Trunc(fval * 100)
}

type ProgressFunc func(sentBytes int)

func (tt *TransactionUploader) ConcurrentOnce(ctx context.Context, concurrentNum int) error {
	// post tx info
	if err := tt.postTransaction(); err != nil {
		return err
	}

	if tt.IsComplete() {
		return nil
	}

	var wg sync.WaitGroup
	if concurrentNum <= 0 {
		concurrentNum = types.DEFAULT_CHUNK_CONCURRENT_NUM
	}
	p, _ := ants.NewPoolWithFunc(concurrentNum, func(i interface{}) {
		defer wg.Done()
		// process submit chunk
		idx := i.(int)

		select {
		case <-ctx.Done():
			log.Warn("ctx.done", "chunkIdx", idx)
			return
		default:
		}
		chunk, err := utils.GetChunk(*tt.Transaction, idx, tt.Data)
		if err != nil {
			log.Error("GetChunk error", "err", err, "idx", idx)
			return
		}
		body, statusCode, err := tt.Client.SubmitChunks(chunk) // always body is errMsg
		if statusCode == 200 {
			if tt.ProgressCallback != nil {
				tt.progressLock.Lock()
				tt.ProgressCallback(chunk.ChunkSize)
				tt.progressLock.Unlock()
			}

			return
		}

		log.Error("concurrent submitChunk failed", "chunkIdx", idx, "statusCode", statusCode, "gatewayErr", body, "httpErr", err)
		// try again
		retryCount := 0
		for {
			select {
			case <-ctx.Done():
				log.Warn("ctx.done", "chunkIdx", idx)
				return
			default:
			}

			retryCount++
			if statusCode == 429 {
				time.Sleep(1 * time.Second)
			} else {
				time.Sleep(200 * time.Millisecond)
			}

			body, statusCode, err = tt.Client.SubmitChunks(chunk)
			if statusCode == 200 {
				if tt.ProgressCallback != nil {
					tt.progressLock.Lock()
					tt.ProgressCallback(chunk.ChunkSize)
					tt.progressLock.Unlock()
				}
				return
			}
			log.Warn("retry submitChunk failed", "retryCount", retryCount, "chunkIdx", idx, "statusCode", statusCode, "gatewayErr", body, "httpErr", err)
		}
	})

	defer p.Release()
	for i := 0; i < len(tt.Transaction.Chunks.Chunks); i++ {
		wg.Add(1)
		if err := p.Invoke(i); err != nil {
			log.Error("p.Invoke(i)", "err", err, "i", i)
			return err
		}
	}

	wg.Wait()
	return nil
}

/**
 * Uploads the next part of the Transaction.
 * On the first call this posts the Transaction
 * itself and on any subsequent calls uploads the
 * next chunk until it completes.
 */
func (tt *TransactionUploader) UploadChunk() error {
	defer func() {
		// if tt.TotalChunks() > 0 {
		// 	log.Debug("chunks", "uploads", fmt.Sprintf("%f%% completes, %d/%d", tt.PctComplete(), tt.UploadedChunks(), tt.TotalChunks()))
		// }
	}()
	if tt.IsComplete() {
		return errors.New("Upload is already complete.")
	}

	if tt.LastResponseError != "" {
		tt.TotalErrors++
	} else {
		tt.TotalErrors = 0
	}

	// We have been trying for about an hour receiving an
	// error every time, so eventually bail.
	if tt.TotalErrors == 100 {
		return errors.New(fmt.Sprintf("Unable to complete upload: %d:%s", tt.LastResponseStatus, tt.LastResponseError))
	}

	var delay = 0.0
	if tt.LastResponseError != "" {
		delay = math.Max(float64(tt.LastRequestTimeEnd+types.ERROR_DELAY-time.Now().UnixNano()/1000000), float64(types.ERROR_DELAY))
	}
	if delay > 0.0 {
		// Jitter delay bcoz networks, subtract up to 30% from 40 seconds
		delay = delay - delay*0.3*rand.Float64()
		time.Sleep(time.Duration(delay) * time.Millisecond) // 休眠
	}

	tt.LastResponseError = ""

	if !tt.TxPosted {
		return tt.postTransaction()
	}

	chunk, err := utils.GetChunk(*tt.Transaction, tt.ChunkIndex, tt.Data)
	if err != nil {
		return err
	}
	path, err := utils.Base64Decode(chunk.DataPath)
	if err != nil {
		return err
	}
	offset, err := strconv.Atoi(chunk.Offset)
	if err != nil {
		return err
	}
	dataSize, err := strconv.Atoi(chunk.DataSize)
	if err != nil {
		return err
	}
	_, chunkOk := utils.ValidatePath(tt.Transaction.Chunks.DataRoot, offset, 0, dataSize, path)
	if !chunkOk {
		return errors.New(fmt.Sprintf("Unable to validate chunk %d ", tt.ChunkIndex))
	}
	// Catch network errors and turn them into objects with status -1 and an error message.
	gc, err := utils.GetChunk(*tt.Transaction, tt.ChunkIndex, tt.Data)
	if err != nil {
		return err
	}
	body, statusCode, err := tt.Client.SubmitChunks(gc) // always body is errMsg
	tt.LastRequestTimeEnd = time.Now().UnixNano() / 1000000
	tt.LastResponseStatus = statusCode
	if statusCode == 200 {
		tt.ChunkIndex++
		if tt.ProgressCallback != nil {
			tt.progressLock.Lock()
			tt.ProgressCallback(chunk.ChunkSize)
			tt.progressLock.Unlock()
		}
	} else {
		errStr := fmt.Sprintf("%s,%v,%d", body, err, statusCode)
		tt.LastResponseError = errStr
		if _, ok := types.FATAL_CHUNK_UPLOAD_ERRORS[body]; ok {
			return errors.New(fmt.Sprintf("Fatal error uploading chunk %d:%v", tt.ChunkIndex, body))
		}
	}
	return nil
}

/**
 * Reconstructs an upload from its serialized state and data.
 * Checks if data matches the expected data_root.
 *
 * @param serialized
 * @param data
 */
func (tt *TransactionUploader) FromSerialized(serialized *SerializedUploader, data io.ReadSeeker, fileSize int64) (*TransactionUploader, error) {
	if serialized == nil {
		return nil, errors.New("Serialized object does not match expected format.")
	}

	// Everything looks ok, reconstruct the TransactionUpload,
	// prepare the chunks again and verify the data_root matches
	upload, err := newUploader(serialized.transaction, tt.Client, data, fileSize)
	if err != nil {
		return nil, err
	}
	// Copy the serialized upload information, and Data passed in.
	upload.ChunkIndex = serialized.chunkIndex
	upload.LastRequestTimeEnd = serialized.lastRequestTimeEnd
	upload.LastResponseError = serialized.lastResponseError
	upload.LastResponseStatus = serialized.lastResponseStatus
	upload.TxPosted = serialized.txPosted
	upload.Data = data
	upload.FileSize = fileSize

	utils.PrepareChunks(upload.Transaction)

	if upload.Transaction.DataRoot != serialized.transaction.DataRoot {
		return nil, errors.New("Data mismatch: Uploader doesn't match provided Data.")
	}

	return upload, nil
}

/**
 * Reconstruct an upload from the tx metadata, ie /tx/<id>.
 *
 * @param api
 * @param id
 * @param data
 */
func (tt *TransactionUploader) FromTransactionId(id string) (*SerializedUploader, error) {
	tx, err := tt.Client.GetTransactionByID(id)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Tx %s not found; error: %v", id, err))
	}
	transaction := tx

	serialized := &SerializedUploader{
		chunkIndex:         0,
		txPosted:           true,
		transaction:        transaction,
		lastRequestTimeEnd: 0,
		lastResponseStatus: 0,
		lastResponseError:  "",
	}
	return serialized, nil
}

func (tt *TransactionUploader) FormatSerializedUploader() *SerializedUploader {
	tx := tt.Transaction
	return &SerializedUploader{
		chunkIndex:         tt.ChunkIndex,
		txPosted:           tt.TxPosted,
		transaction:        tx,
		lastRequestTimeEnd: tt.LastRequestTimeEnd,
		lastResponseStatus: tt.LastResponseStatus,
		lastResponseError:  tt.LastResponseError,
	}
}

// POST to /tx
func (tt *TransactionUploader) postTransaction() error {
	var uploadInBody = tt.TotalChunks() <= types.MAX_CHUNKS_IN_BODY
	return tt.uploadTx(uploadInBody)
}

func (tt *TransactionUploader) uploadTx(withBody bool) error {
	if withBody {
		// Post the Transaction with Data.
		_, err := tt.Data.Seek(0, 0)
		if err != nil {
			return err
		}
		tt.Transaction.Data, err = utils.Base64EncodeReader(tt.Data, tt.FileSize)
		if err != nil {
			return err
		}
	}
	body, statusCode, err := tt.Client.SubmitTransaction(tt.Transaction)
	if err != nil || statusCode >= 400 {
		tt.LastResponseError = fmt.Sprintf("%v,%s", err, body)
		tt.LastResponseStatus = statusCode
		return errors.New(fmt.Sprintf("Unable to upload Transaction: %d, %v, %s", statusCode, err, body))
	}

	tt.LastRequestTimeEnd = time.Now().UnixNano() / 1000000
	tt.LastResponseStatus = statusCode

	if withBody {
		tt.Transaction.Data = ""
	}

	// tx already processed
	if statusCode >= 200 && statusCode < 300 {
		tt.TxPosted = true
		if withBody {
			// We are complete.
			tt.ChunkIndex = types.MAX_CHUNKS_IN_BODY
		}
		return nil
	}

	if withBody {
		tt.LastResponseError = ""
	}
	return nil
}
