package middleware

import (
	"github.com/arbrix/go-test/models"
	"github.com/gin-gonic/gin"
)

func CheckHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerKey := c.Request.Header.Get("Authorization")
		if headerKey != "Bearer testkey123" {
			c.JSON(403, models.NewError("authorization problem"))
			c.AbortWithStatus(403)
		}
		c.Next()
	}
}
