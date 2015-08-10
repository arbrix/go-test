package service

import (
	"github.com/arbrix/go-test/api"
	"github.com/gin-gonic/gin"
)

func CheckHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerKey := c.Request.Header.Get("Authorization")
		if headerKey != "Bearer testkey123" {
			c.JSON(403, api.NewError("authorization problem"))
			c.AbortWithStatus(403)
		}
		c.Next()
	}
}
