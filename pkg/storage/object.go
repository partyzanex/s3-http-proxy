package storage

import (
	"bytes"
	"io"
)

type Object struct {
	Body         []byte
	ContentType  uint
	LastModified int64
}

func (obj Object) GetMimeType(mimeTypes map[uint][]byte) []byte {
	return mimeTypes[obj.ContentType]
}

func (obj Object) NewReadCloser() io.ReadCloser {
	return &readCloser{
		r: bytes.NewReader(obj.Body),
	}
}

type readCloser struct {
	r io.Reader
}

func (rc *readCloser) Read(b []byte) (int, error) {
	return rc.r.Read(b)
}

func (rc *readCloser) Close() error {
	rc.r = nil

	return nil
}
