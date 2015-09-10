package task

import (
	"errors"
	"net/http"
	"time"

	"github.com/arbrix/go-test/common"
	"github.com/arbrix/go-test/model"
	"github.com/labstack/echo"
)

type Service struct {
	a common.App
}

func NewTaskService(a common.App) *Service {
	return &Service{a: a}
}

// Create creates a task.
func (s *Service) Create(c *echo.Context) (int, error) {
	var task model.Task
	var err error

	err = c.Bind(&task)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = s.a.GetDB().Create(&task)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, err
}

// Retrieve retrieves a task.
func (s *Service) Retrieve(c *echo.Context) (*model.Task, int, error) {
	var task model.Task
	id := c.Param("id")
	if s.a.GetDB().First(&task, id) != nil {
		return &task, http.StatusNotFound, errors.New("Task is not found.")
	}
	return &task, http.StatusOK, nil
}

// RetrieveAll retrieves tasks.
func (s *Service) RetrieveAll(c *echo.Context) []*model.Task {
	var tasks []*model.Task
	s.a.GetDB().Find(&tasks, struct{}{})
	return tasks
}

// Update updates a task.
func (s *Service) Update(c *echo.Context) (*model.Task, int, error) {
	id := c.Param("id")
	var task model.Task
	if s.a.GetDB().First(&task, id) != nil {
		return &task, http.StatusNotFound, errors.New("Task is not found.")
	}
	err := c.Bind(&task)
	if err != nil {
		return &task, http.StatusInternalServerError, err
	}
	now := time.Now()
	task.UpdatedAt = &now
	if task.IsCompleted == true {
		task.CompletedAt = &now
	}
	if c.Get("deleted") == true {
		task.IsDeleted = true
	}
	if s.a.GetDB().Save(task) != nil {
		return &task, http.StatusInternalServerError, errors.New("Task is not updated.")
	}
	return &task, http.StatusOK, nil
}

// DeleteTask deletes a task.
func (s *Service) DeleteTask(c *echo.Context) (int, error) {
	id := c.Param("id")
	var task model.Task
	if s.a.GetDB().First(&task, id) != nil {
		return http.StatusNotFound, errors.New("Task is not found.")
	}
	if s.a.GetDB().Delete(&task) != nil {
		return http.StatusInternalServerError, errors.New("Task is not deleted.")
	}
	return http.StatusOK, nil
}
