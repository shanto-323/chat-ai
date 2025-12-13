package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/chat-ai/internal/service"
	"github.com/shanto-323/chat-ai/model"
	"github.com/shanto-323/chat-ai/model/dto"
	"github.com/shanto-323/chat-ai/model/entity"
)

type ChatHandler struct {
	service *service.ChatService
}

func NewChatHandler(s *service.ChatService) *ChatHandler {
	return &ChatHandler{
		service: s,
	}
}

func (h *ChatHandler) ChatHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return Handle(
			func(c echo.Context, req *dto.ChatRequest) (*entity.ConversationLog, error) {
				return h.service.MultimodalChat(c, req)
			},
			http.StatusOK,
			&dto.ChatRequest{},
		)(c)
	}
}

func (h *ChatHandler) ChatHistoryHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return Handle(
			func(c echo.Context, req *dto.ConversationHistoryQuery) (*model.PaginatedResponse[entity.ConversationLog], error) {
				return h.service.MultimodalChatHistory(c, req)
			},
			http.StatusOK,
			&dto.ConversationHistoryQuery{},
		)(c)
	}
}
