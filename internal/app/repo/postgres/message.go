package postgres

//go:generate mockery --dir=$PROJECT_DIR/internal/app/repo/postgres  --name=MessageRepository --filename=$GOFILE --output=$PROJECT_DIR/internal/generated/mock_postgres --outpkg=mock_postgres
import (
	"context"
	"database/sql"

	"message-service-kata/internal/app/repo/postgres/queries"

	"github.com/rs/zerolog/log"
	"go.uber.org/dig"

	"message-service-kata/pkg/domain/entities"
)

type (
	// MessageRepositoryImpl Implementing message repository dependency
	MessageRepositoryImpl struct {
		dig.In
		*sql.DB
	}

	// MessageRepository interfacing Message Repository function
	MessageRepository interface {
		// create
		Create(ctx context.Context, args *entities.MessageData) (messageID int64, err error)
	}
)

// NewMessageRepository initiate message repository
func NewMessageRepository(impl MessageRepositoryImpl) MessageRepository {
	return &impl
}

// Create - function for store conversation message
func (r *MessageRepositoryImpl) Create(ctx context.Context, args *entities.MessageData) (messageID int64, err error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil && tx != nil {
			errs := tx.Rollback()
			if errs != nil {
				log.Error().Any("error", errs).Msg("error process rollback")
			}
		}
	}()

	err = tx.QueryRowContext(
		ctx,
		queries.QueryCreateMessage,
		args.Message,
		args.TriggerBy,
	).Scan(&messageID)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return messageID, nil
}
