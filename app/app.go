package app

import (
	//"time"

	//"github.com/arbrix/go-test/api/v1"
	//"github.com/arbrix/go-test/util/middleware"
	"github.com/gin-gonic/gin"
)

type App struct {
	conf Config
	db   Orm
}

func (app *App) Run() {
	r := gin.New()
	// r.Use(middleware.CORSMiddleware())
	// r.Use(middleware.AccessLogger())
	r.Use(gin.Recovery())

	// app.apiRoute(r)

	if addr, err := app.conf.Get("ListenAddress"); err == nil {
		r.Run(addr.(string))
	}
}

// apiRoute contains router groups for API
// func (app *App) apiRoute(rootRoute gin.Engine) {
// 	route := parentRoute.Group(app.conf.Get("API-URL"))
// 	{
// 		// v1.Users(route)
// 		// v1.Tasks(route)
// 		// v1.Authentications(route)
// 		// v1.Oauth(route)
// 	}
// }
