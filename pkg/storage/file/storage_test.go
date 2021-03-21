package file_test

import (
	"testing"
	"time"

	"github.com/partyzanex/s3-http-proxy/pkg/storage"
	"github.com/partyzanex/s3-http-proxy/pkg/storage/file"
	"github.com/partyzanex/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInFileStorage_Get(t *testing.T) {
	path := "/tmp/file-storage"
	store, err := file.New(path, 100)
	require.NoError(t, err)
	require.NotNil(t, store)

	key := storage.NewKeyFromBytes([]byte(testutils.RandomString(111)))

	t.Cleanup(func() {
		store.Remove(key)
	})

	exp := storage.Object{
		Body:         []byte(testutils.RandomString(1111)),
		ContentType:  uint8(testutils.RandInt(1, 16)),
		LastModified: time.Now().UnixNano(),
	}
	err = store.Set(key, &exp)
	require.NoError(t, err)

	got, err := store.Get(key)
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, &exp, got)
}
