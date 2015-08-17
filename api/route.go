// @APIVersion 1.0.0
// @Title Goyangi API
// @Description Goyangi API usually works as expected. But sometimes its not true
// @Contact api@contact.me
// @TermsOfServiceUrl http://google.com/
// @License BSD
// @LicenseUrl http://opensource.org/licenses/BSD-2-Clause
// @SubApi Authentication [/authentications]
// @SubApi User [/users]
// @SubApi Oauth [/oauth]
// @SubApi Task [/tasks]
package api

import (
	"github.com/gin-gonic/gin"

	"github.com/arbrix/go-test/api/v1"
	"github.com/arbrix/go-test/config"
)

// RouteAPI contains router groups for API
func RouteAPI(parentRoute *gin.Engine) {
	route := parentRoute.Group(config.APIURL)
	{
		v1.Users(route)
		v1.Tasks(route)
		v1.Authentications(route)
		v1.Oauth(route)
	}

}
