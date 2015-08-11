package models

import (
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type Task struct {
	gorm.Model
	ID          int64     `json:"id" gorm:"column:id;primary_key" sql:"AUTO_INCREMENT"`
	Title       string    `json:"title" sql:"type:varchar(128)" gorm:"column:title"`
	Description string    `json:"description" sql:"type:varchar(1024)" gorm:"column:description"`
	Priority    int       `json:"priority" sql:"DEFAULT:0" gorm:"column:priority"`
	CreatedAt   time.Time `json:"created" sql:"type:datetime" gorm:"column:CreatedAt"`
	UpdatedAt   time.Time `json:"updated" sql:"type:datetime" gorm:"column:UpdatedAt"`
	CompletedAt time.Time `json:"completed" sql:"type:datetime" gorm:"column:CompletedAt"`
	IsDeleted   bool      `json:"isDeleted" gorm:"column:isDeleted"`
	IsCompleted bool      `json:"isCompeted" gorm:"column:isComplited"`
}

func (t Task) String() string {
	return "Title: " + t.Title + "; Desc: " + t.Description + "; Prior: " + strconv.Itoa(t.Priority) + "; Crt: " + t.CreatedAt.Format(time.RFC3339)
}
