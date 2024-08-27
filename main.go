package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"task-scheduler/invoker"
	"task-scheduler/server"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func dbConnect(ctx context.Context) (*pgx.Conn, error) {
	dsn := url.URL{
		Scheme: "postgres",
		Host:   os.Getenv("DATABASE_HOST") + ":" + os.Getenv("DATABASE_PORT"),
		User: url.UserPassword(
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
		),
		Path: os.Getenv("DATABASE_NAME"),
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")
	dsn.RawQuery = q.Encode()

	println(dsn.String())

	conn, err := pgx.Connect(ctx, dsn.String())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	conn, err := dbConnect(ctx)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	server := server.New(conn)
	server.RegisterFiberRoutes()

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err = server.Listen(fmt.Sprintf(":%d", port))

	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}

	invoker.Init(conn)

	invoker.Wg.Wait()
}
