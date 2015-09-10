package task

import (
	"errors"
	"net/http"
	"time"

	"github.com/arbrix/go-test/app"
	"github.com/arbrix/go-test/model"
	"github.com/labstack/echo"
)

var taskFields = []string{"title", "description", "priority", "createdAt", "updatedAt"}

type Service struct {
}

// Create creates a task.
func (s *Service) Create(c *echo.Context, a *app.App) (int, error) {
	var task model.Task
	var err error

	err = c.Bind(&task)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = a.GetDB().Create(&task)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, err
}

// Retrieve retrieves a task.
func (s *Service) Retrieve(c *echo.Context, a *app.App) (*model.Task, int, error) {
	var task model.Task
	id := c.Param("id")
	if a.GetDB().First(&task, id) != nil {
		return &task, http.StatusNotFound, errors.New("Task is not found.")
	}
	return &task, http.StatusOK, nil
}

// RetrieveAll retrieves tasks.
func (s *Service) RetrieveAll(c *echo.Context, a *app.App) []*model.Task {
	var tasks []*model.Task
	a.GetDB().Find(&tasks, struct{}{})
	return tasks
}

// Update updates a task.
func (s *Service) Update(c *echo.Context, a *app.App) (*model.Task, int, error) {
	id := c.Param("id")
	var task model.Task
	if a.GetDB().First(&task, id) != nil {
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
	if a.GetDB().Save(task) != nil {
		return &task, http.StatusInternalServerError, errors.New("Task is not updated.")
	}
	return &task, http.StatusOK, nil
}

// DeleteTask deletes a task.
func (s *Service) DeleteTask(c *echo.Context, a *app.App) (int, error) {
	id := c.Param("id")
	var task model.Task
	if a.GetDB().First(&task, id) != nil {
		return http.StatusNotFound, errors.New("Task is not found.")
	}
	if a.GetDB().Delete(&task) != nil {
		return http.StatusInternalServerError, errors.New("Task is not deleted.")
	}
	return http.StatusOK, nil
}
