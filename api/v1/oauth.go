package v1

import (
	"fmt"
	"net/http"

	"github.com/arbrix/go-test/app"
	"github.com/arbrix/go-test/service/oauth"
	"github.com/labstack/echo"
)

type Oauth struct {
	Common
	fb *oauth.Facebook
}

type UrlJSON struct {
	Url string `json:"url"`
}

type TokenJSON struct {
	Token string `json:"token"`
}

//Init Oauth's router group.
func NewOauth(a *app.App, e *echo.Echo) *Oauth {
	facebook := &oauth.Facebook{}
	oa := &Oauth{Common: Common{a: a, e: e}, fb: facebook}
	oa.fb.Init(oa.a)
	g := oa.e.Group("/oauth")
	g.Get("/facebook", oa.facebookAuth)
	g.Get("/facebook/redirect", oa.facebookRedirect)
	return oa
}

//facebookAuth Get facebook oauth url.
func (oa *Oauth) facebookAuth(c *echo.Context) {
	url := oa.fb.URL()
	c.JSON(http.StatusOK, UrlJSON{Url: url})
}

//facebookRedirect Redirect from Facebook oauth.
func (oa *Oauth) facebookRedirect(c *echo.Context) {
	status, err := oa.fb.Oauth(c)
	if err != nil {
		c.JSON(status, fmt.Sprintf("httpStatusCode : %d; error: %v", status, err))
	}
	c.JSON(http.StatusAccepted, TokenJSON{Token: c.Get("jwt").(string)})
}
