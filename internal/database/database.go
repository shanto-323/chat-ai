package database

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/shanto-323/chat-ai/config"
	"github.com/shanto-323/chat-ai/internal/database/mock"
	"github.com/shanto-323/chat-ai/internal/database/postgres"
	"github.com/shanto-323/chat-ai/model/dto"
	"github.com/shanto-323/chat-ai/model/entity"
)

// It contains all methods that database should implement.
type Database interface {
	// Database specific methods
	Ping(ctx context.Context) error
	IsInitialized(ctx context.Context) bool
	Close() error

	// Other methods related to database operation
	CreateUser(ctx context.Context, user *dto.RegisterRequest) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	CreateConversationLog(ctx context.Context, cl *entity.ConversationLog) (*entity.ConversationLog , error)
}

func New(cfg *config.Config, logger *zerolog.Logger) (Database, error) {
	switch cfg.Primary.DatabaseType {
	case "postgres":
		postgres.New(cfg, logger)
		// return this
		return nil, nil
	case "mock":
		return mock.New(cfg, logger)
	default:
		return nil, fmt.Errorf("no database found")
	}
}
