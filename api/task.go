package api

import "time"

type Task struct {
	ID          int64     `json:"id",gorm:"primary_key",sql:"AUTO_INCREMENT"`
	Title       string    `json:"title" binding:"required",sql:"type:varchar(128)"`
	Description string    `json:"description",sql:"type:varchar(1024)"`
	Priority    int       `json:"priority",sql:"DEFAULT:0"`
	CreatedAt   time.Time `json:"created"`
	UpdatedAt   time.Time `json:"updated"`
	CompletedAt time.Time `json:"completed"`
	IsDeleted   bool      `json:"isDeleted"`
	IsCompleted bool      `json:"isCompeted"`
}
