package v1

import (
	"github.com/arbrix/go-test/common"
	"github.com/labstack/echo"
)

type Common struct {
	eg *echo.Group
	a  common.App
}

type TokenJSON struct {
	Token string `json:"token"`
}
