package mem

import (
	"errors"
	"sync/atomic"

	"github.com/DmitriyVTitov/size"
	"github.com/cornelk/hashmap"
	"github.com/partyzanex/s3-http-proxy/pkg/pq"
	"github.com/partyzanex/s3-http-proxy/pkg/storage"
)

const defaultQueCapacity = 1000

type InMemoryStorage struct {
	hm    *hashmap.HashMap
	size  int
	count uint64
	pq    pq.Interface
}

func New(size int) storage.Interface {
	s := &InMemoryStorage{
		hm:    new(hashmap.HashMap),
		size:  size,
		count: 0,
		pq:    pq.New(defaultQueCapacity),
	}

	return s
}

func (s *InMemoryStorage) Get(key uint32) (*storage.Object, error) {
	v, ok := s.hm.Get(key)
	if !ok {
		return nil, errors.New("object not found")
	}

	return v.(*storage.Object), nil
}

func (s *InMemoryStorage) Set(key uint32, value *storage.Object) error {
	s.hm.GetOrInsert(key, value)

	go s.resize(key, value)

	return nil
}

func (s *InMemoryStorage) resize(key uint32, value *storage.Object) {
	sz := uint64(size.Of(value))

	n := pq.Node{
		Value:    key,
		Size:     sz,
		Priority: value.LastModified,
	}
	s.pq.Push(&n)

	if uint64(s.size) > s.count+sz {
		if node := s.pq.Next(); node != nil {
			k := node.Value.(uint32)
			s.Remove(k)
			atomic.AddUint64(&s.count, -node.Size)
		}
	}

	atomic.AddUint64(&s.count, sz)
}

func (s *InMemoryStorage) Remove(key uint32) {
	s.hm.Del(key)
}
