package scheduler

import (
	"time"

	"github.com/mailgun/holster/v3/collections"
)

// FeedQueueItem is a MinHeap item
type FeedQueueItem struct {
	feed Feed
	// timestamp when the feed should be fetched next
	deadline time.Time
}

// A FeedQueue is a MinHeap which pops the next due deadline first
type FeedQueue struct {
	pq *collections.PriorityQueue
}

// NewFeedQueue creates a new MinHeap
func NewFeedQueue() *FeedQueue {
	return &FeedQueue{
		pq: collections.NewPriorityQueue(),
	}
}

// Len gives the element count in the queue
func (f FeedQueue) Len() int { return f.pq.Len() }

// Push is implemented for heap.Interface
func (f *FeedQueue) Push(item *FeedQueueItem) {
	pqi := collections.PQItem{
		Value: item,
		// the collections.PriorityQueue pops the lowest int value first.
		// this is exactly what we want as we want to pop the lowest timestamp.
		Priority: int(item.deadline.Unix()),
	}
	f.pq.Push(&pqi)
}

// Pop is implemented for heap.Interface
func (f *FeedQueue) Pop() FeedQueueItem {
	pqi := f.pq.Pop()
	val := pqi.Value.(*FeedQueueItem)
	return *val
}

// Peek returns the element that would be returned by Pop but without removing it
func (f *FeedQueue) Peek() FeedQueueItem {
	pqi := f.pq.Peek()
	val := pqi.Value.(*FeedQueueItem)
	return *val
}
