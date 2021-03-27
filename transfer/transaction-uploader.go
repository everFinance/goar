package transfer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/everFinance/goar/client"
	"github.com/everFinance/goar/merkle"
	"github.com/everFinance/goar/utils"
	"github.com/shopspring/decimal"
	"github.com/zyjblockchain/sandy_log/log"
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
	transaction        *TransactionChunks
	lastRequestTimeEnd int64
	lastResponseStatus int
	lastResponseError  string
}

type TransactionUploader struct {
	Client             *client.Client `json:"-"`
	ChunkIndex         int
	TxPosted           bool
	Transaction        *TransactionChunks
	Data               []byte
	LastRequestTimeEnd int64
	TotalErrors        int // Not serialized.
	LastResponseStatus int
	LastResponseError  string
}

func NewTransactionUploader(tt *TransactionChunks, client *client.Client) (*TransactionUploader, error) {
	if tt.ID == "" {
		return nil, errors.New("TransactionChunks is not signed.")
	}
	if tt.Chunks == nil {
		log.Warnf("TransactionChunks chunks not perpared.")
	}
	// Make a copy of Transaction, zeroing the Data so we can serialize.
	tu := &TransactionUploader{
		Client: client,
	}
	tu.Data = tt.Data
	tu.Transaction = &TransactionChunks{
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
	tChunks := tt.Transaction.Chunks
	if tChunks == nil {
		return false
	} else {
		return tt.TxPosted && (tt.ChunkIndex == len(tChunks.Chunks))
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

/**
 * Uploads the next part of the Transaction.
 * On the first call this posts the Transaction
 * itself and on any subsequent calls uploads the
 * next chunk until it completes.
 */
func (tt *TransactionUploader) UploadChunk() error {
	defer func() {
		fmt.Printf("%f%% completes, %d/%d \n", tt.PctComplete(), tt.UploadedChunks(), tt.TotalChunks())
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
		delay = math.Max(float64(tt.LastRequestTimeEnd+ERROR_DELAY-time.Now().UnixNano()/1000000), float64(ERROR_DELAY))
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

	chunk, err := tt.Transaction.GetChunk(tt.ChunkIndex, tt.Data)
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
	_, chunkOk := merkle.ValidatePath(tt.Transaction.Chunks.DataRoot, offset, 0, dataSize, path)
	if !chunkOk {
		return errors.New(fmt.Sprintf("Unable to validate chunk %d ", tt.ChunkIndex))
	}
	// Catch network errors and turn them into objects with status -1 and an error message.
	gc, err := tt.Transaction.GetChunk(tt.ChunkIndex, tt.Data)
	if err != nil {
		return err
	}
	byteGc, err := gc.Marshal()
	if err != nil {
		return err
	}
	body, statusCode, err := tt.Client.HttpPost("chunk", byteGc)
	fmt.Println("post tx chunk body: ", string(body))
	tt.LastRequestTimeEnd = time.Now().UnixNano() / 1000000
	tt.LastResponseStatus = statusCode
	if statusCode == 200 {
		tt.ChunkIndex++
	} else if err != nil {
		tt.LastResponseError = err.Error()
		if _, ok := FATAL_CHUNK_UPLOAD_ERRORS[err.Error()]; ok {
			return errors.New(fmt.Sprintf("Fatal error uploading chunk %d:%v", tt.ChunkIndex, err))
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
	// Copy the serialized upload information, and Data passed in.
	upload.ChunkIndex = serialized.chunkIndex
	upload.LastRequestTimeEnd = serialized.lastRequestTimeEnd
	upload.LastResponseError = serialized.lastResponseError
	upload.LastResponseStatus = serialized.lastResponseStatus
	upload.TxPosted = serialized.txPosted
	upload.Data = data

	upload.Transaction.PrepareChunks(data)

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
	body, statusCode, err := tt.Client.HttpGet(fmt.Sprintf("tx/%s", id))
	if err != nil || string(body) == "Pending" || statusCode/100 != 2 {
		return nil, errors.New(fmt.Sprintf("Tx %s not found: %d", id, statusCode))
	}
	transaction := &TransactionChunks{}
	if err := json.Unmarshal(body, transaction); err != nil {
		log.Errorf("json.Unmarshal(body, Transaction) error; body: %s", string(body))
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
	var uploadInBody = tt.TotalChunks() <= MAX_CHUNKS_IN_BODY

	if uploadInBody {
		// Post the Transaction with Data.
		tt.Transaction.Data = tt.Data
		byteTx, err := json.Marshal(tt.Transaction)
		if err != nil {
			return err
		}
		body, status, err := tt.Client.HttpPost("tx", byteTx)
		fmt.Printf("post tx with data;body: %s, status: %d, txId: %s \n", string(body), status, tt.Transaction.ID)
		if err != nil {
			fmt.Printf("tt.Client.SubmitTransaction(&tt.Transaction) error: %v", err)
			return err
		}
		tt.LastRequestTimeEnd = time.Now().UnixNano() / 1000000
		tt.LastResponseStatus = status
		tt.Transaction.Data = make([]byte, 0)

		if status >= 200 && status < 300 {
			// We are complete.
			tt.TxPosted = true
			tt.ChunkIndex = MAX_CHUNKS_IN_BODY
			return nil
		}
		tt.LastResponseError = ""
		return nil
	}

	byteTx, err := json.Marshal(tt.Transaction)
	if err != nil {
		return err
	}

	// else
	// Post the Transaction with no Data.
	body, status, err := tt.Client.HttpPost("tx", byteTx)
	fmt.Printf("post tx with no data; body: %s, status: %d, txId: %s \n", string(body), status, tt.Transaction.ID)
	tt.LastRequestTimeEnd = time.Now().UnixNano() / 1000000
	tt.LastResponseStatus = status
	if !(status >= 200 && status < 300) {
		if err != nil {
			tt.LastResponseError = err.Error()
		}
		return errors.New(fmt.Sprintf("Unable to upload Transaction: %d, %v", status, err))
	}
	tt.TxPosted = true
	return nil
}
