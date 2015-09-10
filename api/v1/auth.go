package v1

import (
	"github.com/arbrix/go-test/app"
	"github.com/arbrix/go-test/service/user"
	"github.com/labstack/echo"
)

type Auth struct {
	Common
}

//Init Auth's route.
func NewOauth(a *app.App, e *echo.Echo) *Auth {
	au := &Auth{Common: Common{a: a, e: e}}
	au.e.Router().Add(echo.POST, "/auth", au.login, ua.e)
	return au
}

//login provide JWT in response if login success.
func login(c *echo.Context) {
	status, err := user.Login(c)
	messageTypes := &response.MessageTypes{OK: "login.done",
		Unauthorized: "login.error.passwordIncorrect",
		NotFound:     "login.error.userNotFound"}
	messages := &response.Messages{OK: "User logged in successfully."}
	if err != nil {
		response.JSON(c, status, messageTypes, messages, err)
		return
	}
	c.JSON(status, gin.H{"token": user.JwtToken})
}
