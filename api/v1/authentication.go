package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/arbrix/go-test/api/response"
	"github.com/arbrix/go-test/service/userService"
)

// @Title Authentications
// @Description Authentications's router group.
func Authentications(parentRoute *gin.RouterGroup) {
	route := parentRoute.Group("/authentications")
	route.POST("", createUserAuthentication)
}

// @Title createUserAuthentication
// @Description Create a user session.
// @Accept  json
// @Param   loginEmail        form   string     true        "User email."
// @Param   loginPassword        form   string  true        "User password."
// @Success 200 {object} token string
// @Success 201 {object} response.BasicResponse "User authentication created"
// @Failure 401 {object} response.BasicResponse "Password incorrect"
// @Failure 404 {object} response.BasicResponse "User is not found"
// @Resource /authentications
// @Router /authentications [post]
func createUserAuthentication(c *gin.Context) {
	status, err := userService.CreateUserAuthentication(c)
	messageTypes := &response.MessageTypes{OK: "login.done",
		Unauthorized: "login.error.passwordIncorrect",
		NotFound:     "login.error.userNotFound"}
	messages := &response.Messages{OK: "User logged in successfully."}
	if err != nil {
		response.JSON(c, status, messageTypes, messages, err)
		return
	}
	c.JSON(status, gin.H{"token": userService.JwtToken})
}
