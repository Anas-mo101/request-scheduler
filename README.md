Task Scheduler
- Inspired from GCP Task Task Scheduler
- Schedules GET/POST reqeuts based on time
- Stores schedules in time order
- Invokes requests within 1 mins
- "/schedule" endpoint to add new schedules
- runs loop to check on upcoming invokations

Schedule Model
 - id
 - date/time of invokaion with timezone
 - created at
 - request method GET|POST
 - request body
 - requset body type: TEST|JSON 
 - request header
 - status: Scheduled|Invoked|Failed
 - retries no.
 - max retries
 - failure reason

Storing Data
- system is read heavy
- therefore effective data retrieval strategy is needed to avoid querying too often
- most recent schedules wil be stored in memory
- PostgresSQL for perminant data storage

Steps to run
- run "make build"
- run "make run"
- on a separate terminal run "make goose/migrate"


Schedule a request using /POST
```bash
curl --location 'http://172.22.0.3:46427/api/schedule' \
--header 'Content-Type: application/json' \
--data '{
    "invocation_timestamp": "2024-08-28T11:19:00.000+08:00",
    "request_method": "GET",
    "request_body": "{}",
    "request_header": {
        "User-Agent": "PostmanRuntime/7.41.2"
    },
    "request_query": {
        "results": "50",
        "seed": "foo"
    },
    "max_retries": 3,
    "request_url": "https://randomuser.me/api/sd",
    "request_body_type": "JSON"
}
'
```

