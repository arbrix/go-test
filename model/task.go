package model

import (
	"strconv"
	"time"
)

type Task struct {
	ID          int64      `json:"id" gorm:"column:id;primary_key" sql:"AUTO_INCREMENT"`
	Title       string     `json:"title" sql:"type:varchar(128)" gorm:"column:title"`
	Description string     `json:"desc" sql:"type:varchar(1024)" gorm:"column:description"`
	Priority    int        `json:"priority" sql:"DEFAULT:0" gorm:"column:priority"`
	CreatedAt   *time.Time `json:"created" sql:"DEFAULT:current_timestamp" gorm:"column:created"`
	UpdatedAt   *time.Time `json:"updated" gorm:"column:updated"`
	CompletedAt *time.Time `json:"completed" gorm:"column:completed"`
	IsDeleted   bool       `json:"isDeleted" gorm:"column:isDeleted"`
	IsCompleted bool       `json:"isCompleted" gorm:"column:isCompleted"`
}

func (t Task) String() string {
	return "Title: " + t.Title + "; Desc: " + t.Description + "; Prior: " +
		strconv.Itoa(t.Priority) + "; Crt: " + t.CreatedAt.Format(time.RFC3339) +
		"; upd: " + t.UpdatedAt.Format(time.RFC3339) +
		"; complete: " + strconv.FormatBool(t.IsCompleted) +
		"; at: " + t.CompletedAt.Format(time.RFC3339)
}
