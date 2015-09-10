package app

import (
	//"time"
	"flag"

	//"github.com/arbrix/go-test/api/v1"
	//"github.com/arbrix/go-test/util/middleware"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

type App struct {
	conf Config
	db   Orm
}

func (a *App) GetDB() Orm {
	if a.db.IsConnected() == false {
		a.db.Connect(a.conf)
	}
	return a.db
}

func (a *App) GetConfig() Config {
	return a.conf
}

func (app *App) Run() {
	var env string
	flag.StringVar(&env, "env", "dev", "define environment: dev, prod, test (place config file *.json with the same name in ./config folder)")
	flag.Parse()

	app.conf.Load(env)

	r := echo.New()
	// r.Use(middleware.CORSMiddleware())
	// r.Use(middleware.AccessLogger())
	r.Use(mw.Recover())

	// app.apiRoute(r)

	if addr, err := app.conf.Get("ListenAddress"); err == nil {
		r.Run(addr.(string))
	}
}

func (a *App) TestInit(db Orm, cnf Config) {
	a.db = db
	a.conf = cnf
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
