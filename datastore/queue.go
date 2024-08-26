package datastore

import (
	"fmt"
	"sync"
	database "task-scheduler/database/sqlc"
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

func (q *ScheduleQueue) SetQueue(items []database.Schedule) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = items
}

// Dequeue removes and returns the item at the front of the queue
func (q *ScheduleQueue) Dequeue() (database.Schedule, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		var zeroValue database.Schedule
		return zeroValue, fmt.Errorf("queue is empty")
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

// Peek returns the item at the front of the queue without removing it
func (q *ScheduleQueue) Peek() (database.Schedule, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		var zeroValue database.Schedule
		return zeroValue, fmt.Errorf("queue is empty")
	}
	return q.items[0], nil
}
