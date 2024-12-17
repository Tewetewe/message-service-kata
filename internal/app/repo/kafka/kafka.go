package kafka

//go:generate mockery --dir=$PROJECT_DIR/internal/app/repo/kafka  --name=KafkaRepository --filename=$GOFILE --output=$PROJECT_DIR/internal/generated/mock_kafka --outpkg=mock_kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/dig"
)

type (
	// PublishData used to Create kafka publish message
	PublishData struct {
		Topic string
		Key   string
		Data  interface{}
	}

	// RepositoryKafkaImpl implementing kafka producer
	RepositoryKafkaImpl struct {
		dig.In
		KafkaProduce *kafka.Producer
	}

	// RepositoryKafka interfacing kafka repository function
	RepositoryKafka interface {
		PublishWithKey(ctx context.Context, args PublishData) (err error)
		PublishWithoutKey(ctx context.Context, args PublishData) (err error)
	}
)

// NewKafkaRepository Initiating Kafka to implement producer
func NewKafkaRepository(impl RepositoryKafkaImpl) RepositoryKafka {
	return &impl
}

// PublishWithKey function to publish kafka message using Key
func (ox *RepositoryKafkaImpl) PublishWithKey(ctx context.Context, args PublishData) (err error) {
	byt, err := json.Marshal(args.Data)
	if err != nil {
		return err
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &args.Topic,
			Partition: kafka.PartitionAny,
		},
		Value: byt,
		Key:   []byte(args.Key),
	}

	deliverChan := make(chan kafka.Event)
	err = ox.KafkaProduce.Produce(msg, deliverChan)
	if err != nil {
		err = fmt.Errorf("[repository][PublishWithKey] while producing : %v", err)
	}

	kafkEvent := <-deliverChan
	if err != nil {
		return err
	}

	msg = kafkEvent.(*kafka.Message)
	if msg.TopicPartition.Error != nil {
		err = msg.TopicPartition.Error
		return err
	}

	return nil
}

// PublishWithoutKey function to publish kafka message without Key
func (ox *RepositoryKafkaImpl) PublishWithoutKey(ctx context.Context, args PublishData) (err error) {
	byt, err := json.Marshal(args.Data)
	if err != nil {
		return err
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &args.Topic,
			Partition: kafka.PartitionAny,
		},
		Value: byt,
	}

	deliverChan := make(chan kafka.Event)
	err = ox.KafkaProduce.Produce(msg, deliverChan)
	if err != nil {
		err = fmt.Errorf("[repository][PublishWithoutKey] while producing : %v", err)
		deliverChan <- nil
	}

	kafkEvent := <-deliverChan
	if err != nil {
		return err
	}

	msg = kafkEvent.(*kafka.Message)
	if msg.TopicPartition.Error != nil {
		err = msg.TopicPartition.Error
		return err
	}
	return nil
}
