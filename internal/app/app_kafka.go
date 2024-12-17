package app

import (
	"context"
	"fmt"
	kafkaCtrl "message-service-kata/internal/app/controller/kafka"
	"message-service-kata/internal/app/infra"
	"os"

	"message-service-kata/pkg/domain/entities"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

type (
	// ConsumerHandlerParams is a consumer handler dependencies
	ConsumerHandlerParams struct {
		dig.In
		Consumer  *kafka.Consumer
		Producer  *kafka.Producer
		KafkaCfg  *infra.KafkaCfg
		KafkaCtrl kafkaCtrl.Processor
	}
)

// Topics is a variable for list of kafka topic.
var Topics = []string{
	string(entities.TopicPublishMessage),
}

func startConsumer(
	args ConsumerHandlerParams,
	shutdownCh <-chan struct{},
	errCh chan<- error,
	topic string,
) {
	topics := Topics
	if topic != "" {
		topics = []string{topic}
	}

	err := args.Consumer.SubscribeTopics(topics, rebalanceCallback)
	if err != nil {
		log.Error().Msgf("SubscribeTopics: %s", err.Error())
		errCh <- err // send error to error channel
		return
	}

	retryCount := make(map[kafka.TopicPartition]int)
	for {
		select {
		case <-shutdownCh:
			log.Info().Msg("shutdown consumer")
			return
		default:
			msg, err := args.Consumer.ReadMessage(-1)
			if err != nil {
				log.Error().Msgf("ReadMessage: %s", err.Error())
				errCh <- err
				return
			}

			isRetryProcessMessage := true
			for isRetryProcessMessage {
				err = handleMessage(msg, args)
				if err != nil {
					log.Error().Any("topic", msg.TopicPartition).Any("value", string(msg.Value)).Any("error", err).Msg("error process kafka message")

					// If the error is retryable and retry count is less than maxRetryCount, don't commit the offset.
					if retryCount[msg.TopicPartition] < args.KafkaCfg.MaxConsumerRetries {
						log.Warn().
							Any("topic", msg.TopicPartition).
							Any("value", string(msg.Value)).
							Any("retry count", retryCount[msg.TopicPartition]).
							Msg("Kafka retry")

						retryCount[msg.TopicPartition]++
						continue
					}

					// If the error is not retryable or retry count has reached maxRetryCount,
					// forward the message to the dead-letter queue and commit the offset.
					// Build the DLQ topic name
					dlqTopic := fmt.Sprintf("%s-dead-letter-queue", *msg.TopicPartition.Topic)
					err = args.Producer.Produce(&kafka.Message{
						TopicPartition: kafka.TopicPartition{Topic: &dlqTopic, Partition: kafka.PartitionAny},
						Value:          msg.Value,
						Key:            msg.Key,
					}, nil)
					if err != nil {
						log.Error().Any("topic", dlqTopic).Any("value", string(msg.Value)).Any("error", err).Msg("error process produce dlq message")
					}
				}

				// If the message has been consumed or forwarded to the dead-letter queue, commit the offset.
				_, err = args.Consumer.CommitMessage(msg)
				if err != nil {
					log.Error().Msgf("Failed to commit message: %s", err.Error())
				}

				// Reset retry count for the message
				retryCount[msg.TopicPartition] = 0

				// delete key
				delete(retryCount, msg.TopicPartition)

				// stop retry process
				isRetryProcessMessage = false
			}
		}
	}
}

func handleMessage(
	msg *kafka.Message,
	args ConsumerHandlerParams,
) (err error) {
	var topic string

	ctx := context.Background()

	if msg.TopicPartition.Topic != nil {
		topic = *msg.TopicPartition.Topic
	}

	switch topic {
	case string(entities.TopicPublishMessage):
		err = args.KafkaCtrl.ProcessMessage(ctx, msg)
	default:
		err = nil
	}

	if err != nil {
		log.Error().Any("topic", topic).Any("value", string(msg.Value)).Any("error", err.Error()).Msg("error handle kafka message")
		return err
	}

	return nil
}

func rebalanceCallback(c *kafka.Consumer, event kafka.Event) error {
	hostName, err := os.Hostname()
	if err != nil {
		log.Error().Msgf("[rebalanceCallback] error while read Hostname: %s", err.Error())
	}

	switch ev := event.(type) {
	case kafka.AssignedPartitions:
		if len(ev.Partitions) > 0 {
			log.Info().Msgf("[rebalanceCallback] %s assigned partitions: %v", hostName, ev.Partitions)
		}
	case kafka.RevokedPartitions:
		if len(ev.Partitions) > 0 {
			log.Info().Msgf("[rebalanceCallback] %s revoked partitions: %v", hostName, ev.Partitions)
		}
	}

	return nil
}
