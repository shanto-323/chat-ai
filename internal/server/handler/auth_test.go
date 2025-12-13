package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/chat-ai/model/dto"
)

type mockAuthService struct {
	loginFn func(c echo.Context, payload *dto.LoginRequest) (*dto.AuthResponse, error)
}

func (m *mockAuthService) Login(
	c echo.Context,
	payload *dto.LoginRequest,
) (*dto.AuthResponse, error) {
	return m.loginFn(c, payload)
}

func (m *mockAuthService) Register(
	c echo.Context,
	payload *dto.RegisterRequest,
) (*dto.AuthResponse, error) {
	return nil, nil
}

func setupEchoContext(method, path string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()

	req := httptest.NewRequest(method, path, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}

func TestLoginHandler(t *testing.T) {
	mockSvc := &mockAuthService{
		loginFn: func(c echo.Context, payload *dto.LoginRequest) (*dto.AuthResponse, error) {
			if payload.Email != "test@example.com" {
				t.Fatalf("unexpected email: %s", payload.Email)
			}
			if payload.Password != "123456" {
				t.Fatalf("unexpected password")
			}

			return &dto.AuthResponse{
				AccessToken: "mock-token",
			}, nil
		},
	}

	handler := &AuthHandler{
		service: mockSvc,
	}

	body := `{
		"email": "test@example.com",
		"password": "123456"
	}`

	c, rec := setupEchoContext(http.MethodPost, "/login", strings.NewReader(body))

	err := handler.LoginHandler()(c)

	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	cookies := rec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected cookie to be set")
	}

	if cookies[0].Name != "access_token" {
		t.Fatalf("unexpected cookie name: %s", cookies[0].Name)
	}

	if cookies[0].Value != "mock-token" {
		t.Fatalf("unexpected token value")
	}
}
