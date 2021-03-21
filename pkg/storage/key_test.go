package storage_test

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"testing"

	"github.com/partyzanex/testutils"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/stretchr/testify/assert"
)

func TestFastHash_HashBytes32(t *testing.T) {
	str1 := []byte("ab")
	str2 := make([]byte, len(str1))
	copy(str2, str1)

	h1 := fnv1a.HashBytes32(str1)
	h2 := fnv1a.HashBytes32(str2)

	assert.Equal(t, h1, h2)
	t.Log(h1, h2)

	str3 := append(str1, str2...)
	h3 := fnv1a.HashBytes32(str3)
	t.Log(h3)
}

func BenchmarkFastHash_HashBytes32(b *testing.B) {
	str := []byte(testutils.RandomString(128))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = fnv1a.HashBytes32(str)
	}
}

func BenchmarkFastHash_HashBytes64(b *testing.B) {
	str := []byte(testutils.RandomString(128))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = fnv1a.HashBytes64(str)
	}
}

func BenchmarkMD5(b *testing.B) {
	str := []byte(testutils.RandomString(128))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = md5.Sum(str)
	}
}

func BenchmarkSHA1(b *testing.B) {
	str := []byte(testutils.RandomString(128))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sha1.Sum(str)
	}
}

func BenchmarkSHA224(b *testing.B) {
	str := []byte(testutils.RandomString(128))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sha256.Sum224(str)
	}
}

func BenchmarkSHA256(b *testing.B) {
	str := []byte(testutils.RandomString(128))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sha256.Sum256(str)
	}
}

func BenchmarkSHA384(b *testing.B) {
	str := []byte(testutils.RandomString(128))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sha512.Sum384(str)
	}
}

func BenchmarkSHA512_224(b *testing.B) {
	str := []byte(testutils.RandomString(128))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sha512.Sum512_224(str)
	}
}

func BenchmarkSHA512_256(b *testing.B) {
	str := []byte(testutils.RandomString(128))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sha512.Sum512_256(str)
	}
}

func BenchmarkSHA512(b *testing.B) {
	str := []byte(testutils.RandomString(128))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sha512.Sum512(str)
	}
}
