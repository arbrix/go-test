package model

import (
	"strconv"
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

func (u User) String() string {
	return "id: " + strconv.FormatInt(u.ID, 10) + "; email: " + u.Email + "; Name: " + u.Name + "; Crt: " + u.CreatedAt.Format(time.RFC3339) + "; upd: " + u.UpdatedAt.Format(time.RFC3339) + "; del: " + u.DeletedAt.Format(time.RFC3339)
}
