package data

import (
	"fmt"
	"sync"
	"task-scheduler/database"
)

// Queue represents a FIFO queue
type ScheduleQueue struct {
	items []database.Schedule
	mu    sync.Mutex
}

var instance *ScheduleQueue
var once sync.Once

func GetQueueInstance() *ScheduleQueue {
	once.Do(func() {
		instance = &ScheduleQueue{}
	})
	return instance
}

// Enqueue adds an item to the end of the queue
func (q *ScheduleQueue) Enqueue(item database.Schedule) {
	q.items = append(q.items, item)
}

func (q *ScheduleQueue) SetQueue(items []database.Schedule) {
	q.items = items
}

// Dequeue removes and returns the item at the front of the queue
func (q *ScheduleQueue) Dequeue() (database.Schedule, error) {
	if len(q.items) == 0 {
		var zeroValue database.Schedule
		return zeroValue, fmt.Errorf("queue is empty")
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

// IsEmpty checks if the queue is empty
func (q *ScheduleQueue) IsEmpty() bool {
	return len(q.items) == 0
}

// Peek returns the item at the front of the queue without removing it
func (q *ScheduleQueue) Peek() (database.Schedule, error) {
	if len(q.items) == 0 {
		var zeroValue database.Schedule
		return zeroValue, fmt.Errorf("queue is empty")
	}
	return q.items[0], nil
}
