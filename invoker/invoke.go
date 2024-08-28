package invoker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	database "task-scheduler/database/sqlc"
	"time"
)

func invoke(schedule database.Schedule) {
	Wg.Add(1)
	defer Wg.Done()

	req, err := http.NewRequest(
		string(schedule.RequestMethod),
		schedule.RequestUrl,
		bytes.NewBufferString(schedule.RequestBody.String),
	)

	if err != nil {
		ch <- InvokedSchedule{
			schedule: schedule,
			err:      errors.New(fmt.Sprintf("failed to create HTTP request: %s", err.Error())),
		}
		return
	}

	if schedule.RequestBodyType == "TEXT" {
		req.Header.Add("Content-Type", "text/plain")
	}

	if schedule.RequestBodyType == "JSON" {
		req.Header.Add("Content-Type", "application/json")
	}

	// Add headers to the request
	var header map[string]string
	jerr := json.Unmarshal(schedule.RequestHeader, &header)
	if jerr != nil {
		ch <- InvokedSchedule{
			schedule: schedule,
			err:      errors.New(fmt.Sprintf("failed to parse HTTP body: %s", jerr.Error())),
		}
		return
	}

	for key, value := range header {
		req.Header.Add(key, value)
	}

	var query map[string]string
	qerr := json.Unmarshal(schedule.RequestQuery, &query)
	if qerr != nil {
		ch <- InvokedSchedule{
			schedule: schedule,
			err:      errors.New(fmt.Sprintf("failed to parse HTTP query: %s", jerr.Error())),
		}
		return
	}

	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	// Create an HTTP client and execute the request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		ch <- InvokedSchedule{
			schedule: schedule,
			err:      errors.New(fmt.Sprintf("failed to execute HTTP request: %s", err.Error())),
		}
		return
	}
	defer resp.Body.Close()

	// Handle the response (for example, logging it)
	fmt.Printf("Invoked: %s\n", resp.Status)

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if statusOK {
		ch <- InvokedSchedule{
			schedule: schedule,
			err:      nil,
		}
		return
	}

	ch <- InvokedSchedule{
		schedule: schedule,
		err:      errors.New(fmt.Sprintf("failed HTTP request %s: %s", resp.StatusCode, err.Error())),
	}
}
