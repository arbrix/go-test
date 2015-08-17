package userService

import (
	"errors"
	"net/http"
	"time"

	"github.com/arbrix/go-test/db"
	"github.com/arbrix/go-test/model"
	"github.com/arbrix/go-test/util/log"
	"github.com/arbrix/go-test/util/timeHelper"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// EmailVerification verifies an email of user.
func EmailVerification(c *gin.Context) (int, error) {
	var user model.User
	var verifyEmailForm VerifyEmailForm
	c.BindWith(&verifyEmailForm, binding.Form)
	log.Debugf("verifyEmailForm.ActivationToken : %s", verifyEmailForm.ActivationToken)
	if db.ORM.Where(&model.User{ActivationToken: verifyEmailForm.ActivationToken}).First(&user).RecordNotFound() {
		return http.StatusNotFound, errors.New("User is not found.")
	}
	isExpired := timeHelper.IsExpired(user.ActivateUntil)
	log.Debugf("passwordResetUntil : %s", user.ActivateUntil.UTC())
	log.Debugf("expired : %t", isExpired)
	if isExpired {
		return http.StatusForbidden, errors.New("token not valid.")
	}
	user.ActivationToken = ""
	user.ActivateUntil = time.Now()
	user.ActivatedAt = time.Now()
	user.Activation = true
	status, err := UpdateUserCore(&user)
	if err != nil {
		return status, err
	}
	status, err = SetCookie(c, user.Token)
	return status, err
}
