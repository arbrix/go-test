package service

import "github.com/arbrix/go-test/api"

func CheckHeader() HandlerFunc {
	return func(c *Context) {
		if headerKey := c.Request.Header.Get("Authorization: Bearer testkey123"); headerKey != nil {
			c.JSON(403, api.NewError("authorization problem"))
		}
	}
}
