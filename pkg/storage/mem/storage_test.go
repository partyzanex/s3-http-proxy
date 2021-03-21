package mem_test

import (
	"testing"
	"time"

	"github.com/cornelk/hashmap"
	"github.com/partyzanex/s3-http-proxy/pkg/storage"
	"github.com/partyzanex/s3-http-proxy/pkg/storage/mem"
	"github.com/partyzanex/testutils"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryStorage_Get(t *testing.T) {
	store := mem.New(10)
	expected := make(map[uint32]*storage.Object)

	for i := 0; i < 100; i++ {
		n := testutils.RandInt(10, 200)
		key := fnv1a.HashBytes32([]byte(testutils.RandomString(n)))
		obj := &storage.Object{
			Body:         []byte(testutils.RandomString(n * 10)),
			ContentType:  uint8(testutils.RandInt(1, 16)),
			LastModified: time.Now().UnixNano(),
		}
		store.Set(key, obj)

		expected[key] = obj
	}

	for key, exp := range expected {
		got, ok := store.Get(key)
		// require.NoError(t, err)
		require.Equal(t, exp, got)
		require.Equal(t, true, ok)
	}
}

func BenchmarkInMemoryStorage_Set(b *testing.B) {
	store := mem.New(10)
	currentTime := time.Now().UnixNano()
	s := []byte(testutils.RandomString(60))
	key := fnv1a.HashBytes32(s)
	obj := storage.Object{
		Body:         []byte(testutils.RandomString(600)),
		ContentType:  uint8(testutils.RandInt(1, 16)),
		LastModified: currentTime,
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.Set(key, &obj)
			_, _ = store.Get(key)
		}
	})
}

func TestHashMap_Get(t *testing.T) {
	hm := &hashmap.HashMap{}
	expected := make(map[storage.Key]*storage.Object)

	for i := 0; i < 100; i++ {
		n := testutils.RandInt(10, 200)
		key := storage.NewKeyFromBytes([]byte(testutils.RandomString(n)))
		obj := &storage.Object{
			Body:         []byte(testutils.RandomString(n * 10)),
			ContentType:  uint8(testutils.RandInt(1, 16)),
			LastModified: time.Now().UnixNano(),
		}
		hm.Set(uint32(key), obj)
		expected[key] = obj
	}

	for key, exp := range expected {
		got, ok := hm.Get(uint32(key))
		assert.Equal(t, true, ok)
		assert.Equal(t, exp, got.(*storage.Object))
	}
}

func BenchmarkHashMap_Set(b *testing.B) {
	hm := &hashmap.HashMap{}
	currentTime := time.Now().UnixNano()
	s := []byte(testutils.RandomString(60))
	key := storage.NewKeyFromBytes(s)
	obj := storage.Object{
		Body:         []byte(testutils.RandomString(600)),
		ContentType:  uint8(testutils.RandInt(1, 16)),
		LastModified: currentTime,
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			hm.Insert(uint32(key), &obj)
			_, _ = hm.Get(uint32(key))
		}
	})
}
