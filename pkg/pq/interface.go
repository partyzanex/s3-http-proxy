package pq

import (
	"container/heap"
)

type Interface interface {
	heap.Interface

	Next() *Node
}
