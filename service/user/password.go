package user

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"time"

	"github.com/arbrix/go-test/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// ResetPassword resets a password of user.
func ResetPassword(c *gin.Context) (int, error) {
	var user model.User
	var passwordResetForm PasswordResetForm
	c.BindWith(&passwordResetForm, binding.Form)
	if db.ORM.Where(&model.User{PasswordResetToken: passwordResetForm.PasswordResetToken}).First(&user).RecordNotFound() {
		return http.StatusNotFound, errors.New("User is not found.")
	}
	isExpired := timeHelper.IsExpired(user.PasswordResetUntil)
	log.Debugf("passwordResetUntil : %s", user.PasswordResetUntil.UTC())
	log.Debugf("expired : %t", isExpired)
	if isExpired {
		return http.StatusForbidden, errors.New("token not valid.")
	}
	newPassword, err := bcrypt.GenerateFromPassword([]byte(passwordResetForm.Password), 10)
	if err != nil {
		return http.StatusInternalServerError, errors.New("User is not updated. Password not Generated.")
	}
	passwordResetForm.Password = string(newPassword)
	log.Debugf("user password before : %s ", user.Password)
	modelHelper.AssignValue(&user, &passwordResetForm)
	user.PasswordResetToken = ""
	user.PasswordResetUntil = time.Now()
	log.Debugf("user password after : %s ", user.Password)
	status, err := UpdateUserCore(&user)
	if err != nil {
		return status, err
	}
	return status, err
}
