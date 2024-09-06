package main

import (
	"context"
	"encoding/json"
	"fmt"
	db "task-scheduler/database"
	database "task-scheduler/database/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	conn, err := db.DbConnect(ctx)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	db := database.New(conn)

	for i := 0; i < 10; i++ {

		// Prepare the request header to be stored as JSONB
		requestHeaderJSON, err := json.Marshal("{}")
		if err != nil {
			fmt.Printf(err.Error())
		}

		requestQueryJSON, err := json.Marshal("{}")
		if err != nil {
			fmt.Printf(err.Error())
		}

		_, err = db.CreateSchedule(ctx, database.CreateScheduleParams{
			InvocationTimestamp: pgtype.Timestamptz{Time: time.Now().Add(time.Duration(i) * time.Minute), Valid: true},
			RequestMethod:       "POST",
			RequestBody:         pgtype.Text{String: "{}"},
			RequestHeader:       requestHeaderJSON,
			MaxRetries:          pgtype.Int4{Int32: 1},
			RequestQuery:        requestQueryJSON,
			RequestUrl:          "https://google.com",
			RequestBodyType:     database.NullBodyType{Valid: true, BodyType: database.BodyTypeTEXT},
		})

		if err != nil {
			fmt.Printf(err.Error())
		}
	}
}
