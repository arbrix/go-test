package app

import (
	"flag"
	"log"
	"net/http"

	"github.com/arbrix/go-test/api/v1"
	"github.com/arbrix/go-test/common"
	"github.com/arbrix/go-test/util/middleware"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

type App struct {
	conf common.Config
	db   common.Orm
}

type ErrorMsg struct {
	Msg string `json:"msg"`
}

func NewApp(c common.Config, db common.Orm) *App {
	return &App{conf: c, db: db}
}

func (a *App) GetDB() common.Orm {
	return a.db
}

func (a *App) GetConfig() common.Config {
	return a.conf
}

func (app *App) Run() error {
	var env string
	flag.StringVar(&env, "env", "dev", "define environment: dev, prod, test (place config file *.json with the same name in ./config folder)")
	flag.Parse()

	err := app.conf.Load(env)
	if err != nil {
		return err
	}
	err = app.db.Connect(app.conf)

	e := echo.New()
	if env == "dev" {
		e.Debug()
	}
	e.SetHTTPErrorHandler(func(err error, c *echo.Context) {
		var code int
		var msg string
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code()
			msg = he.Error()
		} else {
			code = http.StatusInternalServerError
			msg = err.Error()
		}
		method := c.Request().Method
		path := c.Request().URL.Path
		if path == "" {
			path = "/"
		}
		if c.Response().Size() == 0 {
			c.JSON(code, &ErrorMsg{Msg: msg})
		}
		size := c.Response().Size()
		log.Println(method, path, code, msg, size)
	})
	e.Use(middleware.CORSMiddleware())
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	err = app.apiRoute(e)
	if err != nil {
		return err
	}

	addr, err := app.conf.Get("ListenAddress")
	if err != nil {
		return err
	}
	e.Run(addr.(string))
	return nil
}

// apiRoute contains router groups for API
func (app *App) apiRoute(e *echo.Echo) error {
	apiUrl, err := app.conf.Get("api-url")
	if err != nil {
		return err
	}
	g := e.Group(apiUrl.(string))
	v1.NewAuth(app, g)
	v1.NewOauth(app, g)
	v1.NewTask(app, g)
	return nil
}
