Task Scheduler


- Inspired from GCP Task Task Scheduler
- Schedules GET/POST reqeuts based on time
- Stores schedules in time order
- Invokes requests within 1 mins
- "/register" endpoint to add new schedules
- runs loop to check on upcoming invokations



- Schedule Model
 -- id
 -- date/time of invokaion with timezone
 -- created at
 -- request method Get|Post
 -- request body
 -- request header
 -- status: Scheduled|Invoked|Failed
 -- retries no.
 -- failure reason


- Storing Data
 -- this system is read heavy
 -- therefore effective data retraval strategy is needed
 -- to avoid querying too often, most recent schedules wil be stored in memory
 -- PostgresSQL for perminant data storage
