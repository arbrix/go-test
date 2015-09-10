package v1

import (
	"github.com/arbrix/go-test/common"
	"github.com/arbrix/go-test/model"
	"github.com/arbrix/go-test/service/user"
	"github.com/arbrix/go-test/util/jwt"
	"github.com/labstack/echo"
	"net/http"
)

type Auth struct {
	Common
}

//Init Auth's route.
func NewAuth(a common.App, pg *echo.Group) *Auth {
	au := &Auth{Common: Common{a: a, eg: pg}}
	au.eg.Post("/auth", au.login)
	return au
}

//login provide JWT in response if login success.
func (au *Auth) login(c *echo.Context) error {
	var loginData model.LoginJSON
	err := c.Bind(loginData)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}
	us := user.NewUserService(au.a)
	user, status, err := us.Login(loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(status, err)
		return err
	}
	tokenizer := jwt.NewTokenizer(au.a)
	status, err = tokenizer.Create(c, user)
	if err != nil {
		c.JSON(status, err)
		return err
	}
	c.JSON(status, TokenJSON{Token: c.Get("jwt").(string)})
	return nil
}
