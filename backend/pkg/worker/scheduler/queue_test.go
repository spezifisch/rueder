package scheduler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Queue(t *testing.T) {
	timeA, _ := time.Parse("Jan 2 2006", "Jan 1 2000")
	timeB := timeA.Add(10 * time.Minute)
	timeC := timeA.Add(20 * time.Minute)

	pq := NewFeedQueue()
	pq.Push(&FeedQueueItem{deadline: timeA})
	pq.Push(&FeedQueueItem{deadline: timeB})
	pq.Push(&FeedQueueItem{deadline: timeC})

	assert.Equal(t, 3, pq.Len())
	var item FeedQueueItem
	item = pq.Peek()
	assert.Equal(t, timeA, item.deadline)
	assert.Equal(t, 3, pq.Len())
	item = pq.Pop()
	assert.Equal(t, timeA, item.deadline)
	item = pq.Pop()
	assert.Equal(t, timeB, item.deadline)
	item = pq.Pop()
	assert.Equal(t, timeC, item.deadline)
}

func Test_Reverse(t *testing.T) {
	timeA, _ := time.Parse("Jan 2 2006", "Jan 1 2000")
	timeB := timeA.Add(10 * time.Minute)
	timeC := timeA.Add(20 * time.Minute)

	pq := NewFeedQueue()
	pq.Push(&FeedQueueItem{deadline: timeC})
	pq.Push(&FeedQueueItem{deadline: timeB})
	pq.Push(&FeedQueueItem{deadline: timeA})

	assert.Equal(t, 3, pq.Len())
	item := pq.Pop()
	assert.Equal(t, timeA, item.deadline)
	item = pq.Pop()
	assert.Equal(t, timeB, item.deadline)
	item = pq.Pop()
	assert.Equal(t, timeC, item.deadline)
}

func Test_EqualItems(t *testing.T) {
	timeA, _ := time.Parse("Jan 2 2006", "Jan 1 2000")
	timeB := timeA
	timeC := timeA.Add(20 * time.Minute)

	pq := NewFeedQueue()
	pq.Push(&FeedQueueItem{deadline: timeC})
	pq.Push(&FeedQueueItem{deadline: timeB})
	pq.Push(&FeedQueueItem{deadline: timeA})

	assert.Equal(t, 3, pq.Len())
	item := pq.Pop()
	assert.Equal(t, timeA, item.deadline)
	item = pq.Pop()
	assert.Equal(t, timeB, item.deadline)
	item = pq.Pop()
	assert.Equal(t, timeC, item.deadline)
}

func Test_EqualItemsReverse(t *testing.T) {
	timeA, _ := time.Parse("Jan 2 2006", "Jan 1 2000")
	timeB := timeA
	timeC := timeA.Add(20 * time.Minute)

	pq := NewFeedQueue()
	pq.Push(&FeedQueueItem{deadline: timeA})
	pq.Push(&FeedQueueItem{deadline: timeB})
	pq.Push(&FeedQueueItem{deadline: timeC})

	assert.Equal(t, 3, pq.Len())
	item := pq.Pop()
	assert.Equal(t, timeA, item.deadline)
	item = pq.Pop()
	assert.Equal(t, timeB, item.deadline)
	item = pq.Pop()
	assert.Equal(t, timeC, item.deadline)
}

func Test_Push(t *testing.T) {
	timeA, _ := time.Parse("Jan 2 2006", "Jan 1 2000")
	timeB := timeA.Add(10 * time.Minute)
	timeC := timeA.Add(20 * time.Minute)
	timeD := timeA.Add(30 * time.Minute)
	timeE := timeA.Add(40 * time.Minute)

	pq := NewFeedQueue()
	pq.Push(&FeedQueueItem{deadline: timeA})
	pq.Push(&FeedQueueItem{deadline: timeB})
	pq.Push(&FeedQueueItem{deadline: timeC})

	assert.Equal(t, 3, pq.Len())
	pq.Push(&FeedQueueItem{deadline: timeD})
	assert.Equal(t, 4, pq.Len())
	pq.Push(&FeedQueueItem{deadline: timeE})
	assert.Equal(t, 5, pq.Len())

	item := pq.Pop()
	assert.Equal(t, timeA, item.deadline)
	item = pq.Pop()
	assert.Equal(t, timeB, item.deadline)
	item = pq.Pop()
	assert.Equal(t, timeC, item.deadline)
	item = pq.Pop()
	assert.Equal(t, timeD, item.deadline)
	item = pq.Pop()
	assert.Equal(t, timeE, item.deadline)
}
