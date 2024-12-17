# message-service-kata
message service

## Requirements

- Go `>= v1.19`
- Postgres
- Kafka `https://github.com/confluentinc/confluent-kafka-go`

## Environment Variables

- `APP_DEBUG`: bool
  Run in debug mode?
- `APP_ADDRESS`: IP Address
  Address to bind
- `APP_BUILD_ENV`: ENV mode
  The environment build
- `APP_READ_TIMEOUT`: time.Duration
  The amount of time allowed to read request headers
- `APP_WRITE_TIMEOUT`: time.Duration
  The maximum duration before timing out writes of the response

- `PG_DBNAME`: str
  Postgres DB Name
- `PG_HOST`: address
  Postgres DB Host address
- `PG_PORT`: int
  Postgres DB Port
- `PG_DBUSER`: str
  Postgres DB user
- `PG_DBPASS`: str
  Postgres DB passwor
- `PG_SSL_MODE`: str
  Postgres DB ssl mode
- `PG_CONN_MAX_LIFETIME`: str
  Postgres DB connection max lifetime
- `PG_MAX_IDLE_CONNS`: int
  Postgres DB max idle connection allowance
- `PG_MAX_OPEN_CONNS`: int
  Postgres DB max open connection allowance

- `KAFKA_BROKER_ADDR`: address
   Kafka Address

## Available Command
# Running Service
- API server `go run ./cmd/message-service-kata -service rest` or ` make serve-rest`
- Consumer serve `go run ./cmd/message-service-kata -service consumer` or `make serve-consumer`

## Structure Table using PostgreSQL
# Create Table
`CREATE TABLE consumed_messages (
    id SERIAL PRIMARY KEY,
    message JSONB NOT NULL,
    trigger_by VARCHAR(255),
    received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`

## CURL Example
- Kafka Producer : `curl --location 'http://localhost:8089/v1/message/post' \
--header 'Content-Type: application/json' \
--data '{
    "trigger_by": "try",
    "qty": 2
}'`

- Health Check : `curl --location 'http://localhost:8089/v1/message/health'`