package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/chat-ai/internal/service"
	"github.com/shanto-323/chat-ai/model/dto"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

func (h *AuthHandler) LoginHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return HandleNoResponse(
			func(c echo.Context, req *dto.LoginRequest) error {
				resp, err := h.service.Login(c, req)
				if err != nil {
					return err
				}
				c.SetCookie(&http.Cookie{
					Name:  "access_token",
					Value: resp.AccessToken,
					Path:  "/",
				})

				return nil
			},
			http.StatusOK,
			&dto.LoginRequest{},
		)(c)
	}
}

func (h *AuthHandler) RegisterHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return HandleNoResponse(
			func(c echo.Context, req *dto.RegisterRequest) error {
				resp, err := h.service.Register(c, req)
				if err != nil {
					return err
				}
				c.SetCookie(&http.Cookie{
					Name:  "access_token",
					Value: resp.AccessToken,
					Path:  "/",
				})

				return nil
			},
			http.StatusCreated,
			&dto.RegisterRequest{},
		)(c)
	}
}
