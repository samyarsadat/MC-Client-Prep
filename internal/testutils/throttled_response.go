package testutils

import (
	"io"
	"time"
)

type ThrottledReader struct {
	Reader    io.Reader
	ChunkSize uint64
	Delay     time.Duration
}

func (t *ThrottledReader) Read(p []byte) (n int, err error) {
	if uint64(len(p)) > t.ChunkSize {
		p = p[:t.ChunkSize]
	}

	n, err = t.Reader.Read(p)
	if n > 0 {
		time.Sleep(t.Delay)
	}

	return n, err
}
