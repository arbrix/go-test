package api

import (
	"strconv"
	"time"
)

type Task struct {
	ID          int64     `json:"id" gorm:"primary_key" sql:"AUTO_INCREMENT"`
	Title       string    `json:"title" sql:"type:varchar(128)"`
	Description string    `json:"description" sql:"type:varchar(1024)"`
	Priority    int       `json:"priority" sql:"DEFAULT:0"`
	CreatedAt   time.Time `json:"created" sql:"type:datetime"`
	UpdatedAt   time.Time `json:"updated" sql:"type:datetime"`
	CompletedAt time.Time `json:"completed" sql:"type:datetime"`
	IsDeleted   bool      `json:"isDeleted"`
	IsCompleted bool      `json:"isCompeted"`
}

func (t Task) String() string {
	return "Title: " + t.Title + "; Desc: " + t.Description + "; Prior: " + strconv.Itoa(t.Priority) + "; Crt: " + t.CreatedAt.String()
}
