package model

import "time"

type User struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email",sql:"size:255;unique"`
	Password  string    `json:"password",sql:"size:255"`
	Username  string    `json:"username",sql:"size:127;unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}
