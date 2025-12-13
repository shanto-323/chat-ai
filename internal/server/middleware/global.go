package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/shanto-323/chat-ai/internal/server"
	"github.com/shanto-323/chat-ai/internal/server/errs"
	"github.com/shanto-323/chat-ai/sqlerr"
)

type Global struct {
	server *server.Server
}

func NewGlobal(s *server.Server) *Global {
	return &Global{
		server: s,
	}
}

func (g *Global) CROS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: g.server.Config.Server.CORSAllowedOrigins,
	})
}

func (global *Global) RequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogError:   true,
		LogLatency: true,
		LogHost:    true,
		LogMethod:  true,
		LogURIPath: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			statusCode := v.Status

			if v.Error != nil {
				var httpErr *errs.HTTPError
				var echoErr *echo.HTTPError
				if errors.As(v.Error, &httpErr) {
					statusCode = httpErr.Status
				} else if errors.As(v.Error, &echoErr) {
					statusCode = echoErr.Code
				}
			}

			logger := GetLogger(c)

			var e *zerolog.Event

			switch {
			case statusCode >= 500:
				e = logger.Error().Err(v.Error)
			case statusCode >= 400:
				e = logger.Warn()
			default:
				e = logger.Info()
			}

			if requestID := GetRequestID(c); requestID != "" {
				e = e.Str("request_id", requestID)
			}

			if userID := GetUserID(c); userID != "" {
				e = e.Str("user_id", userID)
			}

			e.
				Dur("latency", v.Latency).
				Int("status", statusCode).
				Str("method", v.Method).
				Str("uri", v.URI).
				Str("host", v.Host).
				Str("ip", c.RealIP()).
				Str("user_agent", c.Request().UserAgent()).
				Msg("CHAT-API")

			return nil
		},
	})
}

func (global *Global) Recover() echo.MiddlewareFunc {
	return middleware.Recover()
}

func (global *Global) Secure() echo.MiddlewareFunc {
	return middleware.Secure()
}

func (global *Global) GlobalErrorHandler(err error, c echo.Context) {
	originalErr := err

	var httpErr *errs.HTTPError
	if !errors.As(err, &httpErr) {
		var echoErr *echo.HTTPError
		if errors.As(err, &echoErr) {
			if echoErr.Code == http.StatusNotFound {
				err = errs.NewNotFoundError("Route not found", false, nil)
			}
		} else {
			err = sqlerr.HandleError(err)
		}
	}

	var echoErr *echo.HTTPError
	var status int
	var code string
	var message string
	var fieldErrors []errs.FieldError
	var action *errs.Action

	switch {
	case errors.As(err, &httpErr):
		status = httpErr.Status
		code = httpErr.Code
		message = httpErr.Message
		fieldErrors = httpErr.Errors
		action = httpErr.Action

	case errors.As(err, &echoErr):
		status = echoErr.Code
		code = errs.MakeUpperCaseWithUnderscores(http.StatusText(status))
		if msg, ok := echoErr.Message.(string); ok {
			message = msg
		} else {
			message = http.StatusText(echoErr.Code)
		}

	default:
		status = http.StatusInternalServerError
		code = errs.MakeUpperCaseWithUnderscores(
			http.StatusText(http.StatusInternalServerError))
		message = http.StatusText(http.StatusInternalServerError)
	}

	logger := *GetLogger(c)

	logger.Error().Stack().
		Err(originalErr).
		Int("status", status).
		Str("error_code", code).
		Msg(message)

	if !c.Response().Committed {
		_ = c.JSON(status, errs.HTTPError{
			Code:     code,
			Message:  message,
			Status:   status,
			Override: httpErr != nil && httpErr.Override,
			Errors:   fieldErrors,
			Action:   action,
		})
	}
}
