package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/arbrix/go-test/interfaces"
	"github.com/arbrix/go-test/model"
	"github.com/arbrix/go-test/service/oauth"
	"github.com/arbrix/go-test/util/jwt"
	"github.com/arbrix/go-test/util/middleware"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

type App struct {
	conf interfaces.Config
	db   interfaces.Orm
	fb   *oauth.Facebook
}

type ErrorMsg struct {
	Msg string `json:"msg"`
}

type UrlJSON struct {
	Url string `json:"url"`
}

type TokenJSON struct {
	Token string `json:"token"`
}

func NewApp(c interfaces.Config, db interfaces.Orm) *App {
	return &App{conf: c, db: db}
}

func (a *App) GetDB() interfaces.Orm {
	return a.db
}

func (a *App) GetConfig() interfaces.Config {
	return a.conf
}

func (app *App) Run() error {
	err := app.db.Connect(app.conf)
	if err != nil {
		return err
	}

	e := echo.New()
	env, err := app.conf.Get("env")
	if err != nil {
		return err
	}
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
func (a *App) apiRoute(e *echo.Echo) error {
	apiUrl, err := a.conf.Get("api-url")
	if err != nil {
		return err
	}
	//General API
	g := e.Group(apiUrl.(string))
	//auth
	g.Post("/auth", a.login)
	//oauth
	a.fb, err = oauth.NewFacebook(a)
	if err != nil {
		return err
	}
	fbg := g.Group("/oauth")
	fbg.Get("/facebook", a.facebookAuth)
	fbg.Get("/facebook/redirect", a.facebookRedirect)
	//tasks
	tokenizer := jwt.NewTokenizer(a)
	tg := g.Group("/tasks", tokenizer.Check())
	tg.Post("", a.create)
	tg.Get("/:id", a.retrieve)
	tg.Get("", a.retrieveAll)
	tg.Put("/:id", a.update)
	tg.Delete("/:id", a.delete)
	return nil
}

//login provide JWT in response if login success.
func (a *App) login(c *echo.Context) error {
	loginData := &model.LoginJSON{}
	err := c.Bind(loginData)
	log.Printf("login:%s; passwd:%s\n", loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}
	if loginData.Email == "" {
		err = errors.New("email could't be empty.")
		c.JSON(http.StatusNotFound, err)
		return err
	}
	if loginData.Password == "" {
		err = errors.New("password could't be empty.")
		c.JSON(http.StatusNotFound, err)
		return err
	}
	user := &model.User{}
	user.Email, user.Password = loginData.Email, loginData.Password
	status, err := user.CheckPass(a.GetDB())
	if err != nil {
		c.JSON(status, err)
		return err
	}
	tokenizer := jwt.NewTokenizer(a)
	status, err = tokenizer.Create(c, user)
	if err != nil {
		c.JSON(status, err)
		return err
	}
	return a.sendJWT(c)
}

//facebookAuth Get facebook oauth url.
func (a *App) facebookAuth(c *echo.Context) error {
	url := a.fb.URL()
	c.JSON(http.StatusOK, UrlJSON{Url: url})
	return nil
}

//facebookRedirect Redirect from Facebook oauth.
func (a *App) facebookRedirect(c *echo.Context) error {
	status, err := a.fb.Oauth(c)
	if err != nil {
		c.JSON(status, fmt.Sprintf("httpStatusCode : %d; error: %v", status, err))
		return err
	}
	return a.sendJWT(c)
}

func (a *App) sendJWT(c *echo.Context) error {
	var jwt string
	jwt, ok := c.Get("jwt").(string)
	if ok == false {
		errStr := "JWT generated but not string"
		c.JSON(http.StatusInternalServerError, ErrorMsg{Msg: errStr})
		return errors.New(errStr)
	}

	c.JSON(http.StatusAccepted, TokenJSON{Token: jwt})
	return nil
}

//create Create a task.
func (a *App) create(c *echo.Context) error {
	task := &model.Task{}
	var err error

	err = c.Bind(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMsg{Msg: err.Error()})
		return err
	}
	status, err := task.Create(a.GetDB())
	c.JSON(status, ErrorMsg{Msg: err.Error()})
	return err
}

//rertieve Retrieve a task.
func (a *App) retrieve(c *echo.Context) error {
	id := c.Param("id")
	task := &model.Task{}
	if a.GetDB().First(task, id) != nil {
		err := errors.New("Task is not found.")
		c.JSON(http.StatusNotFound, ErrorMsg{Msg: err.Error()})
		return err
	}
	c.JSON(http.StatusOK, task)
	return nil

}

//retrieve Retrieve task array.
func (a *App) retrieveAll(c *echo.Context) error {
	var tasks []*model.Task
	a.GetDB().Find(&tasks, struct{}{})
	c.JSON(http.StatusOK, tasks)
	return nil
}

//update Update a task.
func (a *App) update(c *echo.Context) error {
	id := c.Param("id")
	task := &model.Task{}
	if a.GetDB().First(task, id) != nil {
		err := errors.New("Task is not found.")
		c.JSON(http.StatusNotFound, ErrorMsg{Msg: err.Error()})
		return err
	}
	status, err := task.Update(a.GetDB(), c)
	if err == nil {
		c.JSON(status, task)
	} else {
		c.JSON(status, ErrorMsg{Msg: err.Error()})
	}
	return err
}

//delete Mark Task as Deleted.
func (a *App) delete(c *echo.Context) error {
	c.Set("deleted", true)
	return a.update(c)
}
