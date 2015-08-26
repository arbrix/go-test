package taskService

import (
	"errors"
	"net/http"
	"time"

	"github.com/arbrix/go-test/db"
	"github.com/arbrix/go-test/model"
	"github.com/arbrix/go-test/util/log"
	"github.com/arbrix/go-test/util/modelHelper"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var taskFields = []string{"title", "description", "priority", "createdAt", "updatedAt"}

// CreateTaskFromForm creates a user from a registration form.
func CreateTaskFromForm(creationForm CreateTaskForm) (model.Task, error) {
	var task model.Task
	log.Debugf("createTaskForm %+v\n", creationForm)
	modelHelper.AssignValue(&task, &creationForm)
	log.Debugf("task %+v\n", task)
	if db.ORM.Create(&task).Error != nil {
		return task, errors.New("Task is not created.")
	}
	return task, nil
}

// CreateTask creates a task.
func CreateTask(c *gin.Context) (int, error) {
	var createForm CreateTaskForm
	var err error

	err = c.BindWith(&createForm, binding.JSON)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	_, err = CreateTaskFromForm(createForm)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, err
}

// RetrieveTask retrieves a task.
func RetrieveTask(c *gin.Context) (*model.Task, int, error) {
	var task model.Task
	id := c.Params.ByName("id")
	if db.ORM.First(&task, id).RecordNotFound() {
		return &task, http.StatusNotFound, errors.New("Task is not found.")
	}
	return &task, http.StatusOK, nil
}

// RetrieveTasks retrieves tasks.
func RetrieveTasks(c *gin.Context) []*model.Task {
	var tasks []*model.Task
	db.ORM.Find(&tasks)
	return tasks
}

// UpdateTaskCore updates a task. (Applying the modifed data of task).
func UpdateTaskCore(task *model.Task) (int, error) {
	if db.ORM.Save(task).Error != nil {
		return http.StatusInternalServerError, errors.New("Task is not updated.")
	}
	return http.StatusOK, nil
}

// UpdateTask updates a task.
func UpdateTask(c *gin.Context) (*model.Task, int, error) {
	id := c.Params.ByName("id")
	var task model.Task
	if db.ORM.First(&task, id).RecordNotFound() {
		return &task, http.StatusNotFound, errors.New("Task is not found.")
	}
	var form UpdateTaskForm
	err := c.BindWith(&form, binding.JSON)
	if err != nil {
		return &task, http.StatusInternalServerError, err
	}
	log.Debugf("form %+v\n", form)
	modelHelper.AssignValue(&task, &form)
	task.UpdatedAt = time.Now()
	if task.IsCompleted == true {
		task.CompletedAt = time.Now()
	}

	log.Debugf("params %+v\n", c.Params)
	status, err := UpdateTaskCore(&task)
	return &task, status, err
}

func MarkAsDeleted(c *gin.Context) (*model.Task, int, error) {
	id := c.Params.ByName("id")
	var task model.Task
	if db.ORM.First(&task, id).RecordNotFound() {
		return &task, http.StatusNotFound, errors.New("Task is not found.")
	}
	task.UpdatedAt = time.Now()
	task.IsDeleted = true
	status, err := UpdateTaskCore(&task)
	return &task, status, err
}

// DeleteTask deletes a task.
func DeleteTask(c *gin.Context) (int, error) {
	id := c.Params.ByName("id")
	var task model.Task
	if db.ORM.First(&task, id).RecordNotFound() {
		return http.StatusNotFound, errors.New("Task is not found.")
	}
	if db.ORM.Delete(&task).Error != nil {
		return http.StatusInternalServerError, errors.New("Task is not deleted.")
	}
	return http.StatusOK, nil
}
