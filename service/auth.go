package service

import (
	"fmt"

	"github.com/arbrix/go-test/api"
	"github.com/gin-gonic/gin"
)

func CheckHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		if headerKey := c.Request.Header.Get("Authorization: Bearer testkey123"); headerKey == "" {
			fmt.Println(headerKey)
			c.JSON(403, api.NewError("authorization problem"))
			c.Abort()
		}
		c.Next()
	}
}
