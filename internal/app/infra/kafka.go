package infra

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type (
	// KafkaCfg used to load kafka config from .env
	KafkaCfg struct {
		BrokerAddress      string `envconfig:"BROKER_ADDR" required:"true" default:"127.0.0.1:9092"`
		GroupID            string `envconfig:"GROUP_ID" required:"true" default:"message-consumer-group"`
		MaxConsumerRetries int    `envconfig:"MAX_CONSUMER_RETRIES" required:"true" default:"3"`
	}
)

// NewConsumer used to connect  to Kafka consumer instance
func NewConsumer(cfg *KafkaCfg) *kafka.Consumer {
	log.Info().Msg(cfg.GroupID)
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":             cfg.BrokerAddress,
		"group.id":                      cfg.GroupID,
		"partition.assignment.strategy": "roundrobin",
		"auto.offset.reset":             "earliest",
		"enable.auto.commit":            false,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create consumer")
	}

	return c
}

// NewProducer used to connect  to Kafka producer instance
func NewProducer(cfg *KafkaCfg) *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.BrokerAddress,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create producer")
	}

	return p
}

// Example producer connection check
func CheckProducerConnection(cfg *KafkaCfg) error {
	// Initialize Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.BrokerAddress,
	})
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}
	defer p.Close()

	testTopic := "test"
	// Produce a dummy message to check connection
	dummyMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &testTopic, // Specify a test topic
			Partition: kafka.PartitionAny,
		},
		Value: []byte("test"), // Test message content
	}

	// Attempt to send the message and wait for the result
	deliveryChan := make(chan kafka.Event)
	err = p.Produce(dummyMessage, deliveryChan)
	if err != nil {
		return fmt.Errorf("failed to send test message: %w", err)
	}

	// Wait for the delivery report
	select {
	case event := <-deliveryChan:
		switch e := event.(type) {
		case *kafka.Message:
			if e.TopicPartition.Error != nil {
				return fmt.Errorf("message delivery failed: %w", e.TopicPartition.Error)
			}
		}
	}

	// If no errors occurred, the connection is good
	return nil
}
