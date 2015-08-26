package web

import (
	//"time"

	"github.com/arbrix/go-test/api"
	"github.com/arbrix/go-test/config"
	"github.com/arbrix/go-test/db"
	"github.com/arbrix/go-test/util/log"
	"github.com/arbrix/go-test/web/middleware"
	"github.com/gin-gonic/gin"
	//cors "github.com/itsjamie/gin-cors"
)

type Service struct {
}

func (s *Service) Run(cfg config.Config) {
	db.ORM = db.GormInit(cfg)
	r := gin.New()
	/*
		r.Use(cors.Middleware(cors.Config{
			Origins:         "*",
			Methods:         "GET, PUT, POST, DELETE, OPTIONS",
			RequestHeaders:  "Origin, Authorization, Content-Type",
			ExposedHeaders:  "",
			MaxAge:          50 * time.Second,
			Credentials:     false,
			ValidateHeaders: false,
		}))
	*/
	r.Use(middleware.CORSMiddleware())
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
