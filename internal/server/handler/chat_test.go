package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/chat-ai/model"
	"github.com/shanto-323/chat-ai/model/dto"
	"github.com/shanto-323/chat-ai/model/entity"
)

type mockChatService struct {
	chatFn        func(c echo.Context, payload *dto.ChatRequest) (*entity.ConversationLog, error)
	chatHistoryFn func(c echo.Context, payload *dto.ConversationHistoryQuery) (*model.PaginatedResponse[entity.ConversationLog], error)
}

func (m *mockChatService) MultimodalChat(
	c echo.Context,
	payload *dto.ChatRequest,
) (*entity.ConversationLog, error) {
	return m.chatFn(c, payload)
}

func (m *mockChatService) MultimodalChatHistory(
	c echo.Context,
	payload *dto.ConversationHistoryQuery,
) (*model.PaginatedResponse[entity.ConversationLog], error) {
	return m.chatHistoryFn(c, payload)
}

func newTestContext(method, path string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()

	req := httptest.NewRequest(method, path, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}

func TestChatHandler(t *testing.T) {
	mockSvc := &mockChatService{
		chatFn: func(c echo.Context, payload *dto.ChatRequest) (*entity.ConversationLog, error) {
			// basic assertions
			if payload.UserMessage != "hello" {
				t.Fatalf("unexpected message: %s", payload.UserMessage)
			}

			if payload.ModelConfig == nil || payload.ModelConfig.LLMModel != "gpt" {
				t.Fatal("model config not bound correctly")
			}

			return &entity.ConversationLog{
				TextQuery:    payload.UserMessage,
				ResponseText: "hi there",
			}, nil
		},
	}

	handler := NewChatHandler(mockSvc)

	body := `{
		"message_query": "hello",
		"model_config": {
			"llm_model": "gpt",
			"vlm_model": "llava"
		}
	}`

	c, rec := newTestContext(http.MethodPost, "/chat", strings.NewReader(body))

	err := handler.ChatHandler()(c)

	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if rec.Body.Len() == 0 {
		t.Fatal("expected response body")
	}
}

func TestChatHistoryHandler(t *testing.T) {
	mockSvc := &mockChatService{
		chatHistoryFn: func(c echo.Context, payload *dto.ConversationHistoryQuery) (*model.PaginatedResponse[entity.ConversationLog], error) {
			if payload.Page == nil || *payload.Page != 1 {
				t.Fatal("page not parsed correctly")
			}

			if payload.Limit == nil || *payload.Limit != 10 {
				t.Fatal("limit not parsed correctly")
			}

			if payload.Order == nil || *payload.Order != "desc" {
				t.Fatal("order not parsed correctly")
			}

			return &model.PaginatedResponse[entity.ConversationLog]{
				Data:  []entity.ConversationLog{},
				Page:  1,
				Limit: 10,
				Total: 0,
			}, nil
		},
	}

	handler := NewChatHandler(mockSvc)

	c, rec := newTestContext(
		http.MethodGet,
		"/chat/history?page=1&limit=10&order=desc",
		nil,
	)

	err := handler.ChatHistoryHandler()(c)

	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
