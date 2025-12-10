package database

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/shanto-323/chat-ai/config"
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
	CreateUser(ctx context.Context, user *dto.CreateUser) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

func New(cfg *config.Config, logger *zerolog.Logger) (Database, error) {
	return nil, nil
}
