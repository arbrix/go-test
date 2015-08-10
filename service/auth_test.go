package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type requestOptions struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    io.Reader
}

func request(server *gin.Engine, options requestOptions) *httptest.ResponseRecorder {
	if options.Method == "" {
		options.Method = "GET"
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(options.Method, options.URL, options.Body)

	if options.Headers != nil {
		for key, value := range options.Headers {
			req.Header.Set(key, value)
		}
	}

	server.ServeHTTP(w, req)

	if err != nil {
		panic(err)
	}

	return w
}

func newServer() *gin.Engine {
	g := gin.New()
	g.Use(CheckHeader())

	return g
}

func TestWithoutAthorizationKey(t *testing.T) {
	g := newServer()
	assert := assert.New(t)

	g.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r := request(g, requestOptions{
		URL: "/test",
		Headers: map[string]string{
			"Authorization": "",
		},
	})

	assert.Equal("", r.Header().Get("Authorization"))
	assert.Equal("{\"error\":\"authorization problem\"}\n", r.Body.String())
}

func TestWithAthorizationKey(t *testing.T) {
	g := newServer()
	assert := assert.New(t)

	g.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r := request(g, requestOptions{
		URL: "/test",
		Headers: map[string]string{
			"Authorization": "Bearer testkey123",
		},
	})

	assert.Equal("", r.Header().Get("Authorization"))
	assert.Equal("OK", r.Body.String())
}
