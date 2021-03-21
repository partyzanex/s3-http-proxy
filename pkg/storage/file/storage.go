package file

import (
	"encoding/gob"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/partyzanex/s3-http-proxy/pkg/pq"
	"github.com/partyzanex/s3-http-proxy/pkg/storage"
	"github.com/pkg/errors"
)

func closer(c io.Closer, err *error) {
	errC := c.Close()
	if errC != nil {
		if *err != nil {
			*err = errors.Wrap(*err, errC.Error())
			return
		}

		*err = errC
	}
}

type InFileStorage struct {
	path  string
	files map[uint32]bool
	mu    *sync.Mutex
	pq    pq.Interface
	size  int
	count uint64
}

func New(path string, size int) (storage.Interface, error) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create path")
	}

	store := &InFileStorage{
		path:  path,
		size:  size,
		files: make(map[uint32]bool),
		pq:    pq.New(1000),
		mu:    new(sync.Mutex),
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read path")
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			return nil, errors.Wrap(err, "cannot get info")
		}

		k, err := strconv.ParseUint(info.Name(), 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse file name")
		}

		key := uint32(k)
		store.files[key] = true
		store.pq.Push(&pq.Node{
			Value:    key,
			Size:     uint64(info.Size()),
			Priority: info.ModTime().UnixNano(),
		})
	}

	return store, nil
}

func (s *InFileStorage) Get(key uint32) (object *storage.Object, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.files[key]; !ok {
		return nil, errors.New("object not found")
	}

	filePath := s.getFilePath(key)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open file '%s'", filePath)
	}

	defer closer(file, &err)

	err = gob.NewDecoder(file).Decode(&object)
	if err != nil {
		return nil, errors.Wrap(err, "cannot decode")
	}

	return
}

func (s *InFileStorage) Set(key uint32, value *storage.Object) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := s.getFilePath(key)

	file, err := os.Create(filePath)
	if err != nil {
		return errors.Wrapf(err, "cannot create file '%s'", filePath)
	}

	defer closer(file, &err)

	err = gob.NewEncoder(file).Encode(value)
	if err != nil {
		return errors.Wrap(err, "cannot encode")
	}

	s.files[key] = true

	stat, err := file.Stat()
	if err != nil {
		return errors.Wrap(err, "cannot get file stat")
	}

	go s.resize(key, uint64(stat.Size()), stat.ModTime().UnixNano())

	return
}

func (s *InFileStorage) resize(key uint32, size uint64, mod int64) {
	n := pq.Node{
		Value:    key,
		Size:     size,
		Priority: mod,
	}
	s.pq.Push(&n)

	if uint64(s.size) > s.count+size {
		if node := s.pq.Next(); node != nil {
			k := node.Value.(uint32)
			s.Remove(k)
			atomic.AddUint64(&s.count, -node.Size)
		}
	}

	atomic.AddUint64(&s.count, size)
}

func (s *InFileStorage) Remove(key uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.Remove(s.getFilePath(key)); err == nil {
		delete(s.files, key)
	}
}

func (s *InFileStorage) getFilePath(key uint32) string {
	return filepath.Join(s.path, strconv.FormatUint(uint64(key), 10))
}
