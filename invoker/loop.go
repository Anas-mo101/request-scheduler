package invoker

import (
	"context"
	database "task-scheduler/database/sqlc"
	"task-scheduler/datastore"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jasonlvhit/gocron"
)

var queries *database.Queries
var queue *datastore.ScheduleQueue

func New(conn *pgx.Conn) {
	queries = database.New(conn)

	Load()

	gocron.Every(1).Minute().Do(PrimaryLoop)
	gocron.Every(10).Minute().Do(Load)
}

func PrimaryLoop() {
	toCheck, err := queue.Peek()

	if err != nil {
		return
	}

	isAfter := toCheck.InvocationTimestamp.Time.After(time.Now())

	if !isAfter {
		return
	}

	schedule, _ := queue.Dequeue()

	go invoke(schedule)
}

func Load() {
	ctx := context.Background()

	// fetch most recent schedule
	schedules, err := queries.ListSchedule(ctx, 10)

	if err != nil {
		return
	}

	queue := datastore.GetQueueInstance()
	queue.SetQueue(schedules)
}
