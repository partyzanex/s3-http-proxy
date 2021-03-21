package storage

type Object struct {
	Body         []byte
	ContentType  uint8
	LastModified int64
}
