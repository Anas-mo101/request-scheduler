package invoker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	database "task-scheduler/database/sqlc"
	"time"
)

func invoke(schedule database.Schedule) {
	ctx := context.Background()

	req, err := http.NewRequest(
		string(schedule.RequestMethod),
		schedule.RequestUrl,
		bytes.NewBufferString(schedule.RequestBody),
	)

	if err != nil {
		fmt.Printf("failed to create HTTP request: %s", err.Error())
	}

	// Add headers to the request
	var header map[string]string
	jerr := json.Unmarshal(schedule.RequestHeader, &header)
	if jerr != nil {
		fmt.Printf("failed to parse HTTP body: %s", jerr.Error())
		return
	}

	for key, value := range header {
		req.Header.Add(key, value)
	}

	// Create an HTTP client and execute the request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		// Update the status to 'Failed'
		updatedSchedule, _ := queries.IncrementFailure(ctx, database.IncrementFailureParams{
			ID:            schedule.ID,
			FailureReason: err.Error(),
		})

		if updatedSchedule.MaxRetries.Int32 < updatedSchedule.RetriesNo.Int32 {
			go invoke(updatedSchedule)
		}

		fmt.Printf("failed to execute HTTP request: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	// Handle the response (for example, logging it)
	fmt.Printf("Response Status: %s\n", resp.Status)

	// Update the status to 'Invoked'
	_, _ = queries.ScheduleSuccss(ctx, schedule.ID)
}
