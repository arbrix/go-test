package v1

import (
	"github.com/labstack/echo"
	"net/http"

	"github.com/arbrix/go-test/common"
	"github.com/arbrix/go-test/service/task"
	"github.com/arbrix/go-test/util/jwt"
)

type Task struct {
	Common
	ts *task.Service
}

//Init Task's route.
func NewTask(a common.App, pg *echo.Group) *Task {
	t := &Task{Common: Common{a: a, eg: pg}, ts: task.NewTaskService(a)}
	tokenizer := jwt.NewTokenizer(a)
	g := t.eg.Group("/tasks", tokenizer.Check())
	g.Post("", t.create)
	g.Get("/:id", t.retrieve)
	g.Get("", t.retrieveAll)
	g.Put("/:id", t.update)
	g.Delete("/:id", t.delete)
	return t
}

//create Create a task.
func (t *Task) create(c *echo.Context) error {
	status, err := t.ts.Create(c)
	c.JSON(status, err)
	return err
}

//rertieve Retrieve a task.
func (t *Task) retrieve(c *echo.Context) error {
	task, status, err := t.ts.Retrieve(c)
	if err == nil {
		c.JSON(status, task)
	} else {
		c.JSON(status, err)
	}
	return err

}

//retrieve Retrieve task array.
func (t *Task) retrieveAll(c *echo.Context) error {
	tasks := t.ts.RetrieveAll(c)
	c.JSON(http.StatusOK, tasks)
	return nil
}

//update Update a task.
func (t *Task) update(c *echo.Context) error {
	task, status, err := t.ts.Update(c)
	if err == nil {
		c.JSON(status, task)
	} else {
		c.JSON(status, err)
	}
	return err
}

//delete Mark Task as Deleted.
func (t *Task) delete(c *echo.Context) error {
	c.Set("deleted", true)
	return t.update(c)
}
