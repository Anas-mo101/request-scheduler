package server

import (
	"context"
	"encoding/json"
	database "task-scheduler/database/sqlc"
	"task-scheduler/datastore"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type RegSchedule struct {
	InvocationTimestamp pgtype.Timestamptz `json:"invocation_timestamp"`
	RequestMethod       database.Method    `json:"request_method"`
	RequestBody         string             `json:"request_body"`
	RequestHeader       []byte             `json:"request_header"`
	RequestQuery        []byte             `json:"request_query"`
	MaxRetries          pgtype.Int4        `json:"max_retries"`
	RequestUrl          string             `json:"request_url"`
}

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Post("/api/schedule", s.RegisterHandler)
}

func (s *FiberServer) RegisterHandler(c *fiber.Ctx) error {
	ctx := context.Background()

	defer func() {
		schedules, err := s.db.ListSchedule(ctx, 10)

		if err != nil {
			return
		}

		queue := datastore.GetQueueInstance()
		queue.SetQueue(schedules)
	}()

	req := new(RegSchedule)

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse request body",
		})
	}

	// Prepare the request header to be stored as JSONB
	requestHeaderJSON, err := json.Marshal(req.RequestHeader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to marshal request header",
		})
	}

	// Create a new schedule record
	schedule, err := s.db.CreateSchedule(ctx, database.CreateScheduleParams{
		InvocationTimestamp: req.InvocationTimestamp,
		RequestMethod:       req.RequestMethod,
		RequestBody:         req.RequestBody,
		RequestHeader:       requestHeaderJSON,
		MaxRetries:          req.MaxRetries,
		RequestQuery:        req.RequestQuery,
		RequestUrl:          req.RequestUrl,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create schedule",
		})
	}

	// Return the created schedule as a JSON response
	return c.Status(fiber.StatusOK).JSON(schedule)
}
