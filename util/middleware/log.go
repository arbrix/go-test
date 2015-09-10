package middleware

import (
	"github.com/labstack/echo"
	"log"
	"os"
	"time"
)

func AccessLogger() gin.HandlerFunc {
	f, err := os.OpenFile("../log/access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	stdlogger := log.New(f, "", 0)

	return func(c *echo.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.Request.RemoteAddr
		method := c.Request.Method
		statusCode := c.Writer.Status()

		stdlogger.Printf("[GIN] %v | %3d | %12v | %s | %s %-7s %s\n%s",
			end.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			c.Request.URL.Path,
			c.Errors.String(),
		)
	}
}
