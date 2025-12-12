package handler

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/chat-ai/internal/server/middleware"
	"github.com/shanto-323/chat-ai/internal/server/validation"
)

type HandlerFunc[Req validation.Validatable, Res any] func(c echo.Context, req Req) (Res, error)

type HandleNoResponseFunc[Req validation.Validatable] func(c echo.Context, req Req) error

type ResponseHandler interface {
	Handle(c echo.Context, result any) error
	GetOperation() string
}

type JSONResponseHandler struct {
	status int
}

func (h JSONResponseHandler) Handle(c echo.Context, result any) error {
	return c.JSON(h.status, result)
}

func (h JSONResponseHandler) GetOperation() string {
	return "handler"
}

type NoResponseHandler struct {
	status int
}

func (h NoResponseHandler) Handle(c echo.Context, result any) error {
	return c.NoContent(h.status)
}

func (h NoResponseHandler) GetOperation() string {
	return "handler_no_response"
}

func handleRequest[Req validation.Validatable](
	c echo.Context,
	req Req,
	handler func(c echo.Context, req Req) (any, error),
	responseHandler ResponseHandler,
) error {
	start := time.Now()
	method := c.Request().Method
	path := c.Path()

	// Get context-enhanced logger
	loggerBuilder := middleware.GetLogger(c).With().
		Str("operation", responseHandler.GetOperation()).
		Str("method", method).
		Str("path", path)

	logger := loggerBuilder.Logger()

	validationStart := time.Now()
	if err := validation.BindAndValidate(c, req); err != nil {
		validationDuration := time.Since(validationStart)

		logger.Error().
			Err(err).
			Dur("validation_duration", validationDuration).
			Msg("request validation failed")

		return err
	}

	validationDuration := time.Since(validationStart)

	logger.Debug().
		Dur("validation_duration", validationDuration).
		Msg("request validation successful")

	handlerStart := time.Now()
	result, err := handler(c, req)
	handlerDuration := time.Since(handlerStart)

	if err != nil {
		totalDuration := time.Since(start)

		logger.Error().
			Err(err).
			Dur("handler_duration", handlerDuration).
			Dur("total_duration", totalDuration).
			Msg("handler execution failed")

		return err
	}

	totalDuration := time.Since(start)

	logger.Info().
		Dur("handler_duration", handlerDuration).
		Dur("validation_duration", validationDuration).
		Dur("total_duration", totalDuration).
		Msg("request completed successfully")

	return responseHandler.Handle(c, result)
}

func Handle[Req validation.Validatable, Res any](
	handler HandlerFunc[Req, Res],
	status int,
	req Req,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handleRequest(c, req, func(c echo.Context, req Req) (any, error) {
			return handler(c, req)
		}, JSONResponseHandler{status: status})
	}
}

func HandleNoResponse[Req validation.Validatable](
	handler HandleNoResponseFunc[Req],
	status int,
	req Req,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handleRequest(c, req, func(c echo.Context, req Req) (any, error) {
			return nil, handler(c, req)
		}, NoResponseHandler{status: status})
	}
}
