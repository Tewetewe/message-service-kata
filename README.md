# Message Service Kata

---

## Prerequisites

1. **Kafka Installed**  
   Make sure Apache Kafka is installed and running on your system. If you haven't set up Kafka yet, follow these steps:
   - Download Kafka: [https://kafka.apache.org/downloads](https://kafka.apache.org/downloads)
   - Extract the archive and navigate to the Kafka folder in your terminal.
   - Start ZooKeeper:
     ```bash
     bin/zookeeper-server-start.sh config/zookeeper.properties
     ```
   - Start Kafka Broker:
     ```bash
     bin/kafka-server-start.sh config/server.properties
     ```

2. **Go Installed**  
   Install Go from [https://golang.org/dl/](https://golang.org/dl/) and set up your environment.

3. **PostgreSQL**
   Ensure PostgreSQL is installed and running.

---

## Setup

1. Clone this repository or copy the code into a new directory:
   ```bash
   git clone https://github.com/Tewetewe/message-service-kata
   cd message-service-kata
   ```

2. Install dependencies:
   ```bash
   go get github.com/segmentio/kafka-go
   go get github.com/lib/pq
   ```

3. Configure the required environment variables:
   ```bash
   export APP_DEBUG=true
   export APP_ADDRESS=localhost:8089
   export PG_DBNAME=yourdbname
   export PG_HOST=localhost
   export PG_PORT=5432
   export PG_DBUSER=yourdbuser
   export PG_DBPASS=yourdbpassword
   export PG_SSL_MODE=disable
   export KAFKA_BROKER_ADDR=localhost:9092
   ```

4. Create the Kafka topic (if it doesn't exist):
   ```bash
   bin/kafka-topics.sh --create --topic message.publish --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
   ```

5. Set up the PostgreSQL table:
   ```sql
   CREATE TABLE consumed_messages (
       id SERIAL PRIMARY KEY,
       message JSONB NOT NULL,
       trigger_by VARCHAR(255),
       received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```

---

## Running the Service

1. Start the API server:
   ```bash
   go run ./cmd/message-service-kata -service rest
   ```
   Or using Make:
   ```bash
   make serve-rest
   ```

2. Start the consumer:
   ```bash
   go run ./cmd/message-service-kata -service consumer
   ```
   Or using Make:
   ```bash
   make serve-consumer
   ```

---

## CURL Examples

### Trigger Kafka Producer:
```bash
curl --location 'http://localhost:8089/v1/message/post' \
--header 'Content-Type: application/json' \
--data '{
    "trigger_by": "try",
    "qty": 2
}'
```

### Health Check:
```bash
curl --location 'http://localhost:8089/v1/message/health'
```

---

## Troubleshooting

- **Kafka connection issues**: Ensure Kafka is running on the specified broker address (`localhost:9092`).
- **Library issues**: Verify that the `kafka-go` library is installed correctly by running:
  ```bash
  go mod tidy
  ```

---


## Process

- **Trigger Producer**: Produce message by hit endpoint /post with qty by request. Every 1 qty will produce this queries. This message will be produce concurrently.
  ```
      var Queries = []string{
      "Hello",
      "Weather update",
      "Tell me a joke",
      "Good morning",
      "What's your name?",
      "How are you?",
    }
  ```

  Sample log info when success produce message to kafka:
  ```
  2024-12-22 21:03:29 INF [MessageSvc][PostMessage][PublishWithoutKey] success publish message with data: {Weather update try}
  ```

  Sample message on Kafka with topic `message.publish`:
  ```
  {
    "message": "Weather update",
    "trigger_by": "try"
  }
  ```

- **Consume Message**: Consume message by queue on kafka, the process devide mapping right response and store data to postgre. On storing postgre, I use goroutine and not mandatory to wait so the consumer will be process next message on queue.

  Sample log info when success consume message to kafka:
  ```
  2024-12-22 21:03:29 INF [MessageSvc][ProcessMessage] finish processing all message with data: map[received_message:Tell me a joke response_message:Why did the chicken cross the road? To get to the other side! 😂]
  ```

  Example data stored on database
  ```
  19	{"received_message": "Hello", "response_message": "Hi there! 😊"}	try	2024-12-18 04:32:04.419
  11	{"received_message": "What's your name?", "response_message": "I'm sorry, I didn't understand that. 🤔"}	try	2024-12-18 04:29:06.722
  ```
---
