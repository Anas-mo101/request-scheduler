package database

import (
	"context"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5"
)

func DbConnect(ctx context.Context) (*pgx.Conn, error) {
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
