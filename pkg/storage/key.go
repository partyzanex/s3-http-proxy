package storage

import "github.com/segmentio/fasthash/fnv1a"

func NewKeyFromBytes(b []byte) uint32 {
	return fnv1a.HashBytes32(b)
}

func NewKeyFromString(s string) uint32 {
	return NewKeyFromBytes([]byte(s))
}
