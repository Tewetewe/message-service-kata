package infra

import (
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
