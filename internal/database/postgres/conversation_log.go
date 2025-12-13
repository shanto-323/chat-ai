package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shanto-323/chat-ai/internal/server/errs"
	"github.com/shanto-323/chat-ai/model"
	"github.com/shanto-323/chat-ai/model/dto"
	"github.com/shanto-323/chat-ai/model/entity"
)

func (db *DB) CreateConversationLog(ctx context.Context, cl *entity.ConversationLog) (*entity.ConversationLog, error) {
	query := `
		INSERT INTO conversation_logs (
			user_id,
			text_query,
			image_url,
			response_text,
			llm_model_name,
			vlm_model_name
		)
		VALUES (
			@user_id,
			@text_query,
			@image_url,
			@response_text,
			@llm_model_name,
			@vlm_model_name
		)	
		RETURNING 
			* 
	`

	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"user_id":        cl.UserID,
		"text_query":     cl.TextQuery,
		"image_url":      cl.ImageURL,
		"response_text":  cl.ResponseText,
		"llm_model_name": cl.LLMModelName,
		"vlm_model_name": cl.VLMModelName,
	}).Scan(
		&cl.ID,
		&cl.UserID,
		&cl.TextQuery,
		&cl.ImageURL,
		&cl.ResponseText,
		&cl.LLMModelName,
		&cl.VLMModelName,
		&cl.Timestamp,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.NewInternalServerError()
		}
		return nil, err
	}

	return cl, nil
}

func (db *DB) GetConversationLogHistory(
	ctx context.Context,
	userId uuid.UUID,
	queryDto *dto.ConversationHistoryQuery,
) (*model.PaginatedResponse[entity.ConversationLog], error) {

	x := *queryDto.Page
	y := *queryDto.Limit
	offset := (x - 1) * y

	query := fmt.Sprintf(`
		SELECT 
			*
		FROM 
			conversation_logs
		WHERE
			user_id=@user_id
		ORDER BY timestamp %s
		LIMIT @limit
		OFFSET @offset
	`, strings.ToUpper(*queryDto.Order))

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"user_id": userId,
		"order":   *queryDto.Order,
		"limit":   y,
		"offset":  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get query")
	}

	logs, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.ConversationLog])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[entity.ConversationLog]{
				Data:       []entity.ConversationLog{},
				Page:       *queryDto.Page,
				Limit:      *queryDto.Limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to collect rows")
	}

	// total count
	count := `
		SELECT
			COUNT(*)
		FROM
			conversation_logs	
		WHERE
			user_id=@user_id
	`

	countArgs := pgx.NamedArgs{
		"user_id": userId,
	}

	var total int
	err = db.pool.QueryRow(ctx, count, countArgs).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count")
	}

	db.logger.Info().Str("event", "chat-history").Str("user_id", userId.String()).Int("total", total).Int("offset", offset).Msg("chat history found")

	return &model.PaginatedResponse[entity.ConversationLog]{
		Data:       logs,
		Page:       *queryDto.Page,
		Limit:      *queryDto.Limit,
		Total:      total,
		TotalPages: (total + *queryDto.Limit - 1) / *queryDto.Limit,
	}, nil
}
