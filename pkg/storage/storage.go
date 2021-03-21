package storage

type Interface interface {
	Get(key uint32) (*Object, error)
	Set(key uint32, value *Object) error
	Remove(key uint32)
}
