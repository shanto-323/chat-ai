package service

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/chat-ai/config"
	"github.com/shanto-323/chat-ai/internal/database"
	"github.com/shanto-323/chat-ai/internal/server/errs"
	"github.com/shanto-323/chat-ai/internal/server/middleware"
	"github.com/shanto-323/chat-ai/model/dto"
	"github.com/shanto-323/chat-ai/pkg"
)

type AuthService struct {
	cfg *config.Config
	db  database.Database
}

func NewAuthService(cfg *config.Config, db database.Database) *AuthService {
	return &AuthService{
		cfg: cfg,
		db:  db,
	}
}

func (a *AuthService) Login(c echo.Context, payload *dto.LoginRequest) (*dto.AuthResponse, error) {
	logger := middleware.GetLogger(c)

	user, err := a.db.GetUserByEmail(context.Background(), payload.Email)
	if err != nil {
		logger.Error().Err(err).Msg("login failed")
		return nil, err
	}

	if pkg.CompareWithHash(user.PasswordHash, payload.Password) != nil {
		logger.Warn().Msg("login failed, wrong password")
		return nil, errs.NewForbiddenError(
			"Invalid credentials",
			true,
		)
	}

	token, err := pkg.CreateAccessToken(a.cfg, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("create token failed")
		return nil, errs.NewInternalServerError()
	}

	logger.Info().
		Str("event", "login").
		Msg("logged in successfully")

	return &dto.AuthResponse{
		AccessToken: token,
	}, nil
}

func (a *AuthService) Register(c echo.Context, payload *dto.RegisterRequest) (*dto.AuthResponse, error) {
	logger := middleware.GetLogger(c)

	user, err := a.db.GetUserByEmail(context.Background(), payload.Email)
	if err == nil && user != nil { // this is mock
		logger.Warn().Msg("register failed, user exists")
		return nil, err
	}

	hashPassword, err := pkg.CreateHash(payload.Password)
	if err != nil {
		logger.Error().Err(err).Msg("create token failed")
		return nil, errs.NewInternalServerError()
	}

	// Swap with hash one
	payload.Password = hashPassword

	user, err = a.db.CreateUser(context.Background(), payload)
	if err != nil {
		logger.Error().Err(err).Msg("register failed")
		return nil, err
	}

	token, err := pkg.CreateAccessToken(a.cfg, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("create token failed")
		return nil, errs.NewInternalServerError()
	}

	logger.Info().
		Str("event", "login").
		Msg("registered successfully")

	return &dto.AuthResponse{
		AccessToken: token,
	}, nil
}
