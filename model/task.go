package model

import (
	"fmt"
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
	return fmt.Sprintf("id: %d; title: %s; desc: %s; pri: %d; crt: %d; upd: %d; iscomp: %b; cmp: %d", t.ID, t.Title, t.Description, t.Priority, t.CreatedAt, t.UpdatedAt, t.IsCompleted, t.CompletedAt)
}
