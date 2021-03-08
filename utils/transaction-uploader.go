package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/everFinance/goar/client"
	"math"
	"math/rand"
	"strconv"
	"time"
)

// Maximum amount of chunks we will upload in the body.
const MAX_CHUNKS_IN_BODY = 1

// We assume these errors are intermitment and we can try again after a delay:
// - not_joined
// - timeout
// - data_root_not_found (we may have hit a node that just hasn't seen it yet)
// - exceeds_disk_pool_size_limit
// We also try again after any kind of unexpected network errors

// Errors from /chunk we should never try and continue on.
var FATAL_CHUNK_UPLOAD_ERRORS = map[string]struct{}{
	"invalid_json":                     struct{}{},
	"chunk_too_big":                    struct{}{},
	"data_path_too_big":                struct{}{},
	"offset_too_big":                   struct{}{},
	"data_size_too_big":                struct{}{},
	"chunk_proof_ratio_not_attractive": struct{}{},
	"invalid_proof":                    struct{}{},
}

// Amount we will delay on receiving an error response but do want to continue.
const ERROR_DELAY = 1000 * 40

type SerializedUploader struct {
	chunkIndex         int
	txPosted           bool
	transaction        *Transaction
	lastRequestTimeEnd int64
	lastResponseStatus int
	lastResponseError  string
}

type TransactionUploader struct {
	Client             *client.Client
	chunkIndex         int
	txPosted           bool
	transaction        *Transaction
	data               []byte
	lastRequestTimeEnd int64
	TotalErrors        int // Not serialized.
	LastResponseStatus int
	LastResponseError  string
}

func NewTransactionUploader(tt *Transaction, client *client.Client) (*TransactionUploader, error) {
	if tt.ID == "" {
		return nil, errors.New("Transaction is not signed.")
	}
	if tt.Chunks == nil {
		return nil, errors.New("Transaction chunks not perpared.")
	}
	// Make a copy of transaction, zeroing the data so we can serialize.
	tu := &TransactionUploader{
		Client: client,
	}
	tu.data = tt.Data
	tu.transaction = &Transaction{
		Format:    tt.Format,
		ID:        tt.ID,
		LastTx:    tt.LastTx,
		Owner:     tt.Owner,
		Tags:      tt.Tags,
		Target:    tt.Target,
		Quantity:  tt.Quantity,
		Data:      make([]byte, 0),
		DataSize:  tt.DataSize,
		DataRoot:  tt.DataRoot,
		Reward:    tt.Reward,
		Signature: tt.Signature,
		Chunks:    tt.Chunks,
	}
	return tu, nil
}

func (tt *TransactionUploader) IsComplete() bool {
	tChunks := tt.transaction.Chunks
	if tChunks == nil {
		return false
	} else {
		return tt.txPosted && (tt.chunkIndex == len(tChunks.Chunks))
	}
}

func (tt *TransactionUploader) TotalChunks() int {
	if tt.transaction.Chunks == nil {
		return 0
	} else {
		return len(tt.transaction.Chunks.Chunks)
	}
}

func (tt *TransactionUploader) UploadChunks() int {
	return tt.chunkIndex
}

func (tt *TransactionUploader) PctComplete() float64 {
	return math.Trunc(float64(tt.UploadChunks()/tt.TotalChunks()) * 100)
}

/**
 * Uploads the next part of the transaction.
 * On the first call this posts the transaction
 * itself and on any subsequent calls uploads the
 * next chunk until it completes.
 */
