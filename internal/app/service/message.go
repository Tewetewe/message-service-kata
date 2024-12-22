package service

//go:generate mockery --dir=$PROJECT_DIR/internal/app/service  --name=MessageSvc --filename=$GOFILE --output=$PROJECT_DIR/internal/generated/mock_service --outpkg=mock_service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"message-service-kata/internal/app/repo/kafka"
	"message-service-kata/internal/app/repo/postgres"
	"message-service-kata/pkg/domain/entities"

	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

type (
	// MessageSvc interfacing message service function
	MessageSvc interface {
		PostMessage(ctx context.Context, args *entities.CreateMessageRequest) (err error)
		ProcessMessage(ctx context.Context, args entities.MessageData) (err error)
	}

	// MessageSvcImpl implementing message service dependencies
	MessageSvcImpl struct {
		dig.In
		MessageRepo postgres.MessageRepository
		KafkaRepo   kafka.RepositoryKafka
	}
)

// NewMessageSvc initiating message service
func NewMessageSvc(impl MessageSvcImpl) MessageSvc {
	return &impl
}

// PostMessage service to produce message
func (s *MessageSvcImpl) PostMessage(
	ctx context.Context, args *entities.CreateMessageRequest,
) (err error) {
	log.Info().Msgf("[MessageSvc][PostMessage] incoming request with arg: %v", args)

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	if args.Qty > 0 {
		// Loop with qty on argument
		for i := 0; i < int(args.Qty); i++ {
			for j := 0; j < len(entities.Queries); j++ {
				wg.Add(1)

				// Passing `queries[j]` as a parameter to the goroutine to avoid closure issues
				go func(query string) {
					defer wg.Done()

					// Recover from any panic inside the goroutine
					defer func() {
						if r := recover(); r != nil {
							// Log the panic recovery
							log.Printf("Panic recovered in goroutine: %v", r)
						}
					}()

					// Assign value to model
					message := entities.MessageData{
						TriggerBy: args.TriggerBy,
						Message:   query, // Use the passed query variable
					}

					// Publish message using Kafka
					err := s.KafkaRepo.PublishWithoutKey(ctx, kafka.PublishData{
						Topic: string(entities.TopicPublishMessage),
						Data:  message,
					})
					if err != nil {
						log.Error().Msgf("Error publishing message: %v", err)
						return
					}

					log.Info().Msgf("[MessageSvc][PostMessage][PublishWithoutKey] success publish message with data: %v", message)
				}(entities.Queries[j]) // Pass `queries[j]` as a parameter
			}
		}
	}
	wg.Wait() // Wait for all goroutines to finish

	log.Info().Msgf("[MessageSvc][PostMessage] finish processing all message with arg: %v", args)

	return nil
}

// ProcessMessage service to process message
func (s *MessageSvcImpl) ProcessMessage(
	ctx context.Context, args entities.MessageData,
) (err error) {
	log.Info().Msgf("[MessageSvc][ProcessMessage] incoming request with arg: %v", args)

	// Generate a response
	responseMessage := generateResponse(args.Message)
	log.Info().Msgf("[MessageSvc][ProcessMessage] reply request to : %v", responseMessage)

	// Prepare the JSON object for storage
	data := map[string]interface{}{
		"received_message": args.Message,
		"response_message": responseMessage,
	}

	// Using go routine to make multithread processing
	go func() {
		err = s.storeConsumedMessageAsJSON(ctx, args.TriggerBy, data)
		if err != nil {
			log.Error().Msgf("[MessageSvc][ProcessMessage] error while storeConsumedMessageAsJSON : %v", err)
		}

		log.Info().Msgf("[MessageSvc][ProcessMessage] finish processing all message with data: %v", data)
	}()

	return nil
}

// generateResponse generates a response based on the received message
func generateResponse(message string) string {
	if response, found := entities.Responses[message]; found {
		return response
	}
	return string(entities.FallbackResponse)
}

// storeConsumedMessageAsJSON saves the received message and response to PostgreSQL as JSONB
func (s *MessageSvcImpl) storeConsumedMessageAsJSON(ctx context.Context, triggerBy string, data map[string]interface{}) (err error) {
	// Convert the map to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		err = fmt.Errorf("failed to marshal data to JSON: %w", err)
		log.Error().Msgf("[MessageSvc][ProcessMessage] error while Marshal jsonData : %v", err)
		return err
	}

	message := &entities.MessageData{
		TriggerBy: triggerBy,
		Message:   string(jsonData),
	}

	_, err = s.MessageRepo.Create(ctx, message)
	if err != nil {
		log.Error().Msgf("[MessageSvc][ProcessMessage] error while Create Data in postgre : %v", err)
		return err
	}

	return nil
}
