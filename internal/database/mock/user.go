package mock

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shanto-323/chat-ai/model/dto"
	"github.com/shanto-323/chat-ai/model/entity"
)

func (db *DB) CreateUser(ctx context.Context, user *dto.CreateUser)(*entity.User, error) {
	return nil ,nil
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	entity := entity.User{
		Email:        email,
		PasswordHash: "some-hash",
	}

	entity.ID = uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()

	return &entity, nil
}
