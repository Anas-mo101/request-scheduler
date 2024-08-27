package invoker

import (
	"context"
	"fmt"
	"sync"
	database "task-scheduler/database/sqlc"
	"task-scheduler/datastore"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jasonlvhit/gocron"
)

var queries *database.Queries
var queue *datastore.ScheduleQueue
var ch chan InvokedSchedule
var Wg sync.WaitGroup

type InvokedSchedule struct {
	schedule database.Schedule
	err      error
}

func Init(conn *pgx.Conn) {
	queries = database.New(conn)
	ch = make(chan InvokedSchedule)

	load()
	go listen()

	gocron.Every(1).Minute().Do(loop)
	gocron.Every(10).Minute().Do(load)
}

func listen() {
	ctx := context.Background()

	for {
		invokedSchedule := <-ch

		if invokedSchedule.err != nil {
			// Update the status to 'Invoked'
			_, _ = queries.ScheduleSuccss(ctx, invokedSchedule.schedule.ID)
			continue
		}

		// Update the status to 'Failed'
		updatedSchedule, _ := queries.IncrementFailure(ctx, database.IncrementFailureParams{
			ID:            invokedSchedule.schedule.ID,
			FailureReason: invokedSchedule.err.Error(),
		})

		if updatedSchedule.MaxRetries.Int32 > updatedSchedule.RetriesNo.Int32 {
			go invoke(invokedSchedule.schedule)
		}
	}
}

func loop() {
	fmt.Print("Cron: loop")
	for {
		toCheck, err := queue.Peek()

		if err != nil {
			// in case queue is empty load again
			load()
			return
		}

		isAfter := toCheck.InvocationTimestamp.Time.After(time.Now())

		if !isAfter {
			return
		}

		schedule, _ := queue.Dequeue()

		go invoke(schedule)
	}
}

func load() {
	fmt.Print("Cron: load")

	ctx := context.Background()

	// fetch most recent schedule
	schedules, err := queries.ListSchedule(ctx, 10)

	if err != nil {
		return
	}

	queue := datastore.GetQueueInstance()
	queue.SetQueue(schedules)
}