func (tt *TransactionUploader) UploadChunk() error {
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
		delay = math.Max(float64(tt.lastRequestTimeEnd+ERROR_DELAY-time.Now().UnixNano()/1000000), float64(ERROR_DELAY))
	}
	if delay > 0.0 {
		// Jitter delay bcoz networks, subtract up to 30% from 40 seconds
		delay = delay - delay*0.3*rand.Float64()
		time.Sleep(time.Duration(delay) * time.Millisecond) // 休眠
	}

	tt.LastResponseError = ""

	if !tt.txPosted {
		return tt.postTransaction()
	}

	chunk, err := tt.transaction.GetChunk(tt.chunkIndex, tt.data)
	if err != nil {
		return err
	}
	path, err := Base64Decode(chunk.DataPath)
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
	chunkOk := ValidatePath(tt.transaction.Chunks.DataRoot,
		offset, 0, dataSize, path)

	if !chunkOk {
		return errors.New(fmt.Sprintf("Unable to validate chunk %d", tt.chunkIndex))
	}
	// Catch network errors and turn them into objects with status -1 and an error message.
	gc, err := tt.transaction.GetChunk(tt.chunkIndex, tt.data)
	if err != nil {
		return err
	}
	byteGc, err := gc.Marshal()
	if err != nil {
		return err
	}
	_, statusCode, err := tt.Client.HttpPost("chunk", byteGc)
	tt.lastRequestTimeEnd = time.Now().UnixNano() / 1000000
	tt.LastResponseStatus = statusCode
	if statusCode == 200 {
		tt.chunkIndex++
	} else if err != nil {
		tt.LastResponseError = err.Error()
		if _, ok := FATAL_CHUNK_UPLOAD_ERRORS[err.Error()]; ok {
			return errors.New(fmt.Sprintf("Fatal error uploading chunk %d:%v", tt.chunkIndex, err))
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
func (tt *TransactionUploader) FromSerialized(serialized *SerializedUploader, data []byte) (*TransactionUploader, error) {
	if serialized == nil {
		return nil, errors.New("Serialized object does not match expected format.")
	}

	// Everything looks ok, reconstruct the TransactionUpload,
	// prepare the chunks again and verify the data_root matches
	upload, err := NewTransactionUploader(serialized.transaction, tt.Client)
	if err != nil {
		return nil, err
	}
	// Copy the serialized upload information, and data passed in.
	upload.chunkIndex = serialized.chunkIndex
	upload.lastRequestTimeEnd = serialized.lastRequestTimeEnd
	upload.LastResponseError = serialized.lastResponseError
	upload.LastResponseStatus = serialized.lastResponseStatus
	upload.txPosted = serialized.txPosted
	upload.data = data

	upload.transaction.PrepareChunks(data)

	if upload.transaction.DataRoot != serialized.transaction.DataRoot {
		return nil, errors.New("Data mismatch: Uploader doesn't match provided data.")
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
	body, statusCode, err := tt.Client.HttpGet(fmt.Sprintf("tx/%s", id))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Tx %s not found: %d", id, statusCode))
	}
	transaction := &Transaction{}
	if err := json.Unmarshal(body, transaction); err != nil {
		return nil, err
	}

	transaction.Data = make([]byte, 0)

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

// POST to /tx
func (tt *TransactionUploader) postTransaction() error {
	var uploadInBody = tt.TotalChunks() <= MAX_CHUNKS_IN_BODY

	byteTx, err := json.Marshal(&tt.transaction)
	if err != nil {
		return err
	}

	if uploadInBody {
		// Post the transaction with data.
		tt.transaction.Data = tt.data

		_, status, err := tt.Client.HttpPost("tx", byteTx)
		if err != nil {
			fmt.Printf("tt.Client.SubmitTransaction(&tt.transaction) error: %v", err)
			return err
		}
		tt.lastRequestTimeEnd = time.Now().UnixNano() / 1000000
		tt.LastResponseStatus = status
		tt.transaction.Data = make([]byte, 0)

		if status >= 200 && status < 300 {
			// We are complete.
			tt.txPosted = true
			tt.chunkIndex = MAX_CHUNKS_IN_BODY
			return nil
		}
		tt.LastResponseError = ""
		return nil
	}

	// else
	// Post the transaction with no data.
	_, status, err := tt.Client.HttpPost("tx", byteTx)
	tt.lastRequestTimeEnd = time.Now().UnixNano() / 1000000
	tt.LastResponseStatus = status
	if !(status >= 200 && status < 300) {
		if err != nil {
			tt.LastResponseError = err.Error()
		}
		return errors.New(fmt.Sprintf("Unable to upload transaction: %d, %v", status, err))
	}
	tt.txPosted = true
	return nil
}
