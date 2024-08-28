package server

import (
	"context"
	"encoding/json"
	"fmt"
	database "task-scheduler/database/sqlc"
	"task-scheduler/datastore"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type RegSchedule struct {
	InvocationTimestamp pgtype.Timestamptz `json:"invocation_timestamp"`
	RequestMethod       database.Method    `json:"request_method"`
	RequestBody         pgtype.Text        `json:"request_body"`
	RequestHeader       map[string]string  `json:"request_header"`
	RequestQuery        map[string]string  `json:"request_query"`
	MaxRetries          pgtype.Int4        `json:"max_retries"`
	RequestUrl          string             `json:"request_url"`
	RequestBodyType     database.BodyType  `json:"request_body_type"`
}

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Post("/api/schedule", s.RegisterHandler)
}

func (s *FiberServer) RegisterHandler(c *fiber.Ctx) error {
	ctx := context.Background()

	req := new(RegSchedule)

	if err := c.BodyParser(&req); err != nil {
		fmt.Printf(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse request body",
		})
	}

	// Prepare the request header to be stored as JSONB
	requestHeaderJSON, err := json.Marshal(req.RequestHeader)
	if err != nil {
		fmt.Printf(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to marshal request header",
		})
	}

	requestQueryJSON, err := json.Marshal(req.RequestQuery)
	if err != nil {
		fmt.Printf(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to marshal RequestQuery",
		})
	}

	// Create a new schedule record
	schedule, err := s.db.CreateSchedule(ctx, database.CreateScheduleParams{
		InvocationTimestamp: req.InvocationTimestamp,
		RequestMethod:       req.RequestMethod,
		RequestBody:         req.RequestBody,
		RequestHeader:       requestHeaderJSON,
		MaxRetries:          req.MaxRetries,
		RequestQuery:        requestQueryJSON,
		RequestUrl:          req.RequestUrl,
		RequestBodyType:     req.RequestBodyType,
	})

	if err != nil {
		fmt.Printf(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create schedule",
		})
	}

	queue := datastore.GetQueueInstance()
	queue.EnQueueWithinRange(schedule)

	// Return the created schedule as a JSON response
	return c.Status(fiber.StatusOK).JSON(schedule)
}
