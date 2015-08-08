package service

import (
	"log"
	"strconv"
	"time"

	"github.com/arbrix/go-test/api"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type TaskResource struct {
	db gorm.DB
}

func (tr *TaskResource) CreateTask(c *gin.Context) {
	var task api.Task

	if c.Bind(&task) != nil {
		c.JSON(400, api.NewError("problem decoding body"))
		return
	}
	task.CreatedAt = time.Now()

	tr.db.Save(&task)

	c.JSON(201, task)
}

func (tr *TaskResource) GetAllTasks(c *gin.Context) {
	var tasks []api.Task

	tr.db.Order("created desc").Find(&tasks)

	c.JSON(200, tasks)
}

func (tr *TaskResource) GetTask(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("problem decoding id sent"))
		return
	}

	var task api.Task

	if tr.db.First(&task, id).RecordNotFound() {
		c.JSON(404, gin.H{"error": "not found"})
	} else {
		c.JSON(200, task)
	}
}

func (tr *TaskResource) UpdateTask(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("problem decoding id sent"))
		return
	}

	var task api.Task

	if c.Bind(&task) != nil {
		c.JSON(400, api.NewError("problem decoding body"))
		return
	}
	task.ID = int64(id)
	task.UpdatedAt = time.Now()
	if task.IsCompleted == true {
		task.CompletedAt = time.Now()
	}

	var existing api.Task

	if tr.db.First(&existing, id).RecordNotFound() {
		c.JSON(404, api.NewError("not found"))
	} else {
		tr.db.Save(&task)
		c.JSON(200, task)
	}

}

func (tr *TaskResource) DeleteTask(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("problem decoding id sent"))
		return
	}

	var task api.Task
	task.IsDeleted = true

	if tr.db.First(&task, id).RecordNotFound() {
		c.JSON(404, api.NewError("not found"))
	} else {
		tr.db.Save(&task)
		c.Data(204, "application/json", make([]byte, 0))
	}
}

func (tr *TaskResource) getId(c *gin.Context) (int64, error) {
	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return int64(id), nil
}
