package v1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arbrix/go-test/common"
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

//Init Oauth's router group.
func NewOauth(a common.App, pg *echo.Group) *Oauth {
	facebook := &oauth.Facebook{}
	oa := &Oauth{Common: Common{a: a, eg: pg}, fb: facebook}
	oa.fb.Init(oa.a)
	g := oa.eg.Group("/oauth")
	g.Get("/facebook", oa.facebookAuth)
	g.Get("/facebook/redirect", oa.facebookRedirect)
	return oa
}

//facebookAuth Get facebook oauth url.
func (oa *Oauth) facebookAuth(c *echo.Context) error {
	url := oa.fb.URL()
	c.JSON(http.StatusOK, UrlJSON{Url: url})
	return nil
}

//facebookRedirect Redirect from Facebook oauth.
func (oa *Oauth) facebookRedirect(c *echo.Context) error {
	status, err := oa.fb.Oauth(c)
	if err != nil {
		log.Printf("httpStatusCode : %d; error: %v", status, err)
		c.JSON(status, fmt.Sprintf("httpStatusCode : %d; error: %v", status, err))
		return err
	}
	c.JSON(http.StatusAccepted, TokenJSON{Token: c.Get("jwt").(string)})
	return nil
}
