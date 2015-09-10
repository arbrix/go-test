package middleware

import (
	"github.com/labstack/echo"
)

func CheckHeader() gin.HandlerFunc {
	return func(c *echo.Context) {
		headerKey := c.Request.Header.Get("Authorization")
		if headerKey != "Bearer testkey123" {
			c.JSON(403, gin.H{"error": "authorization problem"})
			c.AbortWithStatus(403)
		}
		c.Next()
	}
}
