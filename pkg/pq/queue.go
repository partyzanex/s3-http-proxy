package pq

import (
	"container/heap"
	"sync"
)

type queue struct {
	nodes []*Node
	mu    *sync.Mutex
}

func New(capacity int) Interface {
	return &queue{
		nodes: make([]*Node, 0, capacity),
		mu:    new(sync.Mutex),
	}
}

func (q *queue) Len() int {
	return len(q.nodes)
}

func (q *queue) Less(i, j int) (ok bool) {
	return q.nodes[i].Priority > q.nodes[j].Priority
}

func (q *queue) Swap(i, j int) {
	q.nodes[i], q.nodes[j] = q.nodes[j], q.nodes[i]
	q.nodes[i].index = i
	q.nodes[j].index = j
}

func (q *queue) Push(x interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()

	n := len(q.nodes)
	node := x.(*Node)
	node.index = n

	q.nodes = append(q.nodes, node)
	heap.Fix(q, node.index)
}

func (q *queue) Pop() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()

	old := q.nodes
	k := len(old) - 1

	if k < 0 {
		return nil
	}

	node := old[k]
	old[k] = nil
	node.index = -1
	q.nodes = old[0:k]

	return node
}

func (q *queue) Next() *Node {
	return q.Pop().(*Node)
}
