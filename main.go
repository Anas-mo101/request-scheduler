package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"task-scheduler/database"
	"task-scheduler/invoker"
	"task-scheduler/server"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	conn, err := database.DbConnect(ctx)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	server := server.New(conn)
	server.RegisterFiberRoutes()

	invoker.Init(conn)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err = server.Listen(fmt.Sprintf(":%d", port))

	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}

	invoker.Wg.Wait()
}
