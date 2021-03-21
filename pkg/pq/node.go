package pq

type Node struct {
	Value    interface{}
	Size     uint64
	Priority int64
	index    int
}
