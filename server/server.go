package server

import (
	"task-scheduler/database"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type FiberServer struct {
	*fiber.App

	db *database.Queries
}

func New(conn *pgx.Conn) *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "scheduler",
			AppName:      "scheduler",
		}),

		db: database.New(conn),
	}

	return server
}
