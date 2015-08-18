package userPermission

import (
	"errors"
	"net/http"

	"github.com/arbrix/go-test/api/response"
	"github.com/arbrix/go-test/service/userService"
	"github.com/arbrix/go-test/util/log"
	"github.com/gin-gonic/gin"
)

// CurrentUserIdentical check that userId is same as current user's Id.
func CurrentUserIdentical(c *gin.Context, userId int64) (int, error) {
	currentUser, err := userService.CurrentUser(c)
	if err != nil {
		return http.StatusUnauthorized, errors.New("Auth failed.")
	}
	if currentUser.Id != userId {
		return http.StatusForbidden, errors.New("User is not identical.")
	}

	return http.StatusOK, nil
}

// AuthRequired run function when user logged in.
func AuthRequired(f func(c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := userService.CurrentUser(c)
		if err != nil {
			log.Error("Auth failed.")
			response.KnownErrorJSON(c, http.StatusUnauthorized, "error.loginPlease", errors.New("Auth failed."))
			return
		}
		f(c)
		return
	}
}
