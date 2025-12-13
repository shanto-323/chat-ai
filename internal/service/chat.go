package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shanto-323/chat-ai/internal/database"
	"github.com/shanto-323/chat-ai/internal/server/errs"
	"github.com/shanto-323/chat-ai/internal/server/manager"
	"github.com/shanto-323/chat-ai/internal/service/image"
	"github.com/shanto-323/chat-ai/model"
	"github.com/shanto-323/chat-ai/model/dto"
	"github.com/shanto-323/chat-ai/model/entity"
)

type ChatService interface {
	MultimodalChat(c echo.Context, payload *dto.ChatRequest) (*entity.ConversationLog, error)
	MultimodalChatHistory(c echo.Context, payload *dto.ConversationHistoryQuery) (*model.PaginatedResponse[entity.ConversationLog], error)
}

type chatService struct {
	is *image.ImageService

	manager *manager.AIManager
	db      database.Database
}

func NewChatService(manager *manager.AIManager, db database.Database, is *image.ImageService) ChatService {
	return &chatService{
		manager: manager,
		db:      db,
		is:      is,
	}
}

func (s *chatService) MultimodalChat(c echo.Context, payload *dto.ChatRequest) (*entity.ConversationLog, error) {
	images, err := s.is.ProcessImage(payload.Images)
	if err != nil {
		return nil, err
	}

	// With a real vlm service we can get structerd info about images.
	// For now this provide nothing as but can be implemented.
	_, err = s.manager.VLMManager.AnalyzeImage(images)

	// Lets imagine we got our info and we insered that in processedText.
	processedText := payload.UserMessage

	ctx, cancel := context.WithTimeout(c.Request().Context(), 60*time.Second)
	defer cancel()

	response, err := s.manager.LLMManager.GenerateResponse(ctx, &dto.LLMRequest{Messages: processedText, Model: payload.ModelConfig})
	if err != nil {
		return nil, err
	}

	userId, ok := c.Get("id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	conversationLog := &entity.ConversationLog{
		UserID:       userId,
		TextQuery:    payload.UserMessage,
		ImageURL:     make([]string, 0),
		ResponseText: response,
	}

	conversationLog.LLMModelName = payload.ModelConfig.LLMModel
	conversationLog.VLMModelName = payload.ModelConfig.VLMModel

	return s.db.CreateConversationLog(ctx, conversationLog)
}

func (s *chatService) MultimodalChatHistory(c echo.Context, payload *dto.ConversationHistoryQuery) (*model.PaginatedResponse[entity.ConversationLog], error) {
	userId, ok := c.Get("id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	return s.db.GetConversationLogHistory(ctx, userId, payload)
}
