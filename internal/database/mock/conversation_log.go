package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/shanto-323/chat-ai/model"
	"github.com/shanto-323/chat-ai/model/dto"
	"github.com/shanto-323/chat-ai/model/entity"
)

func (db *DB) CreateConversationLog(ctx context.Context, cl *entity.ConversationLog) error {
	return nil
}

func (db *DB) GetConversationLogHistory(ctx context.Context, userId uuid.UUID, query *dto.ConversationHistoryQuery) (*model.PaginatedResponse[dto.ConversationLogResponse], error) {
	return nil , nil
}
