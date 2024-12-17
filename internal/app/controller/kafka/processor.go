package kafka

import (
	"context"
	"encoding/json"

	"message-service-kata/internal/app/service"
	"message-service-kata/pkg/domain/entities"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

type (
	// ProcessorImpl Process a set of kafka processors.
	ProcessorImpl struct {
		dig.In
		MessageSvc service.MessageSvc
	}

	// Processor implementator for processing messages.
	Processor interface {
		ProcessMessage(ctx context.Context, message *kafka.Message) (err error)
	}
)

// NewProcessor creates a new instance of the KafkaProcessor
func NewProcessor(impl ProcessorImpl) Processor {
	return &impl
}

// ProcessMessage impelements interface processor
func (op *ProcessorImpl) ProcessMessage(ctx context.Context, message *kafka.Message) (err error) {
	defer func() {
		if err != nil {
			log.Error().Msgf("[ProcessMessage] any error with msg : %v", err)
		}
	}()

	var data entities.MessageData

	err = json.Unmarshal(message.Value, &data)
	if err != nil {
		return err
	}

	err = op.MessageSvc.ProcessMessage(ctx, data)
	if err != nil {
		return err
	}

	return nil
}
