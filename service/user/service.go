package user

import (
	"errors"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"

	"github.com/arbrix/go-test/common"
	"github.com/arbrix/go-test/model"
	"github.com/arbrix/go-test/util/helper"
	"github.com/labstack/echo"
)

type Service struct {
	a common.App
}

func NewUserService(a common.App) *Service {
	us := &Service{a: a}
	return us
}

// Create creates a user.
func (s *Service) Create(c *echo.Context) (*model.User, error) {
	var user *model.User
	var err error

	err = c.Bind(user)
	if err != nil {
		return nil, err
	}
	return s.AddNew(user)
}

func (s *Service) AddNew(u *model.User) (*model.User, error) {
	var checkUser *model.User
	if s.a.GetDB().First(checkUser, u) == nil {
		return checkUser, errors.New("User already exists")
	}
	if len(u.Password) == 0 {
		strHelper := &helper.Str{}
		u.Password = strHelper.GenRand(12)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	u.Password = string(hashPassword)
	if err != nil {
		return nil, err
	}
	err = s.a.GetDB().Create(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) Login(login, paswd string) (*model.User, int, error) {
	if login == "" {
		return nil, http.StatusNotFound, errors.New("email could't be empty.")
	}
	if paswd == "" {
		return nil, http.StatusNotFound, errors.New("password could't be empty.")
	}
	user := &model.User{}
	//From github.com/asaskevich/govalidator
	emailPattern := "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	if regexp.MustCompile(emailPattern).MatchString(login) {
		user.Email = login
		if s.a.GetDB().First(user, user) != nil {
			return nil, http.StatusNotFound, errors.New("User is not found by email.")
		}
	} else {
		user.Name = login
		if s.a.GetDB().First(user, user) != nil {
			return nil, http.StatusNotFound, errors.New("User is not found by name.")
		}
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(paswd))
	if err != nil {
		return nil, http.StatusUnauthorized, errors.New("Password incorrect.")
	}
	return user, http.StatusOK, nil
}
