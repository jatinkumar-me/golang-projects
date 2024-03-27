package main

import (
	"container/heap"
	"fmt"
)

// Item represents an item in the priority queue.
type Item struct {
	Key   int // The key of the item
	Value int // The value of the item
}

// PriorityQueue implements a priority queue based on a min-heap.
type PriorityQueue []*Item

// Len returns the length of the priority queue.
func (pq PriorityQueue) Len() int { return len(pq) }

// Less compares two items based on their values.
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Value < pq[j].Value
}

// Swap swaps two items in the priority queue.
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push adds an item to the priority queue.
func (pq *PriorityQueue) Push(x any) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

// Pop removes and returns the top item from the priority queue.
func (pq *PriorityQueue) Pop() interface{} {
	n := len(*pq)
	item := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return item
}

func main() {
	// Create a priority queue
	pq := make(PriorityQueue, 0)

	// Push items to the priority queue
	heap.Push(&pq, &Item{Key: 1, Value: 5})
	heap.Push(&pq, &Item{Key: 2, Value: 3})
	heap.Push(&pq, &Item{Key: 3, Value: 7})

	// Pop items from the priority queue in order of priority
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("Key: %d, Value: %d\n", item.Key, item.Value)
	}
}
