package task

import (
	"errors"
	"net/http"
	"time"

	"github.com/arbrix/go-test/interfaces"
	"github.com/arbrix/go-test/model"
	"github.com/labstack/echo"
)

type Service struct {
	a interfaces.App
}

func NewTaskService(a interfaces.App) *Service {
	return &Service{a: a}
}

// Create creates a task.
func (s *Service) Create(c *echo.Context) (int, error) {
	task := &model.Task{}
	var err error

	err = c.Bind(task)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if task.CreatedAt == nil {
		now := time.Now()
		task.CreatedAt = &now
	}
	err = s.a.GetDB().Create(task)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, err
}

// Retrieve retrieves a task.
func (s *Service) Retrieve(c *echo.Context) (*model.Task, int, error) {
	task := &model.Task{}
	id := c.Param("id")
	if s.a.GetDB().First(task, id) != nil {
		return task, http.StatusNotFound, errors.New("Task is not found.")
	}
	return task, http.StatusOK, nil
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
	task := &model.Task{}
	if s.a.GetDB().First(task, id) != nil {
		return task, http.StatusNotFound, errors.New("Task is not found.")
	}
	if task.IsCompleted == true || task.IsDeleted == true {
		return task, http.StatusBadRequest, errors.New("This task could not be updated.")
	}
	err := c.Bind(task)
	if err != nil {
		return task, http.StatusInternalServerError, err
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
		return task, http.StatusInternalServerError, errors.New("Task is not updated.")
	}
	return task, http.StatusOK, nil
}

// DeleteTask deletes a task.
func (s *Service) DeleteTask(c *echo.Context) (int, error) {
	id := c.Param("id")
	task := &model.Task{}
	if s.a.GetDB().First(task, id) != nil {
		return http.StatusNotFound, errors.New("Task is not found.")
	}
	if s.a.GetDB().Delete(task) != nil {
		return http.StatusInternalServerError, errors.New("Task is not deleted.")
	}
	return http.StatusOK, nil
}
