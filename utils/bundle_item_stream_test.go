package utils

import (
	"bytes"
	"io"
	"testing"

	"github.com/everFinance/goar/types"
	"github.com/stretchr/testify/assert"
)

// oneByteReader returns data one byte per Read call — a valid io.Reader whose
// Read returns fewer than len(p) bytes (streams, bufio, os.Pipe behave this way).
type oneByteReader struct {
	data []byte
	pos  int
}

func (o *oneByteReader) Read(p []byte) (int, error) {
	if o.pos >= len(o.data) {
		return 0, io.EOF
	}
	p[0] = o.data[o.pos]
	o.pos++
	return 1, nil
}

// A negative tagsBytesLength (from an overflowed 8-byte field) must be rejected,
// not passed to make([]byte, tagsBytesLength) which panics "makeslice: len out of range".
func TestDecodeBundleItemStream_NegativeTagsLenDoesNotPanic(t *testing.T) {
	sc, ok := types.SigConfigMap[1]
	if !ok {
		t.Skip("sig type 1 not configured")
	}
	buf := &bytes.Buffer{}
	buf.WriteByte(1) // sigType 1 (low byte)
	buf.WriteByte(0)
	buf.Write(make([]byte, sc.SigLength))
	buf.Write(make([]byte, sc.PubLength))
	buf.WriteByte(0) // target absent
	buf.WriteByte(0) // anchor absent
	numOfTags := make([]byte, 8)
	numOfTags[0] = 1 // numOfTags = 1 (>0)
	buf.Write(numOfTags)
	tagsLen := make([]byte, 8)
	for i := 0; i < 8; i++ {
		tagsLen[i] = 0xff // overflows to -1
	}
	buf.Write(tagsLen)

	assert.NotPanics(t, func() {
		_, err := DecodeBundleItemStream(bytes.NewReader(buf.Bytes()))
		assert.Error(t, err)
	})
}

// A valid streaming reader that returns fewer bytes than requested per Read
// (io.ReadFull semantics) must decode the same as a bytes.Reader, not be
// rejected as "itemBinary incorrect".
func TestDecodeBundleItemStream_HandlesChunkedReader(t *testing.T) {
	sc, ok := types.SigConfigMap[1]
	if !ok {
		t.Skip("sig type 1 not configured")
	}
	// Minimal well-formed header with no tags and no data.
	item := make([]byte, 2+sc.SigLength+sc.PubLength+1+1+8+8)
	item[0] = 1 // sigType 1

	full, errFull := DecodeBundleItemStream(bytes.NewReader(item))
	assert.NoError(t, errFull)

	stream, errStream := DecodeBundleItemStream(&oneByteReader{data: item})
	assert.NoError(t, errStream, "a chunked reader must not be rejected")
	if errFull == nil && errStream == nil {
		assert.Equal(t, full.Id, stream.Id)
		assert.Equal(t, full.Signature, stream.Signature)
	}
}
