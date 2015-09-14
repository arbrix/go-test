package model

import (
	"errors"
	"fmt"
	"github.com/arbrix/go-test/interfaces"
	"github.com/arbrix/go-test/util/helper"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"regexp"
	"time"
)

//User data for manage access to API
type User struct {
	ID        int64      `json:"id" gorm:"column:id"`
	Email     string     `json:"email" sql:"size:255;unique" gorm:"column:email" binding:"required"`
	Password  string     `json:"pass" sql:"size:64" gorm:"column:pass" binding:"required"`
	Name      string     `json:"name" sql:"size:127" gorm:"column:name" binding:"required"`
	CreatedAt *time.Time `json:"created" sql:"DEFAULT:current_timestamp" gorm:"column:created"`
	UpdatedAt *time.Time `json:"updated" gorm:"column:updated"`
	DeletedAt *time.Time `json:"deleted" gorm:"column:deleted"`
}

// LoginJSON is used when creating a user authentication.
type LoginJSON struct {
	Email    string `json:"login" form:"loginEmail" binding:"required"`
	Password string `json:"pass" form:"loginPassword" binding:"required"`
}

func (u *User) String() string {
	return fmt.Sprintf("id: %d; email: %s; name: %s; crt: %d; upd: %d; del: %d", u.ID, u.Email, u.Name, u.CreatedAt, u.UpdatedAt, u.DeletedAt)
}

func (u *User) CreateNew(db interfaces.Orm) error {
	checkUser := &User{}
	if db.First(checkUser, u) == nil {
		return errors.New("User already exists")
	}
	if len(u.Password) == 0 {
		strHelper := &helper.Str{}
		u.Password = strHelper.GenRand(12)
		log.Println(u.Password)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	u.Password = string(hashPassword)
	if err != nil {
		return err
	}
	now := time.Now()
	u.CreatedAt = &now
	err = db.Create(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) CheckPass(db interfaces.Orm) (int, error) {
	//From github.com/asaskevich/govalidator
	emailPattern := "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	paswd := u.Password
	if regexp.MustCompile(emailPattern).MatchString(u.Email) {
		if db.First(u, u) != nil {
			return http.StatusNotFound, errors.New("User is not found by email.")
		}
	} else {
		u.Name, u.Email = u.Email, ""
		if db.First(u, u) != nil {
			return http.StatusNotFound, errors.New("User is not found by name.")
		}
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(paswd))
	if err != nil {
		return http.StatusUnauthorized, errors.New("Password incorrect.")
	}
	return http.StatusOK, nil
}
