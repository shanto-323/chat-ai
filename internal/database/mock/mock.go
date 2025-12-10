package mock

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/shanto-323/chat-ai/config"
	"github.com/shanto-323/chat-ai/internal/database"
)

type DB struct {
	logger *zerolog.Logger
}

func New(config *config.Config, logger *zerolog.Logger) (database.Database, error) {
	return &DB{
		logger: logger,
	}, nil
}

func (db *DB) Ping(ctx context.Context) error         { return nil }
func (db *DB) IsInitialized(ctx context.Context) bool { return true }
func (db *DB) Close() error                           { return nil }
