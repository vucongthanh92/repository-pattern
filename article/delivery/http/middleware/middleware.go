package middleware

import (
	"github.com/labstack/echo"
)

// ArticleMiddleware struct
type ArticleMiddleware struct{}

// CORS method
func (m *ArticleMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

// InitMiddleware func run first
func InitMiddleware() *ArticleMiddleware {
	return &ArticleMiddleware{}
}
