package userService

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"

	"github.com/arbrix/go-test/db"
	"github.com/arbrix/go-test/model"
	"github.com/arbrix/go-test/util/log"
	"github.com/arbrix/go-test/util/validation"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// CreateUserAuthentication creates user authentication.
func CreateUserAuthentication(c *gin.Context) (int, error) {
	var form LoginForm
	err := c.BindWith(&form, binding.JSON)
	if err != nil {
		return http.StatusNotFound, err
	}
	email := form.Email
	pass := form.Password
	status, err := Login(email, pass)
	if err != nil {
		return status, err
	}
	status, err = CreateToken(email)
	return status, err
}

func Login(email string, pass string) (int, error) {
	if email == "" {
		return http.StatusNotFound, errors.New("email could't be empty.")
	}
	if pass == "" {
		return http.StatusNotFound, errors.New("password could't be empty.")
	}
	log.Debugf("User email : %s , password : %s", email, pass)
	fmt.Printf("User email : %s , password : %s", email, pass)
	var user model.User
	isValidEmail := validation.EmailValidation(email)
	if isValidEmail {
		log.Debug("User entered valid email.")
		if db.ORM.Where("email = ?", email).First(&user).RecordNotFound() {
			return http.StatusNotFound, errors.New("User is not found by email.")
		}
	} else {
		log.Debug("User entered username.")
		if db.ORM.Where("username = ?", email).First(&user).RecordNotFound() {
			return http.StatusNotFound, errors.New("User is not found by username.")
		}
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	if err != nil {
		return http.StatusUnauthorized, errors.New("Password incorrect.")
	}
	return http.StatusOK, nil
}
