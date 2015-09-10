package v1

import (
	"github.com/arbrix/go-test/app"
	"github.com/labstack/echo"
)

type Common struct {
	e *echo.Echo
	a *app.App
}
