package model

import (
	"errors"
	"fmt"
	"github.com/arbrix/go-test/interfaces"
	"github.com/labstack/echo"
	"net/http"
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

func (t *Task) Create(db interfaces.Orm) (int, error) {
	if t.CreatedAt == nil {
		now := time.Now()
		t.CreatedAt = &now
	}
	err := db.Create(t)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, err
}

// Update updates a task.
func (t *Task) Update(db interfaces.Orm, c *echo.Context) (int, error) {
	if t.IsCompleted == true || t.IsDeleted == true {
		return http.StatusBadRequest, errors.New("This task could not be updated.")
	}
	err := c.Bind(t)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	now := time.Now()
	t.UpdatedAt = &now
	if t.IsCompleted == true {
		t.CompletedAt = &now
	}
	if c.Get("deleted") == true {
		t.IsDeleted = true
	}
	if db.Save(t) != nil {
		return http.StatusInternalServerError, errors.New("Task is not updated.")
	}
	return http.StatusOK, nil
}
