package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	RequestIDHeader = "X-request-ID"
	RequestIDKey    = "request_id"
)

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestId := c.Request().Header.Get(RequestIDHeader)
			if requestId == "" {
				requestId = uuid.New().String()
			}

			c.Set(RequestIDKey, requestId)
			c.Response().Header().Set(RequestIDHeader, requestId)

			return next(c)
		}
	}
}

func GetRequestId(c echo.Context) string {
	if requestId, ok := c.Get(RequestIDKey).(string); ok {
		return requestId
	}
	return ""
}
