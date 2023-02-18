package utils

import (
	"errors"
	"io"
)

type ReadBuffer struct {
	bs []byte
	i  int
}

func NewReadBuffer(bs []byte) *ReadBuffer {
	return &ReadBuffer{bs: bs}
}
func (s *ReadBuffer) GetBytes() []byte {
	return s.bs
}
func (s *ReadBuffer) GetIndex() int {
	return s.i
}
func (s *ReadBuffer) Len() int {
	return len(s.bs)
}

// io.Reader
func (s *ReadBuffer) Read(p []byte) (int, error) {
	if s.i >= len(s.bs)-1 {
		return 0, nil
	}
	n := copy(p, s.bs[s.i:])
	s.i += n
	return n, nil
}

func (s *ReadBuffer) Seek(offset int64, whence int) (int64, error) {
	i := 0
	switch whence {
	case io.SeekStart:
		i = int(offset)
	case io.SeekCurrent:
		i += int(offset)
	case io.SeekEnd:
		i += len(s.bs) + int(offset)
	}
	if i < 0 || i >= len(s.bs) {
		return 0, errors.New("seek index out of buffer")
	}
	c := i - s.i
	s.i = i
	return int64(c), nil
}
