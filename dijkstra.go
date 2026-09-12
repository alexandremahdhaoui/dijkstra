package dijkstra

import (
	"container/heap"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

type Node[T comparable, W Number] interface {
	comparable
	Edges() []Edge[T, W]
}

type Edge[T comparable, W Number] struct {
	To     T
	Weight W
}

type State[T comparable, W Number] struct {
	Node     T
	ParentOf T
	IsStart  bool
	Distance W
	Visited  bool
}

type OptimumState[T comparable, W Number] map[T]*State[T, W]

func (o OptimumState[T, W]) PathTo(end T) ([]T, W, bool) {
	targetState, exists := o[end]
	if !exists {
		var zero W
		return nil, zero, false
	}

	var path []T
	curr := end
	for {
		path = append(path, curr)
		currState := o[curr]
		if currState == nil || currState.IsStart {
			break
		}
		curr = currState.ParentOf
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, targetState.Distance, true
}

type NodeSelector[T comparable] func(n T) bool

func By[T comparable](target T) NodeSelector[T] {
	return func(n T) bool { return n == target }
}

type item[T comparable, W Number] struct {
	node T
	dist W
}

type priorityQueue[T comparable, W Number] []*item[T, W]

func (pq priorityQueue[T, W]) Len() int           { return len(pq) }
func (pq priorityQueue[T, W]) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq priorityQueue[T, W]) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue[T, W]) Push(x any)        { *pq = append(*pq, x.(*item[T, W])) }
func (pq *priorityQueue[T, W]) Pop() any {
	old := *pq
	n := len(old)
	it := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return it
}

func Search[T Node[T, W], W Number](root T, start, end NodeSelector[T]) OptimumState[T, W] {
	optimum := make(OptimumState[T, W])
	pq := &priorityQueue[T, W]{}
	heap.Init(pq)

	if start == nil {
		start = By(root)
	}

	visited := make(map[T]bool)
	queue := []T{root}
	visited[root] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if start(curr) {
			optimum[curr] = &State[T, W]{Node: curr, IsStart: true, Distance: 0}
			heap.Push(pq, &item[T, W]{node: curr, dist: 0})
		}

		for _, edge := range curr.Edges() {
			if !visited[edge.To] {
				visited[edge.To] = true
				queue = append(queue, edge.To)
			}
		}
	}

	for pq.Len() > 0 {
		currItem := heap.Pop(pq).(*item[T, W])
		curr := currItem.node
		currState := optimum[curr]

		if currItem.dist > currState.Distance {
			continue
		}
		currState.Visited = true

		if end != nil && end(curr) {
			break
		}

		for _, edge := range curr.Edges() {
			neighbor := edge.To
			neighState, exists := optimum[neighbor]

			newDist := currState.Distance + edge.Weight

			if !exists {
				optimum[neighbor] = &State[T, W]{
					Node:     neighbor,
					ParentOf: curr,
					Distance: newDist,
				}
				heap.Push(pq, &item[T, W]{node: neighbor, dist: newDist})
			} else if !neighState.Visited && newDist < neighState.Distance {
				neighState.Distance = newDist
				neighState.ParentOf = curr
				heap.Push(pq, &item[T, W]{node: neighbor, dist: newDist})
			}
		}
	}

	return optimum
}
