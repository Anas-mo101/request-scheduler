package invoker

import (
	"context"
	"fmt"
	"sync"
	database "task-scheduler/database/sqlc"
	"task-scheduler/datastore"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var queries *database.Queries
var queue *datastore.ScheduleQueue
var ch chan InvokedSchedule
var Wg sync.WaitGroup

var timerChannel chan bool
var mainTimer *time.Ticker
var secondaryTimeer *time.Ticker

type InvokedSchedule struct {
	schedule database.Schedule
	err      error
}

func Init(conn *pgx.Conn) {
	fmt.Println("init: start")

	queries = database.New(conn)
	ch = make(chan InvokedSchedule)
	timerChannel = make(chan bool)

	queue = datastore.GetQueueInstance()

	load()
	go listen()

	mainTimer = schedule(loop, time.Minute, timerChannel)
	secondaryTimeer = schedule(load, 10*time.Minute, timerChannel)
}

func Terminate() {
	close(timerChannel)
	mainTimer.Stop()
	secondaryTimeer.Stop()
}

func listen() {
	ctx := context.Background()
	fmt.Println("listening to incoming invokes")

	for {
		invokedSchedule := <-ch

		if invokedSchedule.err == nil {
			// Update the status to 'Invoked'
			_, _ = queries.ScheduleSuccss(ctx, invokedSchedule.schedule.ID)
			continue
		}

		// Update the status to 'Failed'
		updatedSchedule, _ := queries.IncrementFailure(ctx, database.IncrementFailureParams{
			ID:            invokedSchedule.schedule.ID,
			FailureReason: pgtype.Text{String: invokedSchedule.err.Error()},
		})

		if updatedSchedule.MaxRetries.Int32 > updatedSchedule.RetriesNo.Int32 {
			go invoke(invokedSchedule.schedule)
		}
	}
}

func loop() {
	fmt.Println("loop queue for schedule")
	for {
		toCheck, err := queue.Peek()

		if err != nil {
			// in case queue is empty load again
			fmt.Println("queue is empty, attempting load")
			load()
			return
		}

		fmt.Println("to Check: ", toCheck.ID)

		isAfter := toCheck.InvocationTimestamp.Time.After(time.Now())

		if isAfter {
			fmt.Println("is After")
			return
		}

		schedule, _ := queue.Dequeue()

		fmt.Println("is time to invoke: ", schedule.ID)

		go invoke(schedule)
	}
}

func load() {
	fmt.Println("load queue")

	ctx := context.Background()

	// fetch most recent schedule
	schedules, err := queries.ListSchedule(ctx, queue.GetLimit())

	if err != nil {
		return
	}

	queue.SetQueue(schedules)
}

func schedule(f func(), interval time.Duration, done <-chan bool) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				f()
			case <-done:
				return
			}
		}
	}()
	return ticker
}
