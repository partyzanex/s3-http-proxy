package storage

type Interface interface {
	Get(key uint32) (*Object, bool)
	Set(key uint32, value *Object)
	Remove(key uint32)
}
