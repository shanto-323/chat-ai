package postgres

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/shanto-323/chat-ai/config"
)

type DB struct {
	pool   *pgxpool.Pool
	logger *zerolog.Logger
}

func New(cfg *config.Config, logger *zerolog.Logger) (*DB, error) {
	hostPort := net.JoinHostPort(cfg.Database.Host, strconv.Itoa(cfg.Database.Port))

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		hostPort,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	pgxPoolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx pool cfg: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	logger.Info().Msg("postgres service initialized successfully")

	return &DB{
		pool:   pool,
		logger: logger,
	}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *DB) IsInitialized(ctx context.Context) bool {
	return db.pool != nil
}

func (db *DB) Close() error {
	db.logger.Info().Msg("closing database connection pool")
	db.pool.Close()
	return nil
}
