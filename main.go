package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"task-scheduler/invoker"
	"task-scheduler/server"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "user=pqgotest dbname=pqgotest sslmode=verify-full")
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	server := server.New(conn)

	go invoker.New(conn)

	server.RegisterFiberRoutes()

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err = server.Listen(fmt.Sprintf(":%d", port))

	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
