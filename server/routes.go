package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	database "task-scheduler/database/sqlc"
	"task-scheduler/datastore"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type RegSchedule struct {
	InvocationTimestamp pgtype.Timestamptz    `json:"invocation_timestamp"`
	RequestMethod       database.Method       `json:"request_method"`
	RequestBody         pgtype.Text           `json:"request_body"`
	RequestHeader       map[string]string     `json:"request_header"`
	RequestQuery        map[string]string     `json:"request_query"`
	MaxRetries          pgtype.Int4           `json:"max_retries"`
	RequestUrl          string                `json:"request_url"`
	RequestBodyType     database.NullBodyType `json:"request_body_type"`
}

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Post("/api/schedule", s.RegisterHandler)
	s.App.Delete("/api/schedule/:id", s.DeleteHandler)
	s.App.Put("/api/schedule/:id", s.UpdateHandler)
	s.App.Get("/api/schedule/:id", s.GetHandler)
	s.App.Get("/api/schedule", s.GetHandler)
}

func (s *FiberServer) UpdateHandler(c *fiber.Ctx) error {
	ctx := context.Background()

	id := c.Params("id")
	sid, err := strconv.ParseInt(id, 10, 32)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete schedule",
		})
	}

	req := new(RegSchedule)

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse request body",
		})
	}

	schedule, err := s.db.UpdateSchedule(ctx, database.UpdateScheduleParams{
		ID: int32(sid),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to udpate schedule",
		})
	}

	// Return the created schedule as a JSON response
	return c.Status(fiber.StatusOK).JSON(schedule)
}

func (s *FiberServer) listHandler(c *fiber.Ctx) error {
	ctx := context.Background()

	// Parsing pagination parameters (limit and offset)
	limitStr := c.Query("limit", "10")  // Default limit is 10
	offsetStr := c.Query("offset", "0") // Default offset is 0

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid limit value",
		})
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid offset value",
		})
	}

	// Parsing filters
	invocationTimestampStr := c.Query("invocation_timestamp")
	requestMethod := c.Query("request_method")
	requestUrl := c.Query("request_url")
	maxRetriesStr := c.Query("max_retries")
	status := c.Query("status")
	requestBodyType := c.Query("request_body_type")

	var maxRetries int
	if maxRetriesStr != "" {
		retries, err := strconv.Atoi(maxRetriesStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid max_retries value",
			})
		}
		maxRetries = retries
	}

	// Execute the query with filters
	schedules, err := s.db.ListRegSchedule(ctx, database.ListRegScheduleParams{
		Limit:   int32(limit),
		Offset:  int32(offset),
		Column3: invocationTimestampStr,
		Column4: requestMethod,
		Column5: requestUrl,
		Column6: maxRetries,
		Column7: requestBodyType,
		Column8: status,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to list schedules",
		})
	}

	// Return the list of schedules as a JSON response
	return c.Status(fiber.StatusOK).JSON(schedules)
}

func (s *FiberServer) GetHandler(c *fiber.Ctx) error {
	ctx := context.Background()

	id := c.Params("id")
	sid, err := strconv.ParseInt(id, 10, 32)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete schedule",
		})
	}

	schedule, err := s.db.GetSchedule(ctx, int32(sid))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get schedule",
		})
	}

	// Return the created schedule as a JSON response
	return c.Status(fiber.StatusOK).JSON(schedule)
}

func (s *FiberServer) DeleteHandler(c *fiber.Ctx) error {
	ctx := context.Background()

	id := c.Params("id")
	sid, err := strconv.ParseInt(id, 10, 32)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete schedule",
		})
	}

	schedule, err := s.db.DeletSchedule(ctx, int32(sid))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete schedule",
		})
	}

	return c.Status(fiber.StatusOK).JSON(schedule)
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
