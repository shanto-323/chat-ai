package service

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/chat-ai/internal/database"
	"github.com/shanto-323/chat-ai/model/dto"
)

type AuthService struct {
	db database.Database
}

func NewAuthRepository(db database.Database) *AuthService {
	return &AuthService{
		db: db,
	}
}

func (a *AuthService) Login(c *echo.Context, payload *dto.CreateUser) (*dto.UserAuthResponse, error) {
	user, err := a.db.GetUserByEmail(context.Background(), payload.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("not exist")
	}

	// match password

	token := "new-token"

	return &dto.UserAuthResponse{
		AccessToken: token,
	}, nil
}

func (a *AuthService) Register(c *echo.Context, payload *dto.CreateUser) (*dto.UserAuthResponse, error) {
	user, err := a.db.GetUserByEmail(context.Background(), payload.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, fmt.Errorf("exist")
	}

	hashedPassword := "hash-password"
	payload.Password = hashedPassword

	user, err = a.db.CreateUser(context.Background(), payload)
	if err != nil {
		return nil, err
	}

	token := "new-token"

	return &dto.UserAuthResponse{
		AccessToken: token,
	}, nil
}
