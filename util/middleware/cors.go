package middleware

import (
	"github.com/labstack/echo"
	"net/http"
)

func CORSMiddleware() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "false")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Response().Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

			if c.Request().Method == "OPTIONS" {
				c.JSON(http.StatusNoContent, nil)
				return nil
			}

			if err := h(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}
}
