package web

import (
	"fmt"

	"github.com/arbrix/go-test/api"
	"github.com/arbrix/go-test/db"
	"github.com/arbrix/go-test/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/tommy351/gin-cors"
)

type Config struct {
	ListenAddress string
	DatabaseUri   string
}

func (cfg Config) String() string {
	return cfg.ListenAddress + "; " + cfg.DatabaseUri
}

type TaskService struct {
}

func (s *TaskService) Run(cfg Config) {
	db.ORM = db.GormInit(cfg)
	r := gin.New()
	r.Use(cors.Middleware(cors.Options{}))
	// Global middlewares
	// If use gin.Logger middlewares, it send duplicated request.
	switch config.Environment {
	case "DEVELOPMENT":
		r.Use(gin.Logger())
	case "TEST":
		r.Use(log.AccessLogger())
	case "PRODUCTION":
		r.Use(log.AccessLogger())
	}
	r.Use(gin.Recovery())
	//r.Use(middleware.CheckHeader())

	api.RouteAPI(r)

	r.Run(cfg.ListenAddress)
}
