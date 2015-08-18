package v1

import (
	"fmt"

	"github.com/arbrix/go-test/api/response"
	"github.com/arbrix/go-test/config"
	"github.com/arbrix/go-test/service/oauthService"
	"github.com/arbrix/go-test/service/userService/userPermission"
	"github.com/arbrix/go-test/util/log"
	"github.com/gin-gonic/gin"
)

// @Title Oauth
// @Description Oauth's router group.
func Oauth(parentRoute *gin.RouterGroup) {

	route := parentRoute.Group("/oauth")
	route.GET("", retrieveOauthStatus)

	route.GET("/facebook", facebookAuth)
	route.DELETE("/facebook", userPermission.AuthRequired(facebookRevoke))
	route.GET("/facebook/redirect", facebookRedirect)
}

// @Title retrieveOauthStatus
// @Description Retrieve oauth connections.
// @Accept  json
// @Success 200 {array} oauthService.oauthStatusMap "OK"
// @Failure 401 {object} response.BasicResponse "Authentication required"
// @Resource /oauth
// @Router /oauth [get]
func retrieveOauthStatus(c *gin.Context) {
	oauthStatus, status, err := oauthService.RetrieveOauthStatus(c)
	if err == nil {
		c.JSON(status, gin.H{"oauthStatus": oauthStatus})
	} else {
		messageTypes := &response.MessageTypes{
			Unauthorized: "oauth.error.unauthorized",
		}
		response.ErrorJSON(c, status, messageTypes, err)
	}
}

// @Title facebookAuth
// @Description Get facebook oauth url.
// @Accept  json
// @Success 200 {object} gin.H "{url: oauthURL}"
// @Resource /oauth
// @Router /oauth/facebook [get]
func facebookAuth(c *gin.Context) {
	url, status := oauthService.FacebookURL()
	c.JSON(status, gin.H{"url": url})
}

// @Title facebookRevoke
// @Description Get facebook oauth url.
// @Accept  json
// @Success 200 {object} gin.H "Revoked"
// @Failure 401 {object} response.BasicResponse "Authentication required"
// @Failure 404 {object} response.BasicResponse "Connection is not found"
// @Failure 500 {object} response.BasicResponse "Connection not revoked from user"
// @Resource /oauth
// @Router /oauth/facebook [delete]
func facebookRevoke(c *gin.Context) {
	oauthStatus, status, err := oauthService.RevokeFacebook(c)
	if err == nil {
		c.JSON(status, gin.H{"oauthStatus": oauthStatus})
	} else {
		messageTypes := &response.MessageTypes{
			Unauthorized:        "oauth.error.unauthorized",
			NotFound:            "oauth.error.notFound",
			InternalServerError: "oauth.error.internalServerError"}
		response.ErrorJSON(c, status, messageTypes, err)
	}
}

// @Title facebookRedirect
// @Description Redirect from Facebook oauth.
// @Accept  json
// @Success 303 {object} response.BasicResponse "Connection linked."
// @Failure 401 {object} response.BasicResponse "Authentication required"
// @Failure 404 {object} response.BasicResponse "User is not found"
// @Failure 500 {object} response.BasicResponse "Connection not linked"
// @Resource /oauth
// @Router /oauth/facebook/redirect [get]
func facebookRedirect(c *gin.Context) {
	status, err := oauthService.OauthFacebook(c)
	if err != nil {
		log.CheckErrorWithMessage(err, fmt.Sprintf("httpStatusCode : %d", status))
	}
	c.Redirect(303, config.Config.ListenAddress)
}
