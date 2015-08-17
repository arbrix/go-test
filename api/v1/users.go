package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/arbrix/go-test/api/response"
	"github.com/arbrix/go-test/service/userService"
	"github.com/arbrix/go-test/service/userService/userPermission"
)

// @Title Users
// @Description Users's router group.
func Users(parentRoute *gin.RouterGroup) {
	route := parentRoute.Group("/users")
	route.POST("", createUser)
	route.GET("/:id", retrieveUser)
	route.GET("", retrieveUsers)
	route.PUT("/:id", userPermission.AuthRequired(updateUser))
	route.DELETE("/:id", userPermission.AuthRequired(deleteUser))
}

// @Title createUser
// @Description Create a user.
// @Accept  json
// @Param   registrationEmail        form   string     true        "User Email."
// @Param   registrationPassword        form   string  true        "User Password."
// @Success 201 {object} response.BasicResponse "User is registered successfully"
// @Failure 401 {object} response.BasicResponse "Authentication required"
// @Failure 404 {object} response.BasicResponse "User not logged in."
// @Failure 500 {object} response.BasicResponse "User is not created."
// @Resource /users
// @Router /users [post]
func createUser(c *gin.Context) {
	status, err := userService.CreateUser(c)
	messageTypes := &response.MessageTypes{
		OK:                  "registration.done",
		Unauthorized:        "login.error.fail",
		NotFound:            "registration.error.fail",
		InternalServerError: "registration.error.fail",
	}
	messages := &response.Messages{OK: "User is registered successfully."}
	response.JSON(c, status, messageTypes, messages, err)
}

// @Title retrieveUser
// @Description Retrieve a user.
// @Accept  json
// @Param   id        path    int     true        "User ID"
// @Success 200 {object} model.PublicUser "OK"
// @Failure 404 {object} response.BasicResponse "Not found"
// @Resource /users
// @Router /users/{id} [get]
func retrieveUser(c *gin.Context) {
	user, isAuthor, currentUserId, status, err := userService.RetrieveUser(c)
	if err == nil {
		c.JSON(status, gin.H{"user": user, "isAuthor": isAuthor, "currentUserId": currentUserId})
	} else {
		messageTypes := &response.MessageTypes{
			NotFound: "user.error.notFound",
		}
		response.ErrorJSON(c, status, messageTypes, err)
	}

}

// @Title retrieveUsers
// @Description Retrieve user array.
// @Accept  json
// @Success 200 {array} model.PublicUser "OK"
// @Resource /users
// @Router /users [get]
func retrieveUsers(c *gin.Context) {
	users := userService.RetrieveUsers(c)
	c.JSON(200, gin.H{"users": users})
}

// @Title updateUser
// @Description Update a user.
// @Accept  json
// @Param   id        path    int     true        "User ID"
// @Success 200 {object} model.User "OK"
// @Failure 401 {object} response.BasicResponse "Authentication required"
// @Failure 404 {object} response.BasicResponse "Not found"
// @Failure 500 {object} response.BasicResponse "User is not updated."
// @Resource /users
// @Router /users/{id} [put]
func updateUser(c *gin.Context) {
	user, status, err := userService.UpdateUser(c)
	if err == nil {
		c.JSON(status, gin.H{"user": user})
	} else {
		messageTypes := &response.MessageTypes{
			Unauthorized:        "user.error.unauthorized",
			NotFound:            "user.error.notFound",
			InternalServerError: "user.error.internalServerError",
		}
		response.ErrorJSON(c, status, messageTypes, err)
	}
}

// @Title deleteUser
// @Description Delete a user.
// @Accept  json
// @Param   id        path    int     true        "User ID"
// @Success 200 {object} response.BasicResponse
// @Failure 401 {object} response.BasicResponse "Authentication required"
// @Failure 404 {object} response.BasicResponse "Not found"
// @Failure 500 {object} response.BasicResponse "User is not deleted."
// @Resource /users
// @Router /users/{id} [delete]
func deleteUser(c *gin.Context) {
	status, err := userService.DeleteUser(c)
	if err == nil {
		c.JSON(status, response.BasicResponse{})
	} else {
		messageTypes := &response.MessageTypes{
			Unauthorized:        "user.error.unauthorized",
			NotFound:            "user.error.notFound",
			InternalServerError: "setting.leaveOurService.fail",
		}
		response.ErrorJSON(c, status, messageTypes, err)
	}
}
